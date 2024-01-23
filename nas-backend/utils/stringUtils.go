package utils

import (
	"strings"
)

// ConvertTagsToString 将 []string 转换成逗号分隔的字符串
func ConvertTagsToString(tags []string) string {
	return strings.Join(tags, ",")
}
