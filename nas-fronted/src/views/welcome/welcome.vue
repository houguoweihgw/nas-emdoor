<template>
  <div class="welcome">
    <h1>欢迎 {{this.$store.state.username}} 来到智慧本地相册</h1>
    <p>这里是您的私人相册，您可以管理和分享照片</p>

    <div class="carousel">
      <el-carousel :interval="4000" type="card" height="450px">
        <el-carousel-item v-for="(photo, index) in welcomePhotos" :key="index">
          <div class="photo-container">
            <img :src="getBase64Image(photo.file_content)" :alt="'Image ' + (index + 1)" class="centered-image" />
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "welcome",
  data(){
    return{
      username:  this.$store.state.username,
      welcomePhotos:[]
    }
  },
  mounted() {
    this.getWelcomePhotos()
  },
  methods: {
    getBase64Image(fileContent) {
      return `data:image/jpeg;base64,${fileContent}`;
    },
    async getWelcomePhotos(){
        try {
          // 发起请求，传递当前页码和每页数量
          const response = await axios.get(`/home/welcomePhotos`,
              {
                params: {
                  username: this.username,
                }
              }); // 发起 GET 请求
          if (response.status === 200) {
            this.welcomePhotos = response.data.photos
            this.parsePhotoData()
            this.$message({
              showClose: true,
              message: response.data.message,
              type: 'success',
              center: true
            });
          } else {
            console.error('Failed to fetch photos with status:', response.status);
          }
        } catch (error) {
          console.error('Network error:', error);
          // 或者抛出错误以供调用者处理
          throw error;
        }
    }
  }
}
</script>

<style scoped>
.welcome {
  text-align: center;
  margin-top: 50px;
}

h1 {
  font-size: 32px;
  color: #333;
  margin: 50px;
}

p {
  font-size: 18px;
  color: #666;
  margin: 30px;
}
.centered-image {
  max-width: 100%; /* 图片最大宽度为容器宽度 */
  //max-height: 100%; /* 图片最大高度为容器高度 */
  display: block; /* 让图片水平居中 */
  margin: 0 auto; /* 让图片垂直居中 */
}

.photo-container {
  width: 100%; /* 容器宽度为 Carousel 宽度 */
  height: 100%; /* 容器高度为 Carousel 高度 */
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
}

.carousel{
  padding: 50px;
}

.el-carousel__item h3 {
  color: #475669;
  font-size: 14px;
  opacity: 0.75;
  line-height: 200px;
  margin: 0;
}

.el-carousel__item:nth-child(2n) {
  background-color: #99a9bf;
}

.el-carousel__item:nth-child(2n+1) {
  background-color: #d3dce6;
}
</style>
