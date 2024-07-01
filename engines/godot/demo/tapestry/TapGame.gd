class_name TapGame
extends Node
# Communication between the game and Tapestry

var _frames: Array[TapFrame]       # data received from tapestry in _post_to_endpoint; processed in _process
var _story: TapStory = TapStory.NewDefault() # helpers for user interaction during frame processing
var _turn_running: bool            # block new tapestry requests while processing previous ones

func _init():
	set_process(false) # only runs when needed

# are tapestry commands allowed?
func is_running_turn() -> bool:
	return _turn_running

# restart
func restart(scene: String, story: TapStory = null) -> void:
	# FIX: assumes synchronous....
	# probably want to change to "await"
	# so that users can interact with things, watch state changes, etc.
	self._story = story if story else TapStory.NewDefault()
	_run_turn("restart", scene)

# player has typed some text
func fabricate(text: String) -> void:
	_run_turn("query", [TapCommands.Fabricate(text)])

# get some information from the tapestry runtime.
func query(msgCalls: Array):
	assert(msgCalls.size() & 1  == 0, "expected an equal number of pool and calls")
	var sends: Array[Variant] = []
	var calls: Array[Callable] = []
	for i in range(0, msgCalls.size(), 2):
		sends.push_back(msgCalls[i+0])
		calls.push_back(msgCalls[i+1])

	var frames: Array[TapFrame] = _post("query", sends)
	for i in frames.size():
		var frame = frames[i]
		assert(not frame.events.size(), "not expecting events for queries")
		if frame.error:
			push_warning(frame.error)
		else:
			var cb:Callable = calls[i]
			if cb:
				cb.call(frame.result)

# execute a tapestry command that expects some change in game state.
func _run_turn(endpoint: String, blob: Variant):
	if _turn_running:
		push_warning("cant start a new turn until after the current turn completes")
	else:
		var frames = _post(endpoint, blob)
		_frames.append_array(frames)
		_set_turn_running(true)

func _set_turn_running(go: bool):
	set_process(go) # make sure we wake up to process those frames
	_story.starting_turn(go) # callback notification

# post a request to the tapestry lib
# returns an array of TapFrame responses
func _post(endpoint: String, blob: Variant) -> Array[TapFrame]:
	var out: Array[TapFrame] = []
	var raw = Tapestry.post(endpoint, JSON.stringify(blob))
	var res: Array = raw as Array if raw else []
	for i in res.size():
		var cmd = TapCommands.Parse(res[i])
		out.push_back(TapFrame.New( cmd.sig, cmd.body ))
	return out

# process response frames
func _process(_delta):
	# we're going to await
	# await might not finish before the next process
	# so turn off process until we're done with the current set of frames
	set_process(false)
	var frames : Array[TapFrame] = _frames
	_frames = []
	for frame in frames:
		if frame.error:
			push_warning(frame.error)
		for evt in frame.events:
			await _handle_event(TapCommands.Parse(evt))

	# done with frame.
	_set_turn_running(false)

# each event is a taoestry command from frame.idl
func _handle_event(cmd: TapCommands.Cmd):
	match cmd.sig:
		# printed text; accumulates over multiple events
		"FrameOutput:":
			var text: String = cmd.body as String
			await _story.saying_text(text)

		"SceneStarted:":
			await _story.changing_scenes(cmd.body as Array, true)

		"SceneEnded:":
			await _story.changing_scenes(cmd.body as Array, false)

		# fix: we need the prev state in order to be able to clear it
		"StateChanged noun:aspect:prev:trait:":
			var noun = cmd.body[0]
			var aspect = cmd.body[1]
			var prev = cmd.body[2]
			var state = cmd.body[3] # godot doesn't like you to use the word trait
			TapPool.change_state(noun, prev, state)
			await _story.changing_state(noun, aspect, prev, state)

		# relational change
		#  fix: we dont get both sides of the relation change:
		#  we only get new relations; fine for now.
		"PairChanged a:b:rel:":
			var rel : String = cmd.body[2]
			if rel == "whereabouts":
				var childId : String = cmd.body[1]     # b
				var newParentId : String = cmd.body[0] # a
				# remove from old parentId:
				TapPool.change_hierarchy(newParentId, childId)
				await _story.reparenting_objects(newParentId, childId)

		_:
			print("unhandled %s" % [cmd.sig])
