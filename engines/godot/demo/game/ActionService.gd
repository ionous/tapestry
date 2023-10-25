class_name ActionService

class Action extends RefCounted:
	var name: String
	var kind: String # kind
	var other_kind: String
	var _do: String
	var _includes: Array[String]
	var _excludes: Array[String]

	func _init(name_: String, noun_kind_: String= "", other_kind_: String= ""):
		name = name_
		kind = noun_kind_
		other_kind = other_kind_
		
	func include(state: String) -> Action:
		_includes.push_back(state)
		return self 

	func excludes(state: String) -> Action:
		_excludes.push_back(state)
		return self

	func do(x: String) -> Action:
		self._do = x
		return self

	func matches(obj: TapObject) -> bool:
		if (kind == obj.kind or kind == "things") and other_kind == "":
			return obj.includes(_includes) and obj.excludes(_excludes)
		return false

	func format(id: String) -> String:
		# FIX! shouldnt send "fabricate" if we can send the actual action to call
		# noting that ".play()" passes time
		# so either, we have to duplicate those actions in TapCommand
		# or we need to send a command like Fabricate:action:
		# the latter might be better?
		return "%s %s" % [ _do, id ]

static func get_player_actions() -> Array[Action]:
	return actions.filter(func(a): return (a.other_kind == ""))

static func get_object_actions(object_id: String) -> Array[Action]:
	var obj: TapObject = TapPool.get_by_id(object_id)
	# fix: materialized path instead of "things" hack
	return actions.filter(func(a): return a.matches(obj))
	
static func get_multi_actions(_one_id: String, _other_id: String) -> Array[Action]:
	# var obj: TapObject = TapPool.get_by_id(object_id)
	# fix: materialized path checks
	return actions.filter(func(a): return a.other_kind != "")

#################################################################
# fix. in "alice", this would query the parser for actions
# include action id, name, included noun count ( and type )
static var actions : Array[Action] = [
	# can't open or close anything in cloak anyway.
	# Action.new("open", "openers"),
	# Action.new("close", "openers"),

	# self actions
	Action.new("looking").do("look"),
	Action.new("sniffing").do("sniff"),
	Action.new("jumping").do("jump"),

	# world actions
	Action.new("traveling", "doors").do("go"),
	Action.new("examining", "things").do("x"),
	Action.new("taking", "things").include("portable").do("take"),

	# multiple object actions
	Action.new("storing", "things", "things").do("store"),
]
