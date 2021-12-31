# Tapestry

A story engine for narrative driven games. 

# Overview

Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This would enable quick iteration for commonplace story-oriented tasks such as creating dialog, puzzles, and quests.

Modern game engines, however, while providing great tools for the visual and arcade parts of a game -- anything requiring immediate player feedback -- don't have great tools for building *stories*. Most provide scripting support, but it's inevitably deeply integrated into the graphical engine: development while quick at the outset, often becomes error prone and slow.

So what if we *could* separate the second-by-second gameplay from the overarching game story? What if we could describe the game world quickly and easily *in words*, and then connect that world to the visual assets of the game? What if we could allow the development of story and graphics to proceed independently whenever possible? What would that process look like and how would it work?

Tapestry's goal is to answer those questions.

# Status

Tapestry is a work in progress. Nothing is probably usable by other people yet. During development, I'm committed to tracking tasks on [trello](https://trello.com/b/EEPnJ6ew/tapestry), writing updates on the [tapestry wiki](https://man.sr.ht/~ionous/tapestry/), and occasionally posting on twitter [@theionous](https://twitter.com/theionous).

The flow of story creation is:

1. Write stories using the web-based "composer" or raw script.
2. Import stories and other data (ex. from game assets.)
3. Assemble a game database which can provide a macro view of any game world.
4. Test and iterate on the story using an interactive console.
5. Embed the story engine into the game engine of your choice, connecting story objects to art assets, linking the logic of the game to the visuals.
6. Test and iterate some more :) ( using the game console and game database to track progression, inspect the game world, run unit tests, etc. )

Rough versions of everything except embedding the engine exist. Eventually, the goal is to provide easy integration with Unity, Unreal, etc.

Near term goals include:

* Adding story logic to support small interactive fiction like games.
* Improving the underlying modeling and scripting.
* Improving the composer to more easily write stories and tests.

The intention is to keep the story engine itself open source and [liberally licensed](https://man.sr.ht/~ionous/tapestry/LICENSE.md), while any 3rd party engine integration would likely be separately licensed.

This effort started as a re-implementation of the [Sashimi interactive fiction engine](https://github.com/ionous/sashimi) which was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a 2D point-and-click adventure game. I've used similar techniques for [various projects](https://www.linkedin.com/in/ionous/) I've worked on over the last 20 years, but Tapestry is its own creation based on lessons learned.

# Relation to interactive fiction

The ability to play and test a game without graphics is to my mind a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar level of interactivity with a similar set of game rules as Inform.

It is *not* a goal to attempt to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. And, it is, of course, a key goal to *extend* interactive fiction: bringing similar tools into any game genre.

# Building the code

Tapestry is written in [Go](https://golang.org/), and initially will use web assembly as a means of integration into existing engines. The runtime is small relative to the behind the scenes work, however, so porting to the c family is an eventual possibility.

Tapestry uses [sqlite3](https://www.sqlite.org/index.html) and the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. ( i've used https://jmeubank.github.io/tdm-gcc/ with good success. )
