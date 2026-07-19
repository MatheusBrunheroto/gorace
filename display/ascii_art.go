package display

import (
	"bufio"
	"fmt"
	"gorace/assets"
	"math/rand"
	"strings"
	"unicode/utf8"
)

func findRowSize(art []string) int {
	largest := 0
	for _, a := range art {
		count := utf8.RuneCountInString(a)
		if count > largest {
			largest = count
		}
	}
	return largest
}

func parseLogo(raw string) []string {
	var art []string
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		art = append(art, scanner.Text())
	}
	return art
}

func parseThemes(raw string) [][]string {
	var themes [][]string
	var current []string

	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()

		if line == "+" {
			if len(current) == 3 {
				themes = append(themes, current)
			}
			current = nil
			continue
		}

		current = append(current, line)
	}

	if len(current) == 3 {
		themes = append(themes, current)
	}

	return themes
}

func printAsciiArt(art []string, theme []string) {
	start, end := 0, len(art)-1

	for i, a := range art {
		var line string
		switch i {
		case start:
			line = theme[0] + a
		case end:
			line = theme[2] + a
		default:
			line = theme[1] + a
		}
		fmt.Println(line)
	}
	fmt.Println("")
}

func handleAsciiArt() (int, error) {

	art := parseLogo(assets.Logo)
	themes := parseThemes(assets.Themes)

	r2 := rand.Intn(len(themes))

	largestRow := findRowSize(art) + findRowSize(themes[r2])
	fmt.Printf("\n%s\n\n", strings.Repeat("━", largestRow))
	printAsciiArt(art, themes[r2])

	return largestRow, nil
}
