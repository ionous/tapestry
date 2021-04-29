@0xac0522bccd6b09e7;
using Go = import "/go.capnp";
using  X = import "../options.capnp";

$Go.package("grammar");
$Go.import("git.sr.ht/~ionous/dl/grammar");

struct GrammarMaker { eval @0:AnyPointer; }
struct ScannerMaker { eval @0:AnyPointer; }

struct Action $X.label("As") $X.group("grammar") {
  action       @0  :Text;
}
struct Alias $X.label("Alias") $X.group("grammar") {
  names        @0  :List(Text);
  asNoun       @1  :Text $X.label("as_noun");
}
struct AllOf $X.label("AllOf") $X.group("grammar") {
  series       @0  :List(ScannerMaker);
}
struct AnyOf $X.label("AnyOf") $X.group("grammar") {
  options      @0  :List(ScannerMaker);
}
struct Directive $X.label("Directive") $X.group("grammar") {
  lede         @0  :List(Text);
  scans        @1  :List(ScannerMaker) $X.label("scans");
}
struct GrammarDecl $X.label("Grammar") $X.group("grammar") $X.desc("Read what the player types and turn it into actions.") {
  grammar      @0  :GrammarMaker;
}
struct Noun $X.label("Noun") $X.group("grammar") {
  kind         @0  :Text;
}
struct Retarget $X.label("Retarget") $X.group("grammar") {
  span         @0  :List(ScannerMaker);
}
struct Reverse $X.label("Reverse") $X.group("grammar") {
  reverses     @0  :List(ScannerMaker);
}
struct Self $X.label("Self") $X.group("grammar") {
  player       @0  :Text;
}
struct Words $X.label("Words") $X.group("grammar") {
  words        @0  :List(Text);
}
