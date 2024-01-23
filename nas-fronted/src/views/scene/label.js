import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";
import axios from "axios";

export default {
    name: "label",
    data() {
        return {
            showPic: false,
            dialogVisible: false,
            currentDate: new Date(),
            labelName: this.$route.params.labelName,
            username: this.$store.state.username,
            photosPerRow: 6, // 每行显示的图片数量
            maxRows: 3,// 最大显示的行数
            photos: [],// 空的照片列表，后续会从后端获取并填充
            showPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
            albumsArray: this.$store.state.albumArray,
            currentPage: 1,  // 当前页码
            photosPerPage: 18,  // 每页显示的照片数量
            totalLabelCount: 0,
        };
    },
    mounted() {
        // 请求来获取总的照片数量
        this.fetchTotalLabelsCount()
        // 在组件挂载后，从后端获取照片数据并填充到photos属性中
        this.fetchLabelPhotos(1)
    },
    methods: {
        goBack() {
            this.$router.go(-1);
        },
        parsePhotoData(){
            for(let i=0; i<this.photos.length; i++){
                this.photos[i].metadata.date_taken=formatDateTime(this.photos[i].metadata.date_taken)
                this.photos[i].metadata.fileMB=kbToMb(this.photos[i].metadata.file_size)
                this.photos[i].metadata.latitude=getLatitude(this.photos[i].metadata.latitude)
                this.photos[i].metadata.longitude=getLongitude(this.photos[i].metadata.longitude)
                this.photos[i].metadata.altitude=getAltitude(this.photos[i].metadata.altitude)
            }
        },
        fetchTotalLabelsCount() {
            axios.get(`/home/labelPhotoCount?username=${this.username}&label=${this.labelName}`)
                .then(response => {
                    this.totalLabelCount = response.data.total;
                    // console.log(this.totalLabelCount);
                })
                .catch(error => {
                    console.error('Failed to fetch total photos count:', error);
                });
        },
        async toggleLike(photo) {
            try {
                // 发送请求来切换收藏状态
                const photoId = photo.id
                const response = await axios.put(`/home/toggleCollected?username=${this.username}&photo=${photoId}`);
                if (response.status === 200) {
                    // 收藏状态已成功切换
                    // 更新当前照片对象中的collected属性
                    const photoIndex = this.photos.findIndex(photo => photo.id === photoId);
                    if (photoIndex !== -1) {
                        this.photos[photoIndex].collected = !this.photos[photoIndex].collected;
                    }
                    if (photo.collected) {
                        this.$message({
                            message: response.data.message,
                            type: 'success',
                            center: true
                        });
                    }
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
        async fetchLabelPhotos(newPage) {
            try {
                this.currentPage = newPage; // 更新页码
                // console.log(this.labelName)
                //todo：将token加入请求头中进行验证
                // 发起请求，传递当前页码和每页数量
                const response = await axios.get(`/home/labelPhotos`,
                    {
                        params: {
                            username: this.username,
                            page: this.currentPage,
                            perPage: this.photosPerPage,
                            label: this.labelName
                        }
                    }); // 发起 GET 请求
                if (response.status === 200) {
                    // console.log(response.data)
                    this.photos = response.data.photos
                    // console.log(this.photos)
                    this.parsePhotoData()
                    // console.log(response.data.photos)
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
                    // console.log(this.photos)
                    this.$message({
                        showClose: true,
                        message: response.data.message,
                        type: 'success',
                        center: true
                    })
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
}