import axios from "axios";

export default {
    name: "register",
    data() {
        return {
            registerInfo: {
                username: '',
                password: '',
                email: ''
            }
        }
    },
    methods: {
        register() {
            // 1. 获取用户输入的信息
            const username = this.registerInfo.username;
            const password = this.registerInfo.password;
            const email = this.registerInfo.email;

            // 2. 构建注册请求
            const requestData = {
                username: username,
                password: password,
                email: email
            };

            // 发送注册请求
            axios.post('/api/register', requestData)
                .then(response => {
                    // 3. 处理注册请求的响应
                    if (response.status === 200) {
                        this.$message({
                            showClose: true,
                            message: response.data.message,
                            type: 'success',
                            center: true
                        })
                        // 4. 如果注册成功，导航到登录页面
                        this.$router.push('/login');
                    } else {
                        // 注册失败，可以显示错误消息或采取其他操作
                        // console.error('Registration failed:', response.data.message);
                        this.$message({
                            message: response.data.message,
                            type: 'warning',
                            center: true
                        })
                    }
                })
                .catch(error => {
                    console.error('Registration error:', error);
                    this.$message({
                        message: "注册失败",
                        type: 'warning',
                        center: true
                    })
                });
        }
    }
}