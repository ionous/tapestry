@0x950a507925528279;
using Go = import "/go.capnp";
using  X = import "../options.capnp";
using Rt = import "../rt/rt.capnp";
using Core = import "../core/core.capnp";

$Go.package("list");
$Go.import("git.sr.ht/~ionous/dl/list");

struct ListIterator { eval @0:AnyPointer; }
struct ListSource { eval @0:AnyPointer; }
struct ListTarget { eval @0:AnyPointer; }

struct AsNum $X.label("Num") $X.group("misc") $X.desc("Define the name of a number variable.") {
  str          @0  :Text $X.pool("variable");
}
struct AsRec $X.label("Rec") $X.group("misc") $X.desc("Define the name of a record variable.") {
  str          @0  :Text $X.pool("variable");
}
struct AsTxt $X.label("Txt") $X.group("misc") $X.desc("Define the name of a text variable.") {
  str          @0  :Text $X.pool("variable");
}
struct At $X.label("Get") $X.group("list") $X.desc("Get a value from a list. The first element is is index 1.") {
  list         @0  :Rt.Assignment;
  index        @1  :Rt.NumberEval $X.label("index");
}
struct Each $X.label("Repeating") $X.group("list") $X.desc("Loops over the elements in the passed list, or runs the 'else' activity if empty.") {
  list         @0  :Rt.Assignment $X.label("across");
  as           @1  :ListIterator $X.label("as");
  exe          @2  :List(Rt.Execute) $X.label("do");
  else         @3  :Core.Brancher $X.optional $X.label("else");
}
struct EraseEdge $X.label("Erase") $X.group("misc") $X.desc("Remove one or more values from a list") {
  from         @0  :ListSource;
  atEdge       @1  :Bool $X.label("at_edge");
}
struct EraseIndex $X.label("Erase") $X.group("misc") $X.desc("Remove one or more values from a list") {
  count        @0  :Rt.NumberEval;
  from         @1  :ListSource $X.label("from");
  atIndex      @2  :Rt.NumberEval $X.label("at_index");
}
struct Erasing $X.label("Erasing") $X.group("list") $X.desc("Erase elements from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.") {
  count        @0  :Rt.NumberEval;
  from         @1  :ListSource $X.label("from");
  atIndex      @2  :Rt.NumberEval $X.label("at_index");
  as           @3  :Text $X.label("as");
  exe          @4  :List(Rt.Execute) $X.label("do");
}
struct ErasingEdge $X.label("Erasing") $X.group("list") $X.desc("Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.") {
  from         @0  :ListSource;
  atEdge       @1  :Bool $X.label("at_edge");
  as           @2  :Text $X.label("as");
  exe          @3  :List(Rt.Execute) $X.label("do");
  else         @4  :Core.Brancher $X.optional $X.label("else");
}
struct Find $X.label("Find") $X.group("list") $X.desc("Search a list for a specific value.") {
  value        @0  :Rt.Assignment;
  list         @1  :Rt.Assignment $X.label("list");
}
struct FromNumList $X.label("Nums") $X.group("misc") $X.desc("Uses a list of numbers") {
  str          @0  :Text $X.pool("variable");
}
struct FromRecList $X.label("Recs") $X.group("misc") $X.desc("Uses a list of records") {
  str          @0  :Text $X.pool("variable");
}
struct FromTxtList $X.label("Txts") $X.group("misc") $X.desc("Uses a list of text") {
  str          @0  :Text $X.pool("variable");
}
struct Gather $X.label("Gather") $X.group("list") $X.desc("Transform the values from a list. The named pattern gets called once for each value in the list. It get called with two parameters: 'in' as each value from the list, and 'out' as the var passed to the gather.") {
  str          @0  :Text $X.pool("variable");
  from         @1  :ListSource $X.label("from");
  using        @2  :Text $X.label("using") $X.pool("pattern");
}
struct IntoNumList $X.label("Nums") $X.group("misc") $X.desc("Targets a list of numbers") {
  str          @0  :Text $X.pool("variable");
}
struct IntoRecList $X.label("Recs") $X.group("misc") $X.desc("Targets a list of records") {
  str          @0  :Text $X.pool("variable");
}
struct IntoTxtList $X.label("Txts") $X.group("misc") $X.desc("Targets a list of text") {
  str          @0  :Text $X.pool("variable");
}
struct Len $X.label("Len") $X.group("list") $X.desc("Determines the number of values in a list.") {
  list         @0  :Rt.Assignment;
}
struct Map $X.label("Map") $X.group("list") $X.desc("Transform the values from one list and place the results in another list. The designated pattern is called with each value from the 'from list', one value at a time.") {
  toList       @0  :Text;
  fromList     @1  :Rt.Assignment $X.label("from_list");
  usingPattern @2  :Text $X.label("using") $X.pool("pattern");
}
struct PutEdge $X.label("Put") $X.group("misc") $X.desc("Add a value to a list") {
  from         @0  :Rt.Assignment;
  into         @1  :ListTarget $X.label("into");
  atEdge       @2  :Bool $X.label("at_edge");
}
struct PutIndex $X.label("Put") $X.group("misc") $X.desc("Replace one value in a list with another") {
  from         @0  :Rt.Assignment;
  into         @1  :ListTarget $X.label("into");
  atIndex      @2  :Rt.NumberEval $X.label("at_index");
}
struct Range $X.label("Range") $X.group("flow") $X.desc("Generates a series of numbers.") {
  to           @0  :Rt.NumberEval;
  from         @1  :Rt.NumberEval $X.optional $X.label("from");
  byStep       @2  :Rt.NumberEval $X.optional $X.label("by_step");
}
struct Reduce $X.label("ListReduce") $X.group("list") $X.desc("Transform the values from one list by combining them into a single value. The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).") {
  intoValue    @0  :Text $X.label("into");
  fromList     @1  :Rt.Assignment $X.label("from_list");
  usingPattern @2  :Text $X.label("using") $X.pool("pattern");
}
struct ReverseList $X.label("ListReverse") $X.group("list") $X.desc("Reverse a list.") {
  list         @0  :ListSource;
}
struct Set $X.label("ListSet") $X.group("list") $X.desc("Overwrite an existing value in a list.") {
  list         @0  :Text;
  index        @1  :Rt.NumberEval $X.label("index");
  from         @2  :Rt.Assignment $X.label("from");
}
struct Slice $X.label("Slice") $X.group("list") $X.desc("Create a new list from a section of another list.") {
  list         @0  :Rt.Assignment;
  start        @1  :Rt.NumberEval $X.optional $X.label("start");
  end          @2  :Rt.NumberEval $X.optional $X.label("end");
}
struct SortByField $X.label("ByField") $X.group("list") {
  name         @0  :Text $X.label("by");
}
struct SortNumbers $X.label("Sort") $X.group("list") {
  str          @0  :Text $X.pool("variable");
  name         @1  :Text $X.label("by");
  order        @2  :Bool $X.label("order");
}
struct SortRecords $X.label("Sort") $X.group("list") $X.desc("Rearrange the elements in the named list by using the designated pattern to test pairs of elements.") {
  str          @0  :Text $X.pool("variable");
  using        @1  :Text $X.label("using") $X.pool("pattern");
}
struct SortText $X.label("Sort") $X.group("list") $X.desc("Rearrange the elements in the named list by using the designated pattern to test pairs of elements.") {
  str          @0  :Text $X.pool("variable");
  name         @1  :Text $X.label("by");
  order        @2  :Bool $X.label("order");
  case         @3  :Bool $X.label("case");
}
struct Splice $X.label("Splice") $X.group("list") $X.desc("Modify a list by adding and removing elements. Note: the type of the elements being added must match the type of the list. Text cant be added to a list of numbers, numbers cant be added to a list of text. If the starting index is negative, it will begin that many elements from the end of the array. If list's length + the start is less than 0, it will begin from index 0. If the remove count is missing, it removes all elements from the start to the end; if it is 0 or negative, no elements are removed.") {
  str          @0  :Text $X.pool("variable");
  start        @1  :Rt.NumberEval $X.label("start");
  remove       @2  :Rt.NumberEval $X.label("remove");
  insert       @3  :Rt.Assignment $X.label("insert");
}
