package tmmd

import (
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark/parser"
)

type ThreatData struct {
	Slug string `yaml:"slug" json:"slug"`
	Title string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description"`
}

var defaultThreatParser = newBlockParser("threat", KindThreat)

// NewThreatParser return a new BlockParser that parses
// strikethrough expressions.
func NewThreatParser() parser.BlockParser {
	return defaultThreatParser
}


// Returns a list of threats
func GetThreats(pc parser.Context) []ThreatData {
	threats := make([]ThreatData, 0)	
	v := pc.Get(contextKey)
	if v == nil {
		return threats
	}
	d := v.(*TMData)
	l := d.Data[KindThreat]	
	for _, v := range l {
		var threat ThreatData
		mapstructure.Decode(v, &threat)
		threats = append(threats, threat)		
	}
	return threats
}



