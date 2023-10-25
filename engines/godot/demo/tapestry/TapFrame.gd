extends RefCounted
class_name TapFrame

var events: Array
var error: String
var _callback: Variant # null or Callable
var _result: Variant # json friendly object; sent to the callback

func _init(result_: String, events_: Variant, error_: String, cb_: Variant):
	assert( (not cb_) or (cb_ is Callable) )
	self.events = events_ if events_ is Array else [] # can be nil
	self.error = error_ 
	self._callback = cb_
	# ick: we debug.Stringify the results to support "any value"
	# so we have to unpack that too.
	self._result = JSON.parse_string(result_) if result_ else null

func report_result():
	if self._callback:
		self._callback.call(self._result)

static func New(sig: String, body: Variant, cb: Variant) -> TapFrame:
	if not sig.begins_with("Frame result:"):
		push_error("unhandled message", sig)
		return null
	var args: Array = body as Array
	return TapFrame.new(
		args[0],
		args[1], 
		args[2] if args.size() > 2 else "", 
		cb)
