<template
><div><lv-markup
  v-for="line in lines"
  :msg="line"
></lv-markup
></div
></template>
<script>
import lvMarkup from './markup/markup.js'

export default {
  emits: ['updated'],
  components: { lvMarkup },
  props: {
    lines: Array 
  },
  updated() {
    this.$emit('updated');
    this.$el.scrollTop = this.$el.scrollHeight;
  },
}
</script>