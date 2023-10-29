extends Button

@export var use_custom_text : bool

func _ready():
	var main = find_parent("Main") # fix?
	var obj: TapObject = TapRequires.find_obj(self)
	self.pressed.connect(func(): main._on_clicked_object(obj.id))
	#
	if not use_custom_text:
		self.set_text(obj.name if obj else "???")
