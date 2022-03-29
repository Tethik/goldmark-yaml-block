package tmmd

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v2"
)

type blockParser struct {
	startRegexp *regexp.Regexp
	endRegexp *regexp.Regexp
	kind ast.NodeKind
}

func newBlockParser(blockType string, kind ast.NodeKind) *blockParser {
	return &blockParser {
		startRegexp: regexp.MustCompile("^```"+`(ya?ml)?\s*`+blockType+`\s*`),
		endRegexp: regexp.MustCompile("^```"+`\s*`),	
		kind: kind,
	}
}

func (s *blockParser) isStartTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))	
	return s.startRegexp.Match(line)
}

func (s *blockParser) isEndTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	return s.endRegexp.Match(line)
}

func (s *blockParser) Trigger() []byte {
	return []byte{'`'}
}

// Open parses the current line and returns a result of parsing.
//
// Open must not parse beyond the current line.
// If Open has been able to parse the current line, Open must advance a reader
// position by consumed byte length.
//
// If Open has not been able to parse the current line, Open should returns
// (nil, NoChildren). If Open has been able to parse the current line, Open
// should returns a new Block node and returns HasChildren or NoChildren.
func (s *blockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	
	if (pos > 0 || !s.isStartTag(line)) {		
		return nil, parser.NoChildren
	}		
	node := NewYamlNode(s.kind)	
	return node, parser.NoChildren
}


// Continue parses the current line and returns a result of parsing.
//
// Continue must not parse beyond the current line.
// If Continue has been able to parse the current line, Continue must advance
// a reader position by consumed byte length.
//
// If Continue has not been able to parse the current line, Continue should
// returns Close. If Continue has been able to parse the current line,
// Continue should returns (Continue | NoChildren) or
// (Continue | HasChildren)
func (s *blockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()	
	pos := pc.BlockOffset()

	if pos == 0 && s.isEndTag(line) {				
		reader.AdvanceLine()		
		return parser.Close
	}	

	if node.Kind() != s.kind {
		return parser.Close
	}

	text := text.NewSegment(segment.Start, segment.Stop)	
	node.Lines().Append(text)
	
	return parser.Continue | parser.NoChildren
}


func (s *blockParser) CanAcceptIndentedLine() bool {
	return true
}

func (s *blockParser) CanInterruptParagraph() bool {
	return false
}

// Close will be called when the parser returns Close.
func (s *blockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {	
	d, ok := pc.Get(contextKey).(*YamlData)
	if !ok {		
		d = &YamlData{
			Data: make(map[string][]interface{}),
		}
	}
	
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}	

	var item map[string]interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &item); err != nil {				
		informativeError := fmt.Errorf("parsing yaml block failed: \n```\n%s\n```\nerror: %s", buf.Bytes(), err)
		d.Errors = append(d.Errors, informativeError)
	} else {			
		if _, ok := d.Data[s.kind.String()]; !ok {
			d.Data[s.kind.String()] = make([]interface{}, 0)
		}
		d.Data[s.kind.String()] = append(d.Data[s.kind.String()], item)
	}

	pc.Set(contextKey, d)
}

