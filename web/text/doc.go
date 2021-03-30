// Package text converts an custom html-like markup into markdown-like text.
//
// Empty paragraphs <p> are elided making <p> a conditional request to start a paragraph.
// The tags <br>, <wbr>, and <p> dont expect a closing tag. ( Their closing tag does nothing. )
// <wbr> is used as a soft-newline even though that doesn't perfectly match html semantics.
// Ordered and unordered lists establish indented regions.
// Unknown tags are left unchanged.
// Attributes are not supported.
package text
