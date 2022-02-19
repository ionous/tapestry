// markup roughly matches the golang `git.sr.ht/~ionous/tapestry/web/text`.
// It takes Tapestry's pseudo-html markup and making it into real html.
// ( it's not clear this is *actual* needed but feels like a better practice
// than accepting game and user generated text directly. )
import writeText from './textWriter.js';

export default {
  props: {
    msg: String,
  },
  render() {
    return writeText(this.msg);
  }
}
