
Types
---

### documents
documents hold a single value.
unlike yaml, there is no multi-document mode.

### values:
any scalar, sequence, or mapping

### whitespace
* whitespace is restricted to the ascii space and newline character; 
* cr/lf is not handled;
* tabs are disallowed everywhere (except inside string scalars). 

these rules differs from yaml. for example, in yaml tabs can appear after indentation.

### scalars:
* num: 64-bit integer or floating point numbers optionally starting with +/-; floats can have exponents `[e|E][|+/-]...`; hex values can be specified with 0x...
* bool: true, false
* interpreted string ( double quotes ): "backslashes indicate escaped characters"
* raw string ( backtick ): `backslashes are backslashes`
* array: comma separated scalars all on one line, ending with an optional fullstop: `1,2,3.` the fullstop is optional, but required for indicating empty inline arrays.


### sequences:
a series of values. 
entries in a sequence start with a dash, followed by whitespace, and a value.
each new entry starts on a new line with the same indentation as the previous entry.
```
	- true
	- false
```

like yaml, whitespace afer the dash can include newlines. like yaml, the above definitions implies that nested sequences can be declared on one line. for example, `- - 5` is equivalent to `[[5]]` in javascript.

### mappings:
dictionaries of signatures to values

### comments:
lines with optional whitespace and `#`, or trailing after any value.
should be round-trippable during de/serialization.

### heredocs
heredocs provide multi-line strings anywhere a scalar string is permitted.

there are two types: 
	1. interpreted strings, indicated by triple quote 
	2. raw strings, indicated by triple backticks

unlike other heredocs, indentation is based on the position of the closing heredoc marker.

future: customizing the closing heredoc tag.
	