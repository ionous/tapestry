<template
	><div id="blockly-area" class="mk-blockly"></div
	><div id="blockly-div" style="position: absolute"></div
></template>

<script>

import Blockly from 'blockly';
Blockly.WorkspaceAudio.prototype.preload= function(){};
		
const toolbox = {
  "kind": "flyoutToolbox",
  "contents": [
    {
      "kind": "block",
      "type": "controls_if"
    },
    {
      "kind": "block",
      "type": "controls_repeat_ext"
    },
    {
      "kind": "block",
      "type": "logic_compare"
    },
    {
      "kind": "block",
      "type": "math_number"
    },
    {
      "kind": "block",
      "type": "math_arithmetic"
    },
    {
      "kind": "block",
      "type": "text"
    },
    {
      "kind": "block",
      "type": "text_print"
    }
  ]
}

// https://developers.google.com/blockly/guides/configure/web/resizable
let workspace, blocklyArea, blocklyDiv;

export default {
  mounted() {
		workspace= Blockly.inject('blockly-div', {toolbox: toolbox});
		blocklyArea = document.getElementById('blockly-area');
  	blocklyDiv = document.getElementById('blockly-div');
  	window.addEventListener('resize',this.onResize);
  	this.onResize();
  	Blockly.svgResize(workspace);
  },
	destroyed: function() {
	  window.removeEventListener('resize', this.onResize, false);
	},
	methods: {
		onResize() {
			// Compute the absolute coordinates and dimensions of blocklyArea.
	    var element = blocklyArea;
	    var x = 0;
	    var y = 0;
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

</script>

