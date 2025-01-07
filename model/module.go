package model

import (
	"fmt"
	"github.com/liwenson/pom_component_parsing/utils"
)

// Module 结构体用于表示一个模块的信息，包括名称、版本、路径、包管理器以及依赖项
type Module struct {
	ModuleName     string           `json:"module_name"`            // 模块的名称
	ModuleVersion  string           `json:"module_version"`         // 模块的版本
	ModulePath     string           `json:"module_path"`            // 模块的路径
	PackageManager string           `json:"package_manager"`        // 使用的包管理器，例如 npm, maven 等
	Dependencies   []DependencyItem `json:"dependencies,omitempty"` // 模块的依赖项列表，如果为空则在 JSON 中省略
}

// String 返回模块的字符串表示，格式为 "[包管理器]模块名称@模块版本"
func (m Module) String() string {
	var s = "[" + m.PackageManager + "]" + m.ModuleName
	if m.ModuleVersion != "" {
		s += "@" + m.ModuleVersion
	}
	return s
}

// IsZero 判断模块是否为空，即没有依赖项、名称和版本都为空
func (m Module) IsZero() bool {
	return len(m.Dependencies) == 0 && m.ModuleName == "" && m.ModuleVersion == ""
}

// ComponentList 返回模块中所有的组件列表，包含所有直接和间接的依赖项
func (m Module) ComponentList() []Component {
	var r = make(map[Component]struct{})
	__componentList(m.Dependencies, r)
	fmt.Printf("r: %v", r)
	return utils.KeysOfMap(r)
}

// __componentList 是一个辅助函数，用于递归遍历依赖项并收集所有的组件
func __componentList(deps []DependencyItem, m map[Component]struct{}) {
	for _, dep := range deps {
		m[dep.Component] = struct{}{}
		__componentList(dep.Dependencies, m)
	}
}