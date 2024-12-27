// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package geo

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestUTM(t *testing.T) {
	for i := range locationTestData {
		loc := locationTestData[i]
		t.Logf("%s: (%s) %si\n", loc.name, loc.pt, loc.utm)
		utm, ok := loc.pt.ToUTM()
		t.Logf("utm=%s", utm)
		assert.True(t, ok)
		assert.True(t, utm.IsClose(loc.utm, 5))
	}
}

func TestInvalid(t *testing.T) {
	pt := Point{Lat: 0, Lon: 180}
	_, ok := pt.ToUTM()
	assert.False(t, ok)
}
