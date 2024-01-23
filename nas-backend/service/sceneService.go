package service

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"nas-backend/entity"
	"nas-backend/utils"
	"net/http"
	"strconv"
)

// GetSceneLabelsHandler 处理用户请求所有场景标签
func GetSceneLabelsHandler(c *gin.Context) {
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
	//根据用户ID查询该用户所有标签数量大于0的标签
	sceneLabels, err := entity.GetSceneLabelsInfoByUserID(utils.DB, userId)
	if err != nil {
		log.Printf("Query User Scene Labels Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询场景标签失败",
		})
	}
	log.Printf("用户%s查询所有标签数量%d", username, len(sceneLabels))
	// 创建一个包含所有图片信息的切片
	var labelData []map[string]interface{}
	//根据查询的标签列表查询每个标签对应的一张照片作为封面
	for _, sceneLabel := range sceneLabels {
		coverURL, err2 := entity.GetCoversForLabel(utils.DB, sceneLabel.LabelName, sceneLabel.UserID)
		if err2 != nil {
			log.Printf("Query User Scene Labels Cover Failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询场景封面失败",
			})
		}
		imageFile, err := ioutil.ReadFile(coverURL)
		if err != nil {
			log.Printf("Failed to read image file: %v", err)
			continue // 如果无法读取图片文件，跳过该图片
		}
		labelInfo := map[string]interface{}{
			"label_name":  sceneLabel.LabelName,
			"label_count": sceneLabel.LabelCount,
			"label_cover": imageFile,
		}
		labelData = append(labelData, labelInfo)
	}
	// 构建响应数据
	responseData := gin.H{
		"labels":  labelData,
		"message": "查询场景标签成功",
	}
	log.Printf("用户%s查询所有标签及其封面成功", username)
	c.JSON(http.StatusOK, responseData)
}

// GetLabelTotalPhotosCountHandler 处理用户查询某标签的照片数量
func GetLabelTotalPhotosCountHandler(c *gin.Context) {
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
	labelName := c.Query("label")
	log.Printf("用户 %d 请求查询场景分类标签 %s 的全部照片总数", userId, labelName)
	//根据用户的ID查询激活照片数量
	userPhotosCount, err := entity.GetLabelTotalPhotosCount(utils.DB, userId, labelName)
	if err != nil {
		log.Printf("Query Label Total Photos Count Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询场景照片数量失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   userPhotosCount,
		"message": "查询标签照片数量成功",
	})
}

// GetSearchLabelTotalPhotosCountHandler 处理用户查询标签的照片数量
func GetSearchLabelTotalPhotosCountHandler(c *gin.Context) {
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
	labelName := c.Query("label")
	log.Printf("用户 %d 请求搜索场景分类标签 %s 的全部照片总数", userId, labelName)
	//根据用户的ID查询激活照片数量
	userPhotosCount, err := entity.GetSearchLabelTotalPhotosCount(utils.DB, userId, labelName)
	if err != nil {
		log.Printf("Query Label Total Photos Count Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询标签照片数量失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   userPhotosCount,
		"message": "查询标签照片数量成功",
	})
}

// GetLabelPhotosHandler 处理用户请求某标签的所有激活照片
func GetLabelPhotosHandler(c *gin.Context) {
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
	labelName := c.Query("label")
	log.Printf("用户 %d 请求查询场景分类标签 %s 的照片，页码 %d 单页照片 %d", userId, labelName, page, perPage)
	//根据用户的ID和标签查询所有照片
	userLabelPhotos, err := entity.GetUserLabelPhotos(utils.DB, userId, labelName, page, perPage)
	if err != nil {
		log.Printf("User Query Label Photos Failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询场景照片失败",
		})
	}
	// 创建一个包含所有图片信息的切片
	var photoData []map[string]interface{}
	for _, photo := range userLabelPhotos {
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
		"message": "查询场景照片成功",
	}
	// 返回 JSON 响应
	c.JSON(http.StatusOK, responseData)
}
