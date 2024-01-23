package entity

import (
	"github.com/jinzhu/gorm"
	"log"
	"nas-backend/tools"
)

type Face struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	ItemID   uint   `json:"item_id"`
	UserID   int    `json:"user_id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	W        int    `json:"w"`
	H        int    `json:"h"`
	Features []byte `json:"features"`
}

// TableName 设置表名
func (Face) TableName() string {
	return "faces"
}

// CreateFaces 插入faces表
func CreateFaces(db *gorm.DB, faceInfo FaceInfo, userID int, itemID int) (Face, error) {
	feature, _ := tools.SerializeEmbeddings(faceInfo.Embeddings)
	newFace := Face{
		ItemID:   uint(itemID),
		UserID:   userID,
		X:        faceInfo.X,
		Y:        faceInfo.Y,
		W:        faceInfo.Width,
		H:        faceInfo.Height,
		Features: feature,
	}
	result := db.Create(&newFace)
	if result.Error != nil {
		log.Printf("ItemID:%d 插入faces表失败\n", itemID)
		return newFace, result.Error
	}
	log.Printf("ItemID:%d 插入faces表成功:id %d\n", itemID, newFace.ID)
	return newFace, nil
}

// GetCoversForCluster 查询人脸聚类的第一张人脸照片作为封面
func GetCoversForCluster(db *gorm.DB, clusterID int) (Face, error) {
	var clusterFace ClusterFace
	var face Face
	if err := db.Model(&ClusterFace{}).
		Where("cluster_id = ?", clusterID).
		First(&clusterFace).Error; err != nil {
		log.Printf("查询ClusterFace映射失败 %d\n", clusterID)
		return face, err
	}
	if err := db.Model(&Face{}).
		Where("id = ?", clusterFace.FaceID).
		First(&face).Error; err != nil {
		log.Printf("查询Face失败 %d\n", clusterFace.FaceID)
		return face, err
	}
	log.Printf("查询人脸聚类 %d 的封面成功成功 %d\n", clusterID, face.ItemID)
	return face, nil
}
