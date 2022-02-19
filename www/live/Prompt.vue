<template
><div class="lv-prompt"
><input 
	class="lv-no-outline"
	v-focus
	v-model.trim="text"
	@change="onInput"
></div></template>
<script>
export default {
	emits: ["changed"],
	data() {
		return {
			text: "",
		}
	},
	directives: {
    focus: {
 		 mounted: (el) => el.focus()
		}
  },
	methods: {
		onInput() {
			const txt= this.text;
			this.text ="";
			this.$emit("changed", txt);
		}
	}
}
</script>