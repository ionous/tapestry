---
geekdocCollapseSection: true
title: Slots

params:
  showindex: true
---

Every Tapestry story file ( `.tell` ) contains a list of Story Statement slots. Various commands can fit into those slots, the complete list of which can be found on the [Story Statement](/api/slot/story_statement) page. Taken together, the story commands allow an author to model a complete game world.

To branch from modeling into changing how a game behaves, some story commands have parameters which specify other types of slots, and those slots allow other types of commands. For instance, to change how a player's input is processed, the story command [Interpret:with:](/api/slot/story_statement#define_leading_grammar") accepts a [Scanner Maker](/api/slot/scanner_maker)Scanner Maker</a>, and scanner maker commands [customize the parser](/api/idl/grammar).

## Index

This is the list of all known slots:
