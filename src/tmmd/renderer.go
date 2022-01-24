package tmmd

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TMBlockRenderer struct {
	e *yamlBlockExt
}

func (e *yamlBlockExt) NewYamlBlockRenderer() renderer.NodeRenderer {
	return &TMBlockRenderer{e}
}

func (r *TMBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	for _, kind := range r.e.kinds {
		reg.Register(kind, r.renderBlock)
	}
}

func (r *TMBlockRenderer) writeLines(w util.BufWriter, source []byte, n gast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.Write(line.Value(source))
	}
}

func (r *TMBlockRenderer) renderBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	// TODO: include kind, render better?
	if entering {
		_, _ = w.WriteString("<pre>\n")
		r.writeLines(w, source, node)
	} else {
		_, _ = w.WriteString("</pre>\n")
	}
	return gast.WalkContinue, nil
}
