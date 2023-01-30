package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/barasher/go-exiftool"
)

var verison = "v0.0.3"

var versionFlag bool
var sourceFlag string
var targetFlag string
var modtimeFlag bool

var exiftoolInstance *exiftool.Exiftool

func main() {
	flag.BoolVar(&versionFlag, "version", false, "version info")
	flag.StringVar(&sourceFlag, "source", "", "source path")
	flag.StringVar(&targetFlag, "target", "", "target path")
	flag.BoolVar(&modtimeFlag, "modtime", false, "modification time fallback")
	flag.Parse()

	if versionFlag {
		fmt.Printf("Version: %s\n", verison)
		os.Exit(0)
	}

	if sourceFlag == "" || targetFlag == "" {
		fmt.Println("Error: --source and --target are required flags.")
		fmt.Println("Usage: img-sort --source /path/to/source --target /path/to/target")
		os.Exit(1)
	}

	// Create exiftool instance
	var exiftoolErr error
	exiftoolInstance, exiftoolErr = exiftool.NewExiftool()
	if exiftoolErr != nil {
		fmt.Println(exiftoolErr)
		os.Exit(1)
	}
	defer exiftoolInstance.Close()

	// Recursively read source directory
	processErr := filepath.Walk(sourceFlag, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, files only
		if fileInfo.IsDir() {
			return nil
		}

		// Allow only these file extensions
		if !isExtAllowed(path, allowedExtensions) {
			return nil
		}

		// Decode file exif data and parse create date
		var fileDate time.Time
		var fileError error
		fileExif, fileError := decodeExif(path)
		fileDate, fileError = parseExifCreateDate(fileExif)
		if fileError != nil {
			fileDate = fileInfo.ModTime()

			if !modtimeFlag {
				return moveFileToUnknown(path, targetFlag)
			}
		}

		return moveFileToTarget(path, targetFlag, fileDate)
	})

	if processErr != nil {
		fmt.Println(processErr)
		os.Exit(1)
	}

	os.Exit(0)
}

func parseExifCreateDate(fileExif exiftool.FileMetadata) (time.Time, error) {
	var fileDate time.Time
	var fileDateErr error
	for _, exifField := range exifDateFields {
		if fileDate, fileDateErr = parseDate(fileExif.Fields[exifField], commonDateFormats); fileDateErr == nil {
			return fileDate, nil
		}
	}

	return time.Time{}, errors.New("Could not parse exif create date")
}

func moveFileToUnknown(path, targetRoot string) error {
	newPath := filepath.Join(targetRoot, "unknown", filepath.Base(path))
	return moveFile(path, newPath)
}

func moveFileToTarget(path string, targetRoot string, fileDate time.Time) error {
	yearDir := fmt.Sprintf("%d", fileDate.Year())
	monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
	fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), filepath.Ext(path))
	newPath := filepath.Join(targetRoot, yearDir, monthDir, fileName)
	return moveFile(path, newPath)
}
