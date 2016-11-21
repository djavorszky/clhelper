package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	sourcePath      = "patching-tool/patches"
	destinationPath = "patching-tool/spit"
)

func main() {
	all := flag.Bool("a", false, "Specify whether to move all files.")

	list := flag.Bool("l", false, "List the contents of the patches folder. When used in conjunction with '-r', lists the contents of the temp folder.")
	listR := flag.Bool("lr", false, "List the contents of the temp directory. Same as specifying both '-r' and '-l'")

	reverse := flag.Bool("r", false, "Specify to switch direction.")

	flag.Usage = func() {
		explanation :=
			`Spit is a CLI app written in Go that copies fixes from the patches folder (patching-tool/patches)
to a temp folder (patching-tool/spit) and vice-versa.
You can specify as many filenames as arguments as you'd like.

Processed arguments:
	1234 -> *hotfix-1234*
	de-6 -> *de-6*
	portal-45 -> *portal-45*`

		fmt.Fprintf(os.Stderr, "%s\n\nUsage of %s:\n", explanation, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *listR {
		listDir(destinationPath)
		return
	}

	sliceStart := 1

	if *reverse {
		sourcePath, destinationPath = destinationPath, sourcePath
		sliceStart++
	}

	if *list {
		listDir(sourcePath)
		return
	}

	createTmp()

	if *all {
		files, err := ioutil.ReadDir(sourcePath)
		if err != nil {
			log.Fatalf("Failed listing directory %s: %s\n", sourcePath, err)
		}

		for _, file := range files {
			moveFile(file.Name())
		}
		return
	}

	args := os.Args[sliceStart:]

	if len(args) == 0 {
		log.Fatal("Nothing to do, exiting.")
	}

	for _, arg := range args {
		moveFile(arg)
	}
}

func moveFile(file string) {
	filename, err := getFileName(file)

	if err != nil {
		fmt.Printf("Ignoring '%s': %s", file, err)
		return
	}

	srcPath := fmt.Sprintf("%s/%s", sourcePath, filename)
	dstPath := fmt.Sprintf("%s/%s", destinationPath, filename)
	if err := os.Rename(srcPath, dstPath); err != nil {
		log.Fatal("Could not move file: ", err)
	}

	fmt.Printf("Moving: %s -> %s\n", srcPath, dstPath)
}

func getFileName(file string) (string, error) {
	pattern := "\\w+-[0-9]+"

	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Couldn't match %s to %s: %s\n", pattern, file, err)
	}

	exists := r.FindString(file)

	if exists == "" {
		pattern := "^[0-9]+$"

		r, err := regexp.Compile(pattern)
		if err != nil {
			log.Fatalf("Couldn't match %s to %s: %s\n", pattern, file, err)
		}

		exists = fmt.Sprintf("hotfix-%s", r.FindString(file))
	}

	if exists == "" {
		return "", fmt.Errorf("Couldn't match '%s' to any known patterns", file)
	}

	p := fmt.Sprintf("%s/*%s*", sourcePath, exists)

	matches, err := filepath.Glob(p)
	if err != nil {
		log.Fatalln("Pattern is malformed: ", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("Didn't find any matching file for '%s'\n", file)
	}

	filename := strings.Split(matches[0], "/")

	return filename[len(filename)-1], nil
}

func listDir(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed listing directory %s: %s\n", path, err)
	}

	if len(files) == 0 {
		fmt.Printf("No files found in %s\n", path)
		return
	}

	fmt.Printf("Listing files in %s:\n", path)
	for _, file := range files {
		fmt.Printf(">> %s\n", file.Name())
	}
}

func exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

func createTmp() {
	if exists(destinationPath) {
		return
	}

	if err := os.Mkdir(destinationPath, os.ModePerm); err != nil {
		log.Fatal("Could not create directory: ", err)
	}
}
