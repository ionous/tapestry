extends ColorRect
class_name GameView
# clickable display region for the game
# things like GameButton look for this to handle interactivity.

var _current : PopupMenu 
var _actions : Array

@export var tap_game : TapGame

# children ( like GameButton ) send these to us via "notify()"
class CustomEvent extends RefCounted:
	var name: String
	var object: String
	var data: Variant

	func _init( id_: String, event_: String, data_: Variant ):
		self.object = id_
		self.name = event_
		self.data = data_

signal custom_events(event: CustomEvent)

func _ready():
	set_mouse_filter(MOUSE_FILTER_STOP)
	custom_events.connect(_on_event)

# queue an event for processing
func notify(obj: String, event: String, data: Variant = null):
	var evt = CustomEvent.new(obj, event, data)
	custom_events.emit(evt)

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

# handle events from notify()
func _on_event(evt: CustomEvent):
	if evt.name == "clicked" and not _current:
		var objid = evt.object
		print("got %s from %s" % [evt.name, objid ])
		# each action has a name, and the kinds it can target
		var pop: PopupMenu = PopupMenu.new()
		_actions = ActionService.get_object_actions(objid )
		for act in _actions:
			# each icon has a text label and an "icon"
			var icon = IconService.find_icon(act.name)
			# fix: are there tooltips in godot?
			var label = "%s %s" % [icon.icon, icon.label]
			pop.add_item(label)

		pop.position = get_global_mouse_position()
		add_child(pop)
		pop.index_pressed.connect(func(idx):
			var act = _actions[idx]
			var text = act.format(objid)
			tap_game.fabricate(text)
			print(text) # should also go in console
			close_popup()
		)
		pop.show()
		_current= pop

