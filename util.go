package main

import (
	"fmt"
	"os"

	"github.com/Nerdmaster/terminal"
	"github.com/stianeikeland/go-rpio/v4"
)

func restoreState(oldState *terminal.State) {
	if err := terminal.Restore(int(os.Stdin.Fd()), oldState); err != nil {
		panic(err)
	}
}

func stateStr(mode rpio.State) string {
	switch mode {
	case rpio.Low:
		return "Low"
	case rpio.High:
		return "High"
	default:
		return "Unknown"
	}
}

func pullStr(mode rpio.Pull) string {
	switch mode {
	case rpio.PullNone:
		return "None"
	case rpio.PullUp:
		return "Up"
	case rpio.PullDown:
		return "Down"
	case rpio.PullOff:
		return "Off"
	default:
		return "Unknown"
	}
}

func printf(format string, args ...interface{}) {
	_, _ = os.Stdout.WriteString(fmt.Sprintf(format, args...))
	_, _ = os.Stdout.WriteString("\n\r")
}
