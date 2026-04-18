package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	base := flag.String("base", "main", "active base (folder)")
	filename := flag.String("file", "index.txt", "select file to note in")
	continued := flag.Bool("c", false, "when continued is true, it keeps taking input until input is 'q'")
	flag.Parse()
	args := flag.Args()

	notePath, err := os.UserHomeDir()
	if err != nil {
		exit(err)
	}
	notePath = filepath.Join(notePath, "Documents", "note")

	var path string = filepath.Join(notePath, *base, *filename)
	pathExists, err := checkPath(path)
	if !pathExists {
		if err != nil {
			exit(err)
		} else {
			err := createPath(notePath, *base, *filename)
			if err != nil {
				exit(err)
			}
 		}
	}
	
	if *continued {
		err = continuedText(path)
		if err != nil { exit(err) }
	} else {
		if len(args) == 0 { return }
		text := strings.Join(args, " ") + "\n"
		err = writeText(path, text)
		if err != nil { exit(err) }
	}

	fmt.Printf("Noted at %s in base %s\n", *filename, *base)
}

func continuedText(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return fmt.Errorf("Error opening file %s\n", path) }
	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		if text == "q" { break }

		if _, err := file.WriteString(text + "\n"); err != nil {
			return fmt.Errorf("Error writing to file %s\n", path)
		}
	}
	return nil
}

func writeText(path, text string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return fmt.Errorf("Error opening file %s\n", path) }
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return fmt.Errorf("Error writing to file %s\n", path)
	}
	return nil
}

func checkPath(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("Error checking path %s\n", path)
	}
}

func createPath(notePath, base, filename string) error {
	path := filepath.Join(notePath, base)
	
	err := os.MkdirAll(path, 0755) 
	if err != nil {
		return err
	}
	
	path = filepath.Join(path, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("Created", path)
	return nil
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}