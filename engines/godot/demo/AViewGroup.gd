extends Node
class_name AViewGroup
# Switches between scenes based on the current game root.
# Should be parented to a TapGame

# The node where ASceneView children  should attach themselves
@export var view_target: Node

var _last : ASceneView 

func _ready():
	var tap: TapGame = get_parent() as TapGame
	assert(tap, "expected parent to be a TapGame")
	if tap:
		tap.root_changed.connect(func(_pool, new_root: TapObject): 
			update_view(new_root.id))

func update_view(room_id: String):
	var kids = self.find_children(room_id,"ASceneView", false)
	assert(kids.size() == 1)
	if kids.size() > 0:
		var view : ASceneView = kids[0]
		if view != _last:
			if _last:
				_last.visible = false   
			view.visible = true  
			_last = view

