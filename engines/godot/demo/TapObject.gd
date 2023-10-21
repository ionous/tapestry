class_name TapObject
extends RefCounted
# Reference to an external Tapestry object
# ( inner references to parents and children are by id to avoid ref counting loops )
# [ alt: potentially we could keep an "off screen" scene of object nodes ]

signal state_changed(aspect: String, prev: String, next: String)
signal relation_changed(rel: String, other: String)

var id: String
var name: String
var kind: String
var traits: Array
var parentId: String
var childIds: Array # id strings

func _init(_id):
	self.id = _id

# returns true if the passed id is the same as this object, or any ancestor
func has_ancestor(ancestor_id: String) -> bool:
	var search: TapObject = self
	while search:
		if search.id == ancestor_id:
			return true
		search = TapPool.get_by_id(search.parentId)
	return false
