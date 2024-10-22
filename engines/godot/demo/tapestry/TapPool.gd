extends Node
# A collection of references to tapestry objects.
# Intended for use as a singleton.

# dictionary key->[obj]
var _all = {}

const player: String = "self"

func get_by_id(id: String) -> TapObject:
	return _all.get(id)

func ensure(id: String) -> TapObject:
	var obj = _all.get(id)
	if !obj and id: 
		obj = TapObject.new(id)
		_all[id] = obj
	return obj


# change the parent-child rel
# "reparent" is a function in Node.
func change_state(objId: String, prev: String, next: String) -> void:
	var obj = get_by_id(objId)
	if obj:
		obj.traits.erase(prev)
		obj.traits.push_back(next)

# change the parent-child rel
# "reparent" is a function in Node.
func change_hierarchy(newParentId: String, childId: String) -> void:
	var child = get_by_id(childId)
	if child:
		var oldParent = get_by_id(child.parentId)
		var newParent = ensure(newParentId)

		if oldParent:
			oldParent.childIds.erase(childId)
		child.parentId = newParentId
		if newParentId:
			newParent.childIds.push_back(child.id)

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
