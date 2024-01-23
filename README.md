## 前提：

    1.linux系统
    2.安装docker和docker-compose环境

## 使用步骤：

    1.运行 docker-compose up --build
    2.浏览器访问 http://localhost:8082 即可进入nas管理系统

## 系统框架

    1.web前端使用Vue2+Element-ui+Router+Axios
    2.后端语言使用Go，框架为：Gin+GORM+NATS，NATS消息队列主要用于和AI算法交互
    3.数据库使用Mysql+Redis，Mysql用于存储用户信息，Redis用于存储用户登录状态和一些缓存数据
    4.AI算法使用Python+PyTorch
    整体框架图如下图所示：
    其中安卓移动端地址为：https://github.com/houguoweihgw/YiSpace
![框架.png](%E6%A1%86%E6%9E%B6.png)

## 系统功能：

    1.用户注册，登录
    2.照片上传及元数据解析
    3.首页照片推荐
    4.全部照片展示
    5.照片详情展示
    6.AI场景分类
    7.AI人脸聚类
    8.照片收藏
    9.照片回收站

## 算法原理
    1.AI场景分类
        使用ResNet50模型+places365数据集，对照片进行场景分类，技术路线如下图；
    2.AI人脸聚类
        使用FaceNet+VGGface数据集+DBSCAN聚类算法，对照片进行人脸聚类，技术路线如下图；
![AI场景识别.png](AI%E5%9C%BA%E6%99%AF%E8%AF%86%E5%88%AB.png)
![AI人脸聚类.png](AI%E4%BA%BA%E8%84%B8%E8%81%9A%E7%B1%BB.png)
        
    
## 注意事项：

1. master分支是日常提交的测试版本
2. stable分支是稳定发行版
