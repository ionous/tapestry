# Getting and Setting Values 

Because (global) objects are such a key feature of interactive stories, for every variable statement there is a corresponding object statement. ex: GetFromVar, GetFromObj.

( They have separate naming scopes, and nothing stops an author from having a local variable with the same name as an object. )

To simplify setting values, there's a common "AssignedValue" that can swap between all of the common inputs: bool, number, text, list, and record. 

Lists and records are values, not objects, and they're values aren't shared across variables. The value of every variable ( including the value held by every field of an object or record ) is unique. Assign means copy.

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

All values are unique and assignment means copy. Objects on the other hand are globally unique. Exposing an object as a value would contradict one of those two foundations. Storing references to objects as values makes sense, storing the actual objects as values less so.

### Why aren't records ( and lists ) first class objects? ( And if they were, could objects be records with user-defined names? )

In my opinion it's easier for people -- especially for non-programmers -- to conceive of variables as completely unique values than to consider that two different variables might refer to the same value. Keeping records and lists as unique values supports that idea, however it does come with some expense. The current runtime -- which doesn't identify temporaries, doesn't have copy-on-write, etc. -- is much less efficient than if values could be shared. ( There might be some wiggle room for revisiting this in some future edition of Tapestry. )