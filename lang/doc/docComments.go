package doc

import (
	"go/doc/comment"
	"html/template"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type docComments struct {
	types []typeinfo.TypeSet
	// current *typeinfo.TypeSet
}

func (d *docComments) formatComment(lines []string) template.HTML {
	// use go's document parser to parse the document header
	// ( since its also used for the package comment, might as well.
	text := strings.Join(lines, "\n")
	p := comment.Parser{
		LookupPackage: d.LookupPackage,
		LookupSym:     d.LookupSym,
		// 	Words: map[string]string{
		// 	"italicword": "",
		// 	"linkedword": "https://example.com/linkedword",
		// },
	}
	doc := p.Parse(text)
	var pr comment.Printer
	return template.HTML(pr.HTML(doc))
}

// LookupPackage resolves a package name to an import path.
//
// If LookupPackage(name) returns ok == true, then [name]
// (or [name.Sym] or [name.Sym.Method])
// is considered a documentation link to importPath's package docs.
// It is valid to return "", true, in which case name is considered
// to refer to the current package.
//
// If LookupPackage(name) returns ok == false,
// then [name] (or [name.Sym] or [name.Sym.Method])
// will not be considered a documentation link,
// except in the case where name is the full (but single-element) import path
// of a package in the standard library, such as in [math] or [io.Reader].
// LookupPackage is still called for such names,
// in order to permit references to imports of other packages
// with the same package names.
//
// Setting LookupPackage to nil is equivalent to setting it to
// a function that always returns "", false.
func (d *docComments) LookupPackage(name string) (importPath string, ok bool) {
	return "", false
}

// LookupSym reports whether a symbol name or method name
// exists in the current package.
//
// If LookupSym("", "Name") returns true, then [Name]
// is considered a documentation link for a const, func, type, or var.
//
// Similarly, if LookupSym("Recv", "Name") returns true,
// then [Recv.Name] is considered a documentation link for
// type Recv's method Name.
//
// Setting LookupSym to nil is equivalent to setting it to a function
// that always returns false.
func (d *docComments) LookupSym(recv, name string) (ok bool) {
	// when true, it generates a local #link
	// some of the valid [] names would be idl links or links to other
	return false
}
