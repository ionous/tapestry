<template
><lv-markup
  v-for="line in lines"
  :msg="line"
></lv-markup
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
    this.$el.parentElement.scrollTop = this.$el.parentElement.scrollHeight;
  },
}
</script>