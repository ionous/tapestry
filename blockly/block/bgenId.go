package block

import "math/rand"

// from google-blockly/core/utils/idgenerator.js
/**
 * Legal characters for the universally unique IDs.  Should be all on
 * a US keyboard.  No characters that conflict with XML or JSON.
 * Requests to remove additional 'problematic' characters from this
 * soup will be denied.  That's your failure to properly escape in
 * your own environment.  Issues #251, #625, #682, #1304.
 */
const soup string = "!#$%()*+,-./:;=?@[]^_`{|}~" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

/**
 * Generate a random unique ID.  This should be globally unique.
 * 87 characters ^ 20 length > 128 bits (better than a UUID).
 * @return {string} A globally unique ID string.
 */
func newId() string {
	var id [20]byte
	for i := 0; i < len(id); i++ {
		at := rand.Intn(len(soup))
		id[i] = soup[at]
	}
	return string(id[:])
}

var NewId = newId
