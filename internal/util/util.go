package util

import (
	"fmt"
	"log"
	"os/exec"
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

// AssertDependenciesInstalled checks if all executable dependencies are present on $PATH, and terminates
// the program with an error message if not.
func AssertDependenciesInstalled() {
	executables := []string{
		"notify-send",
	}
	for _, e := range executables {
		_, err := exec.LookPath(e)
		ExitOnError(err, "Missing dependency")
	}
}
