<template
><div class="lv-prompt"
><label for="prompt">&gt; </label
><input 
  id="prompt"
  class="lv-no-outline lv-blink"
  ref="input"
  name="prompt"
  v-model.trim="text"
  @keyup.enter="onInput"
><i><span class="lv-hidden">{{text}}</span></i
></div></template>
<script>

import CommandHistory from './commandHistory.js'

// fix? would it make more sense to just post using form?
export default {
  emits: ["changed"],
  data() {
    // const testHistory= ["one", "two", "three", "four", "five"];
    return {
      history: new CommandHistory([], 25),
      text: ""
    }
  },
  mounted() {
    this.setFocus();
  },
  methods: {
    setFocus() {
      const el= this.$refs.input;
      if (document.activeElement !== el) {
        console.log("set focus");
        el.focus();
        const end= this.text.length;
        el.setSelectionRange(end,end);
        el.scrollIntoView();
      }
    },
    browseHistory(upArrow) {
      console.log("browseHistory", upArrow);
      const nextText= this.history.browse(upArrow);
      this.text= nextText || "";
      if (nextText) {
        const el= this.$refs.input;
        const len= nextText.length;
        // note: vue's nextTick doesnt work for this :/
        // and using the timeout appears to cause a small flicker *sigh*
        setTimeout(() => {
          el.setSelectionRange(0, len);
        })
      }
    },
    scrollIntoView() {
      this.$refs.input.scrollIntoView(true);
    },
    onInput() {
      const txt= this.text;
      this.text ="";
      this.$emit("changed", txt);
      this.history.remember(txt);
    }
  }
}
</script>