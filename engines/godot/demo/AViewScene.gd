extends Node
class_name AViewScene
# Expects to have a parent node of AViewGroup with a valid scene_target

@export var scene : PackedScene

var _scene : Node

func _ready():
	if not scene:
		push_warning("'%s' is missing a scene" % self.name)

var visible: bool:
	get:
		return _scene != null 

	set(yes):
		if yes and scene:
			var g = self.get_parent() as AViewGroup
			_scene = scene.instantiate()
			g.view_target.add_child(_scene)
		elif !yes and _scene:
			_scene.queue_free()
			_scene = null
