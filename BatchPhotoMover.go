package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	inputFolder := "C:\\Temp\\ToUpload\\From Asus Tab\\ph"
	//inputFolder := os.Args[1]
	//outputFolder := os.Args[2]

	err := filepath.Walk(inputFolder,
		func(path string, info os.FileInfo, err error) error {
			if err == nil {
				if !info.IsDir() {
					//process only files
					f, err := os.Open(path)
					if err == nil {
						exifInfo, err := exif.Decode(f)
						if err == nil {
							exifInfoDateTime, err := exifInfo.DateTime()
							if err == nil {
								fmt.Printf("%s --> %s\n", path, exifInfoDateTime)
							} else {
								fmt.Printf("Error to get DateTime: %s\n%s\n", path, err)
							}
						} else {
							fmt.Printf("Error to get Exif info: %s\n%s\n", path, err)
						}
					} else {
						fmt.Printf("Error to open file: %s\n%s\n", path, err)
					}
				}
			} else {
				fmt.Printf("Error processing folder: %s\n%s\n", path, err)
			}
			return nil
		})

	if err != nil {
		panic(err)
	}
}

func mainOld() {

	// key   --> 2020-01-31
	// value --> c:\**\***\***\2020-01-31
	folderFullPaths := make(map[string]string)

	dirToProcess := os.Args[1]
	//dirToProcess := "C:\\Users\\mumu-in\\OneDrive - FLSmidth\\K\\Photos\\2020-01 Poland\\Mobile"

	files, err := ioutil.ReadDir(dirToProcess)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()
		currentFileLocationPath := filepath.FromSlash(path.Join(dirToProcess, fileName))

		//fmt.Printf(filePath)

		fileStat, err := os.Stat(currentFileLocationPath)
		if err != nil {
			log.Fatal(err)
		}

		fileDateTime := fileStat.ModTime().Format("2006-01-02")

		if _, found := folderFullPaths[fileDateTime]; !found {

			newFolderFullPath := filepath.FromSlash(path.Join(dirToProcess, fileDateTime))

			// create folder
			_ = os.Mkdir(newFolderFullPath, os.ModePerm)

			folderFullPaths[fileDateTime] = newFolderFullPath

			//fmt.Printf("%s\n", newFolder);
		}

		newLocationToMove := filepath.FromSlash(path.Join(folderFullPaths[fileDateTime], fileName))

		fmt.Printf("%s ==> %s\n\n", currentFileLocationPath, newLocationToMove)
		// move file into this folder
		errInMove := os.Rename(currentFileLocationPath, newLocationToMove)
		if errInMove != nil {
			log.Fatal(errInMove)
		}
	}
}
