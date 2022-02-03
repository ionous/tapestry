<template>
<li :class="bemBlock()"
><CatalogButton
  :depth="depth"
  :class="bemElem('button')"
  @activate="onActivated"
>{{name}}</CatalogButton
></li>
</template>
<script>
import CatalogButton from './CatalogButton.vue'
import { CatalogFile } from './catalogItems.js'
import { computed } from 'vue'
import bemMixin from '/src/bemMixin.js'

export default {
  mixins: [ bemMixin("mk-file-item") ],
  components: {
    CatalogButton
  },
  props: {
    file: CatalogFile,
    depth: Number,
  },
  computed: {
    name() {
      const { file }= this;
      const ext= ".if";
      return file && file.name.slice(0, -ext.length);
    },
  },
  emits: ['activated'],
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  }
}
</script>