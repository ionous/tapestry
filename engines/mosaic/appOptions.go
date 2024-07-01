package main

import "io/fs"

// Exposes build configuration for programmatic use
type buildConfig int

//go:generate stringer -type=type buildConfig int
const (
	// expects the npm vite server in the www directory is running.
	Web buildConfig = iota
	// expects the npm vite server in the www directory is running.
	Dev
	// expects the www assets have (already) been built into the www/dist directory.
	Prod
)

// In prod, this points to the www/dist directory; otherwise its nil.
var Frontend fs.FS
