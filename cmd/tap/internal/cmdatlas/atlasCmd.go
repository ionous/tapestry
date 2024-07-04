package cmdatlas

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web"
)

func runAtlas(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if inFile, e := filepath.Abs(cfg.inFile); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else if db, e := tables.OpenModel(inFile); e != nil {
		err = e
	} else {
		defer db.Close()
		listenTo, _ := cfg.listen.GetPort(defaultPort)
		endpoint := fmt.Sprintf("/%s/", defaultEndpoint)
		log.Printf("listening to %d at %s\n", listenTo, endpoint)
		mux := http.NewServeMux()
		newAtlas(mux, &atlasContext{
			db, query.NewDecoder(tapestry.AllSignatures),
		})
		// block forever ish.
		// note: on windows the localhost is required in order to avoid the windows firewall popup
		where := "localhost:" + strconv.Itoa(listenTo)
		err = http.ListenAndServe(where, mux)
	}
	return
}

const defaultPort = 8100
const defaultEndpoint = "atlas" // becomes /atlas/

var CmdAtlas = &base.Command{
	Run:       runAtlas,
	Flag:      buildFlags(),
	UsageLine: "tap atlas",
	Short:     "r/o database server.",
	Long: `
An http server which returns json representations of the database.`,
}

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	inFile string
	listen web.Port
}{}

func buildFlags() (ret flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	ret.StringVar(&cfg.inFile, "in", inFile, "input file name (sqlite3)")
	ret.Var(&cfg.listen, "listen",
		fmt.Sprintf("the port for your web browser. specify a port number; or, 'true' for the default (%d).", defaultPort))
	return
}
