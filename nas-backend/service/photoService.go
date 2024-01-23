package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"nas-backend/entity"
	"nas-backend/utils"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// GetPhotosHandler 处理用户请求所有激活照片
func GetPhotosHandler(c *gin.Context) {
	// 从请求中获取用户名
	username := c.Query("username")
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	// 从请求中获取页码page和每页显示的照片数量perPage
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无效请求页码"})
		return
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("perPage", "24"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无效请求页码"})
		return
	}
	log.Printf("用户 %d 请求查询照片，页码 %d 单页照片 %d", userId, page, perPage)
	// 检查Redis缓存是否存在
	cacheKey := fmt.Sprintf("user_photos:%d:%d:%d", userId, page, perPage)
	cachedData, err := utils.RDB.Get(cacheKey).Result()
	if err == nil {
		// 缓存命中，直接返回缓存数据
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误，请重试",
			})
			return
		}
		log.Printf("用户 %d 命中所有缓存照片 %s", userId, cacheKey)
		c.JSON(http.StatusOK, response)
		return
	}
	// 计算偏移量
	offset := (page - 1) * perPage
	//根据用户的ID查询该页的激活照片
	userPagePhotos, err := entity.GetPagePhotosByUserID(utils.DB, userId, offset, perPage)
	if err != nil {
		log.Printf("User Query Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户查询失败，请重试",
		})
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userPagePhotos {
		// 读取图片文件
		imageFile, err := ioutil.ReadFile(photo.FilePath)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		photoInfo := map[string]interface{}{
			"id":           photo.ID,
			"title":        photo.Title,
			"description":  photo.Description,
			"upload_date":  photo.DateUploaded,
			"file_content": imageFile, // 这里可以是图片的 URL 或者文件路径
			"collected":    photo.Collected,
			"metadata":     photo.ItemMetadata,
		}
		photoData = append(photoData, photoInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"photos":  photoData,
		"message": "查询照片成功",
	}
	// 将数据存入Redis缓存
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Failed to marshal response data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部错误，请重试",
		})
		return
	}
	// 设置缓存过期时间，10 分钟
	utils.RDB.Set(cacheKey, jsonData, 10*time.Minute)
	log.Printf("用户 %d 缓存所有照片 %s", userId, cacheKey)
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}

// GetWelcomePhotosHandler 处理用户请求欢迎界面照片
func GetWelcomePhotosHandler(c *gin.Context) {
	// 从请求中获取用户名
	username := c.Query("username")
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	//根据用户的ID查询的激活照片
	userPagePhotos, err := entity.GetWelcomePhotosByUserID(utils.DB, userId, 6)
	if err != nil {
		log.Printf("User Query Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询欢迎页面失败，请重试",
		})
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userPagePhotos {
		// 读取图片文件
		imageFile, err := ioutil.ReadFile(photo.FilePath)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		photoInfo := map[string]interface{}{
			"id":           photo.ID,
			"file_content": imageFile, // 这里可以是图片的 URL 或者文件路径
		}
		photoData = append(photoData, photoInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"photos":  photoData,
		"message": "欢迎进入云相册",
	}
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}

// GetTotalPhotosCountHandler 处理用户查询照片数量
func GetTotalPhotosCountHandler(c *gin.Context) {
	// 从请求中获取用户名
	username := c.Query("username")
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	//根据用户的ID查询激活照片数量
	userPhotosCount, err := entity.GetTotalPhotosCountByUserID(utils.DB, userId)
	if err != nil {
		log.Printf("Query User Total Photos Count Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户查询照片数量失败，请重试",
		})
		return
	}
	log.Printf("用户 %d 请求查询全部照片总数为 %d", userId, userPhotosCount)
	c.JSON(http.StatusOK, gin.H{
		"total":   userPhotosCount,
		"message": "查询用户所有照片成功",
	})
}

// GetRecycledPhotosHandler 处理用户请求所有回收站照片
func GetRecycledPhotosHandler(c *gin.Context) {
	// 从请求中获取用户名
	username := c.Query("username")
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	// 检查Redis缓存是否存在
	cacheKey := fmt.Sprintf("user_recycles:%d", userId)
	cachedData, err := utils.RDB.Get(cacheKey).Result()
	if err == nil {
		// 缓存命中，直接返回缓存数据
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器内部错误，请重试",
			})
			return
		}
		log.Printf("用户 %d 命中回收站缓存照片 %s", userId, cacheKey)
		c.JSON(http.StatusOK, response)
		return
	}
	//根据用户的ID嘻查询所有回收站照片
	userAllPhotos, err := entity.GetAllRecycledPhotosByUserID(utils.DB, userId)
	if err != nil {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户查询回收站照片失败，请重试",
		})
		return
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userAllPhotos {
		// 读取图片文件
		imageFile, err := ioutil.ReadFile(photo.FilePath)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		photoInfo := map[string]interface{}{
			"id":           photo.ID,
			"title":        photo.Title,
			"description":  photo.Description,
			"upload_date":  photo.DateUploaded,
			"file_content": imageFile, // 这里可以是图片的 URL 或者文件路径
			"metadata":     photo.ItemMetadata,
		}
		photoData = append(photoData, photoInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"recycledPhotos": photoData,
		"message":        "查询回收站照片成功",
	}
	// 将数据存入Redis缓存
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Failed to marshal response data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部错误，请重试",
		})
		return
	}
	// 设置缓存过期时间，10 分钟
	utils.RDB.Set(cacheKey, jsonData, 10*time.Minute)
	log.Printf("用户 %d 缓存回收站照片 %s", userId, cacheKey)
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}

// RemovePhotosHandler 处理用户请求删除照片
func RemovePhotosHandler(c *gin.Context) {
	//收到将到删除的照片id，并强转成int
	photoID := c.Param("photoID")
	username := c.Param("username")
	id, err := strconv.Atoi(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "照片错误，请重试",
		})
		return
	}
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	log.Printf("收到用户 %d 删除照片请求id:%d", userId, id)
	//在items表中删除该照片
	err = entity.RemovePhotoByPhotoID(utils.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除照片错误，请重试"})
		return
	}
	//删除该用户所有照片的缓存
	cachePatternKey := fmt.Sprintf("user_photos:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, cachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	//删除该用户回收站照片的缓存
	recycleCachePatternKey := fmt.Sprintf("user_recycles:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, recycleCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除照片成功"})
}

// RemoveOnePhotoFromAlbumHandler 处理用户将一张照片移出相册
func RemoveOnePhotoFromAlbumHandler(c *gin.Context) {
	//收到将到删除的照片id，并强转成int
	photoID := c.Param("photoID")
	albumID := c.Param("albumID")
	photoId, err := strconv.Atoi(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "照片错误，请重试"})
		return
	}
	albumId, err := strconv.Atoi(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "相册错误，请重试"})
		return
	}
	log.Printf("用户请求将照片 %d 移出到相册 %d", photoId, albumId)
	// 在这里执行将所选照片移出相册的逻辑
	if err := entity.RemovePhotoFromAlbum(utils.DB, photoId, albumId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "照片移除相册失败，请重试",
		})
		return
	}
	// 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "照片成功移出相册"})
}

// RemovePhotosFromAlbumHandler 处理用户将选中照片移出相册
func RemovePhotosFromAlbumHandler(c *gin.Context) {
	// 从请求体中获取相册ID和所选照片ID数组
	var requestData struct {
		AlbumID        int   `json:"albumId"`
		SelectedPhotos []int `json:"selectedPhotos"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据错误，请重试"})
		return
	}
	log.Printf("用户请求将照片 %d 移出到相册 %d", requestData.SelectedPhotos, requestData.AlbumID)
	// 在这里执行将所选照片移出相册的逻辑
	albumID := requestData.AlbumID
	selectedPhotos := requestData.SelectedPhotos
	for _, photoID := range selectedPhotos {
		if err := entity.RemovePhotoFromAlbum(utils.DB, photoID, albumID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "照片移出相册失败，请重试",
			})
			return
		}
	}
	// 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "照片成功移出相册"})
}

// BatchDeletePhotosHandler 处理用户将选中照片从回收站中删除
func BatchDeletePhotosHandler(c *gin.Context) {
	// 从请求体中获取相册ID和所选照片ID数组
	var requestData struct {
		SelectedPhotos []int  `json:"selectedPhotos"`
		Username       string `json:"username"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据错误，请重试"})
		return
	}
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, requestData.Username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	log.Printf("用户请求将照片 %d 从回收站中彻底删除", requestData.SelectedPhotos)
	// todo:在这里执行将所选照片移出相册的逻辑,要开启一个事务来一同删除item_metadata表
	selectedPhotos := requestData.SelectedPhotos
	for _, photoID := range selectedPhotos {
		if err := entity.DeletePhotoByPhotoID(utils.DB, photoID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "照片移出相册失败，请重试",
			})
			return
		}
	}
	//删除该用户回收站照片的缓存
	recycleCachePatternKey := fmt.Sprintf("user_recycles:%d", userId)
	err := utils.DeleteKeysByPattern(utils.RDB, recycleCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	// 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "删除照片成功"})
}

// DeletePhotosHandler 处理用户请求清除照片
func DeletePhotosHandler(c *gin.Context) {
	//收到将到清除的照片id，并强转成int
	photoID := c.Param("photoID")
	id, err := strconv.Atoi(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	log.Printf("收到清除照片请求id:%d", id)
	//在items表中清除该照片
	err = entity.DeletePhotoByPhotoID(utils.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "清除照片成功"})
}

// DeleteSelectedPhotosHandler 处理用户请求清除选中照片
func DeleteSelectedPhotosHandler(c *gin.Context) {
	//解析请求体中的 JSON 数据
	var photos struct {
		PhotoIDs []int  `json:"photoIDs"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&photos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求数据错误",
		})
		return
	}
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, photos.Username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	log.Printf("用户请求批量删除照片：%d", photos)
	// 遍历 request.PhotoIDs，将每个 ID 对应的照片从数据库中删除
	for _, photoID := range photos.PhotoIDs {
		if err := entity.RemovePhotoByPhotoID(utils.DB, photoID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "清除照片失败，请重试",
			})
			return
		}
	}
	//删除该用户所有照片的缓存
	cachePatternKey := fmt.Sprintf("user_photos:%d", userId)
	err := utils.DeleteKeysByPattern(utils.RDB, cachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	//删除该用户回收站照片的缓存
	recycleCachePatternKey := fmt.Sprintf("user_recycles:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, recycleCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "清除照片成功",
	})
}

// AddSelectedPhotosToAlbumHandler 处理用户请求添加选中照片到指定相册
func AddSelectedPhotosToAlbumHandler(c *gin.Context) {
	// 从请求体中获取相册ID和所选照片ID数组
	var requestData struct {
		AlbumID        int   `json:"albumId"`
		SelectedPhotos []int `json:"selectedPhotos"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据错误"})
		return
	}
	log.Printf("用户请求将照片 %d 添加到相册 %d", requestData.SelectedPhotos, requestData.AlbumID)

	// 在这里执行将所选照片添加到相册的逻辑
	albumID := requestData.AlbumID
	selectedPhotos := requestData.SelectedPhotos
	for _, photoID := range selectedPhotos {
		if err := entity.AddPhotoToAlbum(utils.DB, photoID, albumID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "添加照片到相册失败",
			})
			return
		}
	}
	// 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "添加照片到相册成功"})
}

// RecoverSelectedPhotosHandler 处理用户请求从回收站中恢复选中照片
func RecoverSelectedPhotosHandler(c *gin.Context) {
	// 从请求体中获取相册ID和所选照片ID数组
	var requestData struct {
		SelectedPhotos []int  `json:"selectedPhotos"`
		Username       string `json:"username"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据错误"})
		return
	}
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, requestData.Username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	log.Printf("用户请求从回收站中恢复照片 %d ", requestData.SelectedPhotos)
	// 遍历 request.PhotoIDs，将每个 ID 对应的照片从数据库中恢复
	for _, photoID := range requestData.SelectedPhotos {
		if err := entity.RecoverPhotoByPhotoID(utils.DB, photoID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "照片回收失败",
			})
			return
		}
	}
	//删除该用户所有照片的缓存
	cachePatternKey := fmt.Sprintf("user_photos:%d", userId)
	err := utils.DeleteKeysByPattern(utils.RDB, cachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	//删除该用户回收站照片的缓存
	recycleCachePatternKey := fmt.Sprintf("user_recycles:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, recycleCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "照片回收成功",
	})
}

// UploadPhotosHandler 处理用户上传照片
func UploadPhotosHandler(c *gin.Context) {
	// 拿到上传照片的用户名和照片文件
	username := c.PostForm("username")
	file, err := c.FormFile("file") // "file"是前端上传文件字段的名称
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("用户 %s 请求上传照片, 照片名为：%s", username, file.Filename)

	// 构建用户文件夹路径
	userFolder := "./nas-data/items/" + username
	filePath := path.Join(userFolder, file.Filename)
	// 创建用户文件夹（如果不存在）
	if _, err := os.Stat(userFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
			return
		}
	}

	// 根据用户名查询users表获取用户id
	userID := entity.GetUserIdByUsername(utils.DB, username)
	if userID < 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "用户不存在"})
		return
	}

	// 将文件保存到指定目录
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器保存照片失败，请重试"})
		return
	}

	// 解析照片元数据
	itemMeta := &entity.ItemMetadata{}
	err = utils.ParsePictureExif(filePath, itemMeta)
	if err != nil {
		log.Printf("新建一个空的item_metadata")
		//因为有的照片没有元数据,故不做处理返回空的ItemMetadata,但是配置一个虚拟时间,不然插入item_metadata表出错
		date := "1900:01:01 00:00:00"
		parsedDate, _ := time.Parse("2006:01:02 15:04:05", date)
		itemMeta.DateTaken = parsedDate
	}
	itemMeta.UserID = uint(userID)

	item := &entity.Item{
		Title:       file.Filename,
		Description: file.Filename,
		FilePath:    filePath,
		UserID:      uint(userID),
	}
	err = entity.InsertPhotoWithMetadata(utils.DB, item, itemMeta)
	if err != nil {
		return
	}

	// 打开上传的照片
	uploadedFile, err := file.Open()
	if err != nil {
		// 处理错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器打开照片失败，请重试"})
		return
	}
	defer uploadedFile.Close()

	// 读取照片数据
	fileData, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		// 处理错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器读取照片失败，请重试"})
		return
	}
	// 封装场景分类消息，
	message := entity.SCRequest{
		ID:      int(item.ID),
		UserID:  userID,
		Picture: fileData,
	}
	//将场景分类消息结构体转换为 JSON 字符串，并publish
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("无法转换结构体为 JSON: %v", err)
	}
	err = utils.NATSPublish(jsonData)
	err = utils.NATSFCPublish(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "照片AI分类失败"})
		return
	}

	//删除该用户所有照片的缓存
	cachePatternKey := fmt.Sprintf("user_photos:%d", userID)
	err = utils.DeleteKeysByPattern(utils.RDB, cachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}

	log.Printf("照片保存服务器")
	c.JSON(http.StatusOK, gin.H{"message": "上传照片成功"})
}

// TogglePhotoCollectedHandler 处理用户更改收藏状态请求
func TogglePhotoCollectedHandler(c *gin.Context) {
	// 从请求中获取照片id，并类型为int
	photoID := c.Query("photo")
	username := c.Query("username")
	photoId, err := strconv.Atoi(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求数据错误"})
		return
	}
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	log.Printf("用户 %d 请求更改照片：%d 的收藏状态", userId, photoId)
	//根据照片id更改对应照片collected状态
	err = entity.UpdatePhotoCollectedByPhotoID(utils.DB, photoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更改用户收藏状态失败"})
		return
	}
	//删除该用户所有照片的缓存
	allPhotoCachePatternKey := fmt.Sprintf("user_photos:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, allPhotoCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	//删除该用户收藏照片的缓存
	collectedCachePatternKey := fmt.Sprintf("user_collected:%d", userId)
	err = utils.DeleteKeysByPattern(utils.RDB, collectedCachePatternKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误，请重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "收藏照片成功"})
}

// GetCollectedPhotosHandler 处理用户查询收藏照片请求
func GetCollectedPhotosHandler(c *gin.Context) {
	// 从请求中获取用户名
	username := c.Query("username")
	// 根据用户名查询用户的ID
	userId := entity.GetUserIdByUsername(utils.DB, username)
	if userId < 0 {
		log.Printf("User Does NoT Exit")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "用户不存在",
		})
		return
	}
	// 检查Redis缓存是否存在
	cacheKey := fmt.Sprintf("user_collected:%d", userId)
	cachedData, err := utils.RDB.Get(cacheKey).Result()
	if err == nil {
		// 缓存命中，直接返回缓存数据
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器内部错误，请重试",
			})
			return
		}
		log.Printf("用户 %d 命中缓存收藏照片 %s", userId, cacheKey)
		c.JSON(http.StatusOK, response)
		return
	}
	//根据用户的ID查询所有激活照片
	userAllCollectedPhotos, err := entity.GetAllCollectedPhotosByUserID(utils.DB, userId)
	if err != nil {
		log.Printf("User Query Collected Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户查询收藏照片失败",
		})
		return
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userAllCollectedPhotos {
		// 读取图片文件
		imageFile, err := ioutil.ReadFile(photo.FilePath)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		photoInfo := map[string]interface{}{
			"id":           photo.ID,
			"title":        photo.Title,
			"description":  photo.Description,
			"upload_date":  photo.DateUploaded,
			"file_content": imageFile, // 这里可以是图片的 URL 或者文件路径
			"collected":    photo.Collected,
			"metadata":     photo.ItemMetadata,
		}
		photoData = append(photoData, photoInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"photos":  photoData,
		"message": "查询收藏照片成功",
	}
	// 将数据存入Redis缓存
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Failed to marshal response data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部错误，请重试",
		})
		return
	}
	// 设置缓存过期时间，10 分钟
	utils.RDB.Set(cacheKey, jsonData, 10*time.Minute)
	log.Printf("用户 %d 缓存收藏照片 %s", userId, cacheKey)
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}
