package cmdserve

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/web"
)

// exported to package main in cmd/tap
var CmdServe = &base.Command{
	Run:       goServe,
	Flag:      buildFlags(),
	UsageLine: "tap serve [-in dbpath]",
	Short:     "serve a story through http",
	Long:      `Run a server that plays a scene from a previously built story database.`,
}

func goServe(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if lvl, ok := debug.MakeLoggingLevel(cfg.logLevel); !ok {
		err = fmt.Errorf("Unknown log level %s", cfg.logLevel)
	} else {
		listenTo, _ := cfg.listen.GetPort(8080)
		requestFrom, useBrowser := cfg.request.GetPort(3000)
		log.Println("using story files from:", cfg.inFile)
		log.Println("listening to:", listenTo)
		if !useBrowser {
			requestFrom = 0
		} else {
			log.Println("don't forget to run the vite server at port " + strconv.Itoa(requestFrom) + "!")
			log.Println("in the tapestry/www directory run 'npm run dev'")
			log.Printf("then browse to: http://localhost:%d/play/\n", listenTo)
		}
		debug.LogLevel = lvl
		opts := qna.NewOptions()
		err = serveWithOptions(cfg.inFile, opts, listenTo, requestFrom)
	}
	return
}

var cfg = struct {
	inFile          string
	debugging       bool
	logLevel        string
	listen, request web.Port
}{}

// creates a description which writes into the cfg when the base.Command is matched
func buildFlags() (ret flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	ret.StringVar(&cfg.inFile, "in", inFile, "input file name (sqlite3)")
	ret.StringVar(&cfg.logLevel, "log", debug.C_LoggingLevel_Info.String(), levels)
	ret.Var(&cfg.listen, "listen", "the port for your web browser. specify a port number; or, 'true' for the default (8080).")
	ret.Var(&cfg.request, "www", "local vite server where tapestry can find its webapps. specify a port number; or, 'true' to use the default port (3000).")
	return
}
