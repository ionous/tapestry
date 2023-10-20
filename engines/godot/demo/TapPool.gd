class_name TapPool
extends Node
# A collection of references to tapestry objects and queries.
# Intended to be used as singleton.

# dictionary key->[obj]
var _all = {}

func get_by_id(id: String) -> TapObject:
	return _all[id]

func ensure(id: String) -> TapObject:
	var obj = _all.get(id)
	if !obj: 
		obj = TapObject.new(id)
		_all[id] = obj
	return obj

# rebuild from the passed collection
func rebuild(collection: Variant, parent: String = "") -> TapObject:
	var obj = ensure(collection.id)
	obj.parent = parent
	obj.name = collection.name
	obj.kind = collection.kind
	obj.traits = collection.traits
	obj.kids = [] # reset
	var kids = collection.get("kids")
	if kids:
		for kid in kids:
			var child = rebuild(kid, obj.id)
			obj.kids.push_back(child.id)
	return obj
