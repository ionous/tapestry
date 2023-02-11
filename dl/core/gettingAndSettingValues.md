# Getting and Setting Values 

Because (global) objects are such a key feature of interactive stories, for every variable statement there is a corresponding object statement. ( They have separate naming scopes, and nothing stops an author from having a local variable with the same name as an object. )

The reason that explicit variable names ( and object field names ) are required rather than leveraging dot against the local scope ( and object member names ) -- for example: GetObj{ "object", AtField{"member"} } is mainly because i think being explicit is good, and because that would also allow addressing the local scope ( or object members ) by index: and that doesn't make much sense ( ex. GetObj{ "object", AtIndex{5} } ).

To simplify setting values, there's a common "AssignedValue" that can swap between all of the common inputs: bool, number, text, list, and record. Currently, swaps don't support swapping between slots, so each eval has been wrapped in a flow.

## List handling

Where lists produce a value ( ex. reduce ) the result is often stored into a variable ( rather than producing a value. ) This is legacy and maybe should change. (An issue currently is having to write back object values to update them. )


## Auto conversion 

The goal should be to auto-convert between different types of value as best as possible; warning the author on mismatch. For example: int -> bool, etc.

TBD: maybe some auto-conversions list->bool ( is empty ) should be allowed and encouraged?


## Implicit indices 

For now at least, lists must be explicitly appended to. TBD: about sparse lists, sparse list storage and iteration, ...