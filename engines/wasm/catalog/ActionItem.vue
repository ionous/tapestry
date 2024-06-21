<template>
<li :class="bemBlock()"
><mk-catalog-button
  :depth="0"
  :class="bemElem('button')"
  @activate="onActivated"
>{{text}}</mk-catalog-button
></li>
</template>
<script>
import mkCatalogButton from './CatalogButton.vue'
import bemMixin from '/lib/bemMixin.js'

export default {
  mixins: [ bemMixin("mk-action-item") ],
  components: { mkCatalogButton },
  props: {
    text: String,
  },
  emits: ['activated'],
  methods: {
    onActivated() {
      this.$emit("activated", this.text);
    },
  }
}
</script>