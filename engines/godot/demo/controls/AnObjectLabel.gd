extends Label
class_name AnObjectLabel

# Called when the node enters the scene tree for the first time.
func _ready():
	var obj: TapObject = TapRequires.find_obj(self)
	self.set_text(obj.name if obj else "???")
	pressed.connect(_pressed)

func _pressed():
	
