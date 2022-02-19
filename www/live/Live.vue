<template
><lv-results
/><lv-status
/><lv-story :lines="lines"
/><lv-prompt @changed="onPrompt"
/></template>

<script>

import lvPrompt from './Prompt.vue'
import lvResults from './Results.vue'
import lvStatus from './Status.vue'
import lvStory from './Story.vue'

const lines= [
`<b>should write in bold</b>`,
`<div>should write the full tags</div>`,
`<i>should write the italics</i>`,
`<s>should strike</s>`,
`<b>should</s> skip a bad tag</b>`,
`<p>one<p>paragraph<p>then another.`,
`<strong>should close <em>trailing tags <strike>if needed`,
`<hr>`,
`<ol><li>one item</li><ul><li>second item</li></ul></ol>`,
`text with
new line`,
`one line<br>then another<br><wbr>no blank lines<wbr>another new line<br><br>one blank line`,
];

export default {
	components: { lvPrompt,lvResults,lvStatus,lvStory },
	data() {
		return {
			lines,
		}
	},
	methods: {
		onPrompt(txt) {
			console.log("live", txt);
		},
	},
}
</script> 