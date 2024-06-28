import { createApp } from 'vue'
import Play from './Play.vue'         // contains the router-view

const go = new Go();
WebAssembly.instantiateStreaming(fetch("/tap.wasm.gz"), go.importObject).then((result) => {
    go.run(result.instance);
    let app = createApp(Play);
    app.mount('#play');
});