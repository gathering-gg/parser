package gathering

// UploadData encapsulates the data to send to the server
type UploadData struct {
	IsPlaying  bool                  `json:"isPlaying"`
	Collection map[string]int        `json:"collection"`
	Decks      []ArenaDeck           `json:"deck"`
	Inventory  *ArenaPlayerInventory `json:"inventory"`
	Rank       *ArenaRankInfo        `json:"rank"`
	Auth       *ArenaAuthRequest     `json:"auth"`
	Matches    []*ArenaMatch         `json:"matches"`
}
