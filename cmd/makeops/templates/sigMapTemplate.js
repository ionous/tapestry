// sigMap.js
'use strict';
module.exports =`
var Signatures = map[uint]interface{}{
{{#each list}}
  {{@key}}: (*{{Pascal type}})(nil), /* {{sig}} */
{{/each}}
}`;
