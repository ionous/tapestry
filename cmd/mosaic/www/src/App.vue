<template><div class="mk-app"
><img alt="Vue logo" src="./assets/logo.png" 
  /><router-view/><Catalog :catalog="catalog" 
  /><Blockly></Blockly
></div></template>

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

