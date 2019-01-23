package gathering

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestLogFindCollection(t *testing.T) {
	a := assert.New(t)
	files := []string{
		"test-logs/momir-madness.log",
		"test-logs/regular-play.txt",
	}
	for _, f := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		col, err := alog.Collection()
		a.Nil(err)
		a.True(len(col) > 500)
		for k, v := range col {
			a.True(v > 0)
			a.Len(k, 5)
		}
	}
	files = []string{
		"test-logs/claim-weekly-booster.txt",
	}
	for _, f := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		col, err := alog.Collection()
		a.Equal(ErrNotFound.Error(), err.Error())
		a.Len(col, 0)
	}
}

func TestLogFindRank(t *testing.T) {
	a := assert.New(t)
	files := map[string]*ArenaRankInfo{
		"test-logs/momir-madness.log": &ArenaRankInfo{
			PlayerID:                 String("EZIDLEQCFFAMLE27DG4TFGLT5Q"),
			ConstructedSeasonOrdinal: Int(1),
			ConstructedClass:         String("Silver"),
			ConstructedLevel:         Int(2),
			ConstructedStep:          Int(3),
			ConstructedMatchesWon:    Int(16),
			ConstructedMatchesLost:   Int(13),
			ConstructedMatchesDrawn:  Int(0),
			LimitedSeasonOrdinal:     Int(1),
			LimitedClass:             String("Gold"),
			LimitedLevel:             Int(4),
			LimitedStep:              Int(0),
			LimitedMatchesWon:        Int(29),
			LimitedMatchesLost:       Int(31),
			LimitedMatchesDrawn:      Int(0),
		},
		"test-logs/regular-play.txt": &ArenaRankInfo{
			PlayerID:                 String("EZIDLEQCFFAMLE27DG4TFGLT5Q"),
			ConstructedSeasonOrdinal: Int(1),
			ConstructedClass:         String("Gold"),
			ConstructedLevel:         Int(4),
			ConstructedStep:          Int(0),
			ConstructedMatchesWon:    Int(21),
			ConstructedMatchesLost:   Int(17),
			ConstructedMatchesDrawn:  Int(0),
			LimitedSeasonOrdinal:     Int(1),
			LimitedClass:             String("Gold"),
			LimitedLevel:             Int(4),
			LimitedStep:              Int(0),
			LimitedMatchesWon:        Int(29),
			LimitedMatchesLost:       Int(31),
			LimitedMatchesDrawn:      Int(0),
		},
	}
	for f, expected := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		rank, err := alog.Rank()
		a.Nil(err)
		a.True(reflect.DeepEqual(rank, expected))
	}
}

func TestLogFindInventory(t *testing.T) {
	a := assert.New(t)
	files := map[string]*ArenaPlayerInventory{
		"test-logs/momir-madness.log": &ArenaPlayerInventory{
			PlayerID:        "EZIDLEQCFFAMLE27DG4TFGLT5Q",
			WcCommon:        11,
			WcUncommon:      12,
			WcRare:          11,
			WcMythic:        7,
			Gold:            4650,
			Gems:            9220,
			DraftTokens:     0,
			SealedTokens:    0,
			WcTrackPosition: 1,
			VaultProgress:   24.8,
		},
		"test-logs/regular-play.txt": &ArenaPlayerInventory{
			PlayerID:        "EZIDLEQCFFAMLE27DG4TFGLT5Q",
			WcCommon:        6,
			WcUncommon:      8,
			WcRare:          8,
			WcMythic:        7,
			Gold:            6700,
			Gems:            9220,
			DraftTokens:     0,
			SealedTokens:    0,
			WcTrackPosition: 1,
			VaultProgress:   25.1,
		},
	}
	for f, expected := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		inv, err := alog.Inventory()
		a.Nil(err)
		a.True(reflect.DeepEqual(inv, expected))
	}
}

func TestLogParseAuth(t *testing.T) {
	a := assert.New(t)
	raw := fileAsString("test-logs/momir-madness.log", t)
	alog, err := ParseLog(raw)
	a.Nil(err)
	name, err := alog.Auth()
	a.Nil(err)
	a.Equal("Abattoir#66546", name)
}

func TestLogFindDecks(t *testing.T) {
	a := assert.New(t)
	files := map[string]int{
		"test-logs/momir-madness.log": 10,
		"test-logs/regular-play.txt":  11,
	}
	for f, i := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		decks, err := alog.Decks()
		a.Nil(err)
		a.Len(decks, i)
	}
}

func TestLogFindMatches(t *testing.T) {
	a := assert.New(t)
	files := map[string]int{
		"test-logs/momir-madness.log": 2,
		"test-logs/regular-play.txt":  1,
		"test-logs/constructed.txt":   9,
	}
	for f, i := range files {
		raw := fileAsString(f, t)
		alog, err := ParseLog(raw)
		a.Nil(err)
		matches, err := alog.Matches()
		a.Nil(err)
		a.Len(matches, i)
		for _, m := range matches {
			a.NotNil(m.CourseDeck)
			spew.Dump(m.SeenObjects)
		}
	}
}
