import axios from "axios";
export default {
    name: "scene",
    data() {
        return{
            username: this.$store.state.username,
            labels: [],
        }
    },
    mounted() {
        // 在组件挂载后，从后端获取所有场景分类
        this.fetchSceneLabels()
    },
    methods: {
        async fetchSceneLabels() {
            try {
                const response = await axios.get("/home/sceneLabels",
                    {
                        params: {
                            username: this.username,
                        }
                    });
                if (response.status === 200) {
                    this.labels= response.data.labels;
                    // console.log(this.labels);
                    this.$message({
                        showClose: true,
                        message: response.data.message,
                        type: 'success',
                        center: true
                    })
                }
            } catch (error) {
                console.error('Network error:', error);
                // 或者抛出错误以供调用者处理
                throw error;
            }
        },
        getBase64Image(fileContent) {
            return `data:image/jpeg;base64,${fileContent}`;
        },
        gotoLabel(label){
            // console.log(label);
            this.$router.push({path: `/home/labels/${label}`})
        }
    }
}