// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package base

type UsageError struct {
	Cmd   *Command
	Cause error
}

func (e UsageError) Error() (ret string) {
	if e.Cause != nil {
		ret = e.Cause.Error()
	}
	return
}
