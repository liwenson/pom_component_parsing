package main

import (
	"fmt"
	maven "github.com/liwenson/pom_component_parsing"
	"github.com/liwenson/pom_component_parsing/model"
	"github.com/liwenson/pom_component_parsing/utils"
	"log"
	"time"
)

func main() {
	dir := "workspace/newton_buyer"
	modules, e := maven.ScanMavenProject(dir)
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
	// 让程序休眠 1 分钟
	time.Sleep(3 * time.Minute)

	fmt.Println("End:", time.Now())

}
