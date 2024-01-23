<template>
  <div class="footer">
    <p>&copy; 2023 Emdoor. All Rights Reserved.</p>
  </div>
</template>

<script>
export default {
  name: "Footer",
};
</script>

<style>
  .footer {
    //position: fixed; /* 固定定位 */
    bottom: 0; /* 底部距离为0 */
    //width: 100%; /* 页脚宽度占满整个屏幕 */
    text-align: center; /* 文本居中 */
    padding: 5px; /* 可根据需要调整页脚的内边距 */
    background-color: rgb(128, 128, 128);
    color: #fff;
  }
</style>
