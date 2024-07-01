<template>
  <div class="lv-container">
    <div v-if="showDebug" 
      class="lv-debug">
        <h3>Object List</h3>
        <mk-tree :item="currRoom" :list="true"
          @activated="onActivated">
        </mk-tree>
        <div v-if="debugItem">
          <b>{{ debugItem.name }}</b> <span>( {{ debugItem.kind }} ) </span>
          <div class="lv-traits">
            <span v-for="(trait,index) in debugItem.traits">{{index?", ":""}}{{trait}}</span>
          </div>
        </div>
    </div>
    <div class="lv-console"
      v-else-if="narration"
      @click="onContainerClicked">
      <lv-status v-if="!showCredits"
        :status="status"/>
      <div class="lv-output">
        <lv-output 
          :lines="narration" />
        <div v-if="!playing"> 
          <a v-if="!showCredits" href="#" @click="showCredits= !showCredits">Show Credits</a>
          <lv-credits v-else></lv-credits>
        </div>
      </div>
    </div>
    <div class="lv-input"> 
      <lv-prompt v-show="!showDebug && playing"
        @changed="onPrompt"
        ref="prompt" />
      <div v-show="showDebug" class="lv-stub"/>
      <button @click="showDebug = !showDebug">{{ showDebug? "game view": "debug view" }}</button>
    </div>
    <div class="lv-image">
      <mk-tree :item="currRoom">
      </mk-tree>
  </div>
  </div>
  
</template>
<script>
import lvPrompt from "./Prompt.vue";
import lvOutput from "./Output.vue";
import lvStatus from "./Status.vue";
import lvCredits from "./Credits.vue";
import Status from "./status.js";

import { ref, onMounted, onUnmounted, reactive } from "vue";
//
import mkTree from './Tree.vue'
import Query from './query.js'

export default {
  components: { lvOutput, lvPrompt, lvStatus, lvCredits, mkTree },
  // props: {},
  setup(/*props*/) {
    const narration = ref([]); // see textWriter.js & textRender.js
    const status = ref(new Status());
    const playing = ref(true);
    const showDebug = ref(false);
    const prompt = ref(null); // template slot helper.
    const allItems = ref({}); // id and states
    const currRoom = ref({}); // id and states
    const debugItem = ref(false); // id and states
    const showCredits = ref(false);
    
    const q = new Query({
      shuttle: appcfg.shuttle, // gets sent to Io constructor
      narration: narration.value,
      currRoom,
      allItems: allItems.value,
      statusBar: status.value, 
      playing,
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
      currRoom,
      prompt, // template ref
      showDebug,
      showCredits,
      debugItem, // selected item
      playing,
      onActivated(id) {
        console.log("activated", id);
        debugItem.value = allItems.value[id];
      },
      onContainerClicked() {
        const el = prompt.value;
        if (el && el !== document.activeElement) {
          el.setFocus();  
        }
      },
      onPrompt(text) {
        // console.log("onPrompt");
        narration.value.push("> " + text);
        q.fabricate(text).catch(e => {
          narration.value.push("! " + e.message );
          console.error(e);
        });
      },
    };
  },
};
</script>
