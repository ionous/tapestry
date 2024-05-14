---
title: Creating a world
---

A **tapestry game** consists of a bunch of objects placed in one or more rooms. The player can move through those rooms and interact with those objects over a series of turns.

A **tapestry story** can be thought of as having three parts: a **description** of the world, the **actions** a player can take to interact with the world, and the **changes** to the world that occur when the player performs those actions.

The process of describing the world is called "modeling", but before ....

Before jumping into how to declare objects and place them in the world, it's maybe useful to define what *kinds* of objects can exist. 

## Tapestry shared library

Tapestry comes with a "shared library" containing a bunch of different kinds of objects, as well as some common kinds of actions and reactions. 

For instance, the library defines a kind of object called a "container" which can hold other objects. A player can open a container, take objects out or put objects in a container, move a container, etc.

( TODO: link ) The library also -- importantly -- contains a few built-in objects like the player's actor.

Most stories will want to include the shared library called... perhaps, somewhat confusingly.. "Tapestry." The `tap new` tool includes it automatically creating new stories.

## Kinds

If you are familiar with object oriented programming, a **kind** is like a **class**. It defines a set of properties that an object can have, and how objects of that type behave in response to actions. 


There are a dozen or so kinds defined by the Tapestry shared library. The most important of which are probably: **things**, **containers**, **supporters**, and **actors** ( TODO: link ). Followed closely by **rooms**, **doors**, and **directions** which are described here.  (TODO: link ). **Story data**, and **settings** are useful for configuring the game itself, and are described here. (TODO link.) The complete list of kinds can be found here ( TODO: link. )

( TODO: diagram )

You can also define your own kinds, which is described here (TODO )


What distinguishes between the different kinds are their properties and the actions that operate on those properties. 

## Objects 

An object is any kind of abstract or physical noun. Everything that a player can interact with is a kind of object, the kind exists to hold properties common to both abstract and physical nouns: mainly having to do with naming.

An abstract noun is something that helps define the game, but that the player cant interact with directly.

Directions in the game are a good example of abstract nouns. The story ( and the player ) can refer to the idea of "north" or "south" -- but they can't physically touch or take a direction.
 
A door, however, is a physical object in a given location that can be open, shut, locked, etc. It is not abstract; it's "concrete." Concrete objects are usually represented as "things".
 
## Things 

Things are the most common kind of object for both because they are the default kind of object, and because most of the other kinds inherit from things.

> The chair is a thing.

By default, things can be picked up and moved around. And, they appear in the descriptions of the rooms in which they are located.  ( Note: that "rooms" can't be moved around, can't appear in other rooms, etc. They are "objects" but they aren't "things." )

> The chair is fixed in place.

All things can act as wearable objects, which can be taken on and off during the course of the game. ( More on clothing later )

> The player is wearing a hat.

## Containers

Containers have special properties allowing them to hold other objects. They can be open or closed by the author or the player. 

> The treasure chest is a container. The chest contains some coins. The chest is closed and locked. 

They prevented from opening and closing by declaring them to be "unopenable." Another way of thinking about this is that they have no lid.

> The vase contains a rose. The vase is open and unopenable.

Saying one object "contains" another ( or that one is "in" the other )immediately implies the first is a container  For instance, this is fine, if a bit redundant:

> The vase is a container. The vase contains a rose.

## Supporters 

Supporters, like containers, can hold other objects. But, rather than enclosing those objects; those objects rest upon them in some way.

> The table is a supporter. The book is on the table.

There's not much difference between an open, unopenable container and a supporter except that when Tapestry prints information about them in the game it will use the word "in" for containers and "on" for supporters.

So for instance:

> The vase is on the table. The vase contains a rose.

Looking while in its room should print:

> <i>You can see a table (on which is a vase (in which is a rose)).</i>

Like "contains" does for containers the "supports" ( or "on" ) implies a supporter.

## Actors 

An actor is the base kind for anything animate. Every game has at least one actor: the player. ( Also called "you" and "self". )  The only distinguishing property it has is a special "lighting status" to help with determining whether an actor can see.  However, there are a set of actions that only actors can take. ( For instance, taking, dropping, smelling, jumping, etc. )

The verbs "carrying", "wearing", and "holding" imply an actor:

> Alice carries an axe. Bob carries a shield. The player holds the tomato cheese sandwich.

Tapestry doesn't make any special differentiation between men and women, and doesn't plan to. There *is* a plan for being able to apply sets of pronouns for individual actors.... eventually.

# Rooms and Directions

For objects to appear in the game world, they need to have a location. And in Tapestry, defining locations starts with a room. Every story needs to have at least one room. Whether there's more than one depends on the needs of the story in question. But, don't let the name fool you. A "room" could be a cave, it could be a planet, or galaxy. It has no predefined size or scale. It is simply a space with a name.

> The ant hill and the car park are rooms.



# Containment