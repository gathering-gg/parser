package gathering

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestParseSplit(t *testing.T) {
	arenaLog := fileAsString("test-logs/momir-madness.log", t)
	splitLogText2(arenaLog, "")
}

func TestParseCollection(t *testing.T) {
	t.Skip()
	arenaLog := fileAsString("output_log1.txt", t)
	col, err := ParseCollection(arenaLog)
	if err != nil {
		t.Fatalf("error parsing %v", err.Error())
	}
	assert.Len(t, col, 801)
}

/*
func TestParseDecks(t *testing.T) {
	arenaLog := fileAsString("output_log1.txt", t)
	decks, err := ParseDecks(arenaLog)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, d := range decks {
		fmt.Printf("Name: %v\n", d.Name)
	}
	assert.Len(t, decks, 8)
	assert.Equal(t, decks[0].Name, "Selesnya Convoke (C)")
}
*/

func TestParseInventory(t *testing.T) {
	t.Skip()
	arenaLog := fileAsString("output_log1.txt", t)
	inventory, err := ParsePlayerInventory(arenaLog)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := &ArenaPlayerInventory{
		PlayerID:        "EZIDLEQCFFAMLE27DG4TFGLT5Q",
		WcCommon:        15,
		WcUncommon:      19,
		WcRare:          15,
		WcMythic:        7,
		Gold:            1150,
		Gems:            9220,
		DraftTokens:     0,
		SealedTokens:    0,
		WcTrackPosition: 1,
		VaultProgress:   23.9,
	}
	assert.True(t, reflect.DeepEqual(inventory, expected))
}

func TestParseRank(t *testing.T) {
	t.Skip()
	arenaLog := fileAsString("output_log1.txt", t)
	rank, err := ParseRankInfo(arenaLog)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if *rank.LimitedClass != "Gold" || *rank.LimitedLevel != 4 {
		t.Fail()
	}
}

func TestParseAuthRequest(t *testing.T) {
	t.Skip()
	arenaLog := fileAsString("output_log1.txt", t)
	auth, err := ParseAuthRequest(arenaLog)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, auth.Payload.PlayerName, "Abattoir#66546")
}

func TestParseMatches(t *testing.T) {
	t.Skip()
	arenaLog := fileAsString("output_log1.txt", t)
	matches := ParseMatches(arenaLog)
	spew.Dump(matches)
	assert.Len(t, matches, 1)
	assert.Equal(t, *matches[0].SeatID, 1)
	assert.Equal(t, *matches[0].WinningTeamID, 2)
	assert.NotNil(t, matches[0].CourseDeck)
}
