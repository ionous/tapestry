<template>
<li
  ><CatalogButton
    :class="bemElem('button', open?'open':'closed')"
    :depth="depth"
    @activate="onActivated"
  >{{name}}</CatalogButton
  ><slot
  ></slot
></li>
</template>

<script>
import CatalogButton from './CatalogButton.vue'
import { CatalogFolder } from './catalogItems.js'
import bemMixin from '/src/bemMixin.js'

export default {
  mixins: [bemMixin('mk-folder-item')],
  components: {
    CatalogButton
  },
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
