import Blockly from 'blockly';

const DummyInputName = "$DUMMY"

// unifies adding new blocks or items to a block
export default class ShapeWriter {
  constructor(b, addedItem) {
    this.block= b;
    // "Switch off rendering while the block is rebuilt."
    this.savedRendered = b.rendered;
    this.block.rendered = false;
    this.currentIndex = 0;
    this.pendingFields = new PendingFields(addedItem);
    // callbacks
    this.addedItem = addedItem;
  }
  // add an existing field
  // ( it gets held until the next input or shape finalization )
  transferField(field) {
    this.pendingFields.push((input)=>{
        console.assert(input.sourceBlock_ === field.sourceBlock_, `input ${input.sourceBlock_.id} and field blocks ${field.sourceBlock_.id} don't match`);
        input.fieldRow.push(field);
    });
  }
  // add an existing input
  // ( all pending fields get flushed into it )
  transferInput(input) {
    console.assert(input === this.block.inputList[this.currentIndex], "reader and writer need synchronization");
    this.pendingFields.flushInto(input);
    this.currentIndex++; // we added an input
  }
  // create one or more new items
  createItems(itemDef, ofs, want) {
    let prevIndex = this.currentIndex;
    let input = itemDef.type && itemDef.type.startsWith("input_");
    let field = itemDef.type && itemDef.type.startsWith("field_");
    let label = itemDef.label;
    console.assert(input||field||label, `invalid itemDef ${itemDef}`);
    for (let i=0; i<want; i++) {
      if (label && !ofs) {
        this.pendingFields.push((input)=> {
          input.appendField(label); // blockly turns raw strings into labels automatically.
        });
      }
      if (input) {
        const newInput = createInput(this.block, itemDef, ofs++, this.currentIndex++);
        this.pendingFields.flushInto(newInput);
        this.addedItem({input:newInput});
      } else if (field) {
        const fieldNum = ofs++; // pin for our callback.
        this.pendingFields.push((input)=> {
          const field = Blockly.fieldRegistry.fromJson(itemDef);
          field.itemDef = itemDef; // mainly to give swaps access to their choices.
          const fieldName = itemDef.name + (itemDef.repeats? fieldNum: "");
          input.appendField(field, fieldName); // field names are assigned at append time :'(
        });
      }
    }
    return this.currentIndex-prevIndex; // number of inputs  added.
  }
  finalize() {
    // create or destroy the trailing dummy input ( if it exists )
    if (this.pendingFields.empty()) {
      // pass true to avoid throwing an error if there is no dummy.
      this.block.removeInput(DummyInputName, true);
    } else {
      let dummy= this.block.getInput(DummyInputName);
      if (!dummy) {
        dummy = this.block.appendDummyInput(DummyInputName);
      }
      // add the pending items to it.
      this.pendingFields.flushInto(dummy);
    }
    // "Restore rendering and show the changes."
    this.block.rendered = this.savedRendered;
    if (this.block.rendered) {
      this.block.render();
      // if added items ( todo: maybe a tracker? )
      this.block.bumpNeighbours();
    }
  }
}

// hold "constructors" for adding fields until the very next input
// ( blockly arranges fields into rows of inputs )
class PendingFields {
  constructor(addedItem) {
    this.fields= [];
    this.addedItem= addedItem;
  }
  empty() {
    return !this.fields.length;
  }
  push(f) {
    this.fields.push(f);
  }
  flushInto(input) {
    if (this.fields.length) {
      this.fields.forEach((makeField) => {
        makeField(input);
        this.addedItem(input, input.fieldRow[input.fieldRow.length-1]);
      });
      this.fields= [];
    }
  }
}

// allows overrides of the input name to handle repeated elements
function createInput(block, itemDef, ofs, atIndex) {
  const name = itemDef.name + (itemDef.repeats? ofs: "");
  const newInput = appendInput(block, name, itemDef.type);
  newInput.itemDef = itemDef;
  newInput.align = Blockly.ALIGN_RIGHT;
  // blockly can't create inputs directly at designated locations.
  // so move it to the right place after append ( if it's not already there )
  const inputList = block.inputList;
  if (atIndex < inputList.length-1) {
    // blockly has "moveNumberedInputBefore" but this is more straightforward
    inputList.splice(atIndex, 0, inputList.pop());
  }
  // finally, set any input plug assertions.
  if (itemDef.checks) {
    newInput.setCheck(itemDef.checks);
  }
  return newInput;
}

function appendInput(block, name, inputType) {
  // note: the names "statement_input" etc.
  // only have meaning for the json descriptions --
  // and the public interface doesn't have a generic append
  switch (inputType) {
    case 'input_dummy':
      return block.appendDummyInput(name);
    case 'input_value':
      return block.appendValueInput(name);
    case 'input_statement':
      return block.appendStatementInput(name);
    default:
      throw new Error(`couldn't create input ${name} of ${inputType}`);
  }
}