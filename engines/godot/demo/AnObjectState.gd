extends Node
class_name AnObjectState

@export requires : String

var obj: AnObject


func _ready():
	obj = AnObject.Find(self)
	assert( obj, "couldnt find parent")
	obj.tap.state_changed.connect(_state_changed)

func _state_changed(aspect: String, prev: String, next: String):


