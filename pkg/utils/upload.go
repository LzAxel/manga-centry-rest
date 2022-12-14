package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ImageExtansions = "png webp jpg gif"
)

func sortFiles(reader *zip.ReadCloser) {
	sort.Slice(reader.File, func(i, j int) bool {
		iStr := strings.Split(reader.File[i].Name, ".")[0]
		jStr := strings.Split(reader.File[j].Name, ".")[0]
		comp := strings.Compare(iStr, jStr)
		if comp == -1 {
			return true
		} else {
			return false
		}
	})
}

func validateImageExtansion(filename string) error {
	fileNameSplit := strings.Split(filename, ".")
	fileExtansion := fileNameSplit[len(filename)-1]
	if !strings.Contains(ImageExtansions, fileExtansion) {
		return fmt.Errorf("invalid file extansion: %s", fileExtansion)
	}
	return nil
}

func UnzipImages(archivePath, outputPath string) ([]string, error) {
	var files []string

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return files, err
	}
	defer reader.Close()

	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		return files, err
	}
	logrus.Debugln("sorting archive")
	sortFiles(reader)

	logrus.Debugln("extracting files from archive")
	for i, f := range reader.File {
		if err := validateImageExtansion(f.Name); err != nil {
			continue
		}
		fileName := fmt.Sprintf("%s%s", strconv.Itoa(i), filepath.Ext(f.Name))
		filePath := filepath.Join(outputPath, fileName)
		outputFilePath, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return files, err
		}
		defer outputFilePath.Close()

		zippedFile, err := f.Open()
		if err != nil {
			return files, err
		}
		defer zippedFile.Close()

		if _, err := io.Copy(outputFilePath, zippedFile); err != nil {
			return files, err
		}
		files = append(files, fileName)
	}

	return files, nil
}
