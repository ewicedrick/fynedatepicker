package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// copy .txt files from source to dest based on date range
func copyFilesWithDateFilter(sourceDir, destDir, startDate, endDate string) error {
	//parse the start and end dates to ensure valid format
	start, err := time.Parse("2006/01/02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}
	end, err := time.Parse("2006/01/02", startDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %v", err)
	}

	// create destination directory if not exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	//Process .txt files in the source directory
	err = filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {

		}

		// skip directories and non .txt files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".txt") {
			return nil
		}

		// Process each file
		destfile := filepath.Join(destDir, info.Name())
		err = processFile(path, destfile, start, end)
		if err != nil {
			return fmt.Errorf("error processing file %s: %w", path, err)
		}
		return nil
	})
	return err
}

// processFile processes a single file, filtering its contents based on the date range
func processFile(sourceFile, destFile string, startDate, endDate time.Time) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file %v ", err)
	}

	defer src.Close()

	dest, err := os.Create(destFile)

	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dest.Close()

	scanner := bufio.NewScanner(src)
	var header string
	var foundStart, foundEnd bool

	//copy the heaer (assumes first line is the header)
	if scanner.Scan() {
		header = scanner.Text()
		_, err := dest.WriteString(header + "\n")
		if err != nil {
			return fmt.Errorf("failed to write header to destination file: %v", err)
		}
	}

	// filter lines based on the date range
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}

		// parse the date form the line
		lineDate, err := time.Parse("2006/01/02", parts[0])
		if err != nil {
			return fmt.Errorf("invalid date format in line: %s", line)
		}

		// check if the date is within the range
		if lineDate.After(startDate) && lineDate.Before(endDate) {
			foundStart, foundEnd = true, true
			_, err := dest.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("failed to write line to destination file: %v", err)
			}
		}
	}

	// check for error in scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading source file %v", err)
	}

	// return err if start and end date not found
	if !foundStart || !foundEnd {
		return fmt.Errorf("file %s does not contain data within range: %s or %s", sourceFile, startDate.Format("2006/01/02"), endDate.Format("2006/01/02"))
	}

	return nil
}
