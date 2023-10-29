extends Tree
# needs to be connected to the tapestry object so that it can listen for rebuild changes

var _root: TapObject

# whenever the room changes, rebuild the tree view
func _on_main_room_changed(id: String):
	_root = TapPool.get_by_id(id)
	build(_root)

# at the end of every turn, rebuild the tree view
# ( this in lieu of a good object is dirty check )
func _on_main_turn_started(_started):
	if not _started:
		build(_root)

func build(obj: TapObject):
	clear()
	_build(null, obj)

func _build(parentItem: TreeItem, obj: TapObject):
	var item = create_item(parentItem)
	item.set_text(0, obj.name)
	for k in obj.childIds:
		var child = TapPool.get_by_id(k)
		if child:
			_build(item, child)


