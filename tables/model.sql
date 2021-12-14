/**
 * tables describing the the world and its rules.
 */

/* enumerated values used by kinds and nouns */
create table mdl_aspect( domain text, aspect text, trait text, rank int, primary key( aspect, trait ));
/* stored tests, which run a program to verify it produces the expected output. */ 
create table mdl_check( domain text, name text, expect blob, prog blob, at text, primary key( domain, name ));
/* hierarchy of domains */
create table mdl_domain( domain text, path text, at text, primary key( domain ));
/* properties of a kind. type is a PRIM_ */
create table mdl_field( domain text, kind text, field text, affinity text, type text, at text, primary key( domain, kind, field ));
/* statements for user input parsing */
create table mdl_grammar(name text, prog blob, at text, primary key( name ));
/* a class of objects with shared characteristics */
create table mdl_kind( domain text, kind text, path text, at text, primary key( domain, kind ));
/* initialization for pattern local variables */ 
create table  mdl_local( domain text, kind text, field text, value blob, primary key( domain, kind, field ) );
/* words which refer to nouns. in cases where two words may refer to the same noun, 
   the lower rank of the association wins. */
create table mdl_name( domain text, noun text, name text, rank int, at text );
/* a person, place, or thing in the world. 
 * domain tells the scope in which the noun was defined
 * ( the same as - or a child of - the domain of the kind )
 */
create table mdl_noun( domain text, noun text, kind text, at text, primary key( domain, noun ));
/* relation between two specific nouns. these change over the course of a game. */
create table mdl_pair( domain text, noun text, relation text, otherNoun text, at text );
/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */
create table mdl_pat( domain text, kind text, labels text, result text, primary key( domain, kind ) );
/* maps common and uncommon words to their plurals */
create table mdl_plural( domain text, many text, one text, at text, primary key( domain, many ) );
/* relation and constraint between two kinds of nouns */
create table mdl_rel( domain text, relation text, kind text, cardinality text, otherKind text, at text, primary key( domain, relation ));
/*  */
create table mdl_rule( domain text, pattern text, target text, phase int, filter blob, prog blob, at text );
/* initial values for various nouns; changed values are stored in run_start */
create table mdl_value(domain text, noun text, field text, value blob, at text, primary key( domain, noun, field ));
