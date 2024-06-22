package play

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/game"
)

func CaptureInput(pt *Playtime) {
	var done bool
	// fix? used for fabricate; maybe use options instead so that we can have multiple instances?
	debug.Stepper = func(words string) (err error) {
		// FIX: errors for step are getting fmt.Println in playTime.go
		// so expect output can't test for errors ( and on error looks a bit borken )
		var sig game.Signal
		if _, e := pt.Step(words); !done && errors.As(e, &sig) && sig == game.SignalQuit {
			done = true // eat the quit signal on first return; fix? maybe do this on client?
		} else if e != nil {
			err = e
		}
		return
	}
}
