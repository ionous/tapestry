package internal

import "text/template"

// Create a new template and parse the letter into it.
var HeaderProto = template.Must(template.New("headerProto").Parse(headerProto))

//
const headerProto = `syntax = "proto3";
package pb;
option go_package = "git.sr.ht/~ionous/pb";
//
import "options.proto";`
