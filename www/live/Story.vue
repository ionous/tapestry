<template
><div class="lv-story"
><lv-markup
 	v-for="line in lines"
 	:msg="line"
></lv-markup> 
</div></template>
<script>
import lvMarkup from './markup/markup.js'

export default {
	components: { lvMarkup },
	props: {
		lines: Array 
	},
}
</script>