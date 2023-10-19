class_name TapObject
extends RefCounted

# Reference to an external Tapestry object
# ( inner references to parents and children are by id to avoid ref counting loops )
# [ alt: potentially we could keep an "off screen" scene of object nodes ]

var id: String
var name: String
var parent: String
var kind: String
var traits: Array
var kids: Array # id strings

func _init(_id):
	self.id = _id

