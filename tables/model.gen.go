/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// modelTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func modelTemplate() string {
	var tmpl = "/**\n" +
		" * tables describing the the world and its rules.\n" +
		" */\n" +
		"\n" +
		"/* initialization for fields of kinds.\n" +
		" * the value ( an Assignment ) is defined when the field is defined, \n" +
		" * meaning there's no separate domain or origin(at). */ \n" +
		"create table mdl_assign( field int not null, value blob, primary key( field ) );\n" +
		"/* stored tests, which run a program to verify it produces the expected value. \n" +
		" * fix: shouldnt we also be writing class of the value ? */ \n" +
		"create table mdl_check( domain int not null, name text, value blob, affinity text, prog blob, at text, primary key( domain, name ));\n" +
		"/* hierarchy of scopes */\n" +
		"create table mdl_domain( domain text, path text, at text, primary key( domain ));\n" +
		"/* properties for a kind. \n" +
		" * type is most often used for affinities of type \"text\", and usually indicates a kind ( from mdl_kind )\n" +
		" * the scope of a field is the same as its kind. */\n" +
		"create table mdl_field( kind int not null, field text, affinity text, type int, at text, primary key( kind, field ));\n" +
		"/* statements for user input parsing */\n" +
		"create table mdl_grammar( domain int not null, name text, prog blob, at text, primary key( domain, name ));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind( domain int not null, kind text, path text, at text, primary key( domain, kind ));\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		" * the lower rank of the association wins. \n" +
		" * in theory, the scope of a name can be narrower that the scope of its noun */\n" +
		"create table mdl_name( domain int not null, noun int not null, name text, rank int, at text );\n" +
		"/* a person, place, or thing in the world. \n" +
		" * domain tells the scope in which the noun was defined.\n" +
		" * ( the same as - or a child of - the domain of the kind ) */\n" +
		"create table mdl_noun( domain int not null, noun text, kind int not null, at text, primary key( domain, noun ));\n" +
		"/* relation between two specific  nnouns. these change over the course of a game.\n" +
		" * similar to mdl_rule, points back to the relation's kind rather than the entry in the mdl_rel table */\n" +
		"create table mdl_pair( domain int not null, relKind int not null, oneNoun int not null, otherNoun int not null, at text );\n" +
		"/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields \n" +
		" * the scope of the pattern is the same as its kind */\n" +
		"create table mdl_pat( kind int not null, labels text, result text, primary key( kind ) );\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( domain int, many text, one text, at text, primary key( domain, many ) );\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel( relKind int not null, oneKind int not null, otherKind int not null, cardinality text, at text, primary key( relKind ));\n" +
		"/* the scope of a rule can be narrower than its parent kind ( or target )\n" +
		" * fix? can target (kind of nouns the rule applies to) be moved to filter? */\n" +
		"create table mdl_rule( domain int not null, kind int not null, target int, phase int, filter blob, prog blob, at text );\n" +
		"/* initial values for various nouns.\n" +
		" * note: currently, the scope of the value is the same as the noun.\n" +
		" * ( ie. one value per noun per field; not values that change as new domains are activated. )\n" +
		" + the affinity and subtype of the value come from the field.\n" +
		" */\n" +
		"create table mdl_value( noun int not null, field int not null, value blob, at text, primary key( noun, field ));\n" +
		""
	return tmpl
}
