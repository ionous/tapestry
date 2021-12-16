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
		"/* enumerated values used by kinds and nouns.\n" +
		" * fix: the data is duplicated in kinds and fields */\n" +
		"create table mdl_aspect( domain int, aspect int, trait text, rank int, primary key( aspect, trait ));\n" +
		"/* stored tests, which run a program to verify it produces the expected value. \n" +
		" * fix: shouldnt we also be writing class of the value ? */ \n" +
		"create table mdl_check( domain int, name text, value blob, affinity text, prog blob, at text, primary key( domain, name ));\n" +
		"/* hierarchy of domains */\n" +
		"create table mdl_domain( domain text, path text, at text, primary key( domain ));\n" +
		"/* properties of a kind. type is a PRIM_ */\n" +
		"create table mdl_field( domain int, kind int, field text, affinity text, type int, at text, primary key( domain, kind, field ));\n" +
		"/* statements for user input parsing */\n" +
		"create table mdl_grammar( domain int, name text, prog blob, at text, primary key( name ));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind( domain int, kind text, path text, at text, primary key( domain, kind ));\n" +
		"/* initialization for pattern local variables\n" +
		" * the initialization value ( an Assignment ) is defined when the field is defined, so there's no separate origin at here. */ \n" +
		"create table  mdl_local( domain int, kind int, field int, assign blob, primary key( domain, kind, field ) );\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name( domain int, noun int, name int, rank int, at text );\n" +
		"/* a person, place, or thing in the world. \n" +
		" * domain tells the scope in which the noun was defined\n" +
		" * ( the same as - or a child of - the domain of the kind ) */\n" +
		"create table mdl_noun( domain int, noun text, kind text, at text, primary key( domain, noun ));\n" +
		"/* relation between two specific nouns. these change over the course of a game. */\n" +
		"create table mdl_pair( domain int, noun int, relation int, otherNoun int, at text );\n" +
		"/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */\n" +
		"create table mdl_pat( domain int, kind int, labels text, result text, primary key( domain, kind ) );\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( domain int, many text, one text, at text, primary key( domain, many ) );\n" +
		"/* relation and constraint between two kinds of nouns\n" +
		" * fix: the data is duplicated in kinds and fields... should this be removed? */\n" +
		"create table mdl_rel( domain int, rel text, kind int, cardinality text, otherKind int, at text, primary key( domain, rel ));\n" +
		"/*  note: the \"pattern\" id is actually a reference to the core kind */\n" +
		"create table mdl_rule( domain int, kind int, target int, phase int, filter blob, prog blob, at text );\n" +
		"/* initial values for various nouns; changed values are stored in run_start \n" +
		" * fix: shouldnt we also be writing class of the value ? */\n" +
		"create table mdl_value( domain int, noun int, field int, value blob, affinity text, at text, primary key( domain, noun, field ));\n" +
		""
	return tmpl
}
