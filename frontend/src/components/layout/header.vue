<!--
 * @author: 马国胜
 * @since: 2023-04-23
 * header.vue
-->
<template>
  <div></div>
</template>

<script  setup  lang='ts'>
import { onBeforeMount, onMounted } from "vue";
onBeforeMount(() => {
  //console.log('2.组件挂载页面之前执行----onBeforeMount')
});
onMounted(() => {
  //console.log('3.-组件挂载到页面之后执行-------onMounted')
});
</script>
<style scoped lang='scss'>
.text {
  color: $primary-color;
}
</style>