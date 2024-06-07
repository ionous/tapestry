// Copyright 2023 Simon Travis. All rights reserved.
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated via `tap help documentation > alldocs.go`; DO NOT EDIT.
// Edit the source code and rerun 'tap help documentation' to generate this file.

// The 'tap' tool manages Tapestry stories.
//
// Usage:
//
//	tap <command> [arguments]
//
// The commands are:
//
//	check       run tests on existing stories
//	gen         extend tapestry with new golang code
//	edit        run the tapestry story editor
//	play        play an existing story
//	serve       serve a story through http
//	version     print Tapestry version
//	weave       compile a story
//	xform       transform english words
//
// Use "tap help <command>" for more information about a command.
//
// # Run tests on existing stories
//
// Usage:
//
//	tap check [-in path]
//
// Loads an playable database and runs (one or more) test scripts that it contains.
//
// Runs all unit tests by default, use '-run=<name>' to run a specific one.
//
// # Extend tapestry with new golang code
//
// Usage:
//
//	tap gen [-out ../../dl] [-db -dbFile]
//
// Generates .go source code for reading and writing story files from .idl files.
//
// # Run the tapestry story editor
//
// Usage:
//
//	tap edit [-in <directory>] [mosaic flags]
//
// Start the Tapestry story editor.
//
// The 'in' directory should contain two sub-directories:
//  1. "stories" - containing story files ( the target for save/load )
//  2. "ifspec"  - containing interface description files ( these define how to display the story content )
//
// By default, attempts to use a directory called Tapestry in your Documents folder.
//
// # Play an existing story
//
// Usage:
//
//	tap play [-in dbpath] "name of story"
//
// Run a scene within a previously built story database.
//
// Using '-test' can run the list of specified commands as if a player had typed them one by one.
//
// # Serve a story through http
//
// Usage:
//
//	tap serve [-in dbpath]
//
// Run a server that plays a scene from a previously built story database.
//
// # Print Tapestry version
//
// Usage:
//
//	tap version
//
// Version prints the information about the tap tool recorded at the time it was built.
//
// # Compile a story
//
// Usage:
//
//	tap weave [-in path] [-out path]
//
// Turns story files into produces a playable database.
//
// Using '-check' or '-run=<name>' can run all unit tests, or a specific one.
//
// # Transform english words
//
// Usage:
//
//	tap xform [-stem] [-plural] [-singular] [-normal] "word"
//
// Transform English words in various helpful ways.
package main
