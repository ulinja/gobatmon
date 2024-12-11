package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ulinja/gobatmon/internal/battery"
	"github.com/ulinja/gobatmon/internal/config"
	"github.com/ulinja/gobatmon/internal/notification"
	"github.com/ulinja/gobatmon/internal/util"
)

const version = "0.1.0"

// runCheck executes a battery capacity and state check, reading these from the filesystem,
// comparing them to the configured threshold values, and triggering a notification if appropriate.
func runCheck(runtimeConfig *config.RuntimeConfig) (capacity uint, batteryState battery.BatteryState) {
	capacity = uint(battery.ReadCapacity())
	batteryState = battery.ReadBatteryState()

	if batteryState == battery.BatteryStateDischarging {
		if capacity <= runtimeConfig.CriticalWarningPercentageThreshold {
			notification.ShowNotification(notification.Notification{
				Urgency: notification.UrgencyLevelCritical,
				Summary: "Very Low Battery",
				Body:    fmt.Sprintf("Charge is at %d%%.", capacity),
			})
		} else if capacity <= runtimeConfig.NormalWarningPercentageThreshold {
			notification.ShowNotification(notification.Notification{
				Urgency: notification.UrgencyLevelNormal,
				Summary: "Low Battery",
				Body:    fmt.Sprintf("Charge is at %d%%.", capacity),
			})
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
	log.Printf("starting gobatmon v%s\n", version)
	util.AssertDependenciesInstalled()
	runtimeConfig := config.GetRuntimeConfig()
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
