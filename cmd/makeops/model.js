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

const parseData = function(m, dict) {
  for (const k in dict) {
    const data = dict[k];
    // ick.
    const prevGroups= m.currGroups.length;
    if (data.group) {
      m.currGroups= m.currGroups.concat( Array.isArray( data.group )? data.group: data.group );
    }
    switch (data.uses) {
      default:
        throw new Error(`unknown or missing uses ${k}'`);
      case "group": {
        const p = data.specs;
        m.group(k, data.desc || k, function() {
          parseData(m, p);
        });
        break;
      }
      case "str": {
        m.str(k, data.spec, data.desc);
        break;
      }
      case "num": {
        m.num(k, data.spec, data.desc);
        break;
      }
      case "slot": {
        m.slot(k, data.desc);
        break;
      }
      case "swap": {
        m.swap(k, data.spec, data.desc);
        break;
      }
      case "slat": {
        m.slat(k, data.slot || [], data.spec, data.desc || "");
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
          // console.log("no keys for", k);
          // happens for things like "always", "equal_to", etc.
        } else {
          tags.keys.forEach((t, i) => {
            if (t.startsWith("$")) {
              const el = tags.args[t];
              if (!el) {
                console.log(`couldnt find ${t} in ${k}`);
              } else {
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
              }
            }
          });
        }
        const d= {
          name: k,
          uses: "flow",
          group: m.currGroups.slice(),
          with: {
            params: tags.args,
              tokens,
          }
          // todo: roles
        };
        if (data.slot) {
          if (typeof data.slot === "string") {
            d.with.slots = [data.slot];
          } else {
            d.with.slots = data.slot;
          }
        }
        const desc = data.desc;
        if (desc) {
          d.desc = {
            label: pascal(k),
            short: desc,
          };
        }
        m.newFromSpec(d);
        break;
      }
    }
    // remove groups
    const rub= m.currGroups.length - prevGroups;
    if (rub) {
      m.currGroups.splice(prevGroups, rub);
    }
  }
};

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
  const gopack = fn.substring(0, fn.length - ".ifspec".length);
  console.log(`parsing ${path}`);
  m.currGroups.unshift(gopack);
  parseData(m, dict);
  m.currGroups.shift();
});

const sorted = {};
Object.keys(m.types.all).sort().forEach((key) => {
  sorted[key] = m.types.all[key];
});

module.exports = sorted;
