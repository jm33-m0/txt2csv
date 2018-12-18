package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	// cmdline args
	var (
		head   = flag.String("head", "1,2,3", "Name for each column")
		sep    = flag.String("seps", "'SPC , . | /'", "Separators, used to recognize each line")
		inTXT  = flag.String("in", "", "Input TXT file")
		outCSV = flag.String("out", "", "Resulting csv file name")
	)
	flag.Parse()

	if *head == "" || *sep == "" || *inTXT == "" || *outCSV == "" {
		flag.Usage()
		os.Exit(1)
	}

	// how many columns to scan for
	columnCnt := len(strings.Split(*head, ","))

	// separators
	seps := strings.Split(*sep, " ")

	// read input txt, convert each line and write to csv
	txtFile, err := os.Open(*inTXT)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = txtFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// open csv file in append mode
	csvFile, err := OpenFileStream(*outCSV)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = CloseFileStream(csvFile); err != nil {
			log.Fatal(err)
		}
	}()
	err = AppendToFile(csvFile, *head)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup // batch jobs
	txtScanner := bufio.NewScanner(txtFile)
	for txtScanner.Scan() {
		rawStr := txtScanner.Text()

		go func() {
			wg.Add(1)
			defer wg.Done()
			csvStr := convert(*head, seps, rawStr, columnCnt)
			if csvStr == "" {
				return
			}
			err = AppendToFile(csvFile, csvStr)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	// wait until the wait group is empty
	if &wg != nil {
		wg.Wait()
	}
}

func convert(head string, seps []string, rawStr string, columnCnt int) string {
	var result string
	var words []string

	rawStr = strings.Trim(rawStr, "\n")

	// test if SPC works
	words = strings.Fields(rawStr)
	if len(words) == columnCnt {
		result = strings.Join(words, ",")
		log.Println("SPC works")
		return result
	}

	// split by other seps
	for _, sep := range seps {
		if sep == "SPC" {
			continue
		}
		words = strings.Split(rawStr, sep)
		if len(words) != columnCnt {
			return ""
		}

		result = strings.Join(words, ",")
		log.Printf("%s works", sep)
	}
	return result
}
