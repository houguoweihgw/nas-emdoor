<template>
  <div class="app-container">
    <el-container>
      <el-aside width="270px" class="fixed-aside">
        <div class="logo">
          <img :src="require('@/assets/logo.png')" alt="logo"></img>
        </div>
        <el-menu :router='true' class="custom-menu">
          <el-menu-item index="/home/welcome">
            <template slot='title'>
              <span>首页</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/search">
            <template slot='title'>
              <span>智慧搜索</span>
            </template>
          </el-menu-item>

          <el-menu-item index='/home/photos'>
            <template slot='title'>
              <span>全部照片</span>
            </template>
          </el-menu-item>

          <el-submenu
              :default-active="activeMenu"
              index='我的相册'>
            <template slot="title">
              <span>
                {{ submenuTitle }}
              </span>
            </template>
            <el-menu-item
                v-for="menuItem in menuItems"
                :key="menuItem.index"
                :name="menuItem.name"
                :router="true"
                :index="menuItem.index"
                class="manu"
            >
              <router-link :to="`/home/albums/${menuItem.id}`" class="router-link-active">
                <template v-if="!menuItem.editing">
                  <div class="centered-content">{{ menuItem.name }}</div>
                  <div class="right-aligned">
                    <el-dropdown @command="handleCommand">
                      <div>
                        <i class="el-icon-more el-icon--right"></i>
                      </div>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item :command="beforeHandleCommand(menuItem, 'rename')">重命名</el-dropdown-item>
                        <el-dropdown-item :command="beforeHandleCommand(menuItem, 'delete')">删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </el-dropdown>
                  </div>
                </template>
                <template v-else>
                  <el-input ref="inputRef" v-model="menuItem.newName"></el-input>
                  <el-button type="text" @click="saveName(menuItem)">保存</el-button>
                </template>
              </router-link>
            </el-menu-item>
            <div class="centered-button">
              <el-button type="text" @click="addAlbum" class="button-text">
                新建相册
                <i class="el-icon-circle-plus-outline"></i>
              </el-button>
              <el-dialog title="新建相册" :visible.sync="dialogFormVisible">
                <el-form :model="form">
                  <el-form-item label="相册名称" :label-width="formLabelWidth">
                    <el-input v-model="form.name" autocomplete="off"></el-input>
                  </el-form-item>
                  <el-form-item label="相册描述" :label-width="formLabelWidth">
                    <el-input v-model="form.description" autocomplete="off"></el-input>
                  </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                  <el-button @click="dialogFormVisible = false">取 消</el-button>
                  <el-button type="primary" @click="handleCreateAlbum">确 定</el-button>
                </div>
              </el-dialog>
            </div>
          </el-submenu>
          <el-menu-item index="/home/scene">
            <template slot='title'>
              <span>场景分类</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/faces">
            <template slot='title'>
              <span>人物分类</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/collect">
            <template slot='title'>
              <span>我的收藏</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/upload">
            <template slot='title'>
              <span>上传照片</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/recycle">
            <template slot='title'>
              <span>回收站</span>
            </template>
          </el-menu-item>
          <el-menu-item index="/home/about">
            <template slot='title'>
              <span>关于</span>
            </template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container direction="vertical" class="main-content">
        <el-main>
          <router-view :key="$route.fullPath">
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script src="./home.js">

</script>

<style scoped>
/* 自定义导航栏样式 */

.router-link-active {
  /* 添加你的自定义样式 */
  display: flex;
}

.centered-content {
  justify-content: flex-end;
  align-items: flex-end;
}

.right-aligned {
  justify-content: flex-end;
  align-items: flex-end;
}

.icon{
  padding-left: 10px;
}

.dropdown {
  justify-content: flex-end;
}

.logo{
  padding: 10px;
  background-color: rgba(207, 212, 218, 0.49);
}
.el-menu-item, .el-submenu__title {
  height: 60px;
}

.manu {
  display: flex;
  font-size: 13px;
  font-weight: normal;
  text-decoration: none;
  text-align: center;
  background: rgba(236, 238, 241, 0.79);
}

/* 自定义导航栏激活项的背景颜色和文字颜色 */
::v-deep .el-menu-item.is-active {
  background-color: #1890ff;
  color: #fff;
}

::v-deep .el-menu-item {
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 自定义导航栏未激活项的文字颜色 */
::v-deep .el-menu-item:not(.is-active) {
  color: #333;
}

.centered-button {
  display: flex;
  justify-content: center;
  font-weight: normal; /* 取消加粗 */
  align-items: center; /* 垂直居中 */
  background: rgba(236, 238, 241, 0.79);
}

.button-text {
  font-weight: normal; /* 取消加粗 */
  color: #333; /* 修改颜色 */
  font-size: 13px;
  text-align: center;
}

.el-menu {
  background-color: #f0f2f5;
}

/* 自定义侧边栏样式 */
::v-deep .el-aside {
  background-color: #f0f2f5;
  padding: 5px;
  text-align: center;
}

/* 自定义页面主要内容区域样式 */
::v-deep .el-main {
  background-color: #fff;
  padding: 5px;
}

.head-side {
  display: flex;
}

.app-container {
  height: 100vh; /* 100% 视口高度 */
  //overflow: hidden; /* 隐藏超出的内容 */
}
.el-dialog{
  position: relative;
  z-index: 9999; /* 或者任何你需要的较大的数值 */
}
.fixed-aside {
  height: 100%; /* 高度为 100% */
  position: fixed; /* 固定定位 */
  left: 0;
  top: 0;
  bottom: 0;
}

.main-content {
  margin-left: 270px; /* 设置左边距等于侧边栏的宽度 */
  overflow: auto; /* 允许内容滚动 */
  z-index: auto;
}
</style>
