extends HBoxContainer

var showing_inventory: bool 
@export var tap_game : TapGame

# Called when the node enters the scene tree for the first time.
func _ready():
	assert(tap_game)
	build_player_actions()

func build_inventory():
	remove_all()
	var p: TapObject = TapPool.get_by_id(TapPool.player)
	if p.childIds.is_empty():
		var b = add_button("empty", "", func(): pass)
		b.disabled = true
	else:
		for cid in p.childIds:
			var kid: TapObject = TapPool.get_by_id(cid)
			var icon = cloakIcon if cid == "velvet cloak" else unknownIcon
			add_button(icon, kid.name, func(): _on_item(kid))
	# add the close icon:
	add_button(closeIcon, "", func(): build_player_actions())

const unknownIcon = "❔"
const cloakIcon = "🧥"
const inventoryIcon = "🎒"
const closeIcon = "❌"

func build_player_actions():
	remove_all()
	for act in ActionService.get_player_actions():
		var icon = IconService.find_icon(act.name)
		add_button(icon.icon, icon.label, func(): _on_action(act))
	add_button(inventoryIcon, "Inv", func(): build_inventory())

func remove_all():
	while get_child_count() > 0:
		remove_child(get_child(0))

func add_button(icon: String, text: String, cb: Callable):
	var b = Button.new()
	b.text = (icon + " " + text) if text else icon
	b.pressed.connect(cb)
	add_child(b)
	return b

func _on_item(obj: TapObject):
	print("clicked %s" % [obj.name])

func _on_action(act: ActionService.Action):
	tap_game.fabricate(act.format())
	print("clicked %s" % [act.name])

func _on_main_turn_started(started):
	# if the control isn't visible, it doesnt take up space :/
	# ( which means the other controls move around and it looks weird )
	# visible = !started
	for i in get_child_count():
		get_child(i).visible = !started

