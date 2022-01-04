package tmmd

import (
	gast "github.com/yuin/goldmark/ast"
)

// A TaskCheckBox struct represents a checkbox of a task list.
type ThreatNode struct {
	gast.BaseBlock
}

// Dump implements Node.Dump.
func (n *ThreatNode) Dump(source []byte, level int) {
	m := map[string]string{
		
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindThreat is a NodeKind of the TaskCheckBox node.
var KindThreat = gast.NewNodeKind("Threat")

// Kind implements Node.Kind.
func (n *ThreatNode) Kind() gast.NodeKind {
	return KindThreat
}

// IsRaw seems to prevent any child nodes from being parsed
func (n *ThreatNode) IsRaw() bool {
	return true
}

// NewTaskCheckBox returns a new TaskCheckBox node.
func NewThreatNode() *ThreatNode {
	return &ThreatNode{
	}
}


// A TaskCheckBox struct represents a checkbox of a task list.
type ControlNode struct {
	gast.BaseBlock	
}

// Dump implements Node.Dump.
func (n *ControlNode) Dump(source []byte, level int) {
	m := map[string]string{
		
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindThreat is a NodeKind of the TaskCheckBox node.
var KindControl = gast.NewNodeKind("Control")

// Kind implements Node.Kind.
func (n *ControlNode) Kind() gast.NodeKind {
	return KindControl
}

// IsRaw seems to prevent any child nodes from being parsed
func (n *ControlNode) IsRaw() bool {
	return true
}

// NewTaskCheckBox returns a new TaskCheckBox node.
func NewControlNode() *ControlNode {
	return &ControlNode{}
}
