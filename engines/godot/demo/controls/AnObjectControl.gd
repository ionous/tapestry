extends Control
class_name AnObjectControl
# The Godot Control version of a reference to a tapestry object
# Looks at its direct children for requirements on whether to show itself.

# the name of the requested object
# if empty, and the node name starts with #, will use the node name
# if empty and the node name doesnt stat with #, the control will look upwards for a name.
@export var object_id : String
# if not optional, warns if it cant find the named object
@export var optional : bool 
# when specified, the object must have all of these traits for the control to show.
@export var include_traits : Array[String]
# if the object has any of the specified traits, the control will hide.
@export var exclude_traits : Array[String]
# when true, the nearest parent object_id must be an actual parent
@export var contained : bool
# when specified, the object must have the specified parent for the control to show.
@export var include_parent : String
# if the object has the specified parent, the control will hide.
@export var exclude_parent : String

# runtime requirement helper
var requires : TapRequires

func _init():
	# godot doesnt let you set default values for child classes :/ 
	set_mouse_filter(MOUSE_FILTER_IGNORE)

func _enter_tree():
	requires = TapRequires.new(self)
	if optional and not requires.get_object():
		push_warning("couldn't find the requested object '%s' at '%s'" % [requires.id, self.name])

func _process(_dt):
	# fix: possibly a signal, flag, bitfields for dirty instead.
	visible = requires.check_requirements()
