// login.js
import axios from 'axios';

export default {
    name: "Login", data() {
        return {
            loginForm: {
                username: '', password: '',
            }
        }
    }, methods: {
        async login() {
            try {
                const username = String(this.loginForm.username).trim();
                const password = String(this.loginForm.password).trim();

                if (username === "" || password === "") {
                    alert("请输入有效的用户名和密码。");
                    return;
                }
                const response = await axios.post('/api/login', {username, password}, {responseType: 'json'});
                // 检查响应状态
                if (response.status === 200) {
                    // 保存相关用户信息
                    this.$store.commit('SET_USER_NAME', username)
                    this.$store.commit('SET_LOGIN_STATUS', true)
                    this.$store.commit('SET_TOKEN', response.data.token)
                    // 通知
                    this.$message({
                        showClose: true,
                        message: response.data.message,
                        type: 'success',
                        center: true
                    })
                    // 跳转到
                    await this.$router.push({name: 'welcome'})
                } else {
                    // 处理其他响应状态码
                    console.log(response)
                    this.$message({
                        message: response.data.message,
                        type: 'warning',
                        center: true
                    });
                    // console.error('Login failed with status:', response.status);
                }
            } catch (error) {
                // 处理网络请求错误
                console.error('Network error:', error);
                // 或者抛出错误以供调用者处理
                throw error;
            }
        },
    }
}