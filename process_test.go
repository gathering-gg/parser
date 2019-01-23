package gathering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArenaRunning(t *testing.T) {
	a := assert.New(t)
	running, err := IsArenaRunning()
	a.Nil(err)
	a.False(running)
}
