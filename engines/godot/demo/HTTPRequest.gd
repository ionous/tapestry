extends HTTPRequest

const endpoint: String = "http://127.0.0.1:8080/shuttle/"
var callback : Callable

func sendInput(text: String, cb: Callable):
	_send({"in":text}, cb)

func sendCommand(cmd: String, val: Variant, cb: Callable):
	_send({"cmd":{cmd:val}}, cb)
		
func _send(msg: Variant, cb: Callable):
	assert(callback.is_null(), "previous request still pending")
	if callback.is_null():
		callback = cb
		var text = JSON.stringify(msg)
		request(endpoint, [], HTTPClient.METHOD_POST, text)		

func _on_request_completed(_result, _response_code, headers : PackedStringArray, body):
	assert(!callback.is_null(), "unexpected response")
	if !callback.is_null(): 
		var hasJson = headers.has("Content-Type:application/json")
		if !hasJson:
			var str = body.get_string_from_utf8()
			var res = str if !hasJson else JSON.parse_string(str) # Returns null if parsing failed.
			var cb = callback 
			callback = Callable()
			cb.call(res)
