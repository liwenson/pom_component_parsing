package pom_component_parsing

import (
	"fmt"
	"log"
	"os"
	"strings"
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

// Run 执行 Maven 图形命令
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
	c := m.MavenCmdInfo.Command(args...)
	c.Dir = m.ScanDir
	utils.SetPGid(c)

	// 将命令的标准输出和标准错误输出指向空设备，忽略输出
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// 打印启动命令的信息
	log.Printf("开始执行命令: %s, 目录: %s\n", c.String(), c.Dir)

	// 启动命令
	if err := c.Start(); err != nil {
		log.Printf("启动 Maven 命令失败: %s\n", err.Error())
		return fmt.Errorf("启动 Maven 命令失败: %w", err)
	}

	// 如果设置了超时时间，则启动一个计时器，超时后终止命令
	if m.Timeout > 0 {
		go func() {
			time.Sleep(m.Timeout)
			log.Println("Maven 命令执行超时，终止命令")
			if c.Process != nil {
				if err := c.Process.Kill(); err != nil {
					log.Printf("终止 Maven 进程失败: %v\n", err)
				}
			}
		}()
	}

	// 等待命令执行完成
	if err := c.Wait(); err != nil {
		exitCode := c.ProcessState.ExitCode()
		log.Printf("执行 mvn 时发生错误: %s. 退出码: %d\n", err.Error(), exitCode)
		return fmt.Errorf("mvn 执行出错: %w", err)
	}

	// 获取退出码
	exitCode := c.ProcessState.ExitCode()
	if exitCode != 0 {
		log.Printf("mvn 执行出错。退出码: %d\n", exitCode)
		return fmt.Errorf("mvn 执行出错，退出码: %d", exitCode)
	}

	// 命令成功完成
	log.Println("Maven 命令执行成功")
	return nil
}
