// Basic script to create file _examples_gallery.md

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const gallerymd = "_examples_gallery.md"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	abs, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	cwd := filepath.Base(abs)
	if cwd != "examples" {
		return fmt.Errorf("you must be in the charts/chartdraw/examples directory")
	}
	gallery, err := os.Create(gallerymd)
	if err != nil {
		return err
	}
	defer gallery.Close()

	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	type example struct {
		dirName   string
		imagePath string
	}
	var examples []example
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		imagePath := filepath.Join(entry.Name(), "output.png")
		_, err := os.Stat(imagePath)
		if err != nil {
			fmt.Printf("%-70s continuing\n", err)
			continue
		}
		examples = append(examples, example{dirName: entry.Name(), imagePath: imagePath})
	}

	fmt.Fprintf(gallery, "# Examples gallery\n<p>\n")
	for _, example := range examples {
		fmt.Fprintf(gallery, "\n## [%s](%s)\n", example.dirName, example.dirName)
		fmt.Fprintf(gallery, "\n![%s](%s)\n<p>\n", example.dirName, example.imagePath)
	}

	fmt.Printf("\nUpdated file %s; remember to commit\n", gallerymd)
	return nil
}
