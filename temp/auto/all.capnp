@0xad22bd0042f92910;
using Go = import "/go.capnp";
using  X = import "./options.capnp";

$Go.package("auto");
$Go.import("git.sr.ht/~ionous/pb/auto");

struct Pos {
	source @0 :Text; 
	offset @1 :Text;
}

struct Assignment { eval @0:AnyPointer; }
struct AssignmentImpl $X.desc("Helper for setting variables.") {
  union {
	bool                           @0   :FromBool $X.label("Bool:");
	getFrom                        @1   :GetAtField $X.label("Get:from:");
	num                            @2   :FromNum $X.label("Num:");
	nums                           @3   :FromNumbers $X.label("Nums:");
	rec                            @4   :FromRecord $X.label("Rec:");
	recs                           @5   :FromRecords $X.label("Recs:");
	renderArgs                     @6   :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @7   :RenderRef $X.label("RenderRef:flags:");
	txt                            @8   :FromText $X.label("Txt:");
	txts                           @9   :FromTexts $X.label("Txts:");
	var                            @10  :Var $X.label("Var:");
  }
}

struct BoolEval { eval @0:AnyPointer; }
struct BoolEvalImpl $X.desc("Statements which return true/false values.") {
  union {
	allOf                          @0   :AllTrue $X.label("AllOf:");
	always                         @1   :Always $X.label("Always");
	anyOf                          @2   :AnyTrue $X.label("AnyOf:");
	bool                           @3   :Bool $X.label("Bool:");
	cmpIsNum                       @4   :CompareNum $X.label("Cmp:is:num:");
	cmpIsTxt                       @5   :CompareText $X.label("Cmp:is:txt:");
	containsPart                   @6   :Includes $X.label("Contains:part:");
	countOfTrigger                 @7   :CountOf $X.label("CountOf:trigger:");
	determineArgs                  @8   :Determine $X.label("Determine:args:");
	during                         @9   :During $X.label("During:");
	findList                       @10  :Find $X.label("Find:list:");
	getObjTrait                    @11  :HasTrait $X.label("Get obj:trait:");
	getFrom                        @12  :GetAtField $X.label("Get:from:");
	hasDominion                    @13  :HasDominion $X.label("HasDominion:");
	isEmpty                        @14  :IsEmpty $X.label("Is empty:");
	isValid                        @15  :ObjectExists $X.label("Is valid:");
	kindOfIs                       @16  :IsKindOf $X.label("KindOf:is:");
	kindOfIsExactly                @17  :IsExactKindOf $X.label("KindOf:isExactly:");
	matchesTo                      @18  :Matches $X.label("Matches:to:");
	not                            @19  :IsNotTrue $X.label("Not:");
	renderArgs                     @20  :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @21  :RenderRef $X.label("RenderRef:flags:");
	sendToArgs                     @22  :Send $X.label("Send:to:args:");
	var                            @23  :Var $X.label("Var:");
  }
}

struct Brancher { eval @0:AnyPointer; }
struct BrancherImpl $X.desc("Helper for choose action.") {
  union {
	elseDoDo                       @0   :ChooseNothingElse $X.label("ElseDo do:");
	elseIfDo                       @1   :ChooseMore $X.label("ElseIf:do:");
	elseIfDoElse                   @2   :ChooseMore $X.label("ElseIf:do:else:");
	elseIfFromAndDo                @3   :ChooseMoreValue $X.label("ElseIf:from:and:do:");
	elseIfFromAndDoElse            @4   :ChooseMoreValue $X.label("ElseIf:from:and:do:else:");
	ifDo                           @5   :ChooseAction $X.label("If:do:");
	ifDoElse                       @6   :ChooseAction $X.label("If:do:else:");
	ifFromAndDo                    @7   :ChooseValue $X.label("If:from:and:do:");
	ifFromAndDoElse                @8   :ChooseValue $X.label("If:from:and:do:else:");
  }
}

struct Comparator { eval @0:AnyPointer; }
struct ComparatorImpl $X.desc("Helper for comparing values.") {
  union {
	atLeast                        @0   :GreaterOrEqual $X.label("AtLeast");
	atMost                         @1   :LessOrEqual $X.label("AtMost");
	equalTo                        @2   :EqualTo $X.label("EqualTo");
	greaterThan                    @3   :GreaterThan $X.label("GreaterThan");
	lessThan                       @4   :LessThan $X.label("LessThan");
	otherThan                      @5   :NotEqualTo $X.label("OtherThan");
  }
}

struct Execute { eval @0:AnyPointer; }
struct ExecuteImpl $X.desc("Run a series of statements.") {
  union {
	actDo                          @0   :Activity $X.label("Act do:");
	br                             @1   :Newline $X.label("Br");
	break                          @2   :Break $X.label("Break");
	determineArgs                  @3   :Determine $X.label("Determine:args:");
	doNothing                      @4   :DoNothing $X.label("DoNothing");
	doNothingWhy                   @5   :DoNothing $X.label("DoNothing why:");
	eraseAtEdge                    @6   :EraseEdge $X.label("Erase:atEdge:");
	eraseFromAtIndex               @7   :EraseIndex $X.label("Erase:from:atIndex:");
	erasingAtEdgeAsDo              @8   :ErasingEdge $X.label("Erasing:atEdge:as:do:");
	erasingAtEdgeAsDoElse          @9   :ErasingEdge $X.label("Erasing:atEdge:as:do:else:");
	erasingFromAtIndexAsDo         @10  :Erasing $X.label("Erasing:from:atIndex:as:do:");
	ifDo                           @11  :ChooseAction $X.label("If:do:");
	ifDoElse                       @12  :ChooseAction $X.label("If:do:else:");
	ifFromAndDo                    @13  :ChooseValue $X.label("If:from:and:do:");
	ifFromAndDoElse                @14  :ChooseValue $X.label("If:from:and:do:else:");
	letBe                          @15  :Assign $X.label("Let:be:");
	listReduceIntoFromListUsing    @16  :Reduce $X.label("ListReduce into:fromList:using:");
	listReverse                    @17  :ReverseList $X.label("ListReverse:");
	listSetIndexFrom               @18  :Set $X.label("ListSet:index:from:");
	logValue                       @19  :Log $X.label("Log:value:");
	mapFromListUsing               @20  :Map $X.label("Map:fromList:using:");
	next                           @21  :Next $X.label("Next");
	p                              @22  :Paragraph $X.label("P");
	putObjTrait                    @23  :SetTrait $X.label("Put obj:trait:");
	putFromAt                      @24  :PutAtField $X.label("Put:from:at:");
	putIntoAtEdge                  @25  :PutEdge $X.label("Put:into:atEdge:");
	putIntoAtIndex                 @26  :PutIndex $X.label("Put:into:atIndex:");
	relateToRel                    @27  :Relate $X.label("Relate:to:rel:");
	renderArgs                     @28  :RenderPattern $X.label("Render:args:");
	repeatingAcrossAsDo            @29  :Each $X.label("Repeating across:as:do:");
	repeatingAcrossAsDoElse        @30  :Each $X.label("Repeating across:as:do:else:");
	repeatingDo                    @31  :While $X.label("Repeating:do:");
	sayText                        @32  :Say $X.label("SayText:");
	sendToArgs                     @33  :Send $X.label("Send:to:args:");
	sortByOrder                    @34  :SortNumbers $X.label("Sort:by:order:");
	sortByOrderCase                @35  :SortText $X.label("Sort:by:order:case:");
	sortUsing                      @36  :SortRecords $X.label("Sort:using:");
	spliceStartRemoveInsert        @37  :Splice $X.label("Splice:start:remove:insert:");
	wbr                            @38  :Softline $X.label("Wbr");
  }
}

struct FromSourceFields { eval @0:AnyPointer; }
struct FromSourceFieldsImpl $X.desc("Helper for getting fields.") {
  union {
	obj                            @0   :FromObj $X.label("Obj:");
	rec                            @1   :FromRec $X.label("Rec:");
	renderField                    @2   :RenderField $X.label("RenderField:");
	var                            @3   :FromVar $X.label("Var:");
  }
}

struct GrammarMaker { eval @0:AnyPointer; }
struct GrammarMakerImpl $X.desc("Helper for defining parser grammars.") {
  union {
	aliasAsNoun                    @0   :Alias $X.label("Alias:asNoun:");
	directiveScans                 @1   :Directive $X.label("Directive:scans:");
  }
}

struct IntoTargetFields { eval @0:AnyPointer; }
struct IntoTargetFieldsImpl $X.desc("Helper for setting fields.") {
  union {
	obj                            @0   :IntoObj $X.label("Obj:");
	var                            @1   :IntoVar $X.label("Var:");
  }
}

struct ListIterator { eval @0:AnyPointer; }
struct ListIteratorImpl $X.desc("Helper for accessing lists.") {
  union {
	num                            @0   :AsNum $X.label("Num:");
	rec                            @1   :AsRec $X.label("Rec:");
	txt                            @2   :AsTxt $X.label("Txt:");
  }
}

struct ListSource { eval @0:AnyPointer; }
struct ListSourceImpl $X.desc("Helper for accessing lists.") {
  union {
	nums                           @0   :FromNumList $X.label("Nums:");
	recs                           @1   :FromRecList $X.label("Recs:");
	txts                           @2   :FromTxtList $X.label("Txts:");
  }
}

struct ListTarget { eval @0:AnyPointer; }
struct ListTargetImpl $X.desc("Helper for accessing lists.") {
  union {
	nums                           @0   :IntoNumList $X.label("Nums:");
	recs                           @1   :IntoRecList $X.label("Recs:");
	txts                           @2   :IntoTxtList $X.label("Txts:");
  }
}

struct NumListEval { eval @0:AnyPointer; }
struct NumListEvalImpl $X.desc("Statements which return a list of numbers.") {
  union {
	determineArgs                  @0   :Determine $X.label("Determine:args:");
	getFrom                        @1   :GetAtField $X.label("Get:from:");
	nums                           @2   :Numbers $X.label("Nums:");
	range                          @3   :Range $X.label("Range:");
	rangeByStep                    @4   :Range $X.label("Range:byStep:");
	rangeFrom                      @5   :Range $X.label("Range:from:");
	rangeFromByStep                @6   :Range $X.label("Range:from:byStep:");
	renderArgs                     @7   :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @8   :RenderRef $X.label("RenderRef:flags:");
	slice                          @9   :Slice $X.label("Slice:");
	sliceEnd                       @10  :Slice $X.label("Slice:end:");
	sliceStart                     @11  :Slice $X.label("Slice:start:");
	sliceStartEnd                  @12  :Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @13  :Splice $X.label("Splice:start:remove:insert:");
	var                            @14  :Var $X.label("Var:");
  }
}

struct NumberEval { eval @0:AnyPointer; }
struct NumberEvalImpl $X.desc("Statements which return a number.") {
  union {
	dec                            @0   :DiffOf $X.label("Dec:");
	decBy                          @1   :DiffOf $X.label("Dec:by:");
	determineArgs                  @2   :Determine $X.label("Determine:args:");
	divBy                          @3   :QuotientOf $X.label("Div:by:");
	during                         @4   :During $X.label("During:");
	findList                       @5   :Find $X.label("Find:list:");
	getFrom                        @6   :GetAtField $X.label("Get:from:");
	getIndex                       @7   :At $X.label("Get:index:");
	inc                            @8   :SumOf $X.label("Inc:");
	incBy                          @9   :SumOf $X.label("Inc:by:");
	len                            @10  :Len $X.label("Len:");
	modBy                          @11  :RemainderOf $X.label("Mod:by:");
	mulBy                          @12  :ProductOf $X.label("Mul:by:");
	num                            @13  :Number $X.label("Num:");
	numIfElse                      @14  :ChooseNum $X.label("Num:if:else:");
	renderArgs                     @15  :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @16  :RenderRef $X.label("RenderRef:flags:");
	var                            @17  :Var $X.label("Var:");
  }
}

struct RecordEval { eval @0:AnyPointer; }
struct RecordEvalImpl $X.desc("Statements which return a record.") {
  union {
	determineArgs                  @0   :Determine $X.label("Determine:args:");
	getFrom                        @1   :GetAtField $X.label("Get:from:");
	getIndex                       @2   :At $X.label("Get:index:");
	makeArgs                       @3   :Make $X.label("Make:args:");
	renderArgs                     @4   :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @5   :RenderRef $X.label("RenderRef:flags:");
	var                            @6   :Var $X.label("Var:");
  }
}

struct RecordListEval { eval @0:AnyPointer; }
struct RecordListEvalImpl $X.desc("Statements which return a list of records.") {
  union {
	determineArgs                  @0   :Determine $X.label("Determine:args:");
	getFrom                        @1   :GetAtField $X.label("Get:from:");
	renderArgs                     @2   :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @3   :RenderRef $X.label("RenderRef:flags:");
	slice                          @4   :Slice $X.label("Slice:");
	sliceEnd                       @5   :Slice $X.label("Slice:end:");
	sliceStart                     @6   :Slice $X.label("Slice:start:");
	sliceStartEnd                  @7   :Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @8   :Splice $X.label("Splice:start:remove:insert:");
	var                            @9   :Var $X.label("Var:");
  }
}

struct ScannerMaker { eval @0:AnyPointer; }
struct ScannerMakerImpl $X.desc("Helper for defining input scanners.") {
  union {
	allOf                          @0   :AllOf $X.label("AllOf:");
	anyOf                          @1   :AnyOf $X.label("AnyOf:");
	as                             @2   :Action $X.label("As:");
	noun                           @3   :Noun $X.label("Noun:");
	retarget                       @4   :Retarget $X.label("Retarget:");
	reverse                        @5   :Reverse $X.label("Reverse:");
	self                           @6   :Self $X.label("Self:");
	words                          @7   :Words $X.label("Words:");
  }
}

struct TextEval { eval @0:AnyPointer; }
struct TextEvalImpl $X.desc("Statements which return text.") {
  union {
	bracketTextDo                  @0   :Bracket $X.label("BracketText do:");
	bufferTextDo                   @1   :Buffer $X.label("BufferText do:");
	capitalize                     @2   :Capitalize $X.label("Capitalize:");
	commaTextDo                    @3   :Commas $X.label("CommaText do:");
	cycle                          @4   :CycleText $X.label("Cycle:");
	determineArgs                  @5   :Determine $X.label("Determine:args:");
	getFrom                        @6   :GetAtField $X.label("Get:from:");
	getIndex                       @7   :At $X.label("Get:index:");
	idOf                           @8   :IdOf $X.label("IdOf:");
	joinParts                      @9   :Join $X.label("Join:parts:");
	kindOf                         @10  :KindOf $X.label("KindOf:");
	lower                          @11  :MakeLowercase $X.label("Lower:");
	nameOf                         @12  :NameOf $X.label("NameOf:");
	numeralWords                   @13  :PrintNumWord $X.label("Numeral words:");
	numeral                        @14  :PrintNum $X.label("Numeral:");
	pluralize                      @15  :MakePlural $X.label("Pluralize:");
	reciprocalRelObject            @16  :ReciprocalOf $X.label("Reciprocal rel:object:");
	relativeRelObject              @17  :RelativeOf $X.label("Relative rel:object:");
	renderArgs                     @18  :RenderPattern $X.label("Render:args:");
	renderName                     @19  :RenderName $X.label("RenderName:");
	renderRefFlags                 @20  :RenderRef $X.label("RenderRef:flags:");
	renderTemplate                 @21  :RenderTemplate $X.label("RenderTemplate");
	response                       @22  :Response $X.label("Response:");
	responseText                   @23  :Response $X.label("Response:text:");
	reverse                        @24  :MakeReversed $X.label("Reverse:");
	rowDo                          @25  :Row $X.label("Row do:");
	rowsDo                         @26  :Rows $X.label("Rows do:");
	sentence                       @27  :MakeSentenceCase $X.label("Sentence:");
	shuffle                        @28  :ShuffleText $X.label("Shuffle:");
	singularize                    @29  :MakeSingular $X.label("Singularize:");
	slashTextDo                    @30  :Slash $X.label("SlashText do:");
	spanTextDo                     @31  :Span $X.label("SpanText do:");
	stopping                       @32  :StoppingText $X.label("Stopping:");
	title                          @33  :MakeTitleCase $X.label("Title:");
	txt                            @34  :Text $X.label("Txt:");
	txtIfElse                      @35  :ChooseText $X.label("Txt:if:else:");
	upper                          @36  :MakeUppercase $X.label("Upper:");
	var                            @37  :Var $X.label("Var:");
  }
}

struct TextListEval { eval @0:AnyPointer; }
struct TextListEvalImpl $X.desc("Statements which return a list of text.") {
  union {
	determineArgs                  @0   :Determine $X.label("Determine:args:");
	getFrom                        @1   :GetAtField $X.label("Get:from:");
	kindsOf                        @2   :KindsOf $X.label("KindsOf:");
	reciprocalsRelObject           @3   :ReciprocalsOf $X.label("Reciprocals rel:object:");
	relativesRelObject             @4   :RelativesOf $X.label("Relatives rel:object:");
	renderArgs                     @5   :RenderPattern $X.label("Render:args:");
	renderRefFlags                 @6   :RenderRef $X.label("RenderRef:flags:");
	slice                          @7   :Slice $X.label("Slice:");
	sliceEnd                       @8   :Slice $X.label("Slice:end:");
	sliceStart                     @9   :Slice $X.label("Slice:start:");
	sliceStartEnd                  @10  :Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @11  :Splice $X.label("Splice:start:remove:insert:");
	txts                           @12  :Texts $X.label("Txts:");
	var                            @13  :Var $X.label("Var:");
  }
}

struct Trigger { eval @0:AnyPointer; }
struct TriggerImpl $X.desc("Helper for counting values.") {
  union {
	after                          @0   :TriggerSwitch $X.label("After");
	at                             @1   :TriggerOnce $X.label("At");
	every                          @2   :TriggerCycle $X.label("Every");
  }
}


struct Action $X.label("as") $X.group("grammar")  {
  action       @0  :Text;
}

struct Activity $X.label("act") $X.group("hidden")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Alias $X.label("alias") $X.group("grammar")  {
  names        @0  :List(Text);
  asNoun       @1  :Text $X.label("as_noun");
}

struct AllOf $X.label("allOf") $X.group("grammar")  {
  series       @0  :List(ScannerMaker);
}

struct AllTrue $X.label("allOf") $X.group("logic")  $X.desc("Returns true if all of the evaluations are true.")  {
  test         @0  :List(BoolEval);
}

struct Always $X.label("always") $X.group("logic")  $X.desc("Returns true always.")  {
}

struct AnyOf $X.label("anyOf") $X.group("grammar")  {
  options      @0  :List(ScannerMaker);
}

struct AnyTrue $X.label("anyOf") $X.group("logic")  $X.desc("Returns true if any of the evaluations are true.")  {
  test         @0  :List(BoolEval);
}

struct Argument $X.label("arg") $X.group("patterns")  {
  name         @0  :Text;
  from         @1  :Assignment $X.label("from");
}

struct Arguments $X.label("arguments") $X.group("patterns")  {
  args         @0  :List(Argument);
}

struct AsNum $X.label("num") $X.group("misc")  $X.desc("Define the name of a number variable.")  {
  str          @0  :Text $X.pool("variable");
}

struct AsRec $X.label("rec") $X.group("misc")  $X.desc("Define the name of a record variable.")  {
  str          @0  :Text $X.pool("variable");
}

struct AsTxt $X.label("txt") $X.group("misc")  $X.desc("Define the name of a text variable.")  {
  str          @0  :Text $X.pool("variable");
}

struct Assign $X.label("let") $X.group("variables")  $X.desc("Assigns a variable to a value.")  {
  str          @0  :Text $X.pool("variable");
  from         @1  :Assignment $X.label("be");
}

struct At $X.label("get") $X.group("list")  $X.desc("Get a value from a list. The first element is is index 1.")  {
  list         @0  :Assignment;
  index        @1  :NumberEval $X.label("index");
}

struct Bool $X.label("bool") $X.group("literals")  $X.desc("Specify an explicit true or false value.")  {
  bool         @0  :Bool;
}

struct Bracket $X.label("bracketText") $X.group("printing")  $X.desc("Sandwiches text printed during a block and puts them inside parenthesis '()'.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Break $X.label("break") $X.group("flow")  $X.desc("In a repeating loop, exit the loop.")  {
}

struct Buffer $X.label("bufferText") $X.group("printing")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Capitalize $X.label("capitalize") $X.group("format")  $X.desc("Returns new text, with the first letter turned into uppercase.")  {
  text         @0  :TextEval;
}

struct ChooseAction $X.label("if") $X.group("misc")  $X.desc("An if statement.")  {
  if           @0  :BoolEval;
  exe          @1  :List(Execute) $X.label("do");
  else         @2  :Brancher $X.optional $X.label("else");
}

struct ChooseMore $X.label("elseIf") $X.group("misc")  {
  if           @0  :BoolEval;
  exe          @1  :List(Execute) $X.label("do");
  else         @2  :Brancher $X.optional $X.label("else");
}

struct ChooseMoreValue $X.label("elseIf") $X.group("misc")  {
  assign       @0  :Text;
  from         @1  :Assignment $X.label("from");
  filter       @2  :BoolEval $X.label("and");
  exe          @3  :List(Execute) $X.label("do");
  else         @4  :Brancher $X.optional $X.label("else");
}

struct ChooseNothingElse $X.label("elseDo") $X.group("misc")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct ChooseNum $X.label("num") $X.group("math")  $X.desc("Pick one of two numbers based on a boolean test.")  {
  true         @0  :NumberEval;
  if           @1  :BoolEval $X.label("if");
  false        @2  :NumberEval $X.label("else");
}

struct ChooseText $X.label("txt") $X.group("format")  $X.desc("Pick one of two strings based on a boolean test.")  {
  true         @0  :TextEval;
  if           @1  :BoolEval $X.label("if");
  false        @2  :TextEval $X.label("else");
}

struct ChooseValue $X.label("if") $X.group("misc")  $X.desc("An if statement with local assignment.")  {
  assign       @0  :Text;
  from         @1  :Assignment $X.label("from");
  filter       @2  :BoolEval $X.label("and");
  exe          @3  :List(Execute) $X.label("do");
  else         @4  :Brancher $X.optional $X.label("else");
}

struct Commas $X.label("commaText") $X.group("printing")  $X.desc("Separates words with commas, and 'and'.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct CompareNum $X.label("cmp") $X.group("logic")  $X.desc("True if eq,ne,gt,lt,ge,le two numbers.")  {
  a            @0  :NumberEval;
  is           @1  :Comparator $X.label("is");
  b            @2  :NumberEval $X.label("num");
}

struct CompareText $X.label("cmp") $X.group("logic")  $X.desc("True if eq,ne,gt,lt,ge,le two strings ( lexical. )")  {
  a            @0  :TextEval;
  is           @1  :Comparator $X.label("is");
  b            @2  :TextEval $X.label("txt");
}

struct CountOf $X.label("countOf") $X.group("logic")  $X.desc("A guard which returns true based on a counter. Counters start at zero and are incremented every time the guard gets checked.")  {
  at           @0  :Pos $X.internal;
  num          @1  :NumberEval;
  trigger      @2  :Trigger $X.label("trigger");
}

struct CycleText $X.label("cycle") $X.group("output")  $X.desc("When called multiple times, returns each of its inputs in turn.")  {
  seq          @0  :Text $X.internal;
  parts        @1  :List(TextEval);
}

struct Determine $X.label("determine") $X.group("patterns")  $X.desc("Runs a pattern, and potentially returns a value.")  {
  pattern      @0  :Text $X.pool("pattern");
  args         @1  :List(Argument) $X.label("args");
}

struct DiffOf $X.label("dec") $X.group("math")  $X.desc("Subtract two numbers.")  {
  a            @0  :NumberEval;
  b            @1  :NumberEval $X.optional $X.label("by");
}

struct Directive $X.label("directive") $X.group("grammar")  {
  lede         @0  :List(Text);
  scans        @1  :List(ScannerMaker) $X.label("scans");
}

struct DoNothing $X.label("doNothing") $X.group("flow")  $X.desc("Statement which does nothing.")  {
  reason       @0  :Text $X.optional $X.label("why");
}

struct During $X.label("during") $X.group("patterns")  $X.desc("Decide whether a pattern is running.")  {
  pattern      @0  :Text $X.pool("pattern");
}

struct Each $X.label("repeating") $X.group("list")  $X.desc("Loops over the elements in the passed list, or runs the 'else' activity if empty.")  {
  list         @0  :Assignment $X.label("across");
  as           @1  :ListIterator $X.label("as");
  exe          @2  :List(Execute) $X.label("do");
  else         @3  :Brancher $X.optional $X.label("else");
}

struct EqualTo $X.label("equalTo") $X.group("comparison")  $X.desc("Two values exactly match.")  {
}

struct EraseEdge $X.label("erase") $X.group("misc")  $X.desc("Remove one or more values from a list")  {
  from         @0  :ListSource;
  atEdge       @1  :Bool $X.label("at_edge");
}

struct EraseIndex $X.label("erase") $X.group("misc")  $X.desc("Remove one or more values from a list")  {
  count        @0  :NumberEval;
  from         @1  :ListSource $X.label("from");
  atIndex      @2  :NumberEval $X.label("at_index");
}

struct Erasing $X.label("erasing") $X.group("list")  $X.desc("Erase elements from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.")  {
  count        @0  :NumberEval;
  from         @1  :ListSource $X.label("from");
  atIndex      @2  :NumberEval $X.label("at_index");
  as           @3  :Text $X.label("as");
  exe          @4  :List(Execute) $X.label("do");
}

struct ErasingEdge $X.label("erasing") $X.group("list")  $X.desc("Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.")  {
  from         @0  :ListSource;
  atEdge       @1  :Bool $X.label("at_edge");
  as           @2  :Text $X.label("as");
  exe          @3  :List(Execute) $X.label("do");
  else         @4  :Brancher $X.optional $X.label("else");
}

struct Find $X.label("find") $X.group("list")  $X.desc("Search a list for a specific value.")  {
  value        @0  :Assignment;
  list         @1  :Assignment $X.label("list");
}

struct FromBool $X.label("bool") $X.group("variables")  $X.desc("Assigns the calculated boolean value.")  {
  val          @0  :BoolEval;
}

struct FromNum $X.label("num") $X.group("variables")  $X.desc("Assigns the calculated number.")  {
  val          @0  :NumberEval;
}

struct FromNumList $X.label("nums") $X.group("misc")  $X.desc("Uses a list of numbers")  {
  str          @0  :Text $X.pool("variable");
}

struct FromNumbers $X.label("nums") $X.group("variables")  $X.desc("Assigns the calculated numbers.")  {
  vals         @0  :NumListEval;
}

struct FromObj $X.label("obj") $X.group("misc")  $X.desc("Targets an object with a computed name.")  {
  object       @0  :TextEval;
}

struct FromRec $X.label("rec") $X.group("misc")  $X.desc("Targets a record stored in a record.")  {
  rec          @0  :RecordEval;
}

struct FromRecList $X.label("recs") $X.group("misc")  $X.desc("Uses a list of records")  {
  str          @0  :Text $X.pool("variable");
}

struct FromRecord $X.label("rec") $X.group("variables")  $X.desc("Assigns the calculated record.")  {
  val          @0  :RecordEval;
}

struct FromRecords $X.label("recs") $X.group("variables")  $X.desc("Assigns the calculated records.")  {
  vals         @0  :RecordListEval;
}

struct FromText $X.label("txt") $X.group("variables")  $X.desc("Assigns the calculated piece of text.")  {
  val          @0  :TextEval;
}

struct FromTexts $X.label("txts") $X.group("variables")  $X.desc("Assigns the calculated texts.")  {
  vals         @0  :TextListEval;
}

struct FromTxtList $X.label("txts") $X.group("misc")  $X.desc("Uses a list of text")  {
  str          @0  :Text $X.pool("variable");
}

struct FromVar $X.label("var") $X.group("misc")  $X.desc("Targets a record stored in a variable.")  {
  str          @0  :Text $X.pool("variable");
}

struct Gather $X.label("gather") $X.group("list")  $X.desc("Transform the values from a list. The named pattern gets called once for each value in the list. It get called with two parameters: 'in' as each value from the list, and 'out' as the var passed to the gather.")  {
  str          @0  :Text $X.pool("variable");
  from         @1  :ListSource $X.label("from");
  using        @2  :Text $X.label("using") $X.pool("pattern");
}

struct GetAtField $X.label("get") $X.group("variables")  $X.desc("Get a value from a record.")  {
  field        @0  :Text;
  from         @1  :FromSourceFields $X.label("from");
}

struct GrammarDecl $X.label("grammar") $X.group("grammar")  $X.desc("Read what the player types and turn it into actions.")  {
  grammar      @0  :GrammarMaker;
}

struct GreaterOrEqual $X.label("atLeast") $X.group("comparison")  $X.desc("The first value is larger than the second value.")  {
}

struct GreaterThan $X.label("greaterThan") $X.group("comparison")  $X.desc("The first value is larger than the second value.")  {
}

struct HasDominion $X.label("hasDominion") $X.group("logic")  {
  name         @0  :Text $X.pool("domain");
}

struct HasTrait $X.label("get") $X.group("objects")  $X.desc("Return true if the object is currently in the requested state.")  {
  object       @0  :TextEval $X.label("obj");
  trait        @1  :TextEval $X.label("trait");
}

struct IdOf $X.label("idOf") $X.group("objects")  $X.desc("A unique object identifier.")  {
  object       @0  :TextEval;
}

struct Includes $X.label("contains") $X.group("strings")  $X.desc("True if text contains text.")  {
  text         @0  :TextEval;
  part         @1  :TextEval $X.label("part");
}

struct IntoNumList $X.label("nums") $X.group("misc")  $X.desc("Targets a list of numbers")  {
  str          @0  :Text $X.pool("variable");
}

struct IntoObj $X.label("obj") $X.group("misc")  $X.desc("Targets an object with a computed name.")  {
  object       @0  :TextEval;
}

struct IntoRecList $X.label("recs") $X.group("misc")  $X.desc("Targets a list of records")  {
  str          @0  :Text $X.pool("variable");
}

struct IntoTxtList $X.label("txts") $X.group("misc")  $X.desc("Targets a list of text")  {
  str          @0  :Text $X.pool("variable");
}

struct IntoVar $X.label("var") $X.group("misc")  $X.desc("Targets an object or record stored in a variable")  {
  str          @0  :Text $X.pool("variable");
}

struct IsEmpty $X.label("is") $X.group("strings")  $X.desc("True if the text is empty.")  {
  text         @0  :TextEval $X.label("empty");
}

struct IsExactKindOf $X.label("kindOf") $X.group("objects")  $X.desc("True if the object is exactly the named kind.")  {
  object       @0  :TextEval;
  kind         @1  :Text $X.label("is_exactly");
}

struct IsKindOf $X.label("kindOf") $X.group("objects")  $X.desc("True if the object is compatible with the named kind.")  {
  object       @0  :TextEval;
  kind         @1  :Text $X.label("is");
}

struct IsNotTrue $X.label("not") $X.group("logic")  $X.desc("Returns the opposite value.")  {
  test         @0  :BoolEval;
}

struct Join $X.label("join") $X.group("strings")  $X.desc("Returns multiple pieces of text as a single new piece of text.")  {
  sep          @0  :TextEval;
  parts        @1  :List(TextEval) $X.label("parts");
}

struct KindOf $X.label("kindOf") $X.group("objects")  $X.desc("Friendly name of the object's kind.")  {
  object       @0  :TextEval;
}

struct KindsOf $X.label("kindsOf") $X.group("objects")  $X.desc("A list of compatible kinds.")  {
  kind         @0  :Text;
}

struct Len $X.label("len") $X.group("list")  $X.desc("Determines the number of values in a list.")  {
  list         @0  :Assignment;
}

struct LessOrEqual $X.label("atMost") $X.group("comparison")  $X.desc("The first value is larger than the second value.")  {
}

struct LessThan $X.label("lessThan") $X.group("comparison")  $X.desc("The first value is less than the second value.")  {
}

struct Lines $X.label("here") $X.group("literals")  $X.desc("Specify one or more lines of text.")  {
  lines        @0  :Text;
}

struct Log $X.label("log") $X.group("debug")  $X.desc("Debug log")  {
  level        @0  :Int32;
  value        @1  :Assignment $X.label("value");
}

struct Make $X.label("make") $X.group("misc")  {
  name         @0  :Text $X.pool("kind");
  args         @1  :List(Argument) $X.label("args");
}

struct MakeLowercase $X.label("lower") $X.group("format")  $X.desc("Returns new text, with every letter turned into lowercase. For example, 'shout' from 'SHOUT'.")  {
  text         @0  :TextEval;
}

struct MakePlural $X.label("pluralize") $X.group("format")  $X.desc("Returns the plural form of a singular word. (ex. apples for apple. )")  {
  text         @0  :TextEval;
}

struct MakeReversed $X.label("reverse") $X.group("format")  $X.desc("Returns new text flipped back to front. For example, 'elppA' from 'Apple', or 'noon' from 'noon'.")  {
  text         @0  :TextEval;
}

struct MakeSentenceCase $X.label("sentence") $X.group("format")  $X.desc("Returns new text, start each sentence with a capital letter. For example, 'Empire Apple.' from 'Empire apple.'.")  {
  text         @0  :TextEval;
}

struct MakeSingular $X.label("singularize") $X.group("format")  $X.desc("Returns the singular form of a plural word. (ex. apple for apples )")  {
  text         @0  :TextEval;
}

struct MakeTitleCase $X.label("title") $X.group("format")  $X.desc("Returns new text, starting each word with a capital letter. For example, 'Empire Apple' from 'empire apple'.")  {
  text         @0  :TextEval;
}

struct MakeUppercase $X.label("upper") $X.group("format")  $X.desc("Returns new text, with every letter turned into uppercase. For example, 'APPLE' from 'apple'.")  {
  text         @0  :TextEval;
}

struct Map $X.label("map") $X.group("list")  $X.desc("Transform the values from one list and place the results in another list. The designated pattern is called with each value from the 'from list', one value at a time.")  {
  toList       @0  :Text;
  fromList     @1  :Assignment $X.label("from_list");
  usingPattern @2  :Text $X.label("using") $X.pool("pattern");
}

struct Matches $X.label("matches") $X.group("matching")  $X.desc("Determine whether the specified text is similar to the specified regular expression.")  {
  text         @0  :TextEval;
  pattern      @1  :Text $X.label("to");
}

struct NameOf $X.label("nameOf") $X.group("objects")  $X.desc("Full name of the object.")  {
  object       @0  :TextEval;
}

struct Newline $X.label("br") $X.group("printing")  $X.desc("Start a new line.")  {
}

struct Next $X.label("next") $X.group("flow")  $X.desc("In a repeating loop, try the next iteration of the loop.")  {
}

struct NotEqualTo $X.label("otherThan") $X.group("comparison")  $X.desc("Two values don't match exactly.")  {
}

struct Noun $X.label("noun") $X.group("grammar")  {
  kind         @0  :Text;
}

struct Number $X.label("num") $X.group("literals")  $X.desc("Specify a particular number.")  {
  num          @0  :Float64;
}

struct Numbers $X.label("nums") $X.group("literals")  $X.desc("Specify a list of multiple numbers.")  {
  values       @0  :List(Float64);
}

struct ObjectExists $X.label("is") $X.group("objects")  $X.desc("Returns whether there is a object of the specified name.")  {
  object       @0  :TextEval $X.label("valid");
}

struct Paragraph $X.label("p") $X.group("printing")  $X.desc("Add a single blank line following some text.")  {
}

struct PrintNum $X.label("numeral") $X.group("printing")  $X.desc("Writes a number using numerals, eg. '1'.")  {
  num          @0  :NumberEval;
}

struct PrintNumWord $X.label("numeral") $X.group("printing")  $X.desc("Writes a number in plain english: eg. 'one'")  {
  num          @0  :NumberEval $X.label("words");
}

struct ProductOf $X.label("mul") $X.group("math")  $X.desc("Multiply two numbers.")  {
  a            @0  :NumberEval;
  b            @1  :NumberEval $X.label("by");
}

struct PutAtField $X.label("put") $X.group("variables")  $X.desc("Put a value into the field of an record or object")  {
  into         @0  :IntoTargetFields;
  from         @1  :Assignment $X.label("from");
  atField      @2  :Text $X.label("at");
}

struct PutEdge $X.label("put") $X.group("misc")  $X.desc("Add a value to a list")  {
  from         @0  :Assignment;
  into         @1  :ListTarget $X.label("into");
  atEdge       @2  :Bool $X.label("at_edge");
}

struct PutIndex $X.label("put") $X.group("misc")  $X.desc("Replace one value in a list with another")  {
  from         @0  :Assignment;
  into         @1  :ListTarget $X.label("into");
  atIndex      @2  :NumberEval $X.label("at_index");
}

struct QuotientOf $X.label("div") $X.group("math")  $X.desc("Divide one number by another.")  {
  a            @0  :NumberEval;
  b            @1  :NumberEval $X.label("by");
}

struct Range $X.label("range") $X.group("flow")  $X.desc("Generates a series of numbers.")  {
  to           @0  :NumberEval;
  from         @1  :NumberEval $X.optional $X.label("from");
  byStep       @2  :NumberEval $X.optional $X.label("by_step");
}

struct ReciprocalOf $X.label("reciprocal") $X.group("relations")  $X.desc("Returns the implied relative of a noun (ex. the source in a one-to-many relation.)")  {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :TextEval $X.label("object");
}

struct ReciprocalsOf $X.label("reciprocals") $X.group("relations")  $X.desc("Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)")  {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :TextEval $X.label("object");
}

struct Reduce $X.label("listReduce") $X.group("list")  $X.desc("Transform the values from one list by combining them into a single value. The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).")  {
  intoValue    @0  :Text $X.label("into");
  fromList     @1  :Assignment $X.label("from_list");
  usingPattern @2  :Text $X.label("using") $X.pool("pattern");
}

struct Relate $X.label("relate") $X.group("relations")  $X.desc("Relate two nouns.")  {
  object       @0  :TextEval;
  toObject     @1  :TextEval $X.label("to");
  str          @2  :Text $X.label("rel") $X.pool("relation");
}

struct Relation $X.label("relationName") $X.group("misc")  {
  at           @0  :Pos $X.internal;
  str          @1  :Text $X.label("rel") $X.pool("relation");
}

struct RelativeOf $X.label("relative") $X.group("relations")  $X.desc("Returns the relative of a noun (ex. the target of a one-to-one relation.)")  {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :TextEval $X.label("object");
}

struct RelativesOf $X.label("relatives") $X.group("relations")  $X.desc("Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).")  {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :TextEval $X.label("object");
}

struct RemainderOf $X.label("mod") $X.group("math")  $X.desc("Divide one number by another, and return the remainder.")  {
  a            @0  :NumberEval;
  b            @1  :NumberEval $X.label("by");
}

struct RenderField $X.label("renderField") $X.group("internal")  {
  name         @0  :TextEval;
}

struct RenderName $X.label("renderName") $X.group("internal")  {
  name         @0  :Text;
}

struct RenderPattern $X.label("render") $X.group("internal")  {
  pattern      @0  :Text $X.pool("pattern");
  args         @1  :List(Argument) $X.label("args");
}

struct RenderRef $X.label("renderRef") $X.group("internal")  {
  name         @0  :Text $X.pool("variable");
  flags        @1  :Int32 $X.label("flags");
}

struct RenderTemplate $X.label("renderTemplate") $X.group("format")  $X.desc("Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates")  {
  expression   @0  :TextEval $X.internal;
}

struct Response $X.label("response") $X.group("output")  $X.desc("Generate text in a replaceable manner.")  {
  name         @0  :Text;
  text         @1  :TextEval $X.optional $X.label("text");
}

struct Retarget $X.label("retarget") $X.group("grammar")  {
  span         @0  :List(ScannerMaker);
}

struct Reverse $X.label("reverse") $X.group("grammar")  {
  reverses     @0  :List(ScannerMaker);
}

struct ReverseList $X.label("listReverse") $X.group("list")  $X.desc("Reverse a list.")  {
  list         @0  :ListSource;
}

struct Row $X.label("row") $X.group("printing")  $X.desc("A single line as part of a group of lines.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Rows $X.label("rows") $X.group("printing")  $X.desc("Group text into successive lines.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Say $X.label("sayText") $X.group("printing")  $X.desc("Print some bit of text to the player.")  {
  text         @0  :TextEval;
}

struct Self $X.label("self") $X.group("grammar")  {
  player       @0  :Text;
}

struct Send $X.label("send") $X.group("events")  $X.desc("Triggers a event, returns a true/false success value.")  {
  event        @0  :Text $X.pool("event");
  path         @1  :TextListEval $X.label("to");
  args         @2  :List(Argument) $X.label("args");
}

struct Set $X.label("listSet") $X.group("list")  $X.desc("Overwrite an existing value in a list.")  {
  list         @0  :Text;
  index        @1  :NumberEval $X.label("index");
  from         @2  :Assignment $X.label("from");
}

struct SetTrait $X.label("put") $X.group("objects")  $X.desc("Put an object into a particular state.")  {
  object       @0  :TextEval $X.label("obj");
  trait        @1  :TextEval $X.label("trait");
}

struct ShuffleText $X.label("shuffle") $X.group("output")  $X.desc("When called multiple times returns its inputs at random.")  {
  seq          @0  :Text $X.internal;
  parts        @1  :List(TextEval);
}

struct Slash $X.label("slashText") $X.group("printing")  $X.desc("Separates words with left-leaning slashes '/'.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Slice $X.label("slice") $X.group("list")  $X.desc("Create a new list from a section of another list.")  {
  list         @0  :Assignment;
  start        @1  :NumberEval $X.optional $X.label("start");
  end          @2  :NumberEval $X.optional $X.label("end");
}

struct Softline $X.label("wbr") $X.group("printing")  $X.desc("Start a new line ( if not already at a new line. )")  {
}

struct SortByField $X.label("byField") $X.group("list")  {
  name         @0  :Text $X.label("by");
}

struct SortNumbers $X.label("sort") $X.group("list")  {
  str          @0  :Text $X.pool("variable");
  name         @1  :Text $X.label("by");
  order        @2  :Bool $X.label("order");
}

struct SortRecords $X.label("sort") $X.group("list")  $X.desc("Rearrange the elements in the named list by using the designated pattern to test pairs of elements.")  {
  str          @0  :Text $X.pool("variable");
  using        @1  :Text $X.label("using") $X.pool("pattern");
}

struct SortText $X.label("sort") $X.group("list")  $X.desc("Rearrange the elements in the named list by using the designated pattern to test pairs of elements.")  {
  str          @0  :Text $X.pool("variable");
  name         @1  :Text $X.label("by");
  order        @2  :Bool $X.label("order");
  case         @3  :Bool $X.label("case");
}

struct Span $X.label("spanText") $X.group("printing")  $X.desc("Writes text with spaces between words.")  {
  exe          @0  :List(Execute) $X.label("do");
}

struct Splice $X.label("splice") $X.group("list")  $X.desc("Modify a list by adding and removing elements. Note: the type of the elements being added must match the type of the list. Text cant be added to a list of numbers, numbers cant be added to a list of text. If the starting index is negative, it will begin that many elements from the end of the array. If list's length + the start is less than 0, it will begin from index 0. If the remove count is missing, it removes all elements from the start to the end; if it is 0 or negative, no elements are removed.")  {
  str          @0  :Text $X.pool("variable");
  start        @1  :NumberEval $X.label("start");
  remove       @2  :NumberEval $X.label("remove");
  insert       @3  :Assignment $X.label("insert");
}

struct StoppingText $X.label("stopping") $X.group("output")  $X.desc("When called multiple times returns each of its inputs in turn, sticking to the last one.")  {
  seq          @0  :Text $X.internal;
  parts        @1  :List(TextEval);
}

struct SumOf $X.label("inc") $X.group("math")  $X.desc("Add two numbers.")  {
  a            @0  :NumberEval;
  b            @1  :NumberEval $X.optional $X.label("by");
}

struct Text $X.label("txt") $X.group("literals")  $X.desc("Specify a small bit of text.")  {
  text         @0  :Text;
}

struct Texts $X.label("txts") $X.group("literals")  $X.desc("Specifies multiple string values.")  {
  values       @0  :List(Text);
}

struct TriggerCycle $X.label("every") $X.group("comparison")  {
}

struct TriggerOnce $X.label("at") $X.group("comparison")  {
}

struct TriggerSwitch $X.label("after") $X.group("comparison")  {
}

struct Var $X.label("var") $X.group("variables")  $X.desc("Return the value of the named variable.")  {
  name         @0  :Text $X.pool("variable");
}

struct Variable $X.label("var") $X.group("misc")  {
  at           @0  :Pos $X.internal;
  str          @1  :Text $X.pool("variable");
}

struct While $X.label("repeating") $X.group("flow")  $X.desc("Keep running a series of actions while a condition is true.")  {
  true         @0  :BoolEval;
  exe          @1  :List(Execute) $X.label("do");
}

struct Words $X.label("words") $X.group("grammar")  {
  words        @0  :List(Text);
}
