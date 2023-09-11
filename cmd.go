package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	commands map[string]*command
	command  struct {
		name         string
		exec         func([]string)
		requiredArgs int
		help         string
	}
)

func newCommand(name string, fn func([]string)) *command {
	c := command{exec: fn, name: name}
	return &c
}

func (c *command) WithRequiredArgs(n int) *command {
	c.requiredArgs = n
	return c
}

func (c *command) WithHelp(help string) *command {
	c.help = help
	return c
}

func (c *command) helpStr(fields ...string) string {
	if c.help == "" {
		return ""
	}
	if len(fields) == 0 {
		fields = []string{c.name}
	}
	return strings.ReplaceAll(c.help, "$0", fields[0])
}

func newCommands() commands {
	return make(commands)
}

func (c commands) add(name string, fn func([]string), aliases ...string) *command {
	c[name] = newCommand(name, fn)
	for _, alias := range aliases {
		c.alias(name, alias)
	}
	return c[name]
}

func (c commands) addMinArgs(name string, fn func([]string), minArgs int, aliases ...string) {
	c.add(name, fn, aliases...).WithRequiredArgs(minArgs)
}

func (c commands) addHelp(name string, fn func([]string), help string, aliases ...string) {
	c.add(name, fn, aliases...).WithHelp(help)
}

func (c commands) alias(name, alias string) {
	c[alias] = c[name]
}

func sleep(fields []string) {
	if len(fields) != 2 {
		printf(cmds["sleep"].helpStr(fields...))
		return
	}
	var dur string
	switch {
	case strings.HasSuffix(fields[1], "s"):
		dur = fields[1][:len(fields[1])-1]
		num, err := strconv.Atoi(dur)
		if err != nil {
			printf("invalid duration: %s", dur)
			printf(cmds["sleep"].helpStr(fields...))
			return
		}
		time.Sleep(time.Duration(num) * time.Second)
	case strings.HasSuffix(fields[1], "ms"):
		dur = fields[1][:len(fields[1])-2]
		num, err := strconv.Atoi(dur)
		if err != nil {
			printf("invalid duration: %s", dur)
			printf(cmds["sleep"].helpStr(fields...))
			return
		}
		time.Sleep(time.Duration(num) * time.Millisecond)
	}
}

func init() {
	cmds.addMinArgs("ls", ls, 0, "list", "status")
	cmds.addMinArgs("read", readPin, 1, "cat", "get")
	cmds.addMinArgs("write", writePin, 2, "set")
	cmds["write"].WithHelp("usage: $0 <pin> <state|pull> <value>")
	cmds.addMinArgs("sleep", sleep, 1, "wait")
	cmds["sleep"].WithHelp("usage: $0 <number><s|ms>\n\texample: $0 1s")
	cmds.add("exit", func([]string) { os.Exit(0) })
	cmds.add("help", func([]string) {
		for name, cmd := range cmds {
			printf("%s\t%s", name, cmd.help)
		}
	})
}
