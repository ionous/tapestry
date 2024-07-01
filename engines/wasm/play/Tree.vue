<template>
<component :is="list ? 'ul': 'div'"
  :class="traits"
  >
  <component :is="list ? 'li': 'div'">
    <a href="#" 
      v-if="list" role="button" @click="onActivated">{{ item.name }}</a>
    <mk-tree
      v-for="item in item.kids"
      :item="item"
      :key="item.id"
      :list="list"
      @activated="retrigger"
    ></mk-tree>
  </component>
</component>
</template>
<script>
function clean(id) {
  return id.replaceAll(' ', '_')
}
export default {
  name: 'mkTree',
  props: {
    list: {
      Boolean,
      default: false,
    },
    item: {
      type: Object,
      default: {}
    },
  },
  emits: ['activated'],
  computed: {
    traits() {
      const { id = "", kind = "", traits = [] } = this.item;
      return ["id-"+clean(id), "kind-"+clean(kind)].concat(traits.map(clean));
    },
  },
   methods: {
    retrigger(id) { // vue doesn't bubble. *sigh*
      this.$emit('activated', id);
    },
    onActivated() {
      this.$emit('activated', this.item.id);
    }
  }
}

// {
// "id": "",
// "name": "",
// "kind": "",
// "traits": [],
// "kids": []
// }
</script>
