package pom_component_parsing

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestGetJavaHome 测试 GetJavaHome 函数
// 该测试用例会检查在无任何环境变量干扰的情况下，GetJavaHome 是否能够正常返回一个非空值
// 请注意：由于函数依赖系统环境与安装情况，若测试环境中未安装 Java 或路径无法找到，可能会导致测试失败。
func TestGetJavaHome(t *testing.T) {
	home := GetJavaHome()
	if home == "" {
		t.Errorf("期望获取到一个非空的 JAVA_HOME 路径，但实际获取到的是空值")
	} else {
		t.Logf("检测到的 JAVA_HOME 路径: %s", home)
	}
}

// TestGetJavaHomeWithEnvVar 测试在设置 JAVA_HOME 环境变量时，GetJavaHome 是否优先返回该环境变量的值
func TestGetJavaHomeWithEnvVar(t *testing.T) {
	// 先保存原来的环境变量
	originalEnv := os.Getenv("JAVA_HOME")
	defer os.Setenv("JAVA_HOME", originalEnv)

	// 设置一个假的 JAVA_HOME
	mockHome := "/mock/java/home"
	os.Setenv("JAVA_HOME", mockHome)

	home := GetJavaHome()
	if home != mockHome {
		t.Errorf("当环境变量 JAVA_HOME=%s 时，期望返回 %s，实际返回 %s", mockHome, mockHome, home)
	}
}

// TestFindJavaHomeWindows 测试在 Windows 平台上查找 JAVA_HOME 的函数
// 注意：该测试仅在 Windows 系统上有效，其他平台会跳过测试
func TestFindJavaHomeWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("仅在 Windows 平台执行此测试")
	}
	javaHome := findJavaHomeWindows()
	if javaHome == "" {
		t.Log("在当前 Windows 环境中未能找到 JAVA_HOME，可能是未安装 Java 或安装路径不在常规位置")
	} else {
		t.Logf("在 Windows 环境中找到的 JAVA_HOME: %s", javaHome)
	}
}

// TestFindJavaHomeMac 测试在 macOS 平台上查找 JAVA_HOME 的函数
// 注意：该测试仅在 macOS 系统上有效，其他平台会跳过测试
func TestFindJavaHomeMac(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("仅在 macOS 平台执行此测试")
	}
	javaHome := findJavaHomeMac()
	if javaHome == "" {
		t.Log("在当前 macOS 环境中未能找到 JAVA_HOME，可能是未安装 Java 或安装路径不在常规位置")
	} else {
		t.Logf("在 macOS 环境中找到的 JAVA_HOME: %s", javaHome)
	}
}

// TestFindJavaHomeLinux 测试在 Linux 平台上查找 JAVA_HOME 的函数
// 注意：该测试仅在 Linux 系统上有效，其他平台会跳过测试
func TestFindJavaHomeLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("仅在 Linux 平台执行此测试")
	}
	javaHome := findJavaHomeLinux()
	if javaHome == "" {
		t.Log("在当前 Linux 环境中未能找到 JAVA_HOME，可能是未安装 Java 或安装路径不在常规位置")
	} else {
		t.Logf("在 Linux 环境中找到的 JAVA_HOME: %s", javaHome)
	}
}

// TestSearchJavaInDir 测试 searchJavaInDir 函数  耗时太久
// 注意：此函数会遍历给定路径下的所有子目录，可能导致测试耗时较长，或在无权限访问的目录出错
// 示例中仅演示对系统根目录等进行查询往往不合适，可自行改为更小且可控的测试目录
//func TestSearchJavaInDir(t *testing.T) {
//	// 可根据测试需要修改为更小的目录，避免遍历过多文件
//	testDir := "/"
//	if runtime.GOOS == "windows" {
//		testDir = "C:\\"
//	}
//
//	javaHome := searchJavaInDir(testDir)
//	// 若 testDir 较大，这里可能会很耗时或没有找到Java安装，测试时请谨慎选择目录
//	if javaHome == "" {
//		t.Logf("在目录 %s 中未能找到 Java 安装", testDir)
//	} else {
//		t.Logf("在目录 %s 中找到的 Java 安装路径: %s", testDir, javaHome)
//	}
//}

// TestGetJavaVersion 测试 getJavaVersion 函数
// 需要确保本机存在 java 或者指定一个可执行的 java 路径
func TestGetJavaVersion(t *testing.T) {
	// 尝试使用系统 PATH 中的 java，可自行修改为绝对路径进行测试
	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	// 直接使用可执行文件名(依赖系统 PATH)，若无法执行会返回空字符串
	version := getJavaVersion(javaExe)
	if version == "" {
		t.Logf("无法通过可执行文件 %s 获取到 Java 版本信息，可能是系统未安装 Java 或未在 PATH 中", javaExe)
	} else {
		t.Logf("获取到的 Java 版本信息: %s", version)
	}
}

// TestGetJavaVersionWithInvalidPath 测试当给定一个无效的 java 路径时，getJavaVersion 应返回空字符串
func TestGetJavaVersionWithInvalidPath(t *testing.T) {
	invalidPath := filepath.Join(os.TempDir(), "not_exist_java_executable")
	version := getJavaVersion(invalidPath)
	if version != "" {
		t.Errorf("期望在无效路径 %s 上获取到空字符串，但实际返回: %s", invalidPath, version)
	} else {
		t.Log("在无效路径上成功返回了空字符串")
	}
}
