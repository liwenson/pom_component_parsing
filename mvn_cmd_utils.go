// Package pom_component_parsing 提供了 Maven 命令的解析和执行功能
// 该包主要用于检测系统中的 Maven 配置，验证其可用性，并提供命令执行的封装
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

// MvnCommandInfo 存储 Maven 命令的相关配置信息
// 包含了执行 Maven 命令所需的所有必要参数
type MvnCommandInfo struct {
	Path             string `json:"path"`               // Maven 可执行文件的完整路径
	MvnVersion       string `json:"mvn_version"`        // Maven 的版本号（如 3.6.3）
	UserSettingsPath string `json:"user_settings_path"` // Maven 用户配置文件 settings.xml 的路径
	JavaHome         string `json:"java_home"`          // Java 安装目录的路径
}

// String 方法实现了 fmt.Stringer 接口，用于格式化输出 MvnCommandInfo 的信息
// 主要用于调试和日志记录
func (m MvnCommandInfo) String() string {
	return fmt.Sprintf("MavenCommand: %s, JavaHome: %s, MavenVersion: %s, UserSettings: %s",
		m.Path, m.JavaHome, m.MvnVersion, m.UserSettingsPath)
}

// Command 根据传入的参数构建一个可执行的 Maven 命令
// 返回配置好的 exec.Cmd 对象，可直接执行
func (m MvnCommandInfo) Command(args ...string) *exec.Cmd {
	// 预分配足够的切片容量，避免多次扩容
	var cmdArgs = make([]string, 0, len(args)+5)

	// 如果指定了用户配置文件，添加相应参数
	if m.UserSettingsPath != "" {
		cmdArgs = append(cmdArgs, "--settings", m.UserSettingsPath)
	}

	// 添加批处理模式参数，禁用交互式输出
	cmdArgs = append(cmdArgs, "--batch-mode")

	// 添加用户传入的其他参数
	cmdArgs = append(cmdArgs, args...)

	// 创建命令对象
	cmd := exec.Command(m.Path, cmdArgs...)

	// 如果指定了 JAVA_HOME，设置到环境变量中
	if m.JavaHome != "" {
		// 复制当前的环境变量
		cmd.Env = os.Environ()
		// 添加或覆盖 JAVA_HOME 环境变量
		cmd.Env = append(cmd.Env, "JAVA_HOME="+m.JavaHome)
	}

	return cmd
}

// _MvnCommandResult 用于缓存 Maven 命令检查的结果
// 包含了检查结果和可能出现的错误
type _MvnCommandResult struct {
	rs *MvnCommandInfo // Maven 命令的检查结果
	e  error           // 检查过程中可能出现的错误
}

// 全局变量定义
var (
	// 用于缓存 Maven 命令检查的结果，避免重复检查
	cachedMvnCommandResult *_MvnCommandResult
	// 互斥锁，用于保护缓存的并发访问
	mu sync.RWMutex
)

// 错误定义
var (
	// ErrMvnNotFound 表示系统中未找到 Maven 命令
	ErrMvnNotFound = errors.New("Maven command not found")
	// ErrCheckMvnVersion 表示检查 Maven 版本时发生错误
	ErrCheckMvnVersion = errors.New("failed to check Maven version")
)

// CheckMvnCommand 检查并返回系统中的 Maven 命令信息
// 该函数会缓存检查结果，避免重复执行耗时的检查操作
func CheckMvnCommand() (info *MvnCommandInfo, err error) {
	// 尝试从缓存中读取结果
	mu.RLock()
	if cachedMvnCommandResult != nil {
		info, err = cachedMvnCommandResult.rs, cachedMvnCommandResult.e
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	// 如果缓存中没有结果，获取写锁进行检查
	mu.Lock()
	defer mu.Unlock()

	// 双重检查，避免并发情况下的重复初始化
	if cachedMvnCommandResult != nil {
		info, err = cachedMvnCommandResult.rs, cachedMvnCommandResult.e
		return
	}

	// 初始化 Maven 命令信息
	info = &MvnCommandInfo{}

	// 获取 java home
	info.JavaHome = GetJavaHome()

	// 获取 Maven 命令的路径
	info.Path = getMvnCommandOs()
	if info.Path == "" {
		err = ErrMvnNotFound
		cachedMvnCommandResult = &_MvnCommandResult{rs: nil, e: err}
		return
	}

	// 检查 Maven 版本
	ver, e := checkMvnVersion(info.Path, info.JavaHome)
	if e != nil {
		err = e
		cachedMvnCommandResult = &_MvnCommandResult{rs: info, e: err}
		return
	}
	info.MvnVersion = ver

	// 缓存检查结果
	cachedMvnCommandResult = &_MvnCommandResult{
		rs: info,
		e:  nil,
	}
	return
}

// executeMvnVersion 执行 Maven 命令获取版本信息
// 支持超时控制，避免命令执行时间过长
func executeMvnVersion(mvnPath string, javaHome string) (string, error) {
	cmd := exec.Command(mvnPath, "--version", "--batch-mode")

	// 设置环境变量
	cmd.Env = os.Environ()
	if javaHome != "" {
		cmd.Env = append(cmd.Env, "JAVA_HOME="+javaHome)
	}

	// 使用 channel 实现超时控制
	var output []byte
	done := make(chan error, 1)
	go func() {
		var err error
		output, err = cmd.Output()
		done <- err
	}()

	// 等待命令执行完成或超时
	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrCheckMvnVersion, err)
		}
	case <-time.After(8 * time.Second):
		// 超时时强制结束进程
		if err := cmd.Process.Kill(); err != nil {
			return "", fmt.Errorf("无法终止超时的 Maven 进程: %w", err)
		}
		return "", fmt.Errorf("%w: 执行 Maven 版本命令超时", ErrCheckMvnVersion)
	}

	return string(output), nil
}

// checkMvnVersion 检查 Maven 的版本
// 对于 Linux 和 MacOS 系统，如果首次执行失败会尝试修改文件权限后重试
func checkMvnVersion(mvnPath string, javaHome string) (string, error) {
	output, err := executeMvnVersion(mvnPath, javaHome)
	if err != nil {
		// 在 Unix 类系统上尝试修改文件权限后重试
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			_ = os.Chmod(mvnPath, 0755)
			output, err = executeMvnVersion(mvnPath, javaHome)
		}
		if err != nil {
			return "", err
		}
	}

	// 解析版本号
	ver := parseMvnVersion(output)
	if ver == "" {
		return "", fmt.Errorf("%w: 无法解析 Maven 版本信息", ErrCheckMvnVersion)
	}
	return ver, nil
}

// parseMvnVersion 解析 Maven 命令输出中的版本号
// 使用正则表达式匹配版本号信息
func parseMvnVersion(input string) string {
	// 匹配形如 "Apache Maven 3.6.3" 的版本信息
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

// getMvnCommandOs 根据操作系统查找 Maven 命令的路径
// 返回 Maven 可执行文件的绝对路径
func getMvnCommandOs() string {
	// 在系统 PATH 中查找 mvn 命令
	p, err := exec.LookPath("mvn")
	if err != nil {
		return ""
	}

	// 如果已经是绝对路径，直接返回
	if filepath.IsAbs(p) {
		return p
	}

	// 将相对路径转换为绝对路径
	absPath, err := filepath.Abs(p)
	if err == nil {
		return absPath
	}
	return ""
}

// init 函数在包初始化时执行
// 确保缓存变量被正确初始化
func init() {
	cachedMvnCommandResult = nil
}
