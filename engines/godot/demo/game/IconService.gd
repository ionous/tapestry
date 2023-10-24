# do we have to extend anything?
class_name IconService

class Icon extends RefCounted:
	var action: String
	var label: String
	var icon: String
	func _init(action_: String, label_: String, icon_: String):
		action = action_
		label = label_
		icon = icon_


# return the first matching icon, or an error placeholder for unknown actions
static func find_icon(action: String) -> Icon:
	for i in icons:
		if i.action == action:
			return i
	return Icon.new(action, action, "😕")


#################################################################
# in alice, these were font-awesome.....
static var icons = [
	Icon.new("looking", "Look!", "👁"),
	Icon.new("sniffing", "Sniff", "👃"),
	Icon.new("jumping", "Jump", "🏃" ),

	# world actions
	Icon.new("traveling", "Go", "📲"), # 🚪, 📲, ↩️
	Icon.new("examining", "Examine", "🔎"),
	Icon.new("taking", "Take", "🤏"),

	# multiple object actions
	Icon.new("storing", "Store", "👆"),
]
