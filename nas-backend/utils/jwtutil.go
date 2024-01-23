package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// GenerateJWT 生成 JWT 令牌
func GenerateJWT(userName string, userPassword string) (string, error) {
	mySigningKey := []byte("asfasfdafasdfdasfa.")
	// 创建 JWT 负载
	payload := jwt.MapClaims{
		"user-name":     userName,     // 用户ID
		"user-password": userPassword, // 用户名
		"exp":           time.Now().Unix() + 60*60*24,
		// 其他声明...
	}
	// 使用密钥生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 签名 JWT
	jwtString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	//fmt.Println("加密后的token字符串", jwtString)
	log.Println("JWT生成完成")
	return jwtString, nil
	//return "nil", nil
}

// ParseJTW 解析 JWT 令牌
func ParseJTW(tokenString string) (jwt.MapClaims, error) {
	mySigningKey := []byte("asfasfdafasdfdasfa.")
	// 解析 JWT 令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 检查是否有效
	if !token.Valid {
		return nil, fmt.Errorf("Invalid JWT token")
	}
	// 获取负载（claims）
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid JWT claims")
}
