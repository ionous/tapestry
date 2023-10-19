class_name TapObjectPool
extends Resource
# A collection of references to tapestry objects

# test data: uses the same json format as the Tapestry object collection record
var mockObjects = {
  id= "tower",
  name= "tower",
  kind= "rooms",
  traits= ["lit"],
  kids= [{
		id= "mirror",
		name= "mirror",
		kind= "things",
		traits= ["fixed in place"],
		kids= []
  },{
		id= "apple",
		name= "apple",
		kind= "containers",
		traits= ["open","portable"],
		kids= [{
		  id= "worm",
		  name= "worm",
		  kind= "actors",
		  traits= ["unhappy"],
		  kids= []
		}]
  }]
}

# dictionary key->[obj]
var _all = {}
var _root: TapObject

func _init():
	_root = rebuild(mockObjects)

func get_root() -> TapObject:
	return _root

func get_by_id(id: String) -> TapObject:
	return _all[id]

func ensure(id: String) -> TapObject:
	var obj = _all.get(id)
	if !obj: 
		obj = TapObject.new(id)
		_all[id] = obj
	return obj

# rebuild from the passed collection
func rebuild(collection: Variant, parent: String = "") -> TapObject:
	var obj = ensure(collection.id)
	obj.parent = parent
	obj.kind = collection.kind
	obj.traits = collection.traits
	var kids = collection.get("kids")
	if kids:
		for kid in kids:
			var child = rebuild(kid, obj.id)
			obj.kids.push_back(child.id)
	return obj
