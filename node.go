package goflow

type FlowNode struct {
	Prev   []*FlowNode
	Next   []*FlowNode
	Script Script
	In     int
}

type NodeResult struct {
	Data  interface{}
	Error error
}

func NewNode(next []*FlowNode, script Script) *FlowNode {
	return &FlowNode{
		Prev:   make([]*FlowNode, 0),
		Next:   next,
		Script: script,
		In:     0,
	}
}
