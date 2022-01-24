package tmmd

import (
	gast "github.com/yuin/goldmark/ast"
)

// A TaskCheckBox struct represents a checkbox of a task list.
type TMNode struct {
	gast.BaseBlock
	kind gast.NodeKind
}

// Dump implements Node.Dump.
func (n *TMNode) Dump(source []byte, level int) {
	m := map[string]string{
		
	}
	gast.DumpHelper(n, source, level, m, nil)
}


// Kind implements Node.Kind.
func (n *TMNode) Kind() gast.NodeKind {
	return n.kind
}

// IsRaw seems to prevent any child nodes from being parsed
func (n *TMNode) IsRaw() bool {
	return true
}

// Returns a new ThreatNode
func NewTMNode(kind gast.NodeKind) *TMNode {
	return &TMNode{
		kind: kind,
	}
}




