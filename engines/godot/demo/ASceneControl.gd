extends Control
class_name ASceneControl
# create or destroy a scene (containing controls) based on whether this control is visible or not.

@export var scene : PackedScene

var _scene : Node

func _ready():
	if not scene:
		push_warning("'%s' is missing a scene" % self.name)
	else:
		_visibility_changed()
		visibility_changed.connect(_visibility_changed)

func _visibility_changed():
	if visible and not _scene:
		_scene = scene.instantiate()
		self.add_child(_scene)
	elif scene and not visible:
		_scene.queue_free()
		_scene = null
