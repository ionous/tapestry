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
		"/* enumerated values used by kinds and nouns */\n" +
		"create table mdl_aspect( domain text, aspect text, trait text, rank int, primary key( aspect, trait ));\n" +
		"/* stored tests, a work in progress \n" +
		" * ex. determine if \"sources\" should be listed in model ( for debugging )\n" +
		" */ \n" +
		"create table mdl_check( name text, type text, expect text, primary key( name ));\n" +
		"/* hierarchy of domains */\n" +
		"create table mdl_domain( domain text, path text, at text, primary key( domain ));\n" +
		"/* properties of a kind. type is a PRIM_ */\n" +
		"create table mdl_field( domain text, kind text, field text, affinity text, type text, at text, primary key( domain, kind, field ));\n" +
		"/* statements for user input parsing */\n" +
		"create table mdl_grammar(name text, prog blob, at text, primary key( name ));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind( domain text, kind text, path text, at text, primary key( domain, kind ));\n" +
		"/* initialization for pattern local variables */ \n" +
		"create table  mdl_local( domain text, kind text, field text, value blob, primary key( domain, kind, field ) );\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name( domain text, noun text, name text, rank int, at text );\n" +
		"/* a person, place, or thing in the world. \n" +
		" * domain tells the scope in which the noun was defined\n" +
		" * ( the same as - or a child of - the domain of the kind )\n" +
		" */\n" +
		"create table mdl_noun( domain text, noun text, kind text, at text, primary key( domain, noun ));\n" +
		"/* relation between two specific nouns. these change over the course of a game. */\n" +
		"create table mdl_pair( domain text, noun text, relation text, otherNoun text, at text );\n" +
		"/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */\n" +
		"create table mdl_pat( domain text, kind text, labels text, result text, primary key( domain, kind ) );\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( domain text, many text, one text, at text, primary key( domain, many ) );\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel( domain text, relation text, kind text, cardinality text, otherKind text, at text, primary key( domain, relation ));\n" +
		"/*  */\n" +
		"create table mdl_rule( domain text, pattern text, phase int, filter blob, prog blob, at text );\n" +
		"/* initial values for various nouns; changed values are stored in run_start */\n" +
		"create table mdl_value(domain text, noun text, field text, value blob, at text, primary key( domain, noun, field ));\n" +
		""
	return tmpl
}
