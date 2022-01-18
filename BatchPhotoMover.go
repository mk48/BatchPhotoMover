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
					photoTaken, err := getPhotoTakenDateTime(inputFileFolderPath)

					if err == nil && photoTaken != "" {
						outputFolderFullPath := filepath.FromSlash(path.Join(outputFolder, photoTaken))

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
						fmt.Printf("\nErr\n")
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

func getPhotoTakenDateTime(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("\nError to open file: %s\n%s\n\n", fileName, err)
		return "", err
	}

	exifInfo, err := exif.Decode(f)
	defer f.Close()

	if err != nil {
		//fmt.Printf("\nError to get Exif info, so using file time: %s\n%s\n\n", fileName, err)
		return getFileModifiedDateTime(fileName)
	} else {
		exifInfoDateTime, err := exifInfo.DateTime()
		if err != nil {
			fmt.Printf("\nError in retrieving the photo taken time -> exifInfo.DateTime(): %s\n%s\n\n", fileName, err)
			return getFileModifiedDateTime(fileName)
		}
		return exifInfoDateTime.Format("2006-01-02"), nil
	}
}

//get file modified datetime
func getFileModifiedDateTime(fileName string) (string, error) {
	info, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("\nError in getting file modified time: %s\n%s\n\n", fileName, err)
		return "", err
	}

	fmt.Println("File time")
	dt := info.ModTime().Format("2006-01-02")
	return dt, nil
}
