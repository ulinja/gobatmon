package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const pathToCapacityFile = "/sys/class/power_supply/BAT0/capacity"

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
	currentCapacityPercent, err := strconv.ParseUint(text, 10, 64)
	exitOnError(err, "Failed to parse battery capacity file")

	return currentCapacityPercent
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
	defaultPollRate                           = 30
	defaultNormalWarningPercentageThreshold   = 20
	defaultCriticalWarningPercentageThreshold = 10
	defaultNormalWarningReminderTimeout       = 15 * 60
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

// runCapacityCheck executes a battery capacity check, reading the current capacity from the filesystem,
// comparing it to the configured threshold values and triggering a notification if appropriate.
func runCapacityCheck(config *runtimeConfig) {
	// TODO: return the next scheduled run's timestamp as a value

	capacity := uint(readCapacity())
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

func main() {
	assertDependenciesInstalled()

	// Just use defaults for now
	config := runtimeConfig{
		pollRate:                           defaultPollRate,
		normalWarningPercentageThreshold:   defaultNormalWarningPercentageThreshold,
		criticalWarningPercentageThreshold: defaultCriticalWarningPercentageThreshold,
		normalWarningReminderTimeout:       defaultNormalWarningReminderTimeout,
		criticalWarningReminderTimeout:     defaultCriticalWarningPercentageThreshold,
	}
	config = cleanRuntimeConfig(config)

	runCapacityCheck(&config)
}
