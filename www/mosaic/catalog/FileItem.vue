<template>
<li :class="bemBlock()"
><mk-catalog-button
  :depth="depth"
  :class="bemElem('button')"
  @activate="onActivated"
>{{name}}</mk-catalog-button
></li>
</template>
<script>
import mkCatalogButton from './CatalogButton.vue'
import { CatalogFile } from './catalogItems.js'
import { computed } from 'vue'
import bemMixin from '/lib/bemMixin.js'

export default {
  mixins: [ bemMixin("mk-file-item") ],
  components: { mkCatalogButton },
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