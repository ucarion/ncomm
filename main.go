package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ucarion/cli"
)

func main() {
	cli.Run(context.Background(), ncomm)
}

type args struct {
	ShowHeader bool     `cli:"--show-header" usage:"output a header to help explain what each column represents"`
	Files      []string `cli:"files..."`
}

func (args args) ExtendedDescription() string {
	return "like comm(1), but for any number of files"
}

func ncomm(_ context.Context, args args) error {
	if len(args.Files) == 0 {
		return errors.New("at least one file is required")
	}

	if args.ShowHeader {
		for i := 1; i < 1<<len(args.Files); i++ {
			for j := 0; j < len(args.Files); j++ {
				if i&(1<<j) != 0 {
					fmt.Print("x")
				} else {
					fmt.Print("-")
				}
			}

			fmt.Print("\t")
		}

		fmt.Println()
	}

	var files []*bufio.Scanner
	for _, f := range args.Files {
		file, err := os.Open(f)
		if err != nil {
			return err
		}

		defer file.Close()
		files = append(files, bufio.NewScanner(file))
	}

	// init each scanner
	for _, f := range files {
		if !f.Scan() {
			if err := f.Err(); err != nil {
				return err
			}
		}
	}

	for {
		var min []byte
		var hasMin []int
		for i, f := range files {
			if f.Bytes() == nil {
				continue
			}

			if min == nil {
				min = f.Bytes()
			}

			c := bytes.Compare(f.Bytes(), min)
			if c < 0 {
				min = f.Bytes()
				hasMin = []int{i}
			} else if c == 0 {
				hasMin = append(hasMin, i)
			}
		}

		if min == nil {
			return nil
		}

		var tabs int
		for _, i := range hasMin {
			tabs = tabs + 1<<i
			if !files[i].Scan() {
				if err := files[i].Err(); err != nil {
					return err
				}
			}
		}

		fmt.Printf("%s%s\n", strings.Repeat("\t", tabs-1), min)
	}
}
