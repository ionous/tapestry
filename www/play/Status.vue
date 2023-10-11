<template>
  <div class="lv-status">
    <div class="lv-status__title">{{title}}</div>
    <div class="lv-status__room">{{location}}</div>
    <div class="lv-status__score" v-if="useScoring">{{score}}/{{turns}}</div>
    <div style="clear: both;"></div>
  </div>
</template>
<script>
	export default {
   props: {
      title: String,
      location: String,
      useScoring: Boolean,
      score: Number,
      turns: Number, 
    }
  }
</script>