package files

import (
	"fmt"
	"runtime/debug"
)

// for lack of a better place... here we arere.
func GetVersion(details bool) (ret string) {
	if b, ok := debug.ReadBuildInfo(); !ok {
		ret = "???" // only should happen if not built with modules
	} else {
		if details {
			ret = b.String()
		} else {
			if m := b.Main.Version; len(m) > 0 {
				ret = m
			} else {
				ret = "unknown" //provisionally
				for _, d := range b.Deps {
					if d != nil && d.Path == "git.sr.ht/~ionous/tapestry" {
						ret = d.Version
						break
					}
				}
			}
			ret = fmt.Sprintf("%s (%s)", ret, b.GoVersion)
		}
	}
	return
}
