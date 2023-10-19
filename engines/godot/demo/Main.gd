extends Control

@onready var input : LineEdit = find_child("TextInput")
@onready var output : RichTextLabel = find_child("TextOutput")
@onready var scroll : ScrollContainer = find_child("ScrollContainer")
@onready var tap : TapGame = find_child("TapGame")
var objects : TapObjectPool = TapObjectPool.new()

# Called when the node enters the scene tree for the first time.
func _ready():
	input.editable = true
	input.text_submitted.connect(_on_input)
	# objects.root_changed.connect
	# tap.title_changed.connect
	# tap.score_changed.connect
	# tap.turns_changed.connect

	# after everything is connected:
	tap.restart("cloak", objects)

# scroll to the end whenever the end changes
# ( there's some sort of paint delay after append_text; but this works fine )
var last_max 
func _process(_delta):
	var vbar = scroll.get_v_scroll_bar()
	if last_max != vbar.max_value:
		vbar.value = vbar.max_value
		last_max = vbar.max_value
	
func _on_input(text: String):
	assert(input.editable)
	output.append_text("> " + text.replace("[", "[lb]") + "\n")
	tap.fabricate(text)
