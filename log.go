package gathering

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

// ErrNotFound is the error returned when a log item is not found
var ErrNotFound = errors.New("not found")
var segmentStartRegex = regexp.MustCompile(`\[UnityCrossThreadLogger\].*|\[Client GRE\]`)
var clientGRE = []byte("Client GRE")
var newLineDelimiter = []byte("\n")
var findDate = regexp.MustCompile(`(?m)\d+\/\d+\/\d{4}\s\d+:\d+:\d+\s[APM]{2}`)
var dateLayout = "1/2/2006 15:04:05 PM"

// Log is the well-structured format of the output_log.txt, parsed into Segments
type Log struct {
	Segments []*Segment
}

// ParseLog returns a log file parsed into Segments
func ParseLog(f *os.File) (*Log, error) {
	reader := bufio.NewReader(f)
	var segments []*Segment
	previous := &Segment{
		Range: []int{0},
	}
	// This won't succeed with lines longer than 4096 bytes, but a brief
	// look through the log files show that lines generally aren't longer than
	// 1000 characters, and those that _are_ we can ignore.
	// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
	i := 0
	var buffer bytes.Buffer
	for {
		i++
		b, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err == bufio.ErrBufferFull {
			continue
		}
		if err != nil {
			log.Printf("unexpected error reading log file: %v\n", err.Error())
			continue
		}
		if segmentStartRegex.Match(b) {
			t := UnityLogger
			if bytes.Contains(b, clientGRE) {
				t = ClientGRE
			}
			// buffer.WriteTo(previous.Text)
			previous.Text = make([]byte, buffer.Len())
			copy(previous.Text, buffer.Bytes())
			// previous.Text = buffer.Bytes()
			previous.Range = append(previous.Range, i)
			previous.SegmentType = parseType(previous.Text)
			segments = append(segments, previous)
			buffer.Reset()
			previous = &Segment{
				LoggerType: t,
				Time:       parseDate(b),
				Line:       b,
				Range:      []int{i},
			}
		} else {
			buffer.Write(b)
		}
	}
	previous.Text = buffer.Bytes()
	previous.Range = append(previous.Range, i)
	return &Log{
		Segments: segments,
	}, nil
}

// Collection finds a collection
func (l *Log) Collection() (map[string]int, error) {
	// TODO: Put these all in the same for loop
	for i := len(l.Segments) - 1; i >= 0; i-- {
		s := l.Segments[i]
		if s.IsCollection() {
			return s.ParseCollection()
		}
	}
	return nil, ErrNotFound
}

// Rank finds the rank information
// The game doesn't ask for the entire rank info often, so we
// go through the log and update the parsed rank with changes
// so we return the most up to date version
func (l *Log) Rank() (rank *ArenaRankInfo, err error) {
	for _, s := range l.Segments {
		if s.IsRankInfo() {
			rank, err = s.ParseRankInfo()
			if err != nil {
				return
			}
		}
		if s.IsRankUpdated() && rank != nil {
			updated, err := s.ParseRankUpdated()
			if err != nil {
				continue
			}
			rank.Update(updated)
		}
	}
	return
}

// Inventory finds the player inventory information
func (l *Log) Inventory() (*ArenaPlayerInventory, error) {
	// TODO: Put in same loop
	for i := len(l.Segments) - 1; i >= 0; i-- {
		s := l.Segments[i]
		if s.IsPlayerInventory() {
			return s.ParsePlayerInventory()
		}
	}
	return nil, ErrNotFound
}

// Auth finds the player's ingame name
func (l *Log) Auth() ([]byte, error) {
	// TODO: Put in same loop
	for _, s := range l.Segments {
		if s.IsPlayerAuth() {
			return s.ParseAuth()
		}
	}
	return nil, ErrNotFound
}

// Decks finds the player decks
func (l *Log) Decks() ([]ArenaDeck, error) {
	// TODO: Put in same loop
	for i := len(l.Segments) - 1; i >= 0; i-- {
		s := l.Segments[i]
		if s.IsArenaDecks() {
			return s.ParseArenaDecks()
		}
	}
	return nil, ErrNotFound
}

// Boosters finds all the opened boosters
func (l *Log) Boosters() ([]*Booster, error) {
	// TODO: Put in same loop
	boosters := make([]*Booster, 0)
	for _, s := range l.Segments {
		if s.IsCrackBooster() {
			b, err := s.ParseCrackBooster()
			if err == nil {
				boosters = append(boosters, b)
			}
		}
	}
	return boosters, nil
}

// Matches finds the player matches
// This is a little more involved, since we really need 3 pieces of information
// and they are not together.
// First, we look for the start of a match. Once we find that, we look backward
// until we figure out what deck they are using. Finally, we look forward again
// to find the result of the match.
// The result may not be known when we start parsing, so those values are all
// optional. The server only needs the MatchID to tie together the data.
func (l *Log) Matches() ([]*ArenaMatch, error) {
	matches := make(map[string]*ArenaMatch)
	var match *ArenaMatch
	for i, s := range l.Segments {
		if s.IsMatchStart() {
			var err error
			match, err = s.ParseMatchStart()
			if err != nil {
				log.Printf("error parsing match start: %v\n", err.Error())
				continue
			}
			// Depending on the mode the player is playing we need to
			// look for a different thing.
			for j := i; j >= 0; j-- {
				s := l.Segments[j]
				if s.JoinedEvent() {
					course, err := s.ParseJoinedEvent()
					if err != nil {
						log.Printf("error parsing get player course: %v\n", err.Error())
						break
					}
					match.CourseDeck = course.CourseDeck
					break
				}
			}
			matches[match.MatchID] = match
		}
		if match != nil && s.IsMatchEvent() {
			event, err := s.ParseMatchEvent()
			if err != nil {
				log.Printf("error getting match event: %v\n", err.Error())
				continue
			}
			match.LogMatchEvent(event)
		}
		if s.IsMatchEnd() {
			end, err := s.ParseMatchEnd()
			if err != nil {
				log.Printf("error parsing match end: %v\n", err.Error())
				break
			}
			m := end.Params.PayloadObject
			match = matches[m.MatchID]
			// This is bad.
			if m.MatchID != match.MatchID {
				log.Printf("error: end.MatchID != start.MatchID. Unsure state of what happened to last match. Skipping.")
				continue
			}
			match.SeatID = m.SeatID
			match.TeamID = m.TeamID
			match.GameNumber = m.GameNumber
			match.WinningTeamID = m.WinningTeamID
			match.WinningReason = m.WinningReason
			match.TurnCount = m.TurnCount
			match.SecondsCount = m.SecondsCount
			matches[m.MatchID] = match
			match = nil
		}
	}
	var found []*ArenaMatch
	for _, value := range matches {
		found = append(found, value)
	}
	return found, nil
}

// Events finds Arena Events in the logs
func (l *Log) Events() ([]*ArenaEvent, error) {
	events := make([]*ArenaEvent, 0)
	for i, s := range l.Segments {
		if s.IsClaimPrize() {
			event := &ArenaEvent{}
			claim, err := s.ParseEventClaimPrize()
			if err != nil {
				continue
			}
			event.ClaimPrize = claim
			// Go back a few segments to find the inventory change
			for j := i; j > i-10; j-- {
				back := l.Segments[j]
				if back.IsInventoryUpdate() {
					update, err := back.ParseInventoryUpdate()
					if err != nil {
						break
					}
					event.Prize = update
					break
				}
			}
			events = append(events, event)
		}
	}
	return events, nil
}

// Find and Parse 1/8/2019 2:07:00 PM
func parseDate(line []byte) *time.Time {
	date := findDate.Find(line)
	t, err := time.ParseInLocation(dateLayout, string(date), time.Local)
	if err != nil {
		return nil
	}
	t = t.UTC()
	return &t
}

func parseType(b []byte) SegmentType {
	for s, r := range segmentTypeChecks {
		if r.Match(b) {
			return s
		}
	}
	return Unknown
}

// String to pointer
func String(s string) *string {
	return &s
}

// Int to pointer
func Int(i int) *int {
	return &i
}
