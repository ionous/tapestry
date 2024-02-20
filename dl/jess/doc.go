// Package jess helps to turn english like sentences into story fragments destined for a story db.
// For example: `The kitchen is a room. The closed container called the cabinet is in the kitchen.
// The cabinet contains a mug. The mug is transparent.`
//
// Matching uses parse trees defined by Tapestry commands defined in jess.tells.
// The top level command -- the root of all parse trees -- is the type "MatchingPhrases".
// Each successfully matched sentence results a single MatchingPhrases object
// with one ( and only one ) of its members containing the parsed data.
package jess
