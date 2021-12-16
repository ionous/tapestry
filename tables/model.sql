/**
 * tables describing the the world and its rules.
 */

/* enumerated values used by kinds and nouns.
 * fix: the data is duplicated in kinds and fields */
create table mdl_aspect( domain int, aspect int, trait text, rank int, primary key( aspect, trait ));
/* stored tests, which run a program to verify it produces the expected output. */ 
create table mdl_check( domain int, name text, expect blob, prog blob, at text, primary key( domain, name ));
/* hierarchy of domains */
create table mdl_domain( domain text, path text, at text, primary key( domain ));
/* properties of a kind. type is a PRIM_ */
create table mdl_field( domain int, kind int, field text, affinity text, type int, at text, primary key( domain, kind, field ));
/* statements for user input parsing */
create table mdl_grammar( domain int, name text, prog blob, at text, primary key( name ));
/* a class of objects with shared characteristics */
create table mdl_kind( domain int, kind text, path text, at text, primary key( domain, kind ));
/* initialization for pattern local variables
 * the initialization value is defined when the field is defined, so there's no separate origin at here. */ 
create table  mdl_local( domain int, kind int, field int, value blob, primary key( domain, kind, field ) );
/* words which refer to nouns. in cases where two words may refer to the same noun, 
   the lower rank of the association wins. */
create table mdl_name( domain int, noun int, name int, rank int, at text );
/* a person, place, or thing in the world. 
 * domain tells the scope in which the noun was defined
 * ( the same as - or a child of - the domain of the kind ) */
create table mdl_noun( domain int, noun text, kind text, at text, primary key( domain, noun ));
/* relation between two specific nouns. these change over the course of a game. */
create table mdl_pair( domain int, noun int, relation int, otherNoun int, at text );
/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */
create table mdl_pat( domain int, kind int, labels text, result text, primary key( domain, kind ) );
/* maps common and uncommon words to their plurals */
create table mdl_plural( domain int, many text, one text, at text, primary key( domain, many ) );
/* relation and constraint between two kinds of nouns
 * fix: the data is duplicated in kinds and fields... should this be removed? */
create table mdl_rel( domain int, rel text, kind int, cardinality text, otherKind int, at text, primary key( domain, rel ));
/*  note: the "pattern" id is actually a reference to the core kind */
create table mdl_rule( domain int, kind int, target int, phase int, filter blob, prog blob, at text );
/* initial values for various nouns; changed values are stored in run_start */
create table mdl_value( domain int, noun int, field int, value blob, at text, primary key( domain, noun, field ));
