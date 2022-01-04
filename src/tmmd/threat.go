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

var threatStartRegexp = regexp.MustCompile("^```"+`ya?ml\s*threat\s*`)
var threatEndRegexp = regexp.MustCompile("^```"+`\s*`)

type threatParser struct {
}

var defaultThreatParser = &threatParser{}

func (s *threatParser) isStartTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))	
	return threatStartRegexp.Match(line)
}

func (s *threatParser) isEndTag(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	return threatEndRegexp.Match(line)
}


// NewThreatParser return a new BlockParser that parses
// strikethrough expressions.
func NewThreatParser() parser.BlockParser {
	return defaultThreatParser
}

func (s *threatParser) Trigger() []byte {
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
func (s *threatParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	
	if (pos > 0 || !s.isStartTag(line)) {		
		return nil, parser.NoChildren
	}		
	node := NewThreatNode()	
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
func (s *threatParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()	
	pos := pc.BlockOffset()

	if pos == 0 && s.isEndTag(line) {				
		reader.AdvanceLine()		
		return parser.Close
	}	

	if node.Kind() != KindThreat {
		return parser.Close
	}

	text := text.NewSegment(segment.Start, segment.Stop)	
	node.Lines().Append(text)
	
	return parser.Continue | parser.NoChildren
}


func (s *threatParser) CanAcceptIndentedLine() bool {
	return true
}

func (s *threatParser) CanInterruptParagraph() bool {
	return false
}

type ThreatData struct {
	Slug string `yaml:"slug,omitempty"`
	Title string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
}

// Close will be called when the parser returns Close.
func (s *threatParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {	
	d, ok := pc.Get(contextKey).(*data)
	if !ok {
		d = &data{}
	} 
	
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}	
	fmt.Println(buf.String())
	
	var threat ThreatData
	if err := yaml.Unmarshal(buf.Bytes(), &threat); err != nil {
		d.Error = err
	} else {	
		d.Threats = append(d.Threats, threat)
	}

	pc.Set(contextKey, d)
}

