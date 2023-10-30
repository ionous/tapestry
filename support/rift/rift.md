
Types
---

### scalars:
. num
. bool
. interpreted string ( double quotes )
. raw string ( backtick )
. array ( comma separated, ending with fullstop )


### sequences:
arrays of values 

### mappings:
dictionaries of signatures to values

### comments:

### values: 
the "right hand side" elements stored in a sequence or mapping


### heredocs
heredocs provide multi-line strings anywhere a scalar string is permitted.

there are two types: 
	1. interpreted strings, indicated by triple quote 
	2. raw strings, indicated by triple backticks

unlike other heredocs, indentation is based on the position of the closing heredoc marker.

future: customizing the closing heredoc tag.
	