package display

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"unicode/utf8"
)

//go:embed arts/*
var artsFS embed.FS

//go:embed themes/*
var themesFS embed.FS

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

func readAsciiArt(arts *[][]string, fsys embed.FS, path string) error {
	entries, err := fsys.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		var art []string

		file, err := fsys.Open(path + "/" + e.Name())
		if err != nil {
			return errors.New("Unable to open ascii art -> " + path + "/" + e.Name())
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			art = append(art, scanner.Text())
		}
		*arts = append(*arts, art)
	}
	return nil
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
	var arts [][]string
	var themes [][]string

	if err := readAsciiArt(&arts, artsFS, "arts"); err != nil {
		return 0, err
	}
	if err := readAsciiArt(&themes, themesFS, "themes"); err != nil {
		return 0, err
	}

	r1 := rand.Intn(len(arts))
	r2 := rand.Intn(len(themes))
	largestRow := findRowSize(arts[r1]) + findRowSize(themes[r2])

	fmt.Printf("\n%s\n\n", strings.Repeat("━", largestRow))
	printAsciiArt(arts[r1], themes[r2])

	return largestRow, nil
}
