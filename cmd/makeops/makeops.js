// node makeops.js
// reads from spec.js to generate golang models
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes = require('./model.js'); // iffy language file
const fs = require('fs'); // filesystem for loading iffy language file
const child_process = require('child_process');
const path = require('path');
const fnv1a = require('./fnv1a.js');

// unbox: when generating some kinds of simple types...
// replace the specified typename with specified primitive.
// types that map to numbers, etc. are added as unbox automatically.
const unbox = { "text": "string", "bool": "bool" };

// change to tokenized like name
const tokenize = function(name) {
  return '$' + name.toUpperCase();
};

// change to lower case name
const lower = function(name) {
  if (name && name[0] === '$') {
    name = name.slice(1);
  }
  return name? name.toLowerCase(): "xxxx";
};

// change to pascal-cased ( golang public )
const pascal = function(name) {
  const els = lower(name).split('_').map(el => el.charAt(0).toUpperCase() + el.slice(1));
  return els.join('');
};

const strChoices = function(token, strType) {
  const out = [];
  const { with: { params = {}, tokens = [] } = {} } = strType;
  return tokens.filter(t => t[0] == '$' && t !== token).map((t, i) => {
    const d = params[t];
    return { label: d.label || d, index: i, token: t, value: d.value || d };
  });
};

const isClosed = function(strType) {
  const token = tokenize(strType.name);
  const { with: { tokens = [] } = {} } = strType; // safely extract tokens
  return tokens.indexOf(token) < 0;
};

const isPositioned = function(t) {
  return Array.isArray(t.group) ? t.group.includes("positioned") : t.group === "positioned";
};

const paramsOf = function(t) {
  let { with: { params = {} } = {} } = t; //safely extract params
  if (isPositioned(t)) {
    params = Object.assign({ "$AT": { "type": "position", "label": "-" } }, params);
  }
  return params;
};

const groups = {};
const nameToGroup = {};
let currentGroup;

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('IsPositioned', isPositioned);
// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0] === '$');
});
Handlebars.registerHelper('NoHelpers', function(name) {
  return name === "position";
});
// generate a custom compact format command name for a type
const ledeName = function(t) {
  // note: we dont give plain english commands a special lede
  // ( those are the story commands not marked as "modeling" )
  let ret;
  if (t.uses === "flow") {
    const plainEnglish = t.group.includes("story") && !t.group.includes("modeling");
    if (plainEnglish) {
      ret= t.lede;
    } else {
      // get the first token
      const lede = t && t.with && t.with.tokens && t.with.tokens.length > 0 && t.with.tokens[0];
      // and if the first token is english text, then use that as the lede
      ret= (lede && lede.length > 0 && lede[0] !== "$" && lede !== t.name) ? lede : "";
    }
  }
  return ret;
}
Handlebars.registerHelper('LedeName', ledeName);

const scopeOf = function(name) {
  let n = "";
  const g = nameToGroup[name];
  if (g && g !== currentGroup) {
    n = `${g}.`;
  }
  return n;
};

const scopedName = function(name, ignoreUnboxing) {
  let n = pascal(name);
  const unboxType = !ignoreUnboxing && unbox[name];
  if (unboxType) {
    n = unboxType; // stuffed in during makeops startup.
  } else {
    n = scopeOf(name) + n;
  }
  return n;
};

Handlebars.registerHelper('ScopeOf', scopeOf);

Handlebars.registerHelper('LowerNameOf', function(key, param) {
  const el = pascal(key) || pascal(param.type);
  return el.charAt(0).toLowerCase() + el.slice(1);
});
Handlebars.registerHelper('IsBool', function(name) {
  return unbox[name] === 'bool';
});
Handlebars.registerHelper('Unboxed', function(name) {
  return unbox[name];
});
const isPrim= function(type) {
  return ["num", "str"].includes(type.uses);
}
Handlebars.registerHelper('TypeOf', function(param) {
  const name = param.type;
  const type = allTypes[name]; // the referenced type
  if (!type && !param.internal) {
    throw new Error(`unknown type ${name}`);
  }
  //
  let qualifier = "";
  if (param.repeats) {
    qualifier += "[]";
  } else if (param.optional && type.uses === "flow") {
    // re: slot, we dont need *interface{}
    // and its okay enough i think if we accidentally collapse empty strings/numbers into unspecified values
    qualifier += "*";
  }
  return qualifier + scopedName(name);
});

Handlebars.registerHelper('IsSlat', function(name) {
  const { uses } = allTypes[name];
  return uses !== 'slot' && uses !== 'group';
});

Handlebars.registerHelper('Uses', function(name, test) {
  const { uses } = allTypes[name];
  return uses === test;
});

Handlebars.registerHelper('IsInternal', function(label) {
  return label === '-';
});

// for uses='str'
Handlebars.registerHelper('IsClosed', function(strType) {
  return isClosed(strType);
});

Handlebars.registerHelper('IsEnumerated', function(type) {
  let okay= false;
  if (type.uses === 'str') {
    const token = tokenize(type.name);
    okay= strChoices(token, type).length > 0;
  }
  return okay;
});
// for uses='str' or 'swap'
Handlebars.registerHelper('Choices', function(strType) {
  const token = tokenize(strType.name);
  return strChoices(token, strType);
});
// flatten desc
Handlebars.registerHelper('DescOf', function(x) {
  let ret = '';
  if (x.desc) {
    const desc = x.desc;
    if (typeof desc == 'string') {
      ret = desc;
    } else if (desc) {
      ret = pascal(desc.label || x.name);
      const rest = ((desc.short || '') + ' ' + (desc.long || '')).trim();
      if (rest) {
        ret += ': ' + rest;
      }
    }
  }
  return ret;
});

const camelize= function(sel) {
  if (sel && sel.length > 1) {
    sel= pascal(sel);
    sel= sel.charAt(0).toLowerCase() + sel.slice(1);
  }
  return sel || "";
}

// loop over a subset of parameters generating signatures for them recursively
// out is an array of signatures, each signature an array of parts.
const sigParts = function(t, commandName) {
  const ps= t.params.filter(p => !p.internal);
  let sets= [[commandName]];
  ps.forEach((p) => {
    const sel= camelize(p.sel);
    const pt= allTypes[p.type];
    const simpleSwap= !p.repeats && pt.uses === 'swap';
    if (!simpleSwap) {
      const rest= sets.map(a => a.concat(sel));
      if (!p.optional) {
        sets= rest;
      } else {
        sets = sets.concat(rest);
      }
    }else {
      // every choice in a swap gets its own selector for each existing set
      let mul= [];
      pt.params.filter(c => !c.internal).forEach((c)=> {
        const choice= `${sel} ${camelize(c.sel)}`;
        mul= mul.concat(sets.map(a=> a.concat(choice)));
      });
      sets= mul;
    }
  });
  return sets;
}

// generate a signature hash for the passed type.
const signType = function(t, out, all) {
  const lede= ledeName(t);
  const commandName= pascal(lede || t.name);
  const sigs= [];
  if (t.uses === 'swap') {
    for (const p of t.params) {
      if (!p.internal) {
        const sel= camelize(p.sel);
        sigs.push( commandName + " " + sel + ":" );
      }
    }
  } else {
    // each part is an array of selectors
    const sets= sigParts(t, commandName);
    // console.log("sets", sets);
    // we need to reduce those arrays into strings
    for (const set of sets) {
      let sig= set[0]; // index 0 is the command name itself
      if (set.length > 1) {
        // the first parameter doesnt always have selector text
        // ( and sometimes it might have some choice text with a leading (though possibly also blank) selector )
        const p1= set[1].trimLeft();
        if (p1) {
          sig += ' ' + p1 + ':';
        }
        // separate all other parameters:
        const rest= set.slice(p1?2:1);
        if (rest.length) {
          sig += rest.join(':') + ':';
        }
      }
      sigs.push(sig);
    }
  }

  for (const n of sigs) {
    if (n.includes("::") || n.includes("_") || n.includes(": ")) {
      console.log(commandName, n);
      throw new Error(n);
    }
    const hash= fnv1a(n, {size:64});
    const was= all[hash];
    if (was) {
      // Error: hash conflict 14222179384726139225 'FromRec:' from_rec vs. 'FromRec:' from_record
      throw new Error(`hash conflict '${was.sig}' ${was.group}.${was.type} vs. '${n}' ${t.group}.${t.name}`)
    }
    const d= {
      sig:n,
      type:t.name,
      group:t.group,
    };
    out[hash]= all[hash]= d;
  }
};

const locationOf = function(x) {
  let where;
  if (x.includes("/")) {
    where = x;
  } else {
    switch (x) {
      case "rt":
        where = `git.sr.ht/~ionous/iffy/rt`;
        break;
      case "story":
        // FIX: should move all the story files to the dl folder instead.
        where = `git.sr.ht/~ionous/iffy/ephemera/story`;
        break;
      default:
        where = `git.sr.ht/~ionous/iffy/dl/${x}`
        break;
    }
  }
  return where;
}
Handlebars.registerHelper('LocationOf', locationOf);

// flatten groups
Handlebars.registerHelper('GroupOf', function(desc) {
  return desc.group.join(', ');
})
// load each js file as a handlebars template
const partials = [
  'repeat', 'sig', 'spec'
];
const sources = ['header', 'slot', 'prim', 'swap', 'flow', 'footer', 'regList', 'sigMap'];
partials.forEach(k => Handlebars.registerPartial(k, require(`./templates/${k}Template.js`)));
const templates = Object.fromEntries(sources.map(k => [k,
  Handlebars.compile(require(`./templates/${k}Template.js`))])
);
templates["str"] = templates["num"] = templates["prim"];

// split types into different categories
for (const typeName in allTypes) {
  const type = allTypes[typeName];
  // fix: maybe carry through the lines the whole way?
  if (Array.isArray(type.desc)) {
    type.desc = type.desc.join("  ");
  }
  // write over existing parameter data
  let ps=[];
  const params= paramsOf(type);
  let pt= 0; // non-internal total
  for (const key in params) {
    let param= params[key];
    // str types can have params which are just a string
    // whenever the label and the value are the same value
    // to keep the template generator simple: normalize those params.
    if (typeof param === 'string') {
      params[key]= param= { label: param, value: param };
    }
    param.key= key;
    param.internal= param.label === '-';
    ps.push(param);
    if (!param.internal) {
      pt++;
    }
  }
  let pi=0; // non-internal indices
  ps.forEach((param,i)=> {
    // story commands are assumed to be written in plain english, except for those few tagged modeling.
    // all of those plain english commands get anonymous first parameters
    // unless that first parameter is optional and there are other trailing parameters.
    // all of the other commands -- the fluent commands -- control their parameters via their spec.
    let tag= lower(param.key);
    if (type.uses !== 'swap') {
      const plainEnglish = type.group.includes("story") && !type.group.includes("modeling");
      const anon= !pi && (!param.optional || pt === 1);
      tag= plainEnglish ? (anon ? '_': tag) : param.label.replaceAll(" ", "_");
    }
    param.tag= tag;
    param.sel= tag !== "_" ? tag: "";
    if (!param.internal) {
      pi++;
    }
  });
  type.params= ps;

  //
  if (type.uses !== "group") {
    // ex. ["story statements"]=> "story"
    const group = type.group[0].split(" ")[0];
    nameToGroup[typeName] = group;
    const g = groups[group] || { slots: [], slats: [], all: [] };
    if (type.uses === "slot") {
      g.slots.push(typeName);
    } else if (type.uses !== "group") {
      //
      if (type.uses === "num") {
        const { with: { tokens = [] } = {} } = type; // safely extract tokens
        if (tokens.length <= 1) {
          unbox[typeName] = "float64";
        }
      }
      g.slats.push(typeName);
    }
    g.all.push(typeName);
    groups[group] = g;
  }
}
// console.log("num groups", Object.keys(groups).length);
// console.log("all types", JSON.stringify(allTypes,0,2));
let regall= {};

// determine includes:
for (currentGroup in groups) {
  console.log(currentGroup);
  const marshal = currentGroup !== "reader";
  const g = groups[currentGroup];
  // look up all the dependencies
  const inc = new Set();
  for (const typeName of g.slats.filter(n=> !unbox[n])) {
    const type = allTypes[typeName];
    if (isPositioned(type)) {
      const o = "reader"; // for forced str and swap position field
      if (o && o !== currentGroup) {
        inc.add(o);
      }
    }

    for (const param of type.params) {
      // when we are marshaling we need to include all types
      // otherwise we only need to include the types we dont unbox out of existence
      if (param && (marshal || !unbox[param.type])) {
        const o = nameToGroup[param.type];
        if (o && o !== currentGroup) {
          inc.add(o);
        }
      }
    }
  }

  // 1. open a file
  const dir = path.join(process.env.GOPATH, "src", locationOf(currentGroup));
  const filepath = path.join(dir, `${currentGroup}_lang.go`);
  console.log("creating", dir, "with", g.slats.length, "cmds");
  fs.mkdirSync(dir, { recursive: true });
  const fd = fs.openSync(filepath, 'w');
  if (g.slats.length) {
    inc.add("composer");
  }
  if (marshal) {
    inc.add("git.sr.ht/~ionous/iffy/jsn");
    inc.add("github.com/ionous/errutil");
  }

  // 2. write the header ( with package name and inc )
  fs.writeSync(fd, templates.header({
    package: currentGroup,
    imports: Array.from(inc.values()).sort(),
  }));
  // #. write slats ( if any )
  for (const typeName of g.all) {
    const type = allTypes[typeName];
    const template = templates[type.uses];
    if (template) {
      // console.log(typeName, marshal);
      const d= {
        marshal,
        type,
      };
      fs.writeSync(fd, template(d));
    }
  }

  // 3. write registration lists
  fs.writeSync(fd, templates.regList({
    which: "Slots",
    list: g.slots.map(n => allTypes[n]),
    RegType: "interface{}",
  }));
  fs.writeSync(fd, templates.regList({
    which: "Slats",
    list: g.slats.map(n => allTypes[n]),
    RegType: "composer.Composer",
  }));
  let signatures= {};
  // console.log(JSON.stringify(allTypes,0,2));
  for (const t of g.slats.map(n => allTypes[n])) {
    if ((t.uses === 'flow') || (t.uses === 'swap')) {
      signType(t, signatures, regall);
    }
  }
  fs.writeSync(fd, templates.sigMap({
    list: signatures,
  }));

  fs.closeSync(fd);
  // re-format the file using go format.
  // actually, use goimports because its hard sometimes to know where errutil is needed
  child_process.execSync(`goimports -e -w ${filepath}`);
  // child_process.execSync(`gofmt -e -s -w ${filepath}`);
}


