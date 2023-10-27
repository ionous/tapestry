extends RefCounted
class_name TapStory
# connects story events to the game.
# used by the somewhat generic TapGame so the more game specific Main.
# handlers can return promises to suspend processing of subsequent events;
# and handlers can issue new queries ( generating new events. )
# those queries are processed immediately:
# before any still pending events
#
# in alice, the handler was the game's state chart:
# pushing the events through the state tree )  user in game specific ways.


# called while the game is processing previous input
var _starting_turn: Callable

# state of a story has changed
var _changing_scenes: Callable

# state of an object has changed
var _changing_state: Callable

# changing hierarchy
var _reparenting_objects: Callable

# ex. pop up a message box, and return a signal when closed
var _saying_text: Callable

func starting_turn(starting: bool) -> void:
	return self._starting_turn.call(starting)

func changing_scenes(names: Array, started: bool): #-> Signal or nothing
	return self._changing_scenes.call(names, started)

func changing_state(noun: String, aspect: String, state: String): #-> Signal or nothing
	return self._changing_state.call(noun, aspect, state)

func reparenting_objects(pid: String, cid: String): #-> Signal or nothing
	return self._reparenting_objects.call(pid, cid)

func saying_text(text: String): #-> Signal or nothing
	return self._saying_text.call(text)

static func NewDefault():
	var story = TapStory.new()
	story._starting_turn = func(starting: bool):
		print("starting turn: %s" % starting)
	story._changing_scenes = func(names: Array, started: bool):
		print("changing scenes: %s %s" % [names, started])
	story._changing_state = func(noun: String, aspect: String, state: String):
		print("changing states: %s<-%s" % [noun, aspect, state])
	story._reparenting_objects = func(pid: String, cid:String): 
		print("reparenting objects: %s<-%s" % [pid, cid])
	story._saying_text = func(text: String):
		print("saying text: %s" % text)
	return story
