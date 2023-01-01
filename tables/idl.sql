/**
 * tables describing the tapestry commands.
 */

/* a tapestry command
 * uses: one of, str, flow, num, swap, group
 * closed is for str and num types indicating whether the user is allowed custom values.
 * ops are currently expected to be globally unique ( unlike golang where names are scoped per package )
 * to do otherwise, the .ifspec(s) themselves would have to contain package disambiguation when they name a type.
 */
create table idl_op( name text, package text, uses text, closed bool, primary key(name) );

/* permissible formats for each command. slot is a reference to an op of slot type.
 * signatures only have to be unique within the scope of each slot.
 * we allow slot to be NULL for concrete types ( as opposed to those that implement slot interfaces )
 * hash is stored as hex text to make the golang sql driver happy - it fails on uint64 with the highbit set.
 */
create table idl_sig( op int not null, slot int, hash text, signature text, primary key(slot, signature) );

/**
 * the predefined values of str and num types.
 * unlike the raw specification where the label can be blank, its expanded here.
 * future: maybe a separate range min, max table for nums?
 */
create table idl_enum( op int not null, label text, value blob, primary key(op, label) );

/** the choices for a swap op. type is an op reference */
create table idl_swap( op int not null, label text, value text, type int, primary key(op, label) );

/** the members of a flow op. type is an op reference. */
create table idl_term( op int not null, field text, label text, type int,
    private bool, optional bool, repeats bool,
    primary key( op, field ) );

/** markup from the serialized data; especially commeRnts */
create table idl_markup( op int not null, term text, markup text, value blob,
    primary key( op, term, markup ) );



