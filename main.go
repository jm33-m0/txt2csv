package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

type item struct {
	username string `csv:"username"`
	password string `csv:"password"`
	email    string `csv:"email"`
}

func main() {
	// cmdline args
	var (
		sep    = flag.String("sep", "whitespace", "Separator")
		inTXT  = flag.String("in", "", "Input TXT file")
		outCSV = flag.String("out", "", "Resulting csv file name")
	)
	flag.Parse()

	if *sep == "" || *inTXT == "" || *outCSV == "" {
		flag.Usage()
		os.Exit(1)
	}

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

	var wg sync.WaitGroup // batch jobs
	txtScanner := bufio.NewScanner(txtFile)
	for txtScanner.Scan() {
		rawStr := txtScanner.Text()

		go func() {
			wg.Add(1)
			defer wg.Done()
			csvStr := convert(*sep, rawStr)
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

func convert(sep string, rawStr string) string {
	var result string

	rawStr = strings.Trim(rawStr, "\n")

	if sep == "whitespace" {
		// convert to comma sep
		words := strings.Fields(rawStr)
		if len(words) <= 1 {
			return ""
		}

		result = strings.Join(words, ",")

		return result
	}

	// split by sep and convert
	words := strings.Split(rawStr, sep)
	if len(words) <= 1 {
		return ""
	}

	result = strings.Join(words, ",")
	return result
}
