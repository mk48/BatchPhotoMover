package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	//inputFolder := "C:\\Temp\\ToUpload\\From Asus Tab"
	inputFolder := os.Args[1]

	//outputFolder := "C:\\Temp\\ToUpload\\n"
	outputFolder := os.Args[2]

	err := filepath.Walk(inputFolder,
		func(inputFileFolderPath string, info os.FileInfo, err error) error {
			if err == nil {
				if !info.IsDir() {
					//process only files
					f, err := os.Open(inputFileFolderPath)
					if err == nil {
						exifInfo, err := exif.Decode(f)
						f.Close()
						if err == nil {
							exifInfoDateTime, err := exifInfo.DateTime()
							if err == nil {
								//fmt.Printf("%s --> %s\n", path, exifInfoDateTime)
								exifInfoDateTimeForPath := exifInfoDateTime.Format("2006-01-02")
								outputFolderFullPath := filepath.FromSlash(path.Join(outputFolder, exifInfoDateTimeForPath))

								if _, err := os.Stat(outputFolderFullPath); os.IsNotExist(err) {
									os.Mkdir(outputFolderFullPath, os.ModeDir)
									fmt.Printf("\nNew folder created: %s\n", outputFolderFullPath)
								}

								//move the file
								newLocationToMove := filepath.FromSlash(path.Join(outputFolderFullPath, info.Name()))
								fmt.Printf("%s --> %s\n", inputFileFolderPath, newLocationToMove)

								errInMove := os.Rename(inputFileFolderPath, newLocationToMove)
								if errInMove != nil {
									fmt.Printf("Error to move: %s\n%s\n", inputFileFolderPath, errInMove)
								}

							} else {
								fmt.Printf("\nError to get DateTime: %s\n%s\n\n", inputFileFolderPath, err)
							}
						} else {
							fmt.Printf("\nError to get Exif info: %s\n%s\n\n", inputFileFolderPath, err)
						}
					} else {
						fmt.Printf("\nError to open file: %s\n%s\n\n", inputFileFolderPath, err)
					}
				}
			} else {
				fmt.Printf("\nError processing folder: %s\n%s\n\n", inputFileFolderPath, err)
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	fmt.Println("--> Completed <--")
}
