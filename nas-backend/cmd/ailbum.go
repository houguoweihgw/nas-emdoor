package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"nas-backend/service"
	"nas-backend/utils"
)

func main() {
	r := gin.Default()
	//初始化GORM
	utils.InitDB()
	//初始化Redis
	err := utils.InitRedisClient()
	if err != nil {
		log.Println("redis 初始化失败")
		return
	}
	//初始化NATS
	utils.NatsInit()
	//启用跨域资源共享 (CORS) 支持
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许的前端地址
	config.AllowMethods = []string{"*"} // 允许的请求方法
	config.AllowHeaders = []string{"*"} // 允许的请求头
	r.Use(cors.New(config))
	// login路由组: login
	api := r.Group("/api")
	{
		api.POST("/login", service.LoginHandler)
		api.POST("/register", service.RegisterHandler)
	}
	// home路由组: home
	home := r.Group("/home")
	{
		home.GET("/photos", service.GetPhotosHandler)
		home.GET("/welcomePhotos", service.GetWelcomePhotosHandler)
		home.GET("/totalPhotosCount", service.GetTotalPhotosCountHandler)
		home.GET("/labelPhotoCount", service.GetLabelTotalPhotosCountHandler)
		home.GET("/clusterPhotoCount", service.GetClusterTotalPhotosCountHandler)
		home.GET("/searchLabelPhotoCount", service.GetSearchLabelTotalPhotosCountHandler)
		home.GET("/labelPhotos", service.GetLabelPhotosHandler)
		home.GET("/clusterPhotos", service.GetClusterPhotosHandler)
		home.GET("/albumPhotos", service.GetAlbumPhotosHandler)
		home.GET("/getAlbums", service.GetAlbumsHandler)
		home.GET("/myCollected", service.GetCollectedPhotosHandler)
		home.GET("/recycledPhotos", service.GetRecycledPhotosHandler)
		home.GET("/sceneLabels", service.GetSceneLabelsHandler)
		home.GET("/faceClusters", service.GetFaceClustersHandler)
		home.DELETE("/photos/:username/:photoID", service.RemovePhotosHandler)
		home.DELETE("/album/:albumID/:photoID", service.RemoveOnePhotoFromAlbumHandler)
		home.DELETE("/recycledPhotos/:photoID", service.DeletePhotosHandler)
		home.DELETE("/deleteAlbum/:albumID", service.DeleteAlbumHandler)
		home.POST("/upload", service.UploadPhotosHandler)
		home.POST("/addAlbum", service.CreateAlbumHandler)
		home.POST("/deleteSelectedPhotos", service.DeleteSelectedPhotosHandler)
		home.POST("/batchAddPhotosToAlbum", service.AddSelectedPhotosToAlbumHandler)
		home.POST("/recoverBatchPhotos", service.RecoverSelectedPhotosHandler)
		home.POST("/batchRemoveFromAlbum", service.RemovePhotosFromAlbumHandler)
		home.POST("/batchDeletePhotos", service.BatchDeletePhotosHandler)
		home.PUT("/toggleCollected", service.TogglePhotoCollectedHandler)
		home.PUT("/updateAlbum", service.UpdateAlbumHandler)
	}
	err = r.Run(":8001")

	if err != nil {
		return
	} //

	// 阻塞主线程
	select {}
}
