package story

// func (k *Importer) Source() string {
// 	return k.source
// }

// --- din Decoder calls din.Registry.NewType ( cin does similar somewhere )
// but it doesnt call this.
//
// func (k *Importer) NewType(s, t string) (ret interface{}, err error) {
// 	if p, e := k.cmds.NewType(t); e != nil {
// 		err = e
// 	} else {
// 		ret = p
// 		k.path.push(s) // slot ( aka interface ) of most recently new'd type
// 	}
// 	return
// }

// --------------- NOT BEING CALLED ANYWHERE
//
// func (k *Importer) Finalize(ptr interface{}) (ret interface{}, err error) {
// 	if stub, ok := ptr.(stubImporter); !ok {
// 		ret = ptr
// 	} else {
// 		ret, err = stub.ImportStub(k)
// 	}
// 	k.path.pop()
// 	return
// }

// type stubImporter interface {
// 	ImportStub(*Importer) (ret interface{}, err error)
// }

type programPath struct {
	stack         []string
	activityDepth int
}

func (p *programPath) inProg() bool {
	return p.activityDepth > 0
}

// func (p *programPath) push(t string) {
// 	p.stack = append(p.stack, t)
// 	if t == "execute" {
// 		p.activityDepth++
// 	}
// }

// func (p *programPath) pop() {
// 	end := len(p.stack) - 1
// 	if last := p.stack[end]; last == "execute" {
// 		p.activityDepth--
// 	}
// 	p.stack = p.stack[:end] // pop most recently new'd type
// }
