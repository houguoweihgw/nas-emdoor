<template>
  <div class="demo-image__preview">
    <div v-if="this.photos.length<=0">
      <el-empty description="你还没有上传照片，试试上传照片吧"></el-empty>
    </div>
    <div class="selected-count" v-if="selectedPhotosCount > 0">已选 {{ selectedPhotosCount }} 张照片
      <el-button type="primary" @click="dialogVisible = true">批量删除选中照片</el-button>
      <el-dialog
          title="提示"
          :visible.sync="dialogVisible"
          width="30%">
        <span>确定删除所选照片到回收站？确定后相册中该照片也会自动删除！！！</span>
        <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="deleteSelectedPhotos">确 定</el-button>
            </span>
      </el-dialog>
      <span>&nbsp&nbsp&nbsp</span>
      <el-dropdown @command="handleCommand">
        <el-button type="primary">
          批量添加照片到相册<i class="el-icon-arrow-down el-icon--right"></i>
        </el-button>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item v-for="(item, index) in albumsArray" :key="item.id"
                            :command="beforeHandleCommand(item.id, item.name)">{{ item.name }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <el-row :gutter="2">
      <el-col v-for="(photo, index) in photos" :key="index" :span="4">
        <el-card shadow="hover" :body-style="{ padding: '2px'}">

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
            <div class="overlay-collected">
              <div @click="toggleLike(photo)">
                <img :src="photo.collected ? require('@/assets/collected.png') : require('@/assets/notCollected.png')"
                     alt="Collected/NotCollected">
              </div>
            </div>
            <div class="overlay-selected">
              <div @click="toggleSelected(photo)">
                <img :src="photo.selected ? require('@/assets/selected.png') : require('@/assets/notSelected.png')"
                     alt="Selected/NotSelected">
              </div>
            </div>
          </div>
          <div class="title">{{ photo.title }}</div>
          <div class="extra_info">
            <div>
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
    <el-pagination
        background
        align="center"
        @current-change="fetchPhotos"
        :current-page="currentPage"
        :page-size="photosPerPage"
        :total="totalPhotosCount"
        class="pagination"
        v-if="this.photos.length>0"
    ></el-pagination>
    <el-backtop>
      <div width="100%"><i class="el-icon-caret-top"></i></div>
    </el-backtop>
  </div>

</template>

<script src="./photos.js">

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
.el-empty {
  margin-top: 200px;;
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
.pagination{
  margin: 20px;
}
</style>
