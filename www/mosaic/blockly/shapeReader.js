// deconstructs fields and inputs from a block in a unified way
// ( "reading" is a misnomer because this does modify the shape as it goes )
export default class ShapeReader {
	constructor(b) {
		this.block= b;
		this.next= 0;
		this.at = this.advance();
	}
  // the number of inputs in the shape have changed.
  adjust(inputCount) {
    this.next += inputCount;
    // we dont have to change our input or fields
    // ex. caller just took 3 repeating items, we're on a new non-repeating item
    // and the caller adds 2 more repeating items
    // our item is the same and so are our fields, there's just more inbetween.
  }
  // call onItem with { el: field|input, removeFrom }
  takeItems(itemName, count, onItem) {
    for (let i = 0; i < count; ) {
      const item = this.getNext();
      if (!item.validate(itemName)) {
        console.warn("couldn't validate item");
        break;
      }
      if (!itemName || item.name()) {
        ++i;
      }
      onItem(item);
    }
  }
  getNext() {
    let item;
    const fields= this.at.fields;
    if (fields.length) {
      item= new FieldEl(fields.shift());
    } else if (this.at.input) {
      item= new InputEl(this.at.input);
      this.at = this.advance();
    } else {
      throw new Error("out of items");
    }
    return item;
  }
	// remove the current fields (without disposing of them)
	advance() {
    let at;
    const inRange = this.next < this.block.inputList.length;
    if (!inRange)  {
      at= { fields: [] };
    } else {
  		const input= this.block.inputList[this.next++];
  		const fields= input.fieldRow;
  		input.fieldRow = [];
  		at= { input, fields };
    }
    return at;
	}
}

class FieldEl {
  constructor(f) {
    this.field= f;
  }
  name() {
    return this.field.name;
  }
  validate(itemName) {
    // all labels will be considered valid for any item name -- we could set itemDef for them and compare maybe?
    return !this.field.name || this.field.name.startsWith(itemName);
  }
  removeFrom(block) {
    this.field.dispose();
  }
}
class InputEl {
  constructor(input) {
    this.input= input;
  }
  name() {
    return this.input.name;
  }
  validate(itemName) {
    // all inputs should have names
    return itemName && this.input.name.startsWith(itemName);
  }
  removeFrom(block) {
    block.removeInput(this.input.name);
  }
}