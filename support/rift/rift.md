Purpose
--------
A yamlish format, one that's close enough to use existing yaml highlighting ( in editors, etc. )

some major differences:
* documents hold a single value
* tabs are always invalid whitespace
* string scalars must be quoted
* simplified literal strings
* comments are important


Types
---

### documents
documents hold a single value.

### values:
any scalar, sequence, mapping, or heredoc.

### collections:
* document: a collection containing some one value.
* sequences: aka lists, arrays, or slices: a series of values.
* mappings: aka ordered dictionaries: relates string keys to arbitrary values. 

### scalars:
a scalar value always appears on a single line.

* **num**: 64-bit int or float numbers optionally starting with +/-; floats can have exponents `[e|E][|+/-]...`; hex values can be specified with 0x...
* **bool**: true, false
* **interpreted string** ( double quotes ): "backslashes indicate escaped characters"
* **raw string** ( backtick ): `backslashes are backslashes`
* **array**: comma separated scalars all on one line, ending with an optional fullstop: `1,2,3.` the fullstop is optional, but required for indicating empty inline arrays.

### heredocs
heredocs provide multi-line strings anywhere a scalar string is permitted. ( for the sake of round trip preservation, these might be indicated by a custom type. alternatively -- or in addition -- they could be stored with t heir markers with helper functions to subslice out the formatted text. )

there are two types: 
	1. interpreted strings, indicated by triple quote: newlines are presentation; double newlines provide structure.
	2. raw strings, indicated by triple backticks: newlines are structure.

unlike other heredocs, indentation is based on the position of the closing heredoc marker.

future: customizing the closing heredoc tag.

#### sequences:
entries in a sequence start with a dash, followed by whitespace, and a value.
each new entry starts on a new line with the same indentation as the previous entry.
```
	- true
	- false
```

like yaml, whitespace afer the dash can include newlines. like yaml, the above definitions implies that nested sequences can be declared on one line. for example, `- - 5` is equivalent to `[[5]]` in javascript.

#### mappings:
dictionaries of signatures to values

### whitespace
* whitespace is restricted to the ascii space and newline character;
* cr/lf is not handled;
* tabs are disallowed everywhere except inside string scalars. tabs are not even allowed in comments. this rules differs from yaml. for example, in yaml tabs can appear after indentation.


### comments:
comments start with the `#` marker and continue to the end of a line.

hate me forever, comments for documents are preserved and are significant. therefore, unlike yaml, they must follow the indentation flow ( the price, is that they can also interfere with that flow. ) 

the comments for a given collection is conceptualized as a single block of text overlaying that collection. within that block, all values of a sequence and all keys in a mapping are indicated with \v. a comment to the right of a value is indicated with \t ( if there is no such comment, no htab is recorded. )

comment markers are recorded in the block so that empty inter-comment content-less lines can also be preserved; both are space trimmed. the block overall can be trimmed of all newlines and tabs.

the resulting block of text gets stored in the zeroth index of arrays, the blank key of mappings, and the comment field of a document -- this means that all sequences are one-indexed.

nested comments would probably behave as follows:
```
-
  # comment for the preceding sequence
  name: Mark McGwire # comment for the current mapping
  hr:   65
  avg:  0.278
```  

a program that wants to read ( or maintain ) comments can split or count by vertical tab to find the comments of particular elements.

#### rationale:
For tapestry, comments should be preserved, presented, and editable by authors: not just in the original story text file: but also in other kinds of editors ( ie. mosaic ) and any extracted documentation. 

