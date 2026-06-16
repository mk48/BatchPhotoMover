package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: BatchPhotoMover <input-folder> <output-folder>")
		os.Exit(1)
	}

	inputFolder := os.Args[1]
	outputFolder := os.Args[2]

	err := filepath.Walk(inputFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error accessing %s: %v\n", filePath, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		photoTaken, err := getPhotoTakenDateTime(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %s: %v\n", filePath, err)
			return nil
		}

		destDir := filepath.Join(outputFolder, photoTaken)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directory %s: %v\n", destDir, err)
			return nil
		}

		dest := filepath.Join(destDir, info.Name())
		fmt.Printf("%s --> %s\n", filePath, dest)

		if err := os.Rename(filePath, dest); err != nil {
			fmt.Fprintf(os.Stderr, "failed to move %s: %v\n", filePath, err)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("--> Completed <--")
}

func getPhotoTakenDateTime(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", fileName, err)
	}
	defer f.Close()

	exifInfo, err := exif.Decode(f)
	if err != nil {
		return fileModifiedDate(fileName)
	}

	t, err := exifInfo.DateTime()
	if err != nil {
		return fileModifiedDate(fileName)
	}

	return t.Format("2006-01-02"), nil
}

func fileModifiedDate(fileName string) (string, error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", fileName, err)
	}
	return info.ModTime().Format("2006-01-02"), nil
}
