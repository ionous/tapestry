
Vue.component('mk-opt-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    ><mk-switch
      v-if="childNode"
      :node="childNode"
    ></mk-switch
    ><mk-pick-inline
      v-else
      :node="node"
      :param="param"
      @picked="onPick"
    ></mk-pick-inline
  ></span>`,
  computed: {
    childNode() {
      return this.node.kid;
    },
  },
  methods: {
    onPick(token) {
      const spec= node.itemType.with;
      if (!token in spec.params) {
        throw new Error(`unknown token picked '${token}'`);
      }
      const param= spec.params[token];
      const typeName= param.type || param; // an opt's param can map straight to their type.
      const newNode= this.$root.nodes.newFromType(typeName);
      this.node.setSwap(token, newNode);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: SwapNode,
    param: Object,
    token: String,
  }
});
