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

// IsRankInfo checks if a segment contains Rank Info
func (s *Segment) IsRankInfo() bool {
	return s.SegmentType == EventGetCombinedRankInfo
}

// ParseRankInfo parses the rank information out of a segment.
func (s *Segment) ParseRankInfo() (*ArenaRankInfo, error) {
	var rank ArenaRankInfo
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &rank)
	return &rank, err
}
