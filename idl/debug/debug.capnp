@0xba446b5cbc6e7ad3;
using Go = import "/go.capnp";
using  X = import "../options.capnp";
using Rt = import "../rt/rt.capnp";

$Go.package("debug");
$Go.import("git.sr.ht/~ionous/dl/debug");


struct DoNothing $X.label("DoNothing") $X.group("flow") $X.desc("Statement which does nothing.") {
  reason       @0  :Text $X.optional $X.label("why");
}
struct Log $X.label("Log") $X.group("debug") $X.desc("Debug log") {
  level        @0  :Int32;
  value        @1  :Rt.Assignment $X.label("value");
}
