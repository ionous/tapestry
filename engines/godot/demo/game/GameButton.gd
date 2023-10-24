extends Button

@export var use_custom_text : bool

func _ready():
	var scene = find_parent("GameView") # fix?
	var obj: TapObject = TapRequires.find_obj(self)
	self.pressed.connect(func(): scene.notify(obj.id, "clicked"))
	#
	if not use_custom_text:
		self.set_text(obj.name if obj else "???")
