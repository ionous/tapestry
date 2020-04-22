/* alternative name for a noun */
create table eph_alias( idNamedAlias int, idNamedActual int );
/* collection of related states ( called traits ) */
create table eph_aspect( idNamedAspect int );
/* likelyhood that a trait applies to a particular kind */
create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );
/* test programs and the results they are expected to produce */
create table eph_check( idNamedTest int, idProg int, expect text );
/* initial values for the properties of nouns of belonging to the specified kind */
create table eph_default( idNamedKind int, idNamedProp int, value blob );
/* collection of related nouns */
create table eph_kind( idNamedKind int, idNamedParent int );
/* user specified appellation and the location that specification came from */
create table eph_named( name text, category text, idSource int, offset text );
/* a named object in the game world */
create table eph_noun( idNamedNoun int, idNamedKind int );
/* rule for the collective name of a singular word */
create table eph_plural( idNamedPlural int, idNamedSingluar int );
/* property name and type associated with a kind of object */
create table eph_field( primType text, idNamedKind int, idNamedField int );
/* connection between two kinds of object */
create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));
/* connection between two object instances */	
create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );
/* uri, file name or other identification for the origin of the various ephemera recorded in the db. */
create table eph_source( src text );
/* only one trait from a given aspect can be true for a noun at a time. */	
create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );
/* initial value for a noun's field, trait, or aspect */
create table eph_value( idNamedNoun int, idNamedProp int, value blob );
/* word indicating a particular relationship between nouns */
create table eph_verb( idNamedStem int, idNamedRelation int, verb text );
/* type is the name of the command container for de-serialization of the prog */
create table eph_prog( idSource int, type text, prog blob );
/* patternType is one of the table primTypes */
create table eph_pattern( idNamedPattern int, patternType text );
/* parameters for a pattern. */
create table eph_pattern_prim( idNamedPattern int, idNamedParam int, evalType text );
/* parameters for a pattern. named kind is a plural kind. */
create table eph_pattern_kind( idNamedPattern int, idNamedParam int, idNamedKind int );
/* function handler for a pattern */
create table eph_filter( idNamedPattern int, idProg int );