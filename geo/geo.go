// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

type LatLon interface {
	IsValid() bool
	Coords() (lat, lon float64)
	ToUTM() (UTM, bool)
}
