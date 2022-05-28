//
// credoc
// create doc documentation
// usage: credoc file.go
//
// author: prr azul software
// date: 20 May 2022
// copyright prr azul software
//
package main

import (
	"os"
	"fmt"
	utilLib "google/gdoc/util"
)

func main() {

	numArg := len(os.Args)

	switch numArg {
	case 0, 1:
		fmt.Printf("no input file provided\n")
		fmt.Printf("usage is: credoc file\n")
		os.Exit(1)

 	case 2:
		fmt.Printf("input file: %s\n", os.Args[1])

	default:
		fmt.Printf("in correct number of command line parameters: %d\n", numArg)
		fmt.Printf("usage is: credoc file\n")
		os.Exit(1)
	}

	err := parseMdFile(os.Args[1])
	if err != nil {
		fmt.Printf("error - parseMdfile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("*** success ***")
}


func parseMdFile(inpfilnam string)(err error) {
// function that creates doc output file

	var outfilnam string

	// check whether input file is valid

	// check for period
	iper := -1
	for i:=0; i< len(inpfilnam); i++ {
		if inpfilnam[i] == '.' {
			iper = i
			break
		}
	}

	if iper < 0 {return fmt.Errorf(" error parsing input file name, no period for extension found!")}

	// check for md extension
	if inpfilnam[iper+1] != 'm' {return fmt.Errorf(" error extension not md!")}
	if inpfilnam[iper+2] != 'd' {return fmt.Errorf(" error extension not md!")}

	inpfil, err := os.Open(os.Args[1])
	defer inpfil.Close()
	if err != nil {return fmt.Errorf("os.Open: %v\n", err)}

	inpfilInfo,_ := inpfil.Stat()
	inpSize := inpfilInfo.Size()

	fmt.Printf("*** input file:  %s Size: %d\n", inpfilnam, inpSize)

	// create output file name
//	found := false
	for i:=0; i< len(inpfilnam); i++ {
		if inpfilnam[i] == '.' {
//			found = true
			outfilnam = string(inpfilnam[:i])
			break
		}
	}

	fmt.Printf("*** output file: %s\n", outfilnam + ".md")

	outfil, err := os.Create(outfilnam + ".md")
	if err != nil { return fmt.Errorf("os.Create: %v\n", err)}
	defer outfil.Close()

	return nil
}

func IsTyp(buf []byte)(typNam string, res bool) {
// function that checks whether an input line is a type definition.
// It returns the type name if the input line is a typpe definition.

	if len(buf) < 5 { return "", false }

	if buf[0] != 't' {return "", false}
	if buf[1] != 'y' {return "", false}
	if buf[2] != 'p' {return "", false}
	if buf[3] != 'e' {return "", false}
	if buf[4] != ' ' {return "", false}

	fnamSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if utilLib.IsAlpha(buf[i]) {
			fnamSt = i
			break
		}
	}

	if fnamSt == 0 {return "", false}

	fnamEnd := 0
	for i:= fnamSt; i<len(buf); i++ {
		if (buf[i] == ' ') {
			fnamEnd = i
			break
		}
	}

	if fnamEnd == 0 {return "", false}

	typNam = string(buf[fnamSt:fnamEnd])
	return typNam, true
}

func IsFunc(buf []byte)(funcNam string, res bool) {
// function that checks whether an input line is a function.
// It returns the function name if the input line is a  function.

	if len(buf) < 5 { return "", false }

	if buf[0] != 'f' {return "", false}
	if buf[1] != 'u' {return "", false}
	if buf[2] != 'n' {return "", false}
	if buf[3] != 'c' {return "", false}
	if buf[4] != ' ' {return "", false}

	fnamSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if buf[i] == '(' {return "", false}
		if utilLib.IsAlpha(buf[i]) {
			fnamSt = i
			break
		}
	}

	if fnamSt == 0 {return "", false}

	fnamEnd := 0
	for i:= fnamSt; i<len(buf); i++ {
		if (buf[i] == ' ') || (buf[i] == '(') {
			fnamEnd = i
			break
		}
	}

	if fnamEnd == 0 {return "", false}

	funcNam = string(buf[fnamSt:fnamEnd])
	return funcNam, true
}

func IsMethod(buf []byte)(methNam string, typNam string, res bool) {
// function that detemines whether a input line is a  method.
// if so, the function returns the method name and the name of the structure the method is associated with

	if len(buf) < 5 { return "","", false }

	if buf[0] != 'f' {return "", "", false}
	if buf[1] != 'u' {return "", "", false}
	if buf[2] != 'n' {return "", "", false}
	if buf[3] != 'c' {return "", "", false}
	if buf[4] != ' ' {return "", "", false}

	typSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if utilLib.IsAlpha(buf[i]) {return "", "", false}
		if buf[i] == '(' {
			typSt = i
			break
		}
	}

	if typSt == 0 {return "", "", false}

	typEnd := 0
	for i:= typSt; i<len(buf); i++ {
		if buf[i] == ')' {
			typEnd = i
			break
		}
	}

	if typEnd == 0 {return "", "", false}

	typNam = string(buf[typSt+1:typEnd-1])
	methSt := 0
	for i:= typEnd + 1; i<len(buf); i++ {
		if utilLib.IsAlpha(buf[i]) {
			methSt = i
			break
		}
	}

	if methSt == 0 {return "", "", false}

	methEnd := 0
	for i:= methSt; i<len(buf); i++ {
		if (buf[i] == ' ') || (buf[i] == '(') {
			methEnd = i
			break
		}
	}

	if methEnd == 0 {return "", "", false}

	methNam = string(buf[methSt:methEnd])

	return methNam, typNam, true
}

func IsComment (buf []byte)(desc string, res bool) {
// function that determines whether the input line is a comment line.
// If so, it returns the comment in the desc.

	if len(buf) < 2 {return "", false}
	if buf[0] != '/' {return "", false}
	if buf[1] != '/' {return "",  false}

	desc = string(buf[2:])

	return desc, true
}
