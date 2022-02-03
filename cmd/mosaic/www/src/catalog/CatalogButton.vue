<!-- a link that acts like a button. -->
<template>
<a class="mk-cat-button"
  :style="{ paddingLeft: (10*depth)+'px' }"
  role="button"
  @click="onActivated"
  @keydown.prevent.space="onActivated"
  @keydown.prevent.enter="onActivated"
  ><slot
  ></slot
></a>
</template>
<script>
export default {
  props: {
    depth: Number, // depth helps to indent buttons
  },
  emits:[
    'activate'     // uniform event regardless of how the button was pressed.
  ],
  methods: {
    onActivated() {
      // was used for clearing global selection:
      // this.$root.events.$emit("mk-button-activated");
      this.$emit('activate');
    }
  }
}
</script>
