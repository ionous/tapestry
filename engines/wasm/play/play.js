import { createApp } from 'vue'
import Play from './Play.vue'         // contains the router-view

import playGob from '/game/play.gob'
import tapWasm from '/game/tap.wasm'

const go = new Go();
fetch(playGob)
  .then(res => res.arrayBuffer())
  .then(playGob => {
    let storyData = new Uint8Array(playGob); // treat buffer as a sequence of 32-bit integers
    WebAssembly.instantiateStreaming(fetch(tapWasm), go.importObject)
      .then((result) => {
        // tbd: i think run() returns a promise that only resolves on main() exit
        go.run(result.instance);
        // "tapestry" is registered by the wasm main() function
        // fix? it could maybe return the name of the main scene.
        tapestry.play(storyData)
          .then(() => {
            // free up the memory
            storyData = null;
            playGob = null;
            // go play the game
            let app = createApp(Play);
            app.mount('#play');
          });
    });
  });
