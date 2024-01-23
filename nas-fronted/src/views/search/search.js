import {formatDateTime, getAltitude, getLatitude, getLongitude, kbToMb} from "@/utils/dateUtils";
import axios from "axios";

export default {
    name: "search",
    data() {
        return {
            showPic: false,
            dialogVisible: false,
            labelName: this.$route.params.labelName,
            username: this.$store.state.username,
            photosPerRow: 6, // 每行显示的图片数量
            maxRows: 3,// 最大显示的行数
            searchPhotos: [],// 空的照片列表，后续会从后端获取并填充
            showSearchPhotos: [],// 空的大图展示照片列表，后续会从后端获取并填充
            currentPage: 1,  // 当前页码
            photosPerPage: 18,  // 每页显示的照片数量
            totalLabelCount: 0,
            searchQuery: "",
            lastSearchQuery: "",
        };
    },
    props: {
        placeholderText: {
            type: String,
            default: "搜索照片",
        },
        searchButtonText: {
            type: String,
            default: "搜索",
        },
    },
    mounted() {
        // 在input输入框被渲染完毕后再获取焦点
        this.$nextTick(() => {
            // 使用引用的原生DOM对象input的focus方法自动获得焦点
            this.$refs.inputRef.focus();
        });
    },
    methods: {
        handleInput() {
            // 可以在输入时执行一些操作
        },
        handleSearchOnEnter(event){
            if (event.keyCode === 13) {
                this.handleSearch()
            }
        },
        handleSearch() {
            console.log("搜索")
            if (this.searchQuery === ""){
                this.$message({
                    showClose: true,
                    message: '搜索为空',
                    type: 'warning',
                    center: true
                });
            }
            else{
                if (this.lastSearchQuery!==this.searchQuery){
                    this.searchTotalLabelsCount(this.searchQuery)
                    this.searchLabelPhotos(1, this.searchQuery)
                    this.lastSearchQuery=this.searchQuery
                }
            }
        },
        parseSearchPhotoData() {
            for (let i = 0; i < this.searchPhotos.length; i++) {
                this.searchPhotos[i].metadata.date_taken = formatDateTime(this.searchPhotos[i].metadata.date_taken)
                this.searchPhotos[i].metadata.fileMB = kbToMb(this.searchPhotos[i].metadata.file_size)
                this.searchPhotos[i].metadata.latitude = getLatitude(this.searchPhotos[i].metadata.latitude)
                this.searchPhotos[i].metadata.longitude = getLongitude(this.searchPhotos[i].metadata.longitude)
                this.searchPhotos[i].metadata.altitude = getAltitude(this.searchPhotos[i].metadata.altitude)
            }
        },
        searchTotalLabelsCount(label) {
            axios.get(`/home/searchLabelPhotoCount?username=${this.username}&label=${label}`)
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
                    const photoIndex = this.searchPhotos.findIndex(photo => photo.id === photoId);
                    if (photoIndex !== -1) {
                        this.searchPhotos[photoIndex].collected = !this.searchPhotos[photoIndex].collected;
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
        async searchLabelPhotos(newPage, label) {
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
                            label: label
                        }
                    }); // 发起 GET 请求
                if (response.status === 200) {
                    // console.log(response.data)
                    if (response.data.photos!=null){
                        this.searchPhotos = response.data.photos
                        // console.log(this.photos)
                        this.parseSearchPhotoData()
                        // console.log(response.data.photos)
                        // 初始化 showPhotos 数组
                        this.showSearchPhotos = []
                        // 获取照片数据
                        const photosData = response.data.photos;
                        // 处理每张照片
                        for (const photoData of photosData) {
                            // 将二进制数据转换为可显示的图片
                            const img = this.getBase64Image(photoData.file_content) // 或者使用其他适合的库或工具来生成图片 URL
                            // 将转换后的图片添加到 showPhotos 数组中
                            this.showSearchPhotos.push(img);
                        }
                        this.$message({
                            showClose: true,
                            message: response.data.message,
                            type: 'success',
                            center: true
                        });
                    }
                    else {
                        this.searchPhotos=[]
                        this.showSearchPhotos = []
                    }

                    // console.log(this.photos)
                    this.$message({
                        message: response.data.message,
                        type: 'success',
                        center: true
                    });
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
};