<template>
  <h1>Catalog</h1>
  <Folder
      v-if="!!folder.contents"
      :folder="folder"
      :class="bemBlock()"
  ></Folder>
</template>

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
    return {
      folder: new CatalogFolder("")
    }
  },
  mounted() { // mounted happens after data() [ which happens in the created hook. ]
    const { catalog, folder } = this;
    catalog.loadFolder(folder);
  },
  provide() {
    const { catalog } = this;
    return {
      onFolder(folder) {
        if (folder.contents) {
          console.log("closing", folder.name);
          folder.contents= false;
        } else {
          // injects the list of sub-files into the passed folder
          console.log("opening", folder.name);
          catalog.loadFolder(folder);
        }
      },
      onFile(file) {
        // injects the story data into the passed file
        console.log("loading", file.name);
        catalog.loadFile(file);
        //emitter.$emit("opened-file", {file});
      },
    }
  },
}
</script>