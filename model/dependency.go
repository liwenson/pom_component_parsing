package model

// DependencyItem 结构体表示一个依赖项，包括组件信息、子依赖、依赖类型、Maven 范围以及是否在线
type DependencyItem struct {
	Component                     // 嵌入的 Component 结构体，包含组件名称、版本和生态仓库信息
	Dependencies []DependencyItem `json:"dependencies,omitempty"` // 子依赖列表，包含当前依赖项的所有直接或间接依赖
	MavenScope   string           `json:"maven_scope,omitempty"`  // Maven 依赖的范围，例如 "compile", "test", "runtime" 等
	IsOnline     IsOnline         `json:"is_online"`              // 标识依赖项是否在线，true 表示在线仓库，false 表示本地仓库
}
