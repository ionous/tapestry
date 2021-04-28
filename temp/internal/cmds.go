package internal

import (
	r "reflect"
	"strconv"
	"strings"
)

type Cmd struct {
	Type                    r.Type // a struct
	Name, Lede, Group, Desc string
	Places                  []Place
}

func (c *Cmd) Signatures() []Sig {
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
		}
	}
	return sigs
}
