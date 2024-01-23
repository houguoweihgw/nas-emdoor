package tools

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

func PhotoCropping(photoURL string, x int, y int, w int, h int) ([]byte, error) {
	// 打开原始图片
	originalImageFile, err := os.Open(photoURL) // 替换为你的原始图片路径
	if err != nil {
		panic(err)
	}
	defer originalImageFile.Close()
	// 解码原始图片
	originalImage, _, err := image.Decode(originalImageFile)
	if err != nil {
		panic(err)
	}
	// 定义裁剪区域
	cropRect := image.Rect(x, y, w, h) // 以左上角为(100, 100)，右下角为(300, 300)的区域
	// 创建一个与裁剪区域大小相同的画布
	croppedImage := image.NewRGBA(cropRect)
	// 使用 draw 包的 Draw 函数裁剪图片
	draw.Draw(croppedImage, cropRect, originalImage, image.Point{x, y}, draw.Src)
	// 创建一个字节数组缓冲区
	var buffer bytes.Buffer
	// 编码并将裁剪后的图片保存到缓冲区
	err = png.Encode(&buffer, croppedImage)
	if err != nil {
		return nil, err
	}
	// 返回字节数组
	return buffer.Bytes(), nil

	// 创建输出图片文件
	//outputImageFile, err := os.Create("output.JPEG") // 替换为你的输出图片路径
	//if err != nil {
	//	panic(err)
	//}
	//defer outputImageFile.Close()
	//// 编码并保存裁剪后的图片
	//err = png.Encode(outputImageFile, croppedImage)
	//if err != nil {
	//	panic(err)
	//}
}
