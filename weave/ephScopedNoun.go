package weave

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	Requires     // kinds ( when resolved, can have one direct parent )
	domain       *Domain
	names        UniqueNames
	friendlyName string
	aliases      UniqueNames
	aliasat      []string    // origin of each alias
	localRecord  localRecord // store the values of the noun as a record.
}

func (n *ScopedNoun) Resolve() (ret Dependencies, err error) {
	if len(n.at) == 0 {
		err = NounError{n.name, errutil.New("never defined")}
	} else if ks, e := n.resolve(n, (*kindFinder)(n.domain)); e != nil {
		err = NounError{n.name, e}
	} else {
		ret = ks
	}
	return
}

func (n *ScopedNoun) Kind() (ret *ScopedKind, err error) {
	if dep, e := n.GetDependencies(); e != nil {
		err = e
	} else if ks := dep.Parents(); len(ks) != 1 {
		err = errutil.Fmt("noun %q has unexpected %d parents", n.name, len(ks))
	} else {
		ret = ks[0].(*ScopedKind)
	}
	return
}

// returns false if the alias already existed
func (n *ScopedNoun) AddAlias(a, at string) (okay bool) {
	if i := n.aliases.AddName(a); i >= 0 {
		s := append(n.aliasat, "")
		copy(s[i+1:], s[i:])
		n.aliasat, s[i] = s, at
		okay = true
	}
	return
}

// stores literal values because they are serializable
// ( as opposed to generic values which aren't. )
func (n *ScopedNoun) recordValues(at string) (ret localRecord, err error) {
	if n.localRecord.isValid() {
		ret = n.localRecord
	} else if k, e := n.Kind(); e != nil {
		err = e
	} else {
		rv := localRecord{k, new(literal.RecordValue), at}
		ret, n.localRecord = rv, rv
	}
	return
}

func (n *ScopedNoun) Names() []string {
	if len(n.names) == 0 {
		n.names = n.makeNames()
	}
	return n.names
}

func (n *ScopedNoun) UpdateFriendlyName(name string) {
	if len(n.friendlyName) == 0 {
		if clip := strings.TrimSpace(name); clip != n.name {
			n.friendlyName = name
		}
	}
}

func (n *ScopedNoun) makeNames() []string {
	var out []string
	// the ranked 0 name is used for default display when printing nouns
	// (ex. "toy boat")
	defaultName := n.name
	if alt := n.friendlyName; len(alt) > 0 {
		defaultName = alt
	}
	out = append(out, defaultName)
	// write the normalized name if it was different
	if defaultName != n.name {
		out = append(out, n.name)
	}
	// now generate additional names by splitting the lowercase uniform name on the underscores:
	split := strings.FieldsFunc(n.name, lang.IsBreak)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, "_"); breaks != n.name {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := split[i]
			out = append(out, word)
		}
	}
	return out
}
