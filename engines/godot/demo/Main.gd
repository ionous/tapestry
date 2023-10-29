extends Control

# name of the tapestry domain to play
@export var story_name : String
@export var message_box : PackedScene

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
# items have entered/left the player's inventory
signal inventory_changed()
# ui mode change
signal changed_mode(combine:bool)


@onready var _tap_game : TapGame = find_child("TapGame")
@onready var _output : RichTextLabel = find_child("TextOutput")
@onready var _scroll : ScrollContainer = find_child("ScrollContainer")
@onready var _input : LineEdit = find_child("TextInput")
@onready var _game_view : Node = find_child("GameView")

var _story: TapStory
# could be an enum, along with "running turn" maybe
var _combine_target: String 


# When the node enters the scene tree for the first time.
func _ready():
	_story= TapStory.new()
	_story._starting_turn = self._starting_turn
	_story._changing_state = self._changing_state
	_story._changing_scenes = self._changing_scenes
	_story._reparenting_objects = self._reparenting_objects
	_story._saying_text = self._saying_text
	_tap_game.restart(story_name, _story)

func _set_combine_mode(itemId: String):
	if itemId != _combine_target:
		var combining = itemId != ""
		var which = Input.CURSOR_CROSS if combining else Input.CURSOR_ARROW
		Input.set_custom_mouse_cursor(null, which)
		_combine_target = itemId
		changed_mode.emit(combining)

# When the player has entered new text commands
func _on_text_input(text: String):
	_output.append_text("> " + text.replace("[", "[lb]") + "\n")
	_tap_game.fabricate(text)

# from GameButton when clicking on an object in the world.
func _on_clicked_object(obj: String):
	_game_view.show_popup(_build_actions(obj, _combine_target))

# from ButtonBar when clicking on an inventory item
func _on_clicked_item(itemId: String):
	if not _tap_game.is_running_turn():
		if _combine_target:
			# alt: run combine on an inventory item.
			_set_combine_mode("")
		else:
			# actions for each item, plus a custom "combine" mode
			var els:Array = _build_actions(itemId)
			#var obj = TapPool.get_by_id(itemId)
			#if obj and obj.excludes(["worn"]):
			# FIX: need to update the object traits
			var combine: String = "%s %s" % [Icons.Combine, "Use"]
			els.push_back([ combine, func(): _set_combine_mode(itemId) ])
			_game_view.show_popup(els)

# from ButtonBar when clicking on a player action
func _on_clicked_action(act: ActionService.Action):
	# block user input while running a turn
	print("clicked %s" % [act.name])
	if not _tap_game.is_running_turn():
		_tap_game.fabricate(act.format())

# helper to build actions for items
# tbd: move to action service....?
func _build_actions(objId: String, otherId: String = ""):
	var a: Array[ActionService.Action] = \
		ActionService.get_object_actions(objId) if otherId == "" else  \
		ActionService.get_multi_actions(objId, otherId)

	return a.map(func(act: ActionService.Action):
		# each icon has a text label and an "icon"
		var icon = Icons.find_icon(act.name)
		# fix: are there tooltips in godot?
		var label = "%s %s" % [icon.icon, icon.label]
		return [label, func(): _on_action(act, objId, otherId)]
	)

func _on_action(act:ActionService.Action, objId: String, otherId:String):
	var input = act.format(objId, otherId)
	_tap_game.fabricate(input)
	_set_combine_mode("") 

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
		_tap_game.query([
			TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
			TapCommands.CurrentTurn, func(turn:int): turn_changed.emit(turn)
		])
	turn_started.emit(started)

# on start of the first scene, query the title, location, score, etc.
func _changing_scenes(_scenes: Array, _started: bool):
	# block this callback from future processing:
	_story._changing_scenes = func(_names, _started): pass
	# then issue our query/ies
	_tap_game.query([
		TapCommands.StoryTitle, func(title:String): title_changed.emit(title)
	])

func _changing_state(noun: String, _aspect: String, state: String):
	if noun == "story" and state == "playing":
		_tap_game.query([
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
	if pid == TapPool.player:
		inventory_changed.emit()
	elif cid == TapPool.player:
		_tap_game.query([
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
			var parent_rect: Rect2i= _game_view.get_rect();
			var pos = parent_rect.position + (parent_rect.size - box.size) / 2;
			self.add_child(box)
			box.popup(Rect2i( pos, box.size ))
			return box.popup_hide
