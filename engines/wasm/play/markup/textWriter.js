/**
 * mimics git.sr.ht/~ionous/tapestry/web/text
 */
import TextRender from './textRender.js'

// tagTypes
const opening = "opening";
const closing = "closing";

class TextWriter {
  constructor() {
    this.buf = ""; // temporary storage while parsing strings
    this.tag = "";  // a subset of buf while checking for tags
    this.out = new TextRender();
  }
  // the potential tag we've been accumulating is not a tag:
  // write any text we've been accumulating to the output
  rejectTag() {
    this.tag= "";
    if (this.buf) { // write any accumulated text as regular text.
      this.out.writeString(this.buf);
      this.buf= "";
    }
  }
  // keep track of the passed character as potentially belonging to a tag
  // ex. "s" for <strike>
  accumTag(q) {
    if (this.tag.length < longestTag.length) {
      this.tag += q;
    } else {
      this.rejectTag();
    }
    return !!this.tag;
  }
  // done with reading a <tag>
  // ( note: it may or not be a valid tag. )
  dispatchTag(tagType) {
    const open = tagType === opening;
    if (!this.out.writeTag(this.tag, open)) {
      this.rejectTag();
    } else {
      // we're done with tag, so clear our tracking data.
      this.tag= "";
      this.buf= "";
    }
  }
}

// returns new vue html
export default function writeText(msg) {
  const c = new TextWriter();
  let next = states.readingText;
  for (const q of msg) {
    next = next(c, q); // advance states
  }
  c.rejectTag(); // if anything was pending
  // mimic print.NewLineSentences()
  return c.out.finalize( msg.match( /[\.!?]$/ ));
}

const tagFriendly= /[a-zA-Z]/;
const longestTag = "blockquote";

// functions which take context and a single character.
const states= {
  // reading normal text; look out for the start of tags.
  readingText(c, q) {
    let next = states.readingText;
    // detected the start of a tag:
    if (q === '<') {
      c.buf += q; // start buffering in case this isnt a valid tag.
      next = states.openingTag;
    } else {
      c.out.writeChar(q);
    }
    return next;
  },

  // we've recently seen a "<"
  // parse subsequent characters...
  openingTag(c, q) {
    let next= states.openingTag;
    c.buf += q;
    // could be a closing tag
    if (q==='/') {
      // yes. looks like a closing tag...
      if (!c.tag.length) {
        next = states.closingTag;
      }
      // no: there was text before the slash:
      // ex. <abc/> not </abc>
      else {
        c.rejectTag();
        next= states.readingText;
      }
    }
    // end of this opening tag.
    else if (q == '>') {
      c.dispatchTag(opening)
      next= states.readingText;
    }
    // continuing on in the current tag
    // ex. <abc...
    else if (tagFriendly.test(q)) {
      // ( only fails to accumulate if the tag was too long )
      if (!c.accumTag(q)) {
        next= states.readingText;
      }
    }
    // some other character
    // ( therefore not a tag )
    else {
      c.rejectTag();
      next= states.readingText;
    }
    return next;
  },

  // we've recently seen a "</"
  // parse subsequent characters...
  closingTag(c, q) {
    let next= states.closingTag;
    c.buf += q;
    // end of this closing tag.
    if (q == '>') {
      c.dispatchTag(closing);
      next= states.readingText;
    }
    // continuing on in the current tag
    // ex. </abc...
    else if (tagFriendly.test(q)) {
      if (!c.accumTag(q)) {
        // ( only fails to accumulate if the tag was too long )
        next= states.readingText;
      }
    }
    // some other character
    // ( therefore not a tag )
    else  {
      c.rejectTag();
      next= states.readingText;
    }
    return next;
  },
};
