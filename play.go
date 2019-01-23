package gathering

// ArenaPlay :12494
type ArenaPlay struct {
	MatchID                        string     `json:"matchId"`
	OpponentScreenName             string     `json:"opponentScreenName"`
	OpponentIsWotc                 bool       `json:"opponentIsWotc"`
	OpponentRankingClass           string     `json:"opponentRankingClass"`
	OpponentRankingTier            int        `json:"opponentRankingTier"`
	OpponentMythicPercentile       float64    `json:"opponentMythicPercentile"`
	OpponentMythicLeaderboardPlace int        `json:"opponentMythicLeaderboardPlace"`
	EventID                        string     `json:"eventId"`
	SeatID                         *int       `json:"seatId"`
	TeamID                         *int       `json:"teamId"`
	GameNumber                     *int       `json:"gameNumber"`
	WinningTeamID                  *int       `json:"winningTeamId"`
	WinningReason                  *string    `json:"winningReason"`
	TurnCount                      *int       `json:"turnCount"`
	SecondsCount                   *int       `json:"secondsCount"`
	CourseDeck                     *ArenaDeck `json:"CourseDeck"`
}
