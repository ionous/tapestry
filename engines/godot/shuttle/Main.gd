extends Control

@onready var io = find_child("Shuttle")
@onready var input : LineEdit = find_child("TextInput")
@onready var output : RichTextLabel = find_child("TextOutput")
@onready var scroll : ScrollContainer = find_child("ScrollContainer")
const TextWriter = preload("res://TextWriter.gd")


enum State { STARTING, PLAYING, WAITING, ENDED }
var state : State  = State.STARTING

# Called when the node enters the scene tree for the first time.
func _ready():
	assert(state == State.STARTING)
	input.editable = false 
	input.text_submitted.connect(_on_input)
	io.sendCommand("$restart", "cloak", _started)

func _on_input(text):
	if state != State.PLAYING:
		output.append_text("error %s" % state)
	else:
		state = State.WAITING
		input.editable = false
		output.append_text("> " + text.replace("[", "[lb]") + "\n")
		io.sendInput(text, _gotText)
		
func _started(initialText):
	assert(state == State.STARTING)
	if state == State.STARTING:
		var bb = TextWriter.WriteText(initialText + "<p>")
		output.append_text(bb)
		input.editable = true
		state = State.PLAYING

func _gotText(newText):
	assert(state == State.WAITING)
	if state == State.WAITING:
		var bb = TextWriter.WriteText(newText + "<p>")
		output.append_text(bb)
		input.editable = true
		state = State.PLAYING
	

# scroll to the end whenever the end changes
# ( there's some sort of paint delay after append_text; but this works fine )
var last_max 
func _process(_delta):
	var vbar = scroll.get_v_scroll_bar()
	if last_max != vbar.max_value:
		vbar.value = vbar.max_value
		last_max = vbar.max_value
