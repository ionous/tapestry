<script setup>
// https://v3.vuejs.org/api/sfc-script-setup.html#basic-syntax
// Any top-level bindings (including variables, function declarations, and imports) declared inside here are usable by the template.
import MockCatalog from './catalog/mockCatalog.js'
import Catalog from './catalog/Catalog.vue'

const catalog = new MockCatalog();
</script>

<template>
  <img alt="Vue logo" src="./assets/logo.png" />
  <Catalog :catalog="catalog" />
</template>

<style>
/*#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}*/
</style>
