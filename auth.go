package gathering

// ArenaAuthRequestPayload is the payload the Arena client sends when
// authenticating. We only are interested in the Player's name
type ArenaAuthRequestPayload struct {
	PlayerName string `json:"PlayerName"`
}

// ArenaAuthRequest is the base structure which wraps a payload
type ArenaAuthRequest struct {
	Payload ArenaAuthRequestPayload `json:"Payload"`
}

// IsPlayerAuth checks if a segment contains an auth statement
func (s *Segment) IsPlayerAuth() bool {
	return s.SegmentType == PlayerAuth
}

// ParseAuth returns the players username
func (s *Segment) ParseAuth() ([]byte, error) {
	re := segmentTypeChecks[PlayerAuth].Copy()
	matches := re.FindSubmatch(s.Text)
	if len(matches) == 2 {
		return matches[1], nil
	}
	return nil, ErrNotFound
}
