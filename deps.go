package pom_component_parsing

import (
	"sort"
)

// DepsMap 依赖关系映射结构
// 用于存储和管理Maven项目中的依赖关系
// 键为依赖项的坐标(Coordinate)，值为依赖项的详细信息(depsElement)
type DepsMap struct {
	m map[Coordinate]depsElement // 内部使用map存储依赖关系
}

// newDepsMap 创建一个新的DepsMap实例
// 返回初始化好的DepsMap指针
func newDepsMap() *DepsMap {
	return &DepsMap{
		m: map[Coordinate]depsElement{}, // 初始化空map
	}
}

// ListAllEntries 列出所有依赖项元素
// 返回按照坐标排序的依赖项切片
// 排序保证了输出的稳定性，便于后续处理和展示
func (d *DepsMap) ListAllEntries() []depsElement {
	// 创建结果切片
	var rs []depsElement

	// 将map中所有元素添加到切片中
	for _, it := range d.m {
		rs = append(rs, it)
	}

	// 根据坐标对结果进行排序
	// 使用Coordinate的Compare方法作为排序依据
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].coordinate.Compare(rs[j].coordinate) < 0
	})

	return rs
}

// depsElement 依赖项元素结构
// 存储单个依赖项的完整信息
type depsElement struct {
	coordinate   Coordinate   // 依赖项的唯一坐标，包含groupId、artifactId和version
	children     []Dependency // 该依赖项的子依赖项列表
	relativePath string       // 依赖项相对于项目根目录的路径
}

// put 添加或更新依赖项
// 参数：
//   - coordinate: 依赖项的坐标
//   - children: 子依赖项列表
//   - path: 相对路径
func (d *DepsMap) put(coordinate Coordinate, children []Dependency, path string) {
	// 创建新的depsElement并存储到map中
	d.m[coordinate] = depsElement{
		coordinate:   coordinate,
		children:     children,
		relativePath: path,
	}
}

// allEmpty 检查是否所有依赖项都没有子依赖
// 返回:
//   - true: 所有依赖项都没有子依赖
//   - false: 存在至少一个有子依赖的依赖项
func (d *DepsMap) allEmpty() bool {
	// 遍历所有依赖项
	for _, it := range d.m {
		// 如果发现任何有子依赖的项，返回false
		if len(it.children) > 0 {
			return false
		}
	}
	// 所有项都没有子依赖，返回true
	return true
}

// Get 根据坐标获取依赖项
func (d *DepsMap) Get(coordinate Coordinate) (depsElement, bool) {
	elem, exists := d.m[coordinate]
	return elem, exists
}

// Size 返回依赖项的总数
func (d *DepsMap) Size() int {
	return len(d.m)
}

// Clear 清空所有依赖项
func (d *DepsMap) Clear() {
	d.m = make(map[Coordinate]depsElement)
}

// HasDependency 检查是否存在指定坐标的依赖项
func (d *DepsMap) HasDependency(coordinate Coordinate) bool {
	_, exists := d.m[coordinate]
	return exists
}
