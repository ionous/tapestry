@0xb5d9a32bef577d63;
using Go = import "/go.capnp";
using  X = import "../options.capnp";
using Rt = import "../rt/rt.capnp";
using Reader = import "../reader/reader.capnp";

$Go.package("rel");
$Go.import("git.sr.ht/~ionous/dl/rel");


struct ReciprocalOf $X.label("Reciprocal") $X.group("relations") $X.desc("Returns the implied relative of a noun (ex. the source in a one-to-many relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rt.TextEval $X.label("object");
}
struct ReciprocalsOf $X.label("Reciprocals") $X.group("relations") $X.desc("Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rt.TextEval $X.label("object");
}
struct Relate $X.label("Relate") $X.group("relations") $X.desc("Relate two nouns.") {
  object       @0  :Rt.TextEval;
  toObject     @1  :Rt.TextEval $X.label("to");
  str          @2  :Text $X.label("rel") $X.pool("relation");
}
struct Relation $X.label("RelationName") $X.group("misc") {
  at           @0  :Reader.Pos $X.internal;
  str          @1  :Text $X.label("rel") $X.pool("relation");
}
struct RelativeOf $X.label("Relative") $X.group("relations") $X.desc("Returns the relative of a noun (ex. the target of a one-to-one relation.)") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rt.TextEval $X.label("object");
}
struct RelativesOf $X.label("Relatives") $X.group("relations") $X.desc("Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).") {
  str          @0  :Text $X.label("rel") $X.pool("relation");
  object       @1  :Rt.TextEval $X.label("object");
}
