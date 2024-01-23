import axios from "axios";
import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";

export default {
    name: "collect",
    data() {
        return {
            showPic: false,
            currentDate: new Date(),
            username: this.$store.state.username,
            photosPerRow: 6, // 每行显示的图片数量
            maxRows: 3,// 最大显示的行数
            photos: [], // 空的照片列表，后续会从后端获取并填充
            showPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
        };
    },
    mounted() {
        // 在组件挂载后，从后端获取照片数据并填充到photos属性中
        this.fetchPhotos()
    },
    methods: {
        parsePhotoData(){
            for(let i=0; i<this.photos.length; i++){
                this.photos[i].metadata.date_taken=formatDateTime(this.photos[i].metadata.date_taken)
                this.photos[i].metadata.fileMB=kbToMb(this.photos[i].metadata.file_size)
                this.photos[i].metadata.latitude=getLatitude(this.photos[i].metadata.latitude)
                this.photos[i].metadata.longitude=getLongitude(this.photos[i].metadata.longitude)
                this.photos[i].metadata.altitude=getAltitude(this.photos[i].metadata.altitude)
                // console.log("before",this.photos[i].metadata.file_size)
                // console.log("after",kbToMb(this.photos[i].metadata.file_size) )
            }
        },
        toggleSelected(photo){
            photo.selected=!photo.selected
        },
        async toggleLike(photoId) {
            try {
                // 发送请求来切换收藏状态
                const response = await axios.put(`/home/toggleCollected?username=${this.username}&photo=${photoId}`);
                if (response.status === 200) {
                    // 收藏状态已成功切换
                    const photoIndex = this.photos.findIndex(photo => photo.id === photoId);
                    if (photoIndex !== -1) {
                        this.photos.splice(photoIndex, 1); // 从数组中删除指定的项
                    }
                    this.$message({
                            message: '取消收藏成功',
                            type: 'success',
                            center: true
                        }
                    )
                } else {
                    console.error('Failed to toggle like status with status:', response.status);
                }
            } catch (error) {
                console.error('Network error:', error);
            }
        },
        getBase64Image(fileContent) {
            return `data:image/jpeg;base64,${fileContent}`;
        },
        async fetchPhotos() {
            try {
                this.loading = true;
                //todo：将token加入请求头中进行验证
                // console.log('this.$route.params.albumName:', this.$route.params.albumName); // 输出 albumID 的值
                // ${this.username}&album=${this.$store.state.albumArray[this.albumID-1]}
                const response = await axios.get(`/home/myCollected?username=${this.username}`, {responseType: 'json'}); // 发起 GET 请求
                if (response.status === 200) {
                    this.photos = response.data.photos; // 将获取到的照片数据填充到 this.photos 中
                    this.parsePhotoData()
                    if(response.data.photos!=null){
                        // 初始化 showPhotos 数组
                        this.showPhotos = [];
                        // 获取照片数据
                        const photosData = response.data.photos;
                        // 处理每张照片
                        for (const photoData of photosData) {
                            // 将二进制数据转换为可显示的图片
                            const img = this.getBase64Image(photoData.file_content); // 或者使用其他适合的库或工具来生成图片 URL
                            // 将转换后的图片添加到 showPhotos 数组中
                            this.showPhotos.push(img);
                        }
                        this.$message({
                            showClose: true,
                            message: response.data.message,
                            type: 'success',
                            center: true
                        })
                    }

                } else {
                    console.error('Failed to fetch photos with status:', response.status);
                }
            } catch (error) {
                console.error('Network error:', error);
                // 或者抛出错误以供调用者处理
                throw error;
            }
        },
    },
    computed: {
        displayedPhotos() {
            // 计算要显示的图片列表
            const maxPhotos = this.photosPerRow * this.maxRows;
            return this.photos.slice(0, maxPhotos);
        }
    }
}