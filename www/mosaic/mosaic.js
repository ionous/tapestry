import { createApp } from 'vue'
import Mosaic from './Mosaic.vue'
import routes from './routes.js';

const app= createApp(Mosaic).
  use(routes.router).
  mount('#mosaic');
