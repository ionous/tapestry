extends Tree

func _on_tap_game_root_changed(id: String):
	clear()
	build(null, TapPool.get_by_id(id))

func build(parentItem: TreeItem, obj: TapObject):
	var item = create_item(parentItem)
	item.set_text(0, obj.name)
	for k in obj.childIds:
		var child = TapPool.get_by_id(k)
		if child:
			build(item, child)

