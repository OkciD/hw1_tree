package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

const (
	pipe      string = "│\t"
	tJunction string = "├───"
	end       string = "└───"
)

// Из слайса файлов оставляет только директории
func filterFiles(files []os.FileInfo) (result []os.FileInfo) {
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		result = append(result, file)
	}

	return
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var rootDir, _ = os.Open(path)
	defer rootDir.Close()

	rootDirContents, err := rootDir.Readdir(0)

	if err != nil {
		return err
	}

	sort.Slice(rootDirContents, func(i, j int) bool {
		return rootDirContents[i].Name() < rootDirContents[j].Name()
	})

	if !printFiles {
		rootDirContents = filterFiles(rootDirContents)
	}

	for idx, file := range rootDirContents {
		var isLastFile = idx == len(rootDirContents)-1

		if !isLastFile {
			fmt.Fprintln(out, tJunction+file.Name())
		} else {
			fmt.Fprintln(out, end+file.Name())
		}
	}

	return nil
}

func main() {
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
