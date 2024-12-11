package dbusnotification

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ulinja/gobatmon/internal/util"
)

const (
	notificationsAppName = "gobatmon"
)

type Urgency uint

const (
	_ = iota
	UrgencyLow
	UrgencyNormal
	UrgencyCritical
)

type Notification struct {
	Id uint32
	//Icon string
	//Sound string
	Summary string
	Body    string
	Urgency Urgency
	SentAt  time.Time
}

func (notification *Notification) Send() {
	conn, err := dbus.ConnectSessionBus()
	util.ExitOnError(err, "Failed to connect to DBUS session bus")
	defer conn.Close()

	dn := dbusNotification{
		AppName:       notificationsAppName,
		ReplacesID:    0,
		AppIcon:       "",
		Summary:       notification.Summary,
		Body:          notification.Body,
		ExpireTimeout: -1,
	}

	switch notification.Urgency {
	case UrgencyLow:
		dn.setUrgency(dbusUrgencyLow)
	case UrgencyCritical:
		dn.setUrgency(dbusUrgencyCritical)
	default:
		dn.setUrgency(dbusUrgencyNormal)
	}
	//dn.setIcon(notification.Icon)
	//dn.setSound(notification.Sound)

	id := sendDbusNotification(conn, dn)
	notification.Id = id
	notification.SentAt = time.Now()
}

type dbusUrgency byte

const (
	dbusUrgencyLow      dbusUrgency = 0
	dbusUrgencyNormal   dbusUrgency = 1
	dbusUrgencyCritical dbusUrgency = 2
)

// dbusNotification holds all information needed for creating a notification
type dbusNotification struct {
	AppName       string // name of the application sending the notification (can be blank)
	ReplacesID    uint32 // id of the notification being replaced (0 means new notification)
	AppIcon       string // icon to display in the notification (can be blank for no icon)
	Summary       string // title of the notification
	Body          string // body of the notification
	Hints         map[string]dbus.Variant
	ExpireTimeout int // time before the notification expires (-1 for let server decide, 0 for never expire)
}

// See: https://specifications.freedesktop.org/notification-spec/latest/ar01s08.html
type hint struct {
	ID      string
	Variant dbus.Variant
}

// TODO: implement setting image
func hintIconFilePath(imageAbsolutePath string) hint {
	uri := fmt.Sprintf("file://%s", imageAbsolutePath)
	return hint{
		ID:      "image-path",
		Variant: dbus.MakeVariant(uri),
	}
}

// TODO: implement setting sound
func hintSoundName(soundName string) hint {
	return hint{
		ID:      "sound-name",
		Variant: dbus.MakeVariant(soundName),
	}
}

func hintUrgency(urgency dbusUrgency) hint {
	return hint{
		ID:      "urgency",
		Variant: dbus.MakeVariant(byte(urgency)),
	}
}

func (dn *dbusNotification) addHint(hint hint) {
	if dn.Hints == nil {
		dn.Hints = map[string]dbus.Variant{}
	}
	dn.Hints[hint.ID] = hint.Variant
}

func (dn *dbusNotification) setUrgency(urgency dbusUrgency) {
	dn.addHint(hintUrgency(urgency))
}

func (dn *dbusNotification) setIcon(iconName string) {
	dn.addHint(hintIconFilePath(iconName))
}

func (dn *dbusNotification) setSound(soundName string) {
	dn.addHint(hintSoundName(soundName))
}

func sendDbusNotification(conn *dbus.Conn, dn dbusNotification) (notificationId uint32) {
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call(
		"org.freedesktop.Notifications.Notify",
		0,
		dn.AppName,
		dn.ReplacesID,
		dn.AppIcon,
		dn.Summary,
		dn.Body,
		[]string{},
		dn.Hints,
		dn.ExpireTimeout,
	)
	util.ExitOnError(call.Err, "Failed to send notification")
	err := call.Store(&notificationId)
	util.ExitOnError(err, "Failed to store notification ID")

	return
}
