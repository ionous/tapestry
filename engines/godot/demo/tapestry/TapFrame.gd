extends RefCounted
class_name TapFrame

# same as struct Frame in go.
var result: Variant
var events: Array   # reports any state changes
var error: String

func _init(result_: String, events_: Variant, error_: String):
	self.events = events_ if events_ is Array else [] # can be nil
	# ick: we debug.Stringify the results to support "any value"
	# so we have to unpack that too.
	self.result = JSON.parse_string(result_) if result_ else null
	self.error = error_

static func New(sig: String, body: Variant) -> TapFrame:
	if not sig.begins_with("Frame result:"):
		push_error("unhandled message", sig)
		return null
	var args: Array = body as Array
	return TapFrame.new(
		args[0],   # result
		args[1],   # events
		args[2] if args.size() > 2 else "" # error
	)
