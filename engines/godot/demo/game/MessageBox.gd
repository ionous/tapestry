extends PopupPanel

@export var label  : RichTextLabel

func _init():
	popup_hide.connect(func(): self.queue_free())

# popup_hide
var dialog_text : String:
	get:
		return label.text
	set(text):
		label.text= text
