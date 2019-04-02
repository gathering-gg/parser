package gathering

import (
	"encoding/json"
)

// ArenaDeck is the log format for an ArenaDeck
type ArenaDeck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Format      string `json:"format"`
	DeckTileID  int    `json:"deckTileId"`
	MainDeck    []int  `json:"mainDeck"`
	Sideboard   []int  `json:"sideboard"`
	// TODO: What types are these?
	// cardSkins
	// cardBack
}

// IsArenaDecks checks if a segment contains Arena Decks
func (s *Segment) IsArenaDecks() bool {
	return s.SegmentType == DeckGetDeckLists
}

// ParseArenaDecks parses out arena decks from a segment if present.
// Note, it is the caller's responsibility to check if this segment contains
// ArenaDecks by calling `IsArenaDecks()`
func (s *Segment) ParseArenaDecks() ([]ArenaDeck, error) {
	var decks []ArenaDeck
	err := json.Unmarshal(stripNonJSON(s.Text), &decks)
	return decks, err
}
