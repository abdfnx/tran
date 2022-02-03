package dfs

import (
	"os"
	"fmt"
	"io/fs"
	"errors"
	"strings"
	"path/filepath"
)

// Directory shortcuts.
const (
	CurrentDirectory  = "."
	PreviousDirectory = ".."
	HomeDirectory     = "~"
	RootDirectory     = "/"
)

// Different types of listings.
const (
	DirectoriesListingType = "directories"
	FilesListingType       = "files"
)

// RenameDirectoryItem renames a directory or files given a source and destination.
func RenameDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return err
}

// CreateDirectory creates a new directory given a name.
func CreateDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDirectoryListing returns a list of files and directories within a given directory.
func GetDirectoryListing(dir string, showHidden bool) ([]fs.DirEntry, error) {
	n := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if !showHidden {
		for _, file := range files {
			// If the file or directory starts with a dot,
			// we know its hidden so dont add it to the array
			// of files to return.
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		}

		// Set files to the list that does not include hidden files.
		files = files[:n]
	}

	return files, nil
}

// GetDirectoryListingByType returns a directory listing based on type (directories | files).
func GetDirectoryListingByType(dir, listingType string, showHidden bool) ([]fs.DirEntry, error) {
	n := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		switch {
			case file.IsDir() && listingType == DirectoriesListingType && !showHidden:
				if !strings.HasPrefix(file.Name(), ".") {
					files[n] = file
					n++
				}

			case file.IsDir() && listingType == DirectoriesListingType && showHidden:
				files[n] = file
				n++

			case !file.IsDir() && listingType == FilesListingType && !showHidden:
				if !strings.HasPrefix(file.Name(), ".") {
					files[n] = file
					n++
				}

			case !file.IsDir() && listingType == FilesListingType && showHidden:
				files[n] = file
				n++
		}
	}

	return files[:n], nil
}

// DeleteDirectory deletes a directory given a name.
func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)

	return err
}

// GetHomeDirectory returns the users home directory.
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// GetWorkingDirectory returns the current working directory.
func GetWorkingDirectory() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return workingDir, nil
}

// DeleteFile deletes a file given a name.
func DeleteFile(name string) error {
	err := os.Remove(name)

	return err
}

// MoveDirectoryItem moves a file from one place to another.
func MoveDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return err
}

// ReadFileContent returns the contents of a file given a name.
func ReadFileContent(name string) (string, error) {
	fileContent, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

// GetDirectoryItemSize calculates the size of a directory or file.
func GetDirectoryItemSize(path string) (int64, error) {
	curFile, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !curFile.IsDir() {
		return curFile.Size(), nil
	}

	var size int64

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if !d.IsDir() {
			size += fileInfo.Size()
		}

		return err
	})

	return size, err
}

func FindFilesByName(name, dir string) ([]string, []fs.DirEntry, error) {
	var paths []string
	var entries []fs.DirEntry

	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if strings.Contains(entry.Name(), name) {
			paths = append(paths, path)
			entries = append(entries, entry)
		}

		return err
	})

	return paths, entries, err
}

// WriteToFile writes content to a file, overwriting content if it exists.
func WriteToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s\n", filepath.Join(workingDir, content)))

	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return err
}
