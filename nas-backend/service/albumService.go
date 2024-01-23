package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"nas-backend/entity"
	"nas-backend/utils"
	"net/http"
	"strconv"
)

func GetAlbumsHandler(c *gin.Context) {
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
	//根据用户的ID查询所有相册列表
	albums, err := entity.QueryAllAlbumsByUserId(utils.DB, userId)
	if err != nil {
		log.Printf("User Query Albums Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询用户相册列表失败",
		})
		return
	}
	// 遍历 albums 切片，排除默认相册，名为 "default" 的相册
	filteredAlbums := make([]entity.Album, 0)
	for _, album := range albums {
		if album.Name != "default" {
			filteredAlbums = append(filteredAlbums, album)
		}
	}
	// 将 filteredAlbums 返回给客户端
	responseData := gin.H{
		"albums":  filteredAlbums,
		"message": "查询用户相册列表成功",
	}
	c.JSON(http.StatusOK, responseData)
}

func GetAlbumPhotosHandler(c *gin.Context) {
	// 从请求中获取用户名和相册Id
	username := c.Query("username")
	albumID := c.Query("album")
	albumId, err := strconv.Atoi(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
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
	//根据用户的ID查询所有激活照片
	//userAlbumAllPhotos, err := entity.GetAllPhotosByUserIDAndAlbumID(utils.DB, userId, albumId)
	userAlbumAllPhotos, err := entity.GetItemsInAlbum(utils.DB, uint(albumId))
	if err != nil {
		log.Printf("User Query Album Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户查询相册照片失败",
		})
		return
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userAlbumAllPhotos {
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
		}
		photoData = append(photoData, photoInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"photos":  photoData,
		"message": "查询相册照片成功",
	}
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}

// UpdateAlbumHandler 处理用户相册重命名
func UpdateAlbumHandler(c *gin.Context) {
	// 获取相册ID和新名称
	albumID := c.Query("album")
	newName := c.Query("newName")
	log.Printf("用户请求将相册 %s 重命名为：%s", albumID, newName)
	albumId, err := strconv.Atoi(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	// 执行数据库更新操作
	err = entity.UpdateAlbumNameByAlbumID(utils.DB, albumId, newName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "相册命名失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "相册命名成功"})
}

func DeleteAlbumHandler(c *gin.Context) {
	//收到将要删除的照片id，并强转成int
	albumID := c.Param("albumID")
	id, err := strconv.Atoi(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	log.Printf("收到删除相册请求id:%d", id)
	//在items表中清除该照片
	err = entity.DeleteAlbumByAlbumID(utils.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除相册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除相册成功"})
}

func CreateAlbumHandler(c *gin.Context) {
	// 解析 JSON 请求正文
	userName := c.Query("username")
	var album entity.Album
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求数据错误",
		})
		return
	}
	log.Printf("用户 %s 请求添加相册 %s ,相册描述为 %s", userName, album.Name, album.Description)
	//根据用户名查询users表获取用户ID
	userID := entity.GetUserIdByUsername(utils.DB, userName)
	album.UserID = uint(userID)
	//在albums表中添加一个新的相册
	err := entity.CreateAlbum(utils.DB, album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "添加相册失败",
		})
		return
	}
	log.Printf("用户 %s 添加相册成功", userName)
	c.JSON(http.StatusOK, gin.H{"message": "添加相册成功"})
}
