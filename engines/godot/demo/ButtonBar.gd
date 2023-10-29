extends HBoxContainer

@export var tap_game : TapGame

signal clicked_action(act: Action)
signal clicked_item(obj: TapObject)

var _showing_inventory: bool = false

# Called when the node enters the scene tree for the first time.
func _ready():
	assert(tap_game)
	_build_actions()

func rebuild():
	return _build_inventory() if _showing_inventory else _build_actions()

func _build_inventory():
	_remove_all()
	var p: TapObject = TapPool.get_by_id(TapPool.player)
	if p.childIds.is_empty():
		var b = _add_button("empty", "", func(): pass)
		b.disabled = true
	else:
		for cid in p.childIds:
			_add_item(cid)
	# add the close icon:
	_add_button(Icons.Close, "", func(): _build_actions())
	_showing_inventory = true

func _build_actions():
	_remove_all()
	for act in Action.get_player_actions():
		var icon = Icons.find_icon(act.name)
		_add_button(icon.icon, icon.label, func(): clicked_action.emit(act))
	_add_button(Icons.Inventory, "Inv", func(): _build_inventory())
	_showing_inventory = false

func _remove_all():
	while get_child_count() > 0:
		remove_child(get_child(0))

func _add_item(id: String):
	var kid: TapObject = TapPool.get_by_id(id)
	var icon = Icons.CloakOfDarkness if id == "velvet cloak" else Icons.Unknown
	return _add_button(icon, kid.name, func(): clicked_item.emit(id))

func _add_button(icon: String, text: String, cb: Callable):
	var b = Button.new()
	b.text = (icon + " " + text) if text else icon
	b.pressed.connect(cb)
	add_child(b)
	return b

func _on_main_turn_started(started):
	# if the control isn't visible, it collapses
	# so a min size is set, and this hides the children instead.
	if not started:
		rebuild()
	else:
		for i in get_child_count():
			get_child(i).disabled = true

func _on_main_selected_item(itemId: String):
	if not itemId:
		_build_actions()
	else:
		_remove_all()
		var b = _add_item(itemId)
		b.disabled = true
		b = _add_button("use the item with an object in the world", "", func(): pass)
		b.disabled = true
