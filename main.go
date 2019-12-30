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

// возвращает остортированное содержимое папки и, если нужно, исключаем из него файлы
func getSortedDirContents(dir *os.File, shouldFilterFiles bool) (result []os.FileInfo, err error) {
	dirContents, err := dir.Readdir(0)

	if err != nil {
		return nil, err
	}

	// если нужно, убираем из слайса файлы
	if shouldFilterFiles {
		for _, file := range dirContents {
			if !file.IsDir() {
				continue
			}

			result = append(result, file)
		}
	} else {
		result = dirContents
	}

	// сортируем
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})

	return
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var rootDir, _ = os.Open(path)
	defer rootDir.Close()

	rootDirContents, err := getSortedDirContents(rootDir, !printFiles)

	if err != nil {
		return err
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
