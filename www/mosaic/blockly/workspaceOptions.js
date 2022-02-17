// evelopers.google.com/blockly/guides/configure/web/configuration_struct
// unless otherwise specified boolean options are toolbox dependent; defaults to true if the toolbox has categories, false otherwise.
export default {
  // collapse: // can blocks be collapsed or expanded.
  comments: true, // add commenting to the context menu.
  // css:      // when false, CSS becomes the document's responsibility. Defaults to true.
  // disable:  // Allows blocks to be disabled.
  // grid:     // Configures a grid which blocks may snap to.
  // horizontalLayout: // If true toolbox is horizontal, if false toolbox is vertical. Defaults to false.
  // maxBlocks:        // number of blocks that may be created. Useful for student exercises. Defaults to Infinity.
  // maxInstances:     // { block types =>  maximum number of that type }. Undeclared types default to Infinity.
  // media:            //  Path from page (or frame) to the Blockly media directory. Defaults to "https://blockly-demo.appspot.com/static/media/".
  move:   {            // https://developers.google.com/blockly/guides/configure/web/move
    scrollbars: true,
    drag: false,
    wheel: true,
  },          // { scrollbars, drag, wheel }. Configures how users can move around the workspace.
  // oneBasedIndex:    // If true list and string operations should index from 1, if false index from 0. Defaults to true.
  // readOnly:         // If true, prevent the user from editing. Suppresses the toolbox and trashcan. Defaults to false.
  renderer:  "thrasos", //  Determines the renderer used by blockly.  Pre-packaged renderers include
                      // 'geras' (the default),
                      // 'thrasos', and
                      // 'zelos' (a scratch-like renderer) -- has rounded internal blocks
  // rtl: // If true, mirror the editor (for Arabic or Hebrew locales). See RTL demo. Defaults to false.
  scrollbars: true,   // Takes an object where the horizontal property determines if horizontal scrolling is enabled and
              // the vertical property determines if vertical scrolling is enabled.
              // If a boolean is passed then it is equivalent to passing an object with both horizontal and vertical properties set as that value.
              // Defaults to true if the workspace has categories.
  // theme: // Defaults to classic theme if no theme is provided.
  // toolbox: // categories and blocks available to the user.
  // toolboxPosition: //  If "start" toolbox is on top (if horizontal) or left (if vertical and LTR) or right (if vertical and RTL).
                      // If "end" toolbox is on opposite side. Defaults to "start".
  trashcan: true, // Displays or hides the trashcan.
  // maxTrashcanContents: // Maximum number of deleted items that will appear in the trashcan flyout. '0' disables the feature. Defaults to '32'.
  // plugins: // Map of plugin type to name of registered plugin or plugin class. See injecting subclasses.
  zoom:{             // https://developers.google.com/blockly/guides/configure/web/zoom
    controls:true, // show zoom-centre, zoom-in, and zoom-out buttons.
    // wheel:true, // Set to true to allow the mouse wheel to zoom. Defaults to false.
    startScale:0.8,
    maxScale:3,
    minScale:0.3,
    scaleSpeed:1.2,
    pinch: true, // Set to true to enable pinch to zoom support on touch devices. Defaults to true if either the wheel or controls option is set to true.
  },
  // toolbox: toolboxSimple,
  toolbox: {
    'kind': 'flyoutToolbox',
    'contents': [],
  }
};
