

class Make {
  constructor(types) {
    this.types= types.allTypes;
    this.currGroups= [];
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
  run(name, ...slotMsgDesc) {
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
    return this.newType(name, "run", desc,  {
        slots: slots,
        tokens: tags.keys,
        params: tags.args,
    });
  }

  slot( name, desc= null ) {
    return this.newType(name, "slot", desc);
  }

  // displays types inline ( vs. slot and slat dropdowns )
  opt( name, msg, desc= null ) {
    const tags= TagParser.parse(msg);
    return this.newType(name, "opt", desc, {
        tokens: tags.keys,
        params: tags.args
    });
  }
  str( name, msg=null, desc= null ) {
    return this.makeStr(name, "str", msg, desc);
  }

  // pick or enter a small bit of text.
  makeStr( name, uses, msg, desc) {
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

  // multiline text
  txt( name, msg=null, desc= null ) {
    return this.makeStr(name, "txt", msg, desc);
  }

  num( name, desc= null ) {
    return this.newType(name, "num", desc);
  }

  newType(name, uses, desc, withspec=null) {
    const group= this.currGroups.length && this.currGroups;
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
}

