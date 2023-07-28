import Blockly from 'blockly';

// one time registration helpers
import MosaicStrField from "./blockly/mosaicStrField.js";
import MosaicTextField from "./blockly/mosaicTextField.js";
import MosaicMultilineField from "./blockly/mosaicMultilineField.js";
import tapestryMutator from "./blockly/tapestryMutator.js";
import tapestryExtensions from "./blockly/tapestryExtensions.js";
import tapestryHelpers from "./blockly/tapestryHelpers.js";
import "./blockly/customColors.js";
import endpoint from './endpoints.js'
//
import http from '/lib/http.js'


// start loading right away?
// might be easier for debugging if not.
let cache;

export default {
  // promises future data.
  // fix: handle failures in a way that allows reload?
  getShapeData() {
    if (!cache) {
      cache = loadData();
    }
    return cache;
  }
};

function loadData() {
  return Promise.all([
    http.get(endpoint.tools),
    http.get(endpoint.shapes).then((jsonData) => {
      let shapeData= {};
      jsonData.forEach(function(el) {
        Blockly.defineBlocksWithJsonArray([el]);
        let { customData } = el;
        if (customData) { // ironically: mutators use standard blocks.
          shapeData[el.type]= customData;
        }
      });
      return shapeData;
    })
  ]).then((values)=> {
    const tools = values[0];
    const shapes = values[1];

    registerFieldType("mosaic_str_field", MosaicStrField);
    registerFieldType("mosaic_text_field", MosaicTextField);
    registerFieldType("mosaic_multiline_field", MosaicMultilineField);
    registerExtension("tapestry_generic_mutation", tapestryMutator);
    registerExtension(
      "tapestry_generic_extension",
      tapestryExtensions({
        getShapeData(blockType) {
          return shapes[blockType];
        },
      })
    );
    registerExtension("tapestry_generic_mixin", tapestryHelpers);
    return {
      tools, shapes
    }
  });
}

// todo: there are issues with the audio preloading...
Blockly.WorkspaceAudio.prototype.preload = function () {};

// handle vite dev reloads.
function registerExtension(name, obj) {
  const registry = Blockly.Extensions;
  if (registry.isRegistered(name)) {
    registry.unregister(name);
  }
  registry.register(name, obj);
}

function registerFieldType(name, obj) {
  // Blockly.fieldRegistry doesn't directly expose isRegistered :/
  // so use this uses layer underneath it.
  const registry = Blockly.registry;
  const FIELD = registry.Type.FIELD;
  if (registry.hasItem(FIELD, name)) {
    registry.unregister(FIELD, name);
  }
  registry.register(FIELD, name, obj);
}

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
