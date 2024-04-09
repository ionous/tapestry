// Package flex reads tell files that are sectioned into alternating
// blocks of structured and plain text sections.
// The plain text sections are wrapped with commands and
// merged into the structured sections.
// The plain text sections can also "jump out" into structured sections
// on lines ending with colons.
package flex
