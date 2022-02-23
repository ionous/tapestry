package serve

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os/exec"

	"git.sr.ht/~ionous/tapestry/dl/play"
)

// run the passed command
// pass strings from "in" to the command's stdin
// transform its stdout to PlayOut messages; stderr as PlayLog.
// returns after starting the command ( returning any startup errors )
func Play(exe, inFile string, cs *Channels) (err error) {
	var pipes pipes
	defer pipes.Close()

	cmd := exec.Command(exe, "-in", inFile, "-json")
	if outp, e := pipes.AddReader(cmd.StdoutPipe()); e != nil {
		err = e
	} else if logp, e := pipes.AddReader(cmd.StderrPipe()); e != nil {
		err = e // note:  stderr is the default target for package log.
	} else if inp, e := pipes.AddWriter(cmd.StdinPipe()); e != nil {
		err = e
	} else {
		goScanText(logp, func(line string) {
			log.Println("logged:", line)
			cs.msgs <- &play.PlayLog{Log: line}
		})
		goScanJson(outp, func(line string) {
			log.Println("wrote:", line)
			cs.msgs <- &play.PlayOut{Out: line}
		})
		goWrite(inp, cs.input)
		//
		if e := cmd.Run(); e != nil {
			cs.msgs <- &play.PlayLog{Log: e.Error()}
			err = e
		}
	}
	return
}

// reads string from the channel and write them to the writer
func goWrite(w io.Writer, in <-chan string) {
	go func() {
		for str := range in {
			log.Println("posting:", str)
			if _, e := io.WriteString(w, str+"\n"); e != nil {
				log.Println("input error:", e)
				break // just eats the error for now
			}
		}
		log.Println("finished reading user input")
	}()
}

// https://www.yellowduck.be/posts/reading-command-output-line-by-line/
// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func goScanText(r io.Reader, onText func(string)) {
	scan := bufio.NewScanner(r) // by default the scanner separates input line by line
	go func() {
		for scan.Scan() {
			onText(scan.Text())
		}
	}()
}

// reads from a stream of JSON objects.
func goScanJson(r io.Reader, onText func(string)) {
	type Obj struct {
		Out string `json:"out"`
	}
	scan := json.NewDecoder(r)
	go func() {
		for {
			var obj Obj
			if e := scan.Decode(&obj); e != nil {
				log.Println("error scanning:", e)
				break
			} else {
				onText(obj.Out)
			}
		}
	}()
}

// func checkLsExists() {
//     path, err := exec.LookPath("ls")
//     if err != nil {
//         fmt.Printf("didn't find 'ls' executable\n")
//     } else {
//         fmt.Printf("'ls' executable is in '%s'\n", path)
//     }
// }
