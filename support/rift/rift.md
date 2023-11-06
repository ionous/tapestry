Purpose
--------
A simple yaml-like format that's close enough to use existing yaml highlighting ( in editors, etc. )


Some major differences:

* Comments are important.
* String scalars must be quoted.
* Multiline strings use heredocs and only heredocs.
* Except for string literals, tabs are always invalid whitespace.
* No flow style ( although is an array syntax. )
* The order of maps matters.
* No anchors or references.
* Documents hold a single value.

Types
---

### Collections
* **Document**: a collection containing a single **value**.
* **Sequences**: aka lists: a series of one or more **values**.
* **Mappings**: aka ordered dictionaries: relates **signatures** to **values**. 

### Documents
Documents are most often text files. UTF8, no byte order marks. 

Whitespace in documents is restricted to the ascii space and newline character; cr/lf is not handled; tabs are disallowed everywhere except inside string scalars. Tabs are not even allowed in comments. ( This differs from `yaml` where, for example, tabs can appear after single space indentation. )

### Values
Any **scalar**, **array**, **sequence**, **mapping**, or **heredoc**.


### Scalars

* **bool**: `true`, or `false`.
* **raw string** ( backtick ): \`backslashes are backslashes.`
* **interpreted string** ( double quotes ): "backslashes indicate escaped characters."
* **number**: 64-bit int or float numbers optionally starting with `+`/`-`; floats can have exponents `[e|E][|+/-]...`; hex values can be specified with `0x`notation. _( **TBD**: may expand to support https://go.dev/ref/spec#Integer_literals, etc. as needed. )_ It is sad that due to comments hex colors cannot live as `#ffffff`.

A scalar value always appears on a single line. There is no null keyword, null is implicit where no explicit value was provided.

### Arrays
An array is a list of comma separated scalars, ending with an optional fullstop: `1, 2, 3.` 
_( **TBD**: all on one line?  )_  The fullstop is necessary when indicating an empty array. Nested arrays are not a thing; use sequences.

#### Sequences
Sequences define a linear list of values. 
Entries in a sequence start with a dash, whitespace separates the value.
Additional entries in the same sequence start on the next line with the same indentation as the previous entry.
```
  - true
  - false
```

As in `yaml`, whitespace after the dash can include newlines. The lack of differentiation between newline and space implies that nested sequences can be declared on one line. For example, `- - 5` is equivalent to the json `[[5]]`.

#### Mappings
Mappings relate signatures to values in an ordered fashion.
**Signatures** are words separated by colons, ending with a colon and whitespace. For example: `Hello:there: `. The first character of each word must be a (unicode) letter, subsequent characters can also include digits and underscores _( **TBD**: this is somewhat arbitrary; what does yaml do? )_

For the same reason that nested sequences can appear inline, mappings can. However, `yaml` does not allow this and it's probably bad style. For example: `Key: Nested: "some value"` is equivalent to the `json` `{"Key:": {"Nested:": "some value" }` 

_( **Note**: [Tapestry](git.sr.ht/~ionous/tapestry) wants those colons so, for now, the interpretation of `key:` is `"key:"` not `"key"`. This feels like an implementation detail because implementations should know what kind of data they are reading anyway. )_

### Heredocs
Heredocs provide multi-line strings wherever a scalar string is permitted ( but not in a multiline array, dear god. )

There are two types, one for each string type:

1. **raw**, triple backticks: newlines are structure; backslashes are backslashes.
2. **interpreted**, triple quotes: newlines are presentation; double newlines provide structure; backslashes are special.

Indentation of the block is based on the position of the closing heredoc marker. Any text to the left of the closing marker is therefore an error.

```
  - """
    i am a heredoc example.
    these lines are run together
    each separated by a single space.
     this sentence has an extra space in front.

    a blank line ^ becomes a single newline.
    trailing spaces in that line, or any line, are eaten.
    """

  - """
    this sentence starts with
1234 spaces.
"""

  - ```END
    i am a heredoc literal using a custom closing tag.
    this sentence is separated from the preceding with a newline.
     this appears on yet another line, with a single leading space.

    a blank line ^ is a blank line.
    in that line, or any line, spaces to the right of the tag are preserved.
    END
```

_( **Note**: for the sake of round trip preservation, heredocs might be indicated by a custom string type. Alternatively -- or in addition -- they could be stored with their markers and helper functions could subslice out the formatted text. )_

### Comments
Hate me forever, comments are preserved, are significant, and can introduce their own indentation level. Therefore, they must follow the collection flow ( and the price is that they can interfere with that flow. ) They are the most complicated part of this format.

**Rationale:** Comments are a good mechanism for communicating human intent. In [Tapestry](git.sr.ht/~ionous/tapestry), since story files can be edited by hand, edited in mosaic/blockly, or even extracted to present documentation: preserving those comments across different transformations matters. 

Comments begin with the `#` hash and continue to the end of a line. Comments cannot appear within a scalar _( **TBD**: comma separated arrays split across lines might be an exception. )_ Here are some examples:

```yaml
- "scalar" # comments can trail a scalar.
           # they establish a new indentation level
           # which any related comment must follow.

# comments can act as a "header" for elements of a collection.
# but only at the indentation level of the collection.
- "scalar"

- # comments can also live in the "padding" between a dash or signature.
  "scalar" 

- # comment attribution can be tricky. for example:
  # is this considered padding? or,
  # is this a header for the next sequence?
  - "first element"
  
  # more obviously, this is a header.
  # we are above the element, and not inside the padding:
  - "second element"
  
  # less obviously, this is a header for nothing.
  # not "null", just literally nothing. no element at all.
  # that's fine. i guess.
```

#### Comment storage:

Internally, the comments for a given collection are stored as a unified "comment block" -- a string of continuous text generated in the following manner:

Individual comments lines are stored when they are encountered. Each line gets trimmed of spaces, but hash marks are kept intact. Meanwhile, the dash ( or signature ) of a collection is recorded as a horizontal tab `\t`, values are ignored, trailing comments ( if any ) are indicated with an additional `\t`, and the end of each entry is indicated with vertical tab `\v`. Newlines (`\n`) separate comments whenever a tab does not.  _( **TBD**: preserve empty lines using `\n`? )_ The resulting block can then be trimmed of trailing newlines and tabs.

For example, the following sequence results in the comment block: `# one\t# two\t# three\n# four\v# five`.

```
# one
- # two
  "value" # three
          # four
# five  
- "other"
```


Each nested collection gets its own separate comment block. The block gets stored in the zeroth index of its sequence, the blank key of its mappings, or the comment field of its document. **This means all sequences are one-indexed.** _(TBD: arrays should probably be one-indexed as well for consistency's sake, and in case they are allowed comments later on.)_

A program that wants to read ( or maintain ) comments can split or count by vertical tab to find the comments of particular elements.

#### Comment splitting:

Regarding what exactly marks the end of an "entry", there is no answer that's correct 100% of the time. 
For instance:

```yaml
# does this text live as part of the document?
# or is it stored in the comment block for the first sequence?
- "do i have a header comment?"

- # and while, presumably this comment is padding....
  # does *this* comment live as part of that padding,
  # or is it considered the start of the nested sequence?
  - "first element"
```

To provide consistency, if not necessarily correct meaning, the default behavior is to continue the previous comment. A blank new line, or a change in indentation ( if permitted by the context ), can be used to enforce a particular interpretation.  

```yaml
# this a document comment.
# this is also a document comment.

# this is part of the sequence due to the intervening blank line.
- "i've always wanted a header."

- # this is padding,
  # and this continues that padding.
    # however, because *this* line has extra indentation:
    # it indicates a comment that will wind up in the subsequent collection:
    - "i've got one too!"
    
- # note that for scalar elements:
    # a change in indentation doesn't mean anything.
    # it's all padding.
    5
    
- long key name: # that behavior is a bit weird
    # but probably has to be allowed
    # because of long key names.
    5
```
