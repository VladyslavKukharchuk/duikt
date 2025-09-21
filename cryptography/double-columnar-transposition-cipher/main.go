package main

import (
	"fmt"
	"strings"
)

type Table struct {
	data [][]rune
	rows int
	cols int
}

func NewTable(plain string, rows, cols int) *Table {
	runes := []rune(plain)

	t := &Table{
		data: make([][]rune, rows),
		rows: rows,
		cols: cols,
	}

	idx := 0
	for r := 0; r < rows; r++ {
		t.data[r] = make([]rune, cols)
		for c := 0; c < cols; c++ {
			if idx < len(runes) {
				t.data[r][c] = runes[idx]
				idx++
			}
		}
	}

	return t
}

func (t *Table) PermuteRows(key []int) {
	newData := make([][]rune, t.rows)

	for oldIdx, newPos := range key {
		newData[oldIdx] = t.data[newPos-1]
	}

	t.data = newData
}

func (t *Table) PermuteCols(key []int) {
	newData := make([][]rune, t.rows)

	for r := 0; r < t.rows; r++ {
		newData[r] = make([]rune, t.cols)

		for oldIdx, newPos := range key {
			newData[r][newPos-1] = t.data[r][oldIdx]
		}
	}

	t.data = newData
}

func (t *Table) ReadByColumns(key []int) string {
	var sb strings.Builder

	t.PermuteCols(key)

	for c := 0; c < t.cols; c++ {
		for r := 0; r < t.rows; r++ {
			sb.WriteRune(t.data[r][c])
		}
	}

	return sb.String()
}

func (t *Table) Print() {
	for _, row := range t.data {
		for _, ch := range row {
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	plain := "перестановки"

	keyRows := []int{3, 1, 2}
	keyCols := []int{4, 1, 3, 2}

	table := NewTable(plain, len(keyRows), len(keyCols))

	fmt.Println("Initial table:")
	table.Print()

	table.PermuteRows(keyRows)

	fmt.Println("After permutations:")
	table.Print()

	cipher := table.ReadByColumns(keyCols)
	fmt.Printf("Cipher text: %s\n", cipher)
}
