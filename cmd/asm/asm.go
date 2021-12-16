// Generates a playable database from a story file.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// ex. go run asm.go -in /Users/ionous/Documents/Iffy/stories/shared -out /Users/ionous/Documents/Iffy/scratch/shared/play.db
func main() {
	var srcPath, outFile string
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	// var printStories bool
	// printStories:= flag.Bool("log", false, "write imported stories to console")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if len(outFile) == 0 {
		dir, _ := filepath.Split(srcPath)
		outFile = filepath.Join(dir, "play.db")
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// lets do this in the dumbest of ways for now.
	var cat eph.Catalog // fix: capture "Dilemmas" and LogWarning?
	var writeErr error  // fix: this seems less than ideal; maybe writer should return err.
	k := story.NewImporter(collectEphemera(&cat, &writeErr), storyMarshaller)
	if e := distill(k, srcPath); e != nil {
		log.Fatalln(e)
	} else if writeErr != nil {
		log.Fatalln(writeErr)
	} else {
		log.Println("assembling....")
		if e := Assemble(&cat, outFile); e != nil {
			log.Fatalln(e)
		}
	}
}

func materialize(key string, arg int) string {
	return fmt.Sprintf(`
-- create a virtual table consisting of the paths parts turned into ids:
with recursive
-- str is a list of comma separated parts, each time dropping the left-most part.
parts(str, ids) AS (
select ?%[2]d || ',',  ''
union all
select substr(str, 1+instr(str, ',')), ids || ( 
	-- turn the left most part into a rowid
	select rowid from mdl_%[1]s where %[1]s is substr(str, 0, instr(str, ','))
) || ','
from parts
-- the last str printed is empty, and it contains the full id path.
where length(str) > 1
-- stop any accidental infinite recursion
limit 10) `, key, arg)
}

// select the id where all of the parts have been consumed, or if there were no parts (the root) select the empty string.
const materialized = `(select ids from parts where length(str) == 0 union all select '' limit 1)`

// selects from mdl_<key> where <key>=?<arg>
func idOf(key string, arg int) string {
	return fmt.Sprintf("(select rowid from mdl_%[1]s where %[1]s = ?%[2]d)", key, arg)
}

func insert(name string, args ...string) string {
	var ins strings.Builder
	ins.WriteString("insert into ")
	ins.WriteString(name)
	ins.WriteRune('(')
	for i, cnt := 0, len(args); i < cnt; i += 2 {
		key := args[i]
		if i > 0 {
			ins.WriteRune(',')
		}
		ins.WriteString(key)
	}
	ins.WriteRune(')')
	ins.WriteString(" values (")
	for i, cnt := 1, len(args); i < cnt; i += 2 {
		arg := args[i]
		if i > 1 {
			ins.WriteRune(',')
		}
		if len(arg) > 0 {
			ins.WriteString(arg)
		} else {
			ins.WriteRune('?')
			ins.WriteString(strconv.Itoa((i + 1) / 2))
		}
	}
	ins.WriteRune(')')
	return ins.String()
}

// select the id of a key ( ex. mdl_kind.kind ) indexed by domain name d and key name n
func simpleScope(key string, d, n int) string {
	return fmt.Sprintf(
		`(select key.rowid 
		from mdl_%[1]s key
		join mdl_domain md
		on (key.domain=md.rowid)
		where (md.domain=?%[2]d) and (key.%[1]s=?%[3]d))`, key, d, n)
}

// same as simple scope, but the domain d can be a child of the key's domain
// so to find the correct key, we have to look through all domains
// where its materialized path contains the key's domain
//   ex. does path ( ",3,2,1," ) contain ( ",3," )
//   or, at the root ( ",1,," ) contain ( ",1," ).
// the prefixed and suffix commas are to avoid partial matches:
//   ex. "21," against "1," or ",12" against ",1"
func derivedScope(key string, d, n int) string {
	return fmt.Sprintf(
		`(select key.rowid
		from mdl_%[1]s key
		join mdl_domain md
		where (key.%[1]s = ?%[3]d)
		and instr(',' || md.rowid || ',' || md.path, ',' || key.domain || ','))`, key, d, n)
}

// if i have something defined in domain 2 it should be visible in 3
// we look through all domains and build paths
// ,1,,
// ,1,2,
// ,1,2,3,
// each time asking if it contains ,2,

const unchanged = ""

// rewrite some tables to use ids
var idswriter = map[string]string{
	// domain name + aspect name select a specific kind entry.
	// fix? should the traits be field ids?
	// you'd do the same as simple scope to get the kind; maybe changing this into a with table
	// FIX: ! the fields of the aspect arent being written
	// they should have a bunch of bool fields
	mdl.Aspect: insert("mdl_aspect",
		"domain", idOf("domain", 1), // redundant, has the same domain as the aspect's kind
		"aspect", simpleScope("kind", 1, 2),
		"trait", unchanged,
		"rank", unchanged,
	),
	// turn domain name into an id
	mdl.Check: insert("mdl_check",
		"domain", idOf("domain", 1),
		"name", unchanged,
		"value", unchanged,
		"affinity", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// turn the materialized path of domain ancestor names into ancestor idOfs
	mdl.Domain: materialize("domain", 2) +
		insert("mdl_domain",
			"domain", unchanged,
			"path", materialized,
			"at", unchanged,
		),
	// domain name + kind name select a specific kind entry.
	// domain is redundant; fields exist per kind, not per domain.
	mdl.Field: insert("mdl_field",
		"domain", idOf("domain", 1), // redundant, has the same domain as the field's kind
		"kind", simpleScope("kind", 1, 2),
		"field", unchanged,
		"affinity", unchanged,
		"type", simpleScope("kind", 1, 5),
		"at", unchanged,
	),
	// turn domain name into an id
	mdl.Grammar: insert("mdl_grammar",
		"domain", idOf("domain", 1),
		"name", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// turn domain name into an id, and materialize the ancestor path
	mdl.Kind: materialize("kind", 3) +
		insert("mdl_kind",
			"domain", idOf("domain", 1),
			"kind", unchanged,
			"path", materialized,
			"at", unchanged,
		),
	// turn domain, kind, field into ids, associated with the local var's initial assignment.
	// domain and kind become redundant b/c fields exist at the scope of the kind.
	mdl.Local: string(`with parts(did, domain, kid, kind, fid, field) as (
		select md.rowid, md.domain, mk.rowid, mk.kind, mf.rowid, mf.field
		from mdl_field mf
		join mdl_kind mk
			on (mk.rowid = mf.kind)
		join mdl_domain md
			on (mk.domain = md.rowid))
		insert into mdl_local(domain, kind, field, assign)
		select did, kid, fid, ?4
		from parts where domain=?1 and kind=?2 and field=?3`,
	),
	mdl.Name: insert("mdl_name",
		"domain", idOf("domain", 1), // currently redundant, names have the same scope as their noun.
		"noun", simpleScope("noun", 1, 2),
		"name", unchanged,
		"rank", unchanged,
		"at", unchanged,
	),
	mdl.Noun: insert("mdl_noun",
		"domain", idOf("domain", 1), // domain where the noun was declared
		"noun", unchanged,
		"kind", derivedScope("kind", 1, 3),
		"at", unchanged,
	),
	mdl.Pair: insert("mdl_pair",
		"domain", idOf("domain", 1), // domain where the pair was declared
		"noun", derivedScope("noun", 1, 2),
		"relation", derivedScope("rel", 1, 3),
		"otherNoun", derivedScope("noun", 1, 4),
		"at", unchanged,
	),
	mdl.Pat: insert("mdl_pat",
		"domain", idOf("domain", 1), // redundant, has the same domain as the pattern's kind
		"kind", simpleScope("kind", 1, 2),
		"labels", unchanged, // fix? this is a materialized field, should it be field ids?
		"result", unchanged, // fix? this is a field, should it be a field id?
	),
	mdl.Plural: insert("mdl_plural",
		"domain", idOf("domain", 1),
		"many", unchanged,
		"one", unchanged,
		"at", unchanged,
	),
	mdl.Rel: insert("mdl_rel",
		"domain", idOf("domain", 1), // redundant, has the same domain as the relation's kind
		"rel", unchanged,
		"kind", simpleScope("kind", 1, 3),
		"cardinality", unchanged,
		"otherKind", simpleScope("kind", 1, 5),
		"at", unchanged,
	),
	mdl.Rule: insert("mdl_rule",
		"domain", idOf("domain", 1), // domain where the rule was declared
		"kind", derivedScope("kind", 1, 2),
		"target", derivedScope("kind", 1, 3),
		"phase", unchanged,
		"filter", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// first: build a virtual [domain, noun, field] table
	// note: values are written per noun, not per domain; so the domain column is redundant once the noun id is known.
	// to get to the field id, we have to look at all possible fields for the noun.
	// given the kind of the noun, accept all fields who's kind is in its materialized path.
	// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
	mdl.Value: string(`with parts(did, domain, nin, noun, fid, field) as (
			select md.rowid, md.domain, mn.rowid, mn.noun, mf.rowid, mf.field
			from mdl_noun mn
			join mdl_domain md
				on (mn.domain = md.rowid)
			join mdl_kind mk
				on (mn.kind = mk.rowid)
			left join mdl_field mf
				where instr(',' || mk.rowid || ',' || mk.path, ',' || mf.kind || ','))
			insert into mdl_value(domain, noun, field, value, affinity, at)
			select did, nin, fid, ?4, ?5, ?6
			from parts where domain=?1 and noun=?2 and field=?3`,
	),
}

// a terrible way to optimize database writes
type qel struct {
	tgt  string
	args []interface{}
}
type qels []qel

func (q *qels) Write(tgt string, args ...interface{}) (err error) {
	(*q) = append(*q, qel{tgt, args})
	return
}

func Assemble(cat *eph.Catalog, outFile string) (err error) {
	var queue qels
	var w eph.Writer = &queue

	// go process all of the ephemera
	if e := cat.AssembleCatalog(eph.PhaseActions{
		eph.AncestryPhase: eph.AncestryPhaseActions,
		eph.NounPhase:     eph.NounPhaseActions,
	}); e != nil {
		err = e
	} else if e := cat.WritePlurals(w); e != nil {
		err = e
	} else if e := cat.WriteDomains(w); e != nil {
		err = e
	} else if e := cat.WriteKinds(w); e != nil {
		err = e
	} else if e := cat.WriteAspects(w); e != nil {
		err = e
	} else if e := cat.WriteFields(w); e != nil {
		err = e
	} else if e := cat.WriteNouns(w); e != nil {
		err = e
	} else if e := cat.WriteNames(w); e != nil {
		err = e
	} else if e := cat.WritePatterns(w); e != nil {
		err = e
	} else if e := cat.WriteLocals(w); e != nil {
		err = e
	} else if e := cat.WriteDirectives(w); e != nil {
		err = e
	} else if e := cat.WriteRelations(w); e != nil {
		err = e
	} else if e := cat.WritePairs(w); e != nil {
		err = e
	} else if e := cat.WriteRules(w); e != nil {
		err = e
	} else if e := cat.WriteValues(w); e != nil {
		err = e
	} else if e := cat.WriteChecks(w); e != nil {
		err = e
	} else {
		log.Println("writing", len(queue), "entries")
		if outFile, e := filepath.Abs(outFile); e != nil {
			err = e
		} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
			err = errutil.New("couldn't clean output file", outFile, e)
		} else {
			// 0755 -> readable by all but only writable by the user
			// 0700 -> read/writable by user
			// 0777 -> ModePerm ... read/writable by all
			os.MkdirAll(path.Dir(outFile), os.ModePerm)
			if db, e := sql.Open(tables.DefaultDriver, outFile); e != nil {
				err = errutil.New("couldn't create output file", outFile, e)
			} else {
				defer db.Close()
				if e := tables.CreateModel(db); e != nil {
					err = errutil.New("couldnt create model", e)
				} else if tx, e := db.Begin(); e != nil {
					err = errutil.New("couldnt create transaction", e)
				} else {
					for _, q := range queue {
						var out string
						if sel, ok := idswriter[q.tgt]; ok {
							out = sel
						} else {
							out = q.tgt
						}
						if _, e := tx.Exec(out, q.args...); e != nil {
							tx.Rollback()
							err = errutil.New("couldnt write to", q.tgt, e)
							break
						}
					}
					if err == nil {
						if e := tx.Commit(); e != nil {
							err = errutil.New("couldnt commit", e)
						}
					}
				}
			}
		}
	}
	return
}

func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(cat *eph.Catalog, out *error) story.WriterFun {
	// fix: needs to be more clever eventually...
	if e := cat.AddEphemera(
		eph.EphAt{
			At:  "asm",
			Eph: &eph.EphBeginDomain{Name: "entire_game"}}); e != nil {
		panic(e)
	}
	// built in kinds -- see ephKinds.go
	// fix? move to an .if file?
	kinds := []string{
		eph.KindsOfAction, eph.KindsOfPattern,
		eph.KindsOfAspect, "",
		eph.KindsOfEvent, eph.KindsOfPattern,
		eph.KindsOfKind, "",
		eph.KindsOfPattern, "",
		eph.KindsOfRecord, "",
		eph.KindsOfRelation, "",
	}
	for i := 0; i < len(kinds); i += 2 {
		k, p := kinds[i], kinds[i+1]
		if e := cat.AddEphemera(
			eph.EphAt{
				At:  "built in kinds",
				Eph: &eph.EphKinds{Kinds: k, From: p}}); e != nil {
			panic(e)
		}
	}
	var i int
	return func(el eph.Ephemera) {
		if e := cat.AddEphemera(eph.EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			*out = errutil.Append(*out, e)
		}
		i++ // temp
	}
}

func distill(k *story.Importer, srcPath string) (err error) {
	if srcPath, e := filepath.Abs(srcPath); e != nil {
		err = e
	} else if e := readPaths(k, srcPath); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	} else {
		k.Flush()
	}
	return
}

// read a comma-separated list of files and directories
func readPaths(k *story.Importer, filePaths string) (err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = errutil.Append(err, e)
		} else {
			which := readOne
			if info.IsDir() {
				which = readMany
			}
			if e := which(k, filePath); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func readMany(k *story.Importer, path string) error {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	return filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() {
			if ext := filepath.Ext(path); ext == CompactExt || ext == DetailedExt {
				if e := readOne(k, path); e != nil {
					err = errutil.New("error reading", path, e)
				}
			}
		}
		return
	})
}

func readOne(k *story.Importer, path string) (err error) {
	log.Println("reading", path)
	if fp, e := os.Open(path); e != nil {
		err = e
	} else {
		defer fp.Close()
		if b, e := io.ReadAll(fp); e != nil {
			err = e
		} else if script, e := decodeStory(path, b); e != nil {
			err = errutil.New("couldn't decode", path, "b/c", e)
		} else if e := k.ImportStory(path, script); e != nil {
			err = errutil.New("couldn't import", path, "b/c", e)
		}
	}
	return
}

func decodeStory(path string, b []byte) (ret *story.Story, err error) {
	var curr story.Story
	switch ext := filepath.Ext(path); ext {
	case CompactExt:
		if e := cin.Decode(&curr, b, iffy.AllSignatures); e != nil {
			err = e
		} else {
			ret = &curr
		}
	case DetailedExt:
		if e := din.Decode(&curr, iffy.Registry(), b); e != nil {
			err = e
		} else {
			ret = &curr
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
