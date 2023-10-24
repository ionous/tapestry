extends Control
class_name HitBox

var obj: TapObject
@export var highlight: ColorRect
var _default_color :Color

@export var highlight_color: Color

func _ready():
	set_mouse_filter(MOUSE_FILTER_STOP)
	obj = TapRequires.find_obj(self)
	assert(obj, "couldn't find tapestry object")
	if highlight:
		_default_color = highlight.color

# note: no method for mouse, but signals and notifications exist	
func _notification(what):
	match what:
		NOTIFICATION_MOUSE_ENTER:
			highlight.color = highlight_color
			print("entered")
		NOTIFICATION_MOUSE_EXIT:
			highlight.color = _default_color
			print("exited")
			

func _gui_input(event):
	#if event is InputEventMouseButton and !event.is_pressed():
	#	obj.notify("clicked")
#	elif event is InputEventMouseMotion:
		#if Rect2(Vector2(), size).has_point(get_local_mouse_position()):
			#print("here")

# Rect2(Vector2(), size).has_point(get_local_mouse_position()):
