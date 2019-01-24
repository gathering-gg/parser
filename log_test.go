package gathering

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fileAsString(f string, t *testing.T) string {
	file, err := os.Open(f)
	if err != nil {
		t.Fatalf("File not found")
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Error reading file")
	}
	return string(raw[:])
}

func TestLogFindCollection(t *testing.T) {
	a := assert.New(t)
	files := []string{
		"test/output_log0.txt",
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
}

func TestLogFindRank(t *testing.T) {
	a := assert.New(t)
	files := map[string]*ArenaRankInfo{
		"test/output_log0.txt": &ArenaRankInfo{
			PlayerID:                 String("EZIDLEQCFFAMLE27DG4TFGLT5Q"),
			ConstructedSeasonOrdinal: Int(1),
			ConstructedClass:         String("Gold"),
			ConstructedLevel:         Int(3),
			ConstructedStep:          Int(5),
			ConstructedMatchesWon:    Int(63),
			ConstructedMatchesLost:   Int(56),
			ConstructedMatchesDrawn:  Int(0),
			LimitedSeasonOrdinal:     Int(1),
			LimitedClass:             String("Gold"),
			LimitedLevel:             Int(4),
			LimitedStep:              Int(1),
			LimitedMatchesWon:        Int(33),
			LimitedMatchesLost:       Int(37),
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
		"test/output_log0.txt": &ArenaPlayerInventory{
			PlayerID:        "EZIDLEQCFFAMLE27DG4TFGLT5Q",
			WcCommon:        15,
			WcUncommon:      22,
			WcRare:          14,
			WcMythic:        9,
			Gold:            9975,
			Gems:            6220,
			DraftTokens:     0,
			SealedTokens:    0,
			WcTrackPosition: 6,
			VaultProgress:   39.4,
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
	raw := fileAsString("test/output_log0.txt", t)
	alog, err := ParseLog(raw)
	a.Nil(err)
	name, err := alog.Auth()
	a.Nil(err)
	a.Equal("Abattoir#66546", name)
}

func TestLogFindDecks(t *testing.T) {
	a := assert.New(t)
	files := map[string]int{
		"test/output_log0.txt": 12,
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
		"test/output_log0.txt": 8,
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
		}
	}
}
