package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if args[0] == "--help" {
		log.Println("Usage: ./csvToString csvOrigin comma minParams stringWithDollarArgs linesPerFileOutput")
		log.Println("Example: ./csvToString 'origin.csv' ';' '2' \"UPDATE cdr SET ddr = '`$2' WHERE ddr IS NULL AND id = '`$1';\" '100'")
		os.Exit(0)
	}
	// Open origin CSV
	file, err := os.Open(args[0])
	checkError(err)
	defer file.Close()

	// Create Readers
	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)

	// Set CSV Reader Configs
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = []rune(args[1])[0]

	minParams, err := strconv.Atoi(args[2])
	checkError(err)
	linesPerFileOutput := -1
	if len(args) >= 5 {
		linesPerFileOutput, err = strconv.Atoi(args[4])
		checkError(err)
	}
	if linesPerFileOutput == 0 {
		log.Panicln("Lines per page must be at least 1!")
	}
	filename := args[0] + ".csvToString"
	if linesPerFileOutput > 0 {
		filename = filename + ".part"
	}
	stringOriginal := args[3]
	currNumberOfLines := 0
	currPage := 0
	firstFilename := filename
	if linesPerFileOutput > 0 {
		firstFilename = firstFilename + strconv.Itoa(currPage)
	}
	fileDestination := createFileDestination(firstFilename)
	defer fileDestination.Close()
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			log.Println("EOF found.")
			fileDestination.Close()
			os.Exit(0)
		}
		checkError(err)
		currentString := stringOriginal
		if len(line) < minParams {
			log.Println("Skipping line with less than min params: ")
			log.Print(line)
			continue
		}
		for index, value := range line {
			currentString = strings.ReplaceAll(currentString, ("$" + strconv.Itoa(index+1)), strings.TrimSpace(value))
		}
		fileDestination.WriteString(currentString + "\n")
		currNumberOfLines++
		if currNumberOfLines == linesPerFileOutput {
			currPage++
			fileDestination.Close()
			fileDestination = createFileDestination(filename + strconv.Itoa(currPage))
			defer fileDestination.Close()
			currNumberOfLines = 0
		}
	}
}

func checkError(e error) {
	if e != nil {
		log.Panicln(e)
	}
}

func createFileDestination(name string) *os.File {
	log.Println("Creating file", name)
	fileDestination, err := os.Create(name)
	checkError(err)
	return fileDestination
}
