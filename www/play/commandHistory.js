export default class CommandHistory {
  constructor(els, maxHistory) {
    // els is a circular buffer.
    this.els= els;
    this.max= maxHistory;
    this.next= 0; // spot for next memento
    this.at=0; // for browsing existing history
  }

  // new elements are pushed
  remember(el, resetCursor=true) {
    console.log("remember", el);
    const next= this.next;
    this.els[next]= el; // javascript autoexpands the array if need be.
    this.next= (next+1) % this.max;
    if (resetCursor) {
      this.at=0;
    }
  }

  // return the (next) most recent item in the history.
  browse( older=true ) {
    var ret;
    const { els, next, max } = this;
    // cant browse history if there's nothing there.
    if (els.length) {
      // at ranges from 0 to (max) history
      let at= this.at + (older?1:-1);
      if (at< 0) {
        at= 0;
      } else if (at> els.length) {
        at= els.length;
      }
      this.at= at; // update the validated cursor pos
      if (at) {
        let i= next-at; // 1 history, 1 cursor= the 0th element.
        if (i < 0) { // noting that els is a circular buffer.
          i += max;
        }
        ret= els[i];
      }
    }
    return ret;
  }
};
