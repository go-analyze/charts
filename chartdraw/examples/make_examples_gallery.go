// Basic script to create file _examples_gallery.md

package main

import (
	"errors"
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
		return errors.New("you must be in the charts/chartdraw/examples directory")
	}
	gallery, err := os.Create(gallerymd)
	if err != nil {
		return err
	}
	defer func() { _ = gallery.Close() }()

	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	type example struct {
		dirName   string
		imagePath string
	}
	examples := make([]example, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		imagePath := filepath.Join(entry.Name(), "output.png")
		_, err := os.Stat(imagePath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%-70s continuing\n", err)
			continue
		}
		examples = append(examples, example{dirName: entry.Name(), imagePath: imagePath})
	}

	_, _ = fmt.Fprintf(gallery, "# Examples gallery\n<p>\n")
	for _, example := range examples {
		_, _ = fmt.Fprintf(gallery, "\n## [%s](%s)\n", example.dirName, example.dirName)
		_, _ = fmt.Fprintf(gallery, "\n![%s](%s)\n<p>\n", example.dirName, example.imagePath)
	}

	_, _ = fmt.Fprintf(os.Stderr, "\nUpdated file %s; remember to commit\n", gallerymd)
	return nil
}
