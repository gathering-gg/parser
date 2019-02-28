package gathering

import (
	"encoding/json"
	"fmt"
	"time"
)

// ArenaMatch is a match. A match contains at least one game, but in best of
// three, main contain up to 3.
type ArenaMatch struct {
	currentGame                    int
	MatchID                        string       `json:"matchId"`
	Games                          []*ArenaGame `json:"games"`
	GameStart                      *time.Time   `json:"gameStart"`
	EventID                        string       `json:"eventId"`
	OpponentScreenName             string       `json:"opponentScreenName"`
	OpponentIsWotc                 bool         `json:"opponentIsWotc"`
	OpponentRankingClass           string       `json:"opponentRankingClass"`
	OpponentRankingTier            int          `json:"opponentRankingTier"`
	OpponentMythicPercentile       float64      `json:"opponentMythicPercentile"`
	OpponentMythicLeaderboardPlace int          `json:"opponentMythicLeaderboardPlace"`
	CourseDeck                     *ArenaDeck   `json:"CourseDeck"`
}

// ArenaGame is a game within a match
type ArenaGame struct {
	GameStart     *time.Time                     `json:"gameStart"`
	Number        *int                           `json:"number"`
	MatchID       *string                        `json:"matchId"`
	SeatID        *int                           `json:"seatId"`
	TeamID        *int                           `json:"teamId"`
	WinningTeamID *int                           `json:"winningTeamId"`
	WinningReason *string                        `json:"winningReason"`
	TurnCount     *int                           `json:"turnCount"`
	SecondsCount  *int                           `json:"secondsCount"`
	CourseDeck    *ArenaDeck                     `json:"CourseDeck"`
	SeenObjects   map[int][]ArenaMatchGameObject `json:"seenObjects"`
}

// UpdateGameEnd updates the latest game with the game result
func (a *ArenaMatch) UpdateGameEnd(end *ArenaGame) {
	game := a.Games[a.currentGame]
	num := a.currentGame + 1
	game.MatchID = &a.MatchID
	game.SeatID = end.SeatID
	game.TeamID = end.TeamID
	game.Number = &num
	game.WinningTeamID = end.WinningTeamID
	game.WinningReason = end.WinningReason
	game.TurnCount = end.TurnCount
	game.SecondsCount = end.SecondsCount
}

// UpdateMatchCompleted updates the match object with the completed
// status
func (a *ArenaMatch) UpdateMatchCompleted(com *ArenaMatchCompleted) {
	// TODO: I'm not sure we need anything here, but maybe some sanity setters?
}

// LogMatchEvent adds an event to the log
func (a *ArenaMatch) LogMatchEvent(event *ArenaMatchEvent) {
	game := a.Games[len(a.Games)-1]
	if game.SeenObjects == nil {
		game.SeenObjects = make(map[int][]ArenaMatchGameObject)
	}
	for _, m := range event.GreToClientEvent.GreToClientMessages {
		gsm := m.GameStateMessage
		for _, o := range gsm.GameObjects {
			game.SeenObjects[o.OwnerSeatID] = append(game.SeenObjects[o.OwnerSeatID], o)
		}
	}
	for k, v := range game.SeenObjects {
		uniq := make(map[string]bool)
		var objects []ArenaMatchGameObject
		for _, o := range v {
			if _, ok := uniq[o.Hash()]; !ok && o.Type == "GameObjectType_Card" {
				uniq[o.Hash()] = true
				objects = append(objects, o)
			}
		}
		game.SeenObjects[k] = objects
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
	GameInfo    *GameInfo              `json:"gameInfo"`
}

// GameInfo contains match info, such as which game this is
type GameInfo struct {
	MatchID            string        `json:"matchID"`
	GameNumber         int           `json:"gameNumber"`
	Stage              string        `json:"stage"`
	Type               string        `json:"type"`
	MatchState         string        `json:"matchState"`
	MatchWinCondition  string        `json:"matchWinCondition"`
	MaxTimeoutCount    int           `json:"maxTimeoutCount"`
	MaxPipCount        int           `json:"maxPipCount"`
	TimeoutDurationSec int           `json:"timeoutDurationSec"`
	Results            []MatchResult `json:"results"`
	SuperFormat        string        `json:"superFormat"`
	MulliganType       string        `json:"mulliganType"`
}

// MatchResult has a list of who won the games in a match
type MatchResult struct {
	Scope         string `json:" scope"`
	Result        string `json:"result"`
	WinningTeamID int    `json:"winningTeamId"`
	Reason        string `json:"reason"`
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

// Start ArenaMatchCompleted

// ArenaMatchCompleted is when a match (and all games) is done
type ArenaMatchCompleted struct {
	TransactionID                  string                         `json:"transactionId"`
	Timestamp                      string                         `json:"timestamp"`
	MatchGameRoomStateChangedEvent MatchGameRoomStateChangedEvent `json:"matchGameRoomStateChangedEvent"`
}

// MatchGameRoomStateChangedEvent has info
type MatchGameRoomStateChangedEvent struct {
	GameRoomInfo MatchGameRoomInfo `json:"gameRoomInfo"`
}

// MatchGameRoomInfo contains the config
type MatchGameRoomInfo struct {
	GameRoomConfig   MatchGameRoomConfig   `json:"gameRoomConfig"`
	StateType        string                `json:"stateType"`
	FinalMatchResult MatchFinalMatchResult `json:"finalMatchResult"`
}

// MatchGameRoomConfig contains info
type MatchGameRoomConfig struct {
	EventID string `json:"eventId"`
	MatchID string `json:"matchId"`
}

// MatchFinalMatchResult contains the final results
type MatchFinalMatchResult struct {
	MatchID              string                 `json:"matchId"`
	MatchCompletedReason string                 `json:"matchCompletedReason"`
	ResultList           []MatchCompletedResult `json:"resultList"`
}

// MatchCompletedResult contains the match results
type MatchCompletedResult struct {
	Scope         string `json:"scope"`
	Result        string `json:"result"`
	WinningTeamID int    `json:"winningTeamId"`
}

/*********************** End ArenaMatchCompleted ***************************/

// Hash returns a unique string for this object
func (a ArenaMatchGameObject) Hash() string {
	return fmt.Sprintf("%d", a.GrpID)
}

// ArenaMatchEndParams are the params which hold the results of the match
type ArenaMatchEndParams struct {
	PayloadObject *ArenaGame `json:"payloadObject"`
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

// IsMatchCompleted is a better metric for a match (including ALL games) being
// completed.
func (s *Segment) IsMatchCompleted() bool {
	return s.SegmentType == MatchCompleted
}

// ParseMatchStart parses out the match start (will return an incomplete
// ArenaMatch object)
func (s *Segment) ParseMatchStart() (*ArenaMatch, error) {
	var match ArenaMatch
	err := json.Unmarshal(stripNonJSON(s.Text), &match)
	if s.Time != nil {
		match.GameStart = s.Time
	}
	return &match, err
}

// ParseMatchEnd parses the match end. Contains the match ID
func (s *Segment) ParseMatchEnd() (*ArenaMatchEnd, error) {
	var match ArenaMatchEnd
	err := json.Unmarshal(stripNonJSON(s.Text), &match)
	return &match, err
}

// ParseMatchCompleted parses the match completed. It means all games in a match
// are finished.
func (s *Segment) ParseMatchCompleted() (*ArenaMatchCompleted, error) {
	var done ArenaMatchCompleted
	err := json.Unmarshal(stripNonJSON(s.Text), &done)
	return &done, err
}

// IsMatchEvent checks if this segment contains anything interesting
// about a currently parsing match
func (s *Segment) IsMatchEvent() bool {
	return s.SegmentType == MatchEvent
}

// IsSideboardStop checks if this is a sideboard end event
func (s *Segment) IsSideboardStop() bool {
	return s.SegmentType == DuelSceneSideboardingStop
}

// ParseMatchEvent looks through the match segments and pulls out
// cards played by whom
func (s *Segment) ParseMatchEvent() (*ArenaMatchEvent, error) {
	var event ArenaMatchEvent
	err := json.Unmarshal(stripNonJSON(s.Text), &event)
	return &event, err
}
