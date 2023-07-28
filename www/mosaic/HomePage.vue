<!-- 
    used by the router state to view the landing page.
-->
<template>
  <div class="mk-home">
    <div class="mk-actions mk-homebox">
      <ul>
        <mk-action-item
            text="Create new story..."
            @activated="onCreateFile"></mk-action-item>
        <mk-action-item
            text="Open existing story..."
            @activated="onOpenFile()"></mk-action-item>
      </ul>
    </div>
    <mk-catalog 
      :catalog="appcat"
      @fileSelected="onOpenFile"/>
  </div>
</template>

<script>
// hierarchical story file/folder control
import mkCatalog from './catalog/Catalog.vue'
// a link like component with an @activated event 
import mkActionItem from './catalog/ActionItem.vue'
import routes from './routes.js'
import appcat from './appcat.js'


export default {
  // exported by the provide() of Mosaic.vue:
  // the global catalog, actions for buttons to create and open an os file dialog.
  //
  components: { mkCatalog, mkActionItem },
  data() {
    return {
      appcat
    }
  },
  methods: {
    onCreateFile() {
      routes.newFile();
    },
    onOpenFile(file) {
      routes.openFile(file);
    },
  }
}
</script>


