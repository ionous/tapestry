

# Premise

Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This enables quick iteration for commonplace tasks such as creating dialog, puzzles, and quests.

Modern game engines, however, while providing amazing tools for the visual and arcade-like parts of user experience don't have great tools for building *stories*. While most provide scripting support, the scripting is inevitably integrated deeply into the graphical engine. As a result, development, while quick at the outset, often becomes time consuming and error prone the more complex a game becomes.

**The goal of Tapestry** is to allow designers to describe a game world -- its people, places, and objects; and their *interactions* -- quickly and easily *using words*. And then provide a way to connect that work to the more dynamic and visual elements of a game in a simple way. 

# Status

Tapestry is a work in progress. Currently, you could probably use it to create some simple interactive fiction stories, but you would doubtless encounter some significant issues. A grab bag of various to-dos is [here](https://todo.sr.ht/).


The flow of story creation in Tapestry is:

1. Write stories using the browser based "Mosaic editor" or with raw script.
2. Import stories and optionally other data (ex. from game assets.)
3. Weave together a game database which provides a macro view of your game world.
4. Test and iterate on the story using an interactive console.
5. Embed the story engine into the game engine of your choice, connecting story objects to art assets, and linking the logic of the game to the visuals.
6. Test and iterate some more 😊 using the game console and game database to track progression, inspect the game world, run unit tests, etc.

Rough versions of everything exist, and there is a minimally playable console version integrated with [godot](https://godotengine.org/)

Near term goals include:

* Improve integration with godot.
* Add features to allow for more complex stories.
* Improve the editor to more easily write stories and tests.
* Improve the underlying modeling and scripting tools.

This effort started as a re-implementation of the [Sashimi interactive fiction engine](https://github.com/ionous/sashimi) which was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a 2D point-and-click adventure game. I've used similar techniques for [various projects](https://www.linkedin.com/in/ionous/) I've worked on over the last 20 years, but Tapestry is its own creation based on lessons learned.

# Relation to interactive fiction

The ability to play and test a game without graphics is, to my mind, a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar level of interactivity with a similar set of game rules as Inform.

It is *not* a goal to attempt to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. The overriding goal is to *extend* interactive fiction to bring similar tools into any game genre.

# Licenses

All original source code for Tapestry is licensed under a BSD-3 license.

Tapestry is written primarily in [Go](https://go.dev/). Go and its libraries use a BSD style license. All existing third party Go language dependencies use the MIT license. Tapestry uses [SQLite](https://www.sqlite.org/) for data storage, SQLite is public domain. Tapestry's "Mosaic" editor and web console use [Blockly](https://developers.google.com/blockly) and [Vue.js](https://vuejs.org/) which are licensed under the Apache 2.0 and MIT licenses respectively. Development of the editor and web console require [Node.js](https://nodejs.org/) and a variety of modules each of which use their own open source licenses. [Godot](https://godotengine.org) is used for demonstrating 3D party engine integration. It uses an MIT license.

Please see the Tapestry [LICENSE](https://git.sr.ht/~ionous/tapestry/tree/main/item/LICENSE) file for details.

# Code of Conduct

We believe that Black lives matter, and that queer and trans lives matter. In that spirit, if you are interested in contributing to the project: please refer to the [CODE OF CONDUCT](https://git.sr.ht/~ionous/tapestry/tree/main/item/CODE_OF_CONDUCT.md). While we have no control over the uses of Tapestry, or the stories created with it, we hope end users take the same goals to heart.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md) 