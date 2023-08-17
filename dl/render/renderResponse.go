package render

import (
	"errors"
	"io"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// prints the response via the runtime's writer.
func (op *RenderResponse) Execute(run rt.Runtime) (err error) {
	if e := op.printResponse(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

// return the rendered response as a text value.
func (op *RenderResponse) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getResponse(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderResponse) printResponse(run rt.Runtime) (err error) {
	if v, e := op.getResponse(run); e != nil {
		err = e
	} else if w := run.Writer(); w == nil {
		err = errutil.New("missing writer")
	} else {
		_, e := io.WriteString(w, v.String())
		err = e
	}
	return
}

func (op *RenderResponse) getResponse(run rt.Runtime) (ret g.Value, err error) {
	var unknown g.Unknown
	if name := op.Name; len(name) == 0 {
		// and unnamed response
		ret, err = op.getLocalText(run)
	} else if v, e := run.GetField(meta.Response, name); e == nil {
		// try to look up the name from the global replacement table
		ret = v
	} else if errors.As(e, &unknown); unknown.Target != meta.Response {
		// some error that wasn't a failed replacement....
		err = e
	} else if op.Text == nil {
		// todo: once warnings are implemented instead of errors
		// this could return the response name instead.
		err = errutil.Fmt("%w and no fallback specified", g.UnknownResponse(name))
	} else {
		ret, err = op.getLocalText(run)
	}
	return
}

func (op *RenderResponse) getLocalText(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetText(run, op.Text)
}
