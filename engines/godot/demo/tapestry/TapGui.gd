extends RefCounted
class_name TapGui
# abstract interface for TapGame
# so that the semi-generic TapGame can interact with the user in game specific ways

# pop up a popup, and return a signal when closed
var _display_text: Callable

# called while the game is processing previous input
var _block_input: Callable

func display_text(text: String):
	return self._display_text.call(text)

func block_input(blocked: bool) -> void:
	return self._block_input.call(blocked)

static func NewDefault():
	var gui = TapGui.new()
	gui._display_text = func(text: String): print("display text: %s" % text)
	gui._block_input = func(blocked: bool): print("block input: %s" % blocked)
	return gui
