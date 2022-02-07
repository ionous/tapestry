 <template
   ><Startup 
    v-if="!shapeData"
    @started="onStarted"
   /><template v-else
      ><img alt="Vue logo" src="./assets/logo.png" 
      ><router-view
      /><Catalog 
        :catalog="catalog"
      /><Blockly 
        :catalog="catalog"
        :shapeData="shapeData"
  /></template
></template>
 
<script>
import BlockCatalog from './catalog/blockCatalog.js'
import Blockly from './blockly/Blockly.vue'
import Catalog from './catalog/Catalog.vue'
import Startup from './Startup.vue'
import { RouterView } from 'vue-router'
import { ref, onMounted } from 'vue'

// dataUrl comes through vite conifg.
const catalog= new BlockCatalog(dataUrl + '/blocks/'); // only need one.

export default {
  components: { Blockly, Catalog, RouterView, Startup },
  setup(props) {
    let shapeData= ref(null);
    return {
      catalog, // unwatched.
      shapeData, 
      onStarted(_shapeData) {
        shapeData.value= Object.freeze(_shapeData);
      }
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
        const parts= file.path.split("/");
        router.push({ name: 'edit', params: { editPath: parts } });
      },
    }
  }
}
</script>

