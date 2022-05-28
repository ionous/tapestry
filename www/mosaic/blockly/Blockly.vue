<template
  ><div id="blockly-area" class="mk-blockly"></div
  ><div id="blockly-div" style="position: absolute"></div
></template>
<script>
import Blockly from 'blockly';
import Cataloger from '/mosaic/catalog/cataloger.js'
import WorkspaceOptions from './workspaceOptions.js'

Blockly.WorkspaceAudio.prototype.preload= function(){};

// https://developers.google.com/blockly/guides/configure/web/resizable
let workspace, blocklyArea, blocklyDiv;
let currentFile;     // the last successfully displayed file
let pendingLoad;     // we only allow one load request at a time
let shapeData= null; // the customData from the shapes file.

export default {
  emits: ["workspaceChanged"],
  props: {
    catalog: Cataloger,
    shapeData: Object, 
    toolboxData: Object,
  },
  mounted() {
    blocklyDiv = document.getElementById('blockly-div');
    if (!blocklyDiv) {
      throw new Error("couldnt find blockly div");
    } else {
      const w= WorkspaceOptions;
      w.toolbox= this.toolboxData; // overwrite placeholder toolbox with the one from the server.
      workspace= Blockly.inject('blockly-div', w);
      blocklyArea = document.getElementById('blockly-area');
      if (!blocklyArea) {
        throw new Error("couldnt find blockly area");
      } else {
        window.addEventListener('resize',this.onResize);
        this.onResize();
        Blockly.svgResize(workspace);
        this.onRouteChanged(this.$route.params);
      }
    }
  },
  created() {
    shapeData= this.shapeData;
    this.$watch(
      () => this.$route.params,
      (toParams, previousParams) => {
        this.onRouteChanged(toParams);
      }
    )
  },
  destroyed: function() {
    window.removeEventListener('resize', this.onResize, false);
    currentFile= null;
    pendingLoad= null;
  },
  methods: {
    // the browser url has changed and a new story file has been requested: 
    // remember the current workspace and retrieve the requested file
    // ( from memory or from the server depending on if we've seen it before )
    onRouteChanged(params) {
      const { catalog } = this;
      if (params && params.editPath !== undefined) {
        const path = params.editPath.join("/");
        if (path && (!currentFile || (currentFile.path !== path))) {
          if (pendingLoad) {     // if something is already loading
            console.log("queuing:", pendingLoad);
            pendingLoad = path;  // remember the most recent request.
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
        currentFile= file;
        pendingLoad= null;
        if (workspace && file.contents) {
          const pod= JSON.parse(file.contents);
          Blockly.serialization.workspaces.load(pod, workspace);
          this.$emit("workspaceChanged",{
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
      if (workspace && file && (file===currentFile)) { 
        const pod= Blockly.serialization.workspaces.save(workspace);
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
        element = element.offsetParent;
      } while (element);
      // Position blocklyDiv over blocklyArea.
      blocklyDiv.style.left = x + 'px';
      blocklyDiv.style.top = y + 'px';
      blocklyDiv.style.width = blocklyArea.offsetWidth - 3 + 'px';
      blocklyDiv.style.height = blocklyArea.offsetHeight - 3 + 'px';
      Blockly.svgResize(workspace);
    }
  }
}

// handle vite dev reloads. 
function safeRegister(name, fn, obj) {
  if (Blockly.Extensions.isRegistered(name)) {
    Blockly.Extensions.unregister(name);
  }
  fn(name, obj);
}

// custom mutation callbacks for modifying blockly shapes.
// note: save/loadExtraState are *required* by blockly and they can *only* live in the mutator.
// ( ex. the "insertion manager" uses those as part of dragging blocks. )
// this implementation leverages tapestry_generic_mixin ( also needed for non-mutating shapes. )
safeRegister('tapestry_generic_mutation', Blockly.Extensions.registerMutator, {
  // return a copy of the current "extraState" containing the mutation ( called by blockly )
  saveExtraState() {
    return Object.assign({}, this.extraState);
  },
  // called by blockly
  loadExtraState(extraState) {
    this.extraState= Object.assign({}, extraState);
    this.updateShape_(); // lives on tapestry_generic_mixin
  },
  // create the MUI from the block's desired state ( called by blockly )  
  decompose: function(workspace) {
    const self= this; // the block we are creating the mui from.
    const customData= shapeData[self.type];
    const mui = workspace.newBlock(customData.mui); // ex. "_text_value_mutator"
    mui.initSvg();
    mui.inputList.forEach(function(min/*, index, array*/) {
      min.fieldRow.forEach(function(field/*, index, array*/) {
        // "extraState" contains our desired block appearance
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
  // called by blockly: modifies the desired block extraState based on the MUI, then updates the shape from that state.
  compose: function(mui) {
    const self= this;   // our real, workspace, block.
    const customData= shapeData[self.type];
    const shapeDef= customData.shapeDef;
    // first: modify our desired "extraState" based on the mui editor's current status.
    shapeDef.forEach(function(fieldDefs, index/*, array*/) {
      const inputDef= fieldDefs[fieldDefs.length-1];
      fieldDefs.forEach(function(fieldDef) {
        // find notable fields or inputs in our real block's shape definition
        // ( dummy inputs wind up having the same name as their editing field, 
        //   we dont need to check the same value twice so filter it )
        if (fieldDef === inputDef || (fieldDef.name && fieldDef.name !== inputDef.name)) {
          // get the mui input ( it might not exist, ex. for fields that arent mutable )
          const min = mui.getInput(fieldDef.name);
          // get the primary edit field
          const muiField = min && min.fieldRow.find(function(muiField/*, index, array*/) {
            return !!muiField.name;
          });
          // record the status of the mui field
          if (muiField instanceof Blockly.FieldCheckbox) {
            const isChecked= muiField.getValueBoolean();
            self.setExtraState(min.name, isChecked?1:0);
          } else if (muiField instanceof Blockly.FieldNumber) {
            const itemCount= fmuiField.getValue();
            self.setExtraState(min.name, itemCount);
          }
        }
      });
    });
    // now: build the block from the desired values ( the "extraState" ) we just generated.
    this.updateShape_();  // lives on tapestry_generic_mixin
    this.setShadow(false);
  }
});

// create a shape's non-mutable fields and inputs;
// loadExtraState takes care of the rest.
// relies on the generic mutation and mixin having already been added.
safeRegister('tapestry_generic_extension', Blockly.Extensions.register, function() {
  var self = this;       // refers to the block that the extension is being run on
  self.extraState= {};   // tracks the desired block appearance ( see tapestry_generic_mixin )
  self.itemState= {};    // tracks the actual block appearance ( see tapestry_generic_mixin )
  // self.setCommentText("");
  const customData= shapeData[self.type];
  const shapeDef= customData.shapeDef;
  // an array of field-input sets
  shapeDef.forEach(function(fieldDefs/*, index, array*/) {
    // the last entry is always an input ( ex. possibly a dummy or a connector )
    const inputDef = fieldDefs[fieldDefs.length-1];
    // and we're only worrying about the required fields and inputs here.
    if (!inputDef.optional) {        
      // notable fields (and inputs) are given names and we tell updateShape that they exist.
      fieldDefs.forEach(function(fieldDef) {
        if (fieldDef.name)  {
          self.setItemState(fieldDef.name, 1);    // track what we're about to have.
          self.setExtraState(fieldDef.name, 1);   // track what we want: noting, it gets overwritten by loading saves.
        }
      });
      // create the initial input ( which potentially includes a row of fields )
      let name= inputDef.name;
      if (inputDef.repeats) {        // note: stackable inputs are marked as non-repeating in the shapeDef
        name += "0";                 // counters are "name", inputs are "name#" for updateShape_()
      }
      self.createInput(name, fieldDefs);
    }
  });
});

// mix for helper functions
safeRegister('tapestry_generic_mixin', Blockly.Extensions.registerMixin, {
  // extra state tracks our *desired* block appearance
  // fix? while ideally this would only exist for mutable fields;
  // the block generator doesnt always distinguish between optional and required fields.
  getExtraState(name, defaultVal=0)  {
    return this.extraState[ name ] || defaultVal;
  },
  setExtraState(name, itemCount)  {
    return this.extraState[ name ]= itemCount;
  },
  // item state tracks our actual block appearance
  // ( everything starts at 0, not present. ) 
  getItemState(name)  {
    return this.itemState[ name ] || 0;
  },
  setItemState(name, itemCount)  {
    return this.itemState[ name ]= itemCount;
  },
  // extraState: {}, // created by the generic extension "constructor"...
  // itemState: {},  // so they dont wind up shared across instances.

  // given the named input, return its index.
  // ( blockly's getInput() returns the actual input itself )
  // findInputIndex(name, from=0) {
  //   let found= -1;
  //   for (let i= from; i< this.inputList.length; i++) {
  //     const input= this.inputList[i];
  //     if (input.name === name) {
  //       found= i;
  //       break
  //     }
  //   }
  //   return found;
  // },

  // update the workspace block based on its current desired state
  // ( called from load and the mutation's compose )
  updateShape_: function() {
    const self= this;   // our real, workspace, block.
    const customData= shapeData[self.type];
    const shapeDef= customData.shapeDef;
    let insertAfter= 0;    // index in the ws block. inserting after index 0, the initial header.
    // shapeDef is tapestry's custom data format.
    shapeDef.forEach(function(fieldDefs, index/*, array*/) {
      const inputDef= fieldDefs[fieldDefs.length-1];
      const inputName= inputDef.name;
      // this setup handles some cases where extraState has info where its not needed.
      // ( ex. the block generator records the number of stacked inputs even tho that's... wrong. )
      const want= self.getExtraState(inputName, inputDef.optional?0:1);
      let have= self.getItemState(inputName);
      if (want == have) {
        insertAfter += inputDef.repeats? have: 1;
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
          insertAfter+= have;
          while (have < want) {
            // if we have zero, the first added is name0.
            self.createInput(inputName+have, fieldDefs, ++insertAfter);
            ++have;
          }
        } else if (!want) {
          self.removeInput(inputName, /*opt_quiet:dont error if not there.*/ true);
        } else if (!have) {
          self.createInput(inputName, fieldDefs, ++insertAfter);
        } else {
          ++insertAfter; // want and have and not repeating.
        }
      }
    });
  },

  // disconnect incompatible blocks after the author selects a new option from a swap's dropdown.
  // [ this is a custom validator callback bound in createInput ]
  onSwapChanged(inputName, fieldDef, newValue) {
    const allowedType= fieldDef.swaps[newValue];
    const input= this.getInput(inputName);
    const targetConnection= input.connection && input.connection.targetConnection;
    if (targetConnection) {
      // blockly function which returns a list of compatible value types
      // ( null if everything is allowed; rare. )
      const checks= targetConnection.getCheck();
      const stillCompatible= !checks || checks.includes(allowedType);
      if (!stillCompatible) {
        targetConnection.disconnect(); // bumps the disconnecting block away automatically.
      }
    }
  },

  // blockly change notification ( responding to changes made by authors )
  // this searches for swaps, and changes the swap to match the input type.
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
  // ( ie. when an author tries to connect a block to a swappable input
  //   change the swap's combo box to match. )
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
        // blockly function which returns a list of compatible value types
        // ( null if everything is allowed; rare. )
        const checks= targetConnection.getCheck(); 
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
      // note: the field actually doesnt get its name from the definition; instead you have to set it during append.
      // and when generating names, we have to make sure they are all unique: scoped to the input.
      let fieldName= fieldDef.name;
      if (!fieldName && !(field instanceof Blockly.FieldLabel)) {
        fieldName= inputName;
        if (fieldCount++ > 0) {
          fieldName += "-" + fieldCount;
        }
      }
      // add a validator to disconnect incompatible blocks after a combo box change. 
      if ((fieldDef.swaps) && (field instanceof Blockly.FieldDropdown)) {
        // use bind to give the callback some helpful parameters.
        // ( https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_objects/Function/bind )
        field.setValidator(this.onSwapChanged.bind(this, inputName, fieldDef));
      }
      newInput.appendField(field, fieldName);
    }
    // means that we created something other than a dummy input
    // and we should give it an initial value
    const isToolbox= this.workspace!==Blockly.mainWorkspace;
    if (isToolbox && inputDef.shadow) {
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

</script>


