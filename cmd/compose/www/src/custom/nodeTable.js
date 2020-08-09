// nodeList.js
class NodeTable extends DragList {
  constructor( redux, node, items ) {
    super(items);
    this.redux= redux;
    this.nodes= redux.nodes;
    this.node= node;
  }
  makeBlank() {
    throw new Error("not implemented");
  }
  dropFrom(to, from) {
    const thisList= this.list;
    const thatList= from.list;
    const toIdx= to.idx;
    const fromIdx= from.idx;

    if (thisList !== thatList) {
      const needBlank= Math.abs(fromIdx - toIdx) === 1;
      if (needBlank) {
        this.addBlank(toIdx);
      } else {
        this.move(toIdx, fromIdx);
      }
    } else {
      const { redux } = this;
      redux.invoke({
        apply() {
          const paraEls= thatList.removeFrom(fromIdx);
          thisList.addTo( toIdx, paraEls );
        },
        revoke() {
          const paraEls= thisList.removeFrom( toIdx );
          thatList.addTo( fromIdx, paraEls );
        },
      });
    }
  }
  addBlank(at) {
    const { redux, items } = this;
    const blank= this.makeBlank();
    redux.invoke({
      apply() {
        items.splice(at,0,blank);
      },
      revoke() {
        items.splice(at,1);
      },
    });
  }
  // move items within this same list
  move(src, dst, width, nothrow) {
    const { redux, items } = this;
    if (width<=0) {
      const e= new Error("invalid width");
      if (nothrow) { return e; }
      throw e;
    }
    if ((dst > src) && (dst < src+width)) {
      const e= new Error("invalid dest");
      if (nothrow) { return e; }
      throw e;
    }
    if (src+width> items.length) {
      width= items.length-src;
    }
    if (dst > src) {
      dst -= width;
    }
    redux.invoke({
      apply() {
        const rub= items.splice(src, width);
        items.splice(dst, 0, ...rub);
      },
      revoke() {
        const rev= items.splice(dst, width);
        items.splice(src, 0, ...rev);
      }
    });
  }
}
