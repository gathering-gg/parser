package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCollection(t *testing.T) {
	s := &Segment{
		SegmentType: PlayerInventoryGetPlayerCards,
	}
	assert.True(t, s.IsCollection())
}
