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

var controlStartRegexp = regexp.MustCompile("^```"+`ya?ml\s*control\s*`)
var controlEndRegexp = regexp.MustCompile("^```"+`\s*`)

type controlParser struct {
}

var defaultControlParser = &controlParser{}

func (s *controlParser) isStartTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))	
	return controlStartRegexp.Match(line)
}

func (s *controlParser) isEndTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	return controlEndRegexp.Match(line)
}


// NewcontrolParser return a new BlockParser that parses
// strikethrough expressions.
func NewControlParser() parser.BlockParser {
	return defaultControlParser
}

func (s *controlParser) Trigger() []byte {
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
func (s *controlParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	
	if (pos > 0 || !s.isStartTag(line)) {		
		return nil, parser.NoChildren
	}		
	node := NewControlNode()	
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
func (s *controlParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()	
	pos := pc.BlockOffset()

	if pos == 0 && s.isEndTag(line) {				
		reader.AdvanceLine()		
		return parser.Close
	}	

	if node.Kind() != KindControl {
		return parser.Close
	}

	text := text.NewSegment(segment.Start, segment.Stop)	
	node.Lines().Append(text)
	
	return parser.Continue | parser.NoChildren
}


func (s *controlParser) CanAcceptIndentedLine() bool {
	return true
}

func (s *controlParser) CanInterruptParagraph() bool {
	return false
}

type ControlData struct {
	Slug string `yaml:"slug,omitempty"`
	Title string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Mitigates []string `yaml:"mitigates,omitempty"`
}

// Close will be called when the parser returns Close.
func (s *controlParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {	
	d, ok := pc.Get(contextKey).(*data)
	if !ok {
		fmt.Println("empty")
		d = &data{}
	} 
	fmt.Println(d)
	
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}	

	fmt.Println(buf.String())
	
	var control ControlData
	control.Mitigates = make([]string, 0)
	if err := yaml.Unmarshal(buf.Bytes(), &control); err != nil {
		fmt.Println(err)
		d.Error = err
	} else {	
		d.Controls = append(d.Controls, control)
	}

	pc.Set(contextKey, d)
}
