package tmmd

import (
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark/parser"
)

type ControlData struct {
	Slug string `yaml:"slug" json:"slug"`
	Title string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description"`
	Mitigates []string `yaml:"mitigates,omitempty" json:"mitigates"`
}

var defaultControlParser = newBlockParser("control", KindControl)

// NewControlParser return a new BlockParser that parses
// threat modeling control blocks.
func NewControlParser() parser.BlockParser {
	return defaultControlParser
}

// Get returns a YAML metadata.
func GetControls(pc parser.Context) []ControlData {
	controls := make([]ControlData, 0)
	v := pc.Get(contextKey)
	if v == nil {
		return controls
	}
	d := v.(*TMData)
	l := d.Data[KindControl]
	
	for _, v := range l {
		var control ControlData
		mapstructure.Decode(v, &control)
		controls = append(controls, control)		
	}

	return controls
}






