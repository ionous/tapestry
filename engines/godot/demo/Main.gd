extends Control

# name of the tapestry domain to play
@export var story_name : String
@export var message_box : PackedScene

@onready var _tap : TapGame = find_child("TapGame")
@onready var _output : RichTextLabel = find_child("TextOutput")
@onready var _scroll : ScrollContainer = find_child("ScrollContainer")
@onready var _input : LineEdit = find_child("TextInput")
@onready var _view : Node = find_child("GameView")

var _story: TapStory = TapStory.new()
const player: String  = "self"

func _init():
	_story._saying_text = self._saying_text
	_story._starting_frame = self._starting_frame
	_story._reparenting_objects = self._reparenting_objects

# When the node enters the scene tree for the first time.
func _ready():
	_tap.restart(story_name, _story)
	pass

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

func _starting_frame(started: bool) -> void: 
	# only editable *after* the story events have finished processing
	_input.editable = not starting

func _reparenting_objects(pid: String, cid: String) -> Signal:
	if cid == player:
		pass
		# issue a query

func _saying_text(text: String) -> Signal:
	if text != "<p>":
		# fix: look at semantic tags like "title" for message box title.
		var box = message_box.instantiate()
		box.dialog_text = TapWriter.ConvertToBB(text)
		var parent_rect: Rect2i= _view.get_rect();
		var pos = parent_rect.position + (parent_rect.size - box.size) / 2;
		self.add_child(box)
		box.popup(Rect2i( pos, box.size ))
		return box.popup_hide


