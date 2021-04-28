package internal

import "text/template"

// Create a new template and parse the letter into it.
var HeaderCap = template.Must(template.New("headerCap").Parse(headerCap))

//
const headerCap = `@0xad22bd0042f92910;
using Go = import "/go.capnp";
using  X = import "./options.capnp";

$Go.package("auto");
$Go.import("git.sr.ht/~ionous/pb/auto");
`
