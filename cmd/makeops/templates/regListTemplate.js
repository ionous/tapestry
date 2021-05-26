// slotTemplate.js
'use strict';
module.exports =
`{{#if list.length}}
var {{which}} = []interface{}{
{{#each list}}
  (*{{Pascal name}})(nil),
{{/each}}
}
{{/if}}`;
