 <template
 ><img alt="Vue logo" src="./assets/logo.png" 
   /><router-view/><Catalog :catalog="catalog" 
   /><Blockly></Blockly
></template>
 
 <script>
 import MockCatalog from './catalog/mockCatalog.js'
 import Catalog from './catalog/Catalog.vue'
import Blockly from './blockly/Blockly.vue'
import { RouterView } from 'vue-router'

export default {
  components: { Catalog, Blockly, RouterView},
  data() {
    return {
      catalog: new MockCatalog()
    }
  }
}
</script>



<!--   
the blockly music maker is interesting.
it has classes on "permanent" divs -- example the edit div --
which default to not displaying, and then adds a global "mode" attribute to the body
to select which ones display. nice and simple. similar the ".blockly-editor" class is 
left: -400px; by default, and during mode-blockly it "moves" into place.

they load and save the workspace at every transition.
using: 
  let xml = Blockly.Xml.workspaceToDom(Blockly.getMainWorkspace());
  Blockly.Xml.domToWorkspace(button.blocklyXml, workspace);

tbd: that might now work because im using "extraState"
but could just serialize and save to memory somewhere, 
[ ex. if file persisted ]

- 

.mode-edit, .mode-maker, .mode-blockly {
  display: none;
}

[mode="maker"] .mode-maker,
[mode="edit"] .mode-edit,
[mode="blockly"] .mode-blockly {
  display: block;
}
[mode="blockly"] .blockly-editor {
  left: 0;
}


/Users/ionous/Dev/blockly/google-blockly/demos/codelab/app-complete/index.html
https://blocklycodelabs.dev/codelabs/getting-started/index.html

function enableEditMode() {
    document.body.setAttribute('mode', 'edit');
    document.querySelectorAll('.button').forEach(btn => {
      btn.removeEventListener('click', handlePlay);
      btn.addEventListener('click', enableBlocklyMode);
    });
  }
   provide() {
    
    }
  },
 -->