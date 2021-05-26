// node makeops.js | gofmt -e -s > ../../ephemera/story/iffy_model.go
// reads from spec.js to generate golang models
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes= require('./model.js'); // iffy language file
const fs = require('fs'); // filesystem for loading iffy language file
const child_process= require('child_process');
const path= require('path');

// change to tokenized like name
const tokenize= function(name) {
  return '$'+ name.toUpperCase();
};

// change to lower case name
const lower= function(name) {
  if (name && name[0]=== '$') {
    name= name.slice(1);
  }
  return name.toLowerCase();
};

// change to pascal-cased ( golang public )
const pascal= function(name) {
  const els= lower(name).split('_').map(el=> el.charAt(0).toUpperCase() + el.slice(1));
  return els.join('');
};

const strChoices= function(token, strType) {
  const out=[];
  const { with : { params= {}, tokens=[]}= {} } = strType;
  return tokens.filter(t=>t[0]=='$' && t!==token).map((t,i)=> {
    const d= params[t];
    return { label: d.label || d , index:i, token:t, value: d.value || d };
  });
};

const isClosed= function(strType) {
  const token= tokenize(strType.name);
  const { with : { tokens=[]}= {} } = strType; // safely extract tokens
  return tokens.indexOf(token) < 0;
};

const groups= {};
const nameToGroup= {};
let currentGroup;

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('Lower', lower);

// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0]=== '$');
});

Handlebars.registerHelper('LedeName', function(t) {
  const lede= t && t.with && t.with.tokens && t.with.tokens.length >0 && t.with.tokens[0];
  return (lede && lede.length > 0 && lede[0] !== "$" && lede !== t.name) ? lede: "";
});

Handlebars.registerHelper('NameOf', function(key, param) {
  return pascal(key) || pascal(param.type);
});

Handlebars.registerHelper('TypeOf', function(param) {
  const name= param.type;
  const type = allTypes[name]; // the referenced type
  if (!type && param.label!=='-') {
    throw new Error(`unknown type ${name}`);
  }
  //
  let n= pascal(name);
  if (type && type.override) {
    n= type.override; // stuffed in during makeops startup.
  } else {
    const g= nameToGroup[name];
    if (g && g!==currentGroup) {
      n= `${g}.${n}`;
    }
  }
  //
  let qualifier = "";
  if (param.repeats) {
    qualifier+= "[]";
  } else if (param.optional && type.uses !== "slot") {
    // re: slot, for go we dont need *interface{}
    qualifier+= "*";
  }
  return qualifier+n;
});

// is the passed name a slot
Handlebars.registerHelper('IsSlot', function(name) {
  const { uses }= allTypes[name];
  return uses === 'slot';
});

Handlebars.registerHelper('IsSlat', function(name) {
  const { uses }= allTypes[name];
  return uses !== 'slot' && uses !== 'group';
});

Handlebars.registerHelper('IsStr', function(name) {
  const { uses }= allTypes[name];
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
  const token= tokenize(strType.name);
  return strChoices(token, strType);
});

// flatten desc
Handlebars.registerHelper('DescOf', function (x) {
  let ret='';
  if (x.desc) {
    const desc= x.desc;
    if (typeof desc == 'string') {
      ret= desc;
    } else if (desc) {
      ret= pascal(desc.label || x.name);
      const rest= ((desc.short || '') + ' '+ (desc.long || '')).trim();
      if (rest) {
        ret+= ': ' + rest;
      }
    }
  }
  return ret;
})

const locationOf= function (x) {
  return (x !== "rt") ? `git.sr.ht/~ionous/iffy/dl/${x}`:  `git.sr.ht/~ionous/iffy/rt`
}
Handlebars.registerHelper('LocationOf', locationOf);

// flatten groups
Handlebars.registerHelper('GroupOf', function (desc) {
  return desc.group.join(', ');
})

// load each js file as a handlebars template
const partials= ['spec'];
const sources= ['header', 'num', 'swap', 'flow', 'str',  'footer', 'regList'];
partials.forEach(k=> Handlebars.registerPartial(k, require(`./templates/${k}Partial.js`)));
const templates= Object.fromEntries(sources.map(k=> [k,
  Handlebars.compile(require(`./templates/${k}Template.js`))])
);
templates['txt']= templates['str']; // FIX: txt shouldnt even exist i think
// console.log(templates.header({package:'story'}));

// split types into different categories
for (const typeName in allTypes) {
  const type= allTypes[typeName];
  //
  let group= type.group;
  if (!group)  {
    // ironically, this happens on groups
    if (type.uses !== "group") {
      console.log("no group", JSON.stringify(type,0,2));
    }
  } else {
    // ex. ["story statements"]=> "story"
    group= group[0].split(" ")[0];
    nameToGroup[typeName]= group;
    let g= groups[group];
    if (!g) {
      g= {
        slots: [],
        slats:[],
      };
    }
    if (type.uses==="slot") {
      g.slots.push(typeName);
    } else if (type.uses!=="group") {
      // do a bunch of work to figure out whether to "expand" the type
      // go is pretty strict about its typedefs, and sometimes its nicer
      // just to have a string instead of a wrapper type requiring string access.
      if (type.uses === "str") {
        const { with : { tokens=[]}= {} } = type; // safely extract tokens
        if (tokens.length == 1) {
          type.override= "string";
        } else {
          const token= tokenize(typeName);
          const closedChoices= tokens.indexOf(token)<0;
          // console.log(name, token, tokens);
          if (closedChoices && Object.keys(type.with.params).length===2) {
            type.override= "bool";
          }
        }
      } else if (type.uses === "num"){
        const { with : { tokens=[]}= {} } = type; // safely extract tokens
        if (tokens.length<= 1) {
          type.override= "float64";
        }
      }
      g.slats.push(typeName);
    }
    groups[group]= g;
  }
}
console.log("num groups", Object.keys(groups).length);
// write by group
for (currentGroup in groups) {
  const g= groups[currentGroup];
  // look up all the dependencies
  const inc=[];
  let count=0;
  for (const n of g.slats) {
    count++;
    const type= allTypes[n];
    const ps= type && type.with && type.with.params;
    if (ps && !type.override) {
      for (const p in ps) {
        const param= ps[p];
        const o= nameToGroup[param.type];
        if (o && o !== currentGroup && inc.indexOf(o)<0) {
          inc.push(o);
        }
      }
    }
  }
  //
  // 1. open a file
  const dir= path.join(process.env.GOPATH, "src", locationOf(currentGroup));
  const filepath= path.join(dir, `${currentGroup}_lang.go`);
  console.log("creating", dir, "with", count, "cmds");
  fs.mkdirSync(dir, { recursive: true });
  const fd= fs.openSync(filepath, 'w');
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
    const type= allTypes[n];
    const template= templates[type.uses];
    if (!template) {
      throw new Error(`unknown template for ${n}`);
    } else if (!type.override) {
      fs.writeSync(fd, template(type));
    }
  }
  // write registration lists
  fs.writeSync(fd, templates.regList({
    which: "Slots",
    list: g.slots.map(n => allTypes[n])
  }));
  fs.writeSync(fd, templates.regList({
    which: "Swaps",
    list: g.slats.map(n => allTypes[n]).filter(t=> t.uses==="swap")
  }));
  fs.writeSync(fd, templates.regList({
    which: "Flows",
    list: g.slats.map(n => allTypes[n]).filter(t=> (t.uses==="flow" && !t.override))
  }));
  fs.closeSync(fd);
  // re-format the file using go format.
  child_process.execSync(`gofmt -e -s -w ${filepath}`);
}
