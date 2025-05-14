package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const MAX_DEEP_LEVEL = 5

func buildTree(path string, level int, parentPrefix string) (string, error) {
	if level == MAX_DEEP_LEVEL {
		return "", nil
	}

	var indentString string

	d, err := os.ReadDir(path)

	if err != nil {
		return "", err
	}
	subDirsCount := len(d)

	for i := 0; i < subDirsCount; i++ {

		var prefix string
		var folderConnection string

		f := d[i]

		if subDirsCount-i-1 == 0 {
			prefix = "└───"
		} else {
			prefix = "├───"
			folderConnection = "│"
		}

		name := f.Name()

		if f.IsDir() {
			indentString += parentPrefix + prefix + name + "\n"
			subDirPath := path + "/" + name
			subPrefix := "\t" + folderConnection + parentPrefix

			deep, err := buildTree(subDirPath, level+1, subPrefix)

			if err != nil {
				return "", nil
			}

			indentString += deep
		}
	}

	return indentString, nil
}

func dirTree(out io.Writer, path string, withFiles bool) error {
	//fmt.Fprint(out, path, withFiles)

	tree, err := buildTree(path, 0, "")

	if err != nil {
		return err
	}

	fmt.Fprint(out, tree)

	return nil
}

func _main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	r := bytes.NewReader([]byte("Hello world!"))
	r.Read([]byte("Hi"))
	res, _ := r.ReadByte()

	fmt.Println(string(res))
}
