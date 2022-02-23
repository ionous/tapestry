<template>
<li
  ><mk-catalog-button
    :class="bemElem('button', open?'open':'closed')"
    :depth="depth"
    @activate="onActivated"
  >{{name}}</mk-catalog-button
  ><slot
  ></slot
></li>
</template>

<script>
import mkCatalogButton from './CatalogButton.vue'
import { CatalogFolder } from './catalogItems.js'
import bemMixin from '/lib/bemMixin.js'

export default {
  mixins: [bemMixin('mk-folder-item')],
  components: { mkCatalogButton },
  props: {
    folder: CatalogFolder,
    depth: Number
  },
  emits: ['activated','opened-file'],
  computed: {
    name() {
      const { folder }= this;
      return folder && folder.name;
    },
    open() {
      const { folder }= this;
      return folder && folder.contents;
    }
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  }
}
</script>
