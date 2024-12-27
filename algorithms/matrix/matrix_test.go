// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package matrix

import (
	"os"
	"strconv"
	"testing"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/assert"
)

func TestUnitMatrix(t *testing.T) {
	m := NewMatrix(0, 0, types.Nothing)
	assert.Equal(t, 0, m.NumRows())
	assert.Equal(t, 0, m.NumCols())
}

func TestMatrix(t *testing.T) {
	defVal := 2112
	m := NewMatrix(3, 4, defVal)
	assert.Equal(t, 3, m.NumRows())
	assert.Equal(t, 4, m.NumCols())
	for _, row := range m.Rows() {
		for _, cell := range row {
			assert.Equal(t, defVal, cell)
		}
	}
}

func TestMatrixCellAccess(t *testing.T) {
	m := NewMatrix(3, 4, 0)
	i := 0
	for r := 0; r < m.NumRows(); r++ {
		for c := 0; c < m.NumCols(); c++ {
			assert.Equal(t, 0, m.Get(r, c))
			m.Set(r, c, i)
			i += 1
		}
	}
	i = 0
	for r := 0; r < m.NumRows(); r++ {
		for c := 0; c < m.NumCols(); c++ {
			assert.Equal(t, i, m.Get(r, c))
			i += 1
		}
	}
}

func TestParseIntMatrix(t *testing.T) {
	input := `
		0 1 2 

		4 5 6 7`
	m := ParseMatrix(input, strconv.Atoi)
	assert.Equal(t, 2, m.NumRows())
	assert.Equal(t, 4, m.NumCols())
	m.Fprintf(os.Stderr)
}
