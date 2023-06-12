package mdl

import (
  "github.com/ionous/errutil"
)

// when the definition would repeat existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Duplicate = errutil.Error("Duplicate")

// when the definition would contradict existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Conflict = errutil.Error("Conflict")

// when the definition can't find some required information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Missing = errutil.Error("Missing")
