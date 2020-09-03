package object

// reserved fields
const Name = "name" // name of an object as declared by the user

// internal fields
const Prefix = '$'         // leading character used for all internal fields
const Id = "$id"           // unique identifier for an object, includes its home domain
const Exists = "$exists"   // whether a name refers to a declared game object
const Kind = "$kind"       // type of a game object
const Kinds = "$kinds"     // hierarchy of a game object ( a path )
const Counter = "$counter" // sequence counter
const Aspect = "$aspect"   // name of aspect for noun.trait

// asking for the "rule" field of a named pattern
// returns an aggregated list of "programs" of that type.
const BoolRule = "$bool_rule"
const NumberRule = "$number_rule"
const TextRule = "$text_rule"
const ExecuteRule = "$execute_rule"
const NumListRule = "$num_list_rule"
const TextListRule = "$text_list_rule"
