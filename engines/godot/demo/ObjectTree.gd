extends Tree

@export var tap_game : TapGame

func _ready():
	assert(tap_game, "tap game reference not set")

func _on_tap_game_root_changed(id: String):
	clear()
	build(tap_game.pool, null, tap_game.pool.get_by_id(id))

func build(pool: TapPool, parent: TreeItem, obj: TapObject):
	var item = create_item(parent)
	item.set_text(0, obj.name)
	for k in obj.kids:
		var child = pool.get_by_id(k)
		if child:
			build(pool, item, child)

