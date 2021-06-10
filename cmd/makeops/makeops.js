// node makeops.js | gofmt -e -s > ../../ephemera/story/iffy_model.go
// reads from spec.js to generate golang models
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes = require('./model.js'); // iffy language file
const fs = require('fs'); // filesystem for loading iffy language file
const child_process = require('child_process');
const path = require('path');

const overrides= {"text":"string"};

// change to tokenized like name
const tokenize = function(name) {
  return '$' + name.toUpperCase();
};

// change to lower case name
const lower = function(name) {
  if (name && name[0] === '$') {
    name = name.slice(1);
  }
  return name.toLowerCase();
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

const isPositioned= function(t) {
  return Array.isArray(t.group) ? t.group.includes("positioned") : t.group ==="positioned";
};

const paramsOf = function(t) {
  let { with: { params = {} } = {} } = t; //safely extract params
  if (isPositioned(t)) {
    params= Object.assign({"$AT":{"type":"position", "label":"-"}}, params);
  }
  return params;
};

const groups = {};
const nameToGroup = {};
let currentGroup;
let currentType;

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('Lower', lower);
Handlebars.registerHelper('ParamsOf', paramsOf);
Handlebars.registerHelper('IsPositioned', isPositioned);
// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0] === '$');
});

Handlebars.registerHelper('LedeName', function(t) {
  const m= t.group.includes("modeling");
  if (!m && t.uses === "flow") {
    const lede = t && t.with && t.with.tokens && t.with.tokens.length > 0 && t.with.tokens[0];
    return (lede && lede.length > 0 && lede[0] !== "$" && lede !== t.name) ? lede : "";
  }
});
const scopedName= function(name) {
  let n = pascal(name);
  const override= overrides[name];
  if (override) {
    n = override; // stuffed in during makeops startup.
  } else {
    const g = nameToGroup[name];
    if (g && g !== currentGroup) {
      n = `${g}.${n}`;
    }
  }
  return n;
};

Handlebars.registerHelper('ScopedNameOf', scopedName);

Handlebars.registerHelper('ParamNameOf', function(key, param) {
  return pascal(key) || pascal(param.type);
});

Handlebars.registerHelper('LabelOf', function(key, param, index) {
  const m= currentType.group.includes("modeling");
  return m? (index? lower(key): '_') : param.label.replaceAll(" ", "_");
});
Handlebars.registerHelper('Override', function(param) {
  const name = param.type;
  return overrides[name]? name:false;
});
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
  } else if (param.optional && type.uses !== "slot" && type.uses !== "str") {
    // re: slot, for go we dont need *interface{}
    qualifier += "*";
  }
  return qualifier + scopedName(name);
});

// is the passed name a slot
Handlebars.registerHelper('IsSlot', function(name) {
  const { uses } = allTypes[name];
  return uses === 'slot';
});

Handlebars.registerHelper('IsSlat', function(name) {
  const { uses } = allTypes[name];
  return uses !== 'slot' && uses !== 'group';
});

Handlebars.registerHelper('IsStr', function(name) {
  const { uses } = allTypes[name];
  return uses === 'str';
});

Handlebars.registerHelper('IsInternal', function(label) {
  return label === '-';
});

// for uses='str'
Handlebars.registerHelper('IsClosed', function(strType) {
  return isClosed(strType);
});

// for uses='str'
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
  switch (x) {
    case  "rt":
      where= `git.sr.ht/~ionous/iffy/rt`;
      break;
    case "story":
      // FIX: should move all the story files to the dl folder instead.
      where= `git.sr.ht/~ionous/iffy/ephemera/story`;
      break;
    default:
      where= `git.sr.ht/~ionous/iffy/dl/${x}`
      break;
  }
  return where;
}
Handlebars.registerHelper('LocationOf', locationOf);

// flatten groups
Handlebars.registerHelper('GroupOf', function(desc) {
  return desc.group.join(', ');
})

// load each js file as a handlebars template
const partials = ['spec'];
const sources = ['header', 'num', 'swap', 'flow', 'str', 'footer', 'regList'];
partials.forEach(k => Handlebars.registerPartial(k, require(`./templates/${k}Partial.js`)));
const templates = Object.fromEntries(sources.map(k => [k,
  Handlebars.compile(require(`./templates/${k}Template.js`))])
);
templates['phrase'] = templates['flow'];

// console.log(templates.header({package:'story'}));

// split types into different categories
for (const typeName in allTypes) {
  const type = allTypes[typeName];
  // fix: maybe carry through the lines the whole way?
  if (Array.isArray(type.desc)) {
    type.desc = type.desc.join("  ");
  }
  //
  let group = type.group;
  if (!group) {
    // ironically, this happens on groups
    // if (type.uses !== "group") {
    //   console.log("no group", JSON.stringify(type, 0, 2));
    // }
  } else {
    // fix. i dont know.
    if (type.uses === "txt") {
      type.uses= "str";
    }
    // ex. ["story statements"]=> "story"
    group = group[0].split(" ")[0];
    nameToGroup[typeName] = group;
    let g = groups[group];
    if (!g) {
      g = {
        slots: [],
        slats: [],
      };
    }
    if (type.uses === "slot") {
      g.slots.push(typeName);
    } else if (type.uses !== "group") {
      // do a bunch of work to figure out whether to "expand" the type
      // go is pretty strict about its typedefs, and sometimes its nicer
      // just to have a string instead of a wrapper type requiring string access.
      if (type.uses === "str") {
        const { with: { tokens = [] } = {} } = type; // safely extract tokens
        const token = tokenize(typeName);
        const closedChoices = tokens.indexOf(token) < 0;
        // console.log(name, token, tokens);
        if (closedChoices && Object.keys(type.with.params).length === 2) {
          overrides[typeName]= "bool";
        }
      } else if (type.uses === "num") {
        const { with: { tokens = [] } = {} } = type; // safely extract tokens
        if (tokens.length <= 1) {
          overrides[typeName]= "float64";
        }
      }
      g.slats.push(typeName);
    }
    groups[group] = g;
  }
}
console.log("num groups", Object.keys(groups).length);

// determine includes:
for (currentGroup in groups) {
  const g = groups[currentGroup];
  // look up all the dependencies
  const inc = [];
  let count = 0;
  for (const n of g.slats) {
    count++;
    if (!overrides[n]) {
      const type = allTypes[n];
      if (isPositioned(type)) {
        const o= "reader"; // for forced str and swap position field
        if (o && o !== currentGroup && inc.indexOf(o) < 0) {
          inc.push(o);
        }
      }
      const params= paramsOf(type);
      for (const p in params) {
        const param = params[p];
        if (param && !overrides[param.type]) {
          const o = nameToGroup[param.type];
          if (o && o !== currentGroup && inc.indexOf(o) < 0) {
            inc.push(o);
          }
        }
      }
    }
  }

  // for looking at the full spec(s)
  // console.log(JSON.stringify(allTypes,0,2));
  // return;

  //
  // 1. open a file
  const dir = path.join(process.env.GOPATH, "src", locationOf(currentGroup));
  const filepath = path.join(dir, `${currentGroup}_lang.go`);
  console.log("creating", dir, "with", count, "cmds");
  fs.mkdirSync(dir, { recursive: true });
  const fd = fs.openSync(filepath, 'w');
  if (g.slats.length) {
    inc.push("composer");
  }
  // 2. write the header ( with package name and inc )
  fs.writeSync(fd, templates.header({
    package: currentGroup,
    imports: inc.sort(),
  }));
  // #. write slats ( if any )
  for (const n of g.slats) {
    const type = allTypes[n];
    const template = templates[type.uses];
    if (!template) {
      throw new Error(`unknown template for ${n}`);
    } else {
      currentType= type;
      fs.writeSync(fd, template(type));
    }
  }
  // write registration lists
  fs.writeSync(fd, templates.regList({
    which: "Slots",
    list: g.slots.map(n => allTypes[n]),
    RegType: "interface{}",
  }));
  // fs.writeSync(fd, templates.regList({
  //   which: "Swaps",
  //   list: g.slats.map(n => allTypes[n]).filter(t => t.uses === "swap"),
  //   RegType: "interface{}",
  // }));
  fs.writeSync(fd, templates.regList({
    which: "Slats",
    list: g.slats.map(n => allTypes[n]).filter(t => (t.uses !== "slot")),
    RegType: "composer.Composer",
  }));
  fs.closeSync(fd);
  // re-format the file using go format.
  child_process.execSync(`gofmt -e -s -w ${filepath}`);
}
