/**
 * tables describing the the world and its rules.
 */

/* initialization for fields of kinds.
 * the value ( an Assignment ) is defined when the field is defined, 
 * meaning there's no separate domain or origin(at). */ 
create table mdl_assign( field int not null, value blob, primary key( field ) );
/* stored tests, which run a program to verify it produces the expected value. 
 * fix: shouldnt we also be writing class of the value ? */ 
create table mdl_check( domain int not null, name text, value blob, affinity text, prog blob, at text, primary key( domain, name ));
/* hierarchy of scopes */
create table mdl_domain( domain text, path text, at text, primary key( domain ));
/* properties for a kind. 
 * type is most often used for affinities of type "text", and usually indicates a kind ( from mdl_kind )
 * the scope of a field is the same as its kind. */
create table mdl_field( kind int not null, field text, affinity text, type int, at text, primary key( kind, field ));
/* statements for user input parsing */
create table mdl_grammar( domain int not null, name text, prog blob, at text, primary key( domain, name ));
/* a class of objects with shared characteristics */
create table mdl_kind( domain int not null, kind text, path text, at text, primary key( domain, kind ));
/* words which refer to nouns. in cases where two words may refer to the same noun, 
 * the lower rank of the association wins. 
 * in theory, the scope of a name can be narrower that the scope of its noun */
create table mdl_name( domain int not null, noun int not null, name text, rank int, at text );
/* a person, place, or thing in the world. 
 * domain tells the scope in which the noun was defined.
 * ( the same as - or a child of - the domain of the kind ) */
create table mdl_noun( domain int not null, noun text, kind int not null, at text, primary key( domain, noun ));
/* relation between two specific  nnouns. these change over the course of a game.
 * similar to mdl_rule, points back to the relation's kind rather than the entry in the mdl_rel table */
create table mdl_pair( domain int not null, relKind int not null, oneNoun int not null, otherNoun int not null, at text );
/* pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields 
 * the scope of the pattern is the same as its kind */
create table mdl_pat( kind int not null, labels text, result text, primary key( kind ) );
/* maps common and uncommon words to their plurals */
create table mdl_plural( domain int, many text, one text, at text, primary key( domain, many ) );
/* relation and constraint between two kinds of nouns */
create table mdl_rel( relKind int not null, oneKind int not null, otherKind int not null, cardinality text, at text, primary key( relKind ));
/* opposites */
create table mdl_rev( domain int, oneWord text, otherWord text, at text );
/* the scope of a rule can be narrower than its parent kind ( or target )
 * fix? can target (kind of nouns the rule applies to) be moved to filter? */
create table mdl_rule( domain int not null, kind int not null, target int, phase int, filter blob, prog blob, at text );
/* initial values for various nouns.
 * note: currently, the scope of the value is the same as the noun.
 * ( ie. one value per noun per field; not values that change as new domains are activated. )
 + the affinity and subtype of the value come from the field.
 */
create table mdl_value( noun int not null, field int not null, value blob, at text, primary key( noun, field ));
