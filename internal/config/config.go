package config

import (
	"log"
)

type RuntimeConfig struct {
	PollRate                           uint
	NormalWarningPercentageThreshold   uint
	CriticalWarningPercentageThreshold uint
	NormalWarningReminderTimeout       uint
	CriticalWarningReminderTimeout     uint
	NormalWarningIconName              string
	CriticalWarningIconName            string
	EnableIcons                        bool
}

const (
	DefaultPollRate                           = 60
	DefaultNormalWarningPercentageThreshold   = 20
	DefaultCriticalWarningPercentageThreshold = 10
	DefaultNormalWarningReminderTimeout       = 10 * 60
	DefaultCriticalWarningReminderTimeout     = 5 * 60
	DefaultNormalWarningIconName              = "battery-low"
	DefaultCriticalWarningIconName            = "battery-caution"
	DefaultEnableIcons                        = true
)

func logValidationWarning(msg string, exitOnError bool) {
	if exitOnError {
		log.Fatalf("ERROR: %s. Exiting.\n", msg)
	} else {
		log.Printf("WARNING: %s. Using default value(s) instead.\n", msg)
	}
}

// CleanRuntimeConfig validates, cleans and returns the provided runtime configuration.
// If a validation error is encountered in the provided configuration and exitOnError is false,
// an error message is logged and the affected configuration option is set to its default value
// as a fallback.
// If exitOnError is true, the program terminates with an error message when a validation error
// is encountered.
func CleanRuntimeConfig(config RuntimeConfig, exitOnError bool) (cleanedConfig RuntimeConfig) {

	cleanedConfig = config

	if cleanedConfig.PollRate < 1 {
		logValidationWarning("Poll Rate cannot be smaller than 1s", exitOnError)
		cleanedConfig.PollRate = DefaultPollRate
	}

	if t := cleanedConfig.NormalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Normal Warning Percentage Threshold cannot exceed 100%", exitOnError)
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Normal Warning Percentage Threshold cannot be less than 0%", exitOnError)
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
	}

	if t := cleanedConfig.CriticalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Critical Warning Percentage Threshold cannot exceed 100%", exitOnError)
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Critical Warning Percentage Threshold cannot be less than 0%", exitOnError)
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.CriticalWarningPercentageThreshold >= cleanedConfig.NormalWarningPercentageThreshold {
		logValidationWarning("Critical Warning Percentage Threshold must be smaller than Normal Warning Percentage Threshold", exitOnError)
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.NormalWarningReminderTimeout < 1 {
		logValidationWarning("Normal Warning Reminder Timeout cannot be less than 1s", exitOnError)
		cleanedConfig.NormalWarningReminderTimeout = DefaultNormalWarningReminderTimeout
	}
	if cleanedConfig.CriticalWarningReminderTimeout < 1 {
		logValidationWarning("Critical Warning Reminder Timeout cannot be less than 1s", exitOnError)
		cleanedConfig.CriticalWarningReminderTimeout = DefaultCriticalWarningReminderTimeout
	}

	if len(cleanedConfig.NormalWarningIconName) == 0 {
		logValidationWarning("Normal Warning Icon Name cannot be empty", exitOnError)
		cleanedConfig.NormalWarningIconName = DefaultNormalWarningIconName
	}
	if len(cleanedConfig.CriticalWarningIconName) == 0 {
		logValidationWarning("Critical Warning Icon Name cannot be empty", exitOnError)
		cleanedConfig.CriticalWarningIconName = DefaultCriticalWarningIconName
	}

	return
}

// GetRuntimeConfig returns a default runtime confiuration.
func GetDefaultRuntimeConfig() (runtimeConfig RuntimeConfig) {
	return RuntimeConfig{
		PollRate:                           DefaultPollRate,
		NormalWarningPercentageThreshold:   DefaultNormalWarningPercentageThreshold,
		CriticalWarningPercentageThreshold: DefaultCriticalWarningPercentageThreshold,
		NormalWarningReminderTimeout:       DefaultNormalWarningReminderTimeout,
		CriticalWarningReminderTimeout:     DefaultCriticalWarningReminderTimeout,
		NormalWarningIconName:              DefaultNormalWarningIconName,
		CriticalWarningIconName:            DefaultCriticalWarningIconName,
	}
}
