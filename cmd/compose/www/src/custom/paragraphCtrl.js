class StatementTable extends NodeTable {
  constructor(redux, para) {
    super(redux, para, para.getKid("$STORY_STATEMENT"));
    this.inline= true;
  }
  makeBlank() {
    return this.nodes.newFromType("story_statement");
  }
  // when we drag, we re/move everything from a given statement till the end of line.
  // returns a list of statements
  removeFrom(at, width= Number.MAX_VALUE) {
    return this.items.splice(at, width).map(el=> {
      el.parent= null;
      return el;
    });
  }
  // add a paragraph, or a line of statements
  // at the line of statements targeted
  addTo(at, paraEls) {
    let ret;
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      // tack its elements to the end of the targeted line
      const para= paraEls;
      // remove all the kids from their parent array
      const kids= para.getKid("$STORY_STATEMENT");
      const els= kids.splice(0);
      ret= this.addTo( at, els );
    } else {
      const els= paraEls;
      items.splice(at, 0, ...els.map(el=> {
        el.parent= node;
        return el;
      }));
      ret= els.length;
    }
    return ret;
  }
}

// paragraphs are actually, basically, the discrete lines of a story.
Vue.component('mk-paragraph-ctrl', {
  template:
  `<em-table
      :class="$root.shift && 'em-shift'"
      :list="list"
      :dropper="dropper"
      :grip="'\u2630'"
  ><template
      v-slot="{item, idx}"
    ><mk-switch
      :node="item"
    ></mk-switch
    ></template
  ></em-table>`,
  props: {
    node: Node,
  },
  data() {
    const { node, "$root": root } = this;
    // each item is a story statement slot
    return {
      list: new StatementTable(root.redux, node),
      dropper: root.dropper,
    }
  },
});