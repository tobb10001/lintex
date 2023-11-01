package reader

import (
	"fmt"
	"slices"
)

// Structure to hold a set of paths to files, which are kept unique.
type FileList struct {
	files []string
	index int
}

// Create a new file list, pre-populated with the given files.
func NewFileList(files... string) FileList {
	list := FileList{files: []string{}, index: -1}
	list.Add(files...)
	return list
}

func (fl *FileList) Add(files... string) {
	for _, file := range files {
		if !slices.Contains(fl.files, file) {
			fl.files = append(fl.files, file)
		}
	}
}

func (fl *FileList) Next() bool {
	if fl.index+1 < len(fl.files) {
		fl.index++
		return true
	}
	return false
}

func (fl *FileList) Get() string {
	if fl.index == -1 {
		panic("FileList.Get() was called, but index was still -1. Call FileList.Next() first!")
	}
	return fl.files[fl.index]
}

type FileDoesNotExistError struct {
	search string
	tried  []string
}

func (fdnee *FileDoesNotExistError) Error() string {
	return fmt.Sprint("Couldn't find file ", fdnee.search, ". Tried ", fdnee.tried, ".")
}
