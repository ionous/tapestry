<template
	><div 
  :class="bemBlock()"
  ><button 
  	:disabled="!saveEnabled"
  	@click="$emit('save')">Save</button
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
  },
	emits: ['save'],
  mixins: [ bemMixin("mk-tools") ],
}
</script>