package tmmd

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TMData struct {
	Data map[ast.NodeKind][]interface{}
	// TODO better error handling
	Error error
}

var contextKey = parser.NewContextKey()

// Get returns a YAML metadata.
func Get(pc parser.Context) *TMData {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*TMData)
	return d
}

type threatModelingExt struct {
}

// Threat is an extension to allow for documenting threat modeling threats.
var ThreatModelingExtension = &threatModelingExt{}

func (e *threatModelingExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewThreatParser(), 100),
		util.Prioritized(NewControlParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewThreatBlockRenderer(), 100),		
	))
}
