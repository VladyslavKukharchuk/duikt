package main

import (
	"fmt"
	"sort"
	"strings"
)

var uaAlphabet = []rune("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя")

type Grille struct {
	size int
	data [][]rune
}

func NewGrille(size int) *Grille {
	g := &Grille{
		size: size,
		data: make([][]rune, size),
	}

	for i := 0; i < size; i++ {
		g.data[i] = make([]rune, size)
	}

	return g
}

func rotate(x, y, n int) (int, int) {
	return y, n - 1 - x
}

func (g *Grille) Fill(plain string, mask []string) {
	pt := []rune(plain)
	ptIdx := 0
	alphaIdx := 0

	// collect base hole coordinates from mask (0°)
	baseHoles := make([][2]int, 0)
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if mask[i][j] == '1' {
				baseHoles = append(baseHoles, [2]int{i, j})
			}
		}
	}

	// iterate 4 rotations
	for rot := 0; rot < 4; rot++ {
		rotCoords := make([][2]int, 0, len(baseHoles))
		for _, c := range baseHoles {
			x, y := c[0], c[1]
			for k := 0; k < rot; k++ {
				x, y = rotate(x, y, g.size)
			}
			rotCoords = append(rotCoords, [2]int{x, y})
		}

		sort.Slice(rotCoords, func(i, j int) bool {
			if rotCoords[i][0] == rotCoords[j][0] {
				return rotCoords[i][1] < rotCoords[j][1]
			}
			return rotCoords[i][0] < rotCoords[j][0]
		})

		for _, rc := range rotCoords {
			x, y := rc[0], rc[1]

			if g.data[x][y] != 0 {
				continue
			}
			if ptIdx < len(pt) {
				g.data[x][y] = pt[ptIdx]
				ptIdx++
			} else {
				g.data[x][y] = uaAlphabet[alphaIdx%len(uaAlphabet)]
				alphaIdx++
			}
		}
	}
}

func (g *Grille) ReadRows() string {
	var sb strings.Builder

	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			sb.WriteRune(g.data[i][j])
		}
	}

	return sb.String()
}

func (g *Grille) Print() {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			fmt.Printf("%c ", g.data[i][j])
		}
		fmt.Println()
	}

	fmt.Println()
}

func main() {
	plain := "перестановка"

	mask := []string{
		"0010",
		"0001",
		"0100",
		"1000",
	}

	grille := NewGrille(4)
	grille.Fill(plain, mask)

	fmt.Println("Filled grille:")

	grille.Print()

	fmt.Printf("Cipher text: %s\n", grille.ReadRows())
}
