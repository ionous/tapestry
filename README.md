# iffy
This is a reimplmentation of the Sashimi interactive fiction engine with some lessons learned. 

The flow of new story creation is:
1. Use the web-based "composer" to write stories and supporting scripts.
2. Use the tools to: 
    - first, generate "ephemera" from the story files ( other sources -- ie. art assets --- can generate ephemera, too. )
    - second, assemble a "game database" from the ephemera.
3. The "story engine" reads and writes to the gamedb during play.
4. A "game client" then sends commands to ( and listens for events from ) the story engine to progress play.
    - Clients can be command line like traditional interactive fiction;
    - Custom like ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) which used Sashimi's engine;
    - or someday Unity, etc.

Rough versions of the iffy composer, ephemera and gamedb exist. The runtime for the story engine exists, but it lacks a proper game loop. 

Current goals include:
* ✅ porting parent-child containment ( ex. rocks in a box, or people in a room. )
* ✅ printing the contents of a room. ( ex. `> look.` )
* ✅ creating new event system: user actions, state changes, custom notifications.
* [] an initial complement of useful game actions ( looking, examining, taking, dropping, etc. )
* [] handling the player object.
* [] getting a single room story running with a proper game loop.

Ongoing work includes:
* expanding the abilities of the story engine.
* improving the composer to more easily write stories and tests.
* improving the underlying modeling language and its scripting.

# building iffy
I'm going to try to commit to tracking tasks on [trello](https://trello.com/b/EEPnJ6ew/iffy) and writing updates on the [iffy wiki](https://man.sr.ht/~ionous/iffy/). Nothing is probably usable by other people yet. It is a work in progress.

note: iffy uses [sqlite3](https://www.sqlite.org/index.html), the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. ( i'm using: https://jmeubank.github.io/tdm-gcc/ )
