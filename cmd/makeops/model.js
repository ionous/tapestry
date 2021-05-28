// allTypes.js
'use strict';
const fs = require('fs'); // filesystem for loading iffy language file
const vm = require('vm'); // virtual machine for parsing iffy language file
const Make = require('./directives.js'); // composer directives
const TagParser = require('./tags.js');

// change to pascal-cased ( golang public )
const pascal = function(name, sep = " ") {
  const els = name.split('_').map(el => el.charAt(0).toUpperCase() + el.slice(1));
  return els.join(sep);
};

// load the language file; brings 'localLang()' into global scope.
vm.runInThisContext(fs.readFileSync(`../compose/www/data/lang/iffy.js`));
vm.runInThisContext(fs.readFileSync(`../compose/www/data/lang/spec.js`));

// parse the idl
// shouldnt really be here, but its fine for now.
let m = new Make();
const folderPath = "../../idl";
const specs = [];
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
  const gopack = fn.substring(0, fn.length - ".ifspec".length);
  // m.currGroups.push( gopack );
  specs.push({
    "name": gopack,
    "uses": "group"
  });
  console.log(`parsing ${path}`);
  m.currGroups.unshift(gopack);

  for (const k in dict) {
    const data = dict[k];
    switch (data.uses) {
      default:
        throw new Error( `unknown or missing uses ${k} '${data.uses}'  in ${path}`);
      case "str": {
        const d = m.str(k, data.spec, data.desc);
        specs.push(d);
        break;
      }
      case "num": {
        const d = m.num(k, data.spec, data.desc);
        specs.push(d);
        break;
      }
      case "slot": {
        const d = m.slot(k, data.desc);
        specs.push(d);
        break;
      }
      case "swap": {
        const d = m.swap(k, data.spec, data.desc);
        specs.push(d);
        break;
      }
      case "flow": {
        // note: doesnt use the "maker" -- manually expands it.
        const parts = data.spec.split(" ").filter(p => p.length);
        const lede = parts[0];
        const fmt = parts.filter((p, i) => i > 0).join("")
        const tags = fmt && TagParser.parse(fmt);

        const tokens = [lede];
        if (!tags.keys) {
          // console.log("no keys for", gopack, k);
          // happens for things like "always", "equal_to", etc.
        } else {
          tags.keys.forEach((t, i) => {
            const el = tags.args[t];
            // '$STR': { label: 'rel', type: 'relation' },
            if (el.label !== "_") {
              if (i === 0) {
                tokens.push(" ");
              } else {
                tokens.push(", ");
              }
              tokens.push(el.label);
            }
            tokens.push(": ");
            tokens.push(t);
          });
        }
        const d = {
          name: k,
          uses: "flow",
          group: [gopack].concat(data.group || []),
          with: {
            params: tags.args,
            tokens,
          }
          // todo: roles
        };
        if (data.slots) {
          if (typeof data.slots === "string") {
            d.slots = [data.slots];
          } else {
            d.slots = data.slots;
          }
        }
        const desc = data.desc;
        if (desc) {
          d.desc = {
            label: pascal(k),
            short: desc,
          };
        }
        specs.push(d);
        break;
      }
    }
  }
  m.currGroups.shift();
});

// avoid as many "redefining" errors as we can
m = new Make();
localLang(m);

// console.log(JSON.stringify(specs, 0,2));
specs.forEach((spec) => {
  try {
    m.newFromSpec(spec);
  }
  catch (e) {
    console.log("newFromSpec:", e.message);
  }
});

//
// add stubs ( spec and stub are added by spec.js )
// spec.forEach((spec)=> {
//   if (stub.indexOf(spec.name) >= 0) {
//   m.newFromSpec(spec);
//   }
// });
//
const sorted = {};
Object.keys(m.types.all).sort().forEach((key) => {
  sorted[key] = m.types.all[key];
});

module.exports = sorted;
