extends Node
class_name AViewGroup
# Switches between scenes based on the current game root.
# Should be parented to a TapGame

# The node where AViewScene children  should attach themselves
@export var view_target: Node

var _last : AViewScene

func _on_main_room_changed(id: String):
	var kids = self.find_children(id,"AViewScene", false)
	assert(kids.size() == 1)
	if kids.size() > 0:
		var view : AViewScene = kids[0]
		if view != _last:
			if _last:
				_last.visible = false   
			view.visible = true  
			_last = view
