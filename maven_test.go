package pom_component_parsing

import (
	"fmt"
	"github.com/liwenson/pom_component_parsing/model"
	"github.com/liwenson/pom_component_parsing/utils"
	"reflect"
	"testing"
)

func TestDependency_IsZero(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		want bool
	}{
		{
			name: "Empty dependency",
			dep:  Dependency{},
			want: true,
		},
		{
			name: "Non-empty dependency",
			dep: Dependency{
				Coordinate: Coordinate{
					GroupId:    "com.example",
					ArtifactId: "test",
					Version:    "1.0.0",
				},
			},
			want: false,
		},
		{
			name: "Dependency with children",
			dep: Dependency{
				Children: []Dependency{
					{
						Coordinate: Coordinate{
							GroupId:    "com.example",
							ArtifactId: "child",
							Version:    "1.0.0",
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dep.IsZero(); got != tt.want {
				t.Errorf("Dependency.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convDeps(t *testing.T) {
	tests := []struct {
		name string
		deps []Dependency
		want []model.DependencyItem
	}{
		{
			name: "Empty dependencies",
			deps: []Dependency{},
			want: nil,
		},
		{
			name: "Single dependency",
			deps: []Dependency{
				{
					Coordinate: Coordinate{
						GroupId:    "com.example",
						ArtifactId: "test",
						Version:    "1.0.0",
					},
					Scope: "compile",
				},
			},
			want: []model.DependencyItem{
				{
					Component: model.Component{
						CompName:           "com.example:test",
						CompVersion:        "1.0.0",
						IsDirectDependency: true,
						EcoRepo:            EcoRepo,
					},
					IsOnline:   model.IsOnlineTrue(),
					MavenScope: "compile",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convDeps(tt.deps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convDeps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convDep(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		want *model.DependencyItem
	}{
		{
			name: "Zero dependency",
			dep:  Dependency{},
			want: nil,
		},
		{
			name: "Test scope dependency",
			dep: Dependency{
				Coordinate: Coordinate{
					GroupId:    "com.example",
					ArtifactId: "test",
					Version:    "1.0.0",
				},
				Scope: "test",
			},
			want: &model.DependencyItem{
				Component: model.Component{
					CompName:    "com.example:test",
					CompVersion: "1.0.0",
					EcoRepo:     EcoRepo,
				},
				IsOnline:   model.IsOnlineFalse(),
				MavenScope: "test",
			},
		},
		{
			name: "Compile scope dependency with children",
			dep: Dependency{
				Coordinate: Coordinate{
					GroupId:    "com.example",
					ArtifactId: "parent",
					Version:    "1.0.0",
				},
				Scope: "compile",
				Children: []Dependency{
					{
						Coordinate: Coordinate{
							GroupId:    "com.example",
							ArtifactId: "child",
							Version:    "1.0.0",
						},
						Scope: "compile",
					},
				},
			},
			want: &model.DependencyItem{
				Component: model.Component{
					CompName:    "com.example:parent",
					CompVersion: "1.0.0",
					EcoRepo:     EcoRepo,
				},
				IsOnline:   model.IsOnlineTrue(),
				MavenScope: "compile",
				Dependencies: []model.DependencyItem{
					{
						Component: model.Component{
							CompName:    "com.example:child",
							CompVersion: "1.0.0",
							EcoRepo:     EcoRepo,
						},
						IsOnline:   model.IsOnlineTrue(),
						MavenScope: "compile",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _convDep(tt.dep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_convDep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanMavenProject(t *testing.T) {
	dir := "workspace/newton_buyer"
	modules, e := ScanMavenProject(dir)
	if e != nil {
		t.Errorf("组件解析失败 %v", e)
	} else {
		t.Log("组件解析结束")
	}

	var components []model.Component
	for _, m := range modules {
		components = append(components, m.ComponentList()...)
	}

	components = utils.DistinctSlice(components)

	for _, component := range components {
		if component.IsDirectDependency {
			fmt.Printf("component %v\n", component)
		}
	}

}
