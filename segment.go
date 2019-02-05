package gathering

import (
	"encoding/json"
	"regexp"
	"time"
)

// LoggerType is the logger that created this segment (log section)
type LoggerType int

// The Available Log types
const (
	UnityLogger LoggerType = iota
	ClientGRE
)

// SegmentType is the type of log info
type SegmentType int

// The type of data you can expect to find in the text of the log
const (
	Unknown SegmentType = iota
	PlayerInventoryGetPlayerInventory
	PlayerInventoryGetPlayerCards
	DeckGetDeckLists
	EventGetCombinedRankInfo
	EventJoin
	EventPayEntry
	EventGetPlayerCourse // Momir: 18436
	EventDeckSubmit
	EventMatchCreated
	PlayerAuth
	MatchStart
	MatchEnd
	MatchEvent
	CrackBooster
)

var segmentTypeChecks = map[SegmentType]*regexp.Regexp{
	PlayerInventoryGetPlayerInventory: regexp.MustCompile(`<==\sPlayerInventory\.GetPlayerInventory\(\d*\)`),
	PlayerInventoryGetPlayerCards:     regexp.MustCompile(`<==\sPlayerInventory\.GetPlayerCardsV3\(\d*\)`),
	EventGetCombinedRankInfo:          regexp.MustCompile(`<==\sEvent\.GetCombinedRankInfo\(\d+\)`),
	DeckGetDeckLists:                  regexp.MustCompile(`<==\sDeck\.GetDeckLists\(\d+\)`),
	PlayerAuth:                        regexp.MustCompile(`"screenName":\s"(.*)"`),
	EventGetPlayerCourse:              regexp.MustCompile(`<==\sEvent\.GetPlayerCourse\(\d+\)`),
	MatchStart:                        regexp.MustCompile(`Incoming\sEvent\.MatchCreated`),
	MatchEnd:                          regexp.MustCompile(`DuelScene\.GameStop`),
	MatchEvent:                        regexp.MustCompile(`GREMessageType_GameStateMessage|GameStateType_Diff`),
	EventDeckSubmit:                   regexp.MustCompile(`<==\sEvent\.DeckSubmit\(\d+\)`),
	CrackBooster:                      regexp.MustCompile(`<==\sPlayerInventory\.CrackBoostersV3\(\d+\)`),
}

// Segment is a piece of the log
type Segment struct {
	LoggerType  LoggerType
	Time        *time.Time
	SegmentType SegmentType
	Text        string
	Range       []int
	Line        string
}

// JSON parses the text as JSON
func (s *Segment) JSON(i interface{}) error {
	return json.Unmarshal([]byte(s.Text), i)
}
