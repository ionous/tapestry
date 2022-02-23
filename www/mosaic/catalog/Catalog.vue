<template><div 
  :class="bemBlock()"
  ><mk-folder
      v-if="!!folder.contents"
      :folder="folder"
  ></mk-folder
></div></template>

<script>
import mkFolder from './Folder.vue'
import Cataloger from './cataloger.js'
import { CatalogFolder } from './catalogItems.js'
import bemMixin from '/lib/bemMixin.js'

export default {
  mixins: [ bemMixin("mk-catalog") ],
  props: {
    catalog: Cataloger,
  },
  components: { mkFolder },
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