---
weight: 1
title: The Tell Format
---
Tapestry stories are text files with the `.tell` extension. [Tell](http://github.com/ionous/tell) files a variant of [yaml](https://yaml.org/), consisting of alternating plain text and structured command sections, along with optional comments.

The **plain text sections** always consist of English-like statements such as `The pen is on the desk. The man is carrying the plan called Panama.` and so on. Every story must contain at least one room and the player; so this is the simplest possible story ( and the simplest possible `.tell` file. )


> The Lobby is a room. You are in the Lobby.

The **structured command sections**, on the other hand, contain a series of lines consisting of key-value pairs. Such a section might look something like this: 
```yaml
- Say: "Press 'Q' to quit."
- Interpret:with:
  - - "quit"
    - "q"
  - - Action: "request to quit"
```
 
**Section separators** mark the end of one section type, and the start of the next. They are three dashes on their own line:

> `---`

Finally, both sections can contain **comments**  ( ie. notes to yourself or other authors which don't affect the actual story. ) Comments start with hash and run until the end of the line. 

> `# This is a comment.`

{{< hint type=caution >}}
Comments <em>must</em> can only outside of quoted text.
{{< /hint >}}

```yaml
- Say: "Hello world."  # This is a comment
- Say: "This does not # contain a comment."
```

## What's the difference?

The two section types have different purposes. While anything that can be described in plain text can be described using commands ( although not the other way around... ) most often, the **plain text describes the world**, and the **commands control how play proceeds.**

This example says the same thing using both methods; the plain text version is simpler and easier to understand. 

```
The Lobby is a room. You are in the Lobby.
---
- Define nouns:as:
    - "lobby"
    - "room"
- Relate:to:via:
    - "self"
    - "lobby"
    - "whereabouts"   
```

