---
title: Command Syntax
weight: 1 
---

A structured section of a `.tell` file consists of commands specified as key-value pairs. These commands take the place of functions from a typical programming language. 

<small>( In theory, it's also possible to use a scripting language in place of Tapestry commands; for instance, lua, godot script, etc. For now, that's left as an exercise for the reader... )</small>

## Kinds of commands

There are set of fixed commands provided by the Tapestry engine. Adding new commands requires building Tapestry from source. ( see: `tap code` TODO: link )

Commands are divided into different categories based on the functionality they provide. For instance: commands describe the game world ( modeling ), those that alter the parser, those which listen to game events, those which take actions to change the world, and so on. The [API](/api) documentation will let you catch them all.

Here is an example of doing some math:
```yaml 
- Multiply:value:
  - 5
  - Add:value:
    - 2
    - 3
```

If you are familiar with `yaml`, `.tell` commands are a simplified version of that. Unlike yaml, there are no anchors or references, there is no explicit null value, multi-line strings use a custom "heredoc" syntax, and strings must be quoted. 

Complete documentation for tell lives [here](https://github.com/ionous/tell). The following will cover its use in Tapestry.

## Signatures 

There are two parts to a command: its key ( Tapestry calls this a **signature** ) and its **value**. The signature is on the left side of a command, its value on the right.


Every signature follows the same basic format. A **command name** and one or more **parameter labels**. The name is always capitalized, while the parameters are "camelCase." The first parameter is separated from the command name by a space, while subsequent parameters are separated from each other by colons. Spaces aren't used between labels, not before the colon nor after it.

```yaml
CommandName parameterOne:parameterTwo:
```

Some commands have an **"anonymous"** first parameter. These lack a label, so there's no space between the command name and the first colon. The "Say:" command is a good example of this:

```yaml
- Say: "hello"
```

In general, the number of colons is the number of values needed.

{{< hint title="Parameterless commands" >}}
There are a few rare commands which don't require parameterization. However, due to the way yaml/tell works, these commands still require a trailing colon. You may see them with a value of `true` even though that value is unused.  Breaking out of a loop is a common example:
```
- Break: true
```
{{< /hint >}}

## Values

Commands -- like functions a programming language -- need some sort of input, some sort of values, to operate on. Multiplication, for instance, requires two numbers.

The values available in Tapestry include **quoted text**, **numbers**, **boolean** values ( that is the words `true` or `false` ),  or a **"arrays"** of one of those types. And, values can be supplied using the **results** of  another command.

Tapestry uses **strict typing**. If a declaration command requires a number, then the value supplied must also be a number: It **will not** convert automatically from text to numbers, nor between any other types.

This is good:

```yaml
- Add:value:  # Addition needs two numbers
  - 5 
  - 6
```

This is not good:

```yaml
- Add:value:  # add needs two numbers; neither of these are numbers:
  - "this isn't a number, so it's invalid." 
  - "6"  # <-- this is considered quoted text; not a number.
```

## Multiple values

Values always appear after the colon of the last label, separated from it with either a space or a newline. When there are multiple values needed, the values *must* appear on a new line. In that case, each value appears on its own line, starting with a dash, slightly indented from the signature on the line before.

Dashes indicate a list of values. This works just like a "sequence" in yaml.

## Commands as values

Some commands return a value, and those commands can be used to generate values for other commands.

For example, to add `1 + 4` we can use the `Add:value:` command, specifying the numbers directly.

```yaml 
- Add:value:
  - 1
  - 4
```

Or, if we wanted to add three numbers for the sake of a quick example, we could use:

```yaml 
- Add:value:
  - 1
  - Add:value:
    - 2
    - 2
```

The second `Add:value:` takes the same spot as the number `4`, with its values indented underneath the command.