<template>
  <el-row :gutter="2">
    <el-col v-for="(label, index) in labels" :key="index" :span="4">
      <el-card shadow="hover" :body-style="{ padding: '2px'}">

        <div class="image-container" @click="gotoLabel(label.label_name)">
          <el-image
              style="max-width: 100%; max-height: 100%; object-fit: cover;"
              :src="getBase64Image(label.label_cover)"
              :alt="'Photo ' + (index + 1)"
              lazy>
          </el-image>
          <div class="overlay-selected">
              <img :src="require('@/assets/package.png')" alt="package">
          </div>
        </div>
        <div class="title">{{ label.label_name }}</div>
        <div class="extra_info">
          <div>
            <time class="photoTime"> 包含 {{ label.label_count }} 张照片</time>
          </div>
        </div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script src="./scene.js">
</script>

<style scoped>
.el-card{
  background: #e1e1e1;
}

.image-container {
  height: 200px; /* 固定高度 */
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center; /* 使用 align-items 属性来垂直居中图片 */
  align-content: center;
  justify-items: center;
}

.overlay-collected {
  position: absolute;
  left: 0;
  bottom: 0;
}

.overlay-collected img {
  position: absolute;
  left: 0;
  bottom: 0;
}

.overlay-selected {
  position: absolute;
  right: 0;
  bottom: 0;
}

.overlay-selected img {
  position: absolute;
  right: 0;
  bottom: 0;
}

.overlay-more {
  position: absolute;
  right: 0;
  top: 0;
}

.overlay-more img {
  position: absolute;
  right: 0;
  top: 0;
}


.image-overlay button {
  background-color: #337ab7;
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 5px;
  cursor: pointer;
}

.image-container-extra > div {
  margin: 5px;
}

.extra_info {
  padding-top: 10px;
  text-align: center;
  padding-bottom: 5px;
}

.title {
  font-size: 15px;
  text-align: center;
  padding-top: 8px;
  height: 20px;
}

body {
  margin: 10px;
}

.photoTime {
  font-size: 8px;
  color: rgba(0, 0, 0, 0.80);
}


.el-carousel__item h3 {
  color: #475669;
  font-size: 18px;
  opacity: 0.75;
  line-height: 300px;
  margin: 0;
}
.el-empty {
  margin-top: 200px;;
}
.el-carousel__item:nth-child(2n) {
  background-color: #99a9bf;
}

.el-carousel__item:nth-child(2n+1) {
  background-color: #d3dce6;
}
</style>
