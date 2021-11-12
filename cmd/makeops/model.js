// allTypes.js
'use strict';
const fs = require('fs'); // filesystem for loading iffy language file
const vm = require('vm'); // virtual machine for parsing iffy language file
const Make = require('./directives.js'); // composer directives
const TagParser = require('./tags.js');


// parse the idl
// shouldnt really be here, but its fine for now.
let m = new Make();
const folderPath = "../../idl";
fs.readdirSync(folderPath).filter(fn => fn.endsWith(".ifspec")).forEach(fn => {
  const path = folderPath + '/' + fn;
  const raw = fs.readFileSync(path);
  let dict;
  try {
    dict = JSON.parse(raw);
  }
  catch (e) {
    console.log(`couldnt't parse ${path}`);
    throw e;
  }
  console.log(`parsing ${path}`);
  m.readSpec(dict);
});

const sorted = {};
Object.keys(m.types.all).sort().forEach((key) => {
  sorted[key] = m.types.all[key];
});

module.exports = sorted;
