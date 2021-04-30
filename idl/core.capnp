@0xdb087b68b3c8ece3;
using Go = import "/go.capnp";
using  X = import "options.capnp";
using Rtx = import "rtx.capnp";
using Reader = import "reader.capnp";

$Go.package("core");
$Go.import("git.sr.ht/~ionous/iffy/idl/core");

struct Brancher { eval @0:AnyPointer; }
struct Comparator { eval @0:AnyPointer; }
struct FromSourceFields { eval @0:AnyPointer; }
struct IntoTargetFields { eval @0:AnyPointer; }
struct Trigger { eval @0:AnyPointer; }

struct Activity $X.label("Act") $X.group("hidden") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct AllTrue $X.label("AllOf") $X.group("logic") $X.desc("Returns true if all of the evaluations are true.") {
  test         @0  :List(Rtx.BoolEval);
}
struct Always $X.label("Always") $X.group("logic") $X.desc("Returns true always.") {
}
struct AnyTrue $X.label("AnyOf") $X.group("logic") $X.desc("Returns true if any of the evaluations are true.") {
  test         @0  :List(Rtx.BoolEval);
}
struct Argument $X.label("Arg") $X.group("patterns") {
  name         @0  :Text;
  from         @1  :Rtx.Assignment $X.label("from");
}
struct Arguments $X.label("Arguments") $X.group("patterns") {
  args         @0  :List(Argument);
}
struct Assign $X.label("Let") $X.group("variables") $X.desc("Assigns a variable to a value.") {
  str          @0  :Text $X.pool("variable");
  from         @1  :Rtx.Assignment $X.label("be");
}
struct BoolValue $X.label("Bool") $X.group("literals") $X.desc("Specify an explicit true or false value.") {
  bool         @0  :Bool;
}
struct Bracket $X.label("BracketText") $X.group("printing") $X.desc("Sandwiches text printed during a block and puts them inside parenthesis '()'.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct Break $X.label("Break") $X.group("flow") $X.desc("In a repeating loop, exit the loop.") {
}
struct Buffer $X.label("BufferText") $X.group("printing") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct Capitalize $X.label("Capitalize") $X.group("format") $X.desc("Returns new text, with the first letter turned into uppercase.") {
  text         @0  :Rtx.TextEval;
}
struct ChooseAction $X.label("If") $X.group("misc") $X.desc("An if statement.") {
  if           @0  :Rtx.BoolEval;
  exe          @1  :List(Rtx.Execute) $X.label("do");
  else         @2  :Brancher $X.optional $X.label("else");
}
struct ChooseMore $X.label("ElseIf") $X.group("misc") {
  if           @0  :Rtx.BoolEval;
  exe          @1  :List(Rtx.Execute) $X.label("do");
  else         @2  :Brancher $X.optional $X.label("else");
}
struct ChooseMoreValue $X.label("ElseIf") $X.group("misc") {
  assign       @0  :Text;
  from         @1  :Rtx.Assignment $X.label("from");
  filter       @2  :Rtx.BoolEval $X.label("and");
  exe          @3  :List(Rtx.Execute) $X.label("do");
  else         @4  :Brancher $X.optional $X.label("else");
}
struct ChooseNothingElse $X.label("ElseDo") $X.group("misc") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct ChooseNum $X.label("Num") $X.group("math") $X.desc("Pick one of two numbers based on a boolean test.") {
  true         @0  :Rtx.NumberEval;
  if           @1  :Rtx.BoolEval $X.label("if");
  false        @2  :Rtx.NumberEval $X.label("else");
}
struct ChooseText $X.label("Txt") $X.group("format") $X.desc("Pick one of two strings based on a boolean test.") {
  true         @0  :Rtx.TextEval;
  if           @1  :Rtx.BoolEval $X.label("if");
  false        @2  :Rtx.TextEval $X.label("else");
}
struct ChooseValue $X.label("If") $X.group("misc") $X.desc("An if statement with local assignment.") {
  assign       @0  :Text;
  from         @1  :Rtx.Assignment $X.label("from");
  filter       @2  :Rtx.BoolEval $X.label("and");
  exe          @3  :List(Rtx.Execute) $X.label("do");
  else         @4  :Brancher $X.optional $X.label("else");
}
struct Commas $X.label("CommaText") $X.group("printing") $X.desc("Separates words with commas, and 'and'.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct CompareNum $X.label("Cmp") $X.group("logic") $X.desc("True if eq,ne,gt,lt,ge,le two numbers.") {
  a            @0  :Rtx.NumberEval;
  is           @1  :Comparator $X.label("is");
  b            @2  :Rtx.NumberEval $X.label("num");
}
struct CompareText $X.label("Cmp") $X.group("logic") $X.desc("True if eq,ne,gt,lt,ge,le two strings ( lexical. )") {
  a            @0  :Rtx.TextEval;
  is           @1  :Comparator $X.label("is");
  b            @2  :Rtx.TextEval $X.label("txt");
}
struct CountOf $X.label("CountOf") $X.group("logic") $X.desc("A guard which returns true based on a counter. Counters start at zero and are incremented every time the guard gets checked.") {
  at           @0  :Reader.Pos $X.internal;
  num          @1  :Rtx.NumberEval;
  trigger      @2  :Trigger $X.label("trigger");
}
struct CycleText $X.label("Cycle") $X.group("output") $X.desc("When called multiple times, returns each of its inputs in turn.") {
  seq          @0  :Text $X.internal;
  parts        @1  :List(Rtx.TextEval);
}
struct Determine $X.label("Determine") $X.group("patterns") $X.desc("Runs a pattern, and potentially returns a value.") {
  pattern      @0  :Text $X.pool("pattern");
  args         @1  :List(Argument) $X.label("args");
}
struct DiffOf $X.label("Dec") $X.group("math") $X.desc("Subtract two numbers.") {
  a            @0  :Rtx.NumberEval;
  b            @1  :Rtx.NumberEval $X.optional $X.label("by");
}
struct During $X.label("During") $X.group("patterns") $X.desc("Decide whether a pattern is running.") {
  pattern      @0  :Text $X.pool("pattern");
}
struct EqualTo $X.label("EqualTo") $X.group("comparison") $X.desc("Two values exactly match.") {
}
struct FromBool $X.label("Bool") $X.group("variables") $X.desc("Assigns the calculated boolean value.") {
  val          @0  :Rtx.BoolEval;
}
struct FromNum $X.label("Num") $X.group("variables") $X.desc("Assigns the calculated number.") {
  val          @0  :Rtx.NumberEval;
}
struct FromNumbers $X.label("Nums") $X.group("variables") $X.desc("Assigns the calculated numbers.") {
  vals         @0  :Rtx.NumListEval;
}
struct FromObj $X.label("Obj") $X.group("misc") $X.desc("Targets an object with a computed name.") {
  object       @0  :Rtx.TextEval;
}
struct FromRec $X.label("Rec") $X.group("misc") $X.desc("Targets a record stored in a record.") {
  rec          @0  :Rtx.RecordEval;
}
struct FromRecord $X.label("Rec") $X.group("variables") $X.desc("Assigns the calculated record.") {
  val          @0  :Rtx.RecordEval;
}
struct FromRecords $X.label("Recs") $X.group("variables") $X.desc("Assigns the calculated records.") {
  vals         @0  :Rtx.RecordListEval;
}
struct FromText $X.label("Txt") $X.group("variables") $X.desc("Assigns the calculated piece of text.") {
  val          @0  :Rtx.TextEval;
}
struct FromTexts $X.label("Txts") $X.group("variables") $X.desc("Assigns the calculated texts.") {
  vals         @0  :Rtx.TextListEval;
}
struct FromVar $X.label("Var") $X.group("misc") $X.desc("Targets a record stored in a variable.") {
  str          @0  :Text $X.pool("variable");
}
struct GetAtField $X.label("Get") $X.group("variables") $X.desc("Get a value from a record.") {
  field        @0  :Text;
  from         @1  :FromSourceFields $X.label("from");
}
struct GreaterOrEqual $X.label("AtLeast") $X.group("comparison") $X.desc("The first value is larger than the second value.") {
}
struct GreaterThan $X.label("GreaterThan") $X.group("comparison") $X.desc("The first value is larger than the second value.") {
}
struct HasDominion $X.label("HasDominion") $X.group("logic") {
  name         @0  :Text $X.pool("domain");
}
struct HasTrait $X.label("Get") $X.group("objects") $X.desc("Return true if the object is currently in the requested state.") {
  object       @0  :Rtx.TextEval $X.label("obj");
  trait        @1  :Rtx.TextEval $X.label("trait");
}
struct IdOf $X.label("IdOf") $X.group("objects") $X.desc("A unique object identifier.") {
  object       @0  :Rtx.TextEval;
}
struct Includes $X.label("Contains") $X.group("strings") $X.desc("True if text contains text.") {
  text         @0  :Rtx.TextEval;
  part         @1  :Rtx.TextEval $X.label("part");
}
struct IntoObj $X.label("Obj") $X.group("misc") $X.desc("Targets an object with a computed name.") {
  object       @0  :Rtx.TextEval;
}
struct IntoVar $X.label("Var") $X.group("misc") $X.desc("Targets an object or record stored in a variable") {
  str          @0  :Text $X.pool("variable");
}
struct IsEmpty $X.label("Is") $X.group("strings") $X.desc("True if the text is empty.") {
  text         @0  :Rtx.TextEval $X.label("empty");
}
struct IsExactKindOf $X.label("KindOf") $X.group("objects") $X.desc("True if the object is exactly the named kind.") {
  object       @0  :Rtx.TextEval;
  kind         @1  :Text $X.label("is_exactly");
}
struct IsKindOf $X.label("KindOf") $X.group("objects") $X.desc("True if the object is compatible with the named kind.") {
  object       @0  :Rtx.TextEval;
  kind         @1  :Text $X.label("is");
}
struct IsNotTrue $X.label("Not") $X.group("logic") $X.desc("Returns the opposite value.") {
  test         @0  :Rtx.BoolEval;
}
struct Join $X.label("Join") $X.group("strings") $X.desc("Returns multiple pieces of text as a single new piece of text.") {
  sep          @0  :Rtx.TextEval;
  parts        @1  :List(Rtx.TextEval) $X.label("parts");
}
struct KindOf $X.label("KindOf") $X.group("objects") $X.desc("Friendly name of the object's kind.") {
  object       @0  :Rtx.TextEval;
}
struct KindsOf $X.label("KindsOf") $X.group("objects") $X.desc("A list of compatible kinds.") {
  kind         @0  :Text;
}
struct LessOrEqual $X.label("AtMost") $X.group("comparison") $X.desc("The first value is larger than the second value.") {
}
struct LessThan $X.label("LessThan") $X.group("comparison") $X.desc("The first value is less than the second value.") {
}
struct Lines $X.label("Here") $X.group("literals") $X.desc("Specify one or more lines of text.") {
  lines        @0  :Text;
}
struct Make $X.label("Make") $X.group("misc") {
  name         @0  :Text $X.pool("kind");
  args         @1  :List(Argument) $X.label("args");
}
struct MakeLowercase $X.label("Lower") $X.group("format") $X.desc("Returns new text, with every letter turned into lowercase. For example, 'shout' from 'SHOUT'.") {
  text         @0  :Rtx.TextEval;
}
struct MakePlural $X.label("Pluralize") $X.group("format") $X.desc("Returns the plural form of a singular word. (ex. apples for apple. )") {
  text         @0  :Rtx.TextEval;
}
struct MakeReversed $X.label("Reverse") $X.group("format") $X.desc("Returns new text flipped back to front. For example, 'elppA' from 'Apple', or 'noon' from 'noon'.") {
  text         @0  :Rtx.TextEval;
}
struct MakeSentenceCase $X.label("Sentence") $X.group("format") $X.desc("Returns new text, start each sentence with a capital letter. For example, 'Empire Apple.' from 'Empire apple.'.") {
  text         @0  :Rtx.TextEval;
}
struct MakeSingular $X.label("Singularize") $X.group("format") $X.desc("Returns the singular form of a plural word. (ex. apple for apples )") {
  text         @0  :Rtx.TextEval;
}
struct MakeTitleCase $X.label("Title") $X.group("format") $X.desc("Returns new text, starting each word with a capital letter. For example, 'Empire Apple' from 'empire apple'.") {
  text         @0  :Rtx.TextEval;
}
struct MakeUppercase $X.label("Upper") $X.group("format") $X.desc("Returns new text, with every letter turned into uppercase. For example, 'APPLE' from 'apple'.") {
  text         @0  :Rtx.TextEval;
}
struct Matches $X.label("Matches") $X.group("matching") $X.desc("Determine whether the specified text is similar to the specified regular expression.") {
  text         @0  :Rtx.TextEval;
  pattern      @1  :Text $X.label("to");
}
struct NameOf $X.label("NameOf") $X.group("objects") $X.desc("Full name of the object.") {
  object       @0  :Rtx.TextEval;
}
struct Newline $X.label("Br") $X.group("printing") $X.desc("Start a new line.") {
}
struct Next $X.label("Next") $X.group("flow") $X.desc("In a repeating loop, try the next iteration of the loop.") {
}
struct NotEqualTo $X.label("OtherThan") $X.group("comparison") $X.desc("Two values don't match exactly.") {
}
struct NumList $X.label("Nums") $X.group("literals") $X.desc("Specify a list of multiple numbers.") {
  values       @0  :List(Float64);
}
struct NumValue $X.label("Num") $X.group("literals") $X.desc("Specify a particular number.") {
  num          @0  :Float64;
}
struct ObjectExists $X.label("Is") $X.group("objects") $X.desc("Returns whether there is a object of the specified name.") {
  object       @0  :Rtx.TextEval $X.label("valid");
}
struct Paragraph $X.label("P") $X.group("printing") $X.desc("Add a single blank line following some text.") {
}
struct PrintNum $X.label("Numeral") $X.group("printing") $X.desc("Writes a number using numerals, eg. '1'.") {
  num          @0  :Rtx.NumberEval;
}
struct PrintNumWord $X.label("Numeral") $X.group("printing") $X.desc("Writes a number in plain english: eg. 'one'") {
  num          @0  :Rtx.NumberEval $X.label("words");
}
struct ProductOf $X.label("Mul") $X.group("math") $X.desc("Multiply two numbers.") {
  a            @0  :Rtx.NumberEval;
  b            @1  :Rtx.NumberEval $X.label("by");
}
struct PutAtField $X.label("Put") $X.group("variables") $X.desc("Put a value into the field of an record or object") {
  into         @0  :IntoTargetFields;
  from         @1  :Rtx.Assignment $X.label("from");
  atField      @2  :Text $X.label("at");
}
struct QuotientOf $X.label("Div") $X.group("math") $X.desc("Divide one number by another.") {
  a            @0  :Rtx.NumberEval;
  b            @1  :Rtx.NumberEval $X.label("by");
}
struct RemainderOf $X.label("Mod") $X.group("math") $X.desc("Divide one number by another, and return the remainder.") {
  a            @0  :Rtx.NumberEval;
  b            @1  :Rtx.NumberEval $X.label("by");
}
struct Response $X.label("Response") $X.group("output") $X.desc("Generate text in a replaceable manner.") {
  name         @0  :Text;
  text         @1  :Rtx.TextEval $X.optional $X.label("text");
}
struct Row $X.label("Row") $X.group("printing") $X.desc("A single line as part of a group of lines.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct Rows $X.label("Rows") $X.group("printing") $X.desc("Group text into successive lines.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct Say $X.label("SayText") $X.group("printing") $X.desc("Print some bit of text to the player.") {
  text         @0  :Rtx.TextEval;
}
struct Send $X.label("Send") $X.group("events") $X.desc("Triggers a event, returns a true/false success value.") {
  event        @0  :Text $X.pool("event");
  path         @1  :Rtx.TextListEval $X.label("to");
  args         @2  :List(Argument) $X.label("args");
}
struct SetTrait $X.label("Put") $X.group("objects") $X.desc("Put an object into a particular state.") {
  object       @0  :Rtx.TextEval $X.label("obj");
  trait        @1  :Rtx.TextEval $X.label("trait");
}
struct ShuffleText $X.label("Shuffle") $X.group("output") $X.desc("When called multiple times returns its inputs at random.") {
  seq          @0  :Text $X.internal;
  parts        @1  :List(Rtx.TextEval);
}
struct Slash $X.label("SlashText") $X.group("printing") $X.desc("Separates words with left-leaning slashes '/'.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct Softline $X.label("Wbr") $X.group("printing") $X.desc("Start a new line ( if not already at a new line. )") {
}
struct Span $X.label("SpanText") $X.group("printing") $X.desc("Writes text with spaces between words.") {
  exe          @0  :List(Rtx.Execute) $X.label("do");
}
struct StoppingText $X.label("Stopping") $X.group("output") $X.desc("When called multiple times returns each of its inputs in turn, sticking to the last one.") {
  seq          @0  :Text $X.internal;
  parts        @1  :List(Rtx.TextEval);
}
struct SumOf $X.label("Inc") $X.group("math") $X.desc("Add two numbers.") {
  a            @0  :Rtx.NumberEval;
  b            @1  :Rtx.NumberEval $X.optional $X.label("by");
}
struct TextList $X.label("Txts") $X.group("literals") $X.desc("Specifies multiple string values.") {
  values       @0  :List(Text);
}
struct TextValue $X.label("Txt") $X.group("literals") $X.desc("Specify a small bit of text.") {
  text         @0  :Text;
}
struct TriggerCycle $X.label("Every") $X.group("comparison") {
}
struct TriggerOnce $X.label("At") $X.group("comparison") {
}
struct TriggerSwitch $X.label("After") $X.group("comparison") {
}
struct Var $X.label("Var") $X.group("variables") $X.desc("Return the value of the named variable.") {
  name         @0  :Text $X.pool("variable");
}
struct Variable $X.label("Var") $X.group("misc") {
  at           @0  :Reader.Pos $X.internal;
  str          @1  :Text $X.pool("variable");
}
struct While $X.label("Repeating") $X.group("flow") $X.desc("Keep running a series of actions while a condition is true.") {
  true         @0  :Rtx.BoolEval;
  exe          @1  :List(Rtx.Execute) $X.label("do");
}
