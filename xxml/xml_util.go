// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

import (
	"encoding/xml"
	"fmt"
)

var (
	xmlNS        = "http://www.w3.org/XML/1998/namespace"
	xmlAttrSpace = xml.Name{
		Space: xmlNS,
		Local: "space",
	}
	xmlAttrNS = xml.Name{
		Local: "xmlns",
	}
	xmlAttrSpacePreserve = "preserve"
)

func findStartElement(d *xml.Decoder) (xml.StartElement, error) {
	for {
		tok, err := d.Token()
		if tok == nil || err != nil {
			err = fmt.Errorf("start element not found: %v", err)
			return xml.StartElement{}, err
		}
		if se, ok := tok.(xml.StartElement); ok {
			return se, nil
		}
	}
}

// attrVal returns the value of the given attribute,
// or nil if not found
func attrVal(se xml.StartElement, name xml.Name) *string {
	for _, a := range se.Attr {
		if a.Name == name {
			return &a.Value
		}
	}
	return nil
}

// attrValEquals determines if the given attribute has value that
// is equal to the given value
func attrValEquals(se xml.StartElement, name xml.Name, val string) bool {
	v := attrVal(se, name)
	return v != nil && *v == val
}
