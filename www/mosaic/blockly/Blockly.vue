<template>
  <div id="blockly-area" class="mk-blockly"></div>
  <div id="blockly-div" style="position: absolute"></div>
</template>
<script>
import Blockly from "blockly";
import WorkspaceOptions from "./workspaceOptions.js";
import { CatalogFile } from '/mosaic/catalog/catalogItems.js'

export default {
  emits: ["workspaceLoaded", "changed"],
  props: {
    file: CatalogFile,
    shapeData: Object,
    toolboxData: Object,
  },
  setup() {
    return {
      workspace: null,
      blocklyArea: null, 
      blocklyDiv: null, 
    }
  },
  // use blockly serialization to store the file data to *memory*
  // ( actual save to disk via the server saves the whole catalog at once )
  beforeUnmount() {
    const { workspace, file } = this;
    if (workspace && file) {
      const pod = Blockly.serialization.workspaces.save(workspace);
      const text = JSON.stringify(pod);
      if (file.updateContents(text)) {
        this.$emit("changed", file.path, text);        
      }
    }
    window.removeEventListener("resize", this.onResize, false);
  },
  mounted() {
    this.blocklyDiv = document.getElementById("blockly-div");
    if (!this.blocklyDiv) {
      throw new Error("couldnt find blockly div");
    } else {
      let ws = Blockly.inject("blockly-div", this.workspaceOptions());
      this.workspace = ws;
      ws.toolbox_.flyout_.autoClose = false; // set after the inject

      // hack to override the order of icons so 
      // clicking on a comment icon doesnt cause the mutator to pop in
      const oldBlock = ws.newBlock;
      ws.newBlock = function(p, id) {
        const b = oldBlock.call(ws, p, id);
        const oldIcons = b.getIcons;
        b.getIcons = function() {
          return oldIcons.call(b).reverse();
        };
        return b;
      };
      this.blocklyArea = document.getElementById("blockly-area");
      if (!this.blocklyArea) {
        throw new Error("couldnt find blockly area");
      } else {
        window.addEventListener("resize", this.onResize);
        this.onResize();
        Blockly.svgResize(ws);

        const { file } = this;
        if (file.contents) {
          const pod = JSON.parse(file.contents);
          Blockly.serialization.workspaces.load(pod, ws);
          this.$emit("workspaceLoaded", {
            ws, file,
          });
        }
      }
    }
  },
  methods: {
    workspaceOptions() {
      const w = WorkspaceOptions;
      w.toolbox = this.toolboxData; // overwrite placeholder toolbox with the one from the server.
      return w; 
    },
    // window event to let blockly know to update its canvas.
    // https://developers.google.com/blockly/guides/configure/web/resizable
    onResize() {
      // Compute the absolute coordinates and dimensions of blocklyArea.
      const { blocklyArea, blocklyDiv, workspace } = this;
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

</script>
