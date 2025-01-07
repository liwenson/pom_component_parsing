package model

// Component 结构体表示一个组件，包括名称、版本和生态仓库信息
type Component struct {
	CompName           string `json:"comp_name"`            // 组件名称
	CompVersion        string `json:"comp_version"`         // 组件版本
	IsDirectDependency bool   `json:"is_direct_dependency"` // 是否为直接依赖
	EcoRepo                   // 嵌入的生态仓库信息
}

// EcoRepo 结构体表示组件所属的生态系统及其仓库
type EcoRepo struct {
	Ecosystem  string `json:"ecosystem"`  // 生态系统名称，例如 "npm", "maven"
	Repository string `json:"repository"` // 仓库地址或名称
}
