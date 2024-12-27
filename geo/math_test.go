// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import (
	"testing"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/must"
)

func TestDistance(t *testing.T) {
	allowedError := 5.0 // 50 meters
	for _, test := range distanceTestData {
		assert.True(t, test.l1.IsValid())
		assert.True(t, test.l2.IsValid())
		u1 := must.BeOk(test.l1.ToUTM())
		u2 := must.BeOk(test.l2.ToUTM())
		t.Logf("l1: %s, %s", test.l1, u1)
		t.Logf("l2: %s, %s", test.l2, u2)
		lat1, lon1 := test.l1.Coords()
		lat2, lon2 := test.l2.Coords()
		sH := HaversineDistance(lat1, lon1, lat2, lon2)
		t.Logf("hd=%f, expected=%f", sH, test.sH)
		assert.EqualEpsilon(t, test.sH, sH, allowedError)
		sG := GridDistance(lat1, lon1, lat2, lon2)
		t.Logf("gd=%f, expected=%f", sG, test.sG)
		assert.EqualEpsilon(t, test.sG, sG, allowedError)
	}
}
