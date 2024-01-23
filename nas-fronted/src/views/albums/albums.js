import axios from "axios";
import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";

export default {
    name: "albums",
    data() {
        return {
            showPic: false,
            dialogVisible: false,
            albumID: this.$route.params.albumName,
            currentDate: new Date(),
            username: this.$store.state.username,
            photosPerRow: 6, // 每行显示的图片数量
            maxRows: 3,// 最大显示的行数
            photos: [], // 空的照片列表，后续会从后端获取并填充
            showPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
        };
    },
    created() {
        // 在组件挂载后，从后端获取照片数据并填充到photos属性中
        if (this.$store.state.albumArray[this.albumID - 1].id !== undefined) {
            this.fetchPhotos()
        }
    },
    methods: {
        parsePhotoData() {
            for (let i = 0; i < this.photos.length; i++) {
                this.photos[i].metadata.date_taken = formatDateTime(this.photos[i].metadata.date_taken)
                this.photos[i].metadata.fileMB = kbToMb(this.photos[i].metadata.file_size)
                this.photos[i].metadata.latitude = getLatitude(this.photos[i].metadata.latitude)
                this.photos[i].metadata.longitude = getLongitude(this.photos[i].metadata.longitude)
                this.photos[i].metadata.altitude = getAltitude(this.photos[i].metadata.altitude)
                // console.log("before",this.photos[i].metadata.file_size)
                // console.log("after",kbToMb(this.photos[i].metadata.file_size) )
            }
        },
        toggleSelected(photo) {
            photo.selected = !photo.selected
        },
        removeFromAlbum() {
            // 前端请求方法，将所选照片添加到相册
            const selectedPhotos = this.photos.filter(photo => photo.selected);
            // 构造包含所选照片和相册ID的请求数据
            const requestData = {
                albumId: this.$store.state.albumArray[this.albumID - 1].id,
                selectedPhotos: selectedPhotos.map(photo => photo.id) // 假设每个照片对象有一个唯一的 ID
            };
            // 发送 POST 请求来将所选照片添加到相册
            axios.post('/home/batchRemoveFromAlbum', requestData)
                .then(response => {
                    if (response.status === 200) {
                        // 移除已选中的照片
                        const selectedPhotoIDs = requestData.selectedPhotos;
                        this.photos = this.photos.filter(photo => !selectedPhotoIDs.includes(photo.id));
                        // 处理成功的响应，可能需要更新前端状态
                        this.$message({
                            type: 'success',
                            message: '移出相册成功'
                        });
                    } else {
                        console.error('Failed to remove selected photos from album with status:', response.status);
                    }
                })
                .catch(error => {
                    console.error('Network error:', error);
                });
            this.dialogVisible = false
        },
        getBase64Image(fileContent) {
            return `data:image/jpeg;base64,${fileContent}`;
        },
        async fetchPhotos() {
            try {
                //todo：将token加入请求头中进行验证
                // console.log('this.$route.params.albumName:', this.$route.params.albumName); // 输出 albumID 的值
                // ${this.username}&album=${this.$store.state.albumArray[this.albumID-1]}
                const response = await axios.get(`/home/albumPhotos?username=${this.username}&album=${this.$store.state.albumArray[this.albumID - 1].id}`, {responseType: 'json'}); // 发起 GET 请求
                if (response.status === 200) {
                    // console.log(response.data)
                    this.photos = response.data.photos; // 将获取到的照片数据填充到 this.photos 中
                    this.parsePhotoData()
                    console.log(this.photos)
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
                } else {
                    console.error('Failed to fetch photos with status:', response.status);
                }
            } catch (error) {
                console.error('Network error:', error);
                // 或者抛出错误以供调用者处理
                throw error;
            }
        },
        // 处理确认删除的逻辑
        handleDelete(photoID) {
            // 向后端发送删除请求
            axios.delete(`/home/album/${this.$store.state.albumArray[this.albumID - 1].id}/${photoID}`)
                .then(response => {
                    // 请求成功，执行以下操作
                    if (response.status === 200) {
                        // 更新前端的照片列表，删除对应的照片或更新其状态
                        this.deletePhotoLocally(photoID);
                        this.$message({
                                message: response.data.message,
                                type: 'success',
                                center: true
                            }
                        )
                    } else {
                        console.error('Failed to delete photo with status:', response.status);
                    }
                })
                .catch(error => {
                    console.error('Delete error:', error);
                });
        },

        deletePhotoLocally(photoID) {
            const index = this.photos.findIndex(photo => photo.id === photoID);
            if (index !== -1) {
                this.photos.splice(index, 1);
            }
        },

        // 通过按钮点击来显示确认对话框
        showConfirmDialog() {
            this.$refs.confirmDialog.showPopper = true;
        },
    },
    computed: {
        displayedPhotos() {
            // 计算要显示的图片列表
            const maxPhotos = this.photosPerRow * this.maxRows;
            return this.photos.slice(0, maxPhotos);
        },
        selectedPhotosCount() {
            // 使用 Array.prototype.filter 方法来筛选已选中的照片
            const selectedPhotos = this.photos.filter(photo => photo.selected);
            return selectedPhotos.length;
        }
    }
}