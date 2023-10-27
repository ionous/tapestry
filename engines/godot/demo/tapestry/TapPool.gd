extends Node
# A collection of references to tapestry objects.
# Intended for use as a singleton.

# dictionary key->[obj]
var _all = {}

func get_by_id(id: String) -> TapObject:
	return _all.get(id)

func ensure(id: String) -> TapObject:
	var obj = _all.get(id)
	if !obj and id: 
		obj = TapObject.new(id)
		_all[id] = obj
	return obj

# change the parent-child rel
func reparent(newParentId, childId) -> void: 
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
