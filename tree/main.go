package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	type fileNode struct {
		name  string
		isDir bool
		size  int64
	}

	var buildTree func(string, []fileNode, string) error
	buildTree = func(path string, files []fileNode, prefix string) error {
		for i, f := range files {
			isLast := i == len(files)-1
			branch := "├───"
			if isLast {
				branch = "└───"
			}

			fmt.Fprintf(out, "%s%s%s", prefix, branch, f.name)

			if !f.isDir {
				if f.size == 0 {
					fmt.Fprint(out, " (empty)")
				} else {
					fmt.Fprintf(out, " (%db)", f.size)
				}
			}
			fmt.Fprintln(out)

			if f.isDir {
				nextPrefix := prefix
				if !isLast {
					nextPrefix += "│\t"
				} else {
					nextPrefix += "\t"
				}

				dirPath := filepath.Join(path, f.name)
				dirEntries, err := os.ReadDir(dirPath)
				if err != nil {
					return err
				}

				var children []fileNode
				for _, e := range dirEntries {
					info, _ := e.Info()
					if info.IsDir() || printFiles {
						children = append(children, fileNode{
							name:  e.Name(),
							isDir: e.IsDir(),
							size:  info.Size(),
						})
					}
				}

				sort.Slice(children, func(i, j int) bool {
					return children[i].name < children[j].name
				})

				buildTree(dirPath, children, nextPrefix)
			}
		}
		return nil
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var rootFiles []fileNode
	for _, e := range dirEntries {
		info, _ := e.Info()
		if info.IsDir() || printFiles {
			rootFiles = append(rootFiles, fileNode{
				name:  e.Name(),
				isDir: e.IsDir(),
				size:  info.Size(),
			})
		}
	}

	sort.Slice(rootFiles, func(i, j int) bool {
		return rootFiles[i].name < rootFiles[j].name
	})

	return buildTree(path, rootFiles, "")
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
