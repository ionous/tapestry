import { h } from 'vue'
import selectTag from './selectTag.js'

export default class TextRender {
  constructor() {
    this.stack= [];  // contains objects of {tag:,els:}
    this.curr= {tag:'div', els:[]};
    this.str= "";  // text for the current span
    this.line= 2; // vertical spacing ( ie. are we at the start of a new line )
                  // we start it at 2, to indicate we are at the top of a fresh section.
  }
  finalize() {
    this.flushString();
    while (this.stack.length) {
      this.rollup();
    }
    const c= this.curr;
    this.curr= null;
    return h(c.tag, c.els);
  }
  // write a string without worrying about whether it contains tags or not
  writeString(s) {
    for (const q of s) {
      this.writeChar(q);
    }
  }
  // handles opening and closing tags
  writeTag(tag, open)  {
    // there are a few possible operations:
    // 1. write a standalone tag: ex. newline, or hr.
    // 2. open a tag
    // 3. close a tag
    // 4. ignore unknown tags
    const rule= selectTag(tag);
    if (!rule) {
      // console.log("skipping unknown tag", tag);
    } else if (rule.char) {
      this.writeChar(rule.char);
    } else if (rule.void) { // self-closing tags.
      this.render(rule.tag);
    } else if (open) {
      this.push(rule.tag);
    } else {
      this.pop(rule.tag);
    }
    return !!rule; // was there a rule?
  }
  // watches for newline characters, etc.
  // this *doesnt* watch for tags or anything like that.
  writeChar(q) {
    switch (q) {
      case '\n': // newline
        this.render('br', {class:'n'});
        ++this.line;
      break;
      case '\r': // soft-line
      // FIX? a better way for html might be to hold the pending line type
      // until we get the next character or an explicit flush
      // that way we can write <p> when appropriate
        while(this.line<1) {
          this.render('br', {class:'r'});
          ++this.line;
        }
      break;
      case '\v':       // soft-paragraph
        while(this.line<2) {
          this.render('br', {class:'v'});
          ++this.line;
        }
      break;
      // everything else:
      // note: we dont have to worry about tabs in the javascript/html version.
      default:
        this.str+= q
        this.line= 0;
      break;
    }
  }
  // add a self-closing tag to the current span; doesnt start a new span.
  render(tag, ...args) {
    this.flushString();
    this.curr.els.push( h(tag, ...args) );
  }
  // start a new "child" span of the current span
  push(tag) {
    this.flushString();
    this.stack.push(this.curr);
    this.curr= {tag, els:[]};
  }
  // end the current span *if* the tag matches
  // ( ex. <b></s></b> skips the errant /s... )
  pop(tag) {
    if (this.curr.tag===tag) {
      this.flushString();
      this.rollup();
    }
  }
  rollup() {
    const c= this.curr; // write the current data into its parent
    this.curr= this.stack.pop();
    this.curr.els.push( h(c.tag, c.els) );
  }
  flushString() {
    if (this.str) {
      this.curr.els.push(this.str);
      this.str = "";
    }
  }
}
