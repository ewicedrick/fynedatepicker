package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	datepicker "github.com/sdassow/fyne-datepicker"
)

var StartDate string
var EndDate string

// Source and Destination Test
var SourceFolderGlob string
var DestinationFolderGlob string

var ResultsFolderGlob string

func main() {
	a := app.NewWithID("cm.pad.cedrickewi")
	w := a.NewWindow("Demo")

	CtrlAltS := &desktop.CustomShortcut{fyne.KeyS, fyne.KeyModifierControl | fyne.KeyModifierAlt}
	w.Canvas().AddShortcut(CtrlAltS, func(_ fyne.Shortcut) {
		makeScreenshot(w)
	})

	// DATE PICKER ENTRY
	startDate := widget.NewEntry()
	startDate.SetPlaceHolder("0000/00/00")
	startDate.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		if startDate.Text != "" {
			t, err := time.Parse("2006/01/02", startDate.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				startDate.SetText(when.Format("2006/01/02"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date",
			"Ok",
			"Cancel",
			datepicker,
			datepicker.OnActioned,
			w,
		)
	})

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder("0000/00/00")
	endDate.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		if endDate.Text != "" {
			t, err := time.Parse("2006/01/02", endDate.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				endDate.SetText(when.Format("2006/01/02"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date",
			"Ok",
			"Cancel",
			datepicker,
			datepicker.OnActioned,
			w,
		)
	})
	// END OF DATE PICKER ENTRY

	// select folder
	sourceFolder := selectSourceFolder(w, "Select Source Folder")
	destinationFolder := selectDestinationFolder(w, "Select Destination Folder")
	resultFolder := selectResultDestination(w, "Select folder to display Results")

	// Confirm Data Transfer
	btnTransferData := widget.NewButton("Copy Data", func() {
		err := scanAndCopy(SourceFolderGlob, DestinationFolderGlob, w)
		if err != nil {
			fmt.Println("", err)
			return
		}
	})

	go func() {
		for {
			time.Sleep(10 * time.Minute)

			if SourceFolderGlob != "" && DestinationFolderGlob != "" {
				err := scanAndCopy(SourceFolderGlob, DestinationFolderGlob, w)
				if err != nil {
					fmt.Println("", err)
					return
				}
			} else {
				fmt.Println("No Folders Selected")
			}
		}

	}()

	// CONFIRM FILE FILTER BY DATE RANGE
	btnFilterData := widget.NewButton("Filter", func() {
		fmt.Println("Start Date: ", startDate.Text, "\nEndDate", endDate.Text, "\n")
		err := copyFilesWithDateFilter(DestinationFolderGlob, ResultsFolderGlob, startDate.Text, endDate.Text)
		if err != nil {
			fmt.Println("", err)
			return
		} else {
			fmt.Println("Files processed successfully")
		}

	})

	label := widget.NewLabelWithStyle("Select Time Range", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	info := "Select Source and Destination Folders"
	labelfolder := widget.NewLabelWithStyle(info, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	w.SetContent(container.NewVBox(
		// layout.NewBorderLayout(label, nil, nil, nil),
		labelfolder,
		widget.NewForm(
			widget.NewFormItem("Source", sourceFolder),
			widget.NewFormItem("Destination", destinationFolder),
			widget.NewFormItem("", layout.NewSpacer()),
			widget.NewFormItem("", btnTransferData),
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
		label,
		widget.NewForm(
			widget.NewFormItem("Start", startDate),
			widget.NewFormItem("End", endDate),
			widget.NewFormItem("", layout.NewSpacer()),
			widget.NewFormItem("Results", resultFolder),
			// widget.NewFormItem("Dest", resultFolder),
			widget.NewFormItem("", layout.NewSpacer()),
			widget.NewFormItem("", btnFilterData),
		),
	))

	//

	w.Resize(fyne.Size{
		Width:  640,
		Height: 530,
	})

	w.ShowAndRun()
}

// Function to Select folders
func selectSourceFolder(w fyne.Window, usage string) fyne.CanvasObject {
	// Entry do display selected folder path
	folderPathEntry := widget.NewEntry()
	folderPathEntry.SetPlaceHolder(usage)
	var word string
	// Button to open folder selection dialog
	selectFolderButton := widget.NewButton("Select Folder", func() {

		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				folderPathEntry.SetText("Error selecting folder")
				return
			}

			if uri != nil {
				folderPathEntry.SetText(uri.Path()) //Display the path in the entry
				SourceFolderGlob = folderPathEntry.Text
			}
		}, w)
		fmt.Println("", word)
	})
	containerText := container.NewMax(folderPathEntry)
	containerText.Resize(fyne.NewSize(300, 100))

	// display
	content := container.NewVBox(
		containerText,
		layout.NewSpacer(),
		selectFolderButton,
	)

	return content
}

// Select Destination Folder
func selectDestinationFolder(w fyne.Window, usage string) fyne.CanvasObject {
	// Entry do display selected folder path
	folderPathEntry := widget.NewEntry()
	folderPathEntry.SetPlaceHolder(usage)
	var word string
	// Button to open folder selection dialog
	selectFolderButton := widget.NewButton("Select Folder", func() {

		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				folderPathEntry.SetText("Error selecting folder")
				return
			}

			if uri != nil {
				folderPathEntry.SetText(uri.Path()) //Display the path in the entry
				DestinationFolderGlob = folderPathEntry.Text
			}
		}, w)
		fmt.Println("", word)
	})

	// display
	content := container.NewVBox(
		folderPathEntry,
		layout.NewSpacer(),
		selectFolderButton,
	)

	return content
}

// select the folder for store result files
func selectResultDestination(w fyne.Window, usage string) fyne.CanvasObject {
	// Entry do display selected folder path
	folderPathEntry := widget.NewEntry()
	folderPathEntry.SetPlaceHolder(usage)
	var word string
	// Button to open folder selection dialog
	selectFolderButton := widget.NewButton("Select Folder", func() {

		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				folderPathEntry.SetText("Error selecting folder")
				return
			}

			if uri != nil {
				folderPathEntry.SetText(uri.Path()) //Display the path in the entry
				ResultsFolderGlob = folderPathEntry.Text
			}
		}, w)
		fmt.Println("", word)
	})

	// display
	content := container.NewVBox(
		folderPathEntry,
		layout.NewSpacer(),
		selectFolderButton,
	)

	return content
}

// scan a directory for .txt files and process them
func scanAndCopy(sourceDir, destDir string, w fyne.Window) error {
	if sourceDir == "" || destDir == "" {
		dialog.ShowInformation("Unable to copy files", "Make Sure Folders are Selected", w)
		return nil
	}

	err := filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//skip directories
		if info.IsDir() {
			return nil
		}

		//process only .txt files
		if strings.HasSuffix(info.Name(), ".txt") {
			destFilePath := filepath.Join(destDir, info.Name())
			err := copyOrAppendFile(path, destFilePath)
			if err != nil {
				dialog.ShowInformation("Error Copying Files", "Error Encountered copy files call cedrick", w)
				return fmt.Errorf("error processing file %s:%w", path, &err)
			}
		}
		return nil
	})

	return err
}

// copy or append new lines to destination file.
func copyOrAppendFile(sourceFile, destFile string) error {

	//open source file
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	defer src.Close()

	//if dest file does not exist, create and copy content
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		dest, err := os.Create(destFile)
		if err != nil {
			return err
		}
		defer dest.Close()
		_, err = io.Copy(dest, src)
		return err
	}

	// if destination file exists, append only new lines
	dest, err := os.OpenFile(destFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer dest.Close()

	//Read lines from both files info sets
	destLines := make(map[string]bool)
	scanner := bufio.NewScanner(dest)
	for scanner.Scan() {
		destLines[scanner.Text()] = true
	}

	// scan source file for new lines and append them to destination
	scanner = bufio.NewScanner(src)
	var newLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !destLines[line] {
			newLines = append(newLines, line)
			destLines[line] = true
		}
	}

	if len(newLines) > 0 {
		for _, line := range newLines {
			_, err := dest.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func setDates(w fyne.Window, start, end string) {
	if start == "" || end == "" {
		// dialog.NewConfirm("Both Date Required")
		message := "Start date and End date are Required"
		dialog.ShowInformation("Warning", message, w)
		return
	}

	fmt.Println(start, " ", end)
}
