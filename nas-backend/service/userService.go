package service

//接收登录请求：后端应该暴露一个接口，用于接收前端发送的登录请求。通常，这个接口是一个 POST 请求，包含用户提供的用户名和密码。
//验证用户身份：在接收到登录请求后，后端需要验证用户的身份。这通常涉及检查用户提供的用户名和密码是否匹配数据库中的记录。您可以使用数据库查询来执行此操作。
//生成身份认证令牌JWT：这个令牌将用于标识用户并在后续请求中进行身份验证。
//返回令牌：后端应该将生成的令牌作为响应的一部分返回给前端。通常，令牌会在响应的 JSON 数据中返回。
//保护受限资源：后端应该在需要身份验证的受限资源上实施访问控制。这意味着只有具有有效令牌的用户才能访问这些资源。

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"nas-backend/entity"
	"nas-backend/utils"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	// 处理登录请求
	// 解析 JSON 请求正文
	var user entity.User
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "请求数据错误",
		})
		return
	}
	log.Printf("username: %s ,password: %s  \n", user.Username, user.Password)
	// 验证用户身份
	if !entity.ValidLoginInfo(utils.DB, user) {
		log.Printf("Login Password Wrong \n")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "密码错误",
		})
		return
	}
	// 用户登录成功，生成 JWT 令牌
	userName := user.Username     // 假设这是用户的唯一标识
	userPassword := user.Password // 假设这是用户的用户名
	jwtToken, err := utils.GenerateJWT(userName, userPassword)
	if err != nil {
		// 处理生成令牌失败的情况
		log.Printf("Failed to generate JWT token: %v \n", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "服务器内部错误",
		})
		return
	}
	// 响应成功
	c.Writer.WriteHeader(http.StatusOK)
	// 创建一个包含令牌和用户名的 JSON 对象
	responseData := map[string]interface{}{
		"token":    jwtToken,
		"username": user.Username,
		"message":  "登陆成功",
	}
	log.Printf("用户 %s 登陆成功！\n", userName)
	// 将整个对象写入响应
	c.JSON(http.StatusOK, responseData)
}

func RegisterHandler(c *gin.Context) {
	// 处理注册请求
	// 解析 JSON 请求正文
	var user entity.User
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求数据错误",
		})
	}
	log.Printf("用户注册 username: %s ,password: %s ,email: %s  \n", user.Username, user.Password, user.Email)
	err2 := entity.UserRegister(utils.DB, user)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户注册失败",
		})
	}
	// 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}
