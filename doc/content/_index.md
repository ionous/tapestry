---
title: Tapestry
type: homepage
summary: Tapestry documentation.
geekdocBreadcrumb: false
params:
  hideindex: true # w/o this, the list template will generate bad content
---
A story engine for games.

# Premise

Modern game engines -- while providing amazing support for graphics and game play -- don't have great tools for **building stories**. Game scripting is integrated deeply into the engine. Story development, while simple at the start, quickly becomes time consuming and error prone.

Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This enables quick iteration for commonplace tasks such as creating dialog, puzzles, and quests.

**Tapestry** provides a way to build and test stories independent from the game engine; independent of graphics.

## Working Example

This is part of "Cloak of Darkness" story ( ported from [Robert Firth's original](https://www.ifwiki.org/Cloak_of_Darkness) ): 

> The Foyer of the Opera House is a room. You are in the foyer.
> The description of the Foyer is "You are standing in a spacious hall, splendidly decorated in red and gold, with glittering chandeliers overhead. The entrance from the street is to the north, and there are doorways south and east." <br/><br/>
> The entrance is a  door in the foyer. North from the foyer is the entrance. Through the entrance is the Street. Instead of traveling through the entrance:<br/>
> &nbsp;&nbsp; - Say: "You've only just arrived and besides the weather outside is terrible."

Building this with Tapestry produces a playable story and a SQLite database containing the objects in the game, their interactions, and the game rules. The complete story can be found here: [cloak.tell](https://git.sr.ht/~ionous/tapestry/tree/main/item/content/stories/cloak.tell).

## Relation to interactive fiction

The ability to play a game without graphics is, to my mind, a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar set of game rules as Inform. And, Tapestry tries to provide a similar ( if less extensive ) way to model a game world using English-like sentences.

It is *not* a goal to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. Tapestry *extends* interactive fiction into any game genre.
