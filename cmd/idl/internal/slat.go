package internal

import (
	r "reflect"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
)

type Slat struct {
	Type                    r.Type // a struct
	Name, Lede, Group, Desc string
	Places                  []Place
	Slots                   []*Slot
}

func (c *Slat) FlowText() string {
	var out strings.Builder
	out.WriteString(lang.Underscore(c.Lede))
	for _, p := range c.Places {
		if !p.Internal {
			str := p.FlowText()
			out.WriteRune(' ')
			out.WriteRune('{')
			out.WriteString(str)
			out.WriteRune('}')
		}
	}
	return out.String()
}

func (c *Slat) SlotText() string {
	var out strings.Builder
	cnt := len(c.Slots)
	if cnt != 1 {
		out.WriteRune('[')
	}
	for i := 0; i < cnt; i++ {
		if i > 0 {
			out.WriteString(", ")
		}
		n := lang.Underscore(c.Slots[i].Name)
		out.WriteRune('"')
		out.WriteString(n)
		out.WriteRune('"')

	}
	if cnt != 1 {
		out.WriteRune(']')
	}
	return out.String()
}

func (f *Slat) Format() (ret string) {
	return "%-20s %-20s = %3d"
}

func (f *Slat) Camel() (ret string) {
	return Camel(f.Lede)
}

func (c *Slat) Package() string {
	return PackageOf(c.Type)
}

func (c *Slat) Signatures() []Sig {
	// determines the total number of signature permutations
	perms := 1
	for _, p := range c.Places {
		perms *= p.Choices()
	}
	sigs := make([]Sig, perms)
	pascalLede := Pascal(c.Lede)
	camelLede := Camel(c.Lede)
	//
	for k := 0; k < perms; k++ {
		//
		var nsig, psig, csig strings.Builder
		nsig.WriteString(c.Lede)
		psig.WriteString(pascalLede)
		csig.WriteString(camelLede)
		//
		rem := k
		var runIn bool
		for _, p := range c.Places {
			if !p.Internal {
				if mod := rem % p.Choices(); mod == 0 {
					//
					if l := p.Label; len(l) == 0 {
						psig.WriteRune(':')
						nsig.WriteString(strconv.Itoa(p.Index))
						runIn = true
					} else {
						if !runIn {
							psig.WriteRune(' ')
							nsig.WriteRune('0')
							runIn = true
						}

						csig.WriteString(Pascal(l))
						psig.WriteString(Camel(l))
						psig.WriteRune(':')
						//
						nsig.WriteRune('_')
						nsig.WriteString(l)
						nsig.WriteString(strconv.Itoa(p.Index))
					}
				}
			}
			rem = rem / p.Choices()
		}
		sigs[k] = Sig{
			Numbered: nsig.String(),
			Raw:      psig.String(),
			Camel:    csig.String(),
			Type:     c.Name,
			Package:  c.Package(),
		}
	}
	return sigs
}
