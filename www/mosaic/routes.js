// views:
import Edit from './Edit.vue'
import HomePage from './HomePage.vue'
import UserAction from './UserAction.vue'
//
import * as VueRouter from 'vue-router'

const routes = [
  { name: 'home', path: '/', component: HomePage },
  // tbd: this empties the entire view; maybe it could be more of a modal thing somehow?
  { name: 'file', path: '/file/:action', component: UserAction, meta: { transient:true} },
  // the + assumes that the editPath parameter will be an array of one or more elements to handle slashes
  { name: 'edit', path: '/edit/:editPath+', component: Edit },
  // will match everything else and put it under `$route.params.pathMatch`
  //{ name: 'NotFound', path: '/:pathMatch(.*)*', component: NotFound },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHashHistory(),
  routes,
});

export default {
  router,
  // or null if no on the useractions page
  getAction() {
    return router.currentRoute.value.params.action;
  },
  getEditPath() {
    const path= router.currentRoute.value.params.editPath;
    return path && path.join("/");
  },
  goHome() {
    this.go('home');
  },
  newFile() {
    this.go('file', { action: 'new' });
  },
  openFile(path) {
    if (!path) {
      this.go('file', { action: 'open' });
    } else {
      // a little convoluted: the router changes blockly
      // ( via Blockly.vue onRouteChanged() )
      // blockly tells us when the workspace has changed
      // ( via the on-workspace-changed event )
      const parts= path.split("/");
      this.go('edit', { editPath: parts });
    }
  },
  go(name, params) {
    const transient = router.currentRoute.value.meta.transient;
    const method =  transient? router.replace: router.push;
    method({name, params});
  }
}
