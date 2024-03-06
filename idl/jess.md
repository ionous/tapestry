
Unexpected success
------------------
There are sentences that don't make sense in English, that jess will read perfectly fine. And others that mean something in English slightly different than how they are understood here.

For instance:

* `Some containers are usually closed.`: actually means are containers start closed by default.
* `The box is a container. The box is a closed.`: jess ignores the 'a' in front of the trait.
* `The box is a transparent and container.`: because phrases like: "The box is transparent and a container" are supported, it permits some other less than sensible phrases as well.
* `The container called the trunk and the box is in the lobby.`: creates one single container with an unusual name.
* `The containers called boxes are a kind of closed thing.`: 
* `In a container are five things.`  is allowed, makes sense to English and to jess but is useless. The created objects are out of world, and are unaddressable.
* `The five the containers are in the kitchen.` is permitted.

Inform has these same issues.

Article edge cases 
-----
This is funny, if only because its accidentally the same in both jess and inform. 

Article determination differs based on capitalization and location in the sentence, so a sentence like "The Kitchen is a room" recognizes "The" is an article. However, that's also true if it was written "THE", "tHE", or "the".

As expected, for "A thing is in the Kitchen." 'the' is recognized as an article, whereas for "A thing is in The Kitchen." it becomes part of the name. Here, "THE" becomes "THE Kitchen" while "tHE"? it's seen as a normal article.

Categorization depends on the first letter only.


Unexpected failures
-------------------

`Three things are fixed in place in the kitchen` will not generate three anonymous things; instead it will generate an object called "three things." Instead, adjectives for anonymous objects must sit between the number and the the kind. For example: `Three fixed in place things are in the kitchen.`

Inform works the same way.

Differences from Inform
------------------------

## Functional differences

Inform allows limiting traits to kinds with other traits:
For example, `the closed containers are fixed in place.`makes any containers that are *initially* closed also immovable. Jess is not that clever, and instead (theoretically) generates an error.

## Syntax differences

1. A box is a kind of closed and transparent container.
2. A box is a kind of closed transparent container.
3. A box is a kind of container and thing.

The first two are okay with jess, the third one is not; 
In inform <sup>1.68.1</sup>, the first and third are okay.

* `Buckets and baskets are kinds of container.` in inform becomes two kinds: "Bucketss" and "basketss". The extra s isnt present with jess.

Inform is looser with its article matching in some cases:

*  `The the closed box called the Pandorica is open.` and `The a a a openable container called the vase is in the kitchen.` are valid in inform for reasons i don't entirely understand.



 
# regarding traits there are a few patterns:

1. "(det) traits kind are verbing names."
	ex. `Five closed containers are in the kitchen.`
	this is a specific phrase.
	it combines the traits and the kind and leaves them together.
			

1.  "... traits kind called name ..."
	ex `A closed openable container called the trunk is in the lobby.`
	it takes the traits and kind and applies it to a singular name.


1. "the names are traits ( . | verb nouns | kinds )"
	ex. `the box and the trunk are closed containers.`
			`a casket is a kind of closed container.`
	    `the box is closed (in the kitchen).`
	  
	 when you see a kind, or the end of the sentence;
	 take the accumulated traits and add them to the current target.

1. ( kinds are "usually" traits ): this is a specific phrase
xxx: ( traits noun ) isnt allowed ( ex. the unhappy man )



# Noun Value Assignment

This is a phrase that could really use some love.

some notes from inform:
* there are two phrase variants: 
	- property noun value: 
	  - `The age of the bottle is 42.`
	  
	- noun property value: 
		- `The bottle is/are/has age 42.`
		- `The container called the box has the description ...`
	
* unknown nouns are allowed:
  - ok: `The description of the nep is "mightier than the sword."`
  
* supports multiple values for the same noun:
  - ok: `The bottle in the kitchen is openable and has age 42.`
  
* "called the" is only supported in leading phrases:
	- ok: `The container called the box has the description ... `
	- ng: `The description of the container called the box is....`
	
* doesnt support multiple nouns:
	- ng: `The description of the box and the top is...`
	- ng: `The description of the pen is "mightier than the sword" and the description of the box is "hello."`
	
* undesirable results when using verb phrases:
 - ng: `The age of the bottle in the kitchen is 42.` - generates an object called "bottle in the kitchen".
 - ng: `The bottle in the kitchen has the age 42.` - ditto.
 - ng: `The bottle is fixed in place and has age 42.` -- generates an object called "bottle is fixed in place".
 
 
 # Understandings
 
inform supports defining plurals, aliases, command synonyms, actions, references for actions,  trait, kind, and value substitutions....and possibly some others.
 
 In inform, it seems aliases (even for kinds) are applied to the noun name, while plurals add to the rules of parsing that noun:
 
 ```
 		with name 'cardboard' 'box' 'containers//p' 
    with parse_name Parse_Name_GV93  ! = a routine
 ```
 
 For inform, since `Understand "uniqueword [something]" as looking under and box.` both generate parser rules.... if there's a scene described with `The cardboard box is an open container in the kitchen. The rose is in the kitchen.` then the user can say `> x uniqueword rose` and it will say `The cardboard box is empty.` -- which is pretty neat
 
 For `Understand "uniqueword" as the box.` it will add "uniqueword" as an alias to the noun definition -- `Object -> I126_cardboard_box... with name 'cardboard' 'box' 'containers//p' 'uniqueword'`. it analyzes the lhs for the brackets to determine if its grammar or not.
 

Notes on scheduling
-----------

As phrases matched, generation within a "paragraph" ( meaning a "Declare:" block ) are "sorted" so that first kinds are declared, then nouns, then values of nouns, and finally relationships between them. 

Key to the implementation is that when values are assigned to nouns, if the noun does not yet exist: it becomes a `thing.` ( That is except for nouns that were referenced via directional phrases. Those become `rooms` or `doors` depending on the context. )

For instance:

`A thing is either ended or started. The bother is an ended container.` works fine; the reverse does not. Same with `Two boggles are in the kitchen. A boggles is a kind of thing.` `The glass bottle is in the kitchen. The bottle is fixed in place.` creates one noun, the reverse creates two.

However:

`The bottle is in the kitchen. The bottle is a container. The kitchen is a room.` correctly creates one container and one room. And, `Understand "donut" as the doughnut. The doughnut is an animal.` works fine.

Inform doesn't delay processing of kinds. The following, perhaps unexpectedly, generates two objects: a thing called "five trees" and a tree called the "sapling."

```
The five trees are in the kitchen.
A tree is a kind of thing.
The sapling is a tree.
```