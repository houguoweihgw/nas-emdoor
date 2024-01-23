package tools

import (
	"bytes"
	"encoding/gob"
	"math"
)

// SerializeEmbeddings 序列化 Embeddings 切片为字节数组
func SerializeEmbeddings(embeddings []float64) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(embeddings)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DeserializeEmbeddings 从字节数组反序列化 Embeddings 切片
func DeserializeEmbeddings(data []byte) ([]float64, error) {
	var embeddings []float64
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&embeddings)
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

// CalculateEuclideanDistance 计算两个向量之间的欧氏距离
func CalculateEuclideanDistance(embedding1, embedding2 []byte) float64 {
	em1, _ := DeserializeEmbeddings(embedding1)
	em2, _ := DeserializeEmbeddings(embedding2)
	// 确保向量长度相同
	if len(em1) != len(em2) {
		return 0.0
	}
	// 计算欧氏距离
	distance := 0.0
	for i := 0; i < len(em1); i++ {
		diff := em1[i] - em2[i]
		distance += diff * diff
	}
	distance = math.Sqrt(distance)
	return distance
}

// EuclideanDistance 计算两个特征向量之间的欧氏距离
func EuclideanDistance(e1, e2 []float64) float64 {
	// 确保两个向量具有相同的长度
	if len(e1) != len(e2) {
		return 0.0 // 或者其他适当的错误处理方式
	}

	// 计算欧氏距离
	sumOfSquares := 0.0
	for i := 0; i < len(e1); i++ {
		diff := e1[i] - e2[i]
		sumOfSquares += diff * diff
	}

	distance := math.Sqrt(sumOfSquares)
	return distance
}
