package story

import (
	"database/sql"
	r "reflect"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"

	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// Importer helps read story specific json.
type Importer struct {
	*ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	autoCounter ident.Counters
	entireGame  ephemera.Named
	StoryEnv
	// jsonExp.importerExporter
	source string
	cmds   map[string]r.Type
	path   programPath
}

// low level
func NewImporter(db *sql.DB) *Importer {
	iffy.RegisterGobs()
	k := &Importer{
		Recorder:    ephemera.NewRecorder(db),
		oneTime:     make(map[string]bool),
		autoCounter: make(ident.Counters),
	}
	for _, slats := range iffy.AllSlats {
		k.RegisterTypes(slats)
	}
	k.RegisterTypes(Slats) // add story slats
	return k
}

func (i *Importer) RegisterTypes(cmds []composer.Composer) {
	if i.cmds == nil {
		i.cmds = make(map[string]r.Type)
	}
	for _, cmd := range cmds {
		if spec := cmd.Compose(); len(spec.Name) == 0 {
			panic(errutil.Sprintf("Missing type name %T", cmd))
		} else if was, exists := i.cmds[spec.Name]; exists {
			panic(errutil.Sprintf("Duplicate type name %q now: %T, was: %s", spec.Name, cmd, was.String()))
		} else {
			i.cmds[spec.Name] = r.TypeOf(cmd).Elem()
		}
	}
}

func (k *Importer) ImportStory(src string, b []byte) (ret *Story, err error) {
	k.source = src
	k.Recorder.SetSource(src)
	//
	story := new(Story)
	if e := story.UnmarshalDetailed(k, b); e != nil {
		err = e
	} else if e := story.ImportStory(k); e != nil {
		err = e
	} else {
		ret = story
	}
	return
}

//
func (k *Importer) NewName(name, category, ofs string) ephemera.Named {
	return k.NewDomainName(k.currentDomain(), name, category, ofs)
}

func (k *Importer) gameDomain() ephemera.Named {
	if !k.entireGame.IsValid() {
		k.entireGame = k.Recorder.NewName("entire_game", tables.NAMED_SCENE, "internal")
	}
	return k.entireGame
}

func (k *Importer) currentDomain() ephemera.Named {
	domain := k.Current.Domain
	if !domain.IsValid() {
		domain = k.gameDomain()
	}
	return domain
}

// return true if m is the first time once has been called with the specified string.
func (k *Importer) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// NewImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) NewImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		domain := k.gameDomain()
		kKind := k.NewDomainName(domain, kind, tables.NAMED_KINDS, src)
		kAspect := k.NewDomainName(domain, aspect, tables.NAMED_ASPECT, src)
		k.NewAspect(kAspect)
		k.NewField(kKind, kAspect, tables.PRIM_ASPECT, "")
		for i, trait := range traits {
			kTrait := k.NewDomainName(domain, trait, tables.NAMED_TRAIT, src)
			k.NewTrait(kTrait, kAspect, i)
		}
	}
}
