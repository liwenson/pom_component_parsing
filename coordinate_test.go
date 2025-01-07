package pom_component_parsing

import (
	"testing"
)

// TestCoordinate_IsSnapshotVersion 测试 Coordinate 的 IsSnapshotVersion 方法
func TestCoordinate_IsSnapshotVersion(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected bool
	}{
		{
			name: "Snapshot版本",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0-SNAPSHOT",
			},
			expected: true,
		},
		{
			name: "非Snapshot版本",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: false,
		},
		{
			name: "版本为空",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.IsSnapshotVersion()
			if result != tt.expected {
				t.Errorf("IsSnapshotVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_Normalize 测试 Coordinate 的 Normalize 方法
func TestCoordinate_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected Coordinate
	}{
		{
			name: "包含空白字符",
			coord: Coordinate{
				GroupId:    " com. example ",
				ArtifactId: " example-artifact ",
				Version:    " 1.0.0 ",
			},
			expected: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
		},
		{
			name: "不含空白字符",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
		},
		{
			name: "所有字段为空",
			coord: Coordinate{
				GroupId:    "",
				ArtifactId: "",
				Version:    "",
			},
			expected: Coordinate{
				GroupId:    "",
				ArtifactId: "",
				Version:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.Normalize()
			if result != tt.expected {
				t.Errorf("Normalize() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_HasVersion 测试 Coordinate 的 HasVersion 方法
func TestCoordinate_HasVersion(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected bool
	}{
		{
			name: "包含版本",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: true,
		},
		{
			name: "版本为空",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "",
			},
			expected: false,
		},
		{
			name: "版本为空白字符",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "   ",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.HasVersion()
			if result != tt.expected {
				t.Errorf("HasVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_Name 测试 Coordinate 的 Name 方法
func TestCoordinate_Name(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected string
	}{
		{
			name: "正常情况",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: "com.example:example-artifact",
		},
		{
			name: "包含空白字符",
			coord: Coordinate{
				GroupId:    " com. example ",
				ArtifactId: " example-artifact ",
				Version:    "1.0.0",
			},
			expected: "com.example:example-artifact",
		},
		{
			name: "ArtifactId为空",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "",
				Version:    "1.0.0",
			},
			expected: "com.example:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.Name()
			if result != tt.expected {
				t.Errorf("Name() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_String 测试 Coordinate 的 String 方法
func TestCoordinate_String(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected string
	}{
		{
			name: "包含版本",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: "com.example:example-artifact:1.0.0",
		},
		{
			name: "不包含版本",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "",
			},
			expected: "com.example:example-artifact",
		},
		{
			name: "包含空白字符",
			coord: Coordinate{
				GroupId:    " com. example ",
				ArtifactId: " example-artifact ",
				Version:    " 1.0.0 ",
			},
			expected: "com.example:example-artifact:1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.String()
			if result != tt.expected {
				t.Errorf("String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_IsBad 测试 Coordinate 的 IsBad 方法
func TestCoordinate_IsBad(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected bool
	}{
		{
			name: "有效Coordinate",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: false,
		},
		{
			name: "GroupId以${开头",
			coord: Coordinate{
				GroupId:    "${group.id}",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: true,
		},
		{
			name: "ArtifactId以${开头",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "${artifact.id}",
				Version:    "1.0.0",
			},
			expected: true,
		},
		{
			name: "Version以${开头",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "${version}",
			},
			expected: true,
		},
		{
			name: "Version以[开头",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "[1.0.0]",
			},
			expected: true,
		},
		{
			name: "Version以(开头",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "(1.0.0)",
			},
			expected: true,
		},
		{
			name: "多种无效前缀",
			coord: Coordinate{
				GroupId:    "${group.id}",
				ArtifactId: "[artifact.id]",
				Version:    "(1.0.0)",
			},
			expected: true,
		},
		{
			name: "Version为空",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.IsBad()
			if result != tt.expected {
				t.Errorf("IsBad() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_Complete 测试 Coordinate 的 Complete 方法
func TestCoordinate_Complete(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		expected bool
	}{
		{
			name: "完整且有效的Coordinate",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: true,
		},
		{
			name: "缺少Version",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "",
			},
			expected: false,
		},
		{
			name: "含有无效前缀",
			coord: Coordinate{
				GroupId:    "${group.id}",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: false,
		},
		{
			name: "GroupId为空",
			coord: Coordinate{
				GroupId:    "",
				ArtifactId: "example-artifact",
				Version:    "1.0.0",
			},
			expected: false,
		},
		{
			name: "ArtifactId为空",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "",
				Version:    "1.0.0",
			},
			expected: false,
		},
		{
			name: "Version为空白字符",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "example-artifact",
				Version:    "   ",
			},
			expected: false,
		},
		{
			name: "所有字段为空",
			coord: Coordinate{
				GroupId:    "",
				ArtifactId: "",
				Version:    "",
			},
			expected: false,
		},
		{
			name: "包含空白字符且有效",
			coord: Coordinate{
				GroupId:    " com. example ",
				ArtifactId: " example-artifact ",
				Version:    " 1.0.0 ",
			},
			expected: true,
		},
		{
			name: "包含空白字符且无效",
			coord: Coordinate{
				GroupId:    " com. example ",
				ArtifactId: "${example-artifact}",
				Version:    "1.0.0",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.Complete()
			if result != tt.expected {
				t.Errorf("Complete() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCoordinate_Compare 测试 Coordinate 的 Compare 方法
func TestCoordinate_Compare(t *testing.T) {
	tests := []struct {
		name     string
		coord    Coordinate
		other    Coordinate
		expected int
	}{
		{
			name: "完全相同",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			expected: 0,
		},
		{
			name: "GroupId不同",
			coord: Coordinate{
				GroupId:    "com.alpha",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.beta",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			expected: -1,
		},
		{
			name: "ArtifactId不同",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifactA",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifactB",
				Version:    "1.0.0",
			},
			expected: -1,
		},
		{
			name: "Version不同",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "2.0.0",
			},
			expected: -1,
		},
		{
			name: "GroupId更大",
			coord: Coordinate{
				GroupId:    "com.beta",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.alpha",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			expected: 1,
		},
		{
			name: "ArtifactId更大",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifactB",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifactA",
				Version:    "1.0.0",
			},
			expected: 1,
		},
		{
			name: "Version更大",
			coord: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "2.0.0",
			},
			other: Coordinate{
				GroupId:    "com.example",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			expected: 1,
		},
		{
			name: "多个字段不同",
			coord: Coordinate{
				GroupId:    "com.alpha",
				ArtifactId: "artifactB",
				Version:    "2.0.0",
			},
			other: Coordinate{
				GroupId:    "com.beta",
				ArtifactId: "artifactA",
				Version:    "1.0.0",
			},
			expected: -1,
		},
		{
			name: "包含空白字符",
			coord: Coordinate{
				GroupId:    " com. alpha ",
				ArtifactId: " artifact ",
				Version:    " 1.0.0 ",
			},
			other: Coordinate{
				GroupId:    "com.alpha",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			expected: 0,
		},
		{
			name: "另一方包含空白字符",
			coord: Coordinate{
				GroupId:    "com.alpha",
				ArtifactId: "artifact",
				Version:    "1.0.0",
			},
			other: Coordinate{
				GroupId:    " com. alpha ",
				ArtifactId: " artifact ",
				Version:    " 1.0.0 ",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.coord.Compare(tt.other)
			if result != tt.expected {
				t.Errorf("Compare() = %d, want %d", result, tt.expected)
			}
		})
	}
}
