/**
 * tables describing the the world and its rules.
 */

/* enumerated values used by kinds and nouns */
create table mdl_aspect( aspect text, trait text, rank int, primary key( aspect, trait ));
/* stored tests, a work in progress 
 * ex. determine if "sources" should be listed in model ( for debugging )
 */ 
create table mdl_check( name text, type text, expect text );
/* hierarchy of domains */
create table mdl_domain( domain text, path text, primary key( domain ));
/* properties of a kind. type is a PRIM_ */
create table mdl_field( kind text, field text, type text, affinity text, primary key( kind, field ));
/* a class of objects with shared characteristics */
create table mdl_kind( kind text, path text, primary key( kind ));
/* words which refer to nouns. in cases where two words may refer to the same noun, 
   the lower rank of the association wins. */
create table mdl_name( noun text, name text, rank int );
/* a person, place, or thing in the world. */
create table mdl_noun( noun text, kind text, primary key( noun ));
/* relation between two specific nouns. these change over the course of a game. */
create table mdl_pair( noun text, relation text, otherNoun text, domain text );
/* maps common and uncommon words to their plurals */
create table mdl_plural( one text, many text );
/* stored programs, a work in progress 
   the connection between tests and patterns and progs are, essentially, application knowledge. */ 
create table mdl_prog( name text, type text, bytes blob );
/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields */
create table mdl_pat( name text, labels text, result text, primary key( name ) );
/* relation and constraint between two kinds of nouns */
create table mdl_rel( relation text, kind text, cardinality text, otherKind text, primary key( relation ));
/* 
	owner: a pattern or event name 
	target: for events, the noun or kind to which this rule applies; patterns dont have target: they are called directly.
	domain: the scene when the rule applies ( could potentially live in another table if we wanted to reuse rule across domains )
	rule name: unique if it exists
	prog: a rt.Handler containing an optional filter and executable statement
 */
create table mdl_rule( owner text, target text, domain text, phase int, prog blob, name text unique );
/* documentation for pieces of the model: kinds, nouns, fields, etc. */
create table mdl_spec( type text, name text, spec text, primary key( type, name ));
/* initial values for various noun, kind, or pattern properties. 
   changed values are stored in run_start.. */
create table mdl_start( owner text, field text, value blob );