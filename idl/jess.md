
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

Unexpected failures
-------------------

`Three things are fixed in place in the kitchen` will not generate three anonymous things; instead it will generate an object called "three things." Instead, adjectives for anonymous objects must sit between the number and the the kind. For example: `Three fixed in place things are in the kitchen.`

Inform works the same way.

## Some known differences from Inform

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

