package version

import (
	"strconv"
	"strings"
)

// 该函数比较两个版本号是否相等，是否大于或小于的关系
// 返回值：0表示v1与v2相等；1表示v1大于v2；2表示v1小于v2
func Compare(v1, v2 string) int {
	// 替换一些常见的版本符号
	replaceMap := map[string]string{"V": "", "v": "", "-": "."}
	for k, v := range replaceMap {
		v1 = strings.ReplaceAll(v1, k, v)
		v2 = strings.ReplaceAll(v2, k, v)
	}

	ver1 := strings.Split(v1, ".")
	ver2 := strings.Split(v2, ".")

	// 找出v1和v2哪一个最长
	maxLen := len(ver1)
	if len(ver2) > maxLen {
		maxLen = len(ver2)
	}

	// 循环比较
	for i := 0; i < maxLen; i++ {
		var num1, num2 int
		if i < len(ver1) {
			num1, _ = strconv.Atoi(ver1[i])
		}
		if i < len(ver2) {
			num2, _ = strconv.Atoi(ver2[i])
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return 2
		}
	}
	return 0
}

func VersionCompare(v1, v2, operator string) bool {
	com := Compare(v1, v2)
	switch operator {
	case "==":
		return com == 0
	case "<":
		return com == 2
	case ">":
		return com == 1
	case "<=":
		return com == 0 || com == 2
	case ">=":
		return com == 0 || com == 1
	}
	return false
}
