// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package base

import "errors"

// ex. fmt.Errorf("%w expected at least one word to transform", base.ErrUsage)
var ErrUsage = errors.New("error: ")
