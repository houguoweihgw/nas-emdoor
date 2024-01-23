<template>
  <el-upload
      class="upload-demo"
      ref="upload"
      action="http://localhost:8001/home/upload"
      multiple="multiple"
      :on-preview="handlePreview"
      :on-remove="handleRemove"
      :file-list="fileList"
      :data = "uploadData"
      :auto-upload="false">
    <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
    <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上传到服务器</el-button>
    <div slot="tip" class="el-upload__tip">只能上传jpg/png文件</div>
  </el-upload>
</template>

<script>

export default {
  name:"upload",
  data() {
    return {
      fileList: [],
      uploadData: {  // 在 data 中定义要传递的数据
        username: this.$store.state.username
      }
    };
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleRemove(file, fileList) {
      console.log(file, fileList);
    },
    handlePreview(file) {
      console.log(file);
    }
  }
}
</script>

<style scoped>
  .upload-demo{
    padding-top: 20px;
    text-align: center;
  }
</style>