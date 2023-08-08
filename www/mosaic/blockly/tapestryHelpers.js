import Blockly from "blockly";
import ShapeReader from "./shapeReader.js";
import ShapeWriter from "./shapeWriter.js";

// common helper functions added to each Block.
// requires the tapestryMutator ( for showMutatorIcon )
const tapestryHelpers = {
  // extra state tracks our *desired* block appearance
  // fix? while ideally this would only exist for mutable fields;
  // the block generator doesnt always distinguish between optional and required fields.
  getExtraState(name, defaultVal = 0) {
    return this.extraState[name] || defaultVal;
  },
  setExtraState(name, itemCount) {
    return (this.extraState[name] = itemCount);
  },
  // item state tracks our actual block appearance
  // ( everything starts at 0, not present. )
  getItemState(name) {
    return this.itemState[name] || 0;
  },
  setItemState(name, itemCount) {
    return (this.itemState[name] = itemCount);
  },
  // extraState: {}, // created by the generic extension "constructor"...
  // itemState: {},  // so they dont wind up shared across instances.

  // update the workspace block based on its current desired state
  // ( called from load and the mutation's compose )
  updateShape_: function () {
    const blockType = this.type;
    const customData = this.shapeData;
    const shapeDef = customData.shapeDef;
    //
    const shapeReader = new ShapeReader(this);
    const shapeWriter = new ShapeWriter(this, (item, field) => {
      this.onItemAdded(item, field);
    });
    // traversing in tapestry's shapeDef order to track the desired position of fields.
    // this alg handles some cases where extraState has info where its not needed.
    // ( ex. the block generator records the number of stacked inputs even tho that's wrong. )
    shapeDef.forEach((itemDef) => {
      const itemName = itemDef.name;
      let want, have;
      if (!itemName) {
        // we assume we have 1 of all non-notable items, default to zero of everything else.
        (want = 1), (have = 1);
      } else {
        have = this.getItemState(itemName);
        want = this.getExtraState(itemName, itemDef.optional ? 0 : 1);
        // fix? stacks can report > 1 elements due to the block generator on the golang side.
        if (!itemDef.repeats && want > 1) {
          want = 1;
        }
        this.setItemState(itemName, want);
      }
      // transfer the left/top most items first.
      // ( transfer items even if counts match: mutation can cause items to move b/t inputs )
      const transfer = want > have ? have : want;
      if (transfer) {
        shapeReader.takeItems(itemName, have, (item) => {
          if (item.field) {
            shapeWriter.transferField(item.field);
          } else {
            shapeWriter.transferInput(item.input);
          }
        });
      }
      const destroy = have - want;
      if (destroy > 0) {
        shapeReader.takeItems(itemName, destroy, (item) => {
          item.removeFrom(this);
        });
      }
      const add = want - have;
      if (add) {
        // to keep the writer and reader in sync:
        // push the reader forward by the number of new inputs created.
        const inputCount = shapeWriter.createItems(itemDef, have, add);
        shapeReader.adjust(inputCount);
      }
    });
    shapeWriter.finalize();
  },

  // disconnect incompatible blocks after the author selects a new option from a swap's dropdown.
  // [ this is a custom validator callback bound in createInput ]
  onSwapChanged(inputName, swaps, newValue) {
    const input = this.getInput(inputName);
    const targetConnection =
      input.connection && input.connection.targetConnection;
    if (targetConnection) {
      // blockly function which returns a list of compatible value types
      // ( null if everything is allowed; rare. )
      const checks = targetConnection.getCheck();
      const stillCompatible = !checks || checks.includes(swaps[newValue]);
      if (!stillCompatible) {
        targetConnection.disconnect(); // bumps the disconnecting block away automatically.
      }
    }
  },

  // blockly change notification ( responding to changes made by authors )
  // fix? it might be better to have a custom workspace level dispatch.
  // blockly sends every event to every block :/
  onchange: function (e) {
    // show and hide the mutation
    if (e.type == Blockly.Events.SELECTED) {
      if (this.id === e.newElementId && !this.mutator) {
        this.showMutatorIcon(true);
      } else if (this.id == e.oldElementId && this.mutator) {
        // dont hide the mutator bubble due to deselection
        if (!this.mutator.isVisible()) {
          this.showMutatorIcon(false);
        }
      }
      // oldElementId // blockId or workspaceCommentId
      // newElementId // null for deselect
      // isUiEvent // true if a scroll, click, select, block drag.
      // workspaceId // workspace = Blockly.Workspace.getById(event.workspaceId).
      // blockId // block= workspace.getBlockById(event.blockId); but seems to be nil most of the time.
      // group  used for grouping events together for undo/redo. ( such as inserting a statement in a stack.  )
    } else if (e.type == Blockly.Events.BUBBLE_OPEN) {
      // closing the block?
      if ((this.id === e.blockId) && (!e.isOpen) && (e.bubbleType === 'mutator')) {
        const isSelected= Blockly.common.getSelected() === this;
        if (!isSelected) {
           this.showMutatorIcon(false);
        }
      }
    }
    // move events are used for connections
    // this searches for swaps, and changes the swap to match the input type.
    else if (e.type === Blockly.Events.BLOCK_MOVE) {
      // we only care about events intended for this specific block
      if (e.newParentId === this.id) {
        // find the input that's being connected to
        const input = this.getInput(e.newInputName);
        if (input && input.itemDef) {
          // ends when field is nil because array is exhausted.
          for (let i = 0, field; (field = input.fieldRow[i]); i++) {
            if (field.name === input.name) {
              this._updateSwap(input, field);
              break;
            }
          }
        }
      }
    }
  },

  // update the value of a swap based on its current (new) input block.
  // ( ie. when an author tries to connect a block to a swappable input
  //   change the swap's combo box to match. )
  _updateSwap: function (input, field) {
    // if the input is a swap:
    const swaps = field.itemDef && field.itemDef.swaps;
    if (swaps) {
      const targetConnection =
        input.connection && input.connection.targetConnection;
      // assuming something is connected
      if (targetConnection) {
        // first check the current value ....
        const currOpt = field.getValue();
        const currType = swaps[currOpt];
        // blockly function which returns a list of compatible value types
        // ( null if everything is allowed; rare. )
        const checks = targetConnection.getCheck();
        const stillCompatible = !checks || checks.includes(currType);
        if (!stillCompatible) {
          for (let k in swaps) {
            const checkType = swaps[k];
            if (checks.includes(checkType)) {
              field.setValue(k); // this triggers a "swapChanged" :/
              break;
            }
          }
        }
      }
    }
  },
  // callback from ShapeWriter whenever an input is created or a field is added to an input ( for this block )
  onItemAdded(input, field) {
    if (field) {
      // add a validator to disconnect incompatible blocks after a combo box change.
      const swaps = field.itemDef && field.itemDef.swaps; // ex. { "$KINDS": "plural_kinds", "$NOUN": "named_noun" }
      if (swaps && field instanceof Blockly.FieldDropdown) {
        // use bind to give the callback some helpful parameters.
        // ( https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_objects/Function/bind )
        field.setValidator(this.onSwapChanged.bind(this, input.name, swaps));
      }
    } else if (input && input.itemDef) {
      // means that we created something other than a dummy input
      // and we should give it an initial value
      const itemDef = input.itemDef; // this can be null; tbd: maybe on close?
      const isToolbox = this.workspace !== Blockly.mainWorkspace;
      if (isToolbox && itemDef.shadow) {
        const sub = this.workspace.newBlock(itemDef.shadow);
        // guess at a bunch of random things to make this show up correctly.
        // render is needed otherwise drag crashes trying to access a null location,
        // and initSvg is needed before render.
        // shadow is needed for the toolbox otherwise we get a random extra block when drag starts.
        sub.initSvg(); // needed before render is called
        sub.render(false); // false means: only re/render this block.
        sub.setShadow(true); // shadow cleans up better when done; but you can connect other values
        // sub.setDeletable(false);
        // sub.setMovable(false);
        input.connection.connect(sub.outputConnection);
      }
    }
  },
};

// exports a function for use with 'Blockly.extensions.register'
// ( normally user code would use the special 'registerMixin' -- not sure why that's the recommended practice.
// this way makes the code uniform, which helps with safeRegister )
export default function () {
  // "this" here is the Block instance
  this.mixin(tapestryHelpers);
}
