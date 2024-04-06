

# Premise

Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This enables quick iteration for commonplace tasks such as creating dialog, puzzles, and quests.

Modern game engines, however, while providing amazing tools for the visual and arcade-like parts of user experience don't have great tools for building *stories*. While most provide scripting support, the scripting is inevitably integrated deeply into the graphical engine. Development, while quick at the outset, often becomes time consuming and error prone the more complex a game becomes.

**Tapestry** provides a way to build and test stories independent of a graphics engine.

## Relation to interactive fiction

The ability to play and test a game without graphics is, to my mind, a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar level of interactivity with a similar set of game rules as Inform.

It is *not* a goal to attempt to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. The overriding goal is to *extend* interactive fiction to bring similar tools into any game genre.

# Status

Tapestry is a work in progress. You could use it to create some simple games, but you would doubtless encounter some significant issues. A list of various to-dos is [here](https://todo.sr.ht/).

The flow of development with Tapestry is:

1. Write stories using either text scripts or Tapestry's browser-based editor.
2. Optionally, use Tapestry's game database to connect to, and validate, other game assets.
3. Embed the story engine into the game engine of your choice, linking the logic of your game to its graphics.
4. Iterate using the game console and runtime database to track progression, inspect the game world, run unit tests, etc.

Rough versions of everything exist, and there is a minimally playable console version integrated with [godot](https://godotengine.org/).

Near term goals include:

* Improve integration with godot.
* Add features to allow for more complex stories.
* Improve the editor to more easily write stories and tests.
* Improve the underlying modeling and scripting tools.

This effort started as a re-implementation of the [Sashimi interactive fiction engine](https://github.com/ionous/sashimi) which was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a 2D point-and-click adventure game. I've used similar techniques for [various projects](https://www.linkedin.com/in/ionous/) I've worked on over the last 20 years, but Tapestry is its own creation based on lessons learned.

## Version History

- v0.23.3: expand english-like parsing (aka 'jess'.) Handles directions and room creation. Verbs replace the experimental macro system (ex. 'carrying', 'wearing'.) Improved scheduling during weave so command statements and jess statements can work better with each other.
- v0.23.2: english-like parsing in the style of inform. ( "The bottle is a transparent, open, container.", "Understand "jean" and "genie" as the bottle.", "The bottle has the description "still needs some polish.", etc. ) More phrase parsing still to come.
- v0.23.1: serialization revamp. stories now use [tell](github.com/ionous/tell) instead of `json`. encoding and decoding use reflection and autogenerated typeinfo rather than autogenerated marshaling ( significantly reduces code size, and improves code readability; debugging. )

# Licenses

All original source code for Tapestry is licensed under a BSD-3 license.

Tapestry is written primarily in [Go](https://go.dev/). Go and its libraries use a BSD style license. All existing third party Go language dependencies use the MIT license. Tapestry uses [SQLite](https://www.sqlite.org/) for data storage, SQLite is public domain. Tapestry's "Mosaic" editor and web console use [Blockly](https://developers.google.com/blockly) and [Vue.js](https://vuejs.org/) which are licensed under the Apache 2.0 and MIT licenses respectively. Development of the editor and web console require [Node.js](https://nodejs.org/) and a variety of modules each of which use their own open source licenses. [Godot](https://godotengine.org) is used for demonstrating 3D party engine integration. It uses an MIT license.

Please see the Tapestry [LICENSE](https://git.sr.ht/~ionous/tapestry/tree/main/item/LICENSE) file for details.

# Code of Conduct

We believe that Black lives matter, and that queer and trans lives matter. In that spirit, if you are interested in contributing to the project: please refer to the [CODE OF CONDUCT](https://git.sr.ht/~ionous/tapestry/tree/main/item/CODE_OF_CONDUCT.md). While we have no control over the uses of Tapestry, or the stories created with it, we hope end users take the same goals to heart.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md) 
