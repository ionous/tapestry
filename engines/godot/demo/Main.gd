extends Control

#@onready var io = find_child("Tapestry")
@onready var input : LineEdit = find_child("TextInput")
@onready var output : RichTextLabel = find_child("TextOutput")
@onready var scroll : ScrollContainer = find_child("ScrollContainer")
const TextWriter = preload("res://TextWriter.gd")

# Called when the node enters the scene tree for the first time.
func _ready():
	input.editable = true
	input.text_submitted.connect(_on_input)
	_post("restart", "cloak")

# scroll to the end whenever the end changes
# ( there's some sort of paint delay after append_text; but this works fine )
var last_max 
func _process(_delta):
	var vbar = scroll.get_v_scroll_bar()
	if last_max != vbar.max_value:
		vbar.value = vbar.max_value
		last_max = vbar.max_value
	
func _on_input(text):
	assert(input.editable)
	output.append_text("> " + text.replace("[", "[lb]") + "\n")
	_post("query", [{"FromExe:": {"Fabricate input:":text}}])

func _post(endpoint: String, msg: Variant) -> bool:
	var s = JSON.stringify(msg)
	var res = Tapestry.post(endpoint, s)
	if res:
		var out = _handle_response(res)		 
		var bb = TextWriter.WriteText(out + "<p>")
		output.append_text(bb)
	return res != null

# msgs expects a an array of tapestry commands, each a dictionary.
func _handle_response(msgs: Array) -> String:
	var out: String = ""
	for msg in msgs:
		var cmd = _parse_cmd(msg)
		var sig = cmd[0]
		var body = cmd[1]
		match sig:
			# TODO: look for "SceneStarted" containing the scene we want
			"Frame result:events:error:":
				# var res = body[0]; #-> results from a query
				var events = body[1]
				var err = body[2]
				for evt in events:
					out += _process_event(evt)
				print("error", err) # trace out the error
			"Frame result:events:":
				# var res = body[0]; #-> results from a query
				var events = body[1]
				for evt in events:
					out += _process_event(evt)
			_:
				push_error("unexpected message", cmd)
	return out

func _process_event(evt: Variant) -> String:
	var out: String = ""
	var cmd = _parse_cmd(evt)
	var sig = cmd[0]
	var body = cmd[1]
	match sig:
		"FrameOutput:":
			out += body
		"StateChanged noun:aspect:trait:":
			var noun = body[0]
			var aspect = body[1]
			var traitn = body[2] # doesn't like "trait"???
			print("state changed: '", noun, "' '", aspect, "' '", traitn,"'")
		_:
			print("unhandled event", sig)
	return out

# given a valid tapestry command, return its signature and body in an array of two elements
func _parse_cmd(op: Variant) -> Array:
	var pair: Array
	for k in op:
		if k != "--":
			pair = [k, op[k]]
			break
	return pair
