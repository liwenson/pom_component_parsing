//go:build !windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

// KillProcessGroup 终止指定进程ID对应的整个进程组
// 在Unix/Linux系统中，负数PID表示进程组ID
// 进程组ID等于该组长进程的PID的相反数
func KillProcessGroup(pid int) error {
	// 发送SIGKILL信号到进程组
	// -pid表示发送到整个进程组
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		return fmt.Errorf("failed to kill process group %d: %w", pid, err)
	}
	return nil
}

// SetPGid 设置进程的进程组ID属性
// 这将使新创建的进程成为进程组组长
// 进程组便于进行统一的信号处理和管理
func SetPGid(c *exec.Cmd) {
	// 如果SysProcAttr未初始化则创建它
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}

	// 设置Setpgid为true，使进程在新的进程组中启动
	// 这样可以将信号发送给整个进程组
	c.SysProcAttr.Setpgid = true
}

// 以下是一些扩展的辅助函数，用于更细粒度的进程控制

// SendSignalToGroup 向进程组发送指定信号
func SendSignalToGroup(pid int, sig syscall.Signal) error {
	if err := syscall.Kill(-pid, sig); err != nil {
		return fmt.Errorf("failed to send signal %v to process group %d: %w", sig, pid, err)
	}
	return nil
}

// GracefullyKillGroup 尝试优雅地终止进程组
// 首先发送SIGTERM，等待一段时间后如果进程仍然存在，则发送SIGKILL
func GracefullyKillGroup(pid int, timeout int) error {
	// 首先发送SIGTERM信号
	if err := SendSignalToGroup(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to process group %d: %w", pid, err)
	}

	// 检查进程是否存在
	exists := func(pid int) bool {
		return syscall.Kill(pid, 0) == nil
	}

	// 等待进程退出或超时
	for i := 0; i < timeout; i++ {
		if !exists(pid) {
			return nil
		}
		syscall.Sleep(time.Second)
	}

	// 如果进程仍然存在，发送SIGKILL信号
	return KillProcessGroup(pid)
}

// IsProcessExists 检查指定PID的进程是否存在
func IsProcessExists(pid int) bool {
	// 发送空信号(0)来检查进程是否存在
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}

// GetProcessGroup 获取进程的进程组ID
func GetProcessGroup(pid int) (int, error) {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return 0, fmt.Errorf("failed to get process group for pid %d: %w", pid, err)
	}
	return pgid, nil
}
