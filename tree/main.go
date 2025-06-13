package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walk(out, path, printFiles, "", true)
}

func walk(out io.Writer, currentPath string, printFiles bool, prefix string, isLast bool) error {
	files, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	var list []os.DirEntry
	for _, f := range files {
		if f.IsDir() || printFiles {
			list = append(list, f)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})

	for i, f := range list {
		isLast := i == len(list)-1
		branch := "├───"
		if isLast {
			branch = "└───"
		}
		name := f.Name()

		info, _ := f.Info()
		if !f.IsDir() {
			if info.Size() == 0 {
				name += " (empty)"
			} else {
				name += fmt.Sprintf(" (%db)", info.Size())
			}
		}

		fmt.Fprint(out, prefix+branch+name+"\n")

		if f.IsDir() {
			nextPrefix := prefix
			if !isLast {
				nextPrefix += "│\t"
			} else {
				nextPrefix += "\t"
			}
			err := walk(out, filepath.Join(currentPath, f.Name()), printFiles, nextPrefix, isLast)
			if err != nil {
				return err
			}
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
