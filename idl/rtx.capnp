@0x8167b37f5d971356;
using Go = import "/go.capnp";
using  X = import "options.capnp";

$Go.package("rtx");
$Go.import("git.sr.ht/~ionous/iffy/idl/rtx");

struct Assignment { eval @0:AnyPointer; }
struct BoolEval { eval @0:AnyPointer; }
struct Execute { eval @0:AnyPointer; }
struct NumListEval { eval @0:AnyPointer; }
struct NumberEval { eval @0:AnyPointer; }
struct RecordEval { eval @0:AnyPointer; }
struct RecordListEval { eval @0:AnyPointer; }
struct TextEval { eval @0:AnyPointer; }
struct TextListEval { eval @0:AnyPointer; }

