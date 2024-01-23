<template>
  <div>
    <el-form ref="form" :model="loginForm" class="loginContainer">
      <h3 class="loginTitle">Ailbum</h3>
      <el-form-item label="账号：">
        <el-input type="text" auto-complete="false" v-model="loginForm.username" placeholder="请输入用户名">
        </el-input>
      </el-form-item>
      <el-form-item label="密码：">
        <el-input type="password" auto-complete="false" v-model="loginForm.password" placeholder="请输入密码">
        </el-input>
      </el-form-item>
      <el-form-item>
        <span class="register-link">没有账号？<router-link to="/register">立即注册</router-link></span>
        <el-button type="primary" style="width: 95%; margin: 10px;" @click.stop.prevent="login">登陆</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script src="./login.js">
</script>

<style>
.loginContainer {
  border-radius: 15px;
  backgroud-clip: padding-box;
  margin: 180px auto;
  width: 350px;
  padding: 15px 35px 15px 35px;
  background: #fff;
  border: 1px solid #eaaaaa;
  box-shadow: 0 0 25px #cac6c6;
}

.loginTitle {
  margin: 0px auto 40px auto;
  text-align: center;
}


</style>
