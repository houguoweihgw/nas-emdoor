import Header from "@/components/header.vue";
import Footer from "@/components/footer.vue";
import axios from "axios";
import albums from "@/views/albums/albums";
import Search from "@/views/search/search";

export default {
    name: "home",
    components: {Search, Header, Footer },
    data() {
        return {
            username: this.$store.state.username,
            activeMenu: '1',
            submenuTitle: '我的相册',
            menuItems: [],

            dialogFormVisible: false,
            form: {
                name: '',
                description: '',
            },
            formLabelWidth: '100px'
        };
    },
    methods: {
        addAlbum(){
            this.dialogFormVisible = true
            this.$router.push({name: 'blank'});
        },
        handleCreateAlbum() {
            const albumName=this.form.name.trim()
            const albumDescription=this.form.description.trim()
            // 获取表单数据
            const albumData = {
                name: albumName,
                description: albumDescription
            };
            // 发送 POST 请求来保存相册数据
            axios.post(`/home/addAlbum?username=${this.username}`, albumData)
                .then(response => {
                    // 请求成功，执行成功的逻辑，可能包括刷新相册列表等等
                    this.fetchAlbums()
                    // 在完成逻辑后，关闭对话框
                    this.dialogFormVisible = false;
                })
                .catch(error => {
                    // 请求失败，执行失败的逻辑，例如显示错误消息等等
                });
            // 在完成逻辑后，你可以关闭对话框
            this.dialogFormVisible = false;
        },
        handleCommand(command) {
            switch (command.operation){
                case 'rename':
                    this.startEditing(command.index)
                    // console.log(command.operation,command.index)
                    break;
                case 'delete':
                    this.deleteAlbum(command.index)
                    // console.log(command.operation,command.index)
                    break;
            }
        },
        beforeHandleCommand(index, command){
            return {
                'index': index,
                'operation':command
            }
        },
        startEditing(menuItem) {
            // 进入编辑状态
            menuItem.editing = true;
            menuItem.newName = menuItem.name;
            // todo：焦点获取
        },
        stopEditing(menuItem) {
            // 退出编辑状态
            menuItem.editing = false;
        },
        saveName(menuItem) {
            if (menuItem.name === menuItem.newName) {
                // 新名称与旧名称相同，不需要执行重命名
                this.stopEditing(menuItem);
                return;
            }
            // 构造请求参数
            // 发送重命名请求
            axios.put(`/home/updateAlbum?album=${menuItem.id}&newName=${menuItem.newName}`)
                .then(response => {
                    // 请求成功
                    if (response.status === 200) {
                        // 更新本地数据或执行其他操作
                        menuItem.name = menuItem.newName;
                        menuItem.newName = '';
                        console.log(menuItem);
                        this.stopEditing(menuItem);
                    } else {
                        // 处理请求失败的情况
                        console.error('Rename request failed with status:', response.status);
                    }
                })
                .catch(error => {
                    // 处理网络错误
                    console.error('Rename request error:', error);
                });
        },
        deleteAlbum(menuItem) {
            this.$confirm('确定要删除相册 ' + menuItem.name + ' 吗?', '警告', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                // 用户点击了确定按钮
                this.handleDeleteAlbum(menuItem);
            }).catch(() => {
                // 用户点击了取消按钮
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
        },
        handleDeleteAlbum(menuItem) {
            // 向后端发送删除相册的请求
            axios
                .delete(`/home/deleteAlbum/${menuItem.id}`)
                .then((response) => {
                    if (response.status === 200) {
                        // 请求成功，执行以下操作
                        this.$message({
                            type: 'success',
                            message: '相册删除成功'
                        });
                        // 可以刷新相册列表，更新视图，以确保删除后的数据是最新的
                        this.fetchAlbums(); // 例如，调用一个方法来刷新相册列表
                    } else {
                        // 请求成功，但删除操作失败，你可以根据后端返回的响应进行处理
                        this.$message.error('删除相册失败');
                    }
                })
                .catch((error) => {
                    // 请求失败，处理错误
                    console.error('删除相册失败：', error);
                    this.$message.error('删除相册时发生错误');
                });
        },
        fetchAlbums(){
            axios.get(`/home/getAlbums?username=${this.username}`, {responseType: 'json'}) // 请根据你的后端路由进行修改
                .then(response => {
                    // 将从后端获取的相册数据存储在menuItems中
                    // this.resposeManu = response.data.albums
                    const albumArray=[]
                    this.menuItems=response.data.albums.map((album, ind) => {
                        albumArray.push(album)
                        return {
                            ...album,
                            index: (ind+1).toString(), // 添加 index 属性，从1开始递增
                            editing: false,   // 添加 editing 属性
                            newName: ''       // 添加 newName 属性
                        };
                    })
                    this.$store.commit('SET_Album_Array', albumArray)
                    // console.log("albumsWithAttributes",this.menuItems)
                })
                .catch(error => {
                    console.error('获取相册数据失败:', error)
                });
        }
    },
    mounted() {
        // 在组件挂载后，请求后端获取相册数据
        this.fetchAlbums()
    }
};