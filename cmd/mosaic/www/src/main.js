import { createApp } from 'vue'
import * as VueRouter from 'vue-router'
import App from './App.vue'
import NotFound from './NotFound.vue'
import UrlTest from './UrlTest.vue'


const routes = [
  { path: '/', component: UrlTest },
  // will give an array of options to $router.params
  { name: 'edit', path: '/edit/:editPath+', component: UrlTest },
  // will match everything and put it under `$route.params.pathMatch`
  { name: 'NotFound', path: '/:pathMatch(.*)*', component: NotFound },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHashHistory(),
  routes,
});


const app= createApp(App).
  use(router).
  mount('#app');
