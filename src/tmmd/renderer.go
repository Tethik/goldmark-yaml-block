package tmmd

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TMBlockRenderer struct {}

func NewThreatBlockRenderer() renderer.NodeRenderer {
	return &TMBlockRenderer{}
}

func (r *TMBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindThreat, r.renderThreatBlock)
	reg.Register(KindControl, r.renderControlBlock)
}

func (r *TMBlockRenderer) writeLines(w util.BufWriter, source []byte, n gast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.Write(line.Value(source))
	}
}

func (r *TMBlockRenderer) renderThreatBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre>\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("</pre>\n")
	}
	return gast.WalkContinue, nil
}

func (r *TMBlockRenderer) renderControlBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {	
	if entering {
		_, _ = w.WriteString("<pre>\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("</pre>\n")
	}
	return gast.WalkContinue, nil
}