extends Control

# name of the tapestry domain to play
@export var story_name : String
@export var message_box : PackedScene

@onready var _tap : TapGame = find_child("TapGame")
@onready var _output : RichTextLabel = find_child("TextOutput")
@onready var _scroll : ScrollContainer = find_child("ScrollContainer")
@onready var _input : LineEdit = find_child("TextInput")
@onready var _view : Node = find_child("GameView")

const player: String = "self"
var _story: TapStory

# the player's room ( or enclosure ) has changed
signal room_changed(id: String)
# the name of the player's room changed
# fix? consolidate with room_changed
signal location_changed(name: String)
# the story title has changed
signal title_changed(title: String)
# sent at the start and end of every player interaction
signal turn_started(started: bool)
# the turn counter has incremented
signal turn_changed(turns: int)
# the game score has changed
signal score_changed(score: int)
# show the player some new story text
signal narration_changed(bb_text: String)
#
signal inventory_changed()

# When the node enters the scene tree for the first time.
func _ready():
	_story= TapStory.new()
	_story._starting_turn = self._starting_turn
	_story._changing_state = self._changing_state
	_story._changing_scenes = self._changing_scenes
	_story._reparenting_objects = self._reparenting_objects
	_story._saying_text = self._saying_text
	_tap.restart(story_name, _story)

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

var _turn_text: String # accumulate text over the turn

# if not started, then finished.
func _starting_turn(started: bool) -> void:
	# only editable *after* the story events have finished processing
	_input.editable = not started
	# at the end of a turn:
	if not started:
		# flush accumulated text to the console:
		var bb = TapWriter.ConvertToBB(_turn_text + "<p>")
		_turn_text = ""
		narration_changed.emit(bb)
		# update the current score and turn
		_tap.query([
			TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
			TapCommands.CurrentTurn, func(turn:int): turn_changed.emit(turn)
		])
	turn_started.emit(started)

# on start of the first scene, query the title, location, score, etc.
func _changing_scenes(_scenes: Array, _started: bool):
	# block this callback from future processing:
	_story._changing_scenes = func(_names, _started): pass
	# then issue our query/ies
	_tap.query([
		TapCommands.StoryTitle, func(title:String): title_changed.emit(title)
	])

func _changing_state(noun: String, _aspect: String, state: String):
	if noun == "story" and state == "playing":
		_tap.query([
			TapCommands.StoryTitle, func(title:String): title_changed.emit(title),
			TapCommands.LocationName, func(named:String): location_changed.emit(named),
			TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
			TapCommands.CurrentTurn, func(turn:int): turn_changed.emit(turn),
			TapCommands.CurrentObjects, func(objects: Dictionary):
				var root = TapPool.rebuild(objects)
				room_changed.emit(root.id)
		])

func _reparenting_objects(pid: String, cid: String): # -> void  or Signal
	# issue a query
	if pid == player:
		inventory_changed.emit()
	elif cid == player:
		_tap.query([
			TapCommands.StoryTitle, func(title:String): title_changed.emit(title),
			TapCommands.LocationName, func(named:String): location_changed.emit(named),
			TapCommands.CurrentObjects, func(objects: Dictionary):
				var root = TapPool.rebuild(objects)
				room_changed.emit(root.id)
		])

func _saying_text(text: String): # -> void  or Signal
	_turn_text += text
	if text != "<p>":

		var title = ""
		if text.begins_with("<p>"):
			text = text.substr(3)
		# fix: look at semantic tags like "title" for message box title.
		if text.begins_with("<b>"):
			var i = text.find("</b>")
			if i > 0:
				title = text.substr(3, i-3)
				text =  text.substr(i+4)
		text = TapWriter.ConvertToBB(text).strip_edges()
		if text.length():
			var box = message_box.instantiate()
			box.title = title
			box.dialog_text = text
			var parent_rect: Rect2i= _view.get_rect();
			var pos = parent_rect.position + (parent_rect.size - box.size) / 2;
			self.add_child(box)
			box.popup(Rect2i( pos, box.size ))
			return box.popup_hide

