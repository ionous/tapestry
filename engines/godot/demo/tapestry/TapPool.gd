extends Node
# A collection of references to tapestry objects.
# Intended for use as a singleton.

# dictionary key->[obj]
var _all = {}

func get_by_id(id: String) -> TapObject:
	return _all.get(id)

func ensure(id: String) -> TapObject:
	var obj = _all.get(id)
	if !obj: 
		obj = TapObject.new(id)
		_all[id] = obj
	return obj

# rebuild from the passed collection
func rebuild(collection: Variant, parentId: String = "") -> TapObject:
	var obj = ensure(collection.id)
	obj.parentId = parentId
	obj.name = collection.name
	obj.kind = collection.kind
	obj.traits = collection.traits
	obj.childIds = [] # reset
	var kids = collection.get("kids")
	if kids:
		for kid in kids:
			var child = rebuild(kid, obj.id)
			obj.childIds.push_back(child.id)
	return obj
