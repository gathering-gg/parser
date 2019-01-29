package gathering

import (
	"encoding/json"
	"fmt"
	"time"
)

// ArenaMatch is a match in Arena. May not be completed yet
type ArenaMatch struct {
	MatchID                        string                         `json:"matchId"`
	GameStart                      time.Time                      `json:"gameStart"`
	OpponentScreenName             string                         `json:"opponentScreenName"`
	OpponentIsWotc                 bool                           `json:"opponentIsWotc"`
	OpponentRankingClass           string                         `json:"opponentRankingClass"`
	OpponentRankingTier            int                            `json:"opponentRankingTier"`
	OpponentMythicPercentile       float64                        `json:"opponentMythicPercentile"`
	OpponentMythicLeaderboardPlace int                            `json:"opponentMythicLeaderboardPlace"`
	EventID                        string                         `json:"eventId"`
	SeatID                         *int                           `json:"seatId"`
	TeamID                         *int                           `json:"teamId"`
	GameNumber                     *int                           `json:"gameNumber"`
	WinningTeamID                  *int                           `json:"winningTeamId"`
	WinningReason                  *string                        `json:"winningReason"`
	TurnCount                      *int                           `json:"turnCount"`
	SecondsCount                   *int                           `json:"secondsCount"`
	CourseDeck                     *ArenaDeck                     `json:"CourseDeck"`
	SeenObjects                    map[int][]ArenaMatchGameObject `json:"seenObjects"`
}

// LogMatchEvent adds an event to the log
func (a *ArenaMatch) LogMatchEvent(event *ArenaMatchEvent) {
	if a.SeenObjects == nil {
		a.SeenObjects = make(map[int][]ArenaMatchGameObject)
	}
	for _, m := range event.GreToClientEvent.GreToClientMessages {
		gsm := m.GameStateMessage
		for _, o := range gsm.GameObjects {
			a.SeenObjects[o.OwnerSeatID] = append(a.SeenObjects[o.OwnerSeatID], o)
		}
	}
	for k, v := range a.SeenObjects {
		uniq := make(map[string]bool)
		var objects []ArenaMatchGameObject
		for _, o := range v {
			if _, ok := uniq[o.Hash()]; !ok && o.Type == "GameObjectType_Card" {
				uniq[o.Hash()] = true
				objects = append(objects, o)
			}
		}
		a.SeenObjects[k] = objects
	}
}

// ArenaMatchEvent is an event in the match
type ArenaMatchEvent struct {
	TransactionID    string           `json:"transactionId"`
	Timestamp        string           `json:"timestamp"`
	GreToClientEvent GreToClientEvent `json:"greToClientEvent"`
}

// GreToClientEvent see log
type GreToClientEvent struct {
	GreToClientMessages []GreToClientMessages `json:"greToClientMessages"`
}

// GreToClientMessages see log
type GreToClientMessages struct {
	Type             string           `json:"type"`
	GameStateMessage GameStateMessage `json:"gameStateMessage"`
}

// GameStateMessage see log
type GameStateMessage struct {
	Type        string                 `json:"type"`
	GameObjects []ArenaMatchGameObject `json:"gameObjects"`
	TurnInfo    *TurnInfo              `json:"turnInfo"`
	Players     []PlayerState          `json:"players"`
}

// PlayerState see log
type PlayerState struct {
	LifeTotal        int `json:"lifeTotal"`
	SystemSeatNumber int `json:"systemSeatNumber"`
	TeamID           int `json:"teamId"`
	ControllerSeatID int `json:"controllerSeatId"`
}

// TurnInfo see log
type TurnInfo struct {
	Phase        string `json:"phase"`
	Step         string `json:"step"`
	TurnNumber   int    `json:"turnNumber"`
	ActivePlayer int    `json:"activePlayer"`
}

// ArenaMatchGameObject is a game object in a match
type ArenaMatchGameObject struct {
	InstanceID  int    `json:"instanceId"`
	GrpID       int    `json:"grpId"`
	Type        string `json:"type"`
	ZoneID      int    `json:"zoneId"`
	Visibility  string `json:"visibility"`
	OwnerSeatID int    `json:"ownerSeatId"`
}

// Hash returns a unique string for this object
func (a ArenaMatchGameObject) Hash() string {
	return fmt.Sprintf("%d", a.GrpID)
}

// ArenaMatchEndParams are the params which hold the results of the match
type ArenaMatchEndParams struct {
	PayloadObject *ArenaMatch `json:"payloadObject"`
}

// ArenaMatchEnd is the outer structure
type ArenaMatchEnd struct {
	Params *ArenaMatchEndParams `json:"params"`
}

// IsMatchStart does this segment contain match start
func (s *Segment) IsMatchStart() bool {
	return s.SegmentType == MatchStart
}

// IsMatchEnd does this segment contain a match end
func (s *Segment) IsMatchEnd() bool {
	return s.SegmentType == MatchEnd
}

// ParseMatchStart parses out the match start (will return an incomplete
// ArenaMatch object)
func (s *Segment) ParseMatchStart() (*ArenaMatch, error) {
	var match ArenaMatch
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &match)
	if s.Time != nil {
		match.GameStart = *s.Time
	}
	return &match, err
}

// ParseMatchEnd parses the match end. Contains the match ID
func (s *Segment) ParseMatchEnd() (*ArenaMatchEnd, error) {
	var match ArenaMatchEnd
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &match)
	return &match, err
}

// IsMatchEvent checks if this segment contains anything interesting
// about a currently parsing match
func (s *Segment) IsMatchEvent() bool {
	return s.SegmentType == MatchEvent
}

// ParseMatchEvent looks through the match segments and pulls out
// cards played by whom
func (s *Segment) ParseMatchEvent() (*ArenaMatchEvent, error) {
	// grpID, type
	// ZONEID => match these. 35, hand?
	// ownerSeatID who's info
	// I think we want:
	// greToClientEvent.greToClientMessages[0].gameStateMessage.type ==
	// "type": "GameStateType_Diff",
	// turnInfo.turnNumber
	// zoneID: int
	// type: zoneType
	// GAME OBJECTS
	//payload.PerformActionResp.Actions[].GrpID
	// seat 2 played
	var event ArenaMatchEvent
	err := json.Unmarshal([]byte(stripNonJSON(s.Text)), &event)
	return &event, err
}
