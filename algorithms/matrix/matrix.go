// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package matrix

import (
	"fmt"
	"io"
	"strings"

	"github.com/tommika/gorilla/util"
)

type Matrix[T any] struct {
	numRows int
	numCols int
	data    []T
	rows    [][]T
}

func NewMatrix[T any](numRows, numCols int, defVal T) Matrix[T] {
	m := Matrix[T]{
		numRows: numRows,
		numCols: numCols,
	}
	// Allocate as a contiguous block of memory to hold the matrix
	m.data = make([]T, numRows*numCols)
	// Create row slices from the data
	m.rows = make([][]T, numRows)
	for r, i := 0, 0; r < numRows; r, i = r+1, i+numCols {
		m.rows[r] = m.data[i : i+numCols]
	}
	m.SetAll(defVal)
	return m
}

func (m Matrix[T]) Rows() [][]T {
	return m.rows
}

func (m Matrix[T]) NumRows() int {
	return m.numRows
}

func (m Matrix[T]) NumCols() int {
	return m.numCols
}

func (m Matrix[T]) SetAll(val T) {
	for i := range m.data {
		m.data[i] = val
	}
}

func (m Matrix[T]) Get(row, col int) T {
	return m.rows[row][col]
}

func (m Matrix[T]) Set(row, col int, val T) {
	m.rows[row][col] = val
}

func (m Matrix[T]) Fprintf(out io.Writer) {
	for _, row := range m.rows {
		for _, cell := range row {
			fmt.Fprintf(out, "%v\t", cell)
		}
		fmt.Fprintf(out, "\n")
	}
}

type CellParser[T any] func(string) (T, error)

// ParseMatrix is a very tolerant matrix parser.
func ParseMatrix[T any](s string, parse CellParser[T]) (m Matrix[T]) {
	// Split the input on new lines, ensuring that we handle both unix and dos
	// style line endings
	rowsT := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	rows := []string{}
	maxCols := 0
	// Make a first pass of the input to understand the shape of the matrix
	for _, row := range rowsT {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			// skip blank rows
			continue
		}
		rows = append(rows, row)
		cols := strings.Fields(row)
		if len(cols) > maxCols {
			maxCols = len(cols)
		}
	}
	// Now that we know the shape, construct the matrix, and
	// make a second pass to parse the cells.
	m = NewMatrix[T](len(rows), maxCols, util.Zero[T]())
	for r, row := range rows {
		cells := strings.Fields(row)
		for c, cell := range cells {
			v, err := parse(cell)
			if err == nil {
				m.Set(r, c, v)
			}
		}
	}
	return
}
