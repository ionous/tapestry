// Package markup converts an custom html-like markup into markdown-like text.
//
// If english text is a series characters in left-to-right columns and top-to-bottom rows, then:
//
// - Newlines: move to the the start of the next row.
// - Paragraphs: act as a newline followed by an empty row.
// - Soft-newlines: at the start of a row, these do nothing. Otherwise they act like a newline.
// - Soft-paragraphs: after an empty row, these do nothing. Otherwise they act like a paragraph.
// - Lists: a set of uniformly indented rows which may contain nested lists of deeper indentation.
// - Ordered lists : each row in the list starts with a number, starting with 1. increasing by 1.
// - Unordered lists: each row in the list starts with some uniform marker, ex. a bullet.
//
// The custom html elements include:
// <b>, <strong>, or <mark> - bold text
// <i>, <em>, or <cite> - italic text
// <s>, or <strike> - strikethrough
// <u>   - underlined text
// <hr>  - horizontal divider
// <br>  - new line
// <p>   - soft-paragraph ( there is no "hard" paragraph marker )
// <wbr> - soft-newline
// <ol></ul> - ordered list
// <ul></ul> - unordered list
///
// Unknown tags are left as is.
// Attributes on tags are not supported.
//
package markup
