extends RefCounted
class_name Action

var name: String
var kind: String # kind
var other_kind: String
var _do: String
var _includes: Array[String]
var _excludes: Array[String]
var _filter: Callable

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

func filter(filter_: Callable) -> Action:
	assert(not _filter)
	_filter = filter_
	return self

func do(x: String) -> Action:
	self._do = x
	return self

func matches(obj: TapObject) -> bool:
	if not _filter or _filter.call(obj):
		if (kind == obj.kind or kind == "things") and other_kind == "":
			return obj.includes(_includes) and obj.excludes(_excludes)
	return false

func format(id: String = "", otherId: String = "") -> String:
	# FIX! shouldnt send "fabricate" for these
	# send a command like Fabricate:action: instead
	# ( so that we still know to pass time )
	var args = [ id, otherId ].slice(0, _do.count("{"))
	return _do.format(args)

##################################################################
# static methods

static func get_player_actions() -> Array[Action]:
	return actions.filter(func(a): return (a.kind == ""))

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
	# can't open or close any of the doors in cloak of darkness anyway.
	# Action.new("open", "openers"),
	# Action.new("close", "openers"),

	# self actions
	Action.new("looking").do("look"),
	Action.new("sniffing").do("sniff"),
	Action.new("jumping").do("jump"),

	# world actions
	Action.new("traveling", "doors")
		.do("go {0}"),
	Action.new("examining", "things")
		.do("x {0}"),
	Action.new("taking", "things")
		.include("portable")
		# dont offer to take things the player already has
		.filter(func(obj:TapObject): return obj.parentId != TapPool.player)
		.do("take {0}"),

	Action.new("removing", "things")
		.include("worn")
		.do("take off {0}"),

	Action.new("dropping", "things")
		# fix: alice had a "carried" state; might be worth bringing that back.
		.filter(func(obj:TapObject): return obj.parentId == TapPool.player)
		# fix: the stdlib doesnt support implicit removal yet.
		.excludes("worn")
		.do("drop {0}"),

	# multiple object actions
	Action.new("storing", "things", "things")
		# fix: the stdlib correctly disallows this,
		# but the response is... unusual.
		.excludes("doors")
		.do("put {1} on {0}"),
]
