// Adapted from https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/cmd/go/main.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cfg"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/help"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/idlcmd"
)

func main() {
	flag.Usage = exitBadUsage
	flag.Parse()
	log.SetFlags(0) // https://pkg.go.dev/log#pkg-constants - by default it logs date and time.
	args := flag.Args()
	if len(args) < 1 {
		exitBadUsage()
	}

	// TODO: env for tap home?

	// this is a recurisve drill down into command lists expanded as a loop:
	// it starts with the commands assembled from the sub-packages ( in init() below )
	for topCmds := base.Go.Commands; ; {
		var name string
		name, args = args[0], args[1:]
		if name == "help" {
			// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
			// ( ie. this cuts "help" out of the inputed commands )
			help.Help(os.Stdout, append(cfg.CmdNames, args...))
			return // st: why not base.Exit()? its not clear.

		} else if cmd := findCommand(name, topCmds); cmd == nil {
			exitUnknownCommand()

		} else if len(cmd.Commands) == 0 {
			// originally this would move to the next entry in the list of subcommands...
			// but unless there were two commands with the same name, one runnable and one not:
			// it would eventually wind up as not found; so shortcut the confusion here.
			if !cmd.Runnable() {
				exitUnknownCommand()
			} else {
				invokeAndExit(cmd, args)
			}
		} else if len(args) == 0 {
			// we have sub commands, but no args to process select one of those commands:
			help.PrintUsage(os.Stderr, cmd)
			base.ExitWithStatus(2)

		} else {
			// drill down into the commands:
			topCmds = cmd.Commands
			// remembering the path we took to get there:
			cfg.CmdNames = append(cfg.CmdNames, name)
			// NOTE: this is the only path that continues the outer for loop
		}
	}
}

// returns the matching command if it needs to be expanded into a sub command
// returns nil if no matching commands was found
func findCommand(name string, cmds []*base.Command) (ret *base.Command) {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			ret = cmd
			break
		}
	}
	return //
}

func invokeAndExit(cmd *base.Command, args []string) {
	// 'go env' handles checking the build config
	// if cmd != envcmd.CmdEnv {
	// 	buildcfg.Check()
	// 	if cfg.ExperimentErr != nil {
	// 		base.Fatalf("go: %v", cfg.ExperimentErr)
	// 	}
	// }

	// // Set environment (GOOS, GOARCH, etc) explicitly.
	// // In theory all the commands we invoke should have
	// // the same default computation of these as we do,
	// // but in practice there might be skew
	// // This makes sure we all agree.
	// cfg.OrigEnv = os.Environ()
	// cfg.CmdEnv = envcmd.MkEnv()
	// for _, env := range cfg.CmdEnv {
	// 	if os.Getenv(env.Name) != env.Value {
	// 		os.Setenv(env.Name, env.Value)
	// 	}
	// }

	// cmd.Flag.Usage = func() { cmd.Usage() }
	if !cmd.CustomFlags {
		// base.SetFromGOFLAGS(&cmd.Flag)
		cmd.Flag.Parse(args)
		args = cmd.Flag.Args()
	}
	ctx := context.Background()
	// ctx = maybeStartTrace(ctx)
	// ctx, span := trace.StartSpan(ctx, fmt.Sprint("Running ", cmd.Name(), " command"))
	cmd.Run(ctx, cmd, args)
	// span.Done()
	base.Exit()
}

func exitUnknownCommand() {
	var helpArg string
	if last := len(cfg.CmdNames) - 1; last > 0 {
		helpArg = " " + cfg.CmdNames[last]
	}
	fmt.Fprintf(os.Stderr, "tap %s: unknown command\nRun 'tap help%s' for usage.\n", cfg.CmdNames, helpArg)
	base.ExitWithStatus(2)
}

func exitBadUsage() {
	help.PrintUsage(os.Stderr, base.Go)
	os.Exit(2) // st: why not base.Exit()? its all very confusing
}

func init() {
	// rewrites the main tap command to simplify exitBadUsage
	base.Go.Commands = []*base.Command{
		idlcmd.CmdIdl,
	}
}
