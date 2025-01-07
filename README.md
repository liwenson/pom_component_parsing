## pom文件组件解析

该模块是pom文件解析模块，主要功能是解析pom文件，获取pom文件中的组件信息，包括groupId、artifactId、version、scope、type、classifier、optional、exclusions、dependencies、dependencyManagement、repositories、pluginRepositories、parent、properties等。



### 使用

```bash
go get github.com/liwenson/pom_component_parsing
```

```go
package main

import (
	"fmt"
	"log"
	"github.com/liwenson/pom_component_parsing"
	"github.com/liwenson/pom_component_parsing/model"
	"github.com/liwenson/pom_component_parsing/utils"
)

func main() {
	dir := "workspace/newton_buyer"
	modules, e := ScanMavenProject(dir)
	if e != nil {
		log.Fatalf("组件解析失败 %v", e)
	} else {
		log.Println("组件解析结束")
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
```