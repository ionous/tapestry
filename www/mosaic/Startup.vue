<template>
  <p>loading...</p>
</template>
 
<script>
import shapeData from './shapeData.js';

export default {
  // reports back on shapeData
  // ( the custom data defined by Tapestry for each block )
  emits: ['started'],
  // note: setup has no "this" value
  setup(props, { emit }) {
    // fix: ref some error text or something to the display?
    shapeData.getShapeData().then((b) => {
      emit('started', b);
    }).catch((error) => {
      console.log('error:', error)
    });
  } 
}
</script>


