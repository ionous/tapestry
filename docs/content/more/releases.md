---
# title: "Version History"
#
# for verification:
# *** update this file with new info **
# npm run build 
# go test -count=1 ./...
# npm run tap -- weave check
#
# for sourcehut:
# git checkout main
# git merge <branch>
# ex. git tag -a v0.24.6 -m "Documentation and various changes to make the command language more consistent with itself."
# git push --follow-tags
#
# for github:
# git checkout release
# git merge main
# git push --follow-tags
# verify: 
#
# for documentation:
# npm run build  ( or just the above )
# npm run -w docs publish 
# verify: https://tapestry.ionous.net/
#
--- 

To install Tapestry, please see the [Getting Started](/getting-started/#installing-tapestry) page.

# Status

Tapestry is a work in progress. You could use it to create some simple text based games, but you will encounter bugs and missing features. 

Games are currently playable at the command line, there is a bare bones version running inside the [Godot](https://godotengine.org/) game engine, and its possible to create web playable versions using go to build wasm. A list of various to-dos is [here](https://todo.sr.ht/).

Near term goals include:

* Improving documentation.
* Adding features to allow for more complex stories. ( dialog? )
* Improving engine integration using godot as an example.

# Version History

**v0.24.7**: 

-  WebAssembly! with two limitations:
    1. wasm doesn't support save/load. ( need to write storage into indexeddb, or similar. )
    2. stories must live completely within a single scene. ( most stories can survive as a single scene, but it can be helpful for testing to have multiple. scenes need to be able to add and remove objects. that's handled by sqlite, but sqlite isn't used in the wasm version. additional programing work would be needed. )
- Improved documentation.
- Added Github mirror for discoverability ( and ease of reporting issues, forking, etc. for other people. )
- Cleanup of go module layout to try to improve the `go install` experience. And moved 'moasic' to its own module in "engines" so that `tap` isn't dependent on wails ( nor any other big third-party package. )

**v0.24.6**: 

- Documentation and various language changes to make the commands more consistent.

**v0.24.5**: 

- Implements save/load. The model db is treated as read-only, a separate in-memory "rt" database is mounted and then serialized to/from disk.
- Simplifies object and kind creation. ( As part of save/load support for serializing records. )
- Merges package generic into package rt.  ( Required for the simplified creation. )
- Updates to go1.22.  ( Required for rand/v2, which is needed for save. )
- Small changes to make staticcheck happier.

**v0.24.4**:

- Story files now support mixed plain-text and structured-command sections. For example: `The Empty Space is a room. You are in the space.` is now a valid .tell file. 
- Jess now handles english-like rules ( `Instead of examining the message:` )
- Added new `tap` tool commands: "tap new" to create new stories. "tap version" to report the latest git tag ( only works if using tap install. ) Also changed tap to be more "go" like ( ex. "tap check cloak", instead of "tap check -scene cloak" )
- Story files and shared library scripts now sit side-by-side in the source tree content folder. For the user local document folder: if there is no "shared" folder, `tap` will use a built-in copy. And if there is no "stories" folder, tap will create it.
- Added two very simple "sample" stories.
- **Known Issues**: blockly is (probably) broken because of the .tell changes; that's fine for now.

**v0.23.3**: Expanded english-like parsing (aka 'jess'.) Handles directions and room creation. Verbs replace the experimental macro system (ex. 'carrying', 'wearing'.) Improved scheduling during weave so command statements and jess statements can work better with each other.

**v0.23.2**: English-like parsing in the style of Inform. ( "The bottle is a transparent, open, container.", "Understand "jean" and "genie" as the bottle.", "The bottle has the description "still needs some polish.", etc. ) More phrase parsing still to come.

**v0.23.1**: Serialization revamp. Stories now use [tell](github.com/ionous/tell) (aka. yaml) instead of `json`. Encoding and decoding use reflection and autogenerated typeinfo rather than autogenerated marshaling ( significantly reduces code size, and improves code readability; debugging. )

# The old days 

Tapestry started life as the [Sashimi game engine](https://github.com/ionous/sashimi) which was used for ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) - a point-and-click adventure game. I've used similar techniques on [other projects](https://www.linkedin.com/in/ionous/) over the last 20 years, but Tapestry is its own creation. 
