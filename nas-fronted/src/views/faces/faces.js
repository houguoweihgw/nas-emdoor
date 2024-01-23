import axios from "axios";

export default {
    name:"faces",
    data(){
        return{
            username: this.$store.state.username,
            faces:[],
        }
    },
    mounted() {
        // 在组件挂载后，从后端获取所有场景分类
        this.fetchFaceClusters()
    },
    methods:{
        async fetchFaceClusters() {
            try {
                const response = await axios.get("/home/faceClusters",
                    {
                        params: {
                            username: this.username,
                        }
                    });
                if (response.status === 200) {
                    this.faces= response.data.faces;
                    console.log(this.faces);
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
        gotoFace(cluster){
            // console.log(label);
            this.$router.push({path: `/home/clusters/${cluster}`})
        },
    }
}