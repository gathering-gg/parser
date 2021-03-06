package gathering

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRankInfo(t *testing.T) {
	s := &Segment{
		SegmentType: EventGetCombinedRankInfo,
	}
	assert.True(t, s.IsRankInfo())
}

func TestParseEmptyRank(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Text: []byte(`<== Rank`),
	}
	_, err := s.ParseRankInfo()
	a.NotNil(err)
	a.Equal("unexpected end of JSON input", err.Error())
}

func TestParseRankInfo(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Text: []byte(`
<== Event.GetCombinedRankInfo(11)
{
  "playerId": "EZIDLEQCFFAMLE27DG4TFGLT5Q",
  "constructedSeasonOrdinal": 1,
  "constructedClass": "Silver",
  "constructedLevel": 2,
  "constructedStep": 3,
  "constructedMatchesWon": 16,
  "constructedMatchesLost": 13,
  "constructedMatchesDrawn": 0,
  "limitedSeasonOrdinal": 1,
  "limitedClass": "Gold",
  "limitedLevel": 4,
  "limitedStep": 0,
  "limitedMatchesWon": 29,
  "limitedMatchesLost": 31,
  "limitedMatchesDrawn": 0
}`),
	}
	rank, err := s.ParseRankInfo()
	a.Nil(err)
	a.Equal("EZIDLEQCFFAMLE27DG4TFGLT5Q", *rank.PlayerID)
	a.Equal(1, *rank.ConstructedSeasonOrdinal)
	a.Equal("Silver", *rank.ConstructedClass)
	a.Equal(1, *rank.LimitedSeasonOrdinal)
	a.Equal("Gold", *rank.LimitedClass)
}

func TestParseLogGetRank(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/output_log0.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	var rank *ArenaRankInfo
	for i := len(alog.Segments) - 1; i >= 0; i-- {
		s := alog.Segments[i]
		if s.IsRankInfo() {
			rank, err = s.ParseRankInfo()
			break
		}
	}
	a.Nil(err)
	a.Equal("EZIDLEQCFFAMLE27DG4TFGLT5Q", *rank.PlayerID)
	a.Equal("Gold", *rank.LimitedClass)
	a.Equal(4, *rank.LimitedLevel)
	a.Equal("Gold", *rank.ConstructedClass)
	a.Equal(56, *rank.ConstructedMatchesLost)
}

func TestRankUp(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("test/rank-up.txt")
	a.Nil(err)
	alog, err := ParseLog(f)
	a.Nil(err)
	rank, err := alog.Rank()
	a.Nil(err)
	a.Equal("Gold", *rank.ConstructedClass)
	a.Equal(4, *rank.ConstructedLevel)
	a.Equal(2, *rank.ConstructedStep)
}
