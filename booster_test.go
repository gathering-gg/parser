package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBooster(t *testing.T) {
	s := &Segment{
		SegmentType: CrackBooster,
	}
	assert.True(t, s.IsCrackBooster())
}
