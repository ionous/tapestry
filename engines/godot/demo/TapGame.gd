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

# Create a pool of all known tapestry objects
var pool = TapPool.new()

# Top most object ( ex. the current room )
var _root : TapObject 

func _process(_delta):
	if _root: # first frame it will be null
		root_changed.emit(_root.id)
		set_process(false)

func request_rebuild_signal(yes: bool = true):
	set_process(yes)

# the nearby objects have changed rebuild them
func _rebuildpool(collection: Dictionary) -> void:
	_root = pool.rebuild(collection)
	request_rebuild_signal()

# restart
func restart(scene: String) -> void:
	# FIX: assumes synchronous....
	# probably want to change to "await"
	# so that users can interact with things, watch state changes, etc.
	self._post("restart", scene)
	self._query([
		TapCommands.StoryTitle, func(title:String): title_changed.emit(title),
		TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
		TapCommands.CurrentTurn, func(turn:int): turns_changed.emit(turn),
		TapCommands.LocationName, func(named:String): location_changed.emit(named),
		TapCommands.CurrentObjects, func(root:Dictionary): _rebuildpool(root),
	])

# player has typed some text
func fabricate(text: String) -> void:
	var player = pool.ensure("self")
	var prevLoc = player.parent
	self._query([
		# send the player input; no particular response except to listen to events
		TapCommands.Fabricate(text), null,
		# query for the new score and turn each frame
		TapCommands.CurrentScore, func(score:int): score_changed.emit(score),
		TapCommands.CurrentTurn, func(turn:int): turns_changed.emit(turn),
	])
	# todo: consider using events instead
	if player.parent != prevLoc:
		self._query([
			TapCommands.LocationName, func(named:String): location_changed.emit(named),
			TapCommands.CurrentObjects, func(root:Dictionary): _rebuildpool(root),
		])

# given a valid tapestry command:
# return its signature and args in an array of two elements
func _parse_cmd(op: Variant) -> Array:
	var pair: Array
	for k in op:
		if k != "--":
			pair = [k, op[k]]
			break
	return pair

# send cmds and their response handlers
func _query(msgCalls: Array) -> void:
	assert(msgCalls.size() & 1  == 0, "expected an equal number of pool and calls")
	var sends = []
	var calls = []
	for i in range(0, msgCalls.size(), 2):
		sends.push_back(msgCalls[i+0])
		calls.push_back(msgCalls[i+1])
	self._post("query", sends, calls)

func _post(endpoint: String, blob: Variant, calls: Array=[]):
	var res = Tapestry.post(endpoint, JSON.stringify(blob))
	if res:
		var out = _handle_response(res, calls)
		if out:
			var bb = TapWriter.WriteText(out + "<p>")
			narration_changed.emit(bb)

# msgs expects an array of pairs:
# a tapestry command, and a handler to manage that command
func _handle_response(msgs: Array, calls: Array) -> String:
	var out: String = ""
	for i in msgs.size():
		var cmd = _parse_cmd(msgs[i])
		var callback  = calls[i] if i < calls.size() else null
		var sig = cmd[0]
		var args = cmd[1]
		match sig:
			# TODO: look for "SceneStarted" containing the scene we want
			"Frame result:events:error:":
				# var res = args[0]; #-> results from a query
				var events = args[1]
				var err = args[2]
				if events:
					for evt in events:
						out += _process_event(evt)
				print("error", err) # trace out the error

			"Frame result:events:":
				var result = args[0]  # result from a query
				var events = args[1]
				if events:
					for evt in events:
						out += _process_event(evt)

				if callback:
					# ick: we debug.Stringify the results to support "any value"
					# so we have to unpack that too.
					var res = JSON.parse_string(result) if result else null
					callback.call( res )

			_:
				push_error("unhandled message", sig)
	return out


# each evt is a Tapestry Event
# as of 2023-10-18, the complete set is:
# FrameOutput, PairChanged, SceneEnded, SceneStarted, StateChanged
func _process_event(evt: Variant) -> String:
	var out  = ""
	var cmd = _parse_cmd(evt)
	var sig = cmd[0]
	var args = cmd[1]
	match sig:
		# printed text; accumulates over multiple events
		"FrameOutput:":
			out += args

		# fix: we need the prev state in order to be able to clear it
		"StateChanged noun:aspect:trait:":
			var noun = args[0]
			var aspect = args[1]
			var traitn = args[2] # doesn't like "trait"???
			print("state changed: '", noun, "' '", aspect, "' '", traitn,"'")

		# relational change
		#  fix: we dont get both sides of the relation change:
		#  we only get new relations; fine for now.
		"PairChanged a:b:rel:":
			var rel : String = args[2]
			if rel == "whereabouts":
				var childId : String = args[1]     # b
				var newParentId : String = args[0] # a
				# remove from old parent:
				var child = pool.get_by_id(childId)
				if child:
					var oldParent = pool.get_by_id(child.parent)
					if oldParent:
						oldParent.kids.erase(childId)
					child.parent = newParentId
					if newParentId:
						var newParent = pool.ensure(newParentId)
						newParent.kids.push_back(child.id)
					request_rebuild_signal()

		_:
			print("unhandled event", sig)
	return out

