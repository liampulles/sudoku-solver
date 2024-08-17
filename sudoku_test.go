package sudokusolver_test

import (
	"testing"

	sudokusolver "github.com/liampulles/sudoku-solver"
	"github.com/stretchr/testify/assert"
)

var partial = sudokusolver.Grid{
	{5, 3, 0, 0, 7, 0, 0, 0, 0},
	{6, 0, 0, 1, 9, 5, 0, 0, 0},
	{0, 9, 8, 0, 0, 0, 0, 6, 0},
	{8, 0, 0, 0, 6, 0, 0, 0, 3},
	{4, 0, 0, 8, 0, 3, 0, 0, 1},
	{7, 0, 0, 0, 2, 0, 0, 0, 6},
	{0, 6, 0, 0, 0, 0, 2, 8, 0},
	{0, 0, 0, 4, 1, 9, 0, 0, 5},
	{0, 0, 0, 0, 8, 0, 0, 7, 9},
}

var filled = sudokusolver.Grid{
	{5, 3, 4, 6, 7, 8, 9, 1, 2},
	{6, 7, 2, 1, 9, 5, 3, 4, 8},
	{1, 9, 8, 3, 4, 2, 5, 6, 7},
	{8, 5, 9, 7, 6, 1, 4, 2, 3},
	{4, 2, 6, 8, 5, 3, 7, 9, 1},
	{7, 1, 3, 9, 2, 4, 8, 5, 6},
	{9, 6, 1, 5, 3, 7, 2, 8, 4},
	{2, 8, 7, 4, 1, 9, 6, 3, 5},
	{3, 4, 5, 2, 8, 6, 1, 7, 9},
}

func TestGrid_String(t *testing.T) {
	tests := []struct {
		desc     string
		fixture  sudokusolver.Grid
		expected string
	}{
		{
			"empty",
			sudokusolver.Grid{},
			`_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
#####################
_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
#####################
_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
_ _ _ # _ _ _ # _ _ _ 
`,
		},
		{
			"partial",
			partial,
			`5 3 _ # _ 7 _ # _ _ _ 
6 _ _ # 1 9 5 # _ _ _ 
_ 9 8 # _ _ _ # _ 6 _ 
#####################
8 _ _ # _ 6 _ # _ _ 3 
4 _ _ # 8 _ 3 # _ _ 1 
7 _ _ # _ 2 _ # _ _ 6 
#####################
_ 6 _ # _ _ _ # 2 8 _ 
_ _ _ # 4 1 9 # _ _ 5 
_ _ _ # _ 8 _ # _ 7 9 
`,
		},
		{
			"filled",
			filled,
			`5 3 4 # 6 7 8 # 9 1 2 
6 7 2 # 1 9 5 # 3 4 8 
1 9 8 # 3 4 2 # 5 6 7 
#####################
8 5 9 # 7 6 1 # 4 2 3 
4 2 6 # 8 5 3 # 7 9 1 
7 1 3 # 9 2 4 # 8 5 6 
#####################
9 6 1 # 5 3 7 # 2 8 4 
2 8 7 # 4 1 9 # 6 3 5 
3 4 5 # 2 8 6 # 1 7 9 
`,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.fixture.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCellGroup_Valid(t *testing.T) {
	tests := []struct {
		desc     string
		fixture  sudokusolver.CellGroup
		expected bool
	}{
		{
			"empty is valid",
			sudokusolver.CellGroup{},
			true,
		},
		{
			"partial valid",
			partial[0],
			true,
		},
		{
			"partial invalid",
			sudokusolver.CellGroup{0, 1, 0, 0, 2, 0, 0, 3, 1},
			false,
		},
		{
			"full valid",
			filled[0],
			true,
		},
		{
			"full invalid",
			sudokusolver.CellGroup{1, 2, 3, 4, 5, 6, 7, 8, 1},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.fixture.Valid()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestGrid_Valid(t *testing.T) {
	tests := []struct {
		desc     string
		fixture  sudokusolver.Grid
		expected bool
	}{
		{
			"empty is valid",
			sudokusolver.Grid{},
			true,
		},
		{
			"partial valid",
			partial,
			true,
		},
		{
			"filled valid",
			filled,
			true,
		},
		{
			"partial invalid (by row)",
			sudokusolver.Grid{
				{5, 3, 0, 0, 7, 0, 0, 0, 0},
				{6, 0, 0, 1, 9, 5, 0, 0, 0},
				{0, 9, 8, 0, 0, 0, 0, 6, 0},
				{8, 0, 0, 0, 6, 0, 0, 0, 3},
				{4, 0, 0, 8, 0, 3, 0, 0, 1},
				{7, 0, 0, 0, 2, 0, 0, 0, 6},
				{0, 6, 0, 0, 0, 0, 6, 8, 0}, // -> Double 6 here
				{0, 0, 0, 4, 1, 9, 0, 0, 5},
				{0, 0, 0, 0, 8, 0, 0, 7, 9},
			},
			false,
		},
		{
			"partial invalid (by column)",
			sudokusolver.Grid{
				{5, 3, 0, 0, 7, 0, 0, 0, 1}, // -> Double 1 in last column
				{6, 0, 0, 1, 9, 5, 0, 0, 0},
				{0, 9, 8, 0, 0, 0, 0, 6, 0},
				{8, 0, 0, 0, 6, 0, 0, 0, 3},
				{4, 0, 0, 8, 0, 3, 0, 0, 1},
				{7, 0, 0, 0, 2, 0, 0, 0, 6},
				{0, 6, 0, 0, 0, 0, 2, 8, 0},
				{0, 0, 0, 4, 1, 9, 0, 0, 5},
				{0, 0, 0, 0, 8, 0, 0, 7, 9},
			},
			false,
		},
		{
			"partial invalid (by box)",
			sudokusolver.Grid{
				{5, 3, 0, 0, 7, 0, 0, 0, 0},
				{6, 0, 0, 1, 9, 5, 0, 0, 0},
				{0, 9, 8, 0, 0, 0, 0, 6, 0},
				{8, 0, 0, 0, 6, 0, 0, 0, 3},
				{4, 0, 0, 8, 0, 3, 0, 0, 1},
				{7, 0, 0, 0, 2, 0, 0, 1, 6}, // -> Double 1 in this last box
				{0, 6, 0, 0, 0, 0, 2, 8, 0},
				{0, 0, 0, 4, 1, 9, 0, 0, 5},
				{0, 0, 0, 0, 8, 0, 0, 7, 9},
			},
			false,
		},
		{
			"filled invalid (by row or column or box)",
			sudokusolver.Grid{
				{5, 3, 4, 6, 7, 8, 9, 1, 2},
				{6, 7, 2, 1, 9, 5, 3, 4, 8},
				{1, 9, 8, 3, 4, 2, 5, 6, 7},
				{8, 5, 9, 7, 6, 1, 4, 2, 3},
				{4, 2, 6, 8, 5, 3, 7, 9, 1},
				{7, 1, 3, 9, 2, 4, 8, 5, 6},
				{9, 6, 1, 5, 3, 7, 2, 8, 4},
				{2, 8, 7, 4, 1, 9, 6, 3, 5},
				{3, 4, 5, 2, 8, 6, 1, 7, 3}, // -> Double 3 at the end here
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.fixture.Valid()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestBacktrack_Solvable(t *testing.T) {
	actualGrid, actualSolved := sudokusolver.Backtrack(partial)

	assert.Equal(t, actualGrid, filled)
	assert.Equal(t, actualSolved, true)
}
