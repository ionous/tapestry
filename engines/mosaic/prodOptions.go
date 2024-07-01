//go:build production

package main

import (
	"io/fs"

	"git.sr.ht/~ionous/tapestry/www"
)

const BuildConfig = Prod

func init() {
	if sub, e := fs.Sub(www.Dist, "dist"); e != nil {
		panic(e)
	} else {
		Frontend = sub
	}
}
