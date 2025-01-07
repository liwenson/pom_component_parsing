package pom_component_parsing

import (
	"fmt"
	"github.com/liwenson/pom_component_parsing/model"
	"log"
	"path/filepath"
)

// Dependency 表示一个Maven依赖项，包含坐标信息、子依赖和作用域。
type Dependency struct {
	Coordinate
	Children []Dependency `json:"children,omitempty"` // 子依赖列表
	Scope    string       `json:"scope"`              // 依赖作用域
}

// IsZero 判断Dependency是否为空，即没有子依赖且关键字段为nil。
func (d Dependency) IsZero() bool {
	return len(d.Children) == 0 && d.ArtifactId == "" && d.GroupId == "" && d.Version == ""
}

// String 返回Dependency的字符串表示，包括坐标和子依赖。
func (d Dependency) String() string {
	return fmt.Sprintf("%v: %v", d.Coordinate, d.Children)
}

// ScanMavenProject 扫描指定目录下的Maven项目，返回模块列表或错误。
func ScanMavenProject(dir string) ([]model.Module, error) {
	var modules []model.Module
	var deps *DepsMap

	// 检查Maven命令是否可用，若不可用则跳过扫描
	mvnCmdInfo, err := CheckMvnCommand()
	if err != nil {
		log.Println("检查Maven命令时出错:", err)
	} else {
		// 使用Maven插件命令扫描依赖
		deps, err = ScanDepsByPluginCommand(dir, mvnCmdInfo)
		if err != nil {
			log.Println("使用插件命令扫描依赖时出错:", err)
		}
	}

	// 如果依赖映射为空，返回检查错误
	if deps == nil {
		return nil, ErrInspection
	}

	// 遍历所有依赖项，构建模块信息
	for _, entry := range deps.ListAllEntries() {
		modules = append(modules, model.Module{
			PackageManager: "maven",
			ModuleName:     entry.coordinate.Name(),
			ModuleVersion:  entry.coordinate.Version,
			ModulePath:     filepath.Join(dir, entry.relativePath),
			Dependencies:   convDeps(entry.children),
			// ScanStrategy:   strategy, // 可根据需要启用扫描策略
		})
	}

	return modules, nil
}

// convDeps 将内部的Dependency切片转换为模型层的DependencyItem切片。
func convDeps(deps []Dependency) []model.DependencyItem {
	var rs []model.DependencyItem
	for _, it := range deps {
		d := _convDep(it)
		if d == nil {
			continue
		}
		d.IsDirectDependency = true // 标记为直接依赖
		rs = append(rs, *d)
	}
	return rs
}

// _convDep 辅助函数，将单个Dependency转换为模型层的DependencyItem。
func _convDep(dep Dependency) *model.DependencyItem {
	// 如果Dependency为空，返回nil
	if dep.IsZero() {
		return nil
	}

	// 创建DependencyItem并填充基本信息
	d := &model.DependencyItem{
		Component: model.Component{
			CompName:    dep.Name(),
			CompVersion: dep.Version,
			EcoRepo:     EcoRepo,
		},
		IsOnline:   model.IsOnlineTrue(),
		MavenScope: dep.Scope,
	}

	// 根据作用域决定是否为在线依赖
	if d.MavenScope == "test" || d.MavenScope == "provided" || d.MavenScope == "system" {
		d.IsOnline.SetOnline(false)
	}

	// 递归转换子依赖
	for _, it := range dep.Children {
		dd := _convDep(it)
		if dd == nil {
			continue
		}
		d.Dependencies = append(d.Dependencies, *dd)
	}
	return d
}

// EcoRepo 定义了生态系统和仓库信息，用于DependencyItem。
var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}

// ErrInspection 表示依赖检查失败的错误。
var ErrInspection = fmt.Errorf("依赖检查失败")
