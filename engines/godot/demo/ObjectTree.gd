extends Tree

func _on_tap_game_root_changed(pool:TapObjectPool, root:TapObject):
	clear()
	build(pool, null, root)

func build(pool:TapObjectPool, parent:TreeItem, obj:TapObject):
	var item = create_item(parent)
	item.set_text(0, obj.name)
	for k in obj.kids:
		var child = pool.get_by_id(k)
		if child:
			build(pool, item, child)

