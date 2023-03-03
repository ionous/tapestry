// alter json (.if) files using regexp
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"git.sr.ht/~ionous/tapestry/web/files"
	"github.com/ionous/errutil"
)

func main() {
	var inPath, outPath, match, inExts string
	var verbose, recurse bool
	flag.StringVar(&outPath, "out", "", "output directory")
	flag.StringVar(&inPath, "in", "", "input file(s) or paths(s) (comma separated)")
	flag.BoolVar(&recurse, "recurse", false, "scan input sub-directories")
	flag.StringVar(&inExts, "filter", ".if",
		`extension(s) for directory scanning.
ignored if 'in' refers to a specific file`)
	flag.StringVar(&match, "match", "", ",regexp file: alternating lines of regexp and replacement")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.BoolVar(&verbose, "verbose", false, "write line matches and their results")
	flag.Parse()

	if len(inPath) == 0 || len(outPath) == 0 || len(match) == 0 {
		flag.Usage()
	} else if reps, e := readRexp(match); e != nil {
		log.Fatal("error reading match file: ", e)
	} else {
		process := func(inFile string) (err error) {
			if b, e := os.ReadFile(inFile); e != nil {
				err = errutil.New("couldnt read", inFile, e)
			} else {
				var data map[string]any
				if e := json.Unmarshal(b, &data); e != nil {
					err = errutil.New("couldnt unmarshal", inFile, e)
				} else {
					// deprettify
					var deprettify bytes.Buffer
					js := json.NewEncoder(&deprettify)
					js.SetEscapeHTML(false)
					if e := js.Encode(data); e != nil {
						err = errutil.New("couldnt encode", inFile, e)
					} else {
						// replace
						out := reps.ReplaceAllString(deprettify.String(), verbose)

						//re-prettify
						outFile := filepath.Join(outPath, filepath.Base(inFile))
						bytes := []byte(out)
						if e := json.Unmarshal(bytes, &data); e != nil {
							if e := writePlain(outFile, bytes); e != nil {
								err = errutil.New("couldnt write", outFile, e)
							} else {
								err = errutil.New("couldnt prettify", outFile)
							}
						} else {
							if e := writeJson(outFile, data, true); e != nil {
								err = errutil.New("couldnt write", outPath, e)
							}
						}
					}
				}
			}
			return
		}
		if e := files.ReadPaths(inPath, recurse, strings.Split(inExts, ","), process); e != nil {
			log.Fatal("error processing files: ", e)
		}
	}
	return
}

func writePlain(outPath string, bytes []byte) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		fp.Write(bytes)
	}
	return
}

func writeJson(outPath string, data map[string]any, pretty bool) (err error) {
	log.Println("writing", outPath)
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		js := json.NewEncoder(fp)
		js.SetEscapeHTML(false)
		if pretty {
			js.SetIndent("", "  ")
		}
		err = js.Encode(data)
	}
	return
}

func prettify(str string, pretty bool) (ret []byte) {
	ret = []byte(str)
	if pretty {
		var indent bytes.Buffer
		if e := json.Indent(&indent, ret, "", "  "); e != nil {
			log.Println(e)
		} else {
			ret = indent.Bytes()
		}
	}
	return
}

type repSet struct {
	exps []*regexp.Regexp
	reps []string
}

func (rs repSet) ReplaceAllString(out string, verbose bool) string {

	// For each match of the regex in the content.
	for i, exp := range rs.exps {
		if verbose {
			for _, m := range exp.FindAllStringSubmatchIndex(out, -1) {
				// 0: is the whole string that matched
				min, max := m[0], m[1]
				log.Println("matched:", out[min:max])
			}
		}

		out = exp.ReplaceAllString(out, rs.reps[i])
	}
	return out
}

func readRexp(match string) (ret repSet, err error) {
	if c, e := os.ReadFile(match); e != nil {
		err = e
	} else {
		var rep bool
		for i, el := range strings.Split(string(c), "\n") {
			if line := strings.TrimSpace(el); len(line) > 0 && !strings.HasPrefix(line, "--") {
				if rep {
					ret.reps = append(ret.reps, line)
					rep = false
				} else if x, e := regexp.Compile(line); e != nil {
					err = errutil.New("couldnt compile", i, e)
					break
				} else {
					ret.exps = append(ret.exps, x)
					rep = true
				}
			}
		}

	}
	return
}
