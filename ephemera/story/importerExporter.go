package story

import (
	r "reflect"

	"github.com/ionous/errutil"
)

func (k *Importer) Source() string {
	return k.source
}

func (k *Importer) NewType(t string) (ret interface{}, err error) {
	if rtype, ok := k.cmds[t]; !ok {
		err = errutil.New("unknown type", t)
	} else {
		ret = r.New(rtype).Interface()
		k.path.push(t) // push most recently new'd type
	}
	return
}

func (k *Importer) Finalize(ptr interface{}) (ret interface{}, err error) {
	k.path.pop()
	if stub, ok := ptr.(stubImporter); !ok {
		ret = ptr
	} else {
		ret, err = stub.ImportStub(k)
	}
	return
}

type stubImporter interface {
	ImportStub(*Importer) (ret interface{}, err error)
}

type programPath struct {
	stack         []string
	activityDepth int
}

func (p *programPath) inProg() bool {
	return p.activityDepth > 0

}
func (p *programPath) push(t string) {
	p.stack = append(p.stack, t)
	if t == "activity" {
		p.activityDepth++
	}
}

func (p *programPath) pop() {
	end := len(p.stack) - 1
	if last := p.stack[end]; last == "activity" {
		p.activityDepth--
	}
	p.stack = p.stack[:end] // pop most recently new'd type
}
