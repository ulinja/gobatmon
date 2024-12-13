package cliargs

import (
	"flag"
	"fmt"

	"github.com/ulinja/gobatmon/internal/config"
)

func ParseRuntimeConfig() (runtimeConfig config.RuntimeConfig, showHelp bool, showVersion bool) {
	var disableIcons bool

	flag.UintVar(
		&runtimeConfig.PollRate,
		"poll-rate",
		config.DefaultPollRate,
		"Poll rate for battery status in seconds",
	)
	flag.UintVar(
		&runtimeConfig.NormalWarningPercentageThreshold,
		"normal-warning-threshold",
		config.DefaultNormalWarningPercentageThreshold,
		"Threshold percentage below which a normal low battery warning is triggered",
	)
	flag.UintVar(
		&runtimeConfig.CriticalWarningPercentageThreshold,
		"critical-warning-threshold",
		config.DefaultCriticalWarningPercentageThreshold,
		"Threshold percentage below which a critical low battery warning is triggered",
	)
	flag.UintVar(
		&runtimeConfig.NormalWarningReminderTimeout,
		"normal-warning-reminder-timeout",
		config.DefaultNormalWarningReminderTimeout,
		"Timeout in seconds after which a normal low battery warning is repeated",
	)
	flag.UintVar(
		&runtimeConfig.CriticalWarningReminderTimeout,
		"critical-warning-reminder-timeout",
		config.DefaultCriticalWarningReminderTimeout,
		"Timeout in seconds after which a critical low battery warning is repeated",
	)
	flag.StringVar(
		&runtimeConfig.NormalWarningIconName,
		"normal-warning-icon-name",
		config.DefaultNormalWarningIconName,
		"Name of the icon to use for normal low battery warning notifications",
	)
	flag.StringVar(
		&runtimeConfig.CriticalWarningIconName,
		"critical-warning-icon-name",
		config.DefaultCriticalWarningIconName,
		"Name of the icon to use for critical low battery warning notifications",
	)
	flag.BoolVar(
		&disableIcons,
		"disable-icons",
		false,
		"Do not show icons in warning notifications",
	)

	flag.BoolVar(
		&showHelp,
		"help",
		false,
		"Show help message and exit",
	)
	flag.BoolVar(
		&showVersion,
		"version",
		false,
		"Show version information and exit",
	)

	flag.Parse()
	runtimeConfig.EnableIcons = !disableIcons
	runtimeConfig = config.CleanRuntimeConfig(runtimeConfig, true)

	return
}

func PrintHelp() {
	fmt.Printf("Usage: gobatmon [OPTIONS]\n\nOptions:\n")
	flag.PrintDefaults()
}
