'use strict';

// import Blockly from 'blockly';

let workspace = null;

const jsonDefs= {};

function start() {
  // registerFirstContextMenuOptions();
  // registerOutputOption();
  // registerHelpOption();
  // registerDisplayOption();
  // Blockly.ContextMenuRegistry.registry.unregister('workspaceDelete');
  // Create main workspace.

  jsonData.forEach(function(el) {
    jsonDefs[el.type]= el;
    Blockly.defineBlocksWithJsonArray([el]);
  });
  // Blockly.defineBlocksWithJsonArray(exp);
  Blockly.WorkspaceAudio.prototype.preload= function(){};
  workspace = Blockly.inject('blocklyDiv',{
      // toolbox: toolboxSimple,
      toolbox: {
        'kind': 'flyoutToolbox',
        'contents': [
          // {
          //   'kind': 'block',
          //   'type': 'controls_repeat_ext',
          //   'inputs': {
          //     'TIMES': {
          //       'shadow': {
          //         'type': 'math_number',
          //         'fields': {
          //           'NUM': 5,
          //         },
          //       },
          //     },
          //   },
          // },

          {
            'kind': 'block',
            'type': 'test_str',
          },
          {
            'kind': 'block',
            'type': 'text_value',
          },
          {
            'kind': 'block',
            'type': 'text_values',
          },
          {
            'kind': 'block',
            'type': 'field_values',
          },

          {
            'kind': 'block',
            'type': 'test_flow',
          },

          {
            'kind': 'block',
            'type': 'test_swap',
          },
          {
            'kind': 'block',
            'type': 'test_txt',
          },
          {
            'kind': 'block',
            'type': 'test_embed',
          }
        ],
      }
    });
}


// name: name for the mutator, referenced by the block json.
// mixinObj: Contains mutator methods, and any other data to copy to a new block.
// opt_helperFn: Called after a block with this mutation has first been created.
// opt_blockList: A list of blocks to appear in the flyout of the mutator dialog.
Blockly.Extensions.registerMutator(
  'tapestry_generic_mutation', {
    // cant put them as members or they'll share
    // extraState: {}, // serialized data, required by blockly
    // itemState:  {}, // a cache to avoid counting actual inputs
    // these are required functions, even if we decide not to serialize into blockly format.
    // ( ex. the "insertion manager" uses this as part of dragging blocks. )
    saveExtraState() {
      return Object.assign({}, this.extraState);
    },
    loadExtraState(extraState) {
      this.extraState= Object.assign({}, extraState);
      this.updateShape_();
    },
    // create the MUI from the block's desired state
    decompose: function(workspace) {
      const self= this; // the block we are creating the mui from.
      const mui = workspace.newBlock(`_${self.type}_mutator`); // ex. "_text_value_mutator"
      mui.initSvg();
      mui.inputList.forEach(function(min/*, index, array*/) {
        min.fieldRow.forEach(function(field/*, index, array*/) {
          const wants= self.getExtraState(min.name);
          if (field instanceof Blockly.FieldCheckbox) {
            field.setValue(!!wants);
          } else if (field instanceof Blockly.FieldNumber) {
            field.setValue(wants);
          }
        });
      });
      return mui;
    },
    // modifies the real block based on changes from the MUI.
    // 1: modify the "extra state" based on the mui status.
    // 2: update the block from the extra state.
    compose: function(mui) {
      const self= this;   // our real, workspace, block.
      const jsonDef= jsonDefs[self.type];
      const blockDef= jsonDef.customData.blockDef;
      //
      blockDef.forEach(function(fieldDefs, index/*, array*/) {
        const inputDef= fieldDefs[fieldDefs.length-1];
        // get the mui input ( it might not exist, ex. for fields that arent mutable )
        const min= mui.getInput(inputDef.name);
        // get the primary edit field
        const field= min && min.fieldRow.find(function(field/*, index, array*/) {
          return !!field.name;
        });
        // record the edited status
        if (field instanceof Blockly.FieldCheckbox) {
          const isChecked= field.getValueBoolean();
          self.setExtraState(min.name, isChecked?1:0);
        } else if (field instanceof Blockly.FieldNumber) {
          const itemCount= field.getValue();
          self.setExtraState(min.name, itemCount);
        }
      });
      this.updateShape_();
      this.setShadow(false);
    },
    // update the workspace block based on its current desired state
    updateShape_: function() {
      const self= this;   // our real, workspace, block.
      const jsonDef= jsonDefs[self.type];
      const blockDef= jsonDef.customData.blockDef;
      let insertAt= 1;    // index in the ws block, skipping the initial dummy header.
      //
      blockDef.forEach(function(fieldDefs, index/*, array*/) {
        const inputDef= fieldDefs[fieldDefs.length-1];
        const inputName= inputDef.name;
        // for every desired input, search the existing block to keep the insertion point updated
        const existsAt= self.findInputIndex(inputName, insertAt);
        if (existsAt>=0) {
          insertAt= existsAt;
        }
        const want= self.getExtraState(inputName);
        let have= self.getItemState(inputName);
        if (want == have) {
          insertAt+= have;
        } else {
          self.setItemState(inputName, want);
          // note: the 'repeat' status is disabled for "stackable slots" by the tapestry block generator
          // and other repeating fields are represented by numbers even when they are optional.
          // ie. zero means a non-existent optional repeating field.
          if (inputDef.repeats) {
            while (have > want) {
              --have; // if we want zero, the last removed is name0.
              self.removeInput(inputName + have);
            }
            insertAt+= have;
            while (have < want) {
              // if we have zero, the first added is name0.
              self.createInput(inputName+have, fieldDefs, ++insertAt);
              ++have;
            }
          } else if (!want) {
            self.removeInput(inputName, /*opt_quiet*/ true);
          } else {
            self.createInput(inputName, fieldDefs, ++insertAt);
          }
        }
      });
    },
  }, undefined,
  [/* we dont have any blocks */]
);

// uses the same data as the generic mutation to initialize the non-mutable fields
// [ this gives us strict ordering for all the fields ]
Blockly.Extensions.register(
  'tapestry_generic_extension',
  // note: the generic mutation and mixin have already been added.
  // ( its not clear from the documentation if all three of these objects could live as one unit somewhere )
  function() {
    var self = this; // this refers to the block that the extension is being run on
    self.extraState= {};
    self.itemState= {};
    const jsonDef= jsonDefs[self.type];
    const blockDef= jsonDef.customData.blockDef;
    // an array of field-input sets
    blockDef.forEach(function(fieldDefs/*, index, array*/) {
      const inputDef= fieldDefs[fieldDefs.length-1];
      if (!inputDef.optional) {
        let name= inputDef.name;
        if (inputDef.repeats) {        // a non-optional repeating field needs at least one element
          self.setExtraState(name, 1); // track what we want.
          self.setItemState(name, 1);  // track what we're about to have.
          name += "0";                 // the first field is ...0 to match updateShape_.
        }
        // create the initial input
        self.createInput(name, fieldDefs);
      }
    });
  });

// mix for helper functions
Blockly.Extensions.registerMixin(
  'tapestry_generic_mixin', {
    // extra state tracks our *desired* block appearance
    getExtraState(name)  {
      return this.extraState[ name ] || 0;
    },
    setExtraState(name, itemCount)  {
      return this.extraState[ name ]= itemCount;
    },
    // item state tracks our *desired* block appearance
    getItemState(name)  {
      return this.itemState[ name ] || 0;
    },
    setItemState(name, itemCount)  {
      return this.itemState[ name ]= itemCount;
    },

    // given the named input, return its index.
    // ( blockly's getInput() returns the actual input itself )
    findInputIndex(name, from=0) {
      let found= -1;
      for (let i= from; i< this.inputList.length; i++) {
        const input= this.inputList[i];
        if (input.name === name) {
          found= i;
          break
        }
      }
      return found;
    },

    // called by swap dropdowns to disconnect incompatible blocks after the user swaps the value.
    onSwapChanged(inputName, fieldDef, newValue) {
      const allowedType= fieldDef.swaps[newValue];
      const input= this.getInput(inputName);
      const targetConnection= input.connection && input.connection.targetConnection;
      if (targetConnection) {
        const checks= targetConnection.getCheck(); // null if everything is allowed; rare.
        const stillCompatible= !checks || checks.includes(allowedType);
        if (!stillCompatible) {
          targetConnection.disconnect(); // bumps the disconnecting block away automatically.
        }
      }
    },

    // 1. search for swaps, and change the swap to match the input type.
    onchange: function(e) {
      // move events are used for connections
      if (e.type === Blockly.Events.BLOCK_MOVE) {
        // we only care about events intended for this specific block
        if (e.newParentId === this.id) {
          // find the input that's being connected to
          const input= this.getInput(e.newInputName);
          if (input && input.inputDef) {
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
    _updateSwap: function(input, field) {
      // if the input is a swap:
      const swaps= field.fieldDef && field.fieldDef.swaps;
      if (swaps) {
        const targetConnection= input.connection && input.connection.targetConnection;
        // assuming something is connected
        if (targetConnection) {
          // first check the current value ....
          const currOpt= field.getValue();
          const currType= swaps[currOpt];
          const checks= targetConnection.getCheck(); // null if everything is allowed; rare.
          const stillCompatible= !checks || checks.includes(currType);
          if (!stillCompatible) {
            for (let k in swaps) {
              const checkType= swaps[k];
              if (checks.includes(checkType)) {
                field.setValue(k); // this triggers a "swapChanged" :/
                break;
              }
            }
          }
        }
      }
    },

    // allows overrides of the input name to handle repeated elements
    createInput: function(inputName, fieldDefs, atIndex) {
      const inputDef= fieldDefs[fieldDefs.length-1];
      const appendFn= {
        // note: the names "statement_input" etc.
        // only have meaning for the json descriptions --
        // and the public interface doesn't have a generic append
        'input_dummy': 'appendDummyInput',
        'input_value': 'appendValueInput',
        'input_statement': 'appendStatementInput',
      }[inputDef.type];
      if (!appendFn) {
        throw new Error(`Tapestry mutation couldn't create ${inputName} of ${inputDef.type}`);
      }
      const newInput= this[appendFn](inputName);
      newInput.inputDef= inputDef;
      const newIndex= this.inputList.length-1;
      if (atIndex && atIndex < newIndex) {
        this.moveNumberedInputBefore(newIndex, atIndex);
      }
      if (inputDef.check) {
        newInput.setCheck(inputDef.check);
      }
      let fieldCount=0;
      for (let i=0; i<fieldDefs.length-1; i++) {
        const fieldDef= fieldDefs[i];
        const field= Blockly.fieldRegistry.fromJson(fieldDef);
        field.fieldDef= fieldDef;
        // note: the field doesnt get its name from the definition; instead you have to set it during append.
        // also: dummy inputs dont count when being saved, so we have to scope it
        let fieldName;
        if (!(field instanceof Blockly.FieldLabel)) {
          fieldName= inputName;
          if (fieldCount++ > 0) {
            fieldName += "-" + fieldCount;
          }
        }
        if ((fieldDef.swaps) && (field instanceof Blockly.FieldDropdown)) {
          // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_objects/Function/bind
          field.setValidator(this.onSwapChanged.bind(this, inputName, fieldDef));
        }
        newInput.appendField(field, fieldName);
      }
      // means that we created something other than a dummy input
      // and we should give it an initial value
      // const isToolbox= this.workspace!==Blockly.mainWorkspace;
      if (inputDef.shadow) {
        const sub = this.workspace.newBlock(inputDef.shadow);
        // guess at a bunch of random things to make this show up correctly.
        // render is needed otherwise drag crashes trying to access a null location,
        // and initSvg is needed before render.
        // shadow is needed for the toolbox otherwise we get a random extra block when drag starts.
        sub.initSvg(); // needed before render is called
        sub.render(false); // false means: only re/render this block.
        sub.setShadow(true); // shadow cleans up better when done; but you can connect other values
        // sub.setDeletable(false);
        // sub.setMovable(false);
        newInput.connection.connect(sub.outputConnection);
      }
      return newInput;
    }
  });

// function registerFirstContextMenuOptions() {
//   // This context menu item shows how to use a precondition function to set the visibility of the item.
//   const workspaceItem = {
//     displayText: 'Hello World',
//     // Precondition: Enable for the first 30 seconds of every minute; disable for the next 30 seconds.
//     preconditionFn: function(scope) {
//       const now = new Date(Date.now());
//       if (now.getSeconds() < 30) {
//         return 'enabled';
//       }
//       return 'disabled';
//     },
//     callback: function(scope) {
//     },
//     scopeType: Blockly.ContextMenuRegistry.ScopeType.WORKSPACE,
//     id: 'hello_world',
//     weight: 100,
//   };
//   // Register.
//   Blockly.ContextMenuRegistry.registry.register(workspaceItem);

//   // Duplicate the workspace item (using the spread operator).
//   let blockItem = {...workspaceItem}
//   // Use block scope and update the id to a nonconflicting value.
//   blockItem.scopeType = Blockly.ContextMenuRegistry.ScopeType.BLOCK;
//   blockItem.id = 'hello_world_block';
//   Blockly.ContextMenuRegistry.registry.register(blockItem);
// }

// function registerHelpOption() {
//   const helpItem = {
//     displayText: 'Help! There are no blocks',
//     // Use the workspace scope in the precondition function to check for blocks on the workspace.
//     preconditionFn: function(scope) {
//       if (!scope.workspace.getTopBlocks().length) {
//         return 'enabled';
//       }
//       return 'hidden';
//     },
//     // Use the workspace scope in the callback function to add a block to the workspace.
//     callback: function(scope) {
//       Blockly.serialization.blocks.append({
//         'type': 'text',
//         'fields': {
//           'TEXT': 'Now there is a block'
//         }
//       });
//     },
//     scopeType: Blockly.ContextMenuRegistry.ScopeType.WORKSPACE,
//     id: 'help_no_blocks',
//     weight: 100,
//   };
//   Blockly.ContextMenuRegistry.registry.register(helpItem);
// }

// function registerOutputOption() {
//   const outputOption = {
//     displayText: 'I have an output connection',
//     // Use the block scope in the precondition function to hide the option on blocks with no
//     // output connection.
//     preconditionFn: function(scope) {
//       if (scope.block.outputConnection) {
//         return 'enabled';
//       }
//       return 'hidden';
//     },
//     callback: function (scope) {
//     },
//     scopeType: Blockly.ContextMenuRegistry.ScopeType.BLOCK,
//     id: 'block_has_output',
//     // Use a larger weight to push the option lower on the context menu.
//     weight: 200,
//   };
//   Blockly.ContextMenuRegistry.registry.register(outputOption);
// }

// function registerDisplayOption() {
//   const displayOption = {
//     // Use the block scope to set display text dynamically based on the type of the block.
//     displayText: function(scope) {
//       if (scope.block.type.startsWith('text')) {
//         return 'Text block';
//       } else if (scope.block.type.startsWith('controls')) {
//         return 'Controls block';
//       } else {
//         return 'Some other block';
//       }
//     },
//     preconditionFn: function (scope) {
//       return 'enabled';
//     },
//     callback: function (scope) {
//     },
//     scopeType: Blockly.ContextMenuRegistry.ScopeType.BLOCK,
//     id: 'display_text_example',
//     weight: 100,
//   };
//   Blockly.ContextMenuRegistry.registry.register(displayOption);
// }
