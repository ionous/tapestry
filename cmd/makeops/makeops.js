// node makeops.js
// reads from spec.js to generate golang models
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes = require('./model.js'); // iffy language file
const fs = require('fs'); // filesystem for loading iffy language file
const child_process = require('child_process');
const path = require('path');

// unbox: when generating some kinds of simple types...
// replace the specified typename with specified primitive.
// types that map to numbers, etc. are added as unbox automatically.
const unbox = { "text": "string", "bool": "bool" };
// custom serialization of a type
// blocks the generation of the inner most serialization function
const custom = new Set(["get_var", "text"]);

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
let currentType;

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('ParamsOf', paramsOf);
Handlebars.registerHelper('IsPositioned', isPositioned);
// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0] === '$');
});
Handlebars.registerHelper('NoHelpers', function(name) {
  return name === "position";
});

Handlebars.registerHelper('LedeName', function(t) {
  // we exclude modeling ( as a opposed to runtime functions )
  // because currently most of those are in sentence form not expression form.
  // if this get fixed, then the WriteValue WriteChoice should e fixed to pass a lede
  const m = t.group.includes("modeling");
  if (!m && t.uses === "flow") {
    const lede = t && t.with && t.with.tokens && t.with.tokens.length > 0 && t.with.tokens[0];
    return (lede && lede.length > 0 && lede[0] !== "$" && lede !== t.name) ? lede : "";
  }
});

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

const labelOf= function(key, param, index) {
  const m = currentType.group.includes("modeling");
  return m ? (index ? lower(key) : '_') : param.label.replaceAll(" ", "_");
}
Handlebars.registerHelper('LabelOf', labelOf);
Handlebars.registerHelper('SelectorOf', function(key, param, index) {
  const x= labelOf(key, param, index);
  return x !== "_"? x: "";
});
Handlebars.registerHelper('Custom',function(name) {
  return custom.has(name) ? "_Customized": "";
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
  if (!type && param.label !== '-') {
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

Handlebars.registerHelper('IsLiteral', function(group) {
  return group.indexOf('literals') >= 0;
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
})

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
const sources = ['header', 'slot', 'prim', 'swap', 'flow', 'footer', 'regList'];
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
console.log("num groups", Object.keys(groups).length);

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
    const params = paramsOf(type);
    for (const p in params) {
      const param = params[p];
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
    inc.add("git.sr.ht/~ionous/iffy/export/jsn");
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
      currentType = type;
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
  fs.closeSync(fd);
  // re-format the file using go format.
  child_process.execSync(`gofmt -e -s -w ${filepath}`);
}


