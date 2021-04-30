@0xe5343b8be2461313;
using Go = import "/go.capnp";
using  X = import "options.capnp";
using Rtx = import "rtx.capnp";

$Go.package("debug");
$Go.import("git.sr.ht/~ionous/iffy/idl/debug");


struct DoNothing $X.label("DoNothing") $X.group("flow") $X.desc("Statement which does nothing.") {
  reason       @0  :Text $X.optional $X.label("why");
}
struct Log $X.label("Log") $X.group("debug") $X.desc("Debug log") {
  level        @0  :Int32;
  value        @1  :Rtx.Assignment $X.label("value");
}
