package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const pathToCapacityFile = "/sys/class/power_supply/BAT0/capacity"
const pathToBatteryStateFile = "/sys/class/power_supply/BAT0/status"

// exitOnError logs a fatal error message and exits the program with the provided prefix if error is not nil.
func exitOnError(err error, prefix string) {
	if err == nil {
		return
	}

	var msg string
	if len(prefix) > 0 {
		msg = fmt.Sprintf("%s: %s", prefix, err)
	} else {
		msg = fmt.Sprintf("%s", err)
	}
	log.Fatalln(msg)
}

// readCapacity reads the battery capacity file and returns the current charge level in percent as a uint64.
func readCapacity() (capacity uint64) {
	bytes, err := os.ReadFile(pathToCapacityFile)
	exitOnError(err, "Failed to read battery capacity file")

	text := strings.TrimSuffix(string(bytes), "\n")
	capacity, err = strconv.ParseUint(text, 10, 64)
	exitOnError(err, "Failed to parse battery capacity file")

	if capacity > 100 {
		log.Printf("WARNING: system reported battery capacity greater than 100 at %s: %d\n", pathToCapacityFile, capacity)
		capacity = 100
	} else if capacity < 0 {
		exitOnError(fmt.Errorf("capacity: %d", capacity), "Battery capacity is negative")
	}

	return
}

type batteryState uint

const (
	batteryStateDischarging = iota
	batteryStateCharging
	batteryStateFull
)

// readBatteryState reads the battery state file and returns the current charge state.
func readBatteryState() (batteryState batteryState) {
	bytes, err := os.ReadFile(pathToBatteryStateFile)
	exitOnError(err, "Failed to read battery state file")

	text := strings.TrimSuffix(string(bytes), "\n")
	switch text {
	case "Discharging":
		batteryState = batteryStateDischarging
	case "Charging":
		batteryState = batteryStateCharging
	case "Full":
		batteryState = batteryStateFull
	default:
		err = fmt.Errorf("unhandled status value: '%s'", text)
		exitOnError(err, "Failed to parse battery state from file")
	}

	return
}

// assertDependenciesInstalled checks if all executable dependencies are present on $PATH, and terminates
// the program with an error message if not.
func assertDependenciesInstalled() {
	executables := []string{
		"notify-send",
	}
	for _, e := range executables {
		_, err := exec.LookPath(e)
		exitOnError(err, "Missing dependency")
	}
}

type urgencyLevel uint

const (
	_ = iota
	urgencyLevelLow
	urgencyLevelNormal
	urgencyLevelCritical
)

type notification struct {
	urgency urgencyLevel
	summary string
	body    string
}

// showNotification invokes notify-send to display a desktop notification.
func showNotification(n notification) {
	var urgency string
	switch level := n.urgency; level {
	case urgencyLevelLow:
		urgency = "--urgency=low"
	case urgencyLevelCritical:
		urgency = "--urgency=critical"
	default:
		urgency = "--urgency=normal"
	}

	cmd := exec.Command("notify-send", urgency, n.summary, n.body)
	err := cmd.Run()
	exitOnError(err, "Invocation of notify-send failed")
}

type runtimeConfig struct {
	pollRate                           uint
	normalWarningPercentageThreshold   uint
	criticalWarningPercentageThreshold uint
	normalWarningReminderTimeout       uint
	criticalWarningReminderTimeout     uint
}

const (
	defaultPollRate                           = 60
	defaultNormalWarningPercentageThreshold   = 20
	defaultCriticalWarningPercentageThreshold = 10
	defaultNormalWarningReminderTimeout       = 10 * 60
	defaultCriticalWarningReminderTimeout     = 5 * 60
)

func logValidationWarning(msg string) {
	log.Printf("WARNING: %s. Using default value(s) instead.\n", msg)
}

// cleanRuntimeConfig validates, cleans and returns the provided runtime configuration.
// If a validation error is encountered in the provided configuration, an error message is
// logged and the affected configuration option is set to its default value as a fallback.
func cleanRuntimeConfig(config runtimeConfig) (cleanedConfig runtimeConfig) {

	cleanedConfig = config

	if cleanedConfig.pollRate < 1 {
		logValidationWarning("Poll Rate cannot be smaller than 1s")
		cleanedConfig.pollRate = defaultPollRate
	}

	if t := cleanedConfig.normalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Normal Warning Percentage Threshold cannot exceed 100%")
		cleanedConfig.normalWarningPercentageThreshold = defaultNormalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Normal Warning Percentage Threshold cannot be less than 0%")
		cleanedConfig.normalWarningPercentageThreshold = defaultNormalWarningPercentageThreshold
	}

	if t := cleanedConfig.criticalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Critical Warning Percentage Threshold cannot exceed 100%")
		cleanedConfig.criticalWarningPercentageThreshold = defaultCriticalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Critical Warning Percentage Threshold cannot be less than 0%")
		cleanedConfig.criticalWarningPercentageThreshold = defaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.criticalWarningPercentageThreshold >= cleanedConfig.normalWarningPercentageThreshold {
		logValidationWarning("Critical Warning Percentage Threshold must be smaller than Normal Warning Percentage Threshold")
		cleanedConfig.normalWarningPercentageThreshold = defaultNormalWarningPercentageThreshold
		cleanedConfig.criticalWarningPercentageThreshold = defaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.normalWarningReminderTimeout < 1 {
		logValidationWarning("Normal Warning Reminder Timeout cannot be less than 1s")
		cleanedConfig.normalWarningReminderTimeout = defaultNormalWarningReminderTimeout
	}
	if cleanedConfig.criticalWarningReminderTimeout < 1 {
		logValidationWarning("Critical Warning Reminder Timeout cannot be less than 1s")
		cleanedConfig.criticalWarningReminderTimeout = defaultCriticalWarningReminderTimeout
	}

	return
}

// runCheck executes a battery capacity check, reading the current capacity from the filesystem,
// comparing it to the configured threshold values and triggering a notification if appropriate.
func runCheck(config *runtimeConfig) (capacity uint, batteryState batteryState) {
	capacity = uint(readCapacity())
	batteryState = readBatteryState()

	if batteryState == batteryStateDischarging {
		if capacity <= config.criticalWarningPercentageThreshold {
			showNotification(notification{
				urgency: urgencyLevelCritical,
				summary: "Very Low Battery",
				body:    fmt.Sprintf("Charge is at %d%%.", capacity),
			})
		} else if capacity <= config.normalWarningPercentageThreshold {
			showNotification(notification{
				urgency: urgencyLevelNormal,
				summary: "Low Battery",
				body:    fmt.Sprintf("Charge is at %d%%.", capacity),
			})
		}
	}

	return
}

type programState struct {
	batteryState batteryState
	capacity     uint
	lastCheck    time.Time
}

// durationUntilNextCheck calculates and returns the duration until the next battery state check, which depends
// on the current charging status and the configured timeouts and polling rate.
func durationUntilNextCheck(config *runtimeConfig, programState *programState) (duration time.Duration) {
	switch programState.batteryState {
	case batteryStateDischarging:
		break
	default:
		return time.Duration(config.pollRate) * time.Second
	}

	durationSinceLastChecked := time.Now().Sub(programState.lastCheck)
	if c := programState.capacity; c > config.normalWarningPercentageThreshold {
		duration = time.Duration(config.pollRate) * time.Second
	} else if c > config.criticalWarningPercentageThreshold {
		duration = time.Duration(config.normalWarningReminderTimeout)*time.Second - durationSinceLastChecked
	} else {
		duration = time.Duration(config.criticalWarningReminderTimeout)*time.Second - durationSinceLastChecked
	}

	if duration < 0 {
		duration = 0
	}

	return
}

func main() {
	assertDependenciesInstalled()

	// Just use defaults for now
	config := runtimeConfig{
		pollRate:                           defaultPollRate,
		normalWarningPercentageThreshold:   defaultNormalWarningPercentageThreshold,
		criticalWarningPercentageThreshold: defaultCriticalWarningPercentageThreshold,
		normalWarningReminderTimeout:       defaultNormalWarningReminderTimeout,
		criticalWarningReminderTimeout:     defaultCriticalWarningReminderTimeout,
	}
	config = cleanRuntimeConfig(config)

	var programState programState

	for {
		capacity, batteryState := runCheck(&config)
		programState.capacity = capacity
		programState.batteryState = batteryState
		programState.lastCheck = time.Now()

		durationUntilNextCheck := durationUntilNextCheck(&config, &programState)
		time.Sleep(durationUntilNextCheck)
	}
}
