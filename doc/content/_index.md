---
title: Tapestry
type: homepage
summary: Tapestry documentation.
geekdocBreadcrumb: false
params:
  pageclass: "page-main"
---
A story engine for games.

# Premise

Modern game engines -- while providing amazing support for graphics and game play -- don't have great tools for **building stories**. Game scripting is integrated deeply into the engine. Story development, while simple at the start, quickly becomes time consuming and error prone.

Developers should be able to work on the logic of their games separate from the visual elements whenever possible. This enables quick iteration for commonplace tasks such as creating dialog, puzzles, and quests.

**Tapestry** provides a way to build and test stories independent from the game engine; independent of graphics.

## Working Example

This is a port of the "Cloak of Darkness" originally by written by [Roger Firth](https://www.ifwiki.org/Cloak_of_Darkness): 

{{< include file="/static/_includes/cloak.tell" language="yaml" options="linenos=false" >}}

Building this with Tapestry produces a playable story and a SQLite database containing the objects in the game and their interactions. You can play the story here: <br> 
<span class="gdoc-button gdoc-button--regular"><a class="gdoc-button__link" target="_blank" href="/cloak-of-darkness.html">Play ↗</a></span>

## Relation to interactive fiction

The ability to play a game without graphics is, to my mind, a lot like the original [text adventure](https://en.wikipedia.org/wiki/Colossal_Cave_Adventure) games. 

Tapestry is therefore inspired by the world of interactive fiction and owes a lot in particular to [Inform 7](http://inform7.com/). For that reason, the default game world for Tapestry attempts to provide a similar set of game rules as Inform. And, Tapestry tries to provide a similar ( if less extensive ) way to model a game world using English-like sentences.

It is *not* a goal to match Inform's amazing natural language programming environment. ( Nor is it a goal to run on z-machines. )  It *is* however a goal to be able to play some "Inform-like" stories with similar results. Tapestry *extends* interactive fiction into any game genre.
