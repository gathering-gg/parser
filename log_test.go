package gathering

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		o, err := os.Open(f)
		a.Nil(err)
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
		o, err := os.Open(f)
		a.Nil(err)
		alog, err := ParseLog(o)
		a.Nil(err)
		inv, err := alog.Inventory()
		a.Nil(err)
		a.True(reflect.DeepEqual(inv, expected))
	}
}

func TestLogParseAuth(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/output_log0.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	name, err := alog.Auth()
	a.Nil(err)
	a.Equal("Abattoir#66546", string(name))
}

func TestLogFindDecks(t *testing.T) {
	a := assert.New(t)
	files := map[string]int{
		"test/cosmetics.txt": 20,
	}
	for f, i := range files {
		o, err := os.Open(f)
		a.Nil(err)
		alog, err := ParseLog(o)
		a.Nil(err)
		decks, err := alog.Decks()
		a.Nil(err)
		a.Len(decks, i)
		for _, d := range decks {
			a.True(len(d.MainDeck) > 0)
		}
	}
}

func TestLogFindMatches(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/cosmetics.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	matches, err := alog.Matches()
	a.Nil(err)
	a.Len(matches, 1)
	for _, m := range matches {
		a.NotNil(m.MatchID)
		a.NotNil(m.CourseDeck)
	}
}

func TestLogMatchRecap(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/cosmetics.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	matches, err := alog.Matches()
	a.Len(matches, 1)
	var match *ArenaMatch
	for _, m := range matches {
		if m.MatchID == "57f28bff-c2d0-4a06-8b68-24a5c304e15e" {
			match = m
		}
	}
	game := match.Games[len(match.Games)-1]
	a.Len(game.SeenObjects[1], 8)
	a.Len(game.SeenObjects[2], 9)
}

func TestLogCrackBooster(t *testing.T) {
	a := assert.New(t)
	file := "test/new-deck-constructed-7-1-daily-open-booster.txt"
	f, err := os.Open(file)
	a.Nil(err)
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
	file := "test/march-constructed.txt"
	f, err := os.Open(file)
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	eventResults, err := alog.Events()
	a.Nil(err)
	a.Len(eventResults, 1)
	e := eventResults[0]
	a.NotNil(e.ClaimPrize)
	a.NotNil(e.Prize)
	a.Equal("00f2bcfe-e5c3-45ba-9950-05d10d2687ad", e.ClaimPrize.ID)
	a.Equal("Constructed_Event", e.ClaimPrize.InternalEventName)
	a.Equal(7, e.ClaimPrize.ModuleInstanceData.WinLossGate.MaxWins)
	a.Equal(4, e.ClaimPrize.ModuleInstanceData.WinLossGate.CurrentWins)
	a.Equal(3, e.ClaimPrize.ModuleInstanceData.WinLossGate.CurrentLosses)
	a.Equal(67692, e.Prize.Delta.CardsAdded[0])
	a.Equal(66821, e.Prize.Delta.CardsAdded[1])
	a.Equal(500, e.Prize.Delta.GoldDelta)
}

func TestLogBestOfThree(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/bo3-small.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	matches, err := alog.Matches()
	a.Nil(err)
	a.Len(matches, 1)
	for _, m := range matches {
		a.NotNil(m.MatchID)
		a.True(len(m.Games) > 1)
		for _, g := range m.Games {
			a.NotNil(g.GameStart)
			a.NotNil(g.Number)
			a.NotNil(g.MatchID)
			a.NotNil(g.SeatID)
			a.NotNil(g.TeamID)
			a.NotNil(g.WinningTeamID)
			a.NotNil(g.WinningReason)
		}
	}
}
