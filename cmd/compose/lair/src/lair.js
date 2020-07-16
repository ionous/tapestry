// use a pair of numbers for the gutter to manage the sizing.
Vue.component('elem', {
  template:
  `<em-table
      :items="subitems"
      :dropper="dropper"
      :grip="'$'"
    ><template
      v-slot="{item,idx}"
      ><div
      >{{item.text}}</div
    ></template
  ></em-table>`,
  props: {
    idx:Number,
    item:Object,
    dropper:Object,
  },
  data(){
    const id= this.item.id;
    const text= this.item.text;
    const subitems= text.split(". ").map((x, i) => {
      // fake items:
      return {
        id: `${id}-${i}`,
        text: x,
      }
    });
    return {
      subitems: subitems,
    }
  }
});

// use a pair of numbers for the gutter to manage the sizing.
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-max"
    >{{max}}</div
    ><div class="em-num"
    >{{grip || num}}</div
  ></div>`,
  props: {
    grip: String,
    num: Number,
    max: Number,
  },
  beforeDestroy() {
    console.log(`gutter ${this.num}/${this.max} being destroyed`);
  }
});

Vue.component('em-table', {
  data() {
    const list= new DragList(this.items, ()=> new Lipsum());
    const dropper= this.dropper;
    return {
      drag: dropper.newGroup({
        serializeItem(at)  {
          const item= list.items[at];
          return {
            'text/plain': item.text,
          };
        },
        removeItem(src, dst, width=1, newGroup) {
          return list.removeFrom(src, dst, width, newGroup);
        },
        // note: addItem might happen in a group other than serialize and removeItem.
        addItem(src, dst, rub, newGroup) {
          list.addTo(src, dst, rub, newGroup);
        },
      })
    }
  },
  props: {
    grip:String,
    items: Array,
    dropper: Dropper,
  },
  mounted() {
    this.drag.handler.listen(this.$el);
  },
  beforeDestroy() {
    this.drag.handler.silence();
  },
  template:
  `<div class="em-table"
    ><div
      class="em-table__header"
      :data-drag-idx="-1"
      :data-drag-edge="0"
    ></div
    ><div v-for="(item,idx) in items"
        :class="drag.highlight(idx)"
        :data-drag-idx="idx"
        :key="item.id"
      ><em-gutter
        :num="idx+1"
        :grip="grip"
        :max="items.length"
        draggable="true"
        ></em-gutter
      ><slot
        :idx="idx"
        :item="item"
      ></slot
      ></em-gutter
    ></div
    ><div
      class="em-table__footer"
      :data-drag-idx="items.length"
      :data-drag-edge="items.length-1"
    ></div
  ></div>`
});
//
const app= new Vue({
  el: '#app',
  created() {
    document.addEventListener("keydown", (e) => {
      console.log("keydown", e.key === "Shift");
      this.shift= true;
    });
    document.addEventListener("keyup", (e) => {
      console.log("keyup", e.key === "Shift");
      this.shift= false;
    });
    window.addEventListener("blur", (e) => {
      this.shift= false;
    });
  },
  data: {
    groups: [
          Lipsum.list(15, 31, 3, 5, 8, 17),
          Lipsum.list(8, 12, 5, 42, 2, 17),
    ],
    dropper: new Dropper(),
    shift: false,
  },
});
