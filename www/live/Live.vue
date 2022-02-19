<template
><lv-results
/><lv-status
/><lv-story 
	:lines="lines"
/><lv-prompt 
	@changed="onPrompt"
	:len="lines && lines.length"
/></template>

<script>

import lvPrompt from './Prompt.vue'
import lvResults from './Results.vue'
import lvStatus from './Status.vue'
import lvStory from './Story.vue'
import Io from './io.js'  

import { reactive } from 'vue'

const junk= JSON.stringify([
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
]);

export default {
	components: { lvPrompt,lvResults,lvStatus,lvStory },
	// props: {},
	setup(/*props*/) {
		const lines= reactive(JSON.parse(junk));
		let io= new Io(lines, appcfg.live+"/io/");
		return {
			lines,
			io,
			onPrompt(txt) {
				lines.push("> "+ txt);
				// fix? tapestry commands?
				io.send({
					cmd: txt,
				});
			}
		}
	},
	mounted() {
		this.io.startPolling();
	},
	unmounted() {
		this.io.stopPolling();
	},
}
</script> 