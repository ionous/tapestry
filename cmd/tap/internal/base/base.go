// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package base defines shared basic pieces of the go command,
// in particular logging and the Command structure.
// see also: https://pkg.go.dev/cmd/go/internal/base
package base

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cfg"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/str"
)

// A Command is an implementation of a go command
// like go build or go fix.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(ctx context.Context, cmd *Command, args []string) error

	// UsageLine is the one-line usage message.
	// The words between "go" and the first flag or argument in the line are taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool

	// Commands lists the available commands and help topics.
	// The order here is the order in which they are printed by 'go help'.
	// Note that subcommands are in general best avoided.
	Commands []*Command
}

var Exe string = "tap"

var Go = &Command{
	UsageLine: Exe,
	Long:      `Tap manages Tapestry stories.`,
	// Commands initialized in package main init()
}

// hasFlag reports whether a command or any of its subcommands contain the given
// flag.
func hasFlag(c *Command, name string) bool {
	if f := c.Flag.Lookup(name); f != nil {
		return true
	}
	for _, sub := range c.Commands {
		if hasFlag(sub, name) {
			return true
		}
	}
	return false
}

// LongName returns the command's long name: all the words in the usage line between "go" and a flag or argument,
func (c *Command) LongName() string {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}
	if name == Exe {
		return ""
	}
	return strings.TrimPrefix(name, Exe+" ")
}

// Name returns the command's short name: the last word in the usage line before a flag or argument.
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "Run '%s help %s' for details.\n", Exe, c.LongName())
	SetExitStatus(2)
	Exit()
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var atExitFuncs []func()

func AtExit(f func()) {
	atExitFuncs = append(atExitFuncs, f)
}

func ExitWithStatus(status int) {
	SetExitStatus(status)
	Exit()
}

func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func Fatalf(format string, args ...any) {
	Errorf(format, args...)
	Exit()
}

func Errorf(format string, args ...any) {
	log.Printf(format, args...)
	SetExitStatus(1)
}

func ExitIfErrors() {
	if exitStatus != 0 {
		Exit()
	}
}

var exitStatus = 0
var exitMu sync.Mutex

func SetExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func GetExitStatus() int {
	return exitStatus
}

// Run runs the command, with stdout and stderr
// connected to the go command's own stdout and stderr.
// If the command fails, Run reports the error using Errorf.
func Run(cmdargs ...any) {
	cmdline := str.StringList(cmdargs...)

	// -n and -x both print the commands; -n doesn't run them.
	if cfg.BuildN || cfg.BuildX {
		fmt.Printf("%s\n", strings.Join(cmdline, " "))
		if cfg.BuildN {
			return
		}
	}

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Errorf("%v", err)
	}
}

// RunStdin is like run but connects Stdin.
func RunStdin(cmdline []string) {
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = cfg.OrigEnv
	StartSigHandlers()
	if err := cmd.Run(); err != nil {
		Errorf("%v", err)
	}
}
