<template>
<div class="lv-credits">
  <h1>Credits</h1>
  <p>
    The original "Cloak of Darkness" is a story by <a target="_blank" href="https://www.ifwiki.org/Cloak_of_Darkness">Roger Firth</a>. This version was ported from <a target="_blank" href="https://inform7.com">Inform 7.</a>
  </p>
  <p>
    The images used are licensed through the Creative Commons:
    <ul>
      <li>
        <a target="_blank" href="https://commons.wikimedia.org/w/index.php?curid=13126364
        ">Majestic Theatre Lobby</a> by takomabibelot
      </li>
      <li>
        <a target="_blank" href="https://commons.wikimedia.org/wiki/File:MFS_Lichtenberg_Garderobe.JPG">MFS Lichtenberg Garderobe</a> by Fridolin freudenfett 
      </li>
      <li>
        <a target="_blank" href="https://commons.wikimedia.org/wiki/File:2015_London-Woolwich,_Equitable_pub,_side_room.jpg">Interior of The Equitable pub</a> by Kleon3    
      </li>
      <li>
        <a target="_blank" href="https://commons.wikimedia.org/wiki/File:Pub_interior_-_geograph.org.uk_-_555584.jpg">English Pub interior</a> by Howard Selina
      </li>
      <li>
        <a target="_blank" href="https://commons.wikimedia.org/wiki/File:Hangout_(113454357).jpg">Hangout</a> by hobvias sudoneighm
      </li>
    </ul>
  </p>
  <p>
    This game was built using the <a target="_blank" href="tapestry.ionous.net">Tapestry</a> story engine.
  </p>
  
  <div ref="scrollIntoView">Thanks for playing!</div>
</div>
</template>
<script>
import { ref, onMounted } from "vue";
export default {
  setup() {
    const scrollIntoView = ref(); // template slot helper.
    onMounted(() => {
        scrollIntoView.value.scrollIntoView();
    });
    return {
      scrollIntoView,  
      
    }
  }
}
</script>