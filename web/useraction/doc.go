// Package useraction: handle long running actions initiated from the browser.
// For example, selecting a file to open.
// POST returns a token and value
// if the request times out before a valid value is available
// the client can poll via GET passing the token until the requested value becomes available.
package useraction
