// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package gpx

// Summarize a GPX document

import (
	"math"

	"github.com/tommika/gorilla/geo"
	"github.com/tommika/gorilla/must"
)

type Extent struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

func (e *Extent) init() {
	e.MinLat = math.NaN()
	e.MaxLat = math.NaN()
	e.MinLon = math.NaN()
	e.MaxLon = math.NaN()
}

type TrackSummary struct {
	NumPoints    int
	StartPt      Waypoint
	EndPt        Waypoint
	Extent       Extent
	Duration     int64
	MinElePt     Waypoint
	MaxElePt     Waypoint
	EleGain      float64
	Distance     float64
	GridDistance float64
	Area         float64
	Loopiness    float64
	IsClosed     bool
}

func newTrackSummary() *TrackSummary {
	s := TrackSummary{}
	s.MinElePt.Ele = math.NaN() // math.MaxFloat64
	s.MaxElePt.Ele = math.NaN() // -s.MinElevationPt.Ele
	s.EleGain = 0
	s.Extent.init()
	s.IsClosed = true // closed until proved otherwise
	return &s
}

// Summarize computes a summary of the documents GPS data.
func (doc *Document) Summarize() *TrackSummary {
	s := newTrackSummary()
	for i := range doc.Trk {
		s.combine(doc.Trk[i].Summarize())
	}
	return s
}

// Summarize computes a summary of the documents GPS data.
func (trk *Track) Summarize() (summary *TrackSummary) {
	s := newTrackSummary()
	for i := range trk.Trkseg {
		s.combine(trk.Trkseg[i].Summarize())
	}
	return s
}

const (
	eleGainThreshold = 2   // meters
	sClosedThreshold = 100 // meters
)

func (trkseg *TrackSegment) Summarize() *TrackSummary {
	sum := newTrackSummary()
	numPts := len(trkseg.Trkpt)
	if numPts == 0 {
		return sum
	}
	sum.NumPoints += numPts

	sEndpoints := haversineDistance(&trkseg.Trkpt[0], &trkseg.Trkpt[numPts-1])
	sum.IsClosed = sEndpoints < sClosedThreshold

	var trkptPrev *Waypoint
	var trkptElevPrev *Waypoint
	for i := range trkseg.Trkpt {
		trkpt := trkseg.Trkpt[i]
		if !trkpt.IsValid() {
			continue
		}
		if sum.StartPt.Time.IsZero() || trkpt.Time.Before(sum.StartPt.Time) {
			sum.StartPt = trkpt
		}
		if trkpt.Time.After(sum.EndPt.Time) {
			sum.EndPt = trkpt
		}
		// extent
		updateMin(&sum.Extent.MinLat, trkpt.Lat)
		updateMax(&sum.Extent.MaxLat, trkpt.Lat)
		updateMin(&sum.Extent.MinLon, trkpt.Lon)
		updateMax(&sum.Extent.MaxLon, trkpt.Lon)
		// distance
		if trkptPrev != nil {
			sH := haversineDistance(&trkpt, trkptPrev)
			sum.Distance += sH
		}
		if trkptPrev != nil {
			sG := gridDistance(&trkpt, trkptPrev)
			sum.GridDistance += sG
		}
		// elevation
		if !math.IsNaN(trkpt.Ele) {
			if gt(trkpt.Ele, sum.MaxElePt.Ele) {
				sum.MaxElePt = trkpt
			}
			if lt(trkpt.Ele, sum.MinElePt.Ele) {
				sum.MinElePt = trkpt
			}
			// smooth elevation changes by only considering deltas of
			// at least eleThreshold
			if trkptElevPrev == nil {
				trkptElevPrev = &trkpt
			} else {
				d := trkpt.Ele - trkptElevPrev.Ele
				if math.Abs(d) > eleGainThreshold {
					if d > 0 {
						sum.EleGain += d
					}
					trkptElevPrev = &trkpt
				}
			}
		}
		// bottom of loop
		trkptPrev = &trkpt
	}

	sum.Area = 0
	if sum.IsClosed {
		// Compute area using the "Shoelace" method
		// https://en.wikipedia.org/wiki/Shoelace_formula
		must.BeTrue(len(trkseg.Trkpt) > 0)

		// Convert each point to Cartesian coordinates, using the horiz and vert
		// (haversine) distance from a fixed point. The fixed-point is
		// arbitrary, and we'll use the top-left corner of the extent.
		lat0, lon0 := sum.Extent.MinLat, sum.Extent.MinLon
		x0, y0 := trkseg.Trkpt[0].dXdY(lat0, lon0)
		for i := range trkseg.Trkpt {
			j := (i + 1) % numPts
			x1, y1 := trkseg.Trkpt[j].dXdY(lat0, lon0)
			sum.Area += x0 * y1
			sum.Area -= y0 * x1
			x0, y0 = x1, y1
		}
		sum.Area = math.Abs(sum.Area)
	}

	if numPts > 1 {
		sum.Duration = sum.EndPt.Time.Unix() - sum.StartPt.Time.Unix()
	}
	sum.Loopiness = math.Abs(sum.Area / sum.Distance)

	return sum
}

func (sum *TrackSummary) combine(with *TrackSummary) {
	sum.NumPoints += with.NumPoints
	if sum.StartPt.Time.IsZero() || with.StartPt.Time.Before(sum.StartPt.Time) {
		sum.StartPt = with.StartPt
	}
	if with.EndPt.Time.After(sum.EndPt.Time) {
		sum.EndPt = with.EndPt
	}
	// combine extent
	updateMin(&sum.Extent.MinLat, with.Extent.MinLat)
	updateMax(&sum.Extent.MaxLat, with.Extent.MaxLat)
	updateMin(&sum.Extent.MinLon, with.Extent.MinLon)
	updateMax(&sum.Extent.MaxLon, with.Extent.MaxLon)

	// combine distance
	sum.Distance += with.Distance
	sum.GridDistance += with.GridDistance

	// combine elevation
	if gt(with.MaxElePt.Ele, sum.MaxElePt.Ele) {
		sum.MaxElePt = with.MaxElePt
	}
	if lt(with.MinElePt.Ele, sum.MinElePt.Ele) {
		sum.MinElePt = with.MinElePt
	}
	sum.EleGain += with.EleGain

	// combine duration
	sum.Duration += with.Duration

	sum.Area += with.Area
	sEndpoints := haversineDistance(&sum.StartPt, &sum.EndPt)
	sum.IsClosed = sEndpoints < sClosedThreshold
	sum.Loopiness = math.Abs(sum.Area / sum.Distance)
}

func haversineDistance(wp1, wp2 *Waypoint) float64 {
	return geo.HaversineDistance(wp1.Lat, wp1.Lon, wp2.Lat, wp2.Lon)
}

func gridDistance(wp1, wp2 *Waypoint) float64 {
	return geo.GridDistance(wp1.Lat, wp1.Lon, wp2.Lat, wp2.Lon)
}

func updateMax(max *float64, val float64) {
	if gt(val, *max) {
		*max = val
	}
}

func updateMin(min *float64, val float64) {
	if lt(val, *min) {
		*min = val
	}
}

func gt(a, b float64) bool {
	return math.IsNaN(b) || a > b
}

func lt(a, b float64) bool {
	return math.IsNaN(b) || a < b
}

// dXdY computes this points horizontal and vertical distances from
// the given point (lat0, lon0). This can be used to convert the Waypoint
// to Cartesian coordinates relative to some origin.
func (wp *Waypoint) dXdY(lat0, lon0 float64) (dx, dy float64) {
	dx = geo.HaversineDistance(lat0, lon0, lat0, wp.Lon)
	dy = geo.HaversineDistance(lat0, lon0, wp.Lat, lon0)
	return
}
