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
	archive := reader.File
	logrus.Debugln("sorting archive")
	sort.Slice(archive, func(i, j int) bool {
		iStr := strings.Split(reader.File[i].Name, ".")[0]
		jStr := strings.Split(reader.File[j].Name, ".")[0]
		comp := strings.Compare(iStr, jStr)
		if comp == -1 {
			return true
		} else {
			return false
		}
	})

	logrus.Debugln("extracting files from archive")
	for i, f := range archive {

		fileNameSplit := strings.Split(f.Name, ".")
		fileExtansion := fileNameSplit[len(fileNameSplit)-1]
		if !strings.Contains("png webp jpg gif", fileExtansion) {
			logrus.Debugf("invalid file in archive: %s", f.Name)
			continue
		}

		fileName := fmt.Sprintf("%s.%s", strconv.Itoa(i), fileExtansion)
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
