// Package term helps describe parameters, locals, and return values for patterns.
// While this uses records to hold values that's not a perfect match.
// For instance, to avoid aliasing generic.Value copies record values
// but patterns want use record references to serve as in/out values.
// This package looks suspiciously similar to core.Assignment.
package term
