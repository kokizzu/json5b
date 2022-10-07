package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kokizzu/json5b/encoding/json5b"
)

func main() {
	procArgs()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: json5b \n")
	fmt.Fprintf(os.Stderr, "Compiles JSON5 file into sibling JSON.\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func procArgs() {
	flag.Usage = usage
	json5Path := flag.String("c", "", "path/to/file.json5b, or blank for stdin")
	outputPath := flag.String("o", "", "path/to/file.json, or blank for stdout")
	flag.Parse()
	var file *os.File
	var err error
	if *json5Path == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(*json5Path)
		if err != nil {
			fmt.Println(err)
			usage()
		}
	}

	var data interface{}
	dec := json5b.NewDecoder(file)
	err = dec.Decode(&data)
	if err != nil {
		fmt.Println(err)
		usage()
	}
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
		usage()
	}
	if *outputPath == "" {
		fmt.Println(string(b))
	} else {
		ioutil.WriteFile(*outputPath, b, 0644)
	}
}
