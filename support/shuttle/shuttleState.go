package shuttle

import (
	"database/sql"
	"encoding/json"
	"io"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func NewShuttle(inFile string, opts qna.Options) (ret Shuttle, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else if e := tables.CreateRun(db); e != nil {
		err = e
	} else if query, e := qdb.NewQueries(db, true); e != nil {
		err = e
	} else {
		ret = Shuttle{
			inFile:  inFile,
			db:      db,
			query:   query,
			opts:    opts,
			decoder: decode.NewDecoder(tapestry.AllSignatures),
			play:    nil,
		}
	}
	return
}

func (c *Shuttle) Post(w io.Writer, endpoint string, msg json.RawMessage) (err error) {
	switch endpoint {

	case Restart.String():
		var scene string
		if e := json.Unmarshal(msg, &scene); e != nil {
			err = errutil.New("invalid scene")
		} else if _, e := c.Restart(scene); e != nil {
			err = e
		} else {
			evts := c.out.GetEvents()
			err = writeFrames(w, []frame.Frame{{Events: evts}})
		}

	case Query.String():
		var qs []json.RawMessage
		if e := json.Unmarshal(msg, &qs); e != nil {
			err = errutil.New("invalid query")
		} else if play, e := c.EnsurePlay(); e != nil {
			err = e
		} else {
			// a series of commands:
			var frames []frame.Frame // output
			for _, q := range qs {
				var f frame.Frame
				if a, e := c.decoder.DecodeAssignment(affine.None, q); e != nil {
					f.Error = e.Error()
				} else if v, e := a.GetAssignedValue(play); e != nil {
					f.Error = e.Error()
				} else {
					f.Result = debug.Stringify(v)
				}
				f.Events = c.out.GetEvents() // even on error
				frames = append(frames, f)
			}
			if err == nil {
				err = writeFrames(w, frames)
			}
		}

	default:
		err = errutil.New("Unknown endpoint", endpoint)
	}
	return
}

func writeFrames(w io.Writer, frames frame.Frame_Slice) (err error) {
	if have, e := cout.Marshal(&frames, nil); e != nil {
		err = e
	} else {
		_, err = w.Write([]byte(have))
	}
	return
}
