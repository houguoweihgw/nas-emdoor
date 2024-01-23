package utils

import "testing"

func TestDeleteKeysByPattern(t *testing.T) {
	err := InitRedisClient()
	if err != nil {
		return
	}
	err = DeleteKeysByPattern(RDB, "user_photos")
	if err != nil {
		return
	}
}
