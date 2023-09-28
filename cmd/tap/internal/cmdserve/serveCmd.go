package cmdserve

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
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
	if lvl, e := getLogLevel(cfg.logLevel); e != nil {
		err = e
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
		if cnt, e := serveWithOptions(cfg.inFile, opts, listenTo, requestFrom); e != nil {
			errutil.PrintErrors(e, func(s string) { log.Println(s) })
			if errutil.Panic {
				log.Panic("mismatched")
			}
		} else {
			log.Println("done", cnt, cfg.inFile)
		}
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
func buildFlags() (out flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	spec := debug.LogLevel.Compose()
	levels := strings.Join(spec.Strings, ", ")
	out.StringVar(&cfg.inFile, "in", inFile, "input file name (sqlite3)")
	out.StringVar(&cfg.logLevel, "log", "", levels)
	out.Var(&cfg.listen, "listen", "the port for your web browser. specify a port number; or, 'true' for the default (8080).")
	out.Var(&cfg.request, "www", "local vite server where tapestry can find its webapps. specify a port number; or, 'true' to use the default port (3000).")
	return
}

// translate command line specified debug level into the runtime's debug level string
func getLogLevel(in string) (ret debug.LoggingLevel, err error) {
	if len(in) > 0 {
		spec := ret.Compose()
		if key, idx := spec.IndexOfValue(in); idx < 0 {
			err = errutil.New("Unknown log level")
		} else {
			ret.Str = key
		}
	}
	return
}