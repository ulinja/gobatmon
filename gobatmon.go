package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ulinja/gobatmon/internal/battery"
	"github.com/ulinja/gobatmon/internal/cliargs"
	"github.com/ulinja/gobatmon/internal/config"
	"github.com/ulinja/gobatmon/internal/notification"
	"github.com/ulinja/gobatmon/internal/util"
)

const version = "0.3.0"

// runCheck executes a battery capacity and state check, reading these from the filesystem,
// comparing them to the configured threshold values, and triggering a notification if appropriate.
func runCheck(runtimeConfig *config.RuntimeConfig) (capacity uint, batteryState battery.BatteryState) {
	capacity = uint(battery.ReadCapacity())
	batteryState = battery.ReadBatteryState()

	if batteryState == battery.BatteryStateDischarging {
		var ntf notification.Notification
		if capacity <= runtimeConfig.CriticalWarningPercentageThreshold {
			ntf.Summary = "Very Low Battery"
			ntf.Body = fmt.Sprintf("Charge is at %d%%. Plug in AC power now!", capacity)
			ntf.Urgency = notification.UrgencyCritical
			ntf.Send()
		} else if capacity <= runtimeConfig.NormalWarningPercentageThreshold {
			ntf.Summary = "Low Battery"
			ntf.Body = fmt.Sprintf("Charge is at %d%%.", capacity)
			ntf.Urgency = notification.UrgencyNormal
			ntf.Send()
		}
	}

	return
}

type programState struct {
	batteryState battery.BatteryState
	capacity     uint
	lastCheck    time.Time
}

// durationUntilNextCheck calculates and returns the duration until the next battery state check, which depends
// on the current charging status and the configured timeouts and polling rate.
func durationUntilNextCheck(config *config.RuntimeConfig, programState *programState) (duration time.Duration) {
	switch programState.batteryState {
	case battery.BatteryStateDischarging:
		break
	default:
		return time.Duration(config.PollRate) * time.Second
	}

	durationSinceLastChecked := time.Now().Sub(programState.lastCheck)
	if c := programState.capacity; c > config.NormalWarningPercentageThreshold {
		duration = time.Duration(config.PollRate) * time.Second
	} else if c > config.CriticalWarningPercentageThreshold {
		duration = time.Duration(config.NormalWarningReminderTimeout)*time.Second - durationSinceLastChecked
	} else {
		duration = time.Duration(config.CriticalWarningReminderTimeout)*time.Second - durationSinceLastChecked
	}

	if duration < 0 {
		duration = 0
	}

	return
}

func main() {
	runtimeConfig, showHelpAndExit, showVersionAndExit := cliargs.ParseRuntimeConfig()
	if showHelpAndExit {
		cliargs.PrintHelp()
		os.Exit(0)
	}
	if showVersionAndExit {
		fmt.Printf("gobatmon v%s\n", version)
		os.Exit(0)
	}

	log.Printf("starting gobatmon v%s\n", version)
	util.AssertDbusNotificationsAvailable()
	var programState programState

	for {
		capacity, batteryState := runCheck(&runtimeConfig)
		programState.capacity = capacity
		programState.batteryState = batteryState
		programState.lastCheck = time.Now()

		durationUntilNextCheck := durationUntilNextCheck(&runtimeConfig, &programState)
		time.Sleep(durationUntilNextCheck)
	}
}
