package gathering

import (
	"bytes"
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
	EventGetPlayerCourse
	EventDeckSubmit
	EventMatchCreated
	PlayerAuth
	MatchStart
	MatchEnd
	MatchEvent
	CrackBooster
	InventoryRankUpdated
	EventClaimPrize
	IncomingInventoryUpdate
	DuelSceneSideboardingStart
	DuelSceneSideboardingStop
	MatchCompleted
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
	MatchEvent:                        regexp.MustCompile(`"GREMessageType_GameStateMessage"|GameStateType_Diff`),
	EventDeckSubmit:                   regexp.MustCompile(`<==\sEvent\.DeckSubmit\(\d+\)`),
	CrackBooster:                      regexp.MustCompile(`<==\sPlayerInventory\.CrackBoostersV3\(\d+\)`),
	InventoryRankUpdated:              regexp.MustCompile(`Incoming\sRank\.Updated`),
	EventClaimPrize:                   regexp.MustCompile(`<==\sEvent\.ClaimPrize\(\d+\)`),
	IncomingInventoryUpdate:           regexp.MustCompile(`Incoming\sInventory\.Updated`),
	DuelSceneSideboardingStart:        regexp.MustCompile(`DuelScene\.SideboardingStart`),
	DuelSceneSideboardingStop:         regexp.MustCompile(`DuelScene\.SideboardingStop`),
	MatchCompleted:                    regexp.MustCompile(`MatchGameRoomStateType_MatchCompleted`),
}

var cleaners = []*regexp.Regexp{
	regexp.MustCompile(`<<<<<<<<<<.*`),
	regexp.MustCompile(`\[\w.*`),
	regexp.MustCompile(`\dx[\d\w]+.*`),
	regexp.MustCompile(`(?m)ZoneTransferUXEvent.*`),
	regexp.MustCompile(`(?sm)BIError - GRE.Notification:.*`),
}
var clean = []byte(`$1.$2`)

var trimLeft = func(r rune) bool {
	return r != '{' && r != '['
}
var trimRight = func(r rune) bool {
	return r != '}' && r != ']'
}

// Segment is a piece of the log
type Segment struct {
	LoggerType  LoggerType
	Time        *time.Time
	SegmentType SegmentType
	Text        []byte
	Range       []int
	Line        []byte
}

// JSON parses the text as JSON
func (s *Segment) JSON(v interface{}) error {
	return json.Unmarshal(s.Text, v)
}

func stripNonJSON(b []byte) []byte {
	for _, c := range cleaners {
		b = c.ReplaceAll(b, clean)
	}
	b = bytes.TrimLeftFunc(b, trimLeft)
	b = bytes.TrimRightFunc(b, trimRight)
	return b
}
