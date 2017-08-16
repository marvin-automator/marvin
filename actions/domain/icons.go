package domain

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

// MakeSpriteSheet makes an SVG sprite sheet of all the icons of providers and actions.
func MakeSpriteSheet() ([]byte, error) {
	sprites := &svg{
		XMLNS: "http://www.w3.org/2000/svg",
	}

	addProviderIcons(sprites)

	return xml.Marshal(sprites)
}

func addProviderIcons(s *svg) {
	for _, p := range Registry.Providers() {
		icon := unmarshalIcon(p.SVGIcon)
		s.AddSVGAsSymbol(*icon, "provider_" + p.Key)

		for _, g := range Registry.Provider(p.Key).Groups() {
			icon = unmarshalIcon(g.SVGIcon())
			s.AddSVGAsSymbol(*icon, "group_" + g.Name())

			for _, a := range g.Actions() {
				icon = unmarshalIcon(a.SVGIcon)
				s.AddSVGAsSymbol(*icon, "action_" + p.Key + "_" + a.Key)
			}
		}
	}
}

func unmarshalIcon(data []byte) *svg {
	r := regexp.MustCompile(`\<\?xml.*\?\>`)
	data = r.ReplaceAll(data, []byte{})

	s := &svg{}
	err := xml.Unmarshal(data, s)
	if err != nil {
		panic(err)
	}
	return s
}

type svg struct {
	XMLName xml.Name	`xml:"svg"`
	Width	string		`xml:"width,attr,omitempty"`
	Height	string		`xml:"height,attr,omitempty"`
	XMLNS	string		`xml:"xmlns,attr"`
	ViewBox	string		`xml:"viewBox,attr,omitempty"`
	InnerXML string 	`xml:",innerxml"`
	Symbols	 []Symbol
}

func (s *svg) viewBox() string {
	if s.ViewBox != "" {
		return s.ViewBox
	}
	// Generate one based on dimensions.
	return fmt.Sprintf("0 0 %v %v", s.Width, s.Height)
}

func (s *svg) AddSVGAsSymbol(other svg, id string) {
	sym := Symbol{
		ViewBox: other.viewBox(),
		ID: id,
		InnerXML: other.InnerXML,
	}

	s.Symbols = append(s.Symbols, sym)
}

type Symbol struct {
	XMLName	xml.Name	`xml:"symbol"`
	ViewBox	string		`xml:"viewBox,attr"`
	ID		string		`xml:"id,attr"`
	InnerXML string `xml:",innerxml"`
}