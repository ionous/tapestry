// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const Type_{{Pascal name}} = "{{name}}";
{{/with}}
`;
