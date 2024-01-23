package tools

import (
	"fmt"
	"testing"
)

func TestDeserializeEmbeddings(t *testing.T) {
	// 定义两个向量
	e1 := []float64{1.0, 2.0, 3.0}
	e2 := []float64{4.0, 5.0, 6.0}

	// 计算欧氏距离
	distance := EuclideanDistance(e1, e2)

	// 打印结果
	fmt.Println("欧氏距离:", distance)
}
