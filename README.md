# goflow
golang工作流
支持脚本：lua
### 快速开始
```go
package main

import (
	"fmt"

	"github.com/kappere/goflow"
)

func main() {
	// 节点图: n4 -> n3 -> n2 -> n1
	// 默认运行lua函数run(param)，其中param为初始入参，前一节点运行结果会放入参数param["prev"]中
	// lua返回结果必须包含字段data（返回数据）,success（是否运行成功）,message（失败消息）

	// 创建4个节点：n1,n2,n3,n4
	n1 := goflow.NewNode([]*goflow.FlowNode{}, goflow.LuaScript{`
	function run(param)
		return {
			["data"] = "n1",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n2 := goflow.NewNode([]*goflow.FlowNode{n1}, goflow.LuaScript{`
	function run(param)
		return {
			["data"] = "n2",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n3 := goflow.NewNode([]*goflow.FlowNode{n2}, goflow.LuaScript{`
	function run(param)
		return {
			["data"] = "n3",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n4 := goflow.NewNode([]*goflow.FlowNode{n3}, goflow.LuaScript{`
	function run(param)
		return {
			["data"] = "n4",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	// 创建flow，true参数表示node出现错误停止flow运行
	flow := goflow.NewFlow(n4, true)
	result := flow.Run(map[string]interface{}{
		"v1": "456",
	})
	fmt.Printf("n4: %v\n", result[n4].Data)
	fmt.Printf("n3: %v\n", result[n3].Data)
	fmt.Printf("n2: %v\n", result[n2].Data)
	fmt.Printf("n1: %v\n", result[n1].Data)
}

```
