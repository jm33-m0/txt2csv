package main

import (
	"bufio"
	"log"
	"os"
)

// AppendToFile : append a line to target file
func AppendToFile(file *os.File, line string) (err error) {
	// write appendly
	if _, err = file.Write([]byte(line + "\n")); err != nil {
		log.Print("Write err: ", err, "\nWriting ", line)
		return err
	}
	return nil
}

// OpenFileStream : open file for writing
func OpenFileStream(filepath string) (file *os.File, err error) {
	// open outfile
	file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(filepath, " : Failed to open file\n", err)
		return nil, err
	}
	return file, nil
}

// CloseFileStream : Close file when we are done with it
func CloseFileStream(file *os.File) (err error) {
	err = file.Close()
	return err
}

// GetFileLength : How many lines does a text file contain
func GetFileLength(file string) (int, error) {
	i := 0

	lines, err := FileToLines(file)
	if err != nil {
		log.Printf("Can't open file: %s", err.Error())
	}
	for range lines {
		i++
	}

	return i, err
}

// FileToLines : Read lines from a text file
func FileToLines(filepath string) ([]string, error) {
	f, err := os.Open(filepath)
	if err == nil {
		defer func() {
			if err = f.Close(); err != nil {
				log.Printf(err.Error())
			}
		}()

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if scanner.Err() != nil {
			return nil, scanner.Err()
		}
		return lines, nil
	}
	return nil, err
}
