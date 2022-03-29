package tmmd

import (
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type YamlData struct {
	Data map[string][]interface{}
	// TODO better error handling
	Errors []error
}

var contextKey = parser.NewContextKey()

// Get returns a YAML metadata.
func Get(pc parser.Context) *YamlData {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*YamlData)
	return d
}

type yamlBlockExt struct {
	blocks []string
	kinds []ast.NodeKind
}



// Threat is an extension to allow for documenting threat modeling threats.
func CreateYamlBlockExtension(blocks ...string) *yamlBlockExt {
	return &yamlBlockExt{blocks, make([]ast.NodeKind, len(blocks))}
}

func (e *yamlBlockExt) Extend(m goldmark.Markdown) {
	parsers := make([]util.PrioritizedValue, len(e.blocks))	
	for i, block := range e.blocks {		
		e.kinds[i] = ast.NewNodeKind(block)
		parsers[i] = util.Prioritized(newBlockParser(block, e.kinds[i]), 100)	
	}

	m.Parser().AddOptions(parser.WithBlockParsers(
		parsers...
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(e.NewYamlBlockRenderer(), 100),
	))
}

// Get returns YAML items
func GetItems(name string, factory func()interface{}, pc parser.Context) []interface{} {	
	items := make([]interface{}, 0)
	v := pc.Get(contextKey)
	if v == nil {
		return items
	}
	d := v.(*YamlData)
	l := d.Data[name]
	
	for _, v := range l {
		item := factory()		
		mapstructure.Decode(v, &item)		
		items = append(items, item)		
	}	
	return items
}

// Returns any errors found while parsing YAML
func GetErrors(pc parser.Context) []error {
	v := pc.Get(contextKey)
	if v == nil {
		return []error{}
	}
	d := v.(*YamlData)
	return d.Errors
}
