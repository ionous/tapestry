import { createApp } from 'vue'
import Play from './Play.vue'         // contains the router-view

const go = new Go();
let app; 
WebAssembly.instantiateStreaming(fetch("/tap.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    app = createApp(Play);
    app.mount('#play');
});
