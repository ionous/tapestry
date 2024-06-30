<template>
<component :is="list ? 'ul': 'div'"
  :id="id"
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
      default: {
        id: "",
        traits: []
      }
    },
  },
  emits: ['activated'],
  computed: {
    id() {
      const { id = "_missing_" } = this.item;
      return clean(id);
    },
    traits() {
      const { traits = [], kind = "" } = this.item;
      return ["kind-"+clean(kind)].concat(traits.map(clean));
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
