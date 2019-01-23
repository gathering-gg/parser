package gathering

import (
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
		Text: "",
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
		Text: "<== Deck.GetDeckLists(9) [] RandomOtherText",
	}
	decks, err := s.ParseArenaDecks()
	a.Nil(err)
	a.Empty(decks)
}
