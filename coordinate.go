package pom_component_parsing

import (
	"regexp"
	"strings"
)

// Coordinate 表示 Maven 工件的坐标信息
type Coordinate struct {
	GroupId    string `json:"group_id"`    // 组ID
	ArtifactId string `json:"artifact_id"` // 工件ID
	Version    string `json:"version"`     // 版本
}

// precompiledRegexp 是预编译的正则表达式，用于匹配所有空白字符
var precompiledRegexp = regexp.MustCompile(`\s`)

// IsSnapshotVersion 判断当前版本是否为 SNAPSHOT 版本
func (c Coordinate) IsSnapshotVersion() bool {
	// 检查版本字符串是否以 "-SNAPSHOT" 结尾
	return strings.HasSuffix(c.Version, "-SNAPSHOT")
}

// Normalize 返回一个新的 Coordinate，其中 GroupId、ArtifactId 和 Version 中的所有空白字符已被移除
func (c Coordinate) Normalize() Coordinate {
	return Coordinate{
		GroupId:    precompiledRegexp.ReplaceAllString(c.GroupId, ""),    // 去除 GroupId 前后空白字符
		ArtifactId: precompiledRegexp.ReplaceAllString(c.ArtifactId, ""), // 去除 ArtifactId 前后空白字符
		Version:    precompiledRegexp.ReplaceAllString(c.Version, ""),    // 去除 Version 前后空白字符
	}
}

// HasVersion 检查 Coordinate 是否包含有效的版本信息
func (c Coordinate) HasVersion() bool {
	normalized := c.Normalize()
	return normalized.Version != ""
}

// Name 返回 Coordinate 的 "GroupId:ArtifactId" 字符串表示
func (c Coordinate) Name() string {
	normalized := c.Normalize()
	return normalized.GroupId + ":" + normalized.ArtifactId
}

// String 返回 Coordinate 的完整字符串表示，如果存在版本则包括版本信息
func (c Coordinate) String() string {
	normalized := c.Normalize()
	if normalized.Version == "" {
		return normalized.GroupId + ":" + normalized.ArtifactId
	}
	return normalized.GroupId + ":" + normalized.ArtifactId + ":" + normalized.Version
}

// IsBad 判断 Coordinate 是否包含无效或格式错误的信息
func (c Coordinate) IsBad() bool {
	normalized := c.Normalize()
	// 判断任一字段是否以 "${", "[", "(" 开头，通常表示变量未解析或格式不正确
	if strings.HasPrefix(normalized.GroupId, "${") ||
		strings.HasPrefix(normalized.ArtifactId, "${") ||
		strings.HasPrefix(normalized.Version, "${") ||
		strings.HasPrefix(normalized.Version, "[") ||
		strings.HasPrefix(normalized.Version, "(") {
		return true
	}
	return false
}

// Complete 检查 Coordinate 是否包含完整的 GroupId、ArtifactId 和 Version，且不包含格式错误的信息
func (c Coordinate) Complete() bool {
	normalized := c.Normalize()
	// 检查所有字段是否非空且不包含格式错误的信息
	return normalized.GroupId != "" &&
		normalized.ArtifactId != "" &&
		normalized.Version != "" &&
		!normalized.IsBad()
}

// Compare 比较当前 Coordinate 与另一个 Coordinate 的顺序
// 返回值遵循 strings.Compare 的约定：-1 表示小于，0 表示等于，1 表示大于
func (c Coordinate) Compare(other Coordinate) int {
	// 先对当前 Coordinate 对象进行规范化处理，去除字段中的前后空白字符
	cNormalized := c.Normalize()
	// 对传入的另一个 Coordinate 对象也进行规范化处理
	otherNormalized := other.Normalize()

	// 按照 GroupId、ArtifactId、Version 的顺序逐一比较

	// 首先比较 GroupId 字段
	// 比较结果如果不等于 0，直接返回比较结果
	if cmp := strings.Compare(cNormalized.GroupId, otherNormalized.GroupId); cmp != 0 {
		return cmp
	}

	// 如果 GroupId 相同，则比较 ArtifactId 字段
	if cmp := strings.Compare(cNormalized.ArtifactId, otherNormalized.ArtifactId); cmp != 0 {
		return cmp
	}

	// 如果 ArtifactId 也相同，则比较 Version 字段
	return strings.Compare(cNormalized.Version, otherNormalized.Version)
}
