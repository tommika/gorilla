// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package gpx

import (
	"encoding/json"
	"testing"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/must"
)

func logJson(t *testing.T, v any) {
	jsonData := must.NotBeAnError(json.MarshalIndent(v, "", "  "))
	t.Logf("%s\n", string(jsonData))
}

func TestNotFound(t *testing.T) {
	_, err := ReadGpxDocument("bogus")
	assert.NotNil(t, err)
}

var validWaypoints = [...]Waypoint{
	{Lat: 40.75284012047136, Lon: -73.97727928727946},
	{Lat: 41.42011128991017, Lon: -73.95477634666683},
}

func TestWaypoint(t *testing.T) {
	for _, wp := range validWaypoints {
		assert.True(t, wp.IsValid())
		lat, lon := wp.Coords()
		assert.Equal(t, wp.Lat, lat)
		assert.Equal(t, wp.Lon, lon)
	}
}

func TestPointToPoint(t *testing.T) {
	doc, err := ReadGpxDocument("./test-data/BeaconToColdSpring.gpx")
	assert.Nil(t, err)
	s := doc.Summarize()
	logJson(t, s)
	assert.True(t, s.MinElePt.Ele < s.MaxElePt.Ele)
	assert.True(t, s.Extent.MinLat < s.Extent.MaxLat)
	assert.True(t, s.Extent.MinLon < s.Extent.MaxLon)
	assert.EqualEpsilon(t, 12347.0, s.Distance, 1.0)
	assert.EqualEpsilon(t, 12369.0, s.GridDistance, 1.0)
	assert.EqualEpsilon(t, 973.0, s.EleGain, 1.0)
}

func TestLoop(t *testing.T) {
	doc, err := ReadGpxDocument("./test-data/FortHillLoop.gpx")
	assert.Nil(t, err)
	s := doc.Summarize()
	logJson(t, s)
	assert.True(t, s.MinElePt.Ele < s.MaxElePt.Ele)
	assert.True(t, s.Extent.MinLat < s.Extent.MaxLat)
	assert.True(t, s.Extent.MinLon < s.Extent.MaxLon)
	assert.True(t, s.IsClosed)
	assert.EqualEpsilon(t, 3240.0, s.Distance, 1.0)
	assert.EqualEpsilon(t, 3512.0, s.GridDistance, 1.0)
	assert.EqualEpsilon(t, 173.0, s.EleGain, 1.0)
}

func TestSummaryWithBadTrack_TODO(t *testing.T) {
}
