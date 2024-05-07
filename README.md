# Premise

Modern game engines -- while providing amazing support for graphics and game play -- don't have great tools for **building stories**. Game scripting is integrated deeply into the engine, and story development, while quick at the start, becomes time consuming and error prone the more assets are added to the game.


Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This enables quick iteration for commonplace tasks such as creating dialog, puzzles, and quests.

**Tapestry** provides a way to build and test stories independent from the game engine; independent of graphics.

## Relation to interactive fiction

The ability to play and test a game without graphics is, to my mind, a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar set of game rules as Inform. And, Tapestry tries to provide a similar ( if less extensive ) way to model a game world using English-like sentences.

It is *not* a goal to attempt to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. Tapestry hopes to *extend* interactive fiction development into any game genre.


## Working Example

This is part of "Cloak of Darkness" story ( ported from the Inform7 version of [Robert Firth's original](https://www.ifwiki.org/Cloak_of_Darkness) ): 


> The Foyer of the Opera House is a room. You are in the foyer.
> The description of the Foyer is "You are standing in a spacious hall, splendidly decorated in red and gold, with glittering chandeliers overhead. The entrance from the street is to the north, and there are doorways south and east." <br/><br/>
> The entrance is a  door in the foyer. North from the foyer is the entrance. Through the entrance is the Street. Instead of traveling through the entrance:<br/>
> &nbsp;&nbsp; - Say: "You've only just arrived and besides the weather outside is terrible."

Building this with Tapestry produces a playable story and a SQLite database consisting of all of the objects in the game, their interactions, and the game rules.


# Status

Tapestry is a work in progress. You could use it to create some simple text based games, but you will encounter bugs and missing features. 

The flow of development with Tapestry is:

1. Use `tap new` to create a new story file. 
1. Edit using your favorite text editor ( or Tapestry's experimental [Blockly](https://developers.google.com/blockly/) editor. )
2. Optionally, use Tapestry's game database to bind to, and validate, graphical assets.
3. Embed the story engine into the game engine of your choice.
4. Iterate using the built-in game console. Use the runtime database to track progression, inspect the game world, run unit tests, etc.

Games are currently playable at the command line, and there is a bare bones version running inside the [Godot](https://godotengine.org/) game engine. A list of various to-dos is [here](https://todo.sr.ht/).


Near term goals include:

* Add features to allow for more complex stories.
* Improve the underlying modeling and scripting syntax.
* Improve integration with godot.
* Documentation.


## Version History

Tapestry started as a rework of the [Sashimi engine](https://github.com/ionous/sashimi) which was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a point-and-click adventure game. I've used similar techniques on [other projects](https://www.linkedin.com/in/ionous/) over the last 20 years, but Tapestry is its own creation.

**v0.24.4**:

- Story files now support mixed plain-text and structured-command sections. For example: `The Empty Space is a room. You are in the space.` is now a valid .tell file. 
- Jess now handles english-like rules ( `Instead of examining the message:` )
- Added new `tap` tool commands: "tap new" to create new stories. "tap version" to report the latest git tag ( only works if using tap install. ) Also changed tap to be more "go" like ( ex. "tap check cloak", instead of "tap check -scene cloak" )
- Story files and shared library scripts now sit side-by-side in the source tree content folder. For the user local document folder: if there is no "shared" folder, `tap` will use a built-in copy. And if there is no "stories" folder, tap will create it.
- Added two very simple "sample" stories.
- **Known Issues**: blockly is (probably) broken because of the .tell changes; that's fine for now.

**v0.23.3**: Expanded english-like parsing (aka 'jess'.) Handles directions and room creation. Verbs replace the experimental macro system (ex. 'carrying', 'wearing'.) Improved scheduling during weave so command statements and jess statements can work better with each other.

**v0.23.2**: English-like parsing in the style of Inform. ( "The bottle is a transparent, open, container.", "Understand "jean" and "genie" as the bottle.", "The bottle has the description "still needs some polish.", etc. ) More phrase parsing still to come.

**v0.23.1**: Serialization revamp. Stories now use [tell](github.com/ionous/tell) (aka. yaml) instead of `json`. Encoding and decoding use reflection and autogenerated typeinfo rather than autogenerated marshaling ( significantly reduces code size, and improves code readability; debugging. )


# Licenses

Tapestry is open source, and all original source code is provided under a BSD-3 license.

[Go](https://go.dev/) and its libraries use their own BSD style license. All existing third party Go language dependencies use the MIT license. Tapestry uses [SQLite](https://www.sqlite.org/) for data storage, and SQLite is public domain. Tapestry's "Mosaic" editor and web console use [Blockly](https://developers.google.com/blockly) and [Vue.js](https://vuejs.org/) which are licensed under the Apache 2.0 and MIT licenses respectively. Development of the editor and web console require [Node.js](https://nodejs.org/) and a variety of modules each of which use their own open source licenses. [Godot](https://godotengine.org) is used for demonstrating 3D party engine integration. It uses an MIT license.

Please see the Tapestry [LICENSE](https://git.sr.ht/~ionous/tapestry/tree/main/item/LICENSE) file for details.


## Code of Conduct

Black lives matter, and queer and trans lives matter. In that spirit, if you are interested in contributing to the project, please refer to the [CODE OF CONDUCT](https://git.sr.ht/~ionous/tapestry/tree/main/item/CODE_OF_CONDUCT.md).

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md) 