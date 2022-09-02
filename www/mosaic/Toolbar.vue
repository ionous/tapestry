<template
  ><div 
  :class="bemBlock()"
  ><button 
    :disabled="!saveEnabled"
    @click="$emit('save')">Save</button
  >&nbsp;<button
    :disabled="!playEnabled"
    @click="$emit('play')">Play</button
  ></div
></template>

<script>
import bemMixin from '/lib/bemMixin.js'

export default {
  props: {
    saveEnabled: {
      type: Boolean,
      default: true
    },
    playEnabled: {
      type: Boolean,
      default: true
    }
  },
  emits: ['save', 'play'],
  mixins: [ bemMixin("mk-tools") ],
}
</script>