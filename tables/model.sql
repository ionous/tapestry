/**
 * tables describing the world and its rules.
 */

/* 
 * stored tests, which run a program to verify it produces the expected value. 
 * fix: shouldnt we also be writing class of the value ? 
 */ 
create table mdl_check( domain text not null, name text, value blob, affinity text, prog blob, at text, primary key( domain, name ));
/* 
 * pairs of domain name and (domain) dependencies. 
 * domain names are considered globally unique.
 * a domain can have multiple direct parents; 
 * the application is responsible for ensuring no cyclic dependencies.
 */
create table mdl_domain( domain text not null, requires text, at text, primary key( domain, requires ));
/* 
 * arbitrary key-value information about the game world.
 * theoretically, they could be used by macros to affect the weave, or by the runtime to affect gameplay;
 * currently, they exist to detect "semantic" conflicts:
 * ex. an in-game password specified as "secret" in one place, and "mongoose" some place else.
 * the value is always a string right now, potentially could be expanded to a literal or assignment.
 */ 
create table mdl_fact( domain text not null, fact text, value text, at text, primary key( domain, fact ));
/* 
 * properties for a kind. 
 * type is most often used for affinities of type "text", and usually indicates a kind ( from mdl_kind )
 * currently, the name of the field must be unique withing the kind across all domains.
 * ( at runtime, the scope of the field is is considered the same as its kind;
 *   ie. all possible fields for a kind exist at once. )
 */
create table mdl_field( domain text not null, kind int not null, field text, affinity text, type int, at text, primary key( kind, field ));
/* 
 * statements for user input parsing. this is pretty low-bar right now;
 * typed commands aren't separated into unique rows, so conflicts between words and phrases can't be detected.
 */
create table mdl_grammar( domain text not null, name text, prog blob, at text, primary key( domain, name ));
/* 
 * a class of objects with shared characteristics 
 */
create table mdl_kind( domain text not null, kind text, singular text, path text, at text, primary key( domain, kind ));
/* 
 * words which refer to nouns. in cases where two words may refer to the same noun, 
 * the lower rank of the association wins. 
 * in theory, the scope of a name can be narrower that the scope of its noun 
 */
create table mdl_name( domain text not null, noun int not null, name text, rank int, at text );
/* 
 * a person, place, or thing in the world. 
 * domain tells the scope in which the noun was defined.
 * ( the same as - or a child of - the domain of the kind ) 
 */
create table mdl_noun( domain text not null, noun text, kind int not null, at text, primary key( domain, noun ));
/* 
 * relation between two specific nouns. these change over the course of a game.
 * similar to rulese, points back to the relation's kind rather than the entry in the mdl_rel table 
 */
create table mdl_pair( domain text not null, relKind int not null, oneNoun int not null, otherNoun int not null, at text );
/* 
 * pattern, the field ( in md_field ) used for a return value ( if any ) and comma separated labels for calling/processing fields 
 * the scope of the pattern is the same as its kind 
 */
create table mdl_pat( kind int not null, labels text, result text, primary key( kind ) );
/* 
 * maps common and uncommon words to their plurals. 
 * within a particular domain, a plural word produces a unique singular word;
 * however the same singular word can be used by various plurals.
 * ex. "person" can have the plural "people" or "persons".
 * but the *singular* of "people" must always be "person".
 */
create table mdl_plural( domain text not null, many text, one text, at text, primary key( domain, many ) );
/* 
 * statements for simplifying some kinds of story definitions.
 */
create table mdl_phrase( domain text not null, macro int not null, phrase text, reversed bool, at text, primary key( domain, phrase ));
/* 
 * relation and constraint between two kinds of nouns 
 */
create table mdl_rel( relKind int not null, oneKind int not null, otherKind int not null, cardinality text, at text, primary key( relKind ));
/* 
 * opposites 
 */
create table mdl_rev( domain text not null, oneWord text, otherWord text, at text );
/* 
 * the rules for a given kind of pattern within the specified domain are executed in increase rank, 
 * and within each rank, by last declared rule first ( largest row id. )
 + stop and jump describe the default handling of rules following a match:
 * stop controls whether processing terminates (1) or whether it moves to the next phase (0); 
 * jump controls when that transition takes place:
 *   2: transitions after processing all matching rules;
 *   1: transitions after processing the current element;
 *   0: transitions immediately;
 * there are various ways to change the desired stop/jump for an individual rule at runtime;
 * ( ex. setting the cancel flag; using the continue keyword; ... )
 * however the overall stop/jump for a set of rules can only stay the same or trigger sooner.
 * prog is a slice of rt.Execute statements.
 * updates hints whether there are any counters in the prog that need updating.
 */
create table mdl_rule( domain text not null, kind int not null, name text, rank int, stop int, jump int, updates int, prog blob, at text );
/* 
 * initial values for fields.
 * noun might be 0; in which case its a default value for the entire kind.
 * the scope of the value is the same as the noun or kind; that is the values are global: the same across all scenes.
 * dot contains sub field names separated by full stops.
 * provisional values can be refined during weave.
 */
create table mdl_value( noun int not null, field int not null, dot string, value blob, provisional int, at text, primary key( noun, field, dot ));
