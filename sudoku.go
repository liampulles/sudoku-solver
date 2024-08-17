package sudokusolver

import (
	"strings"
)

// ---
// --- High-level Types
// ---

// From 0 to 9 (0 meaning unset).
type Cell uint8

// 1st index is row, 2nd index is column.
type Grid [9][9]Cell

// E.g.
//
// 5 3 _ # _ 7 _ # _ _ _
// 6 _ _ # 1 9 5 # _ _ _
// _ 9 8 # _ _ _ # _ 6 _
// #####################
// 8 _ _ # _ 6 _ # _ _ 3
// 4 _ _ # 8 _ 3 # _ _ 1
// 7 _ _ # _ 2 _ # _ _ 6
// #####################
// _ 6 _ # _ _ _ # 2 8 _
// _ _ _ # 4 1 9 # _ _ 5
// _ _ _ # _ 8 _ # _ 7 9
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
// This won't confirm if the grid is solvable.
func (g Grid) Valid() bool {
	var numsSeenByRow [9][10]bool
	var numsSeenByCol [9][10]bool
	var numsSeenByBox [9][10]bool

	for row, rowCells := range g {
		for col, cell := range rowCells {
			// Empty cells are always valid
			if cell == 0 {
				continue
			}

			// Check row
			if numsSeenByRow[row][cell] {
				return false
			}

			// Check column
			if numsSeenByCol[col][cell] {
				return false
			}

			// Check box
			// -> Which box is it?
			boxRow, boxCol := row/3, col/3
			box := (boxCol * 3) + boxRow
			// -> Ok check now
			if numsSeenByBox[box][cell] {
				return false
			}

			// Ok, mark and continue
			numsSeenByRow[row][cell] = true
			numsSeenByCol[col][cell] = true
			numsSeenByBox[box][cell] = true
		}
	}

	return true
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
