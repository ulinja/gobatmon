package battery

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ulinja/gobatmon/internal/util"
)

const pathToCapacityFile = "/sys/class/power_supply/BAT0/capacity"
const pathToBatteryStateFile = "/sys/class/power_supply/BAT0/status"

// ReadCapacity reads the battery capacity file and returns the current charge level in percent as a uint64.
func ReadCapacity() (capacity uint64) {
	bytes, err := os.ReadFile(pathToCapacityFile)
	util.ExitOnError(err, "Failed to read battery capacity file")

	text := strings.TrimSuffix(string(bytes), "\n")
	capacity, err = strconv.ParseUint(text, 10, 64)
	util.ExitOnError(err, "Failed to parse battery capacity file")

	if capacity > 100 {
		log.Printf("WARNING: system reported battery capacity greater than 100 at %s: %d\n", pathToCapacityFile, capacity)
		capacity = 100
	} else if capacity < 0 {
		util.ExitOnError(fmt.Errorf("capacity: %d", capacity), "Battery capacity is negative")
	}

	return
}

// BatteryState describes one of three possible battery states: discharging, charging or full.
type BatteryState uint

const (
	BatteryStateDischarging = iota
	BatteryStateCharging
	BatteryStateFull
)

// ReadBatteryState reads the battery state file and returns the current charge state.
func ReadBatteryState() (batteryState BatteryState) {
	bytes, err := os.ReadFile(pathToBatteryStateFile)
	util.ExitOnError(err, "Failed to read battery state file")

	text := strings.TrimSuffix(string(bytes), "\n")
	switch text {
	case "Discharging":
		batteryState = BatteryStateDischarging
	case "Charging":
		batteryState = BatteryStateCharging
	case "Full":
		batteryState = BatteryStateFull
	default:
		err = fmt.Errorf("unhandled status value: '%s'", text)
		util.ExitOnError(err, "Failed to parse battery state from file")
	}

	return
}
