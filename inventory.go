package gathering

import (
	"encoding/json"
)

// ArenaPlayerInventory is your player profile details
type ArenaPlayerInventory struct {
	PlayerID        string  `json:"playerId"`
	WcCommon        int     `json:"wcCommon"`
	WcUncommon      int     `json:"wcUncommon"`
	WcRare          int     `json:"wcRare"`
	WcMythic        int     `json:"wcMythic"`
	Gold            int     `json:"gold"`
	Gems            int     `json:"gems"`
	DraftTokens     int     `json:"draftTokens"`
	SealedTokens    int     `json:"sealedTokens"`
	WcTrackPosition int     `json:"wcTrackPosition"`
	VaultProgress   float64 `json:"vaultProgress"`
}

// ArenaPlayerInventoryBooster is a struct which holds the type of booster and
// how many of that booster a player has.
type ArenaPlayerInventoryBooster struct {
	CollationID int
	Count       int
}

// ArenaInventoryUpdateDelta holds the delta change in a players inventory
type ArenaInventoryUpdateDelta struct {
	GemsDelta          int                           `json:"gemsDelta"`
	BoosterDelta       []ArenaPlayerInventoryBooster `json:"boosterDelta"`
	CardsAdded         []int                         `json:"cardsAdded"`
	DecksAdded         []interface{}                 `json:"decksAdded"`
	VanityItemsAdded   []interface{}                 `json:"vanityItemsAdded"`
	VanityItemsRemoved []interface{}                 `json:"vanityItemsRemoved"`
	DraftTokensDelta   int                           `json:"draftTokensDelta"`
	GoldDelta          int                           `json:"goldDelta"`
	SealedTokensDelta  int                           `json:"sealedTokensDelta"`
	VaultProgressDelta float64                       `json:"vaultProgressDelta"`
	WcCommonDelta      int                           `json:"wcCommonDelta"`
	WcUncommonDelta    int                           `json:"wcUncommonDelta"`
	WcRareDelta        int                           `json:"wcRarreDelta"`
	WcMythicDelta      int                           `json:"wcMythicDelta"`
}

// ArenaInventoryUpdate holds the incoming update for the player
type ArenaInventoryUpdate struct {
	Delta   *ArenaInventoryUpdateDelta `json:"delta"`
	Context string                     `json:"context"`
}

// IsPlayerInventory checks if a segment contains player inventory
func (s *Segment) IsPlayerInventory() bool {
	return s.SegmentType == PlayerInventoryGetPlayerInventory
}

// IsInventoryUpdate checks if a segment is an inventory update
func (s *Segment) IsInventoryUpdate() bool {
	return s.SegmentType == IncomingInventoryUpdate
}

// ParsePlayerInventory parses the player inventory information from a segment
func (s *Segment) ParsePlayerInventory() (*ArenaPlayerInventory, error) {
	var inv ArenaPlayerInventory
	err := json.Unmarshal(stripNonJSON(s.Text), &inv)
	return &inv, err
}

// ParseInventoryUpdate parses an incoming inventory update
func (s *Segment) ParseInventoryUpdate() (*ArenaInventoryUpdate, error) {
	var update ArenaInventoryUpdate
	err := json.Unmarshal(stripNonJSON(s.Text), &update)
	return &update, err
}
