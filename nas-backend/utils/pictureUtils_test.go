package utils

import (
	"log"
	"nas-backend/entity"
	"testing"
)

func TestParsePictureExif(t *testing.T) {
	var itemMeta entity.ItemMetadata
	err := ParsePictureExif("/home/houguowei/GolandProjects/nas-backend/data/IMG_5795.JPG", &itemMeta)
	if err != nil {
		log.Println(itemMeta)
		return
	}
}

func TestParLatitudeAndLongitudeTo(t *testing.T) {
	ParLatitudeAndLongitudeTo()
}
