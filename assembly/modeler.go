package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
	"github.com/reiver/go-porterstemmer"
)

type IssueReport func(pos reader.Position, msg string)

func cat(str ...string) string {
	return strings.Join(str, " ")
}

func NewAssembler(db *sql.DB) *Assembler {
	reportNothing := func(reader.Position, string) {}
	return NewAssemblerReporter(db, reportNothing)
}

func NewAssemblerReporter(db *sql.DB, report IssueReport) *Assembler {
	return &Assembler{tables.NewCache(db), report, 0}
}

type Assembler struct {
	cache      *tables.Cache
	issueFn    IssueReport
	IssueCount int
}

func (m *Assembler) reportIssue(src, ofs, msg string) {
	pos := reader.Position{Source: src, Offset: ofs}
	m.issueFn(pos, msg)
	m.IssueCount++
}
func (m *Assembler) reportIssuef(src, ofs, fmt string, args ...interface{}) {
	m.reportIssue(src, ofs, errutil.Sprintf(fmt, args...))
}

// write kind and comma separated ancestors
func (m *Assembler) WriteAncestor(kind, path string) (err error) {
	_, e := m.cache.Exec(mdl_kind, kind, path)
	return e
}

func (m *Assembler) WriteCheck(name, testType, expect string) error {
	_, e := m.cache.Exec(mdl_check, name, testType, expect)
	return e
}

func (m *Assembler) WriteField(kind, field, fieldType, aff string) error {
	// patch till deeper fixes.
	if len(aff) == 0 && fieldType == "aspect" {
		aff = string(affine.Text)
	}
	_, e := m.cache.Exec(mdl_field, kind, field, fieldType, aff)
	return e
}

func DomainNameOf(domain, noun string) string {
	var b strings.Builder
	b.WriteRune('#')
	if len(domain) > 0 {
		b.WriteString(domain)
		b.WriteString("::")
	}
	b.WriteString(lang.Breakcase(noun))
	return b.String()
}

func (m *Assembler) WriteNoun(noun, kind string) error {
	_, e := m.cache.Exec(mdl_noun, noun, kind)
	return e
}

// WriteName for noun
func (m *Assembler) WriteName(noun, name string, rank int) error {
	_, e := m.cache.Exec(mdl_name, noun, name, rank)
	return e
}

// WriteNounWithNames writes the noun to the model,
// and splits the name into separate space separated words.
// each word is recorded as a possible reference to the noun.
func (m *Assembler) WriteNounWithNames(domain, noun, kind string) (err error) {
	id := DomainNameOf(domain, noun)
	if e := m.WriteNoun(id, kind); e != nil {
		err = errutil.Append(err, e)
	} else {
		// COUNTER:#
		// if counter := strings.Index(noun, "#"); counter > 0 {
		// 	noun = noun[:counter]
		// }
		// we want to divide words on breakcase boundaries
		breaks := lang.Breakcase(noun)
		split := strings.FieldsFunc(breaks, lang.IsBreak)
		spaces := strings.Join(split, " ")

		// the ranked 0 name is used for default display when printing nouns
		var ofs int
		if e := m.WriteName(id, spaces, ofs); e != nil {
			err = errutil.Append(err, e)
		} else if cnt := len(split); cnt > 1 {
			if spaces != breaks {
				ofs++
				if e := m.WriteName(id, breaks, ofs); e != nil {
					err = errutil.Append(err, e)
				}
			}

			// write the individual words of the split ( ex. toy, boat )
			for i, k := range split {
				rank := cnt + ofs - i
				if e := m.WriteName(id, strings.ToLower(k), rank); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}

	}
	return
}

func (m *Assembler) WritePattern(name, result string, labels []string) error {
	ls := strings.Join(labels, ",")
	_, e := m.cache.Exec(mdl_pat, name, result, ls)
	return e
}

func (m *Assembler) WritePlural(one, many string) error {
	_, e := m.cache.Exec(mdl_plural, one, many)
	return e
}

func (m *Assembler) WriteProg(progName, typeName string, bytes []byte) (err error) {
	_, err = m.cache.Exec(mdl_prog, progName, typeName, bytes)
	return
}

func (m *Assembler) EncodeValue(rval r.Value) (ret []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if e := enc.EncodeValue(rval); e != nil {
		err = errutil.New("error encoding value", e)
	} else {
		ret = buf.Bytes()
	}
	return

}
func (m *Assembler) WriteGob(progName string, cmd interface{}) (err error) {
	rval := r.ValueOf(cmd)
	if prog, e := m.EncodeValue(rval); e != nil {
		err = errutil.New("WriteGob, error encoding value", e)
	} else {
		typeName := rval.Elem().Type().Name()
		err = m.WriteProg(progName, typeName, prog)
	}
	return
}

func (m *Assembler) WriteGobs(gobs map[string]interface{}) (err error) {
	for k, v := range gobs {
		if e := m.WriteGob(k, v); e != nil {
			err = e
			break
		}
	}
	return
}

func (m *Assembler) WriteRelation(relation, kind, cardinality, otherKind string) error {
	_, e := m.cache.Exec(mdl_rel, relation, kind, cardinality, otherKind)
	return e
}

func (m *Assembler) WriteRule(name *string, pattern, domain, target string, flags pattern.Flags, prog []byte) error {
	var n sql.NullString
	if name != nil {
		n.Valid = true
		n.String = *name
	}
	_, e := m.cache.Exec(mdl_rule, n, pattern, domain, target, flags.String(), prog)
	return e
}

// WriteStart: store the initial value of a field used at start of play.
func (m *Assembler) WriteStart(owner, field string, value interface{}) error {
	_, e := m.cache.Exec(mdl_start, owner, field, value)
	return e
}

func (m *Assembler) WriteTrait(aspect, trait string, rank int) error {
	_, e := m.cache.Exec(mdl_aspect, aspect, trait, rank)
	return e
}

func (m *Assembler) WriteVerb(relation, verb string) (err error) {
	const asm_verb = `insert into asm_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from asm_verb v
					where v.relation=?1 and v.stem=?2
				)`
	// fix: future. verbs only really matter once we can imply other facts
	// like "supporting" means supporters; worn means set the worn flag; etc.
	// but they are still used for building relations -- so we fake a verb of the same name for now.
	if len(verb) > 0 {
		stem := porterstemmer.StemString(verb)
		_, err = m.cache.Exec(asm_verb, relation, stem)
	} else {
		_, err = m.cache.Exec(asm_verb, relation, relation)
	}
	return
}

var mdl_aspect = tables.Insert("mdl_aspect", "aspect", "trait", "rank")
var mdl_check = tables.Insert("mdl_check", "name", "type", "expect")
var mdl_domain = tables.Insert("mdl_domain", "domain", "path")
var mdl_field = tables.Insert("mdl_field", "kind", "field", "type", "affinity")
var mdl_kind = tables.Insert("mdl_kind", "kind", "path")
var mdl_name = tables.Insert("mdl_name", "noun", "name", "rank")
var mdl_noun = tables.Insert("mdl_noun", "noun", "kind")
var mdl_pair = tables.Insert("mdl_pair", "noun", "relation", "otherNoun", "domain")
var mdl_pat = tables.Insert("mdl_pat", "name", "result", "labels")
var mdl_plural = tables.Insert("mdl_plural", "one", "many")
var mdl_prog = tables.Insert("mdl_prog", "name", "type", "bytes")
var mdl_rel = tables.Insert("mdl_rel", "relation", "kind", "cardinality", "otherKind")
var mdl_rule = tables.Insert("mdl_rule", "name", "pattern", "domain", "target", "phase", "prog")
var mdl_spec = tables.Insert("mdl_spec", "type", "name", "spec")
var mdl_start = tables.Insert("mdl_start", "owner", "field", "value")
