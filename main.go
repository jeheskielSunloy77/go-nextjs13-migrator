package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const folderPath = "../../../nextjs-iakn-kupang/src"

func main() {
	if err := os.Mkdir(folderPath+"/app", 0777); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	fileCount := 0
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tsx" && filepath.Ext(path) != ".jsx" {
			return nil
		}
		if strings.Contains(filepath.Dir(path), folderPath+"/app") {
			return nil
		}
		if info.Name() == "_app.tsx" || info.Name() == "_document.tsx" {
			return nil
		}

		err = processFile(path, info)
		if err != nil {
			log.Printf("Error processing file '%s': %s", path, err)
		}
		fileCount++
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File scanning and modification complete! \n Total files scanned: ", fileCount)
}

func processFile(path string, info os.FileInfo) error {
	var newFile *os.File
	isPages := false

	if strings.Contains(filepath.Dir(path), folderPath+"/pages") {
		isPages = true
		trimedDir := strings.TrimPrefix(filepath.Dir(path), folderPath+"/pages")
		var folderName string
		if trimedDir != "" {
			folderName = folderPath + "/app" + trimedDir
		} else {
			folderName = folderPath + "/app/" + strings.Split(info.Name(), ".")[0]
		}
		if err := os.MkdirAll(folderName, 0777); err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
		file, err := os.Create(folderName + "/page.tsx")
		if err != nil {
			return err
		}
		defer file.Close()
		newFile = file
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	lineCount := 0
	isClientComponnet := false

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		lineCount++
		if lineCount == 1 {
			if line == `'use client'` {
				// fmt.Printf("%s is already a react client component!\n", path)
				isClientComponnet = true
			}
		}

		if !isClientComponnet {
			if strings.Contains(line, "from 'react'") {
				// fmt.Printf("%s is a react client component!\n", path)
				lines = append([]string{"'use client'"}, lines...)
				isClientComponnet = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	var writer *bufio.Writer
	if isPages {
		writer = bufio.NewWriter(newFile)
	} else {
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		writer = bufio.NewWriter(file)
	}
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	// fmt.Printf("Successfully processed file '%s'\n", path)
	return nil
}
