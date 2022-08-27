// custom mutation callbacks for modifying blockly shapes.
// this leverages functions provided by tapestryGenericExtension
// ( that extension needs to be applied to all blocks )
import Blockly from "blockly";

// nothing for now...
class CustomMutator extends Blockly.Mutator {}

// note: save/loadExtraState are *required* by blockly and they can *only* live here.
// ( ex. the "insertion manager" uses those as part of dragging blocks. )
const mutatorMixin = {
  // return a copy of the current "extraState" containing the mutation ( called by blockly )
  saveExtraState() {
    return Object.assign({}, this.extraState);
  },
  // called by blockly
  loadExtraState(extraState) {
    this.extraState = Object.assign({}, extraState);
    this.updateShape_(); // lives on tapestry_generic_mixin
  },
  // create the MUI from the block's desired state ( called by blockly )
  decompose: function (workspace) {
    let baseType = this.type;
    if (baseType.startsWith("_") && baseType.endsWith("_stack")) {
      baseType = baseType.substring(1, baseType.length - "_stack".length);
    }
    const mui = workspace.newBlock(`_${baseType}_mutator`); // ex. "_text_value_mutator"
    mui.initSvg();
    mui.inputList.forEach((min /*, index, array*/) => {
      min.fieldRow.forEach((field /*, index, array*/) => {
        // "extraState" contains our desired block appearance
        const wants = this.getExtraState(min.name);
        if (field instanceof Blockly.FieldCheckbox) {
          field.setValue(!!wants);
        } else if (field instanceof Blockly.FieldNumber) {
          field.setValue(wants);
        }
      });
    });
    return mui;
  },
  // called by blockly to modify the desiredState based on the MUI, then updates the workspace shape from that.
  compose: function (mui) {
    const blockType = this.type;
    const customData = this.shapeData;
    const shapeDef = customData.shapeDef;
    // first: modify our "extraState" based on the mui editor.
    shapeDef.forEach((itemDef) => {
      // get the mui input ( it might not exist, ex. for fields that arent mutable or even notable )
      const min = itemDef.name && mui.getInput(itemDef.name);
      // get the primary edit field
      const muiField =
        min &&
        min.fieldRow.find(function (muiField /*, index, array*/) {
          return !!muiField.name;
        });
      // record the status of the mui field
      if (muiField instanceof Blockly.FieldCheckbox) {
        const isChecked = muiField.getValueBoolean();
        this.setExtraState(min.name, isChecked ? 1 : 0);
      } else if (muiField instanceof Blockly.FieldNumber) {
        const itemCount = muiField.getValue();
        this.setExtraState(min.name, itemCount);
      }
    });
    // now: build the block ( and itemState ) from the extraState we just determined.
    this.updateShape_(); // lives on tapestry_generic_mixin
    this.setShadow(false);
  },
  // custom mixin function
  showMutatorIcon(show) {
    const showing = !!this.mutator;
    if (showing != show) {
      if (!show) {
        this.setMutator(null);
      } else {
        const blockNamesForFlyout = []; // we don't use custom blocks in mutations.
        const mutator = new CustomMutator(blockNamesForFlyout);
        this.setMutator(mutator); // creates the iconGroup_  ( an SVGElement ) ( and shows it )
        // mutator.iconGroup_.classList.add("mk-blockly-mutanticon");
      }
    }
  },
};

// normally we'd export the mutatorMixin and the caller would use 'registerMutator'
// we need to override the Mutator class as well ( for icon handling ) so callers should use 'register' instead.
// this function gets applied to each new block; the block is "this".
export default function () {
  // note: skipping blockly's sanity checks.
  // *and* its automatic creation of the mutator icon.
  this.mixin(mutatorMixin);
}
