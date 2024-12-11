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
}

const (
	DefaultPollRate                           = 60
	DefaultNormalWarningPercentageThreshold   = 20
	DefaultCriticalWarningPercentageThreshold = 10
	DefaultNormalWarningReminderTimeout       = 10 * 60
	DefaultCriticalWarningReminderTimeout     = 5 * 60
)

func logValidationWarning(msg string) {
	log.Printf("WARNING: %s. Using default value(s) instead.\n", msg)
}

// cleanRuntimeConfig validates, cleans and returns the provided runtime configuration.
// If a validation error is encountered in the provided configuration, an error message is
// logged and the affected configuration option is set to its default value as a fallback.
func cleanRuntimeConfig(config RuntimeConfig) (cleanedConfig RuntimeConfig) {

	cleanedConfig = config

	if cleanedConfig.PollRate < 1 {
		logValidationWarning("Poll Rate cannot be smaller than 1s")
		cleanedConfig.PollRate = DefaultPollRate
	}

	if t := cleanedConfig.NormalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Normal Warning Percentage Threshold cannot exceed 100%")
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Normal Warning Percentage Threshold cannot be less than 0%")
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
	}

	if t := cleanedConfig.CriticalWarningPercentageThreshold; t > 100 {
		logValidationWarning("Critical Warning Percentage Threshold cannot exceed 100%")
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	} else if t < 0 {
		logValidationWarning("Critical Warning Percentage Threshold cannot be less than 0%")
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.CriticalWarningPercentageThreshold >= cleanedConfig.NormalWarningPercentageThreshold {
		logValidationWarning("Critical Warning Percentage Threshold must be smaller than Normal Warning Percentage Threshold")
		cleanedConfig.NormalWarningPercentageThreshold = DefaultNormalWarningPercentageThreshold
		cleanedConfig.CriticalWarningPercentageThreshold = DefaultCriticalWarningPercentageThreshold
	}

	if cleanedConfig.NormalWarningReminderTimeout < 1 {
		logValidationWarning("Normal Warning Reminder Timeout cannot be less than 1s")
		cleanedConfig.NormalWarningReminderTimeout = DefaultNormalWarningReminderTimeout
	}
	if cleanedConfig.CriticalWarningReminderTimeout < 1 {
		logValidationWarning("Critical Warning Reminder Timeout cannot be less than 1s")
		cleanedConfig.CriticalWarningReminderTimeout = DefaultCriticalWarningReminderTimeout
	}

	return
}

// GetRuntimeConfig returns the active runtime confiuration.
func GetRuntimeConfig() (runtimeConfig RuntimeConfig) {
	// Just use defaults for now
	runtimeConfig = RuntimeConfig{
		PollRate:                           DefaultPollRate,
		NormalWarningPercentageThreshold:   DefaultNormalWarningPercentageThreshold,
		CriticalWarningPercentageThreshold: DefaultCriticalWarningPercentageThreshold,
		NormalWarningReminderTimeout:       DefaultNormalWarningReminderTimeout,
		CriticalWarningReminderTimeout:     DefaultCriticalWarningReminderTimeout,
	}
	runtimeConfig = cleanRuntimeConfig(runtimeConfig)
	return
}
