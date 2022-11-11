// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cfg holds configuration shared by multiple parts
// of the go command.
package cfg

import (
	"strings"
)

// These are general "build flags" used by build and other commands.
var (
	// BuildA                 bool     // -a flag
	// BuildBuildmode         string   // -buildmode flag
	// BuildBuildvcs          = "auto" // -buildvcs flag: "true", "false", or "auto"
	// BuildContext           = defaultContext()
	// BuildMod               string                  // -mod flag
	// BuildModExplicit       bool                    // whether -mod was set explicitly
	// BuildModReason         string                  // reason -mod was set, if set by default
	// BuildI                 bool                    // -i flag
	// BuildLinkshared        bool                    // -linkshared flag
	// BuildMSan              bool                    // -msan flag
	// BuildASan              bool                    // -asan flag
	BuildN bool // -n flag; prints commands but does not run them.
	// BuildO                 string                  // -o flag
	// BuildP                 = runtime.GOMAXPROCS(0) // -p flag
	// BuildPkgdir            string                  // -pkgdir flag
	// BuildRace              bool                    // -race flag
	// BuildToolexec          []string                // -toolexec flag
	// BuildToolchainName     string
	// BuildToolchainCompiler func() string
	// BuildToolchainLinker   func() string
	// BuildTrimpath          bool // -trimpath flag
	// BuildV                 bool // -v flag; print the names of packages as they are compiled.
	// BuildWork              bool // -work flag
	BuildX bool // -x flag; prints commands before running them.

	// ModCacheRW bool   // -modcacherw flag
	// ModFile    string // -modfile flag

	CmdNames CommandNames // "build", "install", "list", "mod tidy", etc.

	// DebugActiongraph string // -debug-actiongraph flag (undocumented, unstable)
	// DebugTrace       string // -debug-trace flag

	// GoPathError is set when GOPATH is not set. it contains an
	// explanation why GOPATH is unset.
	// GoPathError string
)

// OrigEnv is the original environment of the program at startup.
var OrigEnv []string

type CommandNames []string

func (cs CommandNames) String() string {
	return strings.Join(cs, " ")
}
