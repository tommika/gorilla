// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/must"
)

const goriallaNS = "https://github.com/tommika/gorilla"

var testDataXml = `
	<?xml version="1.0"?>
	<data 
		xmlns="https://github.com/tommika/gorilla"
		xmlns:types="https://github.com/tommika/gorilla/types"
		>	
		Wow
		<desc xml:space="preserve">
			This is my data.
			Do you like it?
		</desc>
		<valid>
			<!-- comment -->
			true
			<ignored/>
		</valid>
		<id>
			2112
		</id>
		<id2>2112</id2>
		<p1>
			<x>1</x>
			<y>2</y>
			<z>3</z>
		</p1>
		<p2>
			<x>10</x>
			<y>20</y>
			<z>30</z>
		</p2>
		<p3>
			<x>100</x>
			<y>200</y>
			<z>300</z>
		</p3>
		<!-- sequence of structs -->
		<point>
			<x>1.0</x>
			<y>1.1</y>
			<z>1.2</z>
		</point>
		<point>
			<x>2.0</x>
			<y>3.1</y>
			<z>3.2</z>
		</point>
		<ignore>
			<ignore/>
		</ignore>
		<!-- sequence of ints -->
		<num>1</num>
		<num>2</num>
		<num>3</num>
		<!-- element with cdata -->
		<types:name types:id="2112">
			Rush
		</types:name>
		<start>2024-11-01T02:15:23+6:00</start>
		<start>2024-11-01T02:15:23Z</start>
	</data>
	`

type testData struct {
	Desc    string
	Valid   bool
	Id      int
	Id2     uint
	P1      point[int] `x:"p1,elem"`
	P2      *point[int]
	P3      point[int]       `x:"-"` // ignored
	Points  []point[float32] `x:"point"`
	Numbers []int            `x:"num"`
	Name    name             `x:",,https://github.com/tommika/gorilla/types"`
	Data    *string          `x:",cdata"` // set only if there's cdata
	Start   *time.Time
	End     *time.Time
}

type point[T any] struct {
	X T // defaults to elem 'x'
	Y T `x:"y"` // defaults to elem
	Z T `x:"z,elem"`
}

type name struct {
	_    struct{} `xmlns:"https://github.com/tommika/gorilla/types"`
	Id   int      `x:",attr"`
	Name string   `x:",cdata"`
}

// stringReader creates an io.Reader for the given string.

func stringReader(s string) io.Reader {
	return bytes.NewReader([]byte(s))
}

func logJson(t *testing.T, v any) {
	jsonData := must.NotBeAnError(json.MarshalIndent(v, "", "  "))
	t.Logf("%s\n", string(jsonData))
}

func TestReadXml(t *testing.T) {
	rootName := XmlName(goriallaNS, "data")
	d := testData{}
	err := ReadXmlWithRootName(stringReader(testDataXml), &d, rootName)
	assert.Nil(t, err)
	logJson(t, d)
	assert.True(t, d.Valid)
	assert.Equal(t, 2112, d.Id)
	assert.NotNil(t, d.P2)
	assert.Equal(t, 3, len(d.Numbers))
}

func TestReadXmlWithWrongRoot(t *testing.T) {
	d := testData{}
	err := ReadXmlWithRootName(stringReader(testDataXml), &d, XmlName(goriallaNS, "fred"))
	assert.NotNil(t, err)
	t.Logf("%+v\n", err)
}

func TestReadInvalidXml(t *testing.T) {
	d := testData{}
	err := Unmarshal([]byte{}, &d)
	assert.NotNil(t, err)
	t.Logf("%+v\n", err)
}

func TestUnmarshal(t *testing.T) {
	d := testData{}
	err := Unmarshal([]byte(testDataXml), &d)
	assert.Nil(t, err)
	t.Logf("%+v\n", d)
}

func TestReadXmlWithWrongValueType(t *testing.T) {
	err := ReadXml(stringReader(testDataXml), testData{})
	assert.NotNil(t, err)
	t.Logf("%+v\n", err)

	v := 5
	err = ReadXml(stringReader(testDataXml), &v)
	assert.NotNil(t, err)
	t.Logf("%+v\n", err)
}

var eventsXml = `
	<?xml version="1.0"?>
	<events>	
		<description>
			This is a series of events
		</description>
		<event id="1">
			<time>2111</time>
			Event 1
		</event>
		<event id="2">
			<time>2112</time>
			Event 3
		</event>
		<event id="2">
		 	<time>2112</time>
			Event 3
		</event>
	</events>
`

type eventSequence struct {
	Description string
	Event       func(env *event)
}

type event struct {
	Id   int `x:",attr"`
	Time int
	Data string `x:",cdata"`
}

func TestReadXmlWithCallbackSequence(t *testing.T) {
	count := 0
	events := eventSequence{}
	events.Event = func(env *event) {
		if count == 0 {
			t.Logf("%s", events.Description)
		}
		t.Logf("got an event: %+v", env)
		count++
	}
	err := ReadXml(stringReader(eventsXml), &events)
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
}

type badEventSequence struct {
	Description string
	Event       func(env event) // wrong signature
}

func TestReadXmlWithBadCallbackSequence(t *testing.T) {
	count := 0
	events := badEventSequence{}
	events.Event = func(env event) {
		if count == 0 {
			t.Logf("%s", events.Description)
		}
		t.Logf("got an event: %+v", env)
		count++
	}
	err := ReadXml(stringReader(eventsXml), &events)
	assert.NotNil(t, err)
	t.Logf("%+v\n", err)
}
