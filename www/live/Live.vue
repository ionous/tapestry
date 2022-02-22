<template
><lv-output
  class="lv-results"
  :lines="logging"
/><div v-if="narration"
  ><lv-status
  /><lv-output
    class="lv-story"
    :lines="narration"
  /><lv-prompt 
    @changed="onPrompt"
    ref="prompt"
  /></div
></template
><script>

import lvPrompt from './Prompt.vue'
import lvOutput from './Output.vue'
import lvStatus from './Status.vue'
import Io from './io.js'  
import { ref, onMounted, onUnmounted, reactive } from 'vue';

const junk= JSON.stringify([
`<b>should write in bold</b>`,
`<div>should write the full tags</div>`,
`<i>should write the italics</i>`,
`<s>should strike</s>`,
`<b>should</s> skip a bad tag</b>`,
`<p>one<p>paragraph<p>then another.`,
`<strong>should close <em>trailing tags <strike>if needed`,
`<hr>`,
`<ol><li>one item</li><ul><li>second item</li></ul></ol>`,
`text with
new line`,
`one line<br>then another<br><wbr>no blank lines<wbr>another new line<br><br>one blank line`,
]);


export default {
  components: { lvOutput,lvPrompt,lvStatus },
  // props: {},
  setup(/*props*/) {
    const logging= reactive([]);//JSON.parse(junk));
    const narration= reactive([]);
    const prompt= ref(null); // template ref 
    const io= new Io(appcfg.live+"/io/", (msgs)=>{
      for (const msg of msgs) {
        for (const k in msg) {
          const v= msg[k];
          switch (k) {
            case "Play log:":
              console.info("log:", v);
              logging.push(v);
            break;
            case "Play out:":
              narration.push(v);
              const el= prompt.value;
              if (el) {
                el.scrollTo();
              }
            break;
            case "Play mode:":
              console.warn("mode:", v);
              switch (v) {
                case 'asm':
                break 
                case 'play':
                break
                case 'complete':
                break;
                case 'error':
                break;
              }
            break;
            default:
              console.error("unknown command", k);
            break;
          }
        }
      }
    });
    const onkey= (evt) => {
      // console.log("key", evt.key);
      const ignore= (evt.metaKey || evt.ctrlKey || evt.altKey);
      if (!ignore) {
        const el= prompt.value;
        if (el) {
          if (el !== document.activeElement) {
            el.setFocus();
          }
          switch (evt.key) {
            case 'ArrowUp':
              el.browseHistory(true);
            break;
            case 'ArrowDown':
              el.browseHistory(false);
            break;
          }
        }
      }
    };
    onMounted(()=> {
      document.body.addEventListener("keydown", onkey);
      io.startPolling();
    });
    onUnmounted(()=> {
      document.body.removeEventListener("keydown", onkey);
      io.stopPolling();
    });
    return {
      narration, // story output 
      logging,   // story debugging
      prompt,   // template ref
      onPrompt(txt) {
        console.log("onPrompt");
        narration.push("> "+ txt);
        // fix? tapestry commands?
        io.send({
          cmd: txt,
        });
      }
    }
  }
}
</script> 