<template>
  <div class="lv-container">
    <div v-if="showDebug" 
      class="lv-container__debug">
        <h3>Object List</h3>
      <mk-folder
        :folder="objTree"
        @fileSelected="onLeafObject"
        @folderSelected="onEnclosingObject"
      ></mk-folder>
      <h3>Selected Object</h3>
      <b class="lv-selected-item"> {{ currentItem.name }} ( {{ currentItem.kind }} )</b>:
        <span v-for="(trait,index) in currentItem.traits">{{index?", ":""}}{{trait}}</span>
    </div>
    <div
      v-else-if="narration"
      @click="onContainerClicked">
      <lv-status 
        :status="status"/>
      <lv-output 
        class="lv-story"
        :lines="narration" />
    </div>
    <div class="lv-input"> 
      <lv-prompt v-show="!showDebug"
        enabled="playing"
        @changed="onPrompt"
        ref="prompt" />
      <div v-show="showDebug" class="lv-stub"/>
      <div class="lv-debug">
        <button @click="showDebug = !showDebug">{{ showDebug? "game view": "debug view" }}</button>
      </div>
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
import mkFolder from '/catalog/Folder.vue'
import Query from './query.js'

const objCatalog = new ObjectCatalog();

export default {
  components: { lvOutput, lvPrompt, lvStatus, mkFolder },
  // props: {},
  setup(/*props*/) {
    const narration = ref([]); // ultimately processed through textWriter.js & textRender.js
    const status = ref(new Status());
    const playing = ref(false);
    const showDebug = ref(false);
    const prompt = ref(null); // template slot helper.
    const currentItem = ref({});
    // replace the obj catalog's root with a vue reactive proxy
    // so vue can see those changes when moving between rooms
    const reactiveRoot = reactive(objCatalog.root);
    objCatalog.root = reactiveRoot;
    // make the contents itself reactive ( recursively )
    // so that new objects being added
    // and modifications to those objects can be seen by vue
    objCatalog.all = reactive(objCatalog.all);

    const q = new Query({
      shuttle: appcfg.shuttle, // gets sent to Io constructor
      objCatalog,
      statusBar: status.value, 
      narration: narration.value,
    }); 
    q.restart(tapestry.story); // a promise

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
    return {
      narration, // story output
      status,
      prompt, // template ref
      objTree: reactiveRoot,
      currentItem,
      showDebug,
      onLeafObject(item) {
        currentItem.value = item.data;
      },
      onEnclosingObject(item) {
        // opening and closing the object "folder" causes problems with query event handling
        // if (!item.contents) {
        //   item.contents = item.backup;
        //   item.backup = false;
        // } else {
        //   item.backup = item.contents;
        //   item.contents = false;
        // }
        currentItem.value = item.data;
      },
      // clicking anywhere below the prompt should focus the prompt
      onContainerClicked() {
        const el = prompt.value;
        if (el && el !== document.activeElement) {
          el.setFocus();  
        }
      },
      onPrompt(text) {
        // console.log("onPrompt");
        narration.value.push("> " + text);
        q.fabricate(text);
      },
    };
  },
};
</script>
