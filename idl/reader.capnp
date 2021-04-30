 @0xcc763cb346d64bd9;
 using Go = import "/go.capnp";

$Go.package("reader");
$Go.import("git.sr.ht/~ionous/iffy/idl/reader");

struct Pos {
	source @0 :Text; 
	offset @1 :Text;
}

