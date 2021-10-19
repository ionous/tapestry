// sigMap.js
'use strict';
module.exports =`
var Signatures = map[uint64]interface{}{
{{#each list}}
  {{@key}}: (*{{Pascal type}})(nil), /* {{sig}} */
{{/each}}
}`;
