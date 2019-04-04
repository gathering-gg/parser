package gathering

import (
	"encoding/json"
	"strconv"
)

// ArenaDeck is our format for an Arena Deck
// Arena has changed the deck format, sometimes it has
// the old version of json objects for cards, sometimes it now contains
// a list of `ints`, pairs of card number and quantity.
// We handle both
type ArenaDeck struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Format      string               `json:"format"`
	DeckTileID  int                  `json:"deckTileId"`
	MainDeck    []ArenaDeckCard      `json:"mainDeck"`
	Sideboard   []ArenaDeckCard      `json:"sideboard"`
	CardSkins   []*ArenaDeckCardSkin `json:"cardSkins"`
	CardBack    string               `json:"cardBack"`
}

// ArenaDeckCardSkin contains which cards have which skins
type ArenaDeckCardSkin struct {
	GrpID int    `json:"grpId"`
	CCV   string `json:"ccv"` // No idea what this is
}

// ArenaDeckCard is our representation of the card
// This is not always what the log has, but we normalize it to
// this structure when we parse the JSON
type ArenaDeckCard struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

// UnmarshalJSON handles both strange cases of the deck format
// and normalizes it to what we expect
func (d *ArenaDeck) UnmarshalJSON(data []byte) error {
	var deck map[string]interface{}
	if err := json.Unmarshal(data, &deck); err != nil {
		return err
	}
	if _, ok := deck["id"]; !ok {
		return nil
	}
	d.ID = deck["id"].(string)
	if name, ok := deck["name"].(string); ok {
		d.Name = name
	}
	if description, ok := deck["description"].(string); ok {
		d.Description = description
	}
	if format, ok := deck["format"].(string); ok {
		d.Format = format
	}
	if deckTile, ok := deck["deckTileId"].(float64); ok {
		d.DeckTileID = int(deckTile)
	}
	if cardBack, ok := deck["cardBack"].(string); ok {
		d.CardBack = cardBack
	}
	d.MainDeck = getCards(deck["mainDeck"])
	d.Sideboard = getCards(deck["sideboard"])
	return nil
}

// getCards takes an ambiguous array from the log and
// turns it into our ArenaDeckCard. This is a little ugly, but the format is
// either: [int, int, int, int] (id, quantity) or,
// { id: "id", quantity: num }
func getCards(cards interface{}) []ArenaDeckCard {
	final := []ArenaDeckCard{}
	if array, ok := cards.([]interface{}); ok {
		for i := 0; i < len(array); i++ {
			val := array[i]
			switch val.(type) {
			case int:
				final = append(final, ArenaDeckCard{
					ID:       val.(int),
					Quantity: array[i+1].(int),
				})
				i++
			case map[string]interface{}:
				id, _ := strconv.Atoi(val.(map[string]interface{})["id"].(string))
				final = append(final, ArenaDeckCard{
					ID:       id,
					Quantity: int(val.(map[string]interface{})["quantity"].(float64)),
				})
			}
		}
	}
	return final
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
