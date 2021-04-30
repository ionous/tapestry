package internal

import (
	r "reflect"
	"unicode"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/export"
	"git.sr.ht/~ionous/iffy/export/tag"
	"git.sr.ht/~ionous/iffy/lang"
)

type Place struct {
	Cmd                          *Cmd // parent
	Label                        string
	Arg                          string
	Index                        int
	Repeated, Optional, Internal bool
	Pool                         string // the string field pulls from a set of similar names
	// for flows, these would become an array;
	// these would be arrays; optional a nil type.
	Type r.Type
}

func (p *Place) Choices() (ret int) {
	if !p.Optional {
		ret = 1
	} else {
		ret = 2
	}
	return
}

//
func (p *Place) ProtoQualifier() (ret string) {
	if p.Repeated {
		ret = "repeated"
	} else {
		ret = "        "
	}
	return
}

func (p *Place) Camel() string {
	return Camel(p.Arg)
}
func (p *Place) CapLabel() (ret string) {
	return p.Label // an empty label means anonymous runin
}

func (p *Place) CapType() (ret string) {
	switch n := p.Type.Name(); n {
	// fix? not exactly sure why these dont expand their content
	case "Relation":
		ret = "Text"
	case "string":
		ret = "Text"
	case "Position":
		ret = "Pos"
	case "Case", "Edge", "Order": // list case, list edge, list order
		ret = "Bool"
	case "TryAsNoun", "Level": // debug level
		ret = "Int32"
	default:
		if len(n) == 0 {
			panic("missing name")
		} else {
			ret = string(unicode.ToUpper(rune(n[0]))) + n[1:]
		}
	}
	//
	pack := p.Cmd.Package()
	if n := PackageOf(p.Type); len(n) > 0 && n != pack {
		ret = Pascal(n) + "." + ret
	}

	if p.Repeated {
		ret = "List(" + ret + ")"
	}
	return
}

func (p *Place) ProtoType() (ret string) {
	switch n := p.Type.Name(); n {
	case "float64":
		ret = "double"
	case "Position":
		ret = "Pos"
	case "PatternName", "EventName":
		ret = "string"
	case "Edge", "Order", "Case":
		ret = "bool"
	case "TryAsNoun":
		ret = "int32"
	case "Level":
		ret = "int32"
	default:
		ret = n
	}
	return
}

type Cmds map[string]*Cmd

func (cs *Cmds) Add(p *Cmd) *Cmd {
	(*cs)[p.Name] = p
	return p
}

func MakeCommand(c composer.Composer) (ret *Cmd) {
	rtype := r.TypeOf(c).Elem()
	if rtype.Kind() == r.Struct {
		cmd := Cmd{
			Name:  rtype.Name(),
			Type:  rtype,
			Lede:  makeLede(c),
			Group: makeGroup(c),
			Desc:  Desc(c),
		}
		//
		var inds Indicies
		export.WalkProperties(rtype, func(f *r.StructField, path []int) (done bool) {
			if inner := collapse(f.Type); inner != nil {
				f = inner
			}
			tags := tag.ReadTag(f.Tag)
			optional := tags.Exists("optional")
			internal := tags.Exists("internal")
			label, pool := makeLabel(f, tags, inds.pub)
			var repeated bool
			el := f.Type
			for {
				if k := el.Kind(); k == r.Ptr {
					el = el.Elem()
				} else if k == r.Slice {
					el = el.Elem()
					repeated = true
				} else {
					// unpack typedefs
					switch k {
					case r.String:
						el = r.TypeOf("")
					case r.Int32, r.Int:
						el = r.TypeOf(int32(0))
					}
					break
				}
			}
			place := Place{
				Cmd:      &cmd,
				Label:    label,
				Arg:      lang.Underscore(f.Name),
				Index:    inds.makeIndex(f, tags),
				Repeated: repeated,
				Optional: optional,
				Internal: internal,
				Pool:     pool,
				Type:     el,
			}
			cmd.Places = append(cmd.Places, place)
			return
		})
		ret = &cmd
	}
	return
}

// should we collapse the target's field(s) directly into its parents?
func collapse(el r.Type) (ret *r.StructField) {
	if el.Kind() == r.Ptr {
		el = el.Elem()
	}
	if el.Kind() == r.Struct {
		export.WalkProperties(el, func(f *r.StructField, path []int) (done bool) {
			tags := tag.ReadTag(f.Tag)
			if !tags.Exists("internal") {
				if ret == nil {
					ret = f
				} else {
					ret = nil
					done = true
				}
			}
			return
		})
	}
	return
}
