package utils

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, _ := GenerateJWT("admin", "123456")
	fmt.Println(token)
}

func TestParseJTW(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5MjI0NzcsInVzZXItbmFtZSI6ImFkbWluIiwidXNlci1wYXNzd29yZCI6IjEyMzQ1NiJ9.5LSMnvQrd1yqCo1zIzREi4m81144qAIWZlUG21NMgTo"
	claims, _ := ParseJTW(tokenString)
	fmt.Println(claims)
}
