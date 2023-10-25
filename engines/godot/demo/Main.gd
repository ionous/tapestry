extends Control

# name of the tapestry domain to play
@export var scene_name : String

@onready var _tap : TapGame = find_child("TapGame")
@onready var _output : RichTextLabel = find_child("TextOutput")
@onready var _scroll : ScrollContainer = find_child("ScrollContainer")
@onready var _input : LineEdit = find_child("TextInput")

var _gui : TapGui = TapGui.new()

func _init():
	_gui._display_text = self._display_text
	_gui._block_input = self._block_input

# When the node enters the scene tree for the first time.
func _ready():
	_tap.restart(scene_name, _gui)

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

func _block_input(blocked: bool) -> void: 
	_input.editable = not blocked

func _display_text(text: String):
	print("got some text %s" % text)
