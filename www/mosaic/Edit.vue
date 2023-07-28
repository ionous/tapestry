  <template
  ><mk-startup 
    v-if="!shapeData || !file"
    @started="onStarted"
  /><template v-else
      ><mk-blockly 
        :file="file"
        :shape-data="shapeData"
        :toolbox-data="toolboxData"
        @changed="onChanged"
  /></template
></template>

<script>
import mkBlockly from './blockly/Blockly.vue'
import mkStartup from './Startup.vue'
import appcat from './appcat.js'
import routes from './routes.js'
import endpoint from './endpoints.js'
import http from '/lib/http.js'
import { ref } from 'vue';

export default {
  components: { mkBlockly, mkStartup },
  methods: {  
    // saves on close:
    // todo: autosave periodically
    // todo: save spinner
    onChanged(path, contents) {
      const where = endpoint.file(path);
      const pod = JSON.parse(contents);
      http.put(where, pod, true).then((res)=> {
        console.log("SAVED:", res.ok, res.status, res.statusText);
      }).catch((reason)=>{
        console.log("failed to save", reason);
      });
    },
  },
  // exposes reactive values:
  // https://vuejs.org/api/composition-api-setup.html#composition-api-setup
  setup(props) {
    let file = ref(null);
    let shapeData= ref(null);
    let toolboxData= ref(null);
    
    // load the requested data
    const path = routes.getEditPath();
    appcat.loadFile(path).then((f) => {
      file.value = f;
    }).catch((e) => {
      console.error(e);
      routes.goHome();
    });

    return {
      file,
      shapeData, 
      toolboxData,
      // after we've run the startup, publish the loaded data
      // (here to help avoid watchers )
      onStarted(b) {
        shapeData.value= Object.freeze(b.shapes);
        toolboxData.value= Object.freeze(b.tools);
      }
    }
  }
}
</script>

