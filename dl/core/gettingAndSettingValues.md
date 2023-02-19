# Getting and Setting Values 

Because (global) objects are such a key feature of interactive stories, for every variable statement there is a corresponding object statement. ex: GetFromVar, GetFromObj.

( Note that hey have separate naming scopes, and nothing stops an author from having a local variable with the same name as an object. )

To simplify setting values, there's a common "AssignedValue" that can swap between all of the common inputs: bool, number, text, list, and record. ( Currently, swaps don't support swapping between slots, so each eval has been wrapped in a flow. )

Lists and records are values not objects, they are not shared. Every variable ( and every object or record field ) is unique. Assign always means copy.

## List handling

List mutation operates on variables; list manipulation on evals; indexing relies on the dot functions. TBD: maybe mutation could be expanded to a more generic reference ( get obj/ variable would also incorporate reference then. )

## Auto conversion 

The goal should be to auto-convert between different types of value as best as possible; warning the author on mismatch. For example: int -> bool, etc.

TBD: maybe some auto-conversions list->bool ( is empty ) should be allowed and encouraged?

## Implicit indices 

For now at least, lists must be explicitly appended to. TBD: about sparse lists, sparse list storage and iteration, ...

# Expressions

Object fields an be given expressions in story declarations, new expressions can't be assigned to object fields at runtime. Similarly, records cannot hold expressions. Expressions are not values, they only produce values. ( Some future version of Tapestry could perhaps grow to support this. )


# Object vs kinds 

Ideally, i think objects would be singletons of a kind and interchangeable with functions that talk about kinds. Tapestry isn't there yet.


# FAQ

### Why not use dot ( AtField, and AtIndex ) directly against current scope? Why do they only exist as members of other commands?

For example, why not:
```
	var msg = "hello world";
	Print{AtField{msg}}
```
or
```
	Print{AtField{"color", AtObj{"vase"}}
```

mainly because that would also allow addressing the local scope ( or object members ) by index: and that doesn't make much sense ( ex. `Print{AtIndex{5}}` )

### Why are object references not a first class type? Why does it use "text"? 

This is somewhat a legacy of interactive fiction: authors expect to be able to refer to objects by name. And since names are just text, then why have some sort of special text-like type that refers to objects when you could just use text. ( Granted, this purity of idea is somewhat diminished by the fact that object ids get stored with a special string syntax. )

### Why is object not its own type ( the way bool, number, text, list, and record are )? 

Exposing an object value would imply that object values -- like all of the other values -- could be stored directly in records. But, objects are globally unique. Storing references to them makes sense, storing the actual objects less so.

### Why aren't records ( and lists ) first class objects? If they were, couldn't objects then just be records with user-defined names?

In my opinion, its easier for people -- especially for non-programmers -- to conceive of variables as completely unique values than to consider the idea that two different variables might point to the same value. 

However, people of all technical levels get by in javascript and lua just fine; so maybe it's _more_ confusing that shared values don't exist. ( And certainly, the current naive version -- which doesn't identify temporaries, doesn't have copy-on-write, etc. -- is much less efficient than otherwise. )

There might be some wiggle room for revisiting this in a future edition of Tapestry.