extends Label

var score = 0
var turns = -1

func use_scoring() -> bool:
	return turns > 0

func set_score(v):
	if score != v:
		score = v 
		set_process(true)

func set_turns(v):
	if turns != v :
		turns = v 
		set_process(true)
	
func _process(_delta):
	if !use_scoring():
		self.visible = false 
	else:
		set_text("%d / %d" % [ score, turns ] ) # godot compares to see if it changed....
	set_process(false)

