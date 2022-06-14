//
// parsemd
// parse markdown file
// usage: parsemd file.go
//
// author: prr azul software
// date: 28 May 2022
// copyright prr azul software
//
package main

import (
	"os"
	"fmt"
	mdParse "google/gdoc/mdParse"
)

func main() {

	numArg := len(os.Args)

	switch numArg {
	case 0, 1:
		fmt.Printf("no input file provided\n")
		fmt.Printf("usage is: ./mdparse file\n")
		os.Exit(1)

 	case 2:
		fmt.Printf("input file: %s\n", os.Args[1])

	default:
		fmt.Printf("too many command line parameters: %d\n", numArg)
		fmt.Printf("usage is: ./mdparse file\n")
		os.Exit(1)
	}

	mdp := mdParse.InitMdParse()
	err := mdp.ParseMdFile(os.Args[1])
	if err != nil {
		fmt.Printf("error - parseMdfile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("*** success ***")
}
