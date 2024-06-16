---
title: Alternate Syntax
weight: 2
---

Some types of values allow for alternate specifications inside the `.tell` file.

## Lists 

Lists contain several values all of the same type. They are only allowed when a command specifically asks for them. ( re: strict typing. ) However, `tell` files support a shortcut so that single values can be specified when a list is required. Tapestry automatically transforms them into a list for you.

This appears most often with commands that require other commands. For instance, the `do` parameter of an `if` statement requires multiple entries:

```yaml
- If:do:
  - true
  - # do two things:
    - Say: "hello"  # dash for the first thing
    - Say: "world"  # dash for the second thing
```

If you only had one thing to say, then you could specify that as a list of one entry:

```yaml
- If:do:
  - true
  - # do one thing:
    - Say: "hello world"
```

Or, you could specify a single entry without the list:

```yaml
- If:do:
  - true
  - Say: "hello world" # no sub list, do just one thing.
```

## Arrays

Arrays show up only rarely. Like lists they are several values specified together, however they only support the so-called "primitive types": inline quoted text, numbers, and boolean values. They must appear all together on the same line, surrounded by square brackets.

For example, assume there's some command that takes a bunch of numbers. You could specify a <em>list</em>em> of those numbers like this:

```yaml
MarsNeedsNumbers: 
  - 2
  - 5
  - 6
```

or, you could specify an <em>array</em> of those numbers like this:

```yaml
MarsNeedsNumbers: [2, 5, 6]
```

Tapestry treats them the same.

## More about Quoted Text

There are actually a couple different **flavors of quoted text**. Every piece of quoted text is either considered "interpreted" or "raw"; and can be either specified 'inline' or as a 'heredoc.'

The specifics of the different string types can be found in the [tell documentation](http://github.com/ionous/tell) 
 
# Json

There is a wholly different, alternative file format to `.tell` which uses the `.if` extension and contains [`json`](https://en.wikipedia.org/wiki/JSON) data.  This is the format that Tapestry uses to save commands to the story database.  It is more verbose than `tell`. For instance, all plain-text sections must be manually wrapped in `Declare:`and `Comment:` commands.

It is unlikely anyone would want to author files using `.if`. But, it exists.

```json
{
    "--": "This is an example comment for an example command.",
    "Say:": [
        "Hello, json."
    ]
}
```