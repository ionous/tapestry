package books

import (
	capnp "zombiezen.com/go/capnproto2"
	schemas "zombiezen.com/go/capnproto2/schemas"
	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

// Map is a lazy index of a registry.
// The zero value is an index of the default registry.
type Map struct {
	reg   *schemas.Registry
	nodes map[uint64]schema.Node
}

func (m *Map) registry() (ret *schemas.Registry) {
	if m.reg != nil {
		ret = m.reg
	} else {
		ret = &schemas.DefaultRegistry
	}
	return
}

func (m *Map) UseRegistry(reg *schemas.Registry) {
	m.reg = reg
	m.nodes = make(map[uint64]schema.Node)
}

// Find returns the node for the given ID.
func (m *Map) Find(id uint64) (ret schema.Node, err error) {
	if n := m.nodes[id]; n.IsValid() {
		ret = n
	} else if data, e := m.registry().Find(id); e != nil {
		err = e
	} else if msg, e := capnp.Unmarshal(data); e != nil {
		err = e
	} else if req, e := schema.ReadRootCodeGeneratorRequest(msg); e != nil {
		err = e
	} else if nodes, e := req.Nodes(); e != nil {
		err = e
	} else {
		if m.nodes == nil {
			m.nodes = make(map[uint64]schema.Node)
		}
		// register all the nodes from the schema
		for i := 0; i < nodes.Len(); i++ {
			n := nodes.At(i)
			m.nodes[n.Id()] = n
		}
		ret = m.nodes[id]
	}
	return
}
