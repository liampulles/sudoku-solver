package sudokusolver

import (
	"strings"
)

// ---
// --- High-level Types
// ---

// From 0 to 9 (0 meaning unset).
type Cell uint8

// Could be a row, column, or box.
type CellGroup [9]Cell

// Should not have duplicates of cells 1 to 9.
func (c CellGroup) Valid() bool {
	var numSeen [10]bool
	for _, cell := range c {
		if cell == 0 {
			continue
		}

		if numSeen[cell] {
			return false
		}

		numSeen[cell] = true
	}
	return true
}

// 1st index is row, 2nd index is column.
type Grid [9]CellGroup

// E.g.
//
// 53_#_7_#___
// 6__#195#___
// _98#___#_6_
// ###########
// 8__#_6_#__3
// 4__#8_3#__1
// 7__#_2_#__6
// ###########
// _6_#___#23_
// ___#419#__5
// ___#_8_#_79
func (g Grid) String() string {
	var w strings.Builder
	for row, rowCells := range g {
		if row%3 == 0 && row > 0 {
			w.WriteString("#####################\n")
		}
		for col, cell := range rowCells {
			if col%3 == 0 && col > 0 {
				w.WriteString("# ")
			}

			if cell == 0 {
				w.WriteByte('_')
			} else {
				w.WriteByte('0' + byte(cell))
			}

			w.WriteByte(' ')
		}
		w.WriteByte('\n')
	}
	return w.String()
}

// Must pass row, column, and box constraints.
func (g Grid) Valid() bool {
	// Check row constraint
	for _, row := range g.Rows() {
		if !row.Valid() {
			return false
		}
	}

	// Check column constraints
	for _, col := range g.Columns() {
		if !col.Valid() {
			return false
		}
	}

	// Check box constraints
	for _, box := range g.Boxes() {
		if !box.Valid() {
			return false
		}
	}

	return true
}

func (g Grid) Rows() [9]CellGroup {
	return g
}

// Effectively transposes the grid.
func (g Grid) Columns() [9]CellGroup {
	var columns [9]CellGroup
	g.Loop(func(row, col int, cell Cell) {
		columns[col][row] = cell
	})
	return columns
}

func (g Grid) Boxes() [9]CellGroup {
	var boxes [9]CellGroup
	g.Loop(func(row, col int, cell Cell) {
		// Which box is it?
		boxRow, boxCol := row/3, col/3
		box := (boxCol * 3) + boxRow

		// Which item in the box is it?
		posRow, posCol := row-(boxRow*3), col-(boxCol*3)
		pos := (posCol * 3) + posRow

		boxes[box][pos] = cell
	})
	return boxes
}

func (g Grid) Loop(fn func(row, col int, cell Cell)) {
	for row, rowCells := range g {
		for col, cell := range rowCells {
			fn(row, col, cell)
		}
	}
}

// ---
// --- Solving
// ---

// Produce a solved grid, given a partial grid. Returns false if
// the grid is not solvable.
type Solver func(Grid) (Grid, bool)

// Solve via a traditional backtracking depth-first search.
func Backtrack(grid Grid) (Grid, bool) {
	// If the input grid is not valid, stop.
	if !grid.Valid() {
		return grid, false
	}

	// Let's try solve each unsolved cell.
	for row, rowCells := range grid {
		for col, cell := range rowCells {
			// Ignore filled cells
			if cell != 0 {
				continue
			}

			// Try various options for the cell
			for i := Cell(1); i <= 9; i++ {
				// Make a copy of the grid with the cell filled in
				variant := grid
				variant[row][col] = i

				// Try solve that variant
				filled, solved := Backtrack(variant)

				// If solved, then we are done.
				if solved {
					return filled, true
				}

				// Else, we'll try another number...
			}

			// We have a partially filled, non-constraint violating grid which is nonetheless invalid,
			// because none of the cell options above worked.
			return grid, false
		}
	}

	// If we get here, then it means we started with a filled and valid grid.
	return grid, true
}

var _ Solver = Backtrack
