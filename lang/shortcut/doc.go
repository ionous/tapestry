// Package shortcut transforms strings into dotted paths.
//
// Some examples:
//
//	"@var.5"
//	"#obj.field"
//	"#`i have some complex name`.field"
//	"@@ starts with a single at sign"
//	"## starts with a single hash mark"
//
// Rules:
//   - names consist of alpha numeric runes, underscores are transformed into spaces,
//     case is ignored but lowercase is preferred.
//   - paths are composed of a target and subsequent path parts.
//   - targets start with either "@" for variables, or "#" for objects followed by a name.
//   - path parts follow the target, begin with  ".", and are either names or numbers.
//   - as a special case, a name can also be enclosed in backtick or doublequotes.
//     ( usually it will be backtick because these are embedded in json style strings,
//     so the double quotes would have to be escaped in the original
package shortcut
