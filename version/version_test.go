package version

import (
	"testing"
)

var (
	v1   = "2.0.1"
	v1_1 = "2.0.1.1"
	v2   = "2.1.1"
)

// 该函数比较两个版本号是否相等，是否大于或小于的关系
// 返回值：0表示v1与v2相等；1表示v1大于v2；2表示v1小于v2
func TestCompare(t *testing.T) {
	if 0 != Compare(v1, v1) {
		t.Errorf("v1 %s == v1 %s ", v1, v1)
	}

	if 1 != Compare(v2, v1) {
		t.Errorf("v2 %s > v1 %s", v2, v1)
	}

	if 2 != Compare(v1, v2) {
		t.Errorf("v1 %s < v2 %s", v1, v2)
	}

	// 测试带V和v前缀的版本号
	if 0 != Compare("V2.0.1", "v2.0.1") {
		t.Errorf("V2.0.1 should equal v2.0.1")
	}

	// 测试带-分隔符的版本号
	if 0 != Compare("2-0-1", "2.0.1") {
		t.Errorf("2-0-1 should equal 2.0.1")
	}

	// 测试不同长度的版本号
	if 0 != Compare("2.1.1.0", "2.1.1") {
		t.Errorf("2.1.1.0 should greater than 2.1.1")
	}

	if 0 != Compare("2.1.1", "2.1.1.0") {
		t.Errorf("2.1.1 should less than 2.1.1.0")
	}

	// 测试大版本号比较
	if 1 != Compare("3.0.0", "2.9.9") {
		t.Errorf("3.0.0 should greater than 2.9.9")
	}

	// 测试小版本号比较
	if 2 != Compare("2.0.9", "2.1.0") {
		t.Errorf("2.0.9 should less than 2.1.0")
	}

	// 测试补丁版本号比较
	if 1 != Compare("2.0.2", "2.0.1") {
		t.Errorf("2.0.2 should greater than 2.0.1")
	}

}

func TestVersionCompare(t *testing.T) {
	if VersionCompare("0.0.12", "0.0.11", "<") {
		t.Errorf("v1 %s < v1_1 %s ", v1, v1_1)
	}

	if VersionCompare("0.0.44", "0.0.7", "<") {
		t.Errorf("v1 %s < v1_1 %s ", v1, v1_1)
	}

	if !VersionCompare(v1, v1_1, "<") {
		t.Errorf("v1 %s < v1_1 %s ", v1, v1_1)
	}

	if !VersionCompare(v2, v1_1, ">") {
		t.Errorf("v2 %s < v1_1 %s ", v2, v1_1)
	}

	if !VersionCompare(v2, v1, ">=") {
		t.Errorf("v2 %s >= v1 %s ", v2, v1)
	}

	if !VersionCompare(v1, v2, "<=") {
		t.Errorf("v1 %s >= v2 %s ", v1, v2)
	}

	if !VersionCompare(v1, v1, "==") {
		t.Errorf("v1 %s == v1 %s ", v1, v1)
	}

	// 测试带有前缀的版本号比较
	if !VersionCompare("v1.2.3", "1.2.3", "==") {
		t.Errorf("v1.2.3 应该等于 1.2.3")
	}

	// 测试带有连字符的版本号比较
	if !VersionCompare("1.2-3", "1.2.3", "==") {
		t.Errorf("1.2-3 应该等于 1.2.3")
	}

	// 测试不同长度的版本号比较
	if !VersionCompare("2.0", "2.0.0", "==") {
		t.Errorf("2.0 应该等于 2.0.0")
	}

	// 测试大版本号差异比较
	if !VersionCompare("10.0.0", "2.0.0", ">") {
		t.Errorf("10.0.0 应该大于 2.0.0")
	}

	// 测试复杂版本号比较
	if !VersionCompare("v2.1-5", "2.1.5", "==") {
		t.Errorf("v2.1-5 应该等于 2.1.5")
	}

	// 测试边界情况
	if VersionCompare("1.0.0", "1.0.0", "<") {
		t.Errorf("1.0.0 不应该小于 1.0.0")
	}

	if VersionCompare("1.0.0", "1.0.0", ">") {
		t.Errorf("1.0.0 不应该大于 1.0.0")
	}

}
