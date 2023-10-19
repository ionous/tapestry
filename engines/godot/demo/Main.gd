extends Control

# name of the tapestry domain to play
@export var scene_name : String

@onready var _tap : TapGame = find_child("TapGame")
@onready var _output : RichTextLabel = find_child("TextOutput")
@onready var _scroll : ScrollContainer = find_child("ScrollContainer")
#@onready var _input : LineEdit = find_child("TextInput")

# When the node enters the scene tree for the first time.
func _ready():
	_tap.restart(scene_name)

# When the player has entered new text commands
func _on_text_input(text: String):
	_output.append_text("> " + text.replace("[", "[lb]") + "\n")
	_tap.fabricate(text)

# scroll to the end whenever the end changes
# ( there's some sort of paint delay after append_text; but this works fine )
var _last_max 
func _process(_delta):
	var vbar = _scroll.get_v_scroll_bar()
	if _last_max != vbar.max_value:
		vbar.value = vbar.max_value
		_last_max = vbar.max_value
