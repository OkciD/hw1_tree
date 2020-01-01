package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

	for _, fileInfo := range dirContents {
		// если нужно, не даём файлам попасть в результирующий слайс
		if shouldFilterFiles && !fileInfo.IsDir() {
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

		// печатаем директорию или файл
		if currentFile.IsDir() || printFiles {
			var prefix string

			for i := 0; i < depth; i++ {
				if branch[i][0].isLast {
					prefix += "\t"
				} else {
					prefix += pipe
				}
			}

			if !currentFile.isLast {
				prefix += tJunction
			} else {
				prefix += end
			}

			var sizeString string
			if !currentFile.IsDir() && printFiles {
				size := currentFile.Size()

				if size == 0 {
					sizeString = " (empty)"
				} else {
					sizeString = " (" + strconv.FormatInt(size, 10) + "b)"
				}
			}

			fmt.Fprintln(out, prefix+currentFile.Name()+sizeString)
		}

		if currentFile.IsDir() {
			branch = append(branch, getDirContents(currentFile.path, !printFiles))
			depth++
		} else {
			branch[depth] = branch[depth][1:]
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
