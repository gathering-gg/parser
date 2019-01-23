package gathering

import "encoding/json"

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
	GemsDelta          int
	BoosterDelta       []ArenaPlayerInventoryBooster
	CardsAdded         []int
	DecksAdded         []interface{}
	VanityItemsAdded   []interface{}
	VanityItemsRemoved []interface{}
	DraftTokensDelta   int
	GoldDelta          int
	SealedTokensDelta  int
	VaultProgressDelta float64
	WcCommonDelta      int
	WcUncommonDelta    int
	WcRareDelta        int
	WcMythicDelta      int
}

// ArenaInventoryUpdate holds the incoming update for the player
type ArenaInventoryUpdate struct {
	Delta   *ArenaInventoryUpdateDelta
	context *string
}

// IsPlayerInventory checks if a segment contains player inventory
func (s *Segment) IsPlayerInventory() bool {
	return s.SegmentType == PlayerInventoryGetPlayerInventory
}

// ParsePlayerInventory parses the player inventory information from a segment
func (s *Segment) ParsePlayerInventory() (*ArenaPlayerInventory, error) {
	var inv ArenaPlayerInventory
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &inv)
	return &inv, err
}
