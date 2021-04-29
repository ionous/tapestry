@0xf445246ceb74f07d;
using Go = import "/go.capnp";
using  X = import "../options.capnp";
using Rt = import "../rt/rt.capnp";
using Core = import "../core/core.capnp";

$Go.package("render");
$Go.import("git.sr.ht/~ionous/dl/render");


struct RenderField $X.label("RenderField") $X.group("internal") {
  name         @0  :Rt.TextEval;
}
struct RenderName $X.label("RenderName") $X.group("internal") {
  name         @0  :Text;
}
struct RenderPattern $X.label("Render") $X.group("internal") {
  pattern      @0  :Text $X.pool("pattern");
  args         @1  :List(Core.Argument) $X.label("args");
}
struct RenderRef $X.label("RenderRef") $X.group("internal") {
  name         @0  :Text $X.pool("variable");
  flags        @1  :Int32 $X.label("flags");
}
struct RenderTemplate $X.label("RenderTemplate") $X.group("format") $X.desc("Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates") {
  expression   @0  :Rt.TextEval $X.internal;
}
