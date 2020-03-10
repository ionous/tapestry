create table eph_alias( idNamedAlias int, idNamedActual int );
create table eph_aspect( idNamedAspect int );
create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );
create table eph_check( idNamedTest text, idProg int, expect text );
create table eph_default( idNamedKind int, idNamedProp int, value blob );
create table eph_kind( idNamedKind int, idNamedParent int );
create table eph_named( name text, category text, idSource int, offset text );
create table eph_noun( idNamedNoun int, idNamedKind int );
create table eph_plural( idNamedPlural int, idNamedSingluar int );
create table eph_primitive( primType text, idNamedKind int, idNamedField int );
create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));
create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );
create table eph_source( src text );
create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );
create table eph_value( idNamedNoun int, idNamedProp int, value blob );
create table eph_verb( idNamedStem int, idNamedRelation int, verb text );
create table eph_prog( idSource int, type text, prog blob );


