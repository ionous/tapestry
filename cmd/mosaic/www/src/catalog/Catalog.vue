<template><div 
  :class="bemBlock()"
  ><h1>Catalog</h1
  ><Folder
      v-if="!!folder.contents"
      :folder="folder"
  ></Folder
></div></template>

<script>
import Folder from './Folder.vue'
import Cataloger from './cataloger.js'
import { CatalogFolder } from './catalogItems.js'
import bemMixin from '/src/bemMixin.js'

export default {
  mixins: [ bemMixin("mk-catalog") ],
  props: {
    catalog: Cataloger,
  },
  components: {
    Folder,
  },
  data() {
    const { catalog } = this;
    return {
      folder: catalog.root,
    }
  },
  created() {
    this.$watch(
      () => this.$route.params,
      (toParams, previousParams) => {
        this.syncToRoute();
      }
    )
  },
  mounted() { // mounted happens after data() [ which happens in the created hook. ]
    const { catalog, folder } = this;
    catalog.loadFolder(folder);
  },
  methods: {
    // fix: this really will only work with the mock catalog.
    // fix: will want to check if a file exists *before* loading...
    syncToRoute() {
      const { "$router": router, "$route": route, catalog } = this;
      const path= route.params.editPath || [];
      for (let i=0; i< path.length-1; i++) {
        const sub= path.slice(0, i+1).join('/');
        const f= catalog.findByPath(sub);  
        catalog.loadFolder(f);
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
        // injects the story data into the passed file
        console.log("loading", file.name);
        catalog.loadFile(file);
        const parts= file.path.split("/");
        router.push({ name: 'edit', params: { editPath: parts } })

        //emitter.$emit("opened-file", {file});
      },
    }
  },
}
</script>