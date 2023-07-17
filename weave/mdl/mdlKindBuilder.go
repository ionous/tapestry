package mdl

import (
	"git.sr.ht/~ionous/tapestry/lang"
)

type KindBuilder struct {
	kinds GetRequiredKind
	Kind
}

type Kind struct {
	name    string
	parent  string
	fs      fieldSet
	classes classCache
}

func (p *Kind) Name() string {
	return p.name
}

func (p *Kind) Parent() string {
	return p.parent
}

func NewKindBuilder(kinds GetRequiredKind, name string, parent string) *KindBuilder {
	// if (parent & kindsOf.Kind) == 0 {
	// 	panic("subtype not a kind")
	// }
	var classes classCache
	classes.addClass(kinds, parent)

	return &KindBuilder{kinds: kinds, Kind: Kind{
		// tbd: feels like it'd be best to have spec flag names that need normalization,
		// and convert all the names at load time ( probably storing the original somewhere )
		// ( ex. store the normalized names in the meta data )
		name:    lang.Normalize(name),
		parent:  parent,
		classes: classes,
	}}
}

// defers execution; so no return value.
func (b *KindBuilder) AddField(fn FieldInfo) {
	b.classes.addClass(b.kinds, fn.Class)
	b.fs.fields[PatternLocals] = append(b.fs.fields[PatternLocals], fn)
}

func (p *Kind) write(m *Pen) (ret KindInfo, err error) {
	// if parent, e := p.classes.getClass(p.parent); e != nil {
	// 	err = e
	// } else if kid, e := m.AddKindById(parent, p.name); e != nil {
	// 	err = e // tbd: AddKind can be reissued multiple times... cache via a resource?
	// } else if !kid.HasParent(parent) {
	// 	err = errutil.New("invalid parent %q of %q", parent.Name, kid.Name)
	// } else if e := p.fs.write(m, p.classes, kid); e != nil {
	// 	err = e
	// } else {
	// 	ret = kid
	// }
	panic("not implemented")
}
