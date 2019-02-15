package gathering

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPlayerInventory(t *testing.T) {
	s := &Segment{
		SegmentType: PlayerInventoryGetPlayerInventory,
	}
	assert.True(t, s.IsPlayerInventory())
}

func TestParseGetPlayerInventory(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Text: []byte(`
<== PlayerInventory.GetPlayerInventory(14)
{
  "playerId": "EZIDLEQCFFAMLE27DG4TFGLT5Q",
  "wcCommon": 11,
  "wcUncommon": 12,
  "wcRare": 11,
  "wcMythic": 7,
  "gold": 5050,
  "gems": 9220,
  "draftTokens": 0,
  "sealedTokens": 0,
  "wcTrackPosition": 1,
  "vaultProgress": 24.8,
  "boosters": [
    {
      "collationId": 100008,
      "count": 0
    },
    {
      "collationId": 100009,
      "count": 5
    },
    {
      "collationId": 100007,
      "count": 0
    }
  ],
  "vanityItems": {
    "pets": [],
    "avatars": [],
    "cardBacks": []
  },
  "vanitySelections": {
    "avatarSelection": null,
    "avatarModSelection": null,
    "cardBackSelection": null,
    "cardBackModSelection": null,
    "petSelection": null,
    "petModSelection": null
  }
}
`),
	}
	inv, err := s.ParsePlayerInventory()
	a.Nil(err)
	a.True(reflect.DeepEqual(&ArenaPlayerInventory{
		PlayerID:        "EZIDLEQCFFAMLE27DG4TFGLT5Q",
		WcCommon:        11,
		WcUncommon:      12,
		WcRare:          11,
		WcMythic:        7,
		Gold:            5050,
		Gems:            9220,
		DraftTokens:     0,
		SealedTokens:    0,
		WcTrackPosition: 1,
		VaultProgress:   24.8,
	}, inv))
}
