package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEventJoin(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Line: []byte(`[UnityCrossThreadLogger]1/8/2019 2:11:59 PM`),
		Text: []byte(`
<== Event.Join(57)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {},
  "CurrentEventState": "state",
  "CurrentModule": "PayEntry",
  "CardPool": null,
  "CourseDeck": null
}`),
	}
	joined, err := s.ParseEventJoin()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", joined.ID)
	a.Equal("Momir_20190107", joined.InternalEventName)
	a.Equal("PayEntry", joined.CurrentModule)
}

func TestParsePayEntry(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Line: []byte(`[UnityCrossThreadLogger]1/8/2019 2:11:59 PM`),
		Text: []byte(`
<== Event.PayEntry(58)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {
    "HasPaidEntry": "Gold"
  },
  "CurrentEventState": "state",
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
        67015,
        12,
        67017,
        12,
        67019,
        12,
        67021,
        12,
        67023,
        12
    ],
    "sideboard": [],
    "lastUpdated": "0001-01-01T00:00:00",
    "lockedForUse": false,
    "lockedForEdit": false,
    "isValid": false
  }
}`),
	}
	payed, err := s.ParseEventPayEntry()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", payed.ID)
	a.Equal("Gold", payed.ModuleInstanceData.HasPaidEntry)
	a.Equal("Momir", payed.CourseDeck.Name)
}

func TestParseEventGetPlayerCourse(t *testing.T) {
	a := assert.New(t)
	s := &Segment{
		Line: []byte(`[UnityCrossThreadLogger]1/8/2019 2:12:00 PM`),
		Text: []byte(`
<== Event.GetPlayerCourse(63)
{
  "Id": "6c236ae7-81ff-4024-a836-cc055194fafe",
  "InternalEventName": "Momir_20190107",
  "ModuleInstanceData": {
    "HasPaidEntry": "Gold"
  },
  "CurrentEventState": "Transition",
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
        67015,
        12,
        67017,
        12,
        67019,
        12,
        67021,
        12,
        67023,
        12
    ],
    "sideboard": [],
    "lastUpdated": "0001-01-01T00:00:00",
    "lockedForUse": false,
    "lockedForEdit": false,
    "isValid": false
  }
}
new prize bar state is: PrizeDisplay
`),
	}
	pc, err := s.ParseJoinedEvent()
	a.Nil(err)
	a.Equal("6c236ae7-81ff-4024-a836-cc055194fafe", pc.ID)
	a.Equal("Momir", pc.CourseDeck.Name)
}
