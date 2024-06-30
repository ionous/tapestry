import { h } from 'vue'

// given some output; what character should it become?
const fs= {
  "hr":   {tag: "hr", void: true},
  "p":    {char: '\v'},
  "br":   {char: '\n'},
  "wbr":  {char: '\r'},
};

// allow these tags to pass through as-is.
const permit= [
  "b", "i", "li", "ul", "ol", "u", "s"
];

// return false or the data to generate html.
function selectTag(tag) {
  return permit.includes(tag) ? { tag } : fs[tag];
}

export default selectTag;

