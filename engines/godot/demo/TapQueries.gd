extends Node
# Tapestry command queries for title, score, and more.

const PlayerLocation = {
	"FromText:": {
		"Determine:args:": [
			"location_of", {
				"Arg:from:": [
					"obj", {
						"FromText:": {
						"Object:field:": ["story", "actor"]
					}
				}
			]}
	]}
}

const CurrentObjects = {
	"FromRecord:": {
		"Determine:args:": [
			"collect_objects", {
				"Arg:from:": ["obj", PlayerLocation]
			}
		]
	}
}

const CurrentScore = {
	"FromNumber:": {
	"Num if:then:else:": [
		{ "Is domain:": "scoring" },
		{"Object:field:": ["story", "score"]},
		-1
	]}
}

const CurrentTurn = {
	"FromNumber:": {
		"Num if:then:else:": [
			{ "Is domain:": "scoring" },
			{"Object:field:": ["story", "turn count"]},
			-1
		]
	}
}

# # fix: what do you mean this is insane?
# # it'd help if we could use the implicit pattern call decoder :/
# # maybe change "print name" to some "get name"
# # and, if possible, get rid of From(s)
const LocationName = {
	"FromText:": {
		"Buffers do:": {
			"Determine:args:": [
				"print_name", {
					"Arg:from:": [
					"obj",	PlayerLocation
					]
				}
			]
		}
	}
}

const StoryTitle =	{
	"FromText:": {
		"Object:field:": ["story", "title"]
	}
}


static func Fabricate(text: String) -> Dictionary:
	return {
		"FromExe:": {
			"Fabricate input:":text
		}
	}
