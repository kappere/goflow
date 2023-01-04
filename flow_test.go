package goflow

import (
	"fmt"
	"sync"
	"testing"
)

func TestFlow(t *testing.T) {
	n1 := NewNode([]*FlowNode{}, LuaScript{`
	function run(param)
		return {
			["data"] = param["iter"],
			["success"] = true,
			["message"] = "err1322"
		};
	end
	`})
	n2 := NewNode([]*FlowNode{n1}, LuaScript{`
	function run(param)
		return {
			["data"] = param["iter"],
			["success"] = true,
			["message"] = "err1322"
		};
	end
	`})
	n3 := NewNode([]*FlowNode{n2}, LuaScript{`
	function run(param)
		return {
			["data"] = param["iter"],
			["success"] = true,
			["message"] = "err1322"
		};
	end
	`})
	n4 := NewNode([]*FlowNode{n3}, LuaScript{`
	function run(param)
		return {
			["data"] = param["iter"],
			["success"] = true,
			["message"] = "err1322"
		};
	end
	`})
	flow := NewFlow(n4, true)
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(a int) {
			result := flow.Run(map[string]interface{}{
				"iter": a,
				"a": map[string]interface{}{
					"c": 333,
					"d": "eee",
				},
				"b": "456",
			})
			if a%100 == 0 {
				t.Logf("iter: %v, result: %v", a, result[n4].Data)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestDemo(t *testing.T) {
	// 节点图: n4 -> n3 -> n2 -> n1
	// 默认运行lua函数run(param)，其中param为初始入参，前一节点运行结果会放入参数param["prev"]中
	// lua返回结果必须包含字段data（返回数据）,success（是否运行成功）,message（失败消息）

	// 创建4个节点：n1,n2,n3,n4
	n1 := NewNode([]*FlowNode{}, LuaScript{`
	function run(param)
		return {
			["data"] = "n1",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n2 := NewNode([]*FlowNode{n1}, LuaScript{`
	function run(param)
		return {
			["data"] = "n2",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n3 := NewNode([]*FlowNode{n2}, LuaScript{`
	function run(param)
		return {
			["data"] = "n3",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	n4 := NewNode([]*FlowNode{n3}, LuaScript{`
	function run(param)
		return {
			["data"] = "n4",
			["success"] = true,
			["message"] = "ok"
		};
	end
	`})
	// 创建flow，true参数表示node出现错误停止flow运行
	flow := NewFlow(n4, true)
	result := flow.Run(map[string]interface{}{
		"v1": "456",
	})
	fmt.Printf("n4: %v\n", result[n4].Data)
	fmt.Printf("n3: %v\n", result[n3].Data)
	fmt.Printf("n2: %v\n", result[n2].Data)
	fmt.Printf("n1: %v\n", result[n1].Data)
}
