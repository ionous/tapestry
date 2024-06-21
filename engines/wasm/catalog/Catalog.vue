<template><div 
  class="mk-homebox"
  :class="bemBlock()"
  ><b>Story file browser</b><mk-folder
      v-if="!!folder.contents"
      :folder="folder"
      @fileSelected="onFile"
      @folderSelected="onFolder"
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
  emits: ['fileSelected'],
  data() {
    const { catalog } = this;
    return {
      folder: catalog.root,
    }
  },
  mounted() {
    const { catalog, folder } = this;
    catalog.loadFolder(folder);
  },
  methods: {
    onFile(file) {
      this.$emit('fileSelected', file.path);
    },
    onFolder(folder) {
      if (folder.contents) {
        console.log("closing", folder.name);
        folder.contents= false;
      } else {
        // injects the list of sub-files into the passed folder
        console.log("opening", folder.name, folder.path);
        this.catalog.loadFolder(folder);
      }
    },
  }
}
</script>