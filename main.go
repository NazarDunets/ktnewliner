package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	pathToTraverse, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	switch len(args) {
	case 0:
		fmt.Println("Using current directory")
	case 1:
		pathToTraverse = path.Join(pathToTraverse, args[0])
	default:
		fmt.Println("Too many arguments\n0 args - program will check all .kt files in all src/ subdirectories of current working dir\n1 arg - only relative path supported. Same logic applies to provided folder")
		os.Exit(1)
	}

	newlinesAppended := 0

	filepath.Walk(pathToTraverse, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(info.Name(), ".kt") || !strings.Contains(path, "src/") {
			return nil
		}

		file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		bytesToRead := 1
		buf := make([]byte, bytesToRead)
		start := info.Size() - int64(bytesToRead)
		_, err = file.ReadAt(buf, start)

		if err != nil {
			return err
		}

		if string(buf[:]) == "\n" {
			return nil
		} else {
			newlinesAppended++
			fmt.Println("File", info.Name(), "doesn't end with newline. Appending...")
			_, err = file.WriteString("\n")

			if err != nil {
				return err
			}
		}

		return nil
	})

	fmt.Println("Total of", newlinesAppended, "newlines appended")
}
