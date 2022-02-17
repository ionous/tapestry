import { createApp } from 'vue'
import Live from './Live.vue'         // contains the router-view

const app= createApp(Live).
  mount('#live');
