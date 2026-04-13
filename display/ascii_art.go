package display

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"unicode/utf8"
)

func findLargestRow(art []string) int {

	largest := 0
	for _, a := range art {

		count := utf8.RuneCountInString(a)
		if count > largest {
			largest = count
		}

	}
	return largest
}

func handleAsciiArt() (int, error) {

	var arts [][]string
	var themes [][]string // Three value arrays that dictates how every line from the logo will begin

	if err := readAsciiArt(&arts, "./display/arts/"); err != nil {
		return 0, err
	}
	if err := readAsciiArt(&themes, "./display/themes/"); err != nil {
		return 0, err
	}

	r1 := rand.Intn(len(arts))
	r2 := rand.Intn(len(themes))

	fmt.Println("")
	printAsciiArt(arts[r1], themes[r2])

	return findLargestRow(arts[r1]), nil
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
}

func readAsciiArt(arts *[][]string, path string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {

		var art []string

		file, err := os.Open(path + e.Name())
		if err != nil {
			return errors.New("Unable to open ascii art -> " + path + e.Name())
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			art = append(art, scanner.Text())
		}
		*arts = append(*arts, art)

	}

	return nil
}
