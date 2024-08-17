package sudokusolver

import (
	"slices"
	"strings"
)

// ---
// --- High-level Types
// ---

// From 0 to 9 (0 meaning unset).
type Cell uint8

type CellGroup [9]Cell

func (c CellGroup) Count() int {
	count := 0
	for _, cell := range c {
		if cell != 0 {
			count++
		}
	}
	return count
}

// 1st index is row, 2nd index is column.
type Grid [9]CellGroup

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
			// -> Ok check it
			if numsSeenByBox[box][cell] {
				return false
			}

			// Ok so far, mark and continue
			numsSeenByRow[row][cell] = true
			numsSeenByCol[col][cell] = true
			numsSeenByBox[box][cell] = true
		}
	}

	return true
}

func (g Grid) RowAt(row int) CellGroup {
	return g[row]
}

func (g Grid) ColumnAt(col int) CellGroup {
	var column CellGroup
	for i, rowCells := range g {
		column[i] = rowCells[col]
	}
	return column
}

func (g Grid) BoxAt(row, col int) CellGroup {
	var box CellGroup
	boxRow, boxCol := row/3, col/3
	i := 0
	for r := boxRow * 3; r < (boxRow+1)*3; r++ {
		for c := boxCol * 3; c < (boxCol+1)*3; c++ {
			box[i] = g[row][col]
			i++
		}
	}
	return box
}

func (g Grid) NumCounts() [10]int {
	var counts [10]int
	for _, rowCells := range g {
		for _, cell := range rowCells {
			counts[cell]++
		}
	}
	return counts
}

type Move struct {
	Row   int
	Col   int
	Value Cell
}

func (g Grid) Apply(m Move) Grid {
	applied := g
	applied[m.Row][m.Col] = m.Value
	return applied
}

type Possibilities []Move

// What possible cells are there for a corresponding grid
// with empty spaces?
func (g Grid) Possibilities() Possibilities {
	var moves []Move

	// Let's look at each empty cell.
	for row, rowCells := range g {
		for col, cell := range rowCells {
			// Ignore filled cells
			if cell != 0 {
				continue
			}

			// Try various options for the cell
			for i := Cell(1); i <= 9; i++ {
				// Make a copy of the grid with the cell filled in
				variant := g
				variant[row][col] = i

				// If it is valid, then mark it as a possibility
				if g.Valid() {
					moves = append(moves, Move{
						Row:   row,
						Col:   col,
						Value: i,
					})
				}
			}
		}
	}

	return moves
}

// Sort the moves based on a subjective ranking.
// The first move is then considered the best.
func (p *Possibilities) Sort(grid Grid) {
	// Calculate and store ranks for each
	type withRank struct {
		Move
		Rank int
	}
	ranks := make([]withRank, len(*p))
	for i, move := range *p {
		ranks[i] = withRank{
			Move: move,
			Rank: move.Rank(grid),
		}
	}

	// Sort by the rank
	slices.SortFunc(ranks, func(a, b withRank) int {
		return a.Rank - b.Rank
	})

	// Now update the possibilities.
	for i, rank := range ranks {
		(*p)[i] = rank.Move
	}
}

// The best move IMO is the one that adds information
// to the most unknown area of a grid.
//
// For example, a move that fills a 9 when all the other
// 9's are known is not very useful. But it is useful
// if very few of the 9's are known.
//
// There is a similar principle for the number of elements
// already given in a row, column, or block.
//
// If the move results in any cell only having one remaining
// possibility, that is given the best ranking.
//
// Lower rank is better.
func (m Move) Rank(grid Grid) int {
	applied := grid.Apply(m)

	// First assess whether the move results in only 1 possibility for a cell.
	moveGrid := groupMoves(applied.Possibilities())
	for _, row := range moveGrid {
		for _, moves := range row {
			if len(moves) == 1 {
				// Nice
				return 0
			}
		}
	}

	// Calculate the number count, row, column, and box constraint reductions.
	rowCount := applied.RowAt(m.Row).Count()
	colCount := applied.ColumnAt(m.Col).Count()
	boxCount := applied.BoxAt(m.Row, m.Col).Count()
	numCount := applied.NumCounts()[m.Value]

	return rowCount + colCount + boxCount + numCount
}

func groupMoves(moves []Move) [9][9][]Move {
	var grouped [9][9][]Move
	for _, move := range moves {
		grouped[move.Row][move.Col] = append(grouped[move.Row][move.Col], move)
	}
	return grouped
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
			// because none of the cell options above worked. I.e. time to backtrack.
			return grid, false
		}
	}

	// If we get here, then it means we started with a filled and valid grid.
	return grid, true
}

var _ Solver = Backtrack

// ---
// --- Hints
// ---

// Produce a hint given an early and later grid. This will select
// the hint that reduces the most possibilities in the early grid
// (so hopefully a "good" hint).
func Hint(early, later Grid) Move {
	// Find the cells that are not in the early grid, but are in
	// the later grid.
	var moves Possibilities
	for row, rowCellsEarly := range early {
		for col, cellEarly := range rowCellsEarly {
			cellLate := later[row][col]
			if cellEarly == 0 && cellLate != 0 {
				moves = append(moves, Move{
					Row:   row,
					Col:   col,
					Value: cellLate,
				})
			}
		}
	}

	// If no moves, fail
	if len(moves) == 0 {
		panic("no moves available - this means early is not a smaller subset of later")
	}

	// Return the best move
	moves.Sort(early)
	return moves[0]
}
