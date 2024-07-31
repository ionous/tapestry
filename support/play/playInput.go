package play

import (
	"git.sr.ht/~ionous/tapestry/dl/debug"
)

// used for letting the unit tests fabricate input as if typed by the player.
func CaptureInput(pt *Playtime) {
	// fix? used for fabricate; maybe use options instead so that we can have multiple instances?
	debug.Stepper = func(words string) (err error) {
		// FIX: errors for step are getting fmt.Println in playTime.go
		// so expect output can't test for errors ( and on error looks a bit borken )
		_, err = pt.HandleTurn(words)
		return
	}
}
