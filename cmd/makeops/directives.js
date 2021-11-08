'use strict';

const TagParser = require('./tags.js');
const TypeSet = require('./typeSet.js');

// change to pascal-cased ( golang public )
const pascal = function(name, sep = " ") {
  const els = name.split('_').map(el => el.charAt(0).toUpperCase() + el.slice(1));
  return els.join(sep);
};


class Make {
  constructor() {
    this.types= new TypeSet();
    this.currGroups= [];
  }

  // this.... needs some work.
  readSpec(spec, group= null) {
    if (group) {
      this.currGroups.unshift(group);
    }
    this._readSpec(spec);
    if (group) {
      this.currGroups.shift();
    }
  }
  _readSpec(dict) {
    for (const k in dict) {
      const data = dict[k];
      // ick.
      const prevGroups= this.currGroups.length;
      if (data.group) {
        this.currGroups= this.currGroups.concat( Array.isArray( data.group )? data.group: data.group );
      }
      switch (data.uses) {
        default:
          throw new Error(`unknown or missing uses ${k}'`);
        case "group": {
          const p = data.specs;
          this.group(k, data.desc || k, () => {
            this._readSpec(p);
          });
          break;
        }
        case "str": {
          this.str(k, data.spec, data.desc);
          break;
        }
        case "num": {
          this.num(k, data.spec, data.desc);
          break;
        }
        case "slot": {
          this.slot(k, data.desc);
          break;
        }
        case "swap": {
          this.swap(k, data.spec, data.desc);
          break;
        }
        case "flow": {
          if (this.currGroups.includes("modeling")) {
            this.flow(k, data.slot || [], data.spec, data.desc || "");
            break;
          }
          // note: doesnt use the "maker" for now
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
            group: this.currGroups.slice(),
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
          this.newFromSpec(d);
          break;
        }
      }
      // remove groups
      const rub= this.currGroups.length - prevGroups;
      if (rub) {
        this.currGroups.splice(prevGroups, rub);
      }
    }
  }


  // introduce the passed group name to types
  // created during the passed function.
  group(name, ...descFn) {
    var desc, fn;
    const [a,b]= descFn;
    if (b===undefined) {
      fn= a;
      desc= "";
    } else {
      desc= a;
      fn= b;
    }
    const n= name.toLowerCase();
    this.currGroups.push(n);
    if (!this.types.has(n)) {
      this.newType( n, "group", desc );
    } else if (desc) {
      throw new Error(`group ${name} already declared`);
    }
    fn();
    this.currGroups.pop();
  }

  // typeName string, [ slot or slots string(s) ], msg string, desc multi-part string.
  // msg is a "format string" -- with token, types, etc.
  // desc is a description with the form: "label: short description. long description...".
  flow(name, ...slotMsgDesc) {
    const [b,c,d]= slotMsgDesc;
    var tags, slots, desc;
    // assume msg is the first parameter
    const firstTags= TagParser.parse(b);
    if (Object.keys(firstTags.args).length) {
      tags= firstTags;
      desc= c; // so desc is second
    } else {
      // for backwards compat with directiveTests; check that c had text before assigning tags
      const secondTags= TagParser.parse(c);
      if (!d && c && !Object.keys(secondTags.args).length) {
        tags= firstTags;
        desc= c;
      } else {
        const slot_or_slots= b;
        tags= secondTags;
        slots= Array.isArray(slot_or_slots)? slot_or_slots: [slot_or_slots],
        desc= d;
      }
    }
    return this.newType(name, "flow", desc,
      // using object assign in case slots dont exist.
      Object.assign({
        tokens: tags.keys,
        params: tags.args
      }, slots && {slots}));
  }

  slot( name, desc= null ) {
    return this.newType(name, "slot", desc);
  }

  // displays types inline ( vs. slot and slat dropdowns )
  swap( name, msg, desc= null ) {
    const tags= TagParser.parse(msg);
    return this.newType(name, "swap", desc, {
        tokens: tags.keys,
        params: tags.args
    });
  }
  str( name, msg=null, desc= null ) {
    return this.makeStr(name, "str", msg, desc);
  }

  // pick or enter a small bit of text.
  makeStr( name, uses, msg, desc) {
    // strTypes dont have normal params:
    // we aren't fitting types into parameters
    // the parameters are instead enumerated choices.
    // ( that should probably be made clearer somehow... )
    const settings={ asValues: true, nullValue: name };
    let tags= TagParser.parse(msg, settings);
    // msg had no tags: it's either a desc, or it was set/left as null.
    if (!Object.keys(tags.args).length) {
      desc= desc || msg;
      msg = `{${name}}`;
      tags= TagParser.parse(msg, settings);
    }
    return this.newType(name, uses, desc, {
        tokens: tags.keys,
        params: tags.args
     });
  }

  num( name, desc= null ) {
    return this.newType(name, "num", desc);
  }

  newType(name, uses, desc, withspec=null) {
    const group= this.currGroups.slice();
    return this.types.newType(Object.assign(
      {name:name},
      desc&&{desc:Make.makeDesc(name, desc)},
      uses&&{uses},
      group&&{group},
      withspec&&{with:withspec}
    ));
  }

  // label, a human readable name
  // short, a short description
  // long, additional details
  static makeDesc(name, desc) {
    var ret;
    if (typeof desc !== 'string') {
      ret= desc;
    } else {
      let label= "";
      let short= "";
      let long= "";
      const i= desc.indexOf('.');
      if (i >= 0) {
        long= desc.substring(i+1).trimLeft();
        desc= desc.substring(0,i+1);
      }
      const j= desc.indexOf(':');
      if (j < 0) {
        short= desc;
      } else {
        label=  desc.substring(0,j);
        short= desc.substring(j+1).trimLeft();
      }
      ret= (label || long)? {
        label: (label || name),
        short,
        long
      }: short;
    }
    return ret;
  }

  // given a text description, add a new type
  // used by autogenerated/autogenerating types.
  newFromSpec(spec) {
    const d= Object.assign(spec,
      spec.desc&&{desc:Make.makeDesc(spec.name, spec.desc)},
    );
    if (d.spec) {
      const tags= TagParser.parse(d.spec);
      const w= d.with || {};
      w.tokens= tags.keys;
      w.params= tags.args;
      d["with"]= w;
      // not sure why this is being deleted
      // maybe just for better logging.
      delete d.spec;
    }
    this.types.newType(d);
  }
}


module.exports = Make;
