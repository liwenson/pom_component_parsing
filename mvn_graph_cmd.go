package pom_component_parsing

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/liwenson/pom_component_parsing/utils"
)

// PluginGraphCmd 用于执行 com.github.ferstl:depgraph-maven-plugin:4.0.1:graph 命令的辅助结构体
type PluginGraphCmd struct {
	Profiles     []string        // Maven 配置文件
	Timeout      time.Duration   // 超时时间
	ScanDir      string          // 扫描目录
	MavenCmdInfo *MvnCommandInfo // Maven 命令信息
}

// RunC 执行 Maven 图形命令，并添加超时控制以防止进程无法释放
func (m PluginGraphCmd) RunC() error {
	// 构建 Maven 命令参数
	args := []string{"com.github.ferstl:depgraph-maven-plugin:4.0.1:graph", "-DgraphFormat=json"}
	// 配置 Maven 参数以允许不安全的 TLS 连接
	args = append(args,
		"-Dmaven.wagon.http.ssl.ignore.validity.dates=true",
		"-Dmaven.resolver.transport=wagon",
		"-Dmaven.wagon.http.ssl.allowall=true",
		"-Dmaven.wagon.http.ssl.insecure=true",
	)

	// 如果有指定配置文件，则添加 -P 参数
	if len(m.Profiles) > 0 {
		args = append(args, "-P")
		args = append(args, strings.Join(m.Profiles, ","))
	}

	// 获取 Maven 命令执行实例
	cmd := m.MavenCmdInfo.Command(args...)
	cmd.Dir = m.ScanDir
	utils.SetPGid(cmd)

	// 将命令的标准输出和标准错误输出指向对应的输出流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 打印启动命令的信息
	log.Printf("开始执行命令: %s, 目录: %s\n", cmd.String(), cmd.Dir)

	// 创建上下文，以便可以在超时后取消命令执行
	ctx := context.Background()
	if m.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.Timeout)
		defer cancel()
	}

	// 关联上下文到命令，以便在上下文取消时终止命令
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	// 重新设置命令的工作目录与 PGid
	cmd.Dir = m.ScanDir
	utils.SetPGid(cmd)

	// 将命令的标准输出和标准错误输出指向对应的输出流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("启动 Maven 命令失败: %s\n", err.Error())
		return fmt.Errorf("启动 Maven 命令失败: %w", err)
	}

	// 通道用于接收命令执行完成的错误信息
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 上下文超时，尝试杀死进程
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("发送SIGTERM信号失败: %s\n", err.Error())
			return fmt.Errorf("发送SIGTERM信号失败: %w", err)
		}
		// 等待进程退出
		select {
		case <-done:
			log.Println("Maven 命令已超时并被终止")
			return fmt.Errorf("Maven 命令超时并被终止")
		case <-time.After(5 * time.Second):
			// 如果进程在5秒内没有响应SIGTERM，则强制杀死
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("强制杀死进程失败: %s\n", err.Error())
				return fmt.Errorf("强制杀死进程失败: %w", err)
			}
			log.Println("Maven 命令被强制杀死")
			return fmt.Errorf("Maven 命令被强制杀死")
		}
	case err := <-done:
		// 命令执行完成
		if err != nil {
			// 命令执行出错，获取退出码
			exitError, ok := err.(*exec.ExitError)
			if ok {
				exitCode := exitError.ExitCode()
				log.Printf("执行 mvn 时发生错误: %s. 退出码: %d\n", err.Error(), exitCode)
				return fmt.Errorf("mvn 执行出错: %w", err)
			}
			// 其他类型的错误
			log.Printf("执行 mvn 时发生未知错误: %s\n", err.Error())
			return fmt.Errorf("mvn 执行出错: %w", err)
		}

		// 获取退出码
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			log.Printf("mvn 执行出错。退出码: %d\n", exitCode)
			return fmt.Errorf("mvn 执行出错，退出码: %d", exitCode)
		}

		// 命令成功完成
		log.Println("Maven 命令执行成功")
		return nil
	}
}
