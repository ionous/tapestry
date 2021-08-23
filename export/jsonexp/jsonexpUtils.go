package jsonexp

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// UnmarshalDetailedSlot - helper for use with code generated unmarshalling,
func UnmarshalDetailedSlot(n Context, b []byte) (ret interface{}, err error) {
	if len(b) > 0 {
		var msg Node
		if e := json.Unmarshal(b, &msg); e != nil {
			err = errutil.New("unmarshaling", e)
		} else if contents := msg.Value; len(contents) > 0 {
			var inner Node // peek to create the appropriate type
			if e := json.Unmarshal(contents, &inner); e != nil {
				err = errutil.New("unpacking", e)
			} else if ptr, e := n.NewType(msg.Type, inner.Type); e != nil {
				err = errutil.New("creating", e)
			} else if imp, ok := ptr.(DetailedMarshaler); !ok {
				err = errutil.New("casting", e)
			} else if e := imp.UnmarshalDetailed(n, contents); e != nil {
				err = errutil.New("reading", e)
			} else if fini, e := n.Finalize(ptr); e != nil {
				err = errutil.New("finalizing", e)
			} else {
				ret = fini
			}
		}
	}
	return
}

// UnmarshalCompactSlot - helper for use with code generated unmarshalling,
func UnmarshalCompactSlot(n Context, b []byte) (ret interface{}, err error) {
	if len(b) > 0 {
		var msg Node
		if e := json.Unmarshal(b, &msg); e != nil {
			err = errutil.New("unmarshaling", e)
		} else if contents := msg.Value; len(contents) > 0 {
			var inner Node // peek to create the appropriate type
			if e := json.Unmarshal(contents, &inner); e != nil {
				err = errutil.New("unpacking", e)
			} else if ptr, e := n.NewType(msg.Type, inner.Type); e != nil {
				err = errutil.New("creating", e)
			} else if imp, ok := ptr.(CompactMarshaler); !ok {
				err = errutil.New("casting", e)
			} else if e := imp.UnmarshalCompact(n, contents); e != nil {
				err = errutil.New("reading", e)
			} else if fini, e := n.Finalize(ptr); e != nil {
				err = errutil.New("finalizing", e)
			} else {
				ret = fini
			}
		}
	}
	return
}
