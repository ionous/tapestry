<!-- 
  the controller for the overall app.
  contains the "view" which swaps per route/state.    
 -->
<template>
  <div class="mk-header">
    <h1><mk-catalog-button
          @activate="onHome">Tapestry</mk-catalog-button></h1>
    <!-- <mk-toolbar  
      v-if="!!workspace"
      :save-enabled="canSave" 
      :play-enabled="canPlay" 
      @save="onSave" @play="onPlay" /> -->
  </div>
  <div class="mk-footer">
    {{statusText}}
  </div>
  <div class ="mk-view">
    <router-view></router-view>
  </div>
</template>
 
<script>
// a link like component with an @activated event
import mkCatalogButton  from './catalog/CatalogButton.vue'
import appcat from './appcat.js'

// import mkToolbar from './Toolbar.vue'
import { ref } from 'vue'
import http from '/lib/http.js'
import routes from './routes.js'
import endpoint from './endpoints.js';

export default {
  components: { mkCatalogButton , /*mkToolbar*/ },
  // our "constructor": exposes values to our component.
  // https://vuejs.org/api/composition-api-setup.html#composition-api-setup
  setup() {
    const statusText = ref("Welcome to Tapestry...");
    return {
      workspace: null,
      statusText,
    }
  },
  beforeUnmount() {
      this._removeRoutes();
  },
  mounted() {
    const self = this;
    
    this._removeRoutes = routes.router.afterEach((to/*, from*/) => {
      let text = "";
      for (let k in to.params) {
        const v = to.params[k];
        if (!Array.isArray(v)) {
          text = v;
        } else {
          text = v.join("/");
        }
      }
      self.statusText = `Tapestry ${to.name || ''} ${text}`; 
    });
  },
  methods: {
    onHome() {
      routes.goHome();
    },
    // onPlay() {
    //   console.log("play");
    //   // maybe have to send a command to the client so it can run the exe
    //   // could i suppose allow *it* to open the new tab/browser window
    //   // technically *what* to play depends on the project... but not for the moment.
    //   this.canPlay= false;
    //   // 
    //   http.post(endpoint.play, {play:true})
    //     .catch((e)=>{
    //       console.warn("play error", e);
    //     }).finally(()=>{
    //       this.canPlay= true;
    //     });
    // },
    // onWorkspaceChanged(ws) {
    //   this.workspace= ws;
    // }
  }
}
</script>

