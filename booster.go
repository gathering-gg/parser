package gathering

import (
	"encoding/json"
	"time"
)

// ==> PlayerInventory.CrackBoostersV3(276):
// <== PlayerInventory.CrackBoostersV3(276)
// 1240429

// BoosterCard is a card opened in a booster
type BoosterCard struct {
	GrpID       int    `json:"grpId"`
	GoldAwarded int    `json:"goldAwarded"`
	GemsAwarded int    `json:"gemsAwarded"`
	Set         string `json:"set"` // Empty
}

// Booster is an opened Booster
type Booster struct {
	CardsOpened            []BoosterCard `json:"cardsOpened"`
	TotalVaultProgress     float64       `json:"totalVaultProgress"`
	WildCardTrackMoves     int           `json:"wildCardTrackMoves"`
	WildCardTrackPosition  int           `json:"wildCardTrackPosition"`
	WildCardTrackCommons   int           `json:"wildCardTrackCommons"`
	WildCardTrackUnCommons int           `json:"wildCardTrackUncommons"`
	WildCardTrackRares     int           `json:"wildCardTrackRares"`
	WildCardTrackMythics   int           `json:"wildCardTrackMythics"`
	OpenedAt               time.Time     `json:"openedAt"`
}

// IsCrackBooster checks to see if the user opened a booster
func (s *Segment) IsCrackBooster() bool {
	return s.SegmentType == CrackBooster
}

// ParseCrackBooster parses a booster from the log
func (s *Segment) ParseCrackBooster() (*Booster, error) {
	var booster Booster
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &booster)
	if s.Time != nil {
		booster.OpenedAt = *s.Time
	}
	return &booster, err
}
