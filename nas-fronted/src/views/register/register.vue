<template>
  <div>
    <el-form ref="registerForm" :model="registerInfo" class="register-container">
      <h3 class="register-title">注册</h3>
      <el-form-item label="用户名：" prop="username">
        <el-input type="text" v-model="registerInfo.username" auto-complete="false" placeholder="请输入用户名"></el-input>
      </el-form-item>
      <el-form-item label="密码：" prop="password">
        <el-input type="password" v-model="registerInfo.password" auto-complete="false" placeholder="请输入密码"></el-input>
      </el-form-item>
      <el-form-item label="邮箱：" prop="email">
        <el-input type="text" v-model="registerInfo.email" auto-complete="false" placeholder="请输入邮箱"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" style="width: 100%" @click="register">注册</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script src="./register.js">

</script>


<style scoped>
.register-container {
  border-radius: 15px;
  backgroud-clip: padding-box;
  margin: 180px auto;
  width: 350px;
  padding: 15px 35px 15px 35px;
  background: #fff;
  border: 1px solid #eaaaaa;
  box-shadow: 0 0 25px #cac6c6;
}

.register-title {
  margin: 0px auto 40px auto;
  text-align: center;
}

</style>