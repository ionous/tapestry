<template>
  <div class="lv-container">
    <div class="lv-container__debug">
      <mk-folder
        :folder="objTree"
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

import { ref, onMounted, onUnmounted, reactive } from "vue";
//
import ObjectCatalog from './objectCatalog.js'
import mkFolder from '/mosaic/catalog/Folder.vue'
import Query from './query.js'

const objCatalog = new ObjectCatalog();

export default {
  components: { lvOutput, lvPrompt, lvStatus, mkFolder },
  // props: {},
  setup(/*props*/) {
    const narration = ref([]);
    const status = ref(new Status());
    const playing = ref(false);
    const prompt = ref(null); // template slot helper.
    const currentTraits = ref([]);
    // replace the obj catalog's root with a vue reactive proxy
    // so that when we change its children ( move between rooms )
    // vue can see those changes.
    // "ref" requires access through .value ( and works with primitives )
    // the .value is basically(?) a reactive proxy -- so take your pick.
    const reactiveRoot = reactive(objCatalog.root);
    objCatalog.root = reactiveRoot;

    const q = new Query({
      shuttle: appcfg.shuttle, // url. by default, http://localhost:8080/shuttle/
      objCatalog,
      statusBar: status.value, 
      narration: narration.value,
    }); 
    q.restart("cloak"); // a promise

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
      objTree: reactiveRoot,
      currentTraits,
      onLeafObject(item) {
        currentTraits.value = item.data.traits;
      },
      onEnclosingObject(item) {
        // opening and closing the object "folder" will cause problems with query event handling
        // if (!item.contents) {
        //   item.contents = item.backup;
        //   item.backup = false;
        // } else {
        //   item.backup = item.contents;
        //   item.contents = false;
        // }
        currentTraits.value = item.data.traits;
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
        q.input(text);
      },
    };
  },
};
</script>
