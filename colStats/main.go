package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)


func main() {
	// verify and parse arguments
	op := flag.String("op", "sum", "Operation to perform")
	column := flag.Int("col", 1, "Column to perform operation on")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	// Validate input parameters

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)	
	}

	// Validate the operation
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}


	consolidate := make([]float64, 0)

	// Loop through all files
	for _, fname := range filenames {
		// Open the file for reading
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("Cannot open file: %w", err)
		}

		// Parse the csv into a slice of float64
		data, err := csv2float(f, column)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err	
		}

		// append the data to consolidate
		consolidate = append(consolidate, data...)
	}

	_, err := fmt.Fprintln(out, opFunc(consolidate))
	return err
}