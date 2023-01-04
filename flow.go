package goflow

import "container/list"

type Flow struct {
	Root         *FlowNode
	BreakOnError bool
}

func NewFlow(root *FlowNode, breakOnError bool) *Flow {
	return &Flow{
		Root:         root,
		BreakOnError: breakOnError,
	}
}

func (flow *Flow) traverse(node *FlowNode, currentNodeSet map[*FlowNode]bool) []*FlowNode {
	nodes := make([]*FlowNode, 0)
	if _, exists := currentNodeSet[node]; exists {
		return nodes
	}
	nodes = append(nodes, node)
	currentNodeSet[node] = true
	for _, nextNode := range node.Next {
		nodes = append(nodes, flow.traverse(nextNode, currentNodeSet)...)
	}
	return nodes
}

func (flow *Flow) topologySort() []*FlowNode {
	nodes := flow.traverse(flow.Root, make(map[*FlowNode]bool))
	// 初始化节点信息
	for _, node := range nodes {
		node.In = 0
		if len(node.Prev) != 0 {
			node.Prev = make([]*FlowNode, 0)
		}
	}
	// 更新节点入度出度
	for _, node := range nodes {
		for _, nextNode := range node.Next {
			nextNode.In++
			nextNode.Prev = append(nextNode.Prev, node)
		}
	}
	q := list.New()
	topologySeq := make([]*FlowNode, 0)
	for _, node := range nodes {
		if node.In == 0 {
			q.PushBack(node)
		}
	}

	for q.Len() > 0 {
		head := q.Front()
		q.Remove(head)
		headValue := head.Value.(*FlowNode)
		topologySeq = append(topologySeq, headValue)
		for _, nextNode := range headValue.Next {
			nextNode.In--
			if nextNode.In == 0 {
				q.PushBack(nextNode)
			}
		}
	}
	return topologySeq
}

func (flow *Flow) Run(param map[string]interface{}) map[*FlowNode]*NodeResult {
	result := map[*FlowNode]*NodeResult{}
	nodes := flow.topologySort()
	for _, node := range nodes {
		param["prev"] = result
		r, err := node.Script.Run(param)
		result[node] = &NodeResult{
			Data:  r,
			Error: err,
		}
		if flow.BreakOnError && err != nil {
			break
		}
	}
	return result
}
