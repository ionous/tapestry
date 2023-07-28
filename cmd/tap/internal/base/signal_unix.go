// Copyright 2012 The Go Authors. All rights reserved.
// Copyright 2022 - Modifications by Simon Travis.
// 
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix || js

package base

import (
	"os"

	"golang.org/x/sys/unix"
)

var signalsToIgnore = []os.Signal{os.Interrupt, unix.SIGQUIT}

// SignalTrace is the signal to send to make a Go program
// crash with a stack trace.
var SignalTrace os.Signal = unix.SIGQUIT
