<template
><div ref="container"
  ><lv-output
    class="lv-results"
    :lines="logging"
  /><div v-if="narration"
    ><lv-status
    /><lv-output 
      class="lv-story"
      :lines="narration"
    /><lv-prompt
      enabled="playing"
      @changed="onPrompt"
      ref="prompt"
    /></div
></div
></template
><script>

import lvPrompt from './Prompt.vue'
import lvOutput from './Output.vue'
import lvStatus from './Status.vue'
import Io from './io.js'  
import { ref, onMounted, onUnmounted } from 'vue';

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
    const logging= ref([]);//JSON.parse(junk));
    const narration= ref([]);
    const playing= ref(false);
    const prompt= ref(null); // template ref 
    const container= ref(null);

    function addToNarration(msg) {
      narration.value.push(msg);
      setTimeout(()=> {
        const el= container.value;
        el.scrollIntoView(false);
      });
    }
    // by default, sends commands to http://localhost:8080/shuttle/
    const io= new Io(appcfg.shuttle, (msg)=>{
      addToNarration(msg);
      // for (const msg of msgs) {
      //   for (const k in msg) {
      //     const v= msg[k];
      //     switch (k) {
      //       case "Play log:":
      //         console.info("log:", v);
      //         logging.value.push(v);
      //         setTimeout(()=> {
      //           const el= container.value;
      //           el.scrollIntoView(false);
      //         });
      //       break;
      //       case "Play out:":
      //          addToNarration(msg);
      //       break;
      //       case "Play mode:":
      //         console.warn("mode:", v);
      //         playing.value= (v === 'play');
      //         switch (v) {
      //           case 'asm':
      //             logging.value= [];
      //             narration.value= [];
      //           break
      //           case 'play':
      //           break
      //           case 'complete':
      //           break;
      //           case 'error':
      //           break;
      //         }
      //       break;
      //       default:
      //         console.error("unknown command", k);
      //       break;
      //     }
      //   }
      // }
    });
    // fix: add a button? read from the path or query string?
    io.post("restart", "cloak");
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
      container,
      onPrompt(txt) {
        console.log("onPrompt");
        narration.value.push("> "+ txt);
        // fix? tapestry commands?
        io.post("input", txt);
      }
    }
  }
}
</script> 