package gathering

import (
	"encoding/json"
)

// ArenaEvent encapsulates the various stages a player may be in an event.
// When the log is parsed, the player may have started an event, be in the middle
// of an event, or just signed on to finish an event.
// The server will use the ID to track individual events.
type ArenaEvent struct {
	ClaimPrize *ArenaEventClaimPrize `json:"claimPrize"`
	Prize      *ArenaInventoryUpdate `json:"prize"`
}

// ArenaEventJoin is the payload when a user joins an event
// CardPool: Used in draft to list the cards available in the deck
type ArenaEventJoin struct {
	ID                string
	InternalEventName string
	CurrentEventState string
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

// ArenaEventClaimPrizeModuleInstanceData has the data in claim prize
// about the event
type ArenaEventClaimPrizeModuleInstanceData struct {
	HasPaidEntry string                           `json:"HasPaidEntry"`
	DeckSelected bool                             `json:"DeckSelected"`
	WinLossGate  *ArenaEventClaimPrizeWinLossGate `json:"WinLossGate"`
}

// ArenaEventClaimPrizeWinLossGate has the player record for an event as well
// as max wins/losses and the matches played.
type ArenaEventClaimPrizeWinLossGate struct {
	MaxWins           int      `json:"MaxWins"`
	MaxLosses         int      `json:"MaxLosses"`
	CurrentWins       int      `json:"CurrentWins"`
	CurrentLosses     int      `json:"CurrentLosses"`
	ProcessedMatchIds []string `json:"ProcessedMatchIds"`
}

// ArenaEventClaimPrize is what is sent to the client when the user claims their
// prize after finishing an event
type ArenaEventClaimPrize struct {
	ID                 string                                  `json:"Id"`
	InternalEventName  string                                  `json:"InternalEventName"`
	ModuleInstanceData *ArenaEventClaimPrizeModuleInstanceData `json:"ModuleinstanceData"`
	CurrentEventState  string                                  `json:"CurrentEventState"`
	CurrentModule      string                                  `json:"CurrentModule"`
	CardPool           []int                                   `json:"Cardpool"`
	CourseDeck         *ArenaDeck                              `json:"CourseDeck"`
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
	CurrentEventState  string
	CurrentModule      string
	CardPool           []int
	CourseDeck         *ArenaDeck
}

// ArenaEventGetPlayerCourse is fired when the player goes to find a new match.
// JoinQueue typically follows
type ArenaEventGetPlayerCourse struct {
	ID                 string
	InternalEventName  string
	ModuleInstanceData ArenaModuleInstanceData
	CurrentEventState  string
	CurrentModule      string
	CardPool           []int
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

// IsClaimPrize checks if this segment claims a prize
func (s *Segment) IsClaimPrize() bool {
	return s.SegmentType == EventClaimPrize
}

// ParseEventJoin parses out an event from JSON
func (s *Segment) ParseEventJoin() (*ArenaEventJoin, error) {
	var join ArenaEventJoin
	err := json.Unmarshal(stripNonJSON(s.Text), &join)
	return &join, err
}

// ParseEventPayEntry parses out a pay entry value
func (s *Segment) ParseEventPayEntry() (*ArenaEventPayEntry, error) {
	var pay ArenaEventPayEntry
	err := json.Unmarshal(stripNonJSON(s.Text), &pay)
	return &pay, err

}

// ParseJoinedEvent gets the player course, which is another good way
// to verify the deck the player is using going into a game.
func (s *Segment) ParseJoinedEvent() (*ArenaEventGetPlayerCourse, error) {
	var course ArenaEventGetPlayerCourse
	err := json.Unmarshal(stripNonJSON(s.Text), &course)
	return &course, err
}

// ParseEventClaimPrize parses an event claim prize
func (s *Segment) ParseEventClaimPrize() (*ArenaEventClaimPrize, error) {
	var prize ArenaEventClaimPrize
	err := json.Unmarshal(stripNonJSON(s.Text), &prize)
	return &prize, err
}
