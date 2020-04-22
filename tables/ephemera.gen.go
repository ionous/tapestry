/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// ephemeraTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func ephemeraTemplate() string {
	var tmpl = "/* alternative name for a noun */\n" +
		"create table eph_alias( idNamedAlias int, idNamedActual int );\n" +
		"/* collection of related states ( called traits ) */\n" +
		"create table eph_aspect( idNamedAspect int );\n" +
		"/* likelyhood that a trait applies to a particular kind */\n" +
		"create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );\n" +
		"/* test programs and the results they are expected to produce */\n" +
		"create table eph_check( idNamedTest int, idProg int, expect text );\n" +
		"/* initial values for the properties of nouns of belonging to the specified kind */\n" +
		"create table eph_default( idNamedKind int, idNamedProp int, value blob );\n" +
		"/* collection of related nouns */\n" +
		"create table eph_kind( idNamedKind int, idNamedParent int );\n" +
		"/* user specified appellation and the location that specification came from */\n" +
		"create table eph_named( name text, category text, idSource int, offset text );\n" +
		"/* a named object in the game world */\n" +
		"create table eph_noun( idNamedNoun int, idNamedKind int );\n" +
		"/* rule for the collective name of a singular word */\n" +
		"create table eph_plural( idNamedPlural int, idNamedSingluar int );\n" +
		"/* property name and type associated with a kind of object */\n" +
		"create table eph_field( primType text, idNamedKind int, idNamedField int );\n" +
		"/* connection between two kinds of object */\n" +
		"create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));\n" +
		"/* connection between two object instances */\t\n" +
		"create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );\n" +
		"/* uri, file name or other identification for the origin of the various ephemera recorded in the db. */\n" +
		"create table eph_source( src text );\n" +
		"/* only one trait from a given aspect can be true for a noun at a time. */\t\n" +
		"create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );\n" +
		"/* initial value for a noun's field, trait, or aspect */\n" +
		"create table eph_value( idNamedNoun int, idNamedProp int, value blob );\n" +
		"/* word indicating a particular relationship between nouns */\n" +
		"create table eph_verb( idNamedStem int, idNamedRelation int, verb text );\n" +
		"/* type is the name of the command container for de-serialization of the prog */\n" +
		"create table eph_prog( idSource int, type text, prog blob );\n" +
		"/* patternType is one of the table primTypes */\n" +
		"create table eph_pattern( idNamedPattern int, patternType text );\n" +
		"/* parameters for a pattern. primType is one of the core primTypes ( text, digi, prog ) */\n" +
		"create table eph_pattern_eval( idNamedPattern int, idNamedParam int, primType text );\n" +
		"/* parameters for a pattern. named kind is a plural kind. */\n" +
		"create table eph_pattern_kind( idNamedPattern int, idNamedParam int, idNamedKind int );\n" +
		"/* function handler for a pattern */\n" +
		"create table eph_filter( idNamedPattern int, idProg int );"
	return tmpl
}
