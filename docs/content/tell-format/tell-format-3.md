---
title: Directory Structure
weight: 3
---

Currently, there are two folders used by the `tap` tool: `shared` and `stories`. The **shared** folder contains a common library of useful definitions. If the shared folder doesn't exist, an internal copy of the shared code is used. The more important folder for most people, therefore, will be the **stories** folder.

Currently, the stories folder will always be created in your local "Documents/Tapestry/stories" folder. You can see the exact location by running:

> tap help new

Along with some basic information, it should print something like:

>  New story files are created in:
>    /Users/your_username/Documents/Tapestry/stories

While the default location can't be changed, `tap` has flags which should allow you to override it while creating stories. ( See also, the "Getting Started" guide. )

# Story Scenes

Scenes group story definitions together. Objects declared in one scene can't see objects defined in another scene unless one of the scenes is explicitly defined to depend on the other. 
See the scene documentation for complete details ( TODO: link. )

By default, each file in the stories folder becomes its own standalone scene named. And unless otherwise specified, all scenes inherit from the default `Tapestry` scene provided by the files in the shared folder. This means that while scenes are isolated from each other, they all use the same base set of rules.

## Per-file customization

You can change the name of your scene ( and the set of scenes that it depends on ) using the `Define scene:` command. It **must** appear as the very first command in your story file, and before all plain text declarations. Since every tell file starts with plain text, when customizing the scene, a section divider should appear at the top of your file:

```yaml
---
# Customize the scene used by this story file:
Define scene:requires:
  - "example"
  - - "tapestry" # use the shared library.
    - "scoring"  # use this pre-defined library for tracking a player's score.
---
# First plain text statements here.
The Lobby is a room. You are in the Lobby.
```

## Per-directory customization

In addition to customizing the scene used for a particular file, you can also define the scene used by a whole directory of files using an `_index` file. The index file should contain a single `Define scene:` command.

For instance, the Tapestry `shared` folder contains the following `_index.tell`:

```yaml
---
# Everything in the shared folder
# is considered part of the tapestry scene.
Define scene: "tapestry"
```

{{< hint type=caution >}}
Except for the index file; filenames starting with a dot ( `.` ) or underscore ( `_` ) are ignored when creating stories.
{{< /hint >}}

# Other folders and files

There is one other folder worth mentioning. Stories files must be woven into a database to play. By default, that database is saved to "Documents/Tapestry/build/play.db".

TODO: save game info.