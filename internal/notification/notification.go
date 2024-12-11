package notification

import (
	"os/exec"

	"github.com/ulinja/gobatmon/internal/util"
)

type UrgencyLevel uint

const (
	_ = iota
	UrgencyLevelLow
	UrgencyLevelNormal
	UrgencyLevelCritical
)

type Notification struct {
	Urgency UrgencyLevel
	Summary string
	Body    string
}

// ShowNotification invokes notify-send to display a desktop notification.
func ShowNotification(n Notification) {
	var urgency string
	switch level := n.Urgency; level {
	case UrgencyLevelLow:
		urgency = "--urgency=low"
	case UrgencyLevelCritical:
		urgency = "--urgency=critical"
	default:
		urgency = "--urgency=normal"
	}

	cmd := exec.Command("notify-send", urgency, n.Summary, n.Body)
	err := cmd.Run()
	util.ExitOnError(err, "Invocation of notify-send failed")
}
