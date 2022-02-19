 <template
  ><mk-startup 
    v-if="!shapeData"
    @started="onStarted"
  /><template v-else
      ><div class="mk-header"
      ><h1>Tapestry</h1></div
      ><mk-toolbar
        @save="onSave"
      /><mk-catalog
        :catalog="catalog"
      /><mk-blockly 
        :catalog="catalog"
        :shape-data="shapeData"
        :toolbox-data="toolboxData"
        @workspace-changed="onWorkspaceChanged"
  /></template
></template>
 
<script>
import BlockCatalog from './catalog/blockCatalog.js'
import mkBlockly from './blockly/Blockly.vue'
import mkCatalog from './catalog/Catalog.vue'
import mkToolbar from './Toolbar.vue'
import mkStartup from './Startup.vue'
// import { RouterView } from 'vue-router'
import { ref, onMounted } from 'vue'

// appcfg comes through vite conifg.
const catalog= new BlockCatalog(appcfg.mosaic + '/blocks/'); // only need one.

export default {
  components: { mkBlockly, mkCatalog, mkToolbar, mkStartup },
  setup(props) {
    let shapeData= ref(null);
    let toolboxData= ref(null);
    return {
      catalog, // unwatched.
      shapeData, 
      toolboxData,
      workspace: null,
      onStarted(b) {
        shapeData.value= Object.freeze(b.shapeData);
        toolboxData.value= Object.freeze(b.toolboxData);
      }
    }
  },
  methods: {
    onSave() {
      // we have to look at the currently focused file as well to see if it has any changes 
      // and if so flush them to its in memory catalog file.
      // but *Blockly* has to be the thing to do that... not really something vue wants us to do.
      console.log("save");
      if (this.workspace) {
        this.workspace.flush();
      }
      catalog.saveStories();
    },
    onWorkspaceChanged(ws) {
      this.workspace= ws;
    }
  },
  provide() {
    const { "$router": router, catalog } = this;
    return {
      onFolder(folder) {
        if (folder.contents) {
          console.log("closing", folder.name);
          folder.contents= false;
        } else {
          // injects the list of sub-files into the passed folder
          console.log("opening", folder.name, folder.path);
          catalog.loadFolder(folder);
        }
      },
      onFile(file) {
        // a little convoluted: we change the router; the router changes blockly
        // blockly tells us when the workspace has changed.
        const parts= file.path.split("/");
        router.push({ name: 'edit', params: { editPath: parts } });
      },
    }
  }
}
</script>

