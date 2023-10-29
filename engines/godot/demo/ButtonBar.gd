extends HBoxContainer

var showing_inventory: bool 
@export var tap_game : TapGame

signal clicked_action(act: ActionService.Action)
signal clicked_item(obj: TapObject)


# Called when the node enters the scene tree for the first time.
func _ready():
	assert(tap_game)
	build_actions()

func build_inventory():
	remove_all()
	var p: TapObject = TapPool.get_by_id(TapPool.player)
	if p.childIds.is_empty():
		var b = add_button("empty", "", func(): pass)
		b.disabled = true
	else:
		for cid in p.childIds:
			var kid: TapObject = TapPool.get_by_id(cid)
			var icon = Icons.CloakOfDarkness if cid == "velvet cloak" else Icons.Unknown
			add_button(icon, kid.name, func(): clicked_item.emit(cid))
	# add the close icon:
	add_button(Icons.Close, "", func(): build_actions())


func build_actions():
	remove_all()
	for act in ActionService.get_player_actions():
		var icon = Icons.find_icon(act.name)
		add_button(icon.icon, icon.label, func(): clicked_action.emit(act))
	add_button(Icons.Inventory, "Inv", func(): build_inventory())

func remove_all():
	while get_child_count() > 0:
		remove_child(get_child(0))

func add_button(icon: String, text: String, cb: Callable):
	var b = Button.new()
	b.text = (icon + " " + text) if text else icon
	b.pressed.connect(cb)
	add_child(b)
	return b

func _on_main_turn_started(started):
	# if the control isn't visible, it collapses
	# so a min size is set, and this hides the children instead.
	for i in get_child_count():
		var b = get_child(i)
		if b.text != "empty":
			b.disabled = started

