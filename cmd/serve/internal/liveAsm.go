package serve

import (
	"os/exec"
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/play"
)

func Asm(exe, srcPath, outFile string, check bool, cs *Channels) (ret string, err error) {
	cmd := exec.Command(
		exe,
		"-in", srcPath,
		"-out", outFile,
		"-check", strconv.FormatBool(check),
	)
	// creates a pipe to read from standard out
	if r, e := cmd.StdoutPipe(); e != nil {
		err = e
	} else {
		cmd.Stderr = cmd.Stdout // assign the same *writer*

		// post output as log messages
		goScanText(r, func(line string) {
			cs.msgs <- &play.PlayLog{Log: line}
		})

		// run blocking
		if e := cmd.Run(); e != nil {
			cs.msgs <- &play.PlayLog{Log: e.Error()}
			err = e
		} else {
			ret = outFile
		}
	}
	return
}
