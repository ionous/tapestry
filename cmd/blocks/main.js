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
            'type': 'text_value',
          },
          {
            'kind': 'block',
            'type': 'text_values',
          },
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
    itemCounts_: {},
    // create the MUI.
    decompose: function(workspace) {
      const self= this; // the block we are creating the mui from.
      const mui = workspace.newBlock(`_${self.type}_mutator`); // ex. "_text_value_mutator"
      mui.initSvg();
      mui.inputList.forEach(function(min/*, index, array*/) {
        min.fieldRow.forEach(function(field/*, index, array*/) {
          if (field instanceof Blockly.FieldCheckbox) {
            const exists= !!self.getInput( min.name );
            field.setValue(exists);
          } else if (field instanceof Blockly.FieldNumber) {
            let itemCount= 0; // fix? could make this faster via "findInputIndex" tapestry_mutation_mixin
            // note: the input names in the mui dont have trailing numbers;
            // the names in the actual block do. ( VALUES vs. VALUES0 )
            while (self.getInput( min.name + itemCount )) {
              itemCount++;
            }
            field.setValue(itemCount);
            self.itemCounts_[min.name]= itemCount;
          }
        });
      });
      return mui;
    },
    // modifies the real block based on changes from the MUI.
    compose: function(mui) {
      // mui.setEnabled(true);
      const self= this;   // our real, workspace, block.
      let insertAt= 1;    // index in the ws block, skipping the initial dummy header.
      const jsonDef= jsonDefs[self.type];
      const muiData= jsonDef.customData.muiData;
      muiData.forEach(function(fieldDefs, index/*, array*/) {
        const inputDef= fieldDefs[fieldDefs.length-1];
        const inputName= inputDef.name;
        // for every desired input, search the existing block to keep the insertion point updated
        const existsAt= self.findInputIndex(inputName, insertAt);
        if (existsAt>=0) {
          insertAt= existsAt;
        }
        // the mui input might not exist (ex. for fields  that arent mutable)
        const min=  mui.getInput(inputName);
        const field= min && min.fieldRow.find(function(field/*, index, array*/) {
          return !!field.name; // find the edit field
        });
        if (field instanceof Blockly.FieldCheckbox) {
          const isChecked= field.getValueBoolean();
          // do we need to remove it?
          if (!isChecked && existsAt>=0) {
            self.removeInput(inputName, /*opt_quiet*/ true);
          } else if (isChecked && existsAt< 0) {
            // do we need to to create it?
            self.createInput(inputName, fieldDefs, insertAt+1);
          }
        } else if (field instanceof Blockly.FieldNumber) {
          // when we created the mui, we cached the block's item count.
          let actual= self.itemCounts_[inputName];
          // the editable number is how many items we want.
          const want= field.getValue();
          while (actual > want) {
            --actual; // if we want zero, the last removed is name0.
            self.removeInput(inputName + actual);
          }
          insertAt+= actual;
          while (actual < want) {
            // if we have zero, the first added is name0.
            self.createInput(inputName+actual, fieldDefs, insertAt+1);
            ++actual;
            ++insertAt;
          }
          self.itemCounts_[inputName]= actual;
        }
      });
    },

    // these are required functions, even if we decide not to serialize into blockly format.
    saveExtraState() {
      return {
        'itemCount': "secrets",
      };
    },
    loadExtraState(state) {},

  }, undefined,
  [/* we dont have any blocks */]
);

// uses the same data as the generic mutation to initialize the non-mutable fields
// [ this gives us strict ordering for all the fields
Blockly.Extensions.register(
  'tapestry_mutation_extension',
  function() { // this refers to the block that the extension is being run on
    var self = this;
    const jsonDef= jsonDefs[self.type];
    const muiData= jsonDef.customData.muiData;
    // an array of field-input sets
    muiData.forEach(function(fieldDefs/*, index, array*/) {
      const inputDef= fieldDefs[fieldDefs.length-1];
      if (!inputDef.optional) {
        let name= inputDef.name;
        if (inputDef.repeats) {
          name += "0"; // ugly, but simplifies counting in tapestry_generic_mutation
        }
        // create the initial input
        self.createInput(name, fieldDefs);
      }
    });
  });

// mix for helper functions
Blockly.Extensions.registerMixin(
  'tapestry_mutation_mixin', {
    // find the named input, return its index
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
    // allows overrides of the input name to handle repeated elements
    createInput: function(name, fieldDefs, atIndex) {
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
        throw new Error(`Tapestry mutation couldn't create ${name} of ${inputDef.type}`);
      }
      const newInput= this[appendFn](name);
      const newIndex= this.inputList.length-1;
      if (atIndex && atIndex < newIndex) {
        this.moveNumberedInputBefore(newIndex, atIndex);
      }
      if (inputDef.check) {
        newInput.setCheck(inputDef.check);
      }
      for (let i=0; i<fieldDefs.length-1; i++) {
        const fieldDef= fieldDefs[i];
        const field= Blockly.fieldRegistry.fromJson(fieldDef);
        newInput.appendField(field);
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
