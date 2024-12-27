// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import (
	"fmt"
	"testing"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/must"
)

// Some test data
// TODO - move into a data file

type location struct {
	name string
	pt   Point
	utm  UTM
}

func parseUTM(s string) UTM {
	var u UTM
	n, err := fmt.Sscanf(s, "%d%c %dmE %dmN", &u.Zone, &u.Hemi, &u.Easting, &u.Northing)
	must.BeNil(err)
	must.BeEqual(4, n)
	return u
}

var locationTestData = [...]location{
	{"Grand Central Terminal, NY, NY", Point{40.75284012047136, -73.97727928727946}, parseUTM("18N 586333mE 4511823mN")},
	{"Cold Spring, NY", Point{41.42011128991017, -73.95477634666683}, parseUTM("18N 587345mE 4585922mN")},
	{"Elephant Island", Point{-61.10881764684619, -55.14362567900252}, parseUTM("21S 600050mE 3223672mN")},
}

type distance struct {
	l1 Point
	l2 Point
	sH float64 // haversine distance in meters between l1 and l2
	sG float64 // grid distance in meters between l1 and l2
}

var distanceTestData = [...]distance{
	// distances between GCT to CSP
	{locationTestData[0].pt, locationTestData[1].pt, 74220.0, 74105},
}

func TestData(t *testing.T) {
	for i := range locationTestData {
		loc := locationTestData[i]
		fmt.Printf("%s - (%s) %s\n", loc.name, loc.pt, loc.utm)
		assert.True(t, loc.pt.IsValid())
	}
}
