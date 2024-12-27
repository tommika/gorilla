// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import "fmt"

// Point is a point on the surface of The Earth represented using latitude and
// longitude coordinates in degrees
type Point struct {
	Lat float64
	Lon float64
}

func (pt Point) String() string {
	return fmt.Sprintf("%f, %f", pt.Lat, pt.Lon)
}

// IsValid determines if the coordinates are valid
func (pt Point) IsValid() (valid bool) {
	valid = pt.Lat >= -80 &&
		pt.Lat <= 84 &&
		pt.Lon >= -180 &&
		pt.Lon < 180
	return
}

// Coords returns the coordinates as a float64 tuple
func (pt Point) Coords() (lat, lon float64) {
	return pt.Lat, pt.Lon
}

// ToUTM converts the coordinates to UTM
func (pt Point) ToUTM() (UTM, bool) {
	if !pt.IsValid() {
		return UTM{}, false
	}
	return ToUTM(pt.Coords()), true
}
