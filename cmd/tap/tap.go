// Adapted from https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/cmd/go/main.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cfg"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/help"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/idlcmd"
	"github.com/ionous/errutil"
)

func main() {
	defer base.Exit()
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if e := cmdLine.Parse(os.Args[1:]); e != nil {
		help.PrintUsage(os.Stderr, base.Go)
		base.SetExitStatus(2)
	} else {
		log.SetFlags(0) // https://pkg.go.dev/log#pkg-constants - by default it logs date and time.
		if e := handleMain(cmdLine.Args()); e != nil {
			if e == UnknownCommand {
				var helpArg string
				if last := len(cfg.CmdNames) - 1; last > 0 {
					helpArg = " " + cfg.CmdNames[last]
				}
				// for reasons i don't understand this path in go doesn't print the status number
				// yet, here it does...
				log.Printf("%s %s: unknown command\nRun 'tap help%s' for usage.\n", base.Exe, cfg.CmdNames, helpArg)
				base.SetExitStatus(2)
			} else {
				if cause := e.Error(); len(cause) > 0 {
					log.Println(cause)
					base.SetExitStatus(2)
				}
				var u base.UsageError
				if errors.As(e, &u) {
					help.PrintUsage(os.Stderr, u.Cmd)
				}
			}
		}
	}
}

func handleMain(args []string) (err error) {
	if len(args) < 1 {
		err = base.UsageError{Cmd: base.Go}
	} else {
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
				break // note: asking for help is not an error

			} else {
				cfg.CmdNames = append(cfg.CmdNames, name)
				if cmd := findCommand(name, topCmds); cmd == nil {
					err = UnknownCommand
					break

				} else if len(cmd.Commands) == 0 {
					// originally this would move to the next entry in the list of subcommands...
					// but unless there were two commands with the same name, one runnable and one not:
					// it would eventually wind up as not found; so shortcut the confusion here.
					if !cmd.Runnable() {
						err = UnknownCommand
						break
					} else {
						err = invoke(cmd, args)
						break
					}
				} else if len(args) == 0 {
					err = base.UsageError{Cmd: cmd}
					break

				} else {
					topCmds = cmd.Commands
					// NOTE: this is the only path that continues the outer for loop
				}
			}
		}
	}
	return
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
	return
}

func invoke(cmd *base.Command, args []string) (err error) {
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
	return cmd.Run(ctx, cmd, args)
	// span.Done()
}

const UnknownCommand errutil.Error = "unknown command"

func init() {
	// rewrites the main tap command to simplify exitBadUsage
	base.Go.Commands = []*base.Command{
		idlcmd.CmdIdl,
	}
}
