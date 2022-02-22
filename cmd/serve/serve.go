package main

import (
	"flag"
	"go/build"
	"log"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/dl/play"

	live "git.sr.ht/~ionous/tapestry/cmd/live/internal"
)

// go run live.go -in /Users/ionous/Documents/Tapestry/stories/shared -out /Users/ionous/Documents/Tapestry/build/play.db
// go run live.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	cfg := composer.DevConfig(build.Default.GOPATH)
	//
	var srcPath, outFile string
	var check bool
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "output filename (sqlite3)")
	flag.BoolVar(&check, "check", true, "run check after importing?")
	flag.Parse()

	if len(srcPath) == 0 || len(outFile) == 0 {
		flag.PrintDefaults()
		return
	}

	// sub-process communication
	cs := live.NewChannels()

	// assemble and play ( reads from and writes to channels )
	go func() {
		cs.ChangeMode(play.PlayModes_Asm)
		log.Println("Assembling", srcPath+"...")
		if dbPath, e := live.Asm(cfg.Assemble, srcPath, outFile, check, cs); e != nil {
			println(e.Error())
			cs.Fatal(e)
		} else {
			log.Println("Playing", dbPath+"...")
			cs.ChangeMode(play.PlayModes_Play)
			if e := live.Play(cfg.Play, dbPath, cs); e != nil {
				println(e.Error())
				cs.Fatal(e)
			} else {
				log.Println("Done.")
				cs.ChangeMode(play.PlayModes_Complete)
			}
		}
	}()

	// the server might last longer than the processes so let it block
	log.Fatal(live.ListenAndServe(":8088", cs))
}
