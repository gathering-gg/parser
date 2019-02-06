package gathering

import (
	"encoding/json"
)

// ArenaRankInfo contains a players rank info
type ArenaRankInfo struct {
	PlayerID                 *string `json:"playerId"`
	ConstructedSeasonOrdinal *int    `json:"constructedSeasonOrdinal"`
	ConstructedClass         *string `json:"constructedClass"`
	ConstructedLevel         *int    `json:"constructedLevel"`
	ConstructedStep          *int    `json:"constructedStep"`
	ConstructedMatchesWon    *int    `json:"constructedMatchesWon"`
	ConstructedMatchesLost   *int    `json:"constructedMatchesLost"`
	ConstructedMatchesDrawn  *int    `json:"constructedMatchesDrawn"`
	LimitedSeasonOrdinal     *int    `json:"limitedSeasonOrdinal"`
	LimitedClass             *string `json:"limitedClass"`
	LimitedLevel             *int    `json:"limitedLevel"`
	LimitedStep              *int    `json:"limitedStep"`
	LimitedMatchesWon        *int    `json:"limitedMatchesWon"`
	LimitedMatchesLost       *int    `json:"limitedMatchesLost"`
	LimitedMatchesDrawn      *int    `json:"limitedMatchesDrawn"`
}

// Update updates rank info with a server update
func (a *ArenaRankInfo) Update(update *RankUpdated) {
	if update.RankUpdateType == "Constructed" {
		a.ConstructedClass = &update.NewClass
		a.ConstructedLevel = &update.NewLevel
		a.ConstructedStep = &update.NewStep
	} else {
		a.LimitedClass = &update.NewClass
		a.LimitedLevel = &update.NewLevel
		a.LimitedStep = &update.NewStep
	}
}

// RankUpdated holds the rank update info from the server
// RankUpdateType: 'Constructed' | 'Limited'
type RankUpdated struct {
	PlayerID         string `json:"playerId"`
	SeasonOrdinal    int    `json:"seasonOrdinal"`
	NewClass         string `json:"newClass"`
	OldClass         string `json:"oldClass"`
	NewLevel         int    `json:"newLevel"`
	OldLevel         int    `json:"oldLevel"`
	NewStep          int    `json:"newStep"`
	OldStep          int    `json:"oldStep"`
	WasLossProtected bool   `json:"wasLossProtected"`
	RankUpdateType   string `json:"rankUpdateType"`
}

// IsRankInfo checks if a segment contains Rank Info
func (s *Segment) IsRankInfo() bool {
	return s.SegmentType == EventGetCombinedRankInfo
}

// IsRankUpdated checks if a segment contains Rank Update
func (s *Segment) IsRankUpdated() bool {
	return s.SegmentType == InventoryRankUpdated
}

// ParseRankInfo parses the rank information out of a segment.
func (s *Segment) ParseRankInfo() (*ArenaRankInfo, error) {
	var rank ArenaRankInfo
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &rank)
	return &rank, err
}

// ParseRankUpdated parses the rank update
func (s *Segment) ParseRankUpdated() (*RankUpdated, error) {
	var update RankUpdated
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &update)
	return &update, err
}
