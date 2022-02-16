<template><div 
  :class="bemBlock()"
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
  mounted() {
    const { catalog, folder } = this;
    catalog.loadFolder(folder);
  }
}
</script>