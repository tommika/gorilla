// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

type Datum struct {
	EquatorialRadiusInMeters float64
	PolarRadiusInMeters      float64
	Flattening               float64
	InverseFlattening        float64
}

// WGS84 geo-spatial reference system parameters (constants).
// World Geogetic System. http://en.wikipedia.org/wiki/World_Geodetic_System
func WGS84() Datum {
	wgs84 := Datum{}
	// Earth's equatorial radius in
	wgs84.EquatorialRadiusInMeters = 6378137
	// inverse flattening
	wgs84.InverseFlattening = 298.257223563
	// flattening
	wgs84.Flattening = 1.0 / wgs84.InverseFlattening
	// Earth's polar radius in
	wgs84.PolarRadiusInMeters = wgs84.EquatorialRadiusInMeters * (1 - wgs84.Flattening)
	return wgs84
}
