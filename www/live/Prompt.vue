<template
><div class="lv-prompt"
><label for="prompt">&gt; </label
><input 
	class="lv-no-outline lv-blink"
	ref="input"
	name="prompt"
	v-model.trim="text"
	@keyup.enter="onInput"
><i><span class="lv-hidden">{{text}}</span></i
></div></template>
<script>
// fix? would it make more sense to just post using form?
export default {
	emits: ["changed"],
	data() {
		return {
			text: ""
		}
	},
	mounted() {
		this.setFocus();
  },
  // when props change.
  updated() {
		this.scrollTo();
  },
	methods: {
		setFocus() {
			const el= this.$refs.input;
			if (document.activeElement !== el) {
				el.focus();
				const end= this.text.length;
				el.setSelectionRange(end,end);
			}
		},
		scrollTo() {
			this.$refs.input.scrollIntoView(true);
		},
		onInput() {
			const txt= this.text;
			this.text ="";
			this.$emit("changed", txt);
		}
	}
}
</script>