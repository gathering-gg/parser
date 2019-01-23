package gathering

import (
	"regexp"

	"github.com/shirou/gopsutil/process"
)

var mtgArenaName = regexp.MustCompile(`(?mi)mtga\.exe`)

// IsArenaRunning checks if Arena is running
func IsArenaRunning() (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		name, _ := p.Name()
		if mtgArenaName.MatchString(name) {
			return true, nil
		}
	}
	return false, nil
}
