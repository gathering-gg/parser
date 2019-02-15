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
	paths := []string{
		"test/output_log0.txt",
	}
	for _, p := range paths {
		f, err := os.Open(p)
		a.Nil(err)
		alog, err := ParseLog(f)
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
		o, _ := os.Open(f)
		alog, err := ParseLog(o)
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
		o, _ := os.Open(f)
		alog, err := ParseLog(o)
		a.Nil(err)
		inv, err := alog.Inventory()
		a.Nil(err)
		a.True(reflect.DeepEqual(inv, expected))
	}
}

func TestLogParseAuth(t *testing.T) {
	a := assert.New(t)
	f, _ := os.Open("test/output_log0.txt")
	alog, err := ParseLog(f)
	a.Nil(err)
	name, err := alog.Auth()
	a.Nil(err)
	a.Equal("Abattoir#66546", string(name))
}

func TestLogFindDecks(t *testing.T) {
	a := assert.New(t)
	files := map[string]int{
		"test/output_log0.txt": 12,
	}
	for f, i := range files {
		o, _ := os.Open(f)
		alog, err := ParseLog(o)
		a.Nil(err)
		decks, err := alog.Decks()
		a.Nil(err)
		a.Len(decks, i)
	}
}

func TestLogFindMatches(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/valentines-2019-update.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	matches, err := alog.Matches()
	a.Nil(err)
	a.Len(matches, 5)
	for _, m := range matches {
		a.NotNil(m.MatchID)
		a.NotNil(m.CourseDeck)
	}
}

func TestLogMatchRecap(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/valentines-2019-update.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	matches, err := alog.Matches()
	a.Len(matches, 5)
	var match *ArenaMatch
	for _, m := range matches {
		if m.MatchID == "93958637-81bb-4b15-a48a-340d264682db" {
			match = m
		}
	}
	a.Len(match.SeenObjects[1], 11)
	a.Len(match.SeenObjects[2], 8)
}

func TestLogCrackBooster(t *testing.T) {
	a := assert.New(t)
	file := "test/new-deck-constructed-7-1-daily-open-booster.txt"
	f, _ := os.Open(file)
	alog, err := ParseLog(f)
	a.Nil(err)
	boosters, err := alog.Boosters()
	a.Nil(err)
	a.Len(boosters, 1)
	b := boosters[0]
	a.Len(b.CardsOpened, 8)
	a.Equal(69167, b.CardsOpened[0].GrpID)
}

func TestLogEvents(t *testing.T) {
	a := assert.New(t)
	file := "test/feb-14-2018-update.txt"
	f, _ := os.Open(file)
	alog, err := ParseLog(f)
	a.Nil(err)
	eventResults, err := alog.Events()
	a.Nil(err)
	a.Len(eventResults, 1)
	e := eventResults[0]
	a.NotNil(e.ClaimPrize)
	a.NotNil(e.Prize)
	a.Equal("d2fcb515-beb0-41c2-a069-6a6a7aa3d099", e.ClaimPrize.ID)
	a.Equal("Valentines_2019", e.ClaimPrize.InternalEventName)
	a.Equal(5, e.ClaimPrize.ModuleInstanceData.WinLossGate.MaxWins)
	a.Equal(5, e.ClaimPrize.ModuleInstanceData.WinLossGate.CurrentWins)
	a.Equal(0, e.ClaimPrize.ModuleInstanceData.WinLossGate.CurrentLosses)
	a.Equal(70140, e.Prize.Delta.CardsAdded[0])
	a.Equal(70141, e.Prize.Delta.CardsAdded[1])
	a.Equal(0, e.Prize.Delta.GoldDelta)
}
