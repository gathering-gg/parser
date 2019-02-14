package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEventJoin(t *testing.T) {
	t.Skip()
	// TODO: Use new log to fix
	a := assert.New(t)
	s := &Segment{
		Line: "[UnityCrossThreadLogger]1/8/2019 2:11:59 PM",
		Text: `
<== Event.Join(57)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {},
  "CurrentEventState": 0,
  "CurrentModule": "PayEntry",
  "CardPool": null,
  "CourseDeck": null
}`,
	}
	joined, err := s.ParseEventJoin()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", joined.ID)
	a.Equal("Momir_20190107", joined.InternalEventName)
	a.Equal("PayEntry", joined.CurrentModule)
}

func TestParsePayEntry(t *testing.T) {
	t.Skip()
	a := assert.New(t)
	s := &Segment{
		Line: "[UnityCrossThreadLogger]1/8/2019 2:11:59 PM",
		Text: `
<== Event.PayEntry(58)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {
    "HasPaidEntry": "Gold"
  },
  "CurrentEventState": 1,
  "CurrentModule": "TransitionToMatches",
  "CardPool": null,
  "CourseDeck": {
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "Momir",
    "description": null,
    "format": null,
    "resourceId": "00000000-0000-0000-0000-000000000000",
    "deckTileId": null,
    "mainDeck": [
      {
        "id": "67015",
        "quantity": 12
      },
      {
        "id": "67017",
        "quantity": 12
      },
      {
        "id": "67019",
        "quantity": 12
      },
      {
        "id": "67021",
        "quantity": 12
      },
      {
        "id": "67023",
        "quantity": 12
      }
    ],
    "sideboard": [],
    "lastUpdated": "0001-01-01T00:00:00",
    "lockedForUse": false,
    "lockedForEdit": false,
    "isValid": false
  }
}`,
	}
	payed, err := s.ParseEventPayEntry()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", payed.ID)
	a.Equal("Gold", payed.ModuleInstanceData.HasPaidEntry)
	a.Equal("Momir", payed.CourseDeck.Name)
}

func TestParseEventGetPlayerCourse(t *testing.T) {
	t.Skip()
	a := assert.New(t)
	s := &Segment{
		Line: "[UnityCrossThreadLogger]1/8/2019 2:12:00 PM",
		Text: `
<== Event.GetPlayerCourse(63)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {
    "HasPaidEntry": "Gold"
  },
  "CurrentEventState": 1,
  "CurrentModule": "TransitionToMatches",
  "CardPool": null,
  "CourseDeck": {
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "Momir",
    "description": null,
    "format": null,
    "resourceId": "00000000-0000-0000-0000-000000000000",
    "deckTileId": null,
    "mainDeck": [
      {
        "id": "67015",
        "quantity": 12
      },
      {
        "id": "67017",
        "quantity": 12
      },
      {
        "id": "67019",
        "quantity": 12
      },
      {
        "id": "67021",
        "quantity": 12
      },
      {
        "id": "67023",
        "quantity": 12
      }
    ],
    "sideboard": [],
    "lastUpdated": "0001-01-01T00:00:00",
    "lockedForUse": false,
    "lockedForEdit": false,
    "isValid": false
  }
}
new prize bar state is: PrizeDisplay
`,
	}
	pc, err := s.ParseJoinedEvent()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", pc.ID)
	a.Equal("Momir", pc.CourseDeck.Name)
}
