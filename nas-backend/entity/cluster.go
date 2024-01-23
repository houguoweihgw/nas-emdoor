package entity

import (
	"github.com/jinzhu/gorm"
	"log"
	"nas-backend/tools"
	"time"
)

type Cluster struct {
	ID       int    `gorm:"primary_key" json:"id"`
	UserID   int    `gorm:"not null" json:"user_id"`
	Name     string `gorm:"not null" json:"name"`
	Features []byte `gorm:"not null" json:"features"`
}

func (Cluster) TableName() string {
	return "clusters"
}

// CreateCluster 插入一个人脸聚类
func CreateCluster(db *gorm.DB, userID int, name string, features []byte) (Cluster, error) {
	cluster := Cluster{
		UserID:   userID,
		Name:     name,
		Features: features,
	}
	if err := db.Create(&cluster).Error; err != nil {
		log.Printf("插入clusters表失败")
		return cluster, err
	}
	log.Printf("插入clusters表成功，clusterID: %d ", cluster.ID)
	return cluster, nil
}

// InsertFaceIntoCluster 人脸验证
func InsertFaceIntoCluster(db *gorm.DB, userID int, face Face) error {
	// 获取所有已有的人脸分类
	var clusters []Cluster
	if err := db.Where("user_id = ?", userID).Find(&clusters).Error; err != nil {
		return err
	}
	differenceThreshold := 0.8
	// 遍历已有的分类
	for _, cluster := range clusters {
		difference := tools.CalculateEuclideanDistance(face.Features, cluster.Features)
		if difference < differenceThreshold {
			// 如果相似度小于阈值，将人脸添加到该分类
			log.Printf("与人脸类 %d %s 的差异度为：%f ,小于阈值被归为一类", cluster.ID, cluster.Name, difference)
			// 创建faces和Cluster的中间映射记录
			err := CreateClusterFace(db, uint(cluster.ID), face.ID)
			if err != nil {
				return err
			}
			return nil
		}
	}
	// 如果没有找到相似的分类，创建一个新分类
	log.Printf("没有找到相似的分类，创建一个新分类")
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05") // 自定义时间格式
	cluster, err := CreateCluster(db, userID, currentTimeString, face.Features)
	if err != nil {
		return err
	}
	// 创建faces和Cluster的中间映射记录
	err = CreateClusterFace(db, uint(cluster.ID), face.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetClustersByUserID 根据用户ID查询所有人脸分类
func GetClustersByUserID(db *gorm.DB, userID int) ([]Cluster, error) {
	var clusters []Cluster
	if err := db.Where("user_id = ?", userID).Find(&clusters).Error; err != nil {
		log.Printf("用户 %d 查询clusters表失败", userID)
		return clusters, err
	}
	log.Printf("用户 %d 查询clusters表查到 %d 个人脸聚类", userID, len(clusters))
	return clusters, nil
}

// GetClusterIDByUserIDAndName 根据用户ID和聚类名查询所有聚类ID
func GetClusterIDByUserIDAndName(db *gorm.DB, userID int, clusterName string) (int, error) {
	var cluster Cluster
	if err := db.Where("user_id = ? And name = ?", userID, clusterName).Find(&cluster).Error; err != nil {
		log.Printf("用户 %d 查询名为 %s 的聚类失败", userID, clusterName)
		return cluster.ID, err
	}
	log.Printf("用户 %d 查询名为 %s 的聚类成功，ID：%d", userID, clusterName, cluster.ID)
	return cluster.ID, nil
}
