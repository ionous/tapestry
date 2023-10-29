extends ColorRect
class_name GameView
# clickable display region for the game
# things like GameButton look for this to handle interactivity.

var _current : PopupMenu 

@export var tap_game : TapGame

func _ready():
	assert(tap_game, "tap_game game reference not set")
	set_mouse_filter(MOUSE_FILTER_STOP)

func close_popup():
	if _current:
		_current.queue_free()
		_current= null

# handle godot input events
func _gui_input(event):
	if event is InputEventMouseButton:
		if event.button_index == MOUSE_BUTTON_LEFT and event.pressed:
			print("clicked outside")
			close_popup()

# called by other nodes to display an action bar
# requires pairs of (label, callable)
func show_popup(labelCalls: Array):
	if tap_game.is_running_turn():
		return
	close_popup()
	# build popup menu
	var pop: PopupMenu = PopupMenu.new()
	_current= pop
	# add items, and record callbacks
	var calls: Array[Callable] = [] 
	for labelCall in labelCalls:
		var label: String= labelCall[0]
		var cb: Callable= labelCall[1]
		calls.push_back(cb) 
		pop.add_item(label)
	# handle clicks:
	pop.index_pressed.connect(func(idx):
		var cb = calls[idx]
		if cb:
			cb.call()
		close_popup()
	)
	# show popup
	pop.position = get_global_mouse_position()
	add_child(pop)
	pop.show()
	
