package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"github.com/ionous/errutil"
)

// ex. go run bgenCmd.go -in  ../../stories/blank.if
// go run bgenCmd.go -in ../../stories/shared/settings.if -pretty -out .
func main() {
	var inFile, outFile string
	var pretty bool
	flag.StringVar(&inFile, "in", "", "input file name (.if|.ifx|.ifspecs)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.blocks)")
	flag.BoolVar(&pretty, "pretty", false, "make the output somewhat human readable")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if e := convert(inFile, outFile, pretty); e != nil {
		println(e.Error())
	}
}

func convert(inFile, outFile string, pretty bool) (err error) {
	inExt := filepath.Ext(inFile)
	if len(inFile) == 0 {
		err = errutil.New("requires an input file")
	} else if decode, ok := decoders[inExt]; !ok {
		err = errutil.New("requires some sort of .if, .ifx, or .ifspecs file")
	} else if b, e := readOne(inFile); e != nil {
		err = e
	} else if m, e := decode(b); e != nil {
		err = e
	} else if x, e := block.Convert(m); e != nil {
		err = e
	} else {
		// create outfile name if needed
		const outExt = BlocksExt
		if len(outFile) == 0 {
			outFile = inFile[:len(inFile)-len(inExt)] + outExt
		} else if ext := filepath.Ext(outFile); len(ext) == 0 || ext == outFile {
			// convert directory
			base := filepath.Base(inFile)
			outFile = filepath.Join(outFile, base[:len(base)-len(inExt)]+outExt)
		}
		log.Println("writing", outFile)
		if x, e := prettify([]byte(x), pretty); e != nil {
			err = e
		} else if fp, e := os.Create(outFile); e != nil {
			err = e
		} else {
			_, err = fp.Write(x)
			fp.Close()
		}
	}
	return
}

var decoders = map[string]func([]byte) (jsn.Marshalee, error){
	SpecExt:     decodeSpec,
	CompactExt:  decodeCompact,
	DetailedExt: decodeDet,
}

const (
	SpecExt     = ".ifspecs"
	DetailedExt = ".ifx"
	CompactExt  = ".if"
	BlocksExt   = ".blocks"
)

func prettify(str []byte, pretty bool) (ret []byte, err error) {
	// if !pretty {
	ret = str
	if pretty {
		var indent bytes.Buffer
		if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
			// err = e
			log.Println(e)
		} else {
			ret = indent.Bytes()
		}
	}
	return
}

func decodeCompact(b []byte) (ret jsn.Marshalee, err error) {
	var dst story.StoryFile
	if e := story.Decode(&dst, b, tapestry.AllSignatures); e != nil {
		err = e
	} else {
		ret = &dst
	}
	return
}

func decodeDet(b []byte) (ret jsn.Marshalee, err error) {
	var dst story.StoryFile
	if e := din.Decode(&dst, tapestry.Registry(), b); e != nil {
		err = e
	} else {
		ret = &dst
	}
	return
}

func decodeSpec(b []byte) (ret jsn.Marshalee, err error) {
	var dst spec.TypeSpec
	if e := cin.Decode(&dst, b, cin.Signatures(tapestry.AllSignatures)); e != nil {
		err = e
	} else {
		ret = &dst
	}
	return
}

func readOne(filePath string) (ret []byte, err error) {
	log.Println("reading", filePath)
	if fp, e := os.Open(filePath); e != nil {
		err = e
	} else {
		ret, err = io.ReadAll(fp)
		fp.Close()
	}
	return
}
