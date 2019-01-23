package gathering

import (
	"encoding/json"
)

// ArenaDeck is the log format for an ArenaDeck
type ArenaDeck struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	ResourceID  string          `json:"resourceId"`
	DeckTileID  int             `json:"deckTileId"`
	MainDeck    []ArenaDeckCard `json:"mainDeck"`
	Sideboard   []ArenaDeckCard `json:"sideboard"`
}

// ArenaDeckCard hold the info of the cards in a deck
type ArenaDeckCard struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
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
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &decks)
	return decks, err
}
