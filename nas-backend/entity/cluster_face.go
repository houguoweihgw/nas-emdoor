package entity

import (
	"github.com/jinzhu/gorm"
	"log"
)

type ClusterFace struct {
	ID        uint `gorm:"primaryKey"`
	ClusterID uint
	FaceID    uint
}

// TableName 设置表名
func (ClusterFace) TableName() string {
	return "cluster_faces"
}

// CreateClusterFace 插入
func CreateClusterFace(db *gorm.DB, clusterID, faceID uint) error {
	clusterFace := ClusterFace{
		ClusterID: clusterID,
		FaceID:    faceID,
	}
	if err := db.Create(&clusterFace).Error; err != nil {
		log.Printf("插入cluster_faces表失败，clusterID: %d faceID: %d", clusterID, faceID)
		return err
	}
	log.Printf("插入cluster_faces表成功，clusterID: %d faceID: %d", clusterID, faceID)
	return nil
}

// GetFaceNumberByClusterID 查询聚类结果的人脸照片数量
func GetFaceNumberByClusterID(db *gorm.DB, clusterID int) (int, error) {
	var count int
	if err := db.Model(&ClusterFace{}).
		Where("cluster_id = ?", clusterID).
		Count(&count).Error; err != nil {
		log.Printf("查询人脸聚类 %d 失败", clusterID)
		return 0, err
	}
	log.Printf("查询人脸聚类 %d 的人脸数量为: %d", clusterID, count)
	return count, nil
}

// GetUserClusterPhotos 获取用户聚类的照片
func GetUserClusterPhotos(db *gorm.DB, userID int, clusterName string, page int, pageSize int) ([]Item, error) {
	var items []Item
	// 构建查询条件
	query := db.Table("items").
		Select("items.*").Preload("ItemMetadataItem").
		Joins("JOIN faces ON items.id = faces.item_id").
		Joins("JOIN cluster_faces ON faces.id = cluster_faces.face_id").
		Joins("JOIN clusters ON cluster_faces.cluster_id = clusters.id").
		Where("items.user_id = ? AND items.status = ? AND clusters.name = ?",
			userID, "active", clusterName)
	// 计算分页偏移
	offset := (page - 1) * pageSize
	// 执行查询并应用分页
	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		log.Printf("用户id %d 查询聚类 %s 的照片为空", userID, clusterName)
		return nil, err
	}
	for i := range items {
		db.Where("item_metadata.id = ? ", items[i].ItemMetadataItem.ItemMetadataID).Find(&items[i].ItemMetadata)
	}
	log.Printf("用户id %d 查询聚类 %s 的照片成功", userID, clusterName)
	return items, nil
}
