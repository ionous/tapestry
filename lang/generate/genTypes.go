package generate

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"github.com/mitchellh/go-wordwrap"
)

// common data for all types
type specData struct {
	Name, Idl string
	Markup    markup
}

func (t specData) Link() (ret string, err error) {
	if a, ok := hackForLinks.linkByName(t.Name); !ok {
		err = fmt.Errorf("unknown type %q creating link", t.Name)
	} else {
		ret = a
	}
	return
}

func (d specData) Comments() []string {
	return d.Markup.Comments()
}

// comment as a single json friendly
func (d specData) SchemaComment() (string, error) {
	return d.Markup.SchemaComment()
}

// because references to types arent scoped but the generated code needs to be:
// the generator has to load all possible types before writing them out.
type flowData struct {
	specData
	Lede  string
	Slots []string
	Terms []termData
}

// signatures along with the terms those signatures use.
func (d flowData) Signatures() []sigTerm {
	return sigTerms(d)
}

type markup map[string]any

// fix: cache to stop the multiple creation?
func (m markup) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "pos" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func (m markup) Comments() (ret []string) {
	if vs, e := compact.ExtractComment(m); e != nil {
		log.Panic(e)
	} else {
		ret = vs
	}
	return
}

// comment as a single json friendly
func (m markup) SchemaComment() (ret string, err error) {
	str := strings.Join(m.Comments(), " ")
	if b, e := json.Marshal(str); e != nil || len(b) == 0 {
		err = e
	} else {
		ret = string(b[1 : len(b)-1])
		ret = wordwrap.WrapString(ret, 50)
		ret = strings.ReplaceAll(ret, "\n", ` <br>`)
	}
	return
}

type slotData struct {
	specData
}

type strData struct {
	specData
	Options        []string
	OptionComments []string
}

type numData struct {
	specData
}

type termData struct {
	Name, Label, Type          string
	Private, Optional, Repeats bool
	Markup                     markup
}

func (t termData) TypeScope() (ret string, err error) {
	if group, ok := hackForLinks.findGroup(t.Type); !ok {
		err = fmt.Errorf("unknown type %q creating link", t.Type)
	} else {
		ret = group + "." + t.Type
	}
	return
}

func (t termData) Link() (ret string, err error) {
	if a, ok := hackForLinks.linkByName(t.Type); !ok {
		err = fmt.Errorf("unknown type %q creating link", t.Type)
	} else {
		ret = a
	}
	return
}

// comment as a single json friendly
func (t termData) SchemaComment() (string, error) {
	return t.Markup.SchemaComment()
}

// handle transforming _ into a blank string
func (t termData) SimpleLabel() (ret string) {
	if t.Label != "_" {
		ret = t.Label
	}
	return
}

// shadows the typeinfo.T interface
// these implementations exist to created the data for that interface.
type typeData interface {
	getName() string
	goType() string
	getMarkup() markup
}

func (f specData) getName() string {
	return f.Name
}

func (f specData) getMarkup() markup {
	return f.Markup
}

func (f specData) goType() string {
	return Pascal(f.Name)
}

func (f numData) goType() string {
	return "float64"
}

func (f strData) goType() (ret string) {
	if len(f.Options) == 0 {
		ret = "string"
	} else {
		switch f.Name {
		case "bool":
			ret = "bool"
		default:
			ret = Pascal(f.Name)
		}
	}
	return
}
