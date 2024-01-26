package story

import (
	"github.com/ionous/errutil"
)

const UnhandledSwap = errutil.Error("unhandled swap")
const MissingSlot = errutil.Error("missing slot")
const InvalidValue = errutil.Error("invalid value")
