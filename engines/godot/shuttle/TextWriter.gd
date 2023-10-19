class_name TapWriter

static func WriteText(text):
	var w = TapWriter.new()
	return w.writeText(text)
	
# tag types ( slightly easier to read then true/false)
enum TagType {
	opening, closing
}

# turn tagged text into bb code
class TextRender:
	# return false, or the data to help generate the bbcode.
	static func selectTag(tag):
		return {"tag":tag} if permit.has(tag) else tagTransforms.get(tag)

	# permit these tags to become bb code tags
	static var permit: Array = [ "b", "i", "s", "u", "ol", "ul" ]

	# transform a tag into some text
	static var tagTransforms :Dictionary = {
		"p":   {"char": '\v'},
		"br":  {"char": '\n'},
		"wbr": {"char": '\r'},
		"hr":  {"text": "-------------------------------", "char": '\r'},
		"li":  {"text": ""} # note: bbcode ol and ul only provide one line items... that may cause some troubles.
	}

	# track the text inside the span of a starting tag, and its matching end tag.
	class Span:
		func _init(t:String = ""):
			tag = t
		var tag : String
		var text: String
		
	# current tag and text between its start and end
	var curr  : Span = Span.new()
	# tracks the begining of spans
	# when a span is closed by its matching tag
	# it is converted into text and added to the previous span
	var stack : Array = []
	# vertical spacing ( ie. are we at the start of a new line )
	# 2 indicates we are starting the top of a new section.
	var line  : int = 2

	func finalize():
		# we're at the end of all things
		# close any and all tags that were still open
		while !stack.is_empty():
			closeTag(curr.tag)
		# at the end, curr is the final span
		var out : String = curr.text
		curr = null # will cause an error if we try to use this again.
		return out

	# return true if the tag was known
	func writeTag(tag, tagType):
		var rule = TextRender.selectTag(tag)
		if rule:
			if "text" in rule:
				curr.text += rule["text"]
			if "char" in rule:
				writeChar(rule["char"])
			if "tag" in rule:
				var ntag = rule["tag"] # transformed tag
				if  tagType == TagType.opening: # tbd: are ternaries early out?
					openTag(ntag)
				else:
					closeTag(ntag)
		# return true if the tag was known
		return true if rule else false

	func writeString(s: String):
		for q in s:
			writeChar(q)

	# watches for newline characters
	# and filters unexpected bbcode brackets
	func writeChar(q):
		match q:
			'\n': # newline
				curr.text += q
				line += 1
			'\r': # softline
				while line < 1:
					curr.text += "\n"
					line += 1
			'\v': # softparagraph
				while line < 2:
					curr.text += "\n"
					line += 1
			'[': # https://docs.godotengine.org/en/stable/tutorials/ui/bbcode_in_richtextlabel.html#handling-user-input-safely
				curr.text += "[lb]"
			_: # default
				curr.text += q
				line = 0

	# start a new span within the current span
	func openTag(tag):
		stack.append(curr)
		curr = Span.new(tag)

	# end the current span but only if the tag matches, otherwise the tag is ignored.
	# (writes the current span into its parent span)
	func closeTag(tag):
		if curr.tag == tag:
			var inner = curr
			curr = stack.pop_back()
			curr.text += "[{tag}]{span}[/{tag}]".format({"tag": inner.tag, "span": inner.text})


# these states act on the text builder
enum States {
	readingText,
	openingTag,
	closingTag
}


# accumulates characters into potential tags....
# ( because we the whole tag word to process it )
var buf : String
var pendingTag : String
var out : TextRender = TextRender.new()
static var longestTag : int = "blockquote".length()
static var tagFriendly : RegEx = RegEx.create_from_string("[a-zA-Z]")
	

# tbd: was having trouble putting this as a static function
func writeText(text):
	var nextState = States.readingText
	for q in text:
		match nextState:
			States.readingText:
				nextState= self.readingText(q)
			States.openingTag:
				nextState= self.openingTag(q)
			States.closingTag:
				nextState= self.closingTag(q)
			_:
				assert(false, "we shouldnt be here")
	self.rejectTag()  # if anything was pending
	return out.finalize() # flatten things

# the potential tag we've been accumulating is not a tag:
# write any text we've been accumulating to the output
func rejectTag():
	out.writeString(buf)
	buf = ""
	pendingTag = ""

# keep track of the passed character as potentially belonging to a tag
func accumTag(q):
	if pendingTag.length() < longestTag:
		pendingTag += q
	else:
		rejectTag()
	return pendingTag.length() > 0

# done with reading a <tag>
# ( unknown tags will be eaten )
func dispatchTag(tagType):
	if !out.writeTag(pendingTag, tagType):
		rejectTag()
	else:
		# we're done with tag, so clear our tracking data.
		buf = ""
		pendingTag = ""

# reading normal text
# changes to "openingTag" when a "<" is seen.
func readingText(q):
	var nextState = States.readingText
	# detected the start of a tag:
	if q != '<':
		self.out.writeChar(q)
	else:
		self.buf += q; # start buffering in case this isnt a valid tag.
		nextState = States.openingTag
	return nextState

# reading some text that might be an opening or closing tag.
# changes to "closingTag" when a "</" is seen
# changes to "readingText" at the end of an opening tag ">"
func openingTag(q):
	var nextState= States.openingTag
	self.buf += q
	# could be a closing tag
	if q=='/':
		# yes. looks like a closing tag...
		if !self.pendingTag.length():
			nextState = States.closingTag

		# no: there was text before the slash:
		# ex. <abc/> not </abc>
		else:
			self.rejectTag()
			nextState= States.readingText

	# end of this opening tag.
	elif q == '>':
		self.dispatchTag(TagType.opening)
		nextState= States.readingText

	# continuing on in the current tag
	# ex. <abc...
	elif tagFriendly.search(q):
		# ( only fails to accumulate if the tag was too long )
		if (!self.accumTag(q)):
			nextState= States.readingText

	# some other character
	# ( therefore not a tag )
	else:
		self.rejectTag()
		nextState= States.readingText

	return nextState

# reading some text of a closing tag.
# changes to "readingText" at the end the tag.
func closingTag(q):
	var nextState= States.closingTag
	self.buf += q
	# end of this closing tag.
	if q == '>':
		self.dispatchTag(TagType.closing)
		nextState= States.readingText
	# continuing on in the current tag
	# ex. </abc...
	elif tagFriendly.search(q):
		if !self.accumTag(q):
			# ( only fails to accumulate if the tag was too long )
			nextState= States.readingText
	# some other character
	# ( therefore not a tag )
	else:
		self.rejectTag()
		nextState= States.readingText

	return nextState
