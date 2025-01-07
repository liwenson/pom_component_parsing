//go:build windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

// KillProcessGroup 终止指定进程ID对应的进程组
// 在Windows中，通过发送CTRL_BREAK_EVENT来终止进程组
func KillProcessGroup(pid int) error {
	// 加载kernel32.dll
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGenerateConsoleCtrlEvent := kernel32.NewProc("GenerateConsoleCtrlEvent")

	// 发送CTRL_BREAK_EVENT (值为1) 到进程组
	r1, _, _ := procGenerateConsoleCtrlEvent.Call(
		uintptr(1), // CTRL_BREAK_EVENT
		uintptr(pid),
	)

	// 如果发送信号失败，尝试直接终止进程
	if r1 == 0 {
		handle, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pid))
		if err != nil {
			return fmt.Errorf("failed to open process %d: %w", pid, err)
		}
		defer syscall.CloseHandle(handle)

		if err := syscall.TerminateProcess(handle, 1); err != nil {
			return fmt.Errorf("failed to terminate process %d: %w", pid, err)
		}
	}

	return nil
}

// SetPGid 设置进程的创建标志和属性
// 在Windows中，通过设置CREATE_NEW_PROCESS_GROUP标志来创建新的进程组
func SetPGid(c *exec.Cmd) {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}

	// 设置进程创建标志
	// CREATE_NEW_PROCESS_GROUP (0x00000200) - 创建新的进程组
	c.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
}

// 以下是一些辅助函数，保持与非Windows版本类似的功能

// SendSignalToGroup 向进程组发送控制事件
func SendSignalToGroup(pid int, signal int) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGenerateConsoleCtrlEvent := kernel32.NewProc("GenerateConsoleCtrlEvent")

	r1, _, err := procGenerateConsoleCtrlEvent.Call(
		uintptr(signal),
		uintptr(pid),
	)

	if r1 == 0 {
		return fmt.Errorf("failed to send signal %d to process group %d: %w", signal, pid, err)
	}

	return nil
}
