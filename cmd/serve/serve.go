// Runs external applications, and pipes http input to stdin; stdout to http output.
package main

// go run serve.go -in /Users/ionous/Documents/Tapestry/shared -out /Users/ionous/Documents/Tapestry/build/play.db
// go run serve.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	// var srcPath, outFile string
	// var check, open bool
	// flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	// flag.StringVar(&outFile, "out", "", "output filename (sqlite3)")
	// flag.BoolVar(&check, "check", true, "run check after importing?")
	// flag.BoolVar(&open, "open", false, "open a new browser window.")
	// flag.Parse()

	// if len(srcPath) == 0 || len(outFile) == 0 {
	// 	flag.PrintDefaults()
	// 	return
	// }

	// cfg := web.DevConfig(build.Default.GOPATH, "")
	// if open {
	// 	support.OpenBrowser(web.Endpoint(8088, "localhost", "live"))
	// }

	// // sub-process communication
	// cs := serve.NewChannels()

	// // assemble and play ( reads from and writes to channels )
	// go func() {
	// 	cs.ChangeMode(play.PlayModes_Asm)
	// 	log.Println("Assembling", srcPath+"...")
	// 	if dbPath, e := serve.Asm(cfg.Cmd("asm"), srcPath, outFile, check, cs); e != nil {
	// 		println(e.Error())
	// 		cs.Fatal(e)
	// 	} else {
	// 		log.Println("Playing", dbPath+"...")
	// 		cs.ChangeMode(play.PlayModes_Play)
	// 		if e := serve.Play(cfg.Cmd("play"), dbPath, cs); e != nil {
	// 			println(e.Error())
	// 			cs.Fatal(e)
	// 		} else {
	// 			log.Println("Done.")
	// 			cs.ChangeMode(play.PlayModes_Complete)
	// 		}
	// 	}
	// }()

	// // the server might last longer than the processes so let it block
	// log.Fatal(serve.ListenAndServe(web.Endpoint(8088), cs))
}
