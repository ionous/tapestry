package mdl

import (
	"encoding/base64"
	"hash/fnv"
	"io"
	"strings"
)

// make a unique enough name for the db.
// the input assumes space normalized names.
// ex. "cKnjqKujFjg-after_closing"
//
// rationale: within each category, uniqueness is *at least* the current domain and name.
// ( its actually unique within all parent domains as well )
// so we can combine domain and name to provide a unique id.
// this hashes the combination ( so it cant be picked apart and
// used for anything other than identification. )
// and adds some bit of the original name to help understand an id during de/bugging.
//
// rowid is not as nice. its not consistent over time. its not human readable.
func makeId(domain, name string) string {
	domain = strings.ReplaceAll(domain, " ", "_")
	name = strings.ReplaceAll(name, " ", "_")
	w := fnv.New64a()
	io.WriteString(w, domain)
	io.WriteString(w, "-")
	io.WriteString(w, name)
	prefix := strEncoding.EncodeToString(w.Sum(nil))
	parts := strings.Split(name, "_")
	if len(parts) > 1 {
		// alt: best word break < 16 characters
		name = parts[0] + "_" + parts[len(parts)-1]
	}
	return prefix + "-" + name
}

// alt: https://github.com/tep/encoding-base56
var strEncoding = base64.
	NewEncoding("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_").
	WithPadding(base64.NoPadding)
