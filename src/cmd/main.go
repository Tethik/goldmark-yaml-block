package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Tethik/threat-modeling-markdown/src/tmmd"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type Output struct {
	Meta map[string]interface{}
	Threats []*ThreatData
	Controls []*ControlData
}

type ThreatData struct {
	Slug string `yaml:"slug" json:"slug"`
	Title string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description"`
}

type ControlData struct {
	Slug string `yaml:"slug" json:"slug"`
	Title string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description"`
	Mitigates []string `yaml:"mitigates,omitempty" json:"mitigates"`
}

func NewControlData() interface{} {
	return &ControlData{}
}

func NewThreatData() interface{} {
	return &ThreatData{}
}



func main() {
	fn := os.Args[1]

	source, err := ioutil.ReadFile(fn)
	if (err != nil) {
		panic(err)
	}

	ext := tmmd.CreateThreatModelingExtension("control", "threat")

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, ext, meta.Meta),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	context := parser.NewContext()
	md.Parser().Parse(text.NewReader(source), parser.WithContext(context))
	
	metaData := meta.Get(context)

	var output Output
	output.Meta = metaData
	items := tmmd.GetItems("threat", NewThreatData, context)	 
	for _, item := range items {
		output.Threats = append(output.Threats, item.(*ThreatData))
	}
	items = tmmd.GetItems("control", NewControlData, context)
	for _, item := range items {
		output.Controls = append(output.Controls, item.(*ControlData))
	}
	

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	
	fmt.Println()
	fmt.Println(buf.String())
	fmt.Println()

	b, _ := json.MarshalIndent(output, "", "    ")
	fmt.Println(string(b))
}