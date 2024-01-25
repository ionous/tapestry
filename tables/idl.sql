/**
 * tables describing the tapestry commands.
 */

/* a tapestry command
 * type: one of, str, flow, num, swap.
 * ops are currently expected to be globally unique ( unlike golang where names are scoped per package )
 * to do otherwise, the .ifspec(s) themselves would have to contain package disambiguation when they name a type.
 */
create table idl_op( name text, package text, spec text, 
    primary key(name) );

/* permissible formats for each command. slot is a reference to an op of slot type.
 * signatures only have to be unique within the scope of each slot.
 * we allow slot to be NULL for concrete types ( as opposed to those that implement slot interfaces )
 * hash is stored as hex text to make the golang sql driver happy - it fails on uint64 with the highbit set.
 */
create table idl_sig( op int not null, slot int, hash text, body text,
    primary key(slot, body),
    unique(hash) );

/**
 * the predefined values for str types.
 */
create table idl_enum( op int not null, value text );

/** the members of a flow type; 'type' is an op reference */
create table idl_term( op int not null, name text, label text, type int,
    private bool, optional bool, repeats bool,
    primary key(op, name) );

/** markup from the serialized data; especially comments */
create table idl_markup( op int not null, key text, value blob,
    primary key(op, key) );

/** the pairings of ops to slots can be determined from the signature data */
create view idl_slot as select op,slot from idl_sig where slot is not null group by op, slot;
