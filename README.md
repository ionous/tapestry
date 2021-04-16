# iffy

A story engine for powering narrative driven games. 

# Overview

Most games can be divided into separate continuous and discrete elements. Similar to the idea of separating client-side display from server-side logic. 

Behind the scenes, the discrete elements of a game are typically represented using some set of objects -- aka the game world. The player interacts with that world using some set of commands. Those commands create changes in the world based on custom game rules, and those changes are reported back to the player using a stream of events ( transformed into sounds, vfx, movement, etc. )

Modern engines provide great tools for the "continuous" parts of games, especially graphics and sound. Iffy is focused on the discrete parts of a game. The "server" logic of a game's graphical "client."


# Status

Iffy is a work in progress. Nothing is probably usable by other people yet. In the meantime, I'm committed to tracking tasks on [trello](https://trello.com/b/EEPnJ6ew/iffy), writing updates on the [iffy wiki](https://man.sr.ht/~ionous/iffy/), and occasionally posting on twitter [@theionous](https://twitter.com/theionous).

The flow of story creation is:

1. Write stories using the web-based "composer".
2. Import stories ( and references to other game assets ) into an "ephemera database".
3. Assemble a "game database" from the ephemera.
4. Optionally, test and iterate your story using iffy's interactive console.
5. Embed the story engine into your own game. 

Rough versions of the composer, ephemera and game databases, and the story engine exist. 

Current goals include:

* ✅ parent-child containment ( ex. rocks in a box, or people in a room. )
* ✅ printing the contents of a room. ( ex. `> look.` )
* ✅ an initial complement of useful game actions ( looking, examining, taking, dropping, etc. )
* ( ) implementing the event system: user actions, state changes, custom notifications.
* (in progress) getting a [single room story](http://www.ifwiki.org/index.php/A_Day_for_Fresh_Sushi) running with a proper game loop.

Ongoing work includes:

* Expanding the abilities of the story engine.
* Improving the composer to more easily write stories and tests.
* Improving the underlying modeling language and its scripting.

Eventually, the idea is to provide easy integration with Unity, Unreal, etc.


# Relation to interactive fiction

It's [my](https://www.linkedin.com/in/ionous/) belief that developers should be able to play and test the logic of their game separate from the visual elements whenever possible. Iterating quickly on dialog, puzzles, and quests. This, to my mind, is a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Iffy is therefore inspired by the world of interactive fiction, and owes a lot in particular to [Inform 7](http://inform7.com/). The default game world for iffy attempts to provide a similar level of interactivity with a similar set of game rules as Inform, built on top of iffy's unique story engine.

It is *not* a goal to even attempt to match Inform's amazing natural language programming environment, nor is it a goal to run on z-machines.  It *is*, however, a goal to be able to play *some* "Inform-like" stories with similar results. 

# Building the code

iffy uses [sqlite3](https://www.sqlite.org/index.html) and the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. ( i've used: https://jmeubank.github.io/tdm-gcc/ )

Iffy started as a reimplementation of the [Sashimi interactive fiction engine](https://github.com/ionous/sashimi) with some lessons learned. That engine was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a point-and-click adventure game written in javascript and go.


