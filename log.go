package gathering

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

// ErrNotFound is the error returned when a log item is not found
var ErrNotFound = errors.New("not found")

// Log is the well-structured format of the output_log.txt, parsed into Segments
type Log struct {
	Segments []*Segment
}

// ParseLog returns a log file parsed into Segments
func ParseLog(raw string) (*Log, error) {
	split := regexp.MustCompile(`\r?\n`).Split(raw, -1)
	r := regexp.MustCompile(`\[UnityCrossThreadLogger\].*|\[Client GRE\]`)
	var segments []*Segment
	var text string
	previous := &Segment{
		Range: []int{0},
	}
	for i, s := range split {
		if r.MatchString(s) {
			t := UnityLogger
			if strings.Contains(s, "Client GRE") {
				t = ClientGRE
			}
			previous.Text = text
			previous.Range = append(previous.Range, i-1)
			previous.SegmentType = parseType(text)
			segments = append(segments, previous)
			text = ""
			previous = &Segment{
				LoggerType: t,
				Time:       parseDate(s),
				Line:       s,
				Range:      []int{i},
			}
		} else {
			text = text + "\n" + s
		}
	}
	previous.Text = text
	previous.Range = append(previous.Range, len(split)-1)
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
func (l *Log) Rank() (*ArenaRankInfo, error) {
	// TODO: Put in same loop
	for i := len(l.Segments) - 1; i >= 0; i-- {
		s := l.Segments[i]
		if s.IsRankInfo() {
			return s.ParseRankInfo()
		}
	}
	return nil, ErrNotFound
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
func (l *Log) Auth() (string, error) {
	// TODO: Put in same loop
	for _, s := range l.Segments {
		if s.IsPlayerAuth() {
			return s.ParseAuth()
		}
	}
	return "", ErrNotFound
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
		value.CondenseLogMatch()
		found = append(found, value)
	}
	return found, nil
}

// Find and Parse 1/8/2019 2:07:00 PM
func parseDate(line string) *time.Time {
	re := regexp.MustCompile(`(?m)\d+\/\d+\/\d{4}\s\d+:\d+:\d+\s[APM]{2}`)
	date := re.FindString(line)
	layout := "1/2/2006 15:04:05 PM"
	t, err := time.ParseInLocation(layout, date, time.Local)
	if err != nil {
		return nil
	}
	t = t.UTC()
	return &t
}

func parseType(text string) SegmentType {
	for s, r := range segmentTypeChecks {
		if r.MatchString(text) {
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
