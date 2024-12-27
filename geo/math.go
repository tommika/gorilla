// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import (
	"math"
)

const (
	R_EARTH_KM = 6371.009 // Mean radius of The Earth
)

// Hsin returns the haversine of the given angle a in degrees.
func Hsin(a float64) float64 {
	t := math.Sin(a / 2)
	return t * t
}

// D2R converts degrees to radians
func D2R(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// HaversineDistance computes the haversine distance ("great-circle" distance)
// between two points on the earth.
// See: http://en.wikipedia.org/wiki/Haversine_formula
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1R := D2R(lat1)
	lat2R := D2R(lat2)
	dLatR := D2R(lat2 - lat1)
	dLonR := D2R(lon2 - lon1)

	a := Hsin(dLatR) + Hsin(dLonR)*math.Cos(lat1R)*math.Cos(lat2R)
	s := 2 * R_EARTH_KM * math.Asin(math.Sqrt(a))
	return s * 1000 // Distance in meters
}

// GridDistance returns the grid distance between two points on The Earth.
func GridDistance(lat1, lon1, lat2, lon2 float64) float64 {
	return math.Sqrt(GridDistanceSq(lat1, lon1, lat2, lon2))
}

// GridDistance returns the squared grid distance between two points on The Earth.
func GridDistanceSq(lat1, lon1, lat2, lon2 float64) float64 {
	utm1 := ToUTM(lat1, lon1)
	utm2 := ToUTM(lat2, lon2)
	dx := float64(utm2.Easting - utm1.Easting)
	dy := float64(utm2.Northing - utm1.Northing)
	return dx*dx + dy*dy
}
