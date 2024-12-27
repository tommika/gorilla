// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package gpx

// This file defines a model for GPX data

import (
	"time"

	"github.com/tommika/gorilla/geo"
	"github.com/tommika/gorilla/xxml"
)

const (
	GPX_1_1_NS = "http://www.topografix.com/GPX/1/1"
)

var (
	GPX_NAME = xxml.XmlName(GPX_1_1_NS, "gpx")
)

type Document struct {
	Creator string `x:",attr" json:",omitempty"`
	Wpt     []Waypoint
	Trk     []Track
}

type Waypoint struct {
	Lat  float64 `x:",attr"`
	Lon  float64 `x:",attr"`
	Ele  float64
	Time time.Time
	Name *string `json:",omitempty"`
	Desc *string `json:",omitempty"`
}

type Track struct {
	Trkseg []TrackSegment
}

type TrackSegment struct {
	Trkpt []Waypoint
}

// Coords returns the Waypoint coordinates as a tuple
func (wp *Waypoint) Coords() (lat, lon float64) {
	return wp.Lat, wp.Lon
}

// Point returns the waypoint as a simple point
func (wp *Waypoint) Point() geo.Point {
	return geo.Point{Lat: wp.Lat, Lon: wp.Lon}
}

// IsValid determines if the given waypoint is valid
func (wp *Waypoint) IsValid() bool {
	return wp.Point().IsValid()
}
