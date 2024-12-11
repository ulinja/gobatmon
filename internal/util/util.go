package util

import (
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
)

// ExitOnError logs a fatal error message and exits the program with the provided prefix if error is not nil.
func ExitOnError(err error, prefix string) {
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

func AssertDbusNotificationsAvailable() {
	conn, err := dbus.ConnectSessionBus()
	ExitOnError(err, "DBUS session bus is not available")
	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.GetCapabilities", 0)
	ExitOnError(call.Err, "DBUS notification service is not available")
	var ret []string
	err = call.Store(&ret)
	ExitOnError(err, "DBUS notification service is not available")
}
