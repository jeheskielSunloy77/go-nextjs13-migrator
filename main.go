package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const folderPath = "src"

func main() {
	fileCount := 0
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		err = processFile(path)
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

func processFile(path string) error {
	if filepath.Ext(path) != ".tsx" || filepath.Ext(path) != ".jsx" {
		fmt.Printf("Skipping file '%s' because it is not a react file\n", path)
		return nil
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
				fmt.Printf("%s is already a react client component!\n", path)
				isClientComponnet = true
			}
		}

		if !isClientComponnet {
			if strings.Contains(line, "from 'react'") {
				fmt.Printf("%s is a react client component!\n", path)
				lines = append([]string{"'use client'"}, lines...)
				isClientComponnet = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	fmt.Printf("Successfully processed file '%s'\n", path)
	return nil
}
