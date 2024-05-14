---
weight: 100
---

You can create a new story file using the `tap` command.  For example: 

```shell
> tap new mystory
```

will create a simple story of one room called "The Empty Space", and place the player in it. The story file will be created in your `Documents` directory at `Documents/Tapestry/stories/mystory.tell`.

```yaml
# This story was created by you using 'tap new'.
Tapestry:
- Define scene:requires:with:
  - "space"          # <-- a unique name your story.
  - "tapestry"       # <-- most stories need the tapestry standard library.
  - - Declare: """
      The title of the story is "Untitled".
      The author of the story is "you".
      The Empty Space is a room. You are in the space.
      The description of the space is "An empty space, waiting to be filled with life."
      """ # "
```

Stories are stored in a text file format called `tell`. The tell format is loosely based on [yaml](https://en.wikipedia.org/wiki/YAML). It consists of commands and the values for those commands.

Commands consist of predefined command names and their parameters, separated by colons. Which commands are valid depend on the context. There are commands for modeling (creating the world), for define rules and reactions to events (TODO: link), for defining new parser grammars (TODO: link), and so on.

The elements of a .tell file are:
* Command Signatures
* Command Values 
* Comments

**Command Signatures** act like function calls. There are a set of predefined names and parameters for each command. The name appears first, followed by colons for the parameters name. Together the command name and parameters are called a "signature."

**Command Values** every parameter in a command signature is matched with value. The values for a command follow the signature in a yaml-style list, indented with a minimum of two spaces. 

**Comments** start with a hash `#` followed by a space, and run to the end of line.

Every tell file starts with the `Tapestry:` command. 

The complete list of commands can be found here. (TODO: link)

