 @0x838375eaedd19910;
using Go = import "/go.capnp";
using  X = import "options.capnp";
using Core = import "core.capnp";
using Debug = import "debug.capnp";
using Grammar = import "grammar.capnp";
using List = import "list.capnp";
using Rel = import "rel.capnp";
using Render = import "render.capnp";
using Rtx = import "rtx.capnp";

$Go.package("all");
$Go.import("git.sr.ht/~ionous/iffy/idl/all");

struct AssignmentImpl $X.desc("Helper for setting variables.") {
  union {
	bool                           @0   :Core.FromBool $X.label("Bool:");
	getFrom                        @1   :Core.GetAtField $X.label("Get:from:");
	num                            @2   :Core.FromNum $X.label("Num:");
	nums                           @3   :Core.FromNumbers $X.label("Nums:");
	rec                            @4   :Core.FromRecord $X.label("Rec:");
	recs                           @5   :Core.FromRecords $X.label("Recs:");
	renderArgs                     @6   :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @7   :Render.RenderRef $X.label("RenderRef:flags:");
	txt                            @8   :Core.FromText $X.label("Txt:");
	txts                           @9   :Core.FromTexts $X.label("Txts:");
	var                            @10  :Core.Var $X.label("Var:");
  }
}
struct BoolEvalImpl $X.desc("Statements which return true/false values.") {
  union {
	allOf                          @0   :Core.AllTrue $X.label("AllOf:");
	always                         @1   :Core.Always $X.label("Always");
	anyOf                          @2   :Core.AnyTrue $X.label("AnyOf:");
	bool                           @3   :Core.BoolValue $X.label("Bool:");
	cmpIsNum                       @4   :Core.CompareNum $X.label("Cmp:is:num:");
	cmpIsTxt                       @5   :Core.CompareText $X.label("Cmp:is:txt:");
	containsPart                   @6   :Core.Includes $X.label("Contains:part:");
	countOfTrigger                 @7   :Core.CountOf $X.label("CountOf:trigger:");
	determineArgs                  @8   :Core.Determine $X.label("Determine:args:");
	during                         @9   :Core.During $X.label("During:");
	findList                       @10  :List.Find $X.label("Find:list:");
	getObjTrait                    @11  :Core.HasTrait $X.label("Get obj:trait:");
	getFrom                        @12  :Core.GetAtField $X.label("Get:from:");
	hasDominion                    @13  :Core.HasDominion $X.label("HasDominion:");
	isEmpty                        @14  :Core.IsEmpty $X.label("Is empty:");
	isValid                        @15  :Core.ObjectExists $X.label("Is valid:");
	kindOfIs                       @16  :Core.IsKindOf $X.label("KindOf:is:");
	kindOfIsExactly                @17  :Core.IsExactKindOf $X.label("KindOf:isExactly:");
	matchesTo                      @18  :Core.Matches $X.label("Matches:to:");
	not                            @19  :Core.IsNotTrue $X.label("Not:");
	renderArgs                     @20  :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @21  :Render.RenderRef $X.label("RenderRef:flags:");
	sendToArgs                     @22  :Core.Send $X.label("Send:to:args:");
	var                            @23  :Core.Var $X.label("Var:");
  }
}
struct BrancherImpl $X.desc("Helper for choose action.") {
  union {
	elseDoDo                       @0   :Core.ChooseNothingElse $X.label("ElseDo do:");
	elseIfDo                       @1   :Core.ChooseMore $X.label("ElseIf:do:");
	elseIfDoElse                   @2   :Core.ChooseMore $X.label("ElseIf:do:else:");
	elseIfFromAndDo                @3   :Core.ChooseMoreValue $X.label("ElseIf:from:and:do:");
	elseIfFromAndDoElse            @4   :Core.ChooseMoreValue $X.label("ElseIf:from:and:do:else:");
	ifDo                           @5   :Core.ChooseAction $X.label("If:do:");
	ifDoElse                       @6   :Core.ChooseAction $X.label("If:do:else:");
	ifFromAndDo                    @7   :Core.ChooseValue $X.label("If:from:and:do:");
	ifFromAndDoElse                @8   :Core.ChooseValue $X.label("If:from:and:do:else:");
  }
}
struct ComparatorImpl $X.desc("Helper for comparing values.") {
  union {
	atLeast                        @0   :Core.GreaterOrEqual $X.label("AtLeast");
	atMost                         @1   :Core.LessOrEqual $X.label("AtMost");
	equalTo                        @2   :Core.EqualTo $X.label("EqualTo");
	greaterThan                    @3   :Core.GreaterThan $X.label("GreaterThan");
	lessThan                       @4   :Core.LessThan $X.label("LessThan");
	otherThan                      @5   :Core.NotEqualTo $X.label("OtherThan");
  }
}
struct ExecuteImpl $X.desc("Run a series of statements.") {
  union {
	actDo                          @0   :Core.Activity $X.label("Act do:");
	br                             @1   :Core.Newline $X.label("Br");
	break                          @2   :Core.Break $X.label("Break");
	determineArgs                  @3   :Core.Determine $X.label("Determine:args:");
	doNothing                      @4   :Debug.DoNothing $X.label("DoNothing");
	doNothingWhy                   @5   :Debug.DoNothing $X.label("DoNothing why:");
	eraseAtEdge                    @6   :List.EraseEdge $X.label("Erase:atEdge:");
	eraseFromAtIndex               @7   :List.EraseIndex $X.label("Erase:from:atIndex:");
	erasingAtEdgeAsDo              @8   :List.ErasingEdge $X.label("Erasing:atEdge:as:do:");
	erasingAtEdgeAsDoElse          @9   :List.ErasingEdge $X.label("Erasing:atEdge:as:do:else:");
	erasingFromAtIndexAsDo         @10  :List.Erasing $X.label("Erasing:from:atIndex:as:do:");
	ifDo                           @11  :Core.ChooseAction $X.label("If:do:");
	ifDoElse                       @12  :Core.ChooseAction $X.label("If:do:else:");
	ifFromAndDo                    @13  :Core.ChooseValue $X.label("If:from:and:do:");
	ifFromAndDoElse                @14  :Core.ChooseValue $X.label("If:from:and:do:else:");
	letBe                          @15  :Core.Assign $X.label("Let:be:");
	listReduceIntoFromListUsing    @16  :List.Reduce $X.label("ListReduce into:fromList:using:");
	listReverse                    @17  :List.ReverseList $X.label("ListReverse:");
	listSetIndexFrom               @18  :List.Set $X.label("ListSet:index:from:");
	logValue                       @19  :Debug.Log $X.label("Log:value:");
	mapFromListUsing               @20  :List.Map $X.label("Map:fromList:using:");
	next                           @21  :Core.Next $X.label("Next");
	p                              @22  :Core.Paragraph $X.label("P");
	putObjTrait                    @23  :Core.SetTrait $X.label("Put obj:trait:");
	putFromAt                      @24  :Core.PutAtField $X.label("Put:from:at:");
	putIntoAtEdge                  @25  :List.PutEdge $X.label("Put:into:atEdge:");
	putIntoAtIndex                 @26  :List.PutIndex $X.label("Put:into:atIndex:");
	relateToRel                    @27  :Rel.Relate $X.label("Relate:to:rel:");
	renderArgs                     @28  :Render.RenderPattern $X.label("Render:args:");
	repeatingAcrossAsDo            @29  :List.Each $X.label("Repeating across:as:do:");
	repeatingAcrossAsDoElse        @30  :List.Each $X.label("Repeating across:as:do:else:");
	repeatingDo                    @31  :Core.While $X.label("Repeating:do:");
	sayText                        @32  :Core.Say $X.label("SayText:");
	sendToArgs                     @33  :Core.Send $X.label("Send:to:args:");
	sortByOrder                    @34  :List.SortNumbers $X.label("Sort:by:order:");
	sortByOrderCase                @35  :List.SortText $X.label("Sort:by:order:case:");
	sortUsing                      @36  :List.SortRecords $X.label("Sort:using:");
	spliceStartRemoveInsert        @37  :List.Splice $X.label("Splice:start:remove:insert:");
	wbr                            @38  :Core.Softline $X.label("Wbr");
  }
}
struct FromSourceFieldsImpl $X.desc("Helper for getting fields.") {
  union {
	obj                            @0   :Core.FromObj $X.label("Obj:");
	rec                            @1   :Core.FromRec $X.label("Rec:");
	renderField                    @2   :Render.RenderField $X.label("RenderField:");
	var                            @3   :Core.FromVar $X.label("Var:");
  }
}
struct GrammarMakerImpl $X.desc("Helper for defining parser grammars.") {
  union {
	aliasAsNoun                    @0   :Grammar.Alias $X.label("Alias:asNoun:");
	directiveScans                 @1   :Grammar.Directive $X.label("Directive:scans:");
  }
}
struct IntoTargetFieldsImpl $X.desc("Helper for setting fields.") {
  union {
	obj                            @0   :Core.IntoObj $X.label("Obj:");
	var                            @1   :Core.IntoVar $X.label("Var:");
  }
}
struct ListIteratorImpl $X.desc("Helper for accessing lists.") {
  union {
	num                            @0   :List.AsNum $X.label("Num:");
	rec                            @1   :List.AsRec $X.label("Rec:");
	txt                            @2   :List.AsTxt $X.label("Txt:");
  }
}
struct ListSourceImpl $X.desc("Helper for accessing lists.") {
  union {
	nums                           @0   :List.FromNumList $X.label("Nums:");
	recs                           @1   :List.FromRecList $X.label("Recs:");
	txts                           @2   :List.FromTxtList $X.label("Txts:");
  }
}
struct ListTargetImpl $X.desc("Helper for accessing lists.") {
  union {
	nums                           @0   :List.IntoNumList $X.label("Nums:");
	recs                           @1   :List.IntoRecList $X.label("Recs:");
	txts                           @2   :List.IntoTxtList $X.label("Txts:");
  }
}
struct NumListEvalImpl $X.desc("Statements which return a list of numbers.") {
  union {
	determineArgs                  @0   :Core.Determine $X.label("Determine:args:");
	getFrom                        @1   :Core.GetAtField $X.label("Get:from:");
	nums                           @2   :Core.NumList $X.label("Nums:");
	range                          @3   :List.Range $X.label("Range:");
	rangeByStep                    @4   :List.Range $X.label("Range:byStep:");
	rangeFrom                      @5   :List.Range $X.label("Range:from:");
	rangeFromByStep                @6   :List.Range $X.label("Range:from:byStep:");
	renderArgs                     @7   :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @8   :Render.RenderRef $X.label("RenderRef:flags:");
	slice                          @9   :List.Slice $X.label("Slice:");
	sliceEnd                       @10  :List.Slice $X.label("Slice:end:");
	sliceStart                     @11  :List.Slice $X.label("Slice:start:");
	sliceStartEnd                  @12  :List.Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @13  :List.Splice $X.label("Splice:start:remove:insert:");
	var                            @14  :Core.Var $X.label("Var:");
  }
}
struct NumberEvalImpl $X.desc("Statements which return a number.") {
  union {
	dec                            @0   :Core.DiffOf $X.label("Dec:");
	decBy                          @1   :Core.DiffOf $X.label("Dec:by:");
	determineArgs                  @2   :Core.Determine $X.label("Determine:args:");
	divBy                          @3   :Core.QuotientOf $X.label("Div:by:");
	during                         @4   :Core.During $X.label("During:");
	findList                       @5   :List.Find $X.label("Find:list:");
	getFrom                        @6   :Core.GetAtField $X.label("Get:from:");
	getIndex                       @7   :List.At $X.label("Get:index:");
	inc                            @8   :Core.SumOf $X.label("Inc:");
	incBy                          @9   :Core.SumOf $X.label("Inc:by:");
	len                            @10  :List.Len $X.label("Len:");
	modBy                          @11  :Core.RemainderOf $X.label("Mod:by:");
	mulBy                          @12  :Core.ProductOf $X.label("Mul:by:");
	num                            @13  :Core.NumValue $X.label("Num:");
	numIfElse                      @14  :Core.ChooseNum $X.label("Num:if:else:");
	renderArgs                     @15  :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @16  :Render.RenderRef $X.label("RenderRef:flags:");
	var                            @17  :Core.Var $X.label("Var:");
  }
}
struct RecordEvalImpl $X.desc("Statements which return a record.") {
  union {
	determineArgs                  @0   :Core.Determine $X.label("Determine:args:");
	getFrom                        @1   :Core.GetAtField $X.label("Get:from:");
	getIndex                       @2   :List.At $X.label("Get:index:");
	makeArgs                       @3   :Core.Make $X.label("Make:args:");
	renderArgs                     @4   :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @5   :Render.RenderRef $X.label("RenderRef:flags:");
	var                            @6   :Core.Var $X.label("Var:");
  }
}
struct RecordListEvalImpl $X.desc("Statements which return a list of records.") {
  union {
	determineArgs                  @0   :Core.Determine $X.label("Determine:args:");
	getFrom                        @1   :Core.GetAtField $X.label("Get:from:");
	renderArgs                     @2   :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @3   :Render.RenderRef $X.label("RenderRef:flags:");
	slice                          @4   :List.Slice $X.label("Slice:");
	sliceEnd                       @5   :List.Slice $X.label("Slice:end:");
	sliceStart                     @6   :List.Slice $X.label("Slice:start:");
	sliceStartEnd                  @7   :List.Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @8   :List.Splice $X.label("Splice:start:remove:insert:");
	var                            @9   :Core.Var $X.label("Var:");
  }
}
struct ScannerMakerImpl $X.desc("Helper for defining input scanners.") {
  union {
	allOf                          @0   :Grammar.AllOf $X.label("AllOf:");
	anyOf                          @1   :Grammar.AnyOf $X.label("AnyOf:");
	as                             @2   :Grammar.Action $X.label("As:");
	noun                           @3   :Grammar.Noun $X.label("Noun:");
	retarget                       @4   :Grammar.Retarget $X.label("Retarget:");
	reverse                        @5   :Grammar.Reverse $X.label("Reverse:");
	self                           @6   :Grammar.Self $X.label("Self:");
	words                          @7   :Grammar.Words $X.label("Words:");
  }
}
struct TextEvalImpl $X.desc("Statements which return text.") {
  union {
	bracketTextDo                  @0   :Core.Bracket $X.label("BracketText do:");
	bufferTextDo                   @1   :Core.Buffer $X.label("BufferText do:");
	capitalize                     @2   :Core.Capitalize $X.label("Capitalize:");
	commaTextDo                    @3   :Core.Commas $X.label("CommaText do:");
	cycle                          @4   :Core.CycleText $X.label("Cycle:");
	determineArgs                  @5   :Core.Determine $X.label("Determine:args:");
	getFrom                        @6   :Core.GetAtField $X.label("Get:from:");
	getIndex                       @7   :List.At $X.label("Get:index:");
	idOf                           @8   :Core.IdOf $X.label("IdOf:");
	joinParts                      @9   :Core.Join $X.label("Join:parts:");
	kindOf                         @10  :Core.KindOf $X.label("KindOf:");
	lower                          @11  :Core.MakeLowercase $X.label("Lower:");
	nameOf                         @12  :Core.NameOf $X.label("NameOf:");
	numeralWords                   @13  :Core.PrintNumWord $X.label("Numeral words:");
	numeral                        @14  :Core.PrintNum $X.label("Numeral:");
	pluralize                      @15  :Core.MakePlural $X.label("Pluralize:");
	reciprocalRelObject            @16  :Rel.ReciprocalOf $X.label("Reciprocal rel:object:");
	relativeRelObject              @17  :Rel.RelativeOf $X.label("Relative rel:object:");
	renderArgs                     @18  :Render.RenderPattern $X.label("Render:args:");
	renderName                     @19  :Render.RenderName $X.label("RenderName:");
	renderRefFlags                 @20  :Render.RenderRef $X.label("RenderRef:flags:");
	renderTemplate                 @21  :Render.RenderTemplate $X.label("RenderTemplate");
	response                       @22  :Core.Response $X.label("Response:");
	responseText                   @23  :Core.Response $X.label("Response:text:");
	reverse                        @24  :Core.MakeReversed $X.label("Reverse:");
	rowDo                          @25  :Core.Row $X.label("Row do:");
	rowsDo                         @26  :Core.Rows $X.label("Rows do:");
	sentence                       @27  :Core.MakeSentenceCase $X.label("Sentence:");
	shuffle                        @28  :Core.ShuffleText $X.label("Shuffle:");
	singularize                    @29  :Core.MakeSingular $X.label("Singularize:");
	slashTextDo                    @30  :Core.Slash $X.label("SlashText do:");
	spanTextDo                     @31  :Core.Span $X.label("SpanText do:");
	stopping                       @32  :Core.StoppingText $X.label("Stopping:");
	title                          @33  :Core.MakeTitleCase $X.label("Title:");
	txt                            @34  :Core.TextValue $X.label("Txt:");
	txtIfElse                      @35  :Core.ChooseText $X.label("Txt:if:else:");
	upper                          @36  :Core.MakeUppercase $X.label("Upper:");
	var                            @37  :Core.Var $X.label("Var:");
  }
}
struct TextListEvalImpl $X.desc("Statements which return a list of text.") {
  union {
	determineArgs                  @0   :Core.Determine $X.label("Determine:args:");
	getFrom                        @1   :Core.GetAtField $X.label("Get:from:");
	kindsOf                        @2   :Core.KindsOf $X.label("KindsOf:");
	reciprocalsRelObject           @3   :Rel.ReciprocalsOf $X.label("Reciprocals rel:object:");
	relativesRelObject             @4   :Rel.RelativesOf $X.label("Relatives rel:object:");
	renderArgs                     @5   :Render.RenderPattern $X.label("Render:args:");
	renderRefFlags                 @6   :Render.RenderRef $X.label("RenderRef:flags:");
	slice                          @7   :List.Slice $X.label("Slice:");
	sliceEnd                       @8   :List.Slice $X.label("Slice:end:");
	sliceStart                     @9   :List.Slice $X.label("Slice:start:");
	sliceStartEnd                  @10  :List.Slice $X.label("Slice:start:end:");
	spliceStartRemoveInsert        @11  :List.Splice $X.label("Splice:start:remove:insert:");
	txts                           @12  :Core.TextList $X.label("Txts:");
	var                            @13  :Core.Var $X.label("Var:");
  }
}
struct TriggerImpl $X.desc("Helper for counting values.") {
  union {
	after                          @0   :Core.TriggerSwitch $X.label("After");
	at                             @1   :Core.TriggerOnce $X.label("At");
	every                          @2   :Core.TriggerCycle $X.label("Every");
  }
}