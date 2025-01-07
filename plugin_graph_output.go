package pom_component_parsing

import (
	"encoding/json"
	"os"
	"sync"

	"fmt"
	"log"
)

// PluginGraphOutput 表示 Maven 依赖图的结构，从 dependency-graph.json 文件中读取
type PluginGraphOutput struct {
	GraphName    string           `json:"graphName"`    // 图的名称
	Artifacts    []Artifact       `json:"artifacts"`    // 构成图的各个工件
	Dependencies []DependencyEdge `json:"dependencies"` // 工件之间的依赖关系
}

// Artifact 表示单个 Maven 工件的信息
type Artifact struct {
	GroupId    string   `json:"groupId"`    // 组ID
	ArtifactId string   `json:"artifactId"` // 工件ID
	Optional   bool     `json:"optional"`   // 是否为可选依赖
	Scopes     []string `json:"scopes"`     // 作用域
	Version    string   `json:"version"`    // 版本
}

// DependencyEdge 表示工件之间的依赖关系
type DependencyEdge struct {
	NumericFrom int `json:"numericFrom"` // 依赖来源工件的索引
	NumericTo   int `json:"numericTo"`   // 依赖目标工件的索引
}

// ReadFromFile 从指定路径读取并解析 dependency-graph.json 文件
func (d *PluginGraphOutput) ReadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取依赖图文件失败: %w", err)
	}

	var graph PluginGraphOutput
	if err := json.Unmarshal(data, &graph); err != nil {
		return fmt.Errorf("解析依赖图文件失败: %w", err)
	}

	*d = graph
	return nil
}

// Tree 构建依赖树，返回根节点的依赖结构
func (d *PluginGraphOutput) Tree() (*Dependency, error) {
	edges := d.buildEdgesMap()
	root, err := d.findRootNode()
	if err != nil {
		return nil, err
	}

	visited := make([]bool, len(d.Artifacts))
	return d.buildDependencyTree(root, visited, edges)
}

// buildDependencyTree 递归构建依赖树，防止循环依赖
func (d *PluginGraphOutput) buildDependencyTree(id int, visited []bool, edges map[int][]int) (*Dependency, error) {
	if visited[id] {
		return nil, fmt.Errorf("检测到循环依赖: 工件索引 %d", id)
	}
	visited[id] = true
	defer func() { visited[id] = false }()

	artifact := d.Artifacts[id]
	dependency := &Dependency{
		Coordinate: Coordinate{
			GroupId:    artifact.GroupId,
			ArtifactId: artifact.ArtifactId,
			Version:    artifact.Version,
		},
		Scope:    getFirstScope(artifact.Scopes),
		Children: []Dependency{},
	}

	for _, toID := range edges[id] {
		child, err := d.buildDependencyTree(toID, visited, edges)
		if err != nil {
			log.Printf("构建子依赖时出错: %v", err)
			continue
		}
		if child != nil {
			dependency.Children = append(dependency.Children, *child)
		}
	}

	return dependency, nil
}

// buildEdgesMap 构建从工件索引到依赖目标索引的映射，确保边的唯一性
func (d *PluginGraphOutput) buildEdgesMap() map[int][]int {
	edges := make(map[int][]int)
	seenEdges := make(map[int64]struct{})
	var mu sync.Mutex // 保护 seenEdges 和 edges 的并发安全

	for _, dep := range d.Dependencies {
		uniqueKey := int64(dep.NumericFrom)<<32 | int64(dep.NumericTo)
		mu.Lock()
		if _, exists := seenEdges[uniqueKey]; exists {
			mu.Unlock()
			continue
		}
		seenEdges[uniqueKey] = struct{}{}
		edges[dep.NumericFrom] = append(edges[dep.NumericFrom], dep.NumericTo)
		mu.Unlock()
	}

	return edges
}

// findRootNode 查找依赖图的根节点，即没有任何依赖来源的工件
func (d *PluginGraphOutput) findRootNode() (int, error) {
	isDependent := make([]bool, len(d.Artifacts))

	for _, dep := range d.Dependencies {
		if dep.NumericTo >= len(isDependent) || dep.NumericTo < 0 {
			return 0, fmt.Errorf("依赖目标索引 %d 超出工件范围", dep.NumericTo)
		}
		isDependent[dep.NumericTo] = true
	}

	var roots []int
	for idx, dependent := range isDependent {
		if !dependent {
			roots = append(roots, idx)
		}
	}

	if len(roots) == 0 {
		return 0, fmt.Errorf("未找到根节点")
	}

	if len(roots) > 1 {
		log.Printf("警告: 依赖图有多个根节点: %v", roots)
	}

	return roots[0], nil
}

// getFirstScope 获取作用域切片中的第一个元素，如果切片为空则返回空字符串
func getFirstScope(scopes []string) string {
	if len(scopes) > 0 {
		return scopes[0]
	}
	return ""
}
