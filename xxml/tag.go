// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

// Functions for handling xml tags on structs

import (
	"encoding/xml"
	"reflect"
	"strings"

	"github.com/tommika/gorilla/util"
)

// findFieldForXmlElement searches the given structure for a field that matches
// the given xml name and kind.
func findFieldForXmlElement(structVal reflect.Value, xmlns string, eName xml.Name, kind xmlFieldKind) *reflect.Value {
	if nsFromStruct := getDefXmlNsForStruct(structVal); len(nsFromStruct) > 0 {
		xmlns = nsFromStruct
	}
	t := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		field := t.Field(i)
		tag, ok := getXmlTag(field, xmlns)
		// TODO: build and cache the field map rather than doing this on each element
		if ok && tag.name == eName && tag.kind == kind {
			fieldVal := structVal.Field(i)
			return &fieldVal
		}
	}
	// not found
	return nil
}

func getDefXmlNsForStruct(structVal reflect.Value) (defNS string) {
	t := structVal.Type()
	if f, ok := t.FieldByName("_"); ok {
		val, ok := f.Tag.Lookup("xmlns")
		if ok {
			defNS = val
		}
	}
	return
}

const (
	xmlTagKey = "x"
)

type xmlFieldKind int

const (
	xmlFieldElem xmlFieldKind = iota
	xmlFieldAttr
	xmlFieldCdata
)

type xmlFieldTag struct {
	name xml.Name
	kind xmlFieldKind
}

func defXmlName(field reflect.StructField) string {
	return strings.ToLower(field.Name)
}

func getXmlTag(field reflect.StructField, defNS string) (tag xmlFieldTag, ok bool) {
	// x:"name,kind,namespace"
	var val string
	tag.name.Space = defNS
	if val, ok = field.Tag.Lookup(xmlTagKey); !ok {
		// use default mapping for all exported fields
		if field.IsExported() {
			tag.name.Local = defXmlName(field)
			tag.kind = xmlFieldElem
			ok = true
		}
	} else if val == "-" {
		// ignore this field
		ok = false
	} else {
		valSplit := util.SplitAndTrim(val, ",")
		splitLen := len(valSplit)
		if splitLen > 0 {
			// name
			tag.name.Local = valSplit[0]
			if len(tag.name.Local) == 0 {
				tag.name.Local = defXmlName(field)
			}
		}
		if splitLen > 1 {
			kind := strings.ToLower(valSplit[1])
			if kind == "cdata" {
				tag.kind = xmlFieldCdata
				tag.name = xml.Name{}
			} else if kind == "attr" {
				tag.kind = xmlFieldAttr
			} else {
				tag.kind = xmlFieldElem
			}
		}
		if splitLen > 2 {
			// namespace
			tag.name.Space = valSplit[2]
		}
	}
	return
}
