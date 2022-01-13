// Package main for 'chart'.
// Generates .dot file that can be turned into an image, a webpage, etc. using the graphiz "dot" tool ( https://graphviz.org/ )
package main

import (
	"flag"
	"log"

	chart "git.sr.ht/~ionous/tapestry/cmd/chart/internal"
	"github.com/ionous/errutil"
)

// usage:
// 1. chart [-d=entire_game] -in play.db [-out chart.dot]
// 2. dot -Tpng chart.dot -o chart.png
// example:
// go run chart.go -d traversal -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	var inFile, outFile, scope string
	flag.StringVar(&inFile, "in", "", "input file name (.db)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.dot)")
	flag.StringVar(&scope, "d", "", "optional domain")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		println("requires an input file")
	} else {
		if len(outFile) == 0 {
			outFile = "chart.dot"
		}
		if len(scope) == 0 {
			scope = "entire_game"
		}
		if _, e := chart.Chart(inFile, outFile, scope); e != nil {
			log.Panic(e)
		} else {
			println("ok")
		}
	}
}
