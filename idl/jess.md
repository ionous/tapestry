
Unexpected success
------------------
There are sentences that don't make sense in English, that jess will read perfectly fine. And others that mean something in English slightly different than how they are understood here.

For instance:

* `Some containers are usually closed.`: actually means are containers start closed by default.
* `The box is a container. The box is a closed.`: jess ignores the 'a' in front of the trait.
* `The box is a transparent and container.`: because phrases like: "The box is transparent and a container" are supported, it permits some other less than sensible phrases as well.
* `The container called the trunk and the box is in the lobby.`: creates one single container with an unusual name.
* `The containers called boxes are a kind of closed thing.`: 

Inform has these same issues.

## Some known differences from Inform

1. A box is a kind of closed and transparent container.
2. A box is a kind of closed transparent container.
3. A box is a kind of container and thing.

The first two are okay with jess, the third one is not; 
In inform <sup>1.68.1</sup>, the first and third are okay.

* `Buckets and baskets are kinds of container.` in inform becomes two kinds: "Bucketss" and "basketss".

 
# regarding traits there are a a few patterns:

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

