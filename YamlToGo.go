// golang program that parses a yaml file and generates a go struct
// author: prr, azul software
// created: 20/12/2022
// copyright 2022 prr, Peter Riemenschneider, Azul Software
//
// start: YamlToGo
//

package main

import (
	"os"
	"fmt"
	"bytes"
    util "google/gdoc/utilLib"
	)

func main() {

    numArgs:= len(os.Args)

    if numArgs < 2 {
        fmt.Printf("error - exit: insufficient command line arguments\n")
        fmt.Printf("usage is: YamlToGo \"input file\" [\\out=] [\\dbg]\n")
        os.Exit(-1)
    }

    inFilNam := os.Args[1]
	idx := bytes.Index([]byte(inFilNam),[]byte(".tpl"))
	if idx == -1 {fmt.Printf("invalid input file name \"%s\": no ext \".tpl\"!\n", inFilNam); os.Exit(-1);}

	infil, err := os.Open(inFilNam)
	if err != nil {fmt.Printf("cannot open input file \"%s\": %v!\n", inFilNam, err); os.Exit(-1);}
	defer infil.Close()

    flags := [] string {"out", "dbg"}

    argmap, err := util.ParseFlagsStart(os.Args, flags,2)
    if err != nil {fmt.Printf("error ParseFlags: %v\n", err); os.Exit(-1);}

    outFilNam, ok := argmap["out"]
	outFilNamStr := ""
	tplFilNamStr := string(inFilNam[:idx]) + ".go"
    if !ok {
		outFilNamStr += tplFilNamStr
		fmt.Printf("no output filename provided! using \"%s\"\n", outFilNamStr)

	} else {
    	outFilNamStr := outFilNam.(string)
		outIdx := bytes.Index([]byte(outFilNamStr),[]byte(".go"))
		if outIdx == -1 {fmt.Printf("invalid output file name \"%s\": no ext \".go\"!\n", outFilNamStr); os.Exit(-1);}
    	if outFilNamStr == "none" {outFilNamStr += tplFilNamStr}
	}

fmt.Printf("out file: %s\n",outFilNamStr)

	outfil, err := os.Create(outFilNamStr)
	if err != nil {fmt.Printf("error creating output file: %v",err); os.Exit(-1)}
	defer outfil.Close()
	outstr := fmt.Sprintf("// go struct for yaml file: %s\n", inFilNam)
	outfil.WriteString(outstr)

	info, err := infil.Stat()
	if err != nil {fmt.Printf("error reading stat of input file: %v",err); os.Exit(-1)}
	fmt.Printf("infile size: %d\n", info.Size())


	fmt.Println("YamlToGo success")
}
