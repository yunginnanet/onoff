package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Nerdmaster/terminal"
	"github.com/stianeikeland/go-rpio/v4"
)

func writePin(fields []string) {
	first := fields[0]
	var (
		pnum uint64
		err  error
	)
	if pnum, err = strconv.ParseUint(fields[1], 10, 8); err != nil {
		printf("invalid pin number")
		printf("usage: %s <pin> <state>", first)
		return
	}
	switch strings.ToLower(fields[2]) {
	case "state":
		if len(fields) < 4 {
			printf("invalid state, hint: high/low, 1/0, on/off, true/false")
			printf(cmds["write"].helpStr(fields...))
			return
		}
		var state rpio.State
		switch strings.ToLower(fields[3]) {
		case "high", "1", "on", "true":
			state = rpio.High
		case "low", "0", "off", "false":
			state = rpio.Low
		default:
			printf("invalid state, hint: high/low, 1/0, on/off, true/false")
			printf(cmds["write"].helpStr(fields...))
			return
		}
		rpio.Pin(pnum).Write(state)
	case "pull":
		if len(fields) < 4 {
			printf("invalid pull, hint: none/up/down/off")
			printf(cmds["write"].helpStr(fields...))
			return
		}
		var pull rpio.Pull
		switch strings.ToLower(fields[3]) {
		case "none":
			pull = rpio.PullNone
		case "up":
			pull = rpio.PullUp
		case "down":
			pull = rpio.PullDown
		case "off":
			pull = rpio.PullOff
		default:
			printf("invalid pull, hint: none/up/down/off")
			printf(cmds["write"].helpStr(fields...))
			return
		}
		rpio.Pin(pnum).Pull(pull)
	case "input":
		rpio.Pin(pnum).Input()
	case "output":
		rpio.Pin(pnum).Output()
	default:
		printf("invalid mode, hint: state/pull/input/output")
		printf(cmds["write"].helpStr(fields...))
		return
	}
	printf("Pin %d:\t%s\t%s", pnum, stateStr(rpio.Pin(pnum).Read()), pullStr(rpio.Pin(pnum).ReadPull()))
}

func readPin(fields []string) {
	first := fields[0]
	if len(fields) < 2 {
		printf("usage: " + first + " <pin>")
		return
	}
	var (
		pnum uint64
		err  error
	)
	if pnum, err = strconv.ParseUint(fields[1], 10, 8); err != nil {
		printf("invalid pin number")
		printf("usage: %s <pin>", first)
		return
	}
	pin := rpio.Pin(pnum)
	printf("Pin %d:\t%s\t%s", pnum, stateStr(pin.Read()), pullStr(pin.ReadPull()))
}

func ls(fields []string) {
	if len(fields) > 1 {
		printf("usage: ls")
		return
	}
	// list all pins
	for i := rpio.Pin(0); i < 28; i++ {
		printf("Pin %d:\t%s\t%s", i, stateStr(i.Read()), pullStr(i.ReadPull()))
	}
}

var cmds = newCommands()

func commandLoop(prompt *terminal.Prompt, ctx context.Context) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		line, err := prompt.ReadLine()
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("error reading command: %v\n", err))
		}
		for _, cmd := range strings.Split(line, ";") {
			fields := strings.Fields(cmd)
			first := fields[0]
			var ran *command
			var ok bool
			if ran, ok = cmds[first]; !ok {
				printf("unknown command: %s", first)
				break
			}
			if len(fields)-1 < ran.requiredArgs {
				printf("not enough arguments, expected %d", ran.requiredArgs)
				if ran.help != "" {
					printf("%s\n\t%s", first, ran.helpStr(fields...))
				}
				break
			}
			ran.exec(fields)
			continue
		}
	}
}

func startUI() {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
		return
	}
	defer restoreState(oldState)
	doneChan := make(chan bool)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	ctx, cancel := context.WithCancel(context.Background())
	prompt := terminal.NewPrompt(os.Stdin, os.Stdout, "gpio@"+host+"> ")
	go commandLoop(prompt, ctx)

	select {
	case <-doneChan:
	case <-sigChan:
	}

	cancel()
}

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		printf(err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := rpio.Close(); err != nil {
			printf(err.Error())
			os.Exit(1)
		}
	}()

	startUI()

	// Set pin to output mode
	// pin.Output()

	// Toggle pin 20 times
	/*	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}*/
}
