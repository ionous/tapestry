<template>
  <div class="lv-status">
    <div class="lv-status__title">{{status.title}}</div>
    <div class="lv-status__room">{{status.location}}</div>
    <div class="lv-status__score" v-if="status.useScoring">{{status.score}}/{{status.turns}}</div>
    <div style="clear: both;"></div>
  </div>
</template>
<script>

import Status from "./status.js";

export default {
 props: {
    status: Status 
  }
}
</script>