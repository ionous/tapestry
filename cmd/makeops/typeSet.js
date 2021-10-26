'use strict';

// inner class for Types
// while a single global Type class simplifies code, it hurts testing.
// this provides a way to have a mixture of both.
class TypeSet {
  constructor() {
    this.all= {};
    this.slots= {}; // slot name => [ runs that implement the slot ]
    this.groups= {}; // group name => [ runs that implement the group ]
  }
  get(typeName) {
    return this.all[typeName];
  }
  has(typeName) {
    return !!this.get(typeName);
  }
  // implements but isnt the same type as...
  implements(doesThis, implementThat) {
    const t= this.get(doesThis);
    return t && (t.with && t.with.slots && t.with.slots.indexOf(implementThat)>=0);
  }
  // implements or is the stame type as...
  areCompatible(doesThis, implementThat) {
    const t= this.get(doesThis);
    return t && ((t.name === implementThat) ||
                ((t.with && t.with.slots && t.with.slots.indexOf(implementThat)>=0)));
  }
  // object { string name; string uses;
  //  union { string; object { label, short, long string } } desc;
  //  object with?; }
  newType(type) {
    const name= type.name;
    if (name in this.all) {
      throw new Error(`redefining type ${name}`);
    }
    this.all[ name ]= type;
    return type;
  }
  newItem(typeName, value) {
    if (!(typeName in this.all)) {
      throw new Error(`expected type, got '${typeName}'`);
    }
    return { type:typeName, value };
  }
}


module.exports = TypeSet;
