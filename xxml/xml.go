// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
// The xxml package reads XML-formatted data files using tagged Go structs.
package xxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

type Name = xml.Name

func XmlName(space, local string) Name {
	return Name{
		Space: space,
		Local: local,
	}
}

// Unmarshal wraps the parser with Go's standard method signature for
// marshalling data.
func Unmarshal(data []byte, val any) error {
	return ReadXml(bytes.NewReader(data), val)
}

func ReadXml(r io.Reader, val any) error {
	return ReadXmlWithRootName(r, val, xml.Name{})
}

// ReadXml reads an XML document from the given input stream into the given value.
// The given value must be a pointer to a struct.
func ReadXmlWithRootName(r io.Reader, val any, rootElemName xml.Name) error {
	// Make sure the given value is of the correct type. Has to be a pointer
	// to a struct.
	t := reflect.ValueOf(val)
	if !(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("val must be a pointer to a struct")
	}
	theStruct := t.Elem()

	d := xml.NewDecoder(r)
	// skip cdata, comments, processing instructions, etc at the start of the XML file
	se, err := findStartElement(d)
	if err == nil {
		if rootElemName != (xml.Name{}) && se.Name != rootElemName {
			// wrong root element name
			err = fmt.Errorf("unexpected start element: expected %s; got %s", rootElemName, se.Name)
		} else {
			err = parseElement(d, "", se, theStruct)
		}
	}
	return err
}

// parseElement implements a recursive-descent parse an XML document.
// The switch statement here covers the different Go types that can be
// mapped from XML values.
func parseElement(d *xml.Decoder, defNS string, se xml.StartElement, fieldValue reflect.Value) error {
	var err error
	switch fieldValue.Kind() {
	case reflect.Slice:
		// parse the element as new item in a sequence
		elemType := fieldValue.Type().Elem()
		newVal := reflect.New(elemType).Elem()
		if err = parseElement(d, defNS, se, newVal); err == nil {
			fieldValue.Set(reflect.Append(fieldValue, newVal))
		}
	case reflect.Pointer:
		// allocate a new value and parse the element into it
		elemType := fieldValue.Type().Elem()
		newVal := reflect.New(elemType)
		if err = parseElement(d, defNS, se, newVal.Elem()); err == nil {
			fieldValue.Set(newVal)
		}
	case reflect.Func:
		// treat the func as a callback for an item in a sequence (this is pretty
		// cool!)
		// require that the function have a single parameter, and that it be a
		// pointer to a value
		t := fieldValue.Type()
		if !(t.NumIn() == 1 && t.In(0).Kind() == reflect.Pointer) {
			err = fmt.Errorf("inappropriate callback func: %s", t)
		} else {
			elemType := t.In(0).Elem()
			newVal := reflect.New(elemType)
			if err = parseElement(d, defNS, se, newVal.Elem()); err == nil {
				// invoke the callback
				fieldValue.Call([]reflect.Value{newVal})
			}
		}
	default:
		if shouldParseAsStruct(fieldValue) {
			err = parseElementAsStruct(d, defNS, se, fieldValue)
		} else {
			// treat the element as cdata
			var cdata xml.CharData
			if cdata, err = parseElementAsCdata(d, se); err == nil {
				cdataString := string(cdata)
				if !attrValEquals(se, xmlAttrSpace, xmlAttrSpacePreserve) {
					cdataString = strings.TrimSpace(cdataString)
				}
				setValueFromCdata(fieldValue, cdataString)
			}
		}
	}
	return err
}

func shouldParseAsStruct(fieldValue reflect.Value) bool {
	if fieldValue.Kind() != reflect.Struct {
		return false
	}
	switch fieldValue.Interface().(type) {
	case time.Time:
		return false
	}
	return true
}

// parseElementAsStruct performs a recursive descent parse of the tokens returned from the xml.Decoder
func parseElementAsStruct(d *xml.Decoder, xmlns string, se xml.StartElement, theStruct reflect.Value) error {
	// Handle attributes
	var xmlnsAttr *string
	for _, a := range se.Attr {
		if a.Name == xmlAttrNS {
			xmlnsAttr = &a.Value
		} else if attrField := findFieldForXmlElement(theStruct, "", a.Name, xmlFieldAttr); attrField != nil {
			setValueFromCdata(*attrField, a.Value)
		}
	}
	if xmlnsAttr != nil {
		xmlns = *xmlnsAttr
	}
	// Accumulate cdata
	var cdata xml.CharData
	// Loop until we read the end token
	var err error
	done := false
	for !done && err == nil {
		var tok any
		tok, err = d.Token()
		switch e := tok.(type) {
		case xml.EndElement:
			cdataString := string(cdata)
			if !attrValEquals(se, xmlAttrSpace, xmlAttrSpacePreserve) {
				cdataString = strings.TrimSpace(cdataString)
			}
			if len(cdataString) > 0 {
				// If there's  is a cdata field, set it
				if cdataField := findFieldForXmlElement(theStruct, xmlns, xml.Name{}, xmlFieldCdata); cdataField != nil {
					setValueFromCdata(*cdataField, cdataString)
				}
			}
			done = true
		case xml.StartElement:
			// See if there's a field for this XML element
			fieldValue := findFieldForXmlElement(theStruct, xmlns, e.Name, xmlFieldElem)
			if fieldValue == nil {
				// No field; eat the element
				err = parseElementAsVoid(d, e)
			} else {
				err = parseElement(d, xmlns, e, *fieldValue)
			}
		case xml.CharData:
			cdata = append(cdata, e...)
		}
	}
	return err
}

func parseElementAsCdata(d *xml.Decoder, _ xml.StartElement) (xml.CharData, error) {
	var cdata xml.CharData
	depth := 1
	var err error
	for depth > 0 && err == nil {
		var tok any
		tok, err = d.Token()
		switch e := tok.(type) {
		case xml.EndElement:
			depth -= 1
		case xml.StartElement:
			depth += 1
		case xml.CharData:
			cdata = append(cdata, e...)
		}
	}
	return cdata, err
}

func parseElementAsVoid(d *xml.Decoder, _ xml.StartElement) error {
	depth := 1
	var err error
	for depth > 0 && err == nil {
		var tok any
		tok, err = d.Token()
		switch tok.(type) {
		case xml.EndElement:
			depth -= 1
		case xml.StartElement:
			depth += 1
		}
	}
	return err
}
