@0x80f117e3931e76db;
using Go = import "/go.capnp";
using  X = import "options.capnp";
using Rtx = import "rtx.capnp";
using Reader = import "reader.capnp";

$Go.package("rel");
$Go.import("git.sr.ht/~ionous/iffy/idl/rel");


struct ReciprocalOf $X.label("Reciprocal") $X.group("relations") $X.desc("Returns the implied relative of a noun (ex. the source in a one-to-many relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rtx.TextEval $X.label("object");
}
struct ReciprocalsOf $X.label("Reciprocals") $X.group("relations") $X.desc("Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rtx.TextEval $X.label("object");
}
struct Relate $X.label("Relate") $X.group("relations") $X.desc("Relate two nouns.") {
  object       @0  :Rtx.TextEval;
  toObject     @1  :Rtx.TextEval $X.label("to");
  str          @2  :Text $X.label("rel") $X.pool("relation");
}
struct Relation $X.label("RelationName") $X.group("misc") {
  at           @0  :Reader.Pos $X.internal;
  str          @1  :Text $X.label("rel") $X.pool("relation");
}
struct RelativeOf $X.label("Relative") $X.group("relations") $X.desc("Returns the relative of a noun (ex. the target of a one-to-one relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rtx.TextEval $X.label("object");
}
struct RelativesOf $X.label("Relatives") $X.group("relations") $X.desc("Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rtx.TextEval $X.label("object");
}
