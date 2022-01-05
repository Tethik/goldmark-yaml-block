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
	Threats []tmmd.ThreatData
	Controls []tmmd.ControlData
}

func main() {
	fn := os.Args[1]

	source, err := ioutil.ReadFile(fn)
	if (err != nil) {
		panic(err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, tmmd.ThreatModelingExtension, meta.Meta),
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


	// if Get(context).Error != nil {
	// 	panic(Get(context).Error)
	// }

	var output Output
	output.Meta = metaData
	output.Threats = tmmd.GetThreats(context)
	output.Controls = tmmd.GetControls(context)

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