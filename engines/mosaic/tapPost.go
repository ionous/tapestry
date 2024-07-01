package mosaic

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ionous/errutil"
)

// allows posts to / of specific predetermined commands ( in json )
func HandleCommands(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		} else if r.URL.Path != "/" {
			http.NotFound(w, r)
		} else if e := handlePost(cfg, r.Body); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		} else {
			// note: writing to output would automatically write StatusOk
			// for now... an empty response.
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
		}
	}
}

// when the server was separate apps:
// still will need something to run commands...
func handlePost(cfg *Config, body io.Reader) (err error) {
	// read a request from the client
	// see Mosaic.vue onPlay which sends {play:true}
	var cmd cmdFromClient
	dec := json.NewDecoder(body)
	if e := dec.Decode(&cmd); e != nil {
		err = e
	} else if !cmd.Play {
		err = errutil.New("unknown request")
	} else {
		err = errutil.New("not handled")
		// okay --
		// so previously, this would open an entirely new browser window
		// and would launch a server dedicated to it.
		// options
		// 1. launch a new instance of tapestry  ( on a random temp port, opening a browser depending on the mode )
		// 2. within tapestry, browse to a new endpoint.
		//    god it'd be cool if it could have tabs, but ... no. that'd have to be another process
		//    for this, you'd have to have tabs across the top like links
		//  you'd have

		// 	cmd := cfg.Cmd("serve")
		// 	args := []string{
		// 		"-in", cfg.PathTo("stories", "shared"),
		// 		"-out", cfg.Scratch("play.db"),
		// 		"-open",
		// 	}
		// 	log.Println("playing", cmd, args)
		// 	go func() {
		// 		cmd := exec.Command(cmd, args...)
		// 		if e := cmd.Run(); e != nil {
		// 			log.Println(e)
		// 		}
		// 	}()
	}
	return
}

type cmdFromClient struct {
	Play bool `json:"play"`
}
