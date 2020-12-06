// node makeops.js > ../../ephemera/story/iffy_model.go
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes= require('./model.js'); // iffy language file
//console.log(JSON.stringify(allTypes, 0,2 ));
//return

// change lower case to pascal-cased ( golang public )
const pascal= function(name) {
  if (name && name[0]=== '$') {
    name= name.slice(1).toLowerCase();
  }
  const els= name.split('_').map(el=> el.charAt(0).toUpperCase() + el.slice(1));
  return els.join('');
};

const lower= function(name) {
  return name.slice(1).toLowerCase();
}

// given a strType with the specified pascal'd name, return its list of choice
const strChoices= function(name, strType) {
  const out=[];
  const { with : {params= {}}= {} } = strType;
  if (params) {
    for (const k in params) {
      const p=lower(k);
      if (p === name) {
        out.unshift(p); // move the dynamic key to the front
      } else {
        out.push(p);
      }
    }
  }
  return out;
};

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('Lower', lower);

// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0]=== '$');
});

// characters preceding a type declaration
  // "label": "trait",
  // "type": "trait",
  // "optional": true,
  // "repeats": true,
  // "filters": [
  //   "comma-and"
  // ]
Handlebars.registerHelper('Lede', function(param) {
  let out = "";
  const name = param.type;
  const type = allTypes[name];
  if (param.optional) {
    out+= "*";
  }
  if (param.repeats) {
    out+= "[]";
  }
  out+= (name.indexOf("_eval") >= 0) ? "rt." :"";
  return out;
});

Handlebars.registerHelper('Tail', function(param) {
  return "";//param.optional? ' `if:"optional"`': "";
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

// for uses='str'
Handlebars.registerHelper('IsClosed', function(strType) {
  const name= lower(strType.name);
  const cs= strChoices(name, strType);
  return cs.length && cs[0] !== name;
});

// for uses='str'
Handlebars.registerHelper('Choices', function(strType) {
  const name= lower(strType.name);
  const cs= strChoices(name, strType);
  return cs[0]===name? cs.slice(1): cs; // remove the dynamic key
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

// flatten groups
Handlebars.registerHelper('GroupOf', function (desc) {
  return desc.group.join(', ');
})

// load each js file as a handlebars template
const partials= ['spec'];
const sources= ['header', 'num', 'txt', 'opt', 'run', 'str', 'slot', 'footer'];
partials.forEach(k=> Handlebars.registerPartial(k, require(`./templates/${k}Partial.js`)));
const templates= Object.fromEntries(sources.map(k=> [k,
Handlebars.compile(require(`./templates/${k}Template.js`))]));

console.log(templates.header({package:'story'}));

// switch to partials?
for (const typeName in allTypes) {
  const type= allTypes[typeName];
  const mytemp= templates[type.uses];
  if (mytemp) {
    console.log(mytemp(type));
  }
}

console.log(templates.footer({package:'story', allTypes}));
