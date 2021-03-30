// Package text converts an custom html-like markup into markdown-like text.
//
// Empty paragraphs <p> are elided making <p> a conditional request to start a paragraph.
// Unclosed tags <br> and <p> dont expect a closed tag.
// Ordered and unordered lists establish indented regions.
// Unknown tags are left unchanged.
// Attributes are not supported.
package text
