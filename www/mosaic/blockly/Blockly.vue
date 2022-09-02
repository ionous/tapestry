<template>
  <div id="blockly-area" class="mk-blockly"></div>
  <div id="blockly-div" style="position: absolute"></div>
</template>
<script>
import Blockly from "blockly";
import Cataloger from "/mosaic/catalog/cataloger.js";
import WorkspaceOptions from "./workspaceOptions.js";

import MosaicStrField from "./mosaicStrField.js";
import MosaicTextField from "./mosaicTextField.js";
import MosaicMultilineField from "./mosaicMultilineField.js";
import tapestryMutator from "./tapestryMutator.js";
import tapestryExtensions from "./tapestryExtensions.js";
import tapestryHelpers from "./tapestryHelpers.js";
import "./customColors.js";

// todo: there are issues with the audio preloading...
Blockly.WorkspaceAudio.prototype.preload = function () {};

// https://developers.google.com/blockly/guides/configure/web/resizable
let workspace, blocklyArea, blocklyDiv;
let currentFile; // the last successfully displayed file
let pendingLoad; // we only allow one load request at a time
let shapeData = null; // the customData from the shapes file.

export default {
  emits: ["workspaceChanged"],
  props: {
    catalog: Cataloger,
    shapeData: Object,
    toolboxData: Object,
  },
  mounted() {
    blocklyDiv = document.getElementById("blockly-div");
    if (!blocklyDiv) {
      throw new Error("couldnt find blockly div");
    } else {
      const w = WorkspaceOptions;
      w.toolbox = this.toolboxData; // overwrite placeholder toolbox with the one from the server.
      workspace = Blockly.inject("blockly-div", w);
      blocklyArea = document.getElementById("blockly-area");
      if (!blocklyArea) {
        throw new Error("couldnt find blockly area");
      } else {
        window.addEventListener("resize", this.onResize);
        this.onResize();
        Blockly.svgResize(workspace);
        this.onRouteChanged(this.$route.params);
      }
    }
  },
  created() {
    shapeData = this.shapeData;
    this.$watch(
      () => this.$route.params,
      (toParams, previousParams) => {
        this.onRouteChanged(toParams);
      }
    );
  },
  destroyed: function () {
    window.removeEventListener("resize", this.onResize, false);
    currentFile = null;
    pendingLoad = null;
  },
  methods: {
    // the browser url has changed and a new story file has been requested:
    // remember the current workspace and retrieve the requested file
    // ( from memory or from the server depending on if we've seen it before )
    onRouteChanged(params) {
      const { catalog } = this;
      if (params && params.editPath !== undefined) {
        const path = params.editPath.join("/");
        if (path && (!currentFile || currentFile.path !== path)) {
          if (pendingLoad) {
            // if something is already loading
            console.log("queuing:", pendingLoad);
            pendingLoad = path; // remember the most recent request.
          } else {
            // currentFile gets updated after the new contents have finished loading.
            if (currentFile) {
              this.storeCurrentWorkspace(currentFile);
            }
            workspace.clear();
            pendingLoad = path;
            catalog.loadFile(path).then(this._onFileLoaded);
          }
        }
      }
    },
    // receives a successfully loaded file ( ex. from a changed route )
    // keeps looping until the result is the file currently desired by the author.
    // ( eg. if the author clicked something different in the meantime )
    _onFileLoaded(file) {
      const { catalog } = this;
      if (file.path !== pendingLoad) {
        console.log("reloading:", pendingLoad);
        catalog.loadFile(pendingLoad).then(this._onFileLoaded);
      } else {
        currentFile = file;
        pendingLoad = null;
        if (workspace && file.contents) {
          const pod = JSON.parse(file.contents);
          Blockly.serialization.workspaces.load(pod, workspace);
          this.$emit("workspaceChanged", {
            file: file,
            flush: () => {
              this.storeCurrentWorkspace(file);
            },
          });
        }
      }
    },
    // use blockly serialization to store the file data to memory
    // ( actual save to disk via the server saves the whole catalog at once )
    storeCurrentWorkspace(file) {
      if (workspace && file && file === currentFile) {
        const pod = Blockly.serialization.workspaces.save(workspace);
        file.updateContents(JSON.stringify(pod));
      }
    },
    // window event to let blockly know to update its canvas.
    onResize() {
      // Compute the absolute coordinates and dimensions of blocklyArea.
      let element = blocklyArea;
      let x = 0;
      let y = 0;
      do {
        x += element.offsetLeft;
        y += element.offsetTop;
        // element = element.offsetParent;
        break; // blockly's example walks all the way up to the top... why?
      } while (element);
      // Position blocklyDiv over blocklyArea.
      blocklyDiv.style.left = x + "px";
      blocklyDiv.style.top = y + "px";
      blocklyDiv.style.width = blocklyArea.offsetWidth - 3 + "px";
      blocklyDiv.style.height = blocklyArea.offsetHeight - 3 + "px";
      Blockly.svgResize(workspace);
    },
  },
};

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

registerFieldType("mosaic_str_field", MosaicStrField);
registerFieldType("mosaic_text_field", MosaicTextField);
registerFieldType("mosaic_multiline_field", MosaicMultilineField);
registerExtension("tapestry_generic_mutation", tapestryMutator);
registerExtension(
  "tapestry_generic_extension",
  tapestryExtensions({
    getShapeData(blockType) {
      return shapeData[blockType];
    },
  })
);
registerExtension("tapestry_generic_mixin", tapestryHelpers);

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
</script>
