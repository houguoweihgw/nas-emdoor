package utils

import (
	"fmt"
	"github.com/jasonwinn/geocoder"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"nas-backend/entity"
	"os"
	"strconv"
	"time"
)

func ParsePictureExif(filePath string, itemMeta *entity.ItemMetadata) error {
	//todo:考虑处理解析失败的情况，如无法解析 Exif 数据或获取元数据字段时。
	// 你可以根据你的应用需求，决定是否记录错误、返回错误信息、跳过处理，或采取其他适当的操作。
	//打开图像文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("解析照片打开失败:%v", err)
		return err
	}
	defer file.Close()
	// 解析 Exif 数据
	x, err := exif.Decode(file)
	if err != nil {
		log.Printf("解析照片Exif数据失败:%v", err)
		return err
	}
	// 获取文件信息以获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// 获取文件大小（以字节为单位）
	itemMeta.FileSize = fileInfo.Size()
	// 获取图像宽度和高度
	if value, err := x.Get(exif.PixelXDimension); err == nil {
		itemMeta.ImageWidth, _ = value.Int(0)
	}
	if value, err := x.Get(exif.PixelYDimension); err == nil {
		itemMeta.ImageLength, _ = value.Int(0)
	}
	// 获取相机品牌和型号
	// 获取相机品牌和型号
	makeValue, _ := x.Get(exif.Make)
	modelValue, _ := x.Get(exif.Model)
	if makeStr, err := makeValue.StringVal(); err == nil {
		itemMeta.Make = makeStr
	}
	if modelStr, err := modelValue.StringVal(); err == nil {
		itemMeta.Model = modelStr
	}
	// 获取曝光时间
	if value, err := x.Get(exif.ExposureTime); err == nil {
		numerator, denominator, _ := value.Rat2(0)
		exposureTime := strconv.FormatInt(numerator, 10) + "/" + strconv.FormatInt(denominator, 10) + "s"
		itemMeta.ExposureTime = exposureTime
	}
	// 获取光圈值
	if value, err := x.Get(exif.FNumber); err == nil {
		numerator, denominator, _ := value.Rat2(0)
		itemMeta.Aperture = float64(numerator) / float64(denominator)
	}
	// 获取ISO值
	isoValue, _ := x.Get(exif.ISOSpeedRatings)
	isoInt, err := isoValue.Int(0)
	if err == nil {
		itemMeta.ISO = isoInt
	}
	// 获取焦距
	if value, err := x.Get(exif.FocalLength); err == nil {
		numerator, denominator, _ := value.Rat2(0)
		itemMeta.FocalLength = float64(numerator) / float64(denominator)
	}
	// 获取拍摄日期和时间
	dateTaken, err := x.Get(exif.DateTime)
	date, err := dateTaken.StringVal()
	if err == nil {
		parsedDate, parseErr := time.Parse("2006:01:02 15:04:05", date)
		//fmt.Println(parsedDate)
		if parseErr == nil {
			itemMeta.DateTaken = parsedDate
		}
	}
	// 获取GPS位置（纬度、经度）北纬为正，东经为正
	latRef, err := x.Get(exif.GPSLatitudeRef)
	if err != nil {
		log.Printf("获取GPSLatitudeRef失败：%v", err)
		// 处理错误，可以返回错误信息或采取其他适当的操作
		return err
	}
	latRefVal, _ := latRef.StringVal()

	LonRef, err := x.Get(exif.GPSLongitudeRef)
	if err != nil {
		log.Printf("获取GPSLongitudeRef失败：%v", err)
		// 处理错误，可以返回错误信息或采取其他适当的操作
		return err
	}
	LonRefVal, _ := LonRef.StringVal()
	//fmt.Println(LonRefVal)
	if lat, long, err := x.LatLong(); err == nil {
		if latRefVal != "N" {
			itemMeta.Latitude = 0 - lat
		} else {
			itemMeta.Latitude = lat
		}
		if LonRefVal != "E" {
			itemMeta.Longitude = 0 - long
		} else {
			itemMeta.Longitude = long
		}

	}
	// 获取GPS高度
	if value, err := x.Get(exif.GPSAltitude); err == nil {
		numerator, denominator, _ := value.Rat2(0)
		itemMeta.Altitude = float64(numerator) / float64(denominator)
	}
	log.Printf("完成照片的元数据解析")
	return nil
}

func ParLatitudeAndLongitudeTo() {
	geocoder.SetAPIKey("cGcnrw9v3Er2RfWPraW67tsi0QG4RmLs")
	address, err := geocoder.ReverseGeocode(47.6064, -122.330803)
	if err != nil {
		panic("THERE WAS SOME ERROR!!!!!")
	}
	fmt.Println(address)
}
