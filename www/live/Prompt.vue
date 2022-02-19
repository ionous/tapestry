<template
><div class="lv-prompt"
><input 
	class="lv-no-outline"
	ref= "input"
	v-model.trim="text"
	@change="onInput"
></div></template>
<script>
export default {
	emits: ["changed"],
	props: {
		len: Number, // we use the side effect of update to stay in view
	},
	data() {
		return {
			text: "",
		}
	},
	mounted() {
		this.$refs.input.focus();
  },
  updated() {
		this.$refs.input.scrollIntoView(true);
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