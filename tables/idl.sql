/**
 * tables describing the tapestry commands.
 */

/* a tapestry command
 * uses: one of, str, flow, num, swap, group
 * closed is for str and num types indicating whether the user is allowed custom values.
 */
create table idl_op( name text, package text, uses text, closed bool, primary key(name, package) );

/* permissable formats for each command. slot is a reference to an op of slot type.
signatures only have to be unique within the scope of each slot.
 */
create table idl_sig( op int not null, slot int not null, signature text, primary key(slot, signature) );

/** for str types with predefined values */
create table idl_str( op int not null, label text, value text, primary key(op, label) );

/** for num types with predefined values
future: maybe a separate range min, max table
 */
create table idl_enum( op int not null, label text, value number, primary key(op, label) );

/** the choices for a swap op. type is an op reference */
create table idl_swap( op int not null, label text, type int, primary key(op, label) );

/** the members of a flow op. type is an op reference. */
create table idl_term( op int not null, term text, label text, type int,
    private bool optional bool, repeats bool,
    primary key( op, term ) );

/** markup from the serialized data; especially comments */
create table idl_markup( op int not null, term text, markup text, value blob,
    primary key( op, term, markup ) );



