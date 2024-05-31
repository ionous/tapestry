extends Node
# Tapestry command queries for title, score, and more.

class Cmd extends RefCounted:
	var sig: String
	var body: Variant # the content type depends on the signature
	func _init(sig_: String, body_: Variant):
		self.sig = sig_
		self.body = body_

# given a valid tapestry command:
# return its signature and body
func Parse(op: Variant) -> Cmd:
	for k in op:
		if k != "--":
			return Cmd.new(k, op[k])
	return null

# --------------------------------------------------
# generate a command using the passed player input
# --------------------------------------------------
func Fabricate(input: String) -> Dictionary:
	return {
		"FromExe:": {
			"Fabricate input:":input
		}
	}

# ---------------------------------------------------
const StoryTitle =	{
	"FromText:": {
		"Object:dot:": ["story", "title"]
	}
}

# --------------------------------------------------
const PlayerLocation = {
	"FromText:": {
		"Determine:args:": [
			"location_of", {
				"Arg:from:": [
					"obj", {
						"FromText:": {
						"Object:dot:": ["story", "actor"]
					}
				}
			]}
	]}
}

# --------------------------------------------------
const CurrentObjects = {
	"FromRecord:": {
		"Determine:args:": [
			"collect_objects", {
				"Arg:from:": ["obj", PlayerLocation]
			}
		]
	}
}

# --------------------------------------------------
const CurrentScore = {
	"FrmoNum:": {
	"Num if:then:else:": [
		{ "Is domain:": "scoring" },
		{"Object:dot:": ["story", "score"]},
		-1
	]}
}

# --------------------------------------------------
const CurrentTurn = {
	"FrmoNum:": {
		"Num if:then:else:": [
			{ "Is domain:": "scoring" },
			{"Object:dot:": ["story", "turn count"]},
			-1
		]
	}
}

# --------------------------------------------------
# fix: what do you mean this is insane?
# it'd help if we could use the implicit pattern call decoder :/
# maybe change "print name" to some "get name"
# and, if possible, get rid of From(s)
const LocationName = {
	"FromText:": {
		"Determine:args:": [
			"print_name", {
				"Arg:from:": [
				"obj",	PlayerLocation
				]
			}
		]
	}
}

