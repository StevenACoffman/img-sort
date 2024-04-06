package exif

import (
	"errors"
	"time"

	"github.com/barasher/go-exiftool"

	"github.com/StevenACoffman/img-sort/internal/config"
	"github.com/StevenACoffman/img-sort/internal/log"
	"github.com/StevenACoffman/img-sort/internal/util"
)

var (
	exif *exiftool.Exiftool
	err  error
)

func init() {
	exif, err = exiftool.NewExiftool()
	if err != nil {
		log.Error().Msg("Failed to initialize exiftool")
		panic(err)
	}
}

func Instance() *exiftool.Exiftool {
	return exif
}

func ExtractData(paths ...string) []exiftool.FileMetadata {
	return exif.ExtractMetadata(paths...)
}

func ParseDate(fileExif exiftool.FileMetadata, exifFields []string) (time.Time, error) {
	for _, field := range exifFields {
		date, err := util.TryParseDate(fileExif.Fields[field], config.EXIF_FIELDS_DATE_FORMAT)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.New("failed to parse exif date")
}
