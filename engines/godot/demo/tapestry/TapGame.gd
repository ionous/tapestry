class_name TapGame
extends Node
# Communication between the game and Tapestry

# hrm.... vue for godot? because this feels messy:
# 1. declare the signal, 2. trigger the signal, 3. add a script, 4. connect the signal, 5. implement the handler;
# i definitely like "set data, connect data" more.
signal title_changed(title: String)
signal score_changed(score: int)
signal turns_changed(turns: int)
signal location_changed(name: String)
signal narration_changed(bb_text: String)
signal root_changed(id: String)

var _root : TapObject               # top most object ( ex. the current room )
var _frames : Array[TapFrame]       # data received from tapestry in _post_to_endpoint; processed in _process
var _story : TapStory = TapStory.NewDefault() # helpers for user interaction during frame processing
var _blocked : bool                  # block new tapestry requests while processing previous ones

func _init():
	set_process(false) # only runs when needed

# are tapestry commands allowed?
func is_blocked() -> bool:
	return _blocked

# the nearby objects have changed rebuild them
func _rebuild(collection: Dictionary) -> void:
	_root = TapPool.rebuild(collection)

# restart
func restart(scene: String, story: TapStory) -> void:
	assert(not is_blocked())
	# FIX: assumes synchronous....
	# probably want to change to "await"
	# so that users can interact with things, watch state changes, etc.
	self._story = story if story else TapStory.NewDefault()
	self._post_to_endpoint("restart", scene)
	self.post([
		TapCommands.StoryTitle, func(title:String): title_changed.emit(title),
		TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
		TapCommands.CurrentTurn, func(turn:int): turns_changed.emit(turn),
		TapCommands.LocationName, func(named:String): location_changed.emit(named),
		TapCommands.CurrentObjects, func(root:Dictionary):
			_rebuild(root),
	])

# player has typed some text
func fabricate(text: String) -> void:
	assert(not is_blocked())
	var player = TapPool.ensure("self")
	var prevLoc = player.parentId
	self.post([
		# send the player input; no particular response except to listen to events
		TapCommands.Fabricate(text), func(_none):
			if player.parentId != prevLoc:
				self.post([
					TapCommands.LocationName, func(named:String): location_changed.emit(named),
					TapCommands.CurrentObjects, func(root:Dictionary): _rebuild(root),
				]),
		# query for the new score and turn each frame
		TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
		TapCommands.CurrentTurn, func(turn:int): turns_changed.emit(turn),
	])


# send cmds and their response handlers
func post(msgCalls: Array) -> void:
	assert(msgCalls.size() & 1  == 0, "expected an equal number of pool and calls")
	var sends = []
	var calls = []
	for i in range(0, msgCalls.size(), 2):
		sends.push_back(msgCalls[i+0])
		calls.push_back(msgCalls[i+1])
	self._post_to_endpoint("query", sends, calls)

# send cmds and queue response frames
func _post_to_endpoint(endpoint: String, blob: Variant, calls: Array= []):
	var res: Array = Tapestry.post(endpoint, JSON.stringify(blob)) as Array
	if res and res.size() > 0:
		for i in res.size():
			var cmd = TapCommands.Parse(res[i])
			var callback = calls[i] if i < calls.size() else null
			var frame = TapFrame.New( cmd.sig, cmd.body, callback )
			_frames.push_back(frame)
		_starting_frame(true)

func _starting_frame(block: bool):
	var was: bool = is_blocked()
	_blocked = block
	var now: bool = is_blocked()
	if was != now:
		_story.starting_frame(now)
		set_process(now) # make sure we wake back up again

# process response frames
func _process(_delta):
	# we're going to await
	# await might not finish before the next process
	# so turn off process until we're done with the current set of frames
	set_process(false)

	var out: String = ""
	while _frames.size():
		var frames : Array[TapFrame] = _frames
		_frames = []
		for frame in frames:
			if frame.error:
				push_error(frame.error)

			for evt in frame.events:
				var cmd = TapCommands.Parse(evt)
				var text: String = await handle_event(cmd)
				out += text
			frame.report_result()

	# we dont have good dirty checking,
	# so broadcast this at the end of every frame
	if _root:
		root_changed.emit(_root.id)
	_starting_frame(false)

	# hrmm... the text writer doesnt properly space trailing <p>
	# it probably needs to be stateful ( and instance )
	# and then maybe all this output return text could be cleaned up
	var bb = TapWriter.ConvertToBB(out + "<p>")
	narration_changed.emit(bb)


# each evt is a Tapestry Event
# as of 2023-10-18, the complete set is:
# FrameOutput, PairChanged, SceneEnded, SceneStarted, StateChanged
func handle_event(cmd: TapCommands.Cmd) -> String:
	var out  = ""
	match cmd.sig:
		# printed text; accumulates over multiple events
		"FrameOutput:":
			var text: String = cmd.body as String
			await _story.saying_text(text)
			out += text

		# fix: we need the prev state in order to be able to clear it
		"StateChanged noun:aspect:trait:":
			var noun = cmd.body[0]
			var aspect = cmd.body[1]
			var traitn = cmd.body[2] # doesn't like "trait"???
			await _story.changing_state(noun, aspect, trait)

		# relational change
		#  fix: we dont get both sides of the relation change:
		#  we only get new relations; fine for now.
		"PairChanged a:b:rel:":
			var rel : String = cmd.body[2]
			if rel == "whereabouts":
				var childId : String = cmd.body[1]     # b
				var newParentId : String = cmd.body[0] # a
				# remove from old parentId:
				TapPool.reparent(newParentId, childId)
				await _story.reparenting_objects(newParentId, childId)

	return out

