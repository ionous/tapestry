import { createApp } from 'vue'
import * as VueRouter from 'vue-router'
import Mosaic from './Mosaic.vue'

const routes = [
  { name: 'edit', path: '/edit/:editPath+', component: Mosaic },
  // will match everything else and put it under `$route.params.pathMatch`
  //{ name: 'NotFound', path: '/:pathMatch(.*)*', component: NotFound },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHashHistory(),
  routes,
});

const app= createApp(Mosaic).
  use(router).
  mount('#mosaic');
