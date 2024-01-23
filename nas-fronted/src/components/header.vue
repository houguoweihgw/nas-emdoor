<template>
  <div class="header">
    <!-- 标语 -->
    <div class="slogan">
      <h1>{{ sloganText }}</h1>
    </div>
    <div class="more">
      <el-dropdown @command="handleCommand">
      <span class="el-dropdown-link">
        {{ username }}<i class="el-icon-arrow-down el-icon--right"></i>
      </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="个人信息">个人信息</el-dropdown-item>
          <el-dropdown-item command="登出">登出</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
export default {
  name: "Header",
  data() { // 正确的data选项
    return {
      sloganText: "NAS智慧本地相册", // 传入的标语文本
      username: this.$store.state.username, // 传入的用户名
    };
  },
  methods: {
    handleCommand(command) {
      switch (command) {
        case "个人信息":
          //todo:点击个人信息的处理逻辑

          break;
        case "登出":
          try {
            // console.log("dengchu")
            // 清除前端的用户状态和令牌信息
            this.$store.commit('SET_LOGIN_STATUS', false);
            this.$store.commit('SET_TOKEN', null);
            this.$store.commit('SET_USER_NAME', null);
            this.$store.commit('SET_Album_Array', []);
            // 导航到登录页面或首页
            this.$router.push({name: 'login'}); // 假设登录页面的路由名称是 'login'
          } catch (error) {
            console.error('Logout error:', error);
          }
          break;
      }

    }
  },
};
</script>

<style>
.header {
  display: flex;
  flex-direction: column; /* 垂直排列 */
//justify-content: space-between; /* 让 slogan 居中，el-dropdown 靠右 */ background-color: rgb(128, 128, 128);
  color: #fff;
  padding: 2px;
  width: 220px;
}

.slogan {
  display: flex;
  align-self: center; /* 垂直居中 slogan */
}

.more{
  display: flex;
  align-self: flex-end; /* 让 el-dropdown 靠右 */
  padding: 5px;
}

.account span {
  margin-right: 5px;
  font-weight: bold;
}

.account button {
  background-color: transparent;
  border: none;
  color: #fff;
  cursor: pointer;
}

.el-dropdown .el-dropdown-link {
  cursor: pointer;
  color: #fff;
  font-weight: bold;
  font-size: 15px;
  margin-bottom: 5px; /* 将底部边距设置为15px */
}

.el-icon-arrow-down {
  font-size: 12px;
}
</style>
