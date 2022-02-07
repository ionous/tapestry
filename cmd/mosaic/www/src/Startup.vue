 <template>
  <p>loading...</p>
</template>
 
<script>
import Blockly from 'blockly';
import { inject } from 'vue'

// shapes is an array of blockly shape data.
const shapesUrl= 'http://localhost:8080/shapes/';

export default {
  // reports back on shapeData
  // ( the custom data defined by Tapestry for each block )
  emits: ['started'],
  // note: setup has no "this" value
  setup(props, { emit /*attrs, slots, emit, expose*/ }) {
    // fix: ref some error text or something to the display?
    // fix: push your json fetch into a shared library?
    fetch(shapesUrl).then((response) => {
      return (!response.ok) ?
        Promise.reject({status: response.status, url}) :
        response.json().then((jsonData) => {
          let shapeData= {};
          jsonData.forEach(function(el) {
            Blockly.defineBlocksWithJsonArray([el]);
            let { customData } = el;
            if (customData) { // ironically: mutators use standard blocks.
              shapeData[el.type]= customData;
            }
          });
          emit('started', shapeData);
        });
    }).catch((error) => {
      console.log('error:', error)
    });
  } 
}
</script>