// Package markup converts a custom html-like markup into markdown-like text.
//
// Using an interpretation of English-like text as a series of characters in
// left-to-right columns and top-to-bottom rows, then:
//   - Newlines: move to the start of the next row.
//   - Paragraphs: act as a newline followed by an empty row.
//   - Soft-newlines: at the start of a row, these do nothing. Otherwise they act like a newline.
//   - Soft-paragraphs: after an empty row, these do nothing. Otherwise they act like a paragraph.
//   - Lists: a set of uniformly indented rows which may contain nested lists of deeper indentation.
//   - Ordered lists : each row in the list starts with a number, starting with 1. increasing by 1.
//   - Unordered lists: each row in the list starts with some uniform marker, ex. a bullet.
//
// The custom html elements include:
//
//	<b>   - bold text
//	<i>   - italic text
//	<s>   - strikethrough
//	<u>   - underlined text
//	<hr>  - horizontal divider
//	<br>  - new line
//	<p>   - soft-paragraph ( there is no "hard" paragraph marker )
//	<wbr> - soft-newline
//	<ol></ul> - ordered list
//	<ul></ul> - unordered list
//	<li></li> - line item in a list
//
// Nested tags are okay; interleaved tags are not supported.
// For example:
//
//	<b><i></i></b>
//
// is fine. While this:
//
//	<b><i></b></i>
//
// will cause trouble.
//
// Unknown tags are left as is.
// Tags attributes ( <b class="..."> ) are not supported.
package markup
