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
		"create table mdl_aspect( aspect text, trait text, rank int, primary key( aspect, trait ));\n" +
		"/* stored tests, a work in progress \n" +
		" * ex. determine if \"sources\" should be listed in model ( for debugging )\n" +
		" */ \n" +
		"create table mdl_check( name text, type text, expect text );\n" +
		"/* hierarchy of domains */\n" +
		"create table mdl_domain( domain text, path text, primary key( domain ));\n" +
		"/* properties of a kind. type is a PRIM_ */\n" +
		"create table mdl_field( kind text, field text, type text, affinity text, primary key( kind, field ));\n" +
		"/* a class of objects with shared characteristics */\n" +
		"create table mdl_kind( kind text, path text, primary key( kind ));\n" +
		"/* words which refer to nouns. in cases where two words may refer to the same noun, \n" +
		"   the lower rank of the association wins. */\n" +
		"create table mdl_name( noun text, name text, rank int );\n" +
		"/* a person, place, or thing in the world. */\n" +
		"create table mdl_noun( noun text, kind text, primary key( noun ));\n" +
		"/* relation between two specific nouns. these change over the course of a game. */\n" +
		"create table mdl_pair( noun text, relation text, otherNoun text, domain text );\n" +
		"/* maps common and uncommon words to their plurals */\n" +
		"create table mdl_plural( one text, many text );\n" +
		"/* stored programs, a work in progress \n" +
		"   the connection between tests and patterns and progs are, essentially, application knowledge. */ \n" +
		"create table mdl_prog( name text, type text, bytes blob );\n" +
		"/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */\n" +
		"create table mdl_pat( name text, labels text, result text, primary key( name ) );\n" +
		"/* relation and constraint between two kinds of nouns */\n" +
		"create table mdl_rel( relation text, kind text, cardinality text, otherKind text, primary key( relation ));\n" +
		"/* \n" +
		"\towner: a pattern or event name \n" +
		"\tdomain: the scene when the rule applies ( could potentially live in another table if we wanted to reuse rule across domains )\n" +
		"\tscope: for events, the noun or kind to which this rule applies; patterns are called directly and dont have scope.\n" +
		"\trule name: unique if it exists\n" +
		"\tprog: a rt.Handler containing an optional filter and executable statement\n" +
		" */\n" +
		"create table mdl_rule( owner text, scope text, domain text, phase int, prog blob, name text unique );\n" +
		"/* documentation for pieces of the model: kinds, nouns, fields, etc. */\n" +
		"create table mdl_spec( type text, name text, spec text, primary key( type, name ));\n" +
		"/* initial values for various noun, kind, or pattern properties. \n" +
		"   changed values are stored in run_start.. */\n" +
		"create table mdl_start( owner text, field text, value blob );"
	return tmpl
}
