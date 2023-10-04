package shuttle

import (
	"encoding/json"
	"io"

	"github.com/ionous/errutil"
)

// read a request from the client; return a string or map.
// see Play.vue... example: io.send({in: txt});
func Decode(r io.Reader) (ret any, err error) {
	var msg struct {
		Input   string         `json:"in"`
		Command map[string]any `json:"cmd"`
	}
	dec := json.NewDecoder(r)
	// decode an array value (Message)
	if e := dec.Decode(&msg); e != nil {
		err = e
	} else if len(msg.Command) > 0 {
		ret = msg.Command
	} else if len(msg.Input) > 0 {
		ret = msg.Input
	} else {
		err = errutil.New("unknown or empty input")
	}
	return
}
