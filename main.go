package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const (
	pipe      string = "│\t"
	tJunction string = "├───"
	end       string = "└───"
)

type ExtendedFileInfo struct {
	os.FileInfo
	path   string
	isLast bool
}

// возвращает остортированное содержимое папки и, если нужно, исключает из него файлы
func getDirContents(path string, shouldFilterFiles bool) (result []ExtendedFileInfo) {
	dir, _ := os.Open(path)
	defer dir.Close()

	dirContents, _ := dir.Readdir(0)

	// сортируем
	sort.Slice(dirContents, func(i, j int) bool {
		return dirContents[i].Name() < dirContents[j].Name()
	})

	for idx, fileInfo := range dirContents {
		// если нужно, убираем из результирующего слайса файлы
		if shouldFilterFiles && !fileInfo.IsDir() {
			dirContents = append(dirContents[:idx], dirContents[idx+1:]...)
			continue
		}

		extendedFileInfo := ExtendedFileInfo{
			FileInfo: fileInfo,
			path:     filepath.Join(path, fileInfo.Name()),
		}

		result = append(result, extendedFileInfo)
	}

	if result != nil {
		result[len(result)-1].isLast = true
	}

	return
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	rootDirContents := getDirContents(path, !printFiles)
	var branch = [][]ExtendedFileInfo{rootDirContents}
	var depth = 0

	for len(branch[0]) > 0 {
		if len(branch[depth]) == 0 {
			branch = branch[:len(branch)-1]
			depth--
			branch[depth] = branch[depth][1:]
			continue
		}

		currentFile := branch[depth][0]

		if currentFile.IsDir() || printFiles {
			if !currentFile.isLast {
				fmt.Fprintln(out, tJunction+currentFile.Name())
			} else {
				fmt.Fprintln(out, end+currentFile.Name())
			}
		}

		if currentFile.IsDir() {
			branch = append(branch, getDirContents(currentFile.path, !printFiles))
			depth++
		} else {
			branch[depth] = branch[depth][1:]
		}
	}

	//for _, file := range rootDirContents {
	//	if !file.isLast {
	//		fmt.Fprintln(out, tJunction+file.Name())
	//	} else {
	//		fmt.Fprintln(out, end+file.Name())
	//	}
	//}

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
