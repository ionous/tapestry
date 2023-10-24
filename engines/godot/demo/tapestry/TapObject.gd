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

# returns true if the passed id is the same as this object, or any parent
func has_parent(parent_id: String) -> bool:
	var search: TapObject = self
	while search:
		if search.id == parent_id:
			return true
		search = TapPool.get_by_id(search.parentId)
	return false

# false if one of the includes is missing
func includes(includes_traits: Array):
	for include in includes_traits:
		if not include in traits:
			return false
	return true

# false if one of the excludes is present
func excludes(excludes_traits: Array):
	for exclude in excludes_traits:
		if exclude in traits:
			return false
	return true
