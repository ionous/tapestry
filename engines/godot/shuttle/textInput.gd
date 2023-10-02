extends LineEdit

# by default isnt editable.
func _ready():
	grab_focus()
	editable = false 	

# other nodes also listen to this signal
func _on_text_submitted(_str):
	clear()
