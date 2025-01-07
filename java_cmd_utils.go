package pom_component_parsing

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetJavaHome 获取系统中的 JAVA_HOME 路径
func GetJavaHome() string {
	// 首先检查环境变量
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		return javaHome
	}

	// 根据不同操作系统查找 Java 安装路径
	switch runtime.GOOS {
	case "windows":
		return findJavaHomeWindows()
	case "darwin":
		return findJavaHomeMac()
	default: // linux 和其他类 Unix 系统
		return findJavaHomeLinux()
	}
}

// findJavaHomeWindows 在 Windows 系统中查找 JAVA_HOME
func findJavaHomeWindows() string {
	// 常见的 Windows Java 安装路径
	commonPaths := []string{
		`C:\Program Files\Java`,
		`C:\Program Files (x86)\Java`,
	}

	// 尝试使用 where java 命令
	cmd := exec.Command("where", "java")
	output, err := cmd.Output()
	if err == nil {
		paths := strings.Split(string(output), "\n")
		if len(paths) > 0 {
			// 通常 java.exe 在 bin 目录下，我们需要其父目录的父目录
			javaPath := strings.TrimSpace(paths[0])
			if binDir := filepath.Dir(javaPath); binDir != "" {
				if javaHome := filepath.Dir(binDir); javaHome != "" {
					return javaHome
				}
			}
		}
	}

	// 搜索常见路径
	for _, basePath := range commonPaths {
		if javaHome := searchJavaInDir(basePath); javaHome != "" {
			return javaHome
		}
	}

	return ""
}

// findJavaHomeMac 在 macOS 系统中查找 JAVA_HOME
func findJavaHomeMac() string {
	// 使用 /usr/libexec/java_home 命令
	cmd := exec.Command("/usr/libexec/java_home")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	// 检查常见的 macOS Java 安装路径
	commonPaths := []string{
		"/Library/Java/JavaVirtualMachines",
		"/System/Library/Java/JavaVirtualMachines",
	}

	for _, path := range commonPaths {
		if javaHome := searchJavaInDir(path); javaHome != "" {
			return javaHome
		}
	}

	return ""
}

// findJavaHomeLinux 在 Linux 系统中查找 JAVA_HOME
func findJavaHomeLinux() string {
	// 使用 which java 命令
	cmd := exec.Command("which", "java")
	output, err := cmd.Output()
	if err == nil {
		// 解析 java 命令的符号链接
		javaPath := strings.TrimSpace(string(output))
		realPath, err := filepath.EvalSymlinks(javaPath)
		if err == nil {
			// 通常 java 在 bin 目录下，我们需要其父目录的父目录
			if binDir := filepath.Dir(realPath); binDir != "" {
				if javaHome := filepath.Dir(binDir); javaHome != "" {
					return javaHome
				}
			}
		}
	}

	// 检查常见的 Linux Java 安装路径
	commonPaths := []string{
		"/usr/lib/jvm",
		"/usr/java",
		"/usr/local/java",
		"/opt/java",
	}

	for _, path := range commonPaths {
		if javaHome := searchJavaInDir(path); javaHome != "" {
			return javaHome
		}
	}

	return ""
}

// searchJavaInDir 在指定目录中搜索 Java 安装
func searchJavaInDir(baseDir string) string {
	// 检查目录是否存在
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return ""
	}

	// 遍历目录查找 java 可执行文件
	var newestJavaHome string
	var newestVersion string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		// 检查是否为 bin 目录且包含 java 可执行文件
		if info.IsDir() && strings.HasSuffix(path, "bin") {
			javaExe := "java"
			if runtime.GOOS == "windows" {
				javaExe = "java.exe"
			}

			if _, err := os.Stat(filepath.Join(path, javaExe)); err == nil {
				possibleHome := filepath.Dir(path)
				version := getJavaVersion(filepath.Join(path, javaExe))

				// 如果找到更新的版本，更新结果
				if version > newestVersion {
					newestVersion = version
					newestJavaHome = possibleHome
				}
			}
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return newestJavaHome
}

// getJavaVersion 获取 Java 版本信息
func getJavaVersion(javaPath string) string {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// 查找版本字符串
	outputStr := string(output)
	if strings.Contains(outputStr, "version") {
		return outputStr
	}
	return ""
}
