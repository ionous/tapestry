/*
Package parser

The parser requires a grammar built of Scanner commands.

Grammar

The parser contains a tree of Scanners. A successful match returns an implementation of the "Results" interface. Branching scanners generally produce a "ResultsList". Terminal nodes each have a matching "Results" type.

Branching scanners:

	- AllOf matches the passed matchers in order.
	- AnyOf matches any one of the passed Scanners; whichever first matches.
	- Focus changes the bounds for subsequent scanners. For instance, searching only though held objects.
	- Target changes the bounds of its first scanner in response to the results of its last scanner.
	  Generally, this means that the last scanner should be Noun{}.

Terminal scanners:

	- Action terminates a matcher sequence, resolving to the named action.
	returns ResolvedAction

	- Multi matches one or more objects.
	returns ResolvedMulti

	- Noun matches one object held by the context.
	returns ResolvedNoun

	- Word matches one word.
	returns ResolvedWords{1}

Context

Scanners read from the world model using "Context".

	- test for plurals
	- the set of objects in reach of the player
	- the set of objects available to another object. ( ex. inside or on )

Results

The result of a successful parsing is most often a ResultList,
the .Last element() of which is usually an action.

*/
package parser
