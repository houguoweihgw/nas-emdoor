import axios from "axios";
import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";

export default {
    name: "photos",
    data() {
        return {
            showPic: false,
            dialogVisible: false,
            currentDate: new Date(),
            username: this.$store.state.username,
            photosPerRow: 6, // 每行显示的图片数量
            maxRows: 3,// 最大显示的行数
            photos: [],// 空的照片列表，后续会从后端获取并填充
            showPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
            albumsArray: this.$store.state.albumArray,
            currentPage: 1,  // 当前页码
            photosPerPage: 18,  // 每页显示的照片数量
            totalPhotosCount: 0,
        };
    },
    mounted() {
        // 请求来获取总的照片数量
        this.fetchTotalPhotosCount()
        // 在组件挂载后，从后端获取照片数据并填充到photos属性中
        this.fetchPhotos(1)
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
        handleCommand(command) {
            this.addSelectedPhotosToAlbum(command.id)
            console.log("添加到相册", command.id, command.name)
        },
        beforeHandleCommand(id, name) {
            return {
                'id': id,
                'name': name
            }
        },
        toggleSelected(photo){
            photo.selected=!photo.selected
        },
        fetchTotalPhotosCount() {
            axios.get(`/home/totalPhotosCount?username=${this.username}`)
                .then(response => {
                    this.totalPhotosCount = response.data.total;
                })
                .catch(error => {
                    console.error('Failed to fetch total photos count:', error);
                });
        },
        addSelectedPhotosToAlbum(albumId) {
            // 前端请求方法，将所选照片添加到相册
            const selectedPhotos = this.photos.filter(photo => photo.selected);
            // 构造包含所选照片和相册ID的请求数据
            const requestData = {
                albumId: albumId,
                selectedPhotos: selectedPhotos.map(photo => photo.id) // 假设每个照片对象有一个唯一的 ID
            };
            // 发送 POST 请求来将所选照片添加到相册
            axios.post('/home/batchAddPhotosToAlbum', requestData)
                .then(response => {
                    if (response.status === 200) {
                        // 处理成功的响应，可能需要更新前端状态
                        this.$message({
                            type: 'success',
                            message: '添加相册成功'
                        });
                    } else {
                        console.error('Failed to add selected photos to album with status:', response.status);
                    }
                    // 找到所有选中的照片，将它们的 selected 状态设置为 false
                    this.photos.forEach(photo => {
                            if (photo.selected) {
                                photo.selected = false;
                            }
                        }
                    )
                })
                .catch(error => {
                    console.error('Network error:', error);
                });
        },
        deleteSelectedPhotos() {
            // 创建一个数组用来存储选中的照片的 ID
            const selectedPhotoIDs = [];
            // 遍历 photos 数组，检查 selected 属性
            this.photos.forEach(photo => {
                if (photo.selected) {
                    selectedPhotoIDs.push(photo.id);
                }
            });
            if (selectedPhotoIDs.length === 0) {
                // 没有选中的照片，不执行删除操作
                return;
            }
            // 发送请求删除选中的照片
            axios.post('/home/deleteSelectedPhotos',
                {
                    photoIDs: selectedPhotoIDs,
                    username: this.username
                })
                .then(response => {
                    if (response.status === 200) {
                        // 删除成功后，更新前端照片列表
                        selectedPhotoIDs.forEach(id => {
                            this.deletePhotoLocally(id);
                        });
                        this.$message({
                            message: response.data.message,
                            type: 'success',
                            center: true
                        });
                    } else {
                        console.error('Failed to delete selected photos with status:', response.status);
                    }
                })
                .catch(error => {
                    console.error('Delete error:', error);
                });
            this.dialogVisible = false
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
        async fetchPhotos(newPage) {
            try {
                this.currentPage = newPage; // 更新页码
                //todo：将token加入请求头中进行验证
                // 发起请求，传递当前页码和每页数量
                const response = await axios.get(`/home/photos`,
                    {
                        params: {
                            username: this.username,
                            page: this.currentPage,
                            perPage: this.photosPerPage,
                        }
                    }); // 发起 GET 请求
                if (response.status === 200) {
                    // console.log(response.data)
                    // this.photos = response.data.photos.map(photo => {
                    //     return {
                    //         ...photo,
                    //         selected: false,
                    //     };
                    // });
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
        // 处理确认删除的逻辑
        handleDelete(photoID) {
            // 向后端发送删除请求
            axios.delete(`/home/photos/${this.username}/${photoID}`)
                .then(response => {
                    // 请求成功，执行以下操作
                    if (response.status === 200) {
                        // 更新前端的照片列表，删除对应的照片或更新其状态
                        this.deletePhotoLocally(photoID);
                        this.$message({
                            message: response.data.message,
                            type: 'success',
                            center: true
                        });
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
        // displayedPhotos() {
        //     // 计算要显示的图片列表
        //     const maxPhotos = this.photosPerRow * this.maxRows;
        //     return this.photos.slice(0, maxPhotos);
        // },
        selectedPhotosCount() {
            // 使用 Array.prototype.filter 方法来筛选已选中的照片
            const selectedPhotos = this.photos.filter(photo => photo.selected);
            return selectedPhotos.length;
        }
    }
}