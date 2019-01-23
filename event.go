package gathering

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

// ArenaEventState is the state of an event
type ArenaEventState int

// ArenaEventStates
const (
	NotJoined ArenaEventState = iota
	Playing
	ClaimPrize
	Complete
)

// ArenaEvent encapsulates the various stages a player may be in an event.
// When the log is parsed, the player may have started an event, be in the middle
// of an event, or just signed on to finish an event.
// The server will use the ID to track individual events.
type ArenaEvent struct {
	ID     string
	Joined *ArenaEventJoin
	Payed  *ArenaEventPayEntry
	// TODO: Matches
	Prize ArenaEventClaimPrizeRequest
}

// ArenaEventJoin is the payload when a user joins an event
// CardPool: Used in draft to list the cards available in the deck
type ArenaEventJoin struct {
	ID                string
	InternalEventName string
	CurrentEventState ArenaEventState
	CurrentModule     string
	CardPool          []int
	CourseDeck        ArenaDeck
}

// ArenaWinLossGate tracks your results as you go through an event
type ArenaWinLossGate struct {
	MaxWins           int
	MaxLosses         int
	CurrentWins       int
	CurrentLosses     int
	ProcessedMatchIDs []string
}

// ArenaEventClaimPrizeRequestParams are the params sent when a client requests
// the prize (Which prize does it request?)
type ArenaEventClaimPrizeRequestParams struct {
	EventName *string `json:"eventName"`
}

// ArenaEventClaimPrizeRequest is the inventory update from an event to know what the
// prizes are.
type ArenaEventClaimPrizeRequest struct {
	params *ArenaEventClaimPrizeRequestParams
}

// ArenaModuleInstanceData is instance data in a request.
// This is used in:
// * ArenaEventPayEntry: To see how they paid.
type ArenaModuleInstanceData struct {
	HasPaidEntry string
	WinLossGate  *ArenaWinLossGate
}

// ArenaEventPayEntry is the payload when the user pays for an event
type ArenaEventPayEntry struct {
	ID                 string
	InternalEventName  string
	ModuleInstanceData ArenaModuleInstanceData
	CurrentEventState  ArenaEventState
	CurrentModule      string
	CardPool           interface{}
	CourseDeck         *ArenaDeck
}

// ArenaEventGetPlayerCourse is fired when the player goes to find a new match.
// JoinQueue typically follows
type ArenaEventGetPlayerCourse struct {
	ID                 string
	InternalEventName  string
	ModuleInstanceData ArenaModuleInstanceData
	CurrentEventState  ArenaEventState
	CurrentModule      string
	CardPool           interface{}
	CourseDeck         *ArenaDeck
}

// IsEventJoin checks if a segment contains an Event Join
func (s *Segment) IsEventJoin() bool {
	return s.SegmentType == EventJoin
}

// IsEventGetPlayerCourse does this segment contain the player course
func (s *Segment) IsEventGetPlayerCourse() bool {
	return s.SegmentType == EventGetPlayerCourse
}

// IsEventDeckSubmit does this segment contain a deck submit for the play queue
func (s *Segment) IsEventDeckSubmit() bool {
	return s.SegmentType == EventDeckSubmit
}

// JoinedEvent is a higher level function to find if you joined
// any queue with a deck. Works with events and "play"
func (s *Segment) JoinedEvent() bool {
	return s.IsEventGetPlayerCourse() || s.IsEventDeckSubmit()
}

// ParseEventJoin parses out an event from JSON
func (s *Segment) ParseEventJoin() (*ArenaEventJoin, error) {
	text := stripNonJSON(s.Text)
	var join ArenaEventJoin
	err := json.Unmarshal([]byte(text), &join)
	if err != nil {
		return nil, err
	}
	return &join, nil
}

// ParseEventPayEntry parses out a pay entry value
func (s *Segment) ParseEventPayEntry() (*ArenaEventPayEntry, error) {
	text := stripNonJSON(s.Text)
	var pay ArenaEventPayEntry
	err := json.Unmarshal([]byte(text), &pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil

}

// ParseJoinedEvent gets the player course, which is another good way
// to verify the deck the player is using going into a game.
func (s *Segment) ParseJoinedEvent() (*ArenaEventGetPlayerCourse, error) {
	var course ArenaEventGetPlayerCourse
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &course)
	return &course, err
}

// ParseMatches finds the matches in a log
func ParseMatches(raw string) []ArenaMatch {
	texts := splitLogText(raw, logSplitRegex)
	var match *ArenaMatch
	var matches []ArenaMatch
	for _, t := range texts {
		isMatchDeck := regexp.MustCompile(isMatchPlayerCourse)
		isMatchStart := regexp.MustCompile(isMatchStartRegex)
		isMatchEnd := regexp.MustCompile(isMatchEndRegex)
		// The Player Course shows what they started searching for, and with
		// which deck, which we need to know what they played with
		if isMatchDeck.MatchString(t) {
			if match != nil {
				match = nil
			}
			playerCourse := strings.SplitN(t, "\n", 3)[2]
			if err := parseJSONBackoff(playerCourse, &match); err != nil {
				log.Printf("Error Parsing Player Course: %v", err.Error())
				continue
			}
		}
		if isMatchStart.MatchString(t) {
			incomingMatchJSON := strings.SplitN(t, "\n", 2)[1]
			// Need to chomp off the first part until we get to the JSON
			incomingMatchJSON = strings.TrimPrefix(incomingMatchJSON, "(-1) Incoming Event.MatchCreated ")
			if err := json.Unmarshal([]byte(incomingMatchJSON), &match); err != nil {
				log.Printf("Error Parsing Match Start: %v", err.Error())
				continue
			}
		}
		// Okay, we have a match, now what was the result?
		if isMatchEnd.MatchString(t) && match != nil {
			matchEndJSON := strings.SplitN(t, "\n", 3)[2]
			var result ArenaMatchEnd
			err := json.Unmarshal([]byte(matchEndJSON), &result)
			if err != nil {
				log.Printf("Error Parsing Match: %v", err.Error())
				continue
			}
			match.SeatID = result.Params.PayloadObject.SeatID
			match.TeamID = result.Params.PayloadObject.TeamID
			match.GameNumber = result.Params.PayloadObject.GameNumber
			match.WinningTeamID = result.Params.PayloadObject.WinningTeamID
			match.WinningReason = result.Params.PayloadObject.WinningReason
			match.TurnCount = result.Params.PayloadObject.TurnCount
			match.SecondsCount = result.Params.PayloadObject.SecondsCount
			matches = append(matches, *match)
			match = nil
		}
	}
	return matches
}
