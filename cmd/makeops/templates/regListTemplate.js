// regListTemlpate.js
'use strict';
module.exports =`
{{#if list.length}}
var {{which}} = []{{RegType}}{
{{#each list}}
  (*{{Pascal name}})(nil),
{{/each}}
}
{{/if}}`;
