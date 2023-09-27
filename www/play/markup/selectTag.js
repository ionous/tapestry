import { h } from 'vue'

const fs= {
  "hr":   {tag: "hr", void: true},
  "p":    {char: '\v'},
  "br":   {char: '\n'},
  "wbr":  {char: '\r'},
};

const permit= [
  "b", "strong", "mark", "i", "em", "cite", "li", "ul", "ol", "u", "s"
];

// return false or the data to generate html.
function selectTag(tag) {
  return permit.includes(tag) ? { tag } : fs[tag];
}

export default selectTag;

