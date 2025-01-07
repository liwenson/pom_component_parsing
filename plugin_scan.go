package pom_component_parsing

import (
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/vifraa/gopom"
)

// ScanDepsByPluginCommand 使用 Maven 插件命令扫描依赖关系。
// 该函数不再使用上下文，并使用默认日志打印日志信息。
func ScanDepsByPluginCommand(projectDir string, mvnCmdInfo *MvnCommandInfo) (*DepsMap, error) {
	// 查找项目的 Pom 配置文件中的 profiles
	profiles, err := findPomProfiles(filepath.Join(projectDir, "pom.xml"))
	if err != nil {
		// 打印错误信息
		log.Printf("查找 Pom profiles 时出错: %v\n", err)
	} else {
		// 打印找到的 profiles 数量
		log.Printf("找到 %d 个 profiles\n", len(profiles))
	}

	// 初始化 PluginGraphCmd 结构体
	c := PluginGraphCmd{
		MavenCmdInfo: mvnCmdInfo,
		Profiles:     profiles,
		Timeout:      time.Duration(60) * time.Second,
		ScanDir:      projectDir,
	}

	// 执行 Maven 图命令
	if err := c.RunC(); err != nil {
		// 打印执行失败的错误信息
		log.Printf("执行 Maven 图命令失败: %v\n", err)
		return nil, err
	}

	// 打印命令执行成功的信息
	log.Println("Maven 图命令执行成功，正在收集图文件...")

	// 收集插件结果文件
	return collectPluginResultFile(projectDir)
}

// collectPluginResultFile 收集项目目录中的 dependency-graph.json 文件并解析依赖关系。
func collectPluginResultFile(projectDir string) (*DepsMap, error) {
	var graphPaths []string

	// 遍历项目目录，查找所有的 dependency-graph.json 文件
	err := filepath.Walk(projectDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}
		if info.Name() == "dependency-graph.json" {
			// 记录找到的图文件路径
			log.Printf("找到图文件: %s\n", path)
			graphPaths = append(graphPaths, path)
		}
		return nil
	})
	if err != nil {
		// 打印遍历目录时的错误信息
		log.Printf("收集图文件时出错: %v\n", err)
	}

	// 初始化 DepsMap 以存储依赖关系
	rs := newDepsMap()

	// 遍历所有找到的图文件，解析并存储依赖关系
	for _, graphPath := range graphPaths {
		log.Printf("正在处理图文件: %s\n", graphPath)
		var g PluginGraphOutput

		// 从文件中读取图数据
		if err := g.ReadFromFile(graphPath); err != nil {
			// 打印读取文件时的错误信息，并继续处理下一个文件
			log.Printf("读取图文件时出错: %v\n", err)
			continue
		}

		// 构建依赖树
		tree, err := g.Tree()
		if err != nil {
			// 打印构建依赖树时的错误信息，并继续处理下一个文件
			log.Printf("构建依赖树时出错: %v\n", err)
			continue
		}

		// 计算图文件所在目录相对于项目根目录的相对路径
		relPath, err := filepath.Rel(projectDir, filepath.Dir(filepath.Dir(graphPath)))
		if err != nil {
			// 打印计算相对路径时的警告信息
			log.Printf("计算相对路径时出错: %v\n", err)
		}

		// 将解析后的依赖关系存储到 DepsMap 中
		rs.put(tree.Coordinate, tree.Children, filepath.Join(relPath, "pom.xml"))
	}

	return rs, nil
}

// findPomProfiles 解析指定的 pom.xml 文件，查找所有的 profiles。
func findPomProfiles(pomPath string) ([]string, error) {
	// 解析 pom.xml 文件
	project, err := gopom.Parse(pomPath)
	if err != nil {
		return nil, err
	}

	// 初始化 profiles 切片
	var profiles []string

	// 如果项目中存在 profiles，提取每个 profile 的 ID
	if project.Profiles != nil {
		for _, profile := range *project.Profiles {
			profiles = append(profiles, *profile.ID)
		}
	}

	return profiles, nil
}
