package gathering

import (
	"encoding/json"
)

// IsCollection checks if a segment contains the collection
func (s *Segment) IsCollection() bool {
	return s.SegmentType == PlayerInventoryGetPlayerCards
}

// ParseCollection parses a collection from a Segment. It is up to the caller to
// check if this segment contains a Collection with `IsCollection()`
func (s *Segment) ParseCollection() (map[string]int, error) {
	var collection map[string]int
	err := json.Unmarshal(stripNonJSON(s.Text), &collection)
	return collection, err
}
