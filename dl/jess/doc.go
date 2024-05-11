// Turns English-like sentences into story fragments destined for a story db.
// For example:
// 
//  The kitchen is a room. The closed container called the cabinet is in the kitchen.
//  The cabinet contains a mug. The mug is transparent.
// 
// The types of sentences jess can process are based on Inform7's modeling language.
// 
// The matching algorithm uses parse trees defined by Tapestry commands.
// The root of all parse trees is the {"MatchingPhrases:"} command.
// Each successfully matched plain-English sentence results a single MatchingPhrases object
// containing exactly one member set to the parsed data.
package jess