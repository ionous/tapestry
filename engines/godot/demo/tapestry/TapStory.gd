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
var _starting_frame: Callable

# ex. pop up a message box, and return a signal when closed
var _saying_text: Callable

# state of an object has changed
var _changing_state: Callable

# changing hierarchy
var _reparenting_objects: Callable

func starting_frame(starting: bool) -> void:
	return self._starting_frame.call(starting)

func saying_text(text: String): #-> Signal or nothing
	return self._saying_text.call(text)

func changing_state(noun: String, aspect: String, trait: String): #-> Signal or nothing
	return self._changing_state(noun, aspect, trait)

func reparenting_objects(pid: String, cid: String): #-> Signal or nothing
	return self._reparenting_objects.call(starting)

static func NewDefault():
	var story = TapStory.new()
	story._starting_frame = func(starting: bool): 
		print("starting frame: %s" % starting)
	story._saying_text = func(text: String): 
		print("saying text: %s" % text)
	story._changing_state = func(noun: String, aspect: String, trait: String):
		print("changing states: %s<-%s" % [noun, aspect, trait])
	story._reparenting_objects = func(pid: String, cid:String): 
		print("reparenting objects: %s<-%s" % [pid, cid])
	return story
