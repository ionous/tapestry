<template>
  <div ref="container">
    <lv-output class="lv-results" :lines="logging" />
    <div v-if="narration">
      <lv-status /><lv-output class="lv-story" :lines="narration" /><lv-prompt
        enabled="playing"
        @changed="onPrompt"
        ref="prompt"
      />
    </div>
  </div>
</template>
<script>
import lvPrompt from "./Prompt.vue";
import lvOutput from "./Output.vue";
import lvStatus from "./Status.vue";
import Io from "./io.js";
import { ref, onMounted, onUnmounted } from "vue";

const junk = JSON.stringify([
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
  components: { lvOutput, lvPrompt, lvStatus },
  // props: {},
  setup(/*props*/) {
    const logging = ref([]); //JSON.parse(junk));
    const narration = ref([]);
    const playing = ref(false);
    const prompt = ref(null); // template ref
    const container = ref(null);

    function addToNarration(msg) {
      narration.value.push(msg);
      setTimeout(() => {
        const el = container.value;
        el.scrollIntoView(false);
      });
    }
    function processEvent(msg) {
      let out = "";
      const sig = Object.keys(msg)[0];
      const body = msg[sig];
      switch (sig) {
        case "StateChanged noun:aspect:trait:":
          const [noun, aspect, trait] = body;
          console.log("state changed", noun, aspect, trait);
          break;
        case "FrameOutput:":
          out += body;
          break;
        default:
          console.log("unhandled", sig);
      }
      return out;
    }
    // by default, sends commands to http://localhost:8080/shuttle/
    const io = new Io(appcfg.shuttle, (msgs) => {
      let out = "";
      if (typeof msgs === 'string') {
        console.error(msgs);
        return;
      }
      for (const msg of msgs) {
        const sig = Object.keys(msg)[0];
        const body = msg[sig];
        if (sig.endsWith("error:")) {
          console.warn(body);
        } else {
          const [result, events] = body;
          for (const evt of events) {
            out += processEvent(evt);
          }
        }
      }
      if (out.length) {
        addToNarration(out);
      }
    });
    // fix: add a button? read from the path or query string?
    io.post("restart", "cloak");
    const onkey = (evt) => {
      // console.log("key", evt.key);
      const ignore = evt.metaKey || evt.ctrlKey || evt.altKey;
      if (!ignore) {
        const el = prompt.value;
        if (el) {
          if (el !== document.activeElement) {
            el.setFocus();
          }
          switch (evt.key) {
            case "ArrowUp":
              el.browseHistory(true);
              break;
            case "ArrowDown":
              el.browseHistory(false);
              break;
          }
        }
      }
    };
    onMounted(() => {
      document.body.addEventListener("keydown", onkey);
      io.startPolling();
    });
    onUnmounted(() => {
      document.body.removeEventListener("keydown", onkey);
      io.stopPolling();
    });
    return {
      narration, // story output
      logging, // story debugging
      prompt, // template ref
      container,
      onPrompt(text) {
        console.log("onPrompt");
        narration.value.push("> " + text);
        const msg = [{
          "FromExe:": {
            "Fabricate input:":text
          }
        }];
        io.post("query", msg);
      },
    };
  },
};
</script>
