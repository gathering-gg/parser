package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMatchStart(t *testing.T) {
	s := &Segment{
		SegmentType: MatchStart,
	}
	assert.True(t, s.IsMatchStart())
}

func TestIsMatchEnd(t *testing.T) {
	s := &Segment{
		SegmentType: MatchEnd,
	}
	assert.True(t, s.IsMatchEnd())
}
