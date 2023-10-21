class_name TapRequires
extends RefCounted

# tapestry object id 
# determined from the owner or a parent owner
var id: String

# uses duck typing; expects owner to have:
# object_id, include_traits, exclude_traits, include_ancestor, exclude_ancestor
var owner: Node

# requires that the parent req(s) if any are valid
# ( so for nodes, expects to be created in "enter tree" )
func _init(owner_: Node):
	owner = owner_
	id = owner.object_id
	if not id:
		var n = owner.get_name()
		if n.begins_with("#"):
			id = n.right(-1)
		else:
			id = TapRequires.find_id(owner.get_parent())
			if not id:
				push_error("couldnt find object id for '%s'" % owner.name)

# find the nearest requirement object
static func find_obj(search: Node) -> TapObject:
	var pr = find_req(search)
	return pr.get_object() if pr else ""

# find the nearest requirement id
static func find_id(search: Node) -> String:
	var pr = find_req(search)
	return pr.id if pr else ""

# find the nearest requirement
static func find_req(search: Node) -> TapRequires:
	var out: TapRequires
	while search:
		var pr = search.get("requires") as TapRequires
		if not pr:
			search = search.get_parent()
		else:
			out = pr
			break
	return out

func get_object() -> TapObject:
	return TapPool.get_by_id(id) if id else null

func check_requirements() -> bool:
	var myobj = get_object()
	return (myobj and _include_traits(myobj) 
				and _exclude_traits(myobj) 
				and _contained(myobj)
				and _include_ancestor(myobj)
				and _exclude_ancestor(myobj))

# false if one of the include_traits is missing
func _include_traits(myobj: TapObject) -> bool:
	var okay = true # provisionally
	var includes: Array = owner.get("include_traits")
	for include in includes:
		if not include in myobj.traits:
			okay = false
			break
	return okay

# false if one of the exclude_traits is present
func _exclude_traits(myobj: TapObject) -> bool:
	var okay = true # provisionally
	var excludes: Array = owner.get("exclude_traits")
	for exclude in excludes:
		if exclude in myobj.traits:
			okay = false
			break
	return okay

# false if the nearest TapRequired object_id isn't an ancestor
func _contained(myobj: TapObject) -> bool:
	var contained: bool = owner.get("contained")
	return (not contained) or myobj.has_ancestor(TapRequires.find_id(owner.get_parent()))

# false if the include_ancestor is missing
func _include_ancestor(myobj: TapObject) -> bool:
	var include: String = owner.get("include_ancestor")
	return (not include) or myobj.has_ancestor(include)

# false if the exclude_ancestor is present
func _exclude_ancestor(myobj: TapObject) -> bool:
	var exclude: String = owner.get("exclude_ancestor")
	return (not exclude) or not myobj.has_ancestor(exclude)
