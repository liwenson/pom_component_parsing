package pom_component_parsing

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// MvnCommandInfo 存储Maven命令的相关信息
type MvnCommandInfo struct {
	Path             string `json:"path"`               // Maven命令的路径
	MvnVersion       string `json:"mvn_version"`        // Maven的版本
	UserSettingsPath string `json:"user_settings_path"` // 用户自定义的settings.xml路径
	JavaHome         string `json:"java_home"`          // Java的安装目录
}

// String 实现了fmt.Stringer接口，便于打印MvnCommandInfo的信息
func (m MvnCommandInfo) String() string {
	return fmt.Sprintf("MavenCommand: %s, JavaHome: %s, MavenVersion: %s, UserSettings: %s", m.Path, m.JavaHome, m.MvnVersion, m.UserSettingsPath)
}

// Command 根据传入的参数构建一个exec.Cmd对象，用于执行Maven命令
func (m MvnCommandInfo) Command(args ...string) *exec.Cmd {
	// 构建Maven命令的参数列表
	var _args = make([]string, 0, len(args)+5)
	if m.UserSettingsPath != "" {
		_args = append(_args, "--settings", m.UserSettingsPath)
	}
	_args = append(_args, "--batch-mode") // 添加批处理模式参数
	_args = append(_args, args...)

	// 创建命令
	cmd := exec.Command(m.Path, _args...)

	// 设置JAVA_HOME环境变量，如果有指定
	if m.JavaHome != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "JAVA_HOME="+m.JavaHome)
	}

	return cmd
}

// _MvnCommandResult 用于缓存Maven命令检查的结果
type _MvnCommandResult struct {
	rs *MvnCommandInfo // Maven命令信息
	e  error           // 可能的错误
}

var (
	cachedMvnCommandResult *_MvnCommandResult // 缓存的Maven命令结果
	mu                     sync.RWMutex       // 互斥锁，确保并发安全
)

// 定义错误变量
var (
	// ErrMvnNotFound 表示未找到Maven命令
	ErrMvnNotFound = errors.New("Maven command not found")

	// ErrCheckMvnVersion 表示检查Maven版本时发生的错误
	ErrCheckMvnVersion = errors.New("failed to check Maven version")
)

// CheckMvnCommand 检查并返回Maven命令的信息，如果之前已经检查过，则返回缓存的结果
func CheckMvnCommand() (info *MvnCommandInfo, err error) {
	mu.RLock()
	if cachedMvnCommandResult != nil {
		if cachedMvnCommandResult.e != nil {
			// 可以在此处添加日志记录错误信息
		}
		if cachedMvnCommandResult.rs != nil {
			// 可以在此处添加日志记录使用缓存信息
			fmt.Println("使用缓存的Maven命令信息")
		}
		info, err = cachedMvnCommandResult.rs, cachedMvnCommandResult.e
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// 再次检查缓存，避免重复初始化
	if cachedMvnCommandResult != nil {
		info, err = cachedMvnCommandResult.rs, cachedMvnCommandResult.e
		return
	}

	// 初始化Maven命令信息
	info = &MvnCommandInfo{
		Path:             "",
		JavaHome:         "",
		UserSettingsPath: "",
	}

	if info.Path == "" {
		info.Path = getMvnCommandOs()
	}
	if info.Path == "" {
		err = ErrMvnNotFound
		cachedMvnCommandResult = &_MvnCommandResult{rs: nil, e: err}
		return
	}

	// 检查Maven版本
	ver, e := checkMvnVersion(info.Path, info.JavaHome)
	if e != nil {
		err = e
		cachedMvnCommandResult = &_MvnCommandResult{rs: info, e: err}
		return
	}
	info.MvnVersion = ver

	// 缓存结果
	cachedMvnCommandResult = &_MvnCommandResult{
		rs: info,
		e:  nil,
	}
	return
}

// locateMvnCmdPath 定位Maven命令的路径
func locateMvnCmdPath() string {
	return getMvnCommandOs()
}

// executeMvnVersion 执行Maven命令以获取其版本信息
func executeMvnVersion(mvnPath string, javaHome string) (string, error) {
	// 构建执行命令
	cmd := exec.Command(mvnPath, "--version", "--batch-mode")

	// 设置环境变量
	cmd.Env = os.Environ()
	if javaHome != "" {
		cmd.Env = append(cmd.Env, "JAVA_HOME="+javaHome)
	}

	// 设置超时
	var output []byte
	done := make(chan error, 1)
	go func() {
		var err error
		output, err = cmd.Output()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrCheckMvnVersion, err)
		}
	case <-time.After(8 * time.Second):
		// 超时处理
		if err := cmd.Process.Kill(); err != nil {
			return "", fmt.Errorf("无法杀死Maven进程: %w", err)
		}
		return "", fmt.Errorf("%w: 执行Maven版本命令超时", ErrCheckMvnVersion)
	}

	return string(output), nil
}

// checkMvnVersion 检查Maven的版本，并返回版本号
func checkMvnVersion(mvnPath string, javaHome string) (string, error) {
	output, err := executeMvnVersion(mvnPath, javaHome)
	if err != nil {
		// 针对Linux和Darwin系统，尝试授予可执行权限后重试
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			_ = os.Chmod(mvnPath, 0755) // 忽略错误
			output, err = executeMvnVersion(mvnPath, javaHome)
		}
		if err != nil {
			return "", err
		}
	}

	// 解析Maven版本
	ver := parseMvnVersion(output)
	if ver == "" {
		return "", fmt.Errorf("%w: 无法解析Maven版本", ErrCheckMvnVersion)
	}
	return ver, nil
}

// parseMvnVersion 解析Maven命令输出，提取版本号
func parseMvnVersion(input string) string {
	// 定义正则表达式匹配版本号，例如：Apache Maven 3.6.3
	versionPattern := regexp.MustCompile(`Apache Maven (\d+(?:\.[\dA-Za-z_-]+)+)`)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if m := versionPattern.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}

// getMvnCommandOs 根据操作系统查找Maven命令的路径
func getMvnCommandOs() string {
	p, err := exec.LookPath("mvn")
	if err != nil {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	absPath, err := filepath.Abs(p)
	if err == nil {
		return absPath
	}
	return ""
}

// 初始化一次性设置，确保缓存的安全性
func init() {
	cachedMvnCommandResult = nil
}
