<template>
  <div class="demo-image__preview">
    <div v-if="this.recycledPhotos.length<=0">
      <el-empty description="回收站还是空的哦"></el-empty>
    </div>
    <div class="selected-count" v-if="selectedPhotosCount > 0">已选 {{ selectedPhotosCount }} 张照片
      <el-button type="primary" @click="recoverPhotos">批量恢复</el-button>
      <el-button type="primary" @click="dialogVisible = true">批量清除照片</el-button>
      <el-dialog
          title="提示"
          :visible.sync="dialogVisible"
          width="30%">
        <span>确定从回收站彻底删除所选照片？</span>
        <span>彻底删除后将无法恢复!!!</span>
        <span slot="footer" class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="deletePhotos">确 定</el-button>
        </span>
      </el-dialog>
    </div>
      <el-row :gutter="2">
        <el-col v-for="(photo, index) in displayedRecycledPhotos" :key="index" :span="4">
          <el-card shadow="hover" :body-style="{ padding: '2px' }">
            <div class="image-container">
              <el-image
                  style="max-width: 100%; max-height: 100%; object-fit: cover;"
                  :src="getBase64Image(photo.file_content)"
                  :alt="'Photo ' + (index + 1)"
                  :preview-src-list="showPhotos"
                  lazy>
              </el-image>
              <div class="overlay-more" v-popover="`node-popover-${photo.id}`">
                <img :src="require('@/assets/more.png')" alt="更多"></img>
              </div>
              <div class="overlay-selected">
                <div @click="toggleSelected(photo)">
                  <img :src="photo.selected ? require('@/assets/selected.png') : require('@/assets/notSelected.png')"
                       alt="Selected/NotSelected">
                </div>
              </div>
            </div>

            <div class="extra_info">
              <div class="title">{{ photo.title }}</div>
              <div class="timeInfo">
                <i class="el-icon-date"></i>
                <time class="photoTime"> {{ photo.metadata.date_taken }}</time>
              </div>
              <div>
                <i class="el-icon-camera-solid"> </i>
                <time class="photoTime"> {{ photo.metadata.model }}
                  {{ photo.metadata.image_width }}&times{{ photo.metadata.image_length }} {{ photo.metadata.fileMB }}MB
                </time>
              </div>
              <div>
                <i class="el-icon-orange"> </i>
                <time class="photoTime"> {{ photo.metadata.focal_length }}mm
                  {{ photo.metadata.exposure_time }} f/{{ photo.metadata.aperture }} ISO{{ photo.metadata.iso }}
                </time>
              </div>
              <div>
                <i class="el-icon-location"> </i>
                <time class="photoTime"> {{ photo.metadata.latitude }} {{ photo.metadata.longitude }}
                  {{ photo.metadata.altitude }}
                </time>
              </div>
            </div>
          </el-card>
          <div>
            <el-popover
                placement="top"
                width="250"
                :ref="`node-popover-${photo.id}`"
                trigger="hover">
              <el-descriptions title="详细信息" column="1">
                <el-descriptions-item label="照片大小">{{ photo.metadata.fileMB }}MB
                  ({{ photo.metadata.file_size }}字节)
                </el-descriptions-item>
                <el-descriptions-item label="照片宽度">{{ photo.metadata.image_width }}像素</el-descriptions-item>
                <el-descriptions-item label="照片高度">{{ photo.metadata.image_length }}像素</el-descriptions-item>
                <el-descriptions-item label="相机品牌">{{ photo.metadata.make }}</el-descriptions-item>
                <el-descriptions-item label="相机型号">{{ photo.metadata.model }}</el-descriptions-item>
                <el-descriptions-item label="曝光时间">{{ photo.metadata.exposure_time }}</el-descriptions-item>
                <el-descriptions-item label="光圈值">f/{{ photo.metadata.aperture }}</el-descriptions-item>
                <el-descriptions-item label="ISO速度等级">{{ photo.metadata.iso }}</el-descriptions-item>
                <el-descriptions-item label="焦距">{{ photo.metadata.focal_length }}mm</el-descriptions-item>
                <el-descriptions-item label="拍摄时间">{{ photo.metadata.date_taken }}</el-descriptions-item>
                <el-descriptions-item label="地理坐标">{{ photo.metadata.latitude }}/{{ photo.metadata.longitude }}
                  ({{ photo.metadata.altitude }})
                </el-descriptions-item>
                <el-descriptions-item label="标签">{{ photo.metadata.scene_tags }}</el-descriptions-item>
              </el-descriptions>
            </el-popover>
          </div>
        </el-col>
      </el-row>
      <el-backtop>
        <div width="100%"><i class="el-icon-caret-top"></i></div>
      </el-backtop>
  </div>
</template>

<script>
import * as assert from "assert";
import axios from "axios";
import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";

export default {
  name: "recycledPhotos",
  data() {
    return {
      dialogVisible: false,
      username: this.$store.state.username,
      photosPerRow: 6, // 每行显示的图片数量
      maxRows: 3,// 最大显示的行数
      recycledPhotos: [], // 空的照片列表，后续会从后端获取并填充
      showPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
    };
  },
  mounted() {
    // 在组件挂载后，从后端获取照片数据并填充到recycledPhotos属性中
    this.fetchRecycledPhotos()
  },
  methods: {
    parsePhotoData(){
      for(let i=0; i<this.recycledPhotos.length; i++){
        this.recycledPhotos[i].metadata.date_taken=formatDateTime(this.recycledPhotos[i].metadata.date_taken)
        this.recycledPhotos[i].metadata.fileMB=kbToMb(this.recycledPhotos[i].metadata.file_size)
        this.recycledPhotos[i].metadata.latitude=getLatitude(this.recycledPhotos[i].metadata.latitude)
        this.recycledPhotos[i].metadata.longitude=getLongitude(this.recycledPhotos[i].metadata.longitude)
        this.recycledPhotos[i].metadata.altitude=getAltitude(this.recycledPhotos[i].metadata.altitude)
        // console.log("before",this.photos[i].metadata.file_size)
        // console.log("after",kbToMb(this.photos[i].metadata.file_size) )
      }
    },
    toggleSelected(photo){
      photo.selected=!photo.selected
    },
    recoverPhotos() {
      // 前端请求方法，将所选照片添加到相册
      const selectedPhotos = this.recycledPhotos.filter(photo => photo.selected);
      // 构造包含所选照片和相册ID的请求数据
      const requestData = {
        // albumId: this.$store.state.albumArray[this.albumID-1].id,
        selectedPhotos: selectedPhotos.map(photo => photo.id), // 假设每个照片对象有一个唯一的 ID
        username: this.username
      };
      // 发送 POST 请求来批量恢复照片
      axios.post('/home/recoverBatchPhotos', requestData)
          .then(response => {
            if (response.status === 200) {
              // 移除已选中的照片
              const selectedPhotoIDs = requestData.selectedPhotos;
              this.recycledPhotos = this.recycledPhotos.filter(photo => !selectedPhotoIDs.includes(photo.id));
              // 处理成功的响应，可能需要更新前端状态
              this.$message({
                type: 'success',
                message: '恢复成功'
              });
            } else {
              console.error('Failed to recover selected photos from recycle with status:', response.status);
            }
          })
          .catch(error => {
            console.error('Network error:', error);
          });
    },
    deletePhotos() {
      // 前端请求方法，将所选照片添加到相册
      const selectedPhotos = this.recycledPhotos.filter(photo => photo.selected);
      // 构造包含所选照片和相册ID的请求数据
      const requestData = {
        // albumId: this.$store.state.albumArray[this.albumID-1].id,
        selectedPhotos: selectedPhotos.map(photo => photo.id), // 假设每个照片对象有一个唯一的 ID
        username: this.username
      };
      // 发送 POST 请求来将所选照片添加到相册
      axios.post('/home/batchDeletePhotos', requestData)
          .then(response => {
            if (response.status === 200) {
              // 移除已选中的照片
              const selectedPhotoIDs = requestData.selectedPhotos;
              this.recycledPhotos = this.recycledPhotos.filter(photo => !selectedPhotoIDs.includes(photo.id));
              // 处理成功的响应，可能需要更新前端状态
              this.$message({
                type: 'success',
                message: '彻底删除成功'
              });
            } else {
              console.error('Failed to recover selected photos from recycle with status:', response.status);
            }
          })
          .catch(error => {
            console.error('Network error:', error);
          });
      this.dialogVisible = false
    },
    getBase64Image(fileContent) {
      return `data:image/jpeg;base64,${fileContent}`;
    },
    async fetchRecycledPhotos() {
      try {
        //todo：将token加入请求头中进行验证
        const response = await axios.get(`/home/recycledPhotos?username=${this.username}`, {responseType: 'json'}); // 发起 GET 请求
        if (response.status === 200) {
          // console.log(response.data.recycledPhotos)
          this.recycledPhotos = response.data.recycledPhotos; // 将获取到的照片数据填充到 this.recycledPhotos 中
          this.parsePhotoData()
          // 初始化 showPhotos 数组
          this.showPhotos = [];
          // 获取照片数据
          const photosData = response.data.recycledPhotos;
          // 处理每张照片
          for (const photoData of photosData) {
            // 将二进制数据转换为可显示的图片
            const img = this.getBase64Image(photoData.file_content); // 或者使用其他适合的库或工具来生成图片 URL
            // 将转换后的图片添加到 showPhotos 数组中
            this.showPhotos.push(img);
          }
          this.$message({
            showClose: true,
            message: response.data.message,
            type: 'success',
            center: true
          })
        } else {
          console.error('Failed to fetch recycledPhotos with status:', response.status);
        }
      } catch (error) {
        console.error('Network error:', error);
        // 或者抛出错误以供调用者处理
        throw error;
      }
    },
    // 处理确认删除的逻辑
    handleDelete(photoID) {
      // 向后端发送删除请求
      axios.delete(`/home/recycledPhotos/${photoID}`)
          .then(response => {
            // 请求成功，执行以下操作
            if (response.status === 200) {
              // 更新前端的照片列表，删除对应的照片或更新其状态
              this.deletePhotoLocally(photoID);
              this.$message({
                message: response.data.message,
                type: 'success',
                center: true
              });
            } else {
              console.error('Failed to delete photo with status:', response.status);
            }
          })
          .catch(error => {
            console.error('Delete error:', error);
          });
    },

    deletePhotoLocally(photoID) {
      const index = this.recycledPhotos.findIndex(photo => photo.id === photoID);
      if (index !== -1) {
        this.recycledPhotos.splice(index, 1);
      }
    },
    // 通过按钮点击来显示确认对话框
    showConfirmDialog() {
      this.$refs.confirmDialog.showPopper = true;
    },
  },

  computed: {
    displayedRecycledPhotos() {
      // 计算要显示的图片列表
      const maxRecycledPhotos = this.photosPerRow * this.maxRows;
      return this.recycledPhotos.slice(0, maxRecycledPhotos);
    },
    selectedPhotosCount() {
      // 使用 Array.prototype.filter 方法来筛选已选中的照片
      const selectedPhotos = this.recycledPhotos.filter(photo => photo.selected);
      return selectedPhotos.length;
    }
  }
}

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
.el-empty {
  margin-top: 200px;;
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

.image-container-extra {
  display: flex;
  align-items: center; /* 垂直居中对齐 */
  justify-content: center; /* 水平居中对齐 */
  height: 25px;
}

.image-container-extra > div {
  margin: 5px;
}

.extra_info {
  padding-left: 25px;
  text-align: start;
}

.title {
  font-size: 13px;
  text-align: center;
  padding-top: 5px;
  height: 17px;
}

body {
  margin: 10px;
}

.photoTime {
  font-size: 8px;
  color: rgba(0, 0, 0, 0.80);
}


.selected-count {
  margin: 5px;
  padding: 5px;
}

.pagination {
  margin-top: 20px;
}

.el-carousel__item h3 {
  color: #475669;
  font-size: 18px;
  opacity: 0.75;
  line-height: 300px;
  margin: 0;
}

.el-carousel__item:nth-child(2n) {
  background-color: #99a9bf;
}

.el-carousel__item:nth-child(2n+1) {
  background-color: #d3dce6;
}
</style>
