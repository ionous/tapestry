package story

import (
	"database/sql"
	r "reflect"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// Importer helps read story specific json.
type Importer struct {
	*ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	decoder     *decode.Decoder
	autoCounter ident.Counters
	entireGame  ephemera.Named
	StoryEnv
}

// low level
func NewImporterDecoder(db *sql.DB, dec *decode.Decoder) *Importer {
	return &Importer{
		Recorder:    ephemera.NewRecorder(db),
		oneTime:     make(map[string]bool),
		decoder:     dec,
		autoCounter: make(ident.Counters),
	}
}

func NewImporter(db *sql.DB, reporter decode.IssueReport) *Importer {
	iffy.RegisterGobs()
	dec := decode.NewDecoderReporter(reporter)
	k := NewImporterDecoder(db, dec)
	//
	for _, slats := range iffy.AllSlats {
		dec.AddDefaultCallbacks(slats)
	}
	dec.AddDefaultCallbacks(core.Slats)
	// add ops from iffy.js, including golang generated stubs via stubs.js
	// anything that implements ImportStub() will get processed during ReadSpec.
	k.AddModel(Model)
	//
	return k
}

func (k *Importer) ImportStory(src string, m reader.Map) (ret *Story, err error) {
	k.SetSource(src)
	if i, e := k.decoder.ReadSpec(m); e != nil {
		err = e
	} else if story, ok := i.(*Story); !ok {
		err = errutil.Fmt("imported spec wasn't a story %T", i)
	} else if e := story.ImportStory(k); e != nil {
		err = e
	} else {
		ret = story
	}
	return
}

func (k *Importer) SetSource(s string) *Importer {
	k.Recorder.SetSource(s)
	k.decoder.SetSource(s)
	return k
}

//
func (k *Importer) AddModel(model []composer.Composer) {
	type stubImporter interface {
		ImportStub(k *Importer) (ret interface{}, err error)
	}
	dec := k.decoder
	for _, cmd := range model {
		if _, ok := cmd.(stubImporter); !ok {
			dec.AddCallback(cmd, nil)
		} else {
			// need to pin the loop variable for the callback
			// so pin the type. why not.
			rtype := r.TypeOf(cmd).Elem()
			dec.AddCallback(cmd, func(m reader.Map) (ret interface{}, err error) {
				// create an instance of the stub
				op, at := r.New(rtype), reader.At(m)
				// read it in
				dec.ReadFields(at, op.Elem(), m.MapOf(reader.ItemValue))
				// convert it
				stub := op.Interface().(stubImporter)
				return stub.ImportStub(k)
			})
		}
	}
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
