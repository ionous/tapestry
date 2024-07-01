---
title: Editor Customization
weight: 4
---

In both Sublime and VSCode, you can tell the editor to use "yaml" highlighting for "tell" files. This will get you some basic syntax colorization for commands, comments, strings, etc.


{{< figure alt="example of syntax highlighting" src="../highlighting.png" width="60%">}}

It's not always perfect when the file contains plain English sections, but it's still very helpful.

# Syntax checking and auto-completion

Although it only works in a preliminary way, it's also possible to ask Sublime ( and VSCode ) to attempt auto-completion and hover documentation for the commands. 

{{< figure alt="example of inline documentation" src="../schema-example.png" width="60%">}}

In Sublime you need use its package manager to install [LSP](https://lsp.sublimetext.io/) and the [LSP-yaml](https://github.com/sublimelsp/LSP-yaml) plugin.  ( The procedure for VSCode is similar. )

Then, at the top of your `.tell` file add the comment:

```
# yaml-language-server: $schema=http:/tapestry.ionous.net/schema/tap.schema.json
```

After that, you should start to get hover and auto-completion for most commands. **However**, until some `tell` specific "language server" is written, it won't play nice when the file includes plain text sections. It will sometimes display the wrong command, or placeholder characters without any actual text. The utility is therefore questionable.