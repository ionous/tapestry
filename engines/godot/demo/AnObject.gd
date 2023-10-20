extends Node
class_name AnObject
# The name of the node should be the id the object

var tap: TapObject

func _init():
	var pool  = find_parent("TapGame").pool
	tap = pool.ensure(self.name)

static Find(from: Node) -> TapObject:
	var obj: AnObject
	var p: Node = get_parent()
	while p and not obj:
		obj = p as AnObject
		p = p.get_parent()
	return obj


