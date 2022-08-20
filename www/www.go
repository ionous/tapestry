//go:build production

// Package www exists in "production" builds to embed frontend resources into the go binary.
// The resources are compiled by vite.js using `npm run build`.
package www

import "embed" // note: Patterns may not contain ‘.’ or ‘..’ or empty path elements, nor may they begin or end with a slash.

//go:embed dist
var Dist embed.FS
