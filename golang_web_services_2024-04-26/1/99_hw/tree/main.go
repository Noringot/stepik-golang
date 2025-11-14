package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

const FILE_PREFIX = "├───"
const FILE_LAST_PREFIX = "└───"
const INDENT_SYMBOL = "\t"
const FOLDER_CONNECTION = "│"

func fileFilterDirs(dirs []os.DirEntry, isIncludeFiles bool) []os.DirEntry {
	if isIncludeFiles {
		return dirs
	}

	filtered := make([]os.DirEntry, 0, len(dirs))

	for _, dir := range dirs {
		if dir.IsDir() {
			filtered = append(filtered, dir)
		}
	}

	return filtered
}

func buildTree(path string, level int, parentPrefix string, isPrintFiles bool) (string, error) {

	var indentString string

	dirs, err := os.ReadDir(path)

	if err != nil {
		return "", err
	}

	dirs = fileFilterDirs(dirs, isPrintFiles)

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	for i, dir := range dirs {
		if !isPrintFiles && !dir.IsDir() {
			continue
		}

		currentPrefix := parentPrefix

		if i == len(dirs)-1 {
			currentPrefix += FILE_LAST_PREFIX
		} else {
			currentPrefix += FILE_PREFIX
		}

		indentString += currentPrefix + dir.Name()

		if dir.IsDir() {
			currentFolderConnection := ""
			subDirPath := path + string(os.PathSeparator) + dir.Name()

			if i != len(dirs)-1 {
				currentFolderConnection = FOLDER_CONNECTION
			}

			subDirPrefix := parentPrefix + currentFolderConnection + INDENT_SYMBOL
			subDirTree, err := buildTree(subDirPath, level+1, subDirPrefix, isPrintFiles)

			if err != nil {
				return "", err
			}

			indentString += "\n" + subDirTree
		} else {
			info, err := dir.Info()

			if err != nil {
				return "", err
			}

			fileSize := info.Size()
			sizeString := ""

			if fileSize > 0 {
				sizeString = fmt.Sprintf(" (%db)", info.Size())
			} else {
				sizeString = " (empty)"
			}

			indentString += sizeString + "\n"
		}
	}

	return indentString, nil
}

func dirTree(out io.Writer, path string, isPrintFiles bool) error {
	tree, err := buildTree(path, 0, "", isPrintFiles)

	if err != nil {
		return err
	}

	fmt.Fprint(out, tree)

	return nil
}

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	isPrintFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, isPrintFiles)

	if err != nil {
		panic(err.Error())
	}
}
