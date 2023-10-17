<template>
  <div class="lv-container">
    <div class="lv-container__debug">
      <mk-folder
        :folder="enclosure"
        @fileSelected="onLeafObject"
        @folderSelected="onEnclosingObject"
      ></mk-folder>
      <b>Traits:</b>
      <div class="lv-traits">
        <span v-for="(trait,index) in currentTraits">{{index?", ":""}}{{trait}}</span>
      </div>
    </div>
    <div class="lv-container__play" 
      v-if="narration"
      @click="onContainerClicked">
      <lv-status 
        :status="status"/>
      <lv-output 
        class="lv-story" 
        :lines="narration" />
      <lv-prompt
        enabled="playing"
        @changed="onPrompt"
        ref="prompt" />
    </div>
  </div>
</template>
<script>
import lvPrompt from "./Prompt.vue";
import lvOutput from "./Output.vue";
import lvStatus from "./Status.vue";
import Status from "./status.js";

import Io from "./io.js";
import { ref, onMounted, onUnmounted } from "vue";
//
import ObjectCatalog from './objectCatalog.js'
import mkFolder from '/mosaic/catalog/Folder.vue'

import cmds from "./cmds.js";

const objCatalog = new ObjectCatalog();

export default {
  components: { lvOutput, lvPrompt, lvStatus, mkFolder },
  // props: {},
  setup(/*props*/) {
    const narration = ref([]);
    const status = ref(new Status());
    const playing = ref(false);
    const prompt = ref(null); // template ref
    const enclosure = ref(objCatalog.root);
    const currentTraits = ref([]);

    function addToNarration(msg) {
      narration.value.push(msg);
    }
    //  given a valid tapestry command, 
    // return its signature and body in an array of two elements
    function parseCommand(op) {
      for (const k in op) {
        if (k!== "--") {
          return [k, op[k]];
        }
      }
    }
    function processEvent(msg) {
      let out = "";
      const [ sig, body ] = parseCommand(msg);
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
    const io = new Io(appcfg.shuttle, (msgs, calls) => {
      let out = "";
      if (typeof msgs === 'string') {
        console.error(msgs);
        return;
      }
      for (let i=0; i< msgs.length; ++i) {
        const msg = msgs[i];
        const call = calls[i];
        //
        const [ sig, body ] = parseCommand(msg);
        switch (sig) {
          // fix: result and events should probably be optional;
          // or, make two commands that satisfy some response interface
          case "Frame result:events:error:":
          {
            const [_res, _evts, error] = body;
            console.warn(error);
            break;
          }
          case "Frame result:events:":
          {
            const [result, events] = body;
            if (events) {
              for (const evt of events) {
                out += processEvent(evt);
              }
            }
            if (call) {
              // ick: we debug.Stringify the results to support "any value"
              // so we have to unpack that too.
              const res = result? JSON.parse(result): "";
              call(res);
            }
            break;
          }
          default:
            console.log("unhandled", sig);
        };
      }
      if (out.length) {
        addToNarration(out);
      }
    });
    io.post("restart", "cloak").then(()=> {
      io.query([
      cmds.storyTitle, (title)=>{
        status.value.title = title;
      },
      cmds.currentScore, (score)=>{
        status.value.score = score;
      },
      cmds.currentTurn, (turn)=>{
        status.value.useScoring = turn >= 0;
        status.value.turns = turn;
      },
      cmds.locationName, (name)=>{
        status.value.location = name;
      },
      cmds.currentObjects, (objs)=>{
        enclosure.value = objCatalog.rebuild(objs);
      }]);
    });
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
    });
    onUnmounted(() => {
      document.body.removeEventListener("keydown", onkey);
    });
    let first = true; 
    return {
      narration, // story output
      status,
      prompt, // template ref
      enclosure,
      currentTraits,
      onLeafObject(item) {
        currentTraits.value = item.data.traits;
      },
      onEnclosingObject(items) {
        if (!items.contents) {
          items.contents = items.backup;
          items.backup = false;
          currentTraits.values = item.data.traits;
        } else {
          items.backup = items.contents;
          items.contents = false;
          currentTraits.values = [];
        }
      },
      // clicking anywhere below the prompt should focus the prompt
      onContainerClicked() {
        const el = prompt.value;
        if (el && el !== document.activeElement) {
          el.setFocus();  
        }
      },
      onPrompt(text) {
        console.log("onPrompt");
        if (!first) {
          // patch for an errant space after the first set of text
          // or there could always be a break
          narration.value.push("<br>");
        }
        first = false;
        narration.value.push("> " + text);
        const msg = [{
          "FromExe:": {
            "Fabricate input:":text
          }
        }];
        io.query(msg, null);
      },
    };
  },
};
</script>
