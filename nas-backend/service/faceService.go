package service

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"nas-backend/entity"
	"nas-backend/tools"
	"nas-backend/utils"
	"net/http"
	"strconv"
)

// GetFaceClustersHandler 处理用户请求所有人脸聚类结果
func GetFaceClustersHandler(c *gin.Context) {
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
	//根据用户ID查询该用户所有人脸聚类
	clusters, err := entity.GetClustersByUserID(utils.DB, userId)
	if err != nil {
		log.Printf("Query User Face Clusters Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询所有人脸聚类失败",
		})
	}
	log.Printf("用户%s查询所有人脸聚类数量%d", username, len(clusters))
	// 创建一个包含所有人脸聚类信息的切片
	var labelData []map[string]interface{}
	//根据查询的人脸聚类列表查询每个人脸聚类对应的一张照片作为封面，以及数量
	for _, cluster := range clusters {
		clusterCount, err := entity.GetFaceNumberByClusterID(utils.DB, cluster.ID)
		if err != nil {
			log.Printf("Query User Face Clusters Count Failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询人脸照片数量失败",
			})
		}
		if clusterCount <= 2 {
			continue
		}
		face, err2 := entity.GetCoversForCluster(utils.DB, cluster.ID)
		if err2 != nil {
			log.Printf("Query User Face Clusters Cover Failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询人脸封面失败",
			})
		}
		var item entity.Item
		utils.DB.Model(&entity.Item{}).Where("id = ?", face.ItemID).Find(&item)
		imageFile, err := tools.PhotoCropping(item.FilePath, face.X, face.Y, face.W, face.H)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		labelInfo := map[string]interface{}{
			"face_name":  cluster.Name,
			"face_count": clusterCount,
			"face_cover": imageFile,
		}
		labelData = append(labelData, labelInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"faces":   labelData,
		"message": "查询人脸聚类成功",
	}
	log.Printf("用户%s查询所有标签及其封面成功", username)
	c.JSON(http.StatusOK, responseData)
}

// GetClusterTotalPhotosCountHandler 处理用户查询某人脸聚类的照片数量
func GetClusterTotalPhotosCountHandler(c *gin.Context) {
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
	// 从请求中获取标签
	clusterName := c.Query("cluster")
	log.Printf("用户 %d 请求查询场景分类标签 %s 的全部照片总数", userId, clusterName)
	//根据用户的ID和聚类名查询聚类ID
	clusterID, err := entity.GetClusterIDByUserIDAndName(utils.DB, userId, clusterName)
	if err != nil {
		log.Printf("Query Cluster ID Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询人脸聚类失败",
		})
		return
	}
	//根据聚类ID查询聚类照片数量
	userClusterCount, err := entity.GetFaceNumberByClusterID(utils.DB, clusterID)
	if err != nil {
		log.Printf("Query Label Total Photos Count Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询人脸照片数量失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   userClusterCount,
		"message": "查询标签照片数量成功",
	})
}

// GetClusterPhotosHandler 处理用户请求某聚类的所有激活照片
func GetClusterPhotosHandler(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("perPage", "24"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "请求数据错误"})
		return
	}
	// 从请求中获取标签
	clusterName := c.Query("cluster")
	log.Printf("用户 %d 请求查询人脸聚类名 %s 的照片，页码 %d 单页照片 %d", userId, clusterName, page, perPage)
	//根据用户的ID和标签查询所有照片
	userClusterPhotos, err := entity.GetUserClusterPhotos(utils.DB, userId, clusterName, page, perPage)
	if err != nil {
		log.Printf("User Query Label Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询人脸所有照片失败",
		})
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userClusterPhotos {
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
		"message": "查询人脸聚类照片成功",
	}
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}
