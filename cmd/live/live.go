package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"go/build"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// go run live.go -in /Users/ionous/Documents/Tapestry/stories/shared [-out /Users/ionous/Documents/Tapestry/build/play.db]
// go run live.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	cfg := composer.DevConfig(build.Default.GOPATH)
	//
	var srcPath, outFile string
	var check bool
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&check, "check", true, "run check after importing?")
	flag.Parse()
	msgs := make(chan string, 1024)

	// read from stdin
	// go func() {
	// r := bufio.NewReader(os.Stdin)
	// 	for {
	// 		if in, e := r.ReadString('\n'); e != nil {
	// 			break
	// 		} else if cnt := len(in); cnt > 0 {
	// 			msgs <- in[:cnt-1] // trim the newline.
	// 		}
	// 	}
	// }()

	buildMsgs := func() (retCnt int, retStr string) {
		var out js.Builder // fix: really need to make this streamable.
		out.Brace(js.Array, func(a *js.Builder) {
			// some arbitrary maximum size
			for i := 0; i < 512; i++ {
				select {
				case msg := <-msgs:
					if i > 0 {
						a.R(js.Comma)
					}
					a.Q(msg) // json quote each string.
					retCnt++
				default:
					break // done
				}
			}
		})
		retStr = out.String()
		return
	}
	writeMsgs := func(w http.ResponseWriter) (ret int, err error) {
		w.Header().Set("Content-Type", "application/json")
		cnt, msgs := buildMsgs()
		if _, e := io.WriteString(w, msgs); e != nil {
			err = e
		} else {
			ret = cnt
		}
		return
	}

	// backend
	u, _ := url.Parse("http://localhost:3000/")
	p := httputil.NewSingleHostReverseProxy(u)
	http.Handle("/", p)
	ready := false

	// input/output processing
	http.HandleFunc("/io/", web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == "io" {
				ret = &web.Wrapper{

					// polling for data.
					Gets: func(ctx context.Context, w http.ResponseWriter) (err error) {
						if cnt, e := writeMsgs(w); e != nil {
							err = e
						} else {
							log.Println("wrote", cnt, "msgs")
						}
						return
					},
					// sending a command.
					Posts: func(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
						if !ready {
							_, err = io.WriteString(w, "(not ready)")
						} else if msg, e := Decode(r); e != nil {
							err = e
						} else {
							select {
							case msgs <- msg:
								log.Println("posted:", msg)
								_, err = writeMsgs(w)
							default:
								log.Println("ignored:", msg)
								_, err = io.WriteString(w, "(server busy)")
							}
						}
						return
					},
				}
			}
			return
		}}))
	//
	if _, e := GoAsm(cfg.Assemble, srcPath, outFile, check, msgs, &ready); e != nil {
		log.Fatal(e)
	} else {
		log.Println("Listening on port", ":8088", "...")
		if e := http.ListenAndServe(":8088", nil); e != nil {
			log.Fatal(e)
		}
	}
}

func Decode(r io.Reader) (ret string, err error) {
	var cmd struct {
		Msg string `json:"cmd"`
	}
	dec := json.NewDecoder(r)
	// decode an array value (Message)
	if e := dec.Decode(&cmd); e != nil {
		err = e
	} else {
		ret = cmd.Msg
	}
	return
}

// https://www.yellowduck.be/posts/reading-command-output-line-by-line/
func GoAsm(asm, srcPath, outFile string, check bool, msgs chan<- string, done *bool) (ret string, err error) {
	log.Println("Assembling", srcPath+"...")
	if len(outFile) == 0 {
		dir, _ := filepath.Split(srcPath)
		outFile = filepath.Join(dir, "play.db")
	}

	cmd := exec.Command(
		asm,
		"-in", srcPath,
		"-out", outFile,
		"-check", strconv.FormatBool(check),
	)

	// Get a pipe to read from standard out
	if r, e := cmd.StdoutPipe(); e != nil {
		err = e
	} else {
		cmd.Stderr = cmd.Stdout        // Use the same pipe for standard error
		scanner := bufio.NewScanner(r) // Create a scanner which scans r in a line-by-line fashion
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				msgs <- line
			}
			if e := cmd.Wait(); e != nil {
				msgs <- e.Error()
			}
			*done = true
		}()
		if e := cmd.Start(); e != nil {
			err = e
		} else {
			ret = outFile
		}
	}
	return
}
