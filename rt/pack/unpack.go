package pack

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// read from json structured as a tapestry command { Record: [] }
func UnpackRecord(ks rt.Kinds, b []byte, expectedType string) (ret *rt.Record, err error) {
	if len(b) > 0 {
		var msg map[string][]json.RawMessage
		if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else {
			for sig, raw := range msg {
				if ret != nil {
					err = errors.New("unexpected fields")
					break
				} else if parts, e := decode.DecodeSignature(sig); e != nil {
					err = e
					break
				} else if cmd := parts[0]; cmd != expectedType {
					err = fmt.Errorf("couldn't unpack record of %q into %q", cmd, expectedType)
					break
				} else if rec, e := safe.NewRecord(ks, cmd); e != nil {
					err = e
					break
				} else {
					var lastField int
					for part, field := range parts[1:] {
						if ft, idx := findField(rec.Kind, field, lastField); idx < 0 {
							err = fmt.Errorf("couldnt find field %q", field)
							break
						} else if v, e := UnpackValue(ks, raw[part], ft.Affinity, ft.Type); e != nil {
							err = e
							break
						} else if e := rec.SetIndexedField(idx, v); e != nil {
							err = e
							break
						} else {
							lastField = idx
						}
					}
					if err == nil {
						ret = rec
						continue
					}
				}
			}
		}
	}
	return
}

func findField(k *rt.Kind, name string, startingAt int) (retField rt.Field, retIndex int) {
	retIndex = -1 // provisionally
	for i, cnt := startingAt, k.FieldCount(); i < cnt; i++ {
		if ft := k.Field(i); ft.Name == name {
			retField, retIndex = ft, i
			break
		}
	}
	return
}

func UnpackValue(ks rt.Kinds, b []byte, a affine.Affinity, t string) (ret rt.Value, err error) {
	switch a {
	case affine.Bool:
		var v bool
		if e := json.Unmarshal(b, &v); e != nil {
			err = e
		} else {
			ret = rt.BoolFrom(v, t)
		}
	case affine.Num:
		var v float64
		if e := json.Unmarshal(b, &v); e != nil {
			err = e
		} else {
			ret = rt.FloatFrom(v, t)
		}
	case affine.Text:
		var v string
		if e := json.Unmarshal(b, &v); e != nil {
			err = e
		} else {
			ret = rt.StringFrom(v, t)
		}
	case affine.Record:
		if rec, e := UnpackRecord(ks, b, t); e != nil {
			err = e
		} else if rec != nil {
			ret = rt.RecordOf(rec)
		}
	case affine.NumList:
		var v []float64
		if e := json.Unmarshal(b, &v); e != nil {
			err = e
		} else {
			ret = rt.FloatsFrom(v, t)
		}
	case affine.TextList:
		var v []string
		if e := json.Unmarshal(b, &v); e != nil {
			err = e
		} else {
			ret = rt.StringsFrom(v, t)
		}
	case affine.RecordList:
		var raw []json.RawMessage
		if e := json.Unmarshal(b, &raw); e != nil {
			err = e
		} else {
			vs := make([]*rt.Record, len(raw))
			for i := range raw {
				if rec, e := UnpackRecord(ks, raw[i], t); e != nil {
					err = e
					break
				} else {
					vs[i] = rec
				}
			}
			if err == nil {
				ret = rt.RecordsFrom(vs, t)
			}
		}
	default:
		err = fmt.Errorf("unhandled affinity %s", a.String())
	}
	return
}
