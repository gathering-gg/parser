package gathering

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsArenaDeck(t *testing.T) {
	s := &Segment{
		SegmentType: DeckGetDeckLists,
	}
	assert.True(t, s.IsArenaDecks())
}

func TestParseEmpty(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Text: []byte(``),
	}
	a.False(s.IsArenaDecks())
	decks, err := s.ParseArenaDecks()
	a.NotNil(err)
	a.Nil(decks)
	a.Equal("unexpected end of JSON input", err.Error())
}

func TestParseEmptyArrayWithText(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Text: []byte(`<== Deck.GetDeckLists(9) [] RandomOtherText`),
	}
	decks, err := s.ParseArenaDecks()
	a.Nil(err)
	a.Empty(decks)
}

func TestUnmarshal(t *testing.T) {
	a := assert.New(t)
	deck := ArenaDeck{
		ID:         "id",
		Name:       "name",
		DeckTileID: 1,
		MainDeck: []ArenaDeckCard{
			ArenaDeckCard{
				ID:       2,
				Quantity: 3,
			},
		},
		Sideboard: []ArenaDeckCard{},
	}
	data, _ := json.Marshal(deck)
	var parsed ArenaDeck
	json.Unmarshal(data, &parsed)
	a.Equal(2, parsed.MainDeck[0].ID)
	a.Equal(3, parsed.MainDeck[0].Quantity)

}
