package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Message map[string]interface{}

func ReadInput(inputFilePath string) ([]Message, error) {
	var contents []byte
	var j []Message
	var scanner *bufio.Scanner
	if inputFilePath != "" {
		file, err := os.Open(inputFilePath)
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(file)

	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		contents = append(contents, scanner.Bytes()...)
	}

	json.Unmarshal(contents, &j)

	return j, nil
}

func getColumns(r Message) []string {
	var collector []string
	for key := range r {
		collector = append(collector, key)
	}
	return collector
}

func BuildWriter(outpath string) (*bufio.Writer, error) {
	var writer *bufio.Writer
	if outpath != "" {
		file, err := os.Create(outpath)
		if err != nil {
			return writer, err
		}
		writer = bufio.NewWriter(file)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	return writer, nil
}

// The main function
func main() {
	var inputFilePath string
	var outputFilePath string
	flag.StringVar(&inputFilePath, "f", "", "Input File path")
	flag.StringVar(&outputFilePath, "o", "", "Output file path")

	flag.Parse()

	contents, err := ReadInput(inputFilePath)

	if err != nil {
		panic(err)
	}

	var columns []string

	if len(contents) > 0 {
		columns = getColumns(contents[0])
	}

	writer, err := BuildWriter(outputFilePath)

	if err != nil {
		error := fmt.Errorf("error building I/O writer: %e", err)
		panic(error)
	}

	defer writer.Flush()

	for idx, col := range columns {
		writer.WriteString(col)
		if idx < len(columns)-1 {
			writer.WriteString(",")
		}

	}
	// Write new line after column headers
	writer.WriteString("\n")

	for _, row := range contents {
		for idx, col := range columns {
			match, ok := row[col]
			if ok {
				writer.WriteString(fmt.Sprintf("%v", match))
			}
			if idx < len(columns)-1 {
				// Only write column 
				writer.WriteString(",")
			}

		}
		writer.WriteString("\n")
	}

}
