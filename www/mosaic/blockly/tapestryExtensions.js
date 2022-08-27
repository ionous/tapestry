import ShapeWriter from "./shapeWriter.js";

// a function to generate a blockly extension for creating a block's non-mutable fields and inputs.
// loadExtraState ( from the tapestryMutation ) takes care of the dynamic (mutable) fields and inputs.
// the passed { shapes } contains a custom interface to access global shape data.
// ( relies on the tapestry mutator and helpers also having been mixed in. )
export default function(shapes) {
  return function() {
    this.extraState = {}; // tracks the mutation state: the desired block; appearance ( see tapestry_generic_mixin )
    this.itemState = {}; // tracks the workspace state: the actual block appearance ( see tapestry_generic_mixin )
    this.shapeData = shapes.getShapeData(this.type);

    // this.setCommentText("");
    const customData = this.shapeData;
    const shapeDef = customData.shapeDef;
    const shapeWriter = new ShapeWriter(this, ({ input, field }) => {
      this.onItemAdded(input, field);
    });
    // create our default workspace state,
    // then blockly will load extraState and we will mutate into the desired shape.
    shapeDef.forEach((itemDef /*, index, array*/) => {
      if (itemDef.optional) {
        console.assert(itemDef.name); // optional elements should have a name.
      } else {
        // we track all notable optional/non-optional items to simplify updateShape
        // ( though we're only worrying about the required fields and inputs here. )
        if (itemDef.name) {
          this.setItemState(itemDef.name, 1); // what we're about to have.
          this.setExtraState(itemDef.name, 1); // what we want: overwritten by loading saves.
        }
        shapeWriter.createItems(itemDef, 0, 1);
      }
    });
    shapeWriter.finalize();
  };
};
