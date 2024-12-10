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

func main() {
	assertDependenciesInstalled()

	capacity := readCapacity()

	n := notification{
		summary: "Current Battery Capacity",
		body:    fmt.Sprintf("%d%%", capacity),
	}
	showNotification(n)
}
