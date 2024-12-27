// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import (
	"fmt"
	"math"

	"github.com/tommika/gorilla/must"
	"github.com/tommika/gorilla/util"
)

type Hemisphere = rune

const (
	SOUTH Hemisphere = 'S'
	NORTH Hemisphere = 'N'
)

// UTM is a location on earth using the UTM representation
type UTM struct {
	Hemi     Hemisphere
	Zone     int
	Easting  int64
	Northing int64
}

func (u UTM) String() string {
	return fmt.Sprintf("%d%c %dmE %dmN", u.Zone, u.Hemi, u.Easting, u.Northing)
}

func (u UTM) IsClose(other UTM, epsilon int64) bool {
	return u.Hemi == other.Hemi &&
		u.Zone == other.Zone &&
		util.AbsInt(u.Easting-other.Easting) < epsilon &&
		util.AbsInt(u.Northing-other.Northing) < epsilon
}

// ToUTM converts the given coordinates to UTM.  The coordinates MUST be valid,
// and it is up to the caller to ensure that they are.
func ToUTM(degLat, degLon float64) (utm UTM) {
	// https://en.wikipedia.org/wiki/Universal_Transverse_Mercator_coordinate_system

	zone := 1 + int((math.Floor((degLon + 180) / 6.0)))
	must.BeTrue(zone >= 1 && zone <= 60) // FIXME - do this in debug only

	datum := WGS84()

	// Ellipsoid parameters. These are constant for a given datum.
	// REVIEW - move these into the datum class, so we don't have to re-calc on every call.
	// random double between [min,max)

	a := datum.EquatorialRadiusInMeters
	b := a * (1.0 - datum.Flattening)
	esq := (1.0 - (b/a)*(b/a)) // first-eccentricity squared
	e2sq := esq / (1.0 - esq)  // second-eccentricity squared

	// Determine central meridian of given zone.
	// Each zone is 6 degs wide with the central meridian being 3 degs from the
	// edge of the zone.
	degLon0 := 3.0 + 6.0*float64(zone-1) - 180.0

	// convert input location degrees to radians
	lat := D2R(degLat)
	lon := D2R(degLon)
	lon0 := D2R(degLon0)

	// pre-calculate some common/intermediate terms used in formulas.
	// these terms depend on the input location
	p := lon - lon0
	p3 := p * p * p
	sinlat := math.Sin(lat)
	sin2lat := sinlat * sinlat
	tanlat := math.Tan(lat)
	tan2lat := tanlat * tanlat
	coslat := math.Cos(lat)
	cos2lat := coslat * coslat
	cos3lat := coslat * coslat * coslat
	nu := a / math.Sqrt(1-(esq*sin2lat))

	// UTM Parameters
	k0 := 0.9996 // scale factor along central meridian

	// compute easting
	K4 := k0 * nu * math.Cos(lat)
	K5 := ((k0 * nu * cos3lat) / 6.0) * (1 - tan2lat + e2sq*cos2lat)
	easting := (K4 * p) + (K5 * p3)

	// compute northing
	N := nu
	T := tan2lat
	C := e2sq * cos2lat
	A := p * coslat
	M := lat * (1.0 - esq*(1.0/4.0+esq*(3.0/64.0+5.0*esq/256.0)))
	M = M - math.Sin(2.0*lat)*(esq*(3.0/8.0+esq*(3.0/32.0+45.0*esq/1024.0)))
	M = M + math.Sin(4.0*lat)*(esq*esq*(15.0/256.0+esq*45.0/1024.0))
	M = M - math.Sin(6.0*lat)*(esq*esq*esq*(35.0/3072.0))
	M = M * a //Arc length along standard meridian
	M0 := 0.0
	northing := k0 * (M - M0 + N*tanlat*(A*A*(1.0/2.0+A*A*((5-T+9*C+4*C*C)/24.0+A*A*(61-58*T+T*T+600*C-330*e2sq)/720.0)))) //Northing from equator

	utm.Zone = zone
	utm.Easting = int64(500000 + easting)
	if northing >= 0 {
		// northern hemisphere
		utm.Hemi = NORTH
		utm.Northing = int64(northing)
	} else {
		// southern hemisphere; adjust for false northing
		utm.Hemi = SOUTH
		utm.Northing = int64(northing + 10000000)
	}

	return
}
