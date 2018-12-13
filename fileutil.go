package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
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
		PrintError("Can't open file: %s", err.Error())
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
				LogError(err.Error())
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

// PrintCyan : print main msg
func PrintCyan(format string, a ...interface{}) {
	color.Set(color.FgCyan)
	defer color.Unset()
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

// PrintRed : print main msg
func PrintRed(format string, a ...interface{}) {
	color.Set(color.FgRed)
	defer color.Unset()
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

// PrintError : print text in red
func PrintError(format string, a ...interface{}) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

// PrintSuccess : print text in red
func PrintSuccess(format string, a ...interface{}) {
	color.Set(color.FgHiGreen, color.Bold)
	defer color.Unset()
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

// LogError : print log in red
func LogError(format string, a ...interface{}) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	log.Printf(format, a...)
	fmt.Print("\n")
}

// LogSuccess : print log in red
func LogSuccess(format string, a ...interface{}) {
	color.Set(color.FgHiGreen, color.Bold)
	defer color.Unset()
	log.Printf(format, a...)
	fmt.Print("\n")
}

// SetCyan : make text following go cyan
func SetCyan() {
	color.Set(color.FgCyan, color.Bold)
}

// UnsetCyan : make text following go back to normal color
func UnsetCyan() {
	color.Unset()
}
