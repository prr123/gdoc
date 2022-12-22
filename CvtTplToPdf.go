// golang program that parses a yaml tpl file, creates a copy and replaces fields with data
// author: prr, azul software
// created: 22/12/2022
// copyright 2022 prr, Peter Riemenschneider, Azul Software
//
// start: CvtTplToPdf
//

package main

import (
    "os"
    "fmt"
    "bytes"
	"gopkg.in/yaml.v3"
    util "google/gdoc/utilLib"
    )

// needs to be adjusted for each template

// go struct for yaml file: output/tpltest/simpleTpl.tpl
type tplObj struct {
  Title string `yaml:"Title"`
  Id string `yaml:"Id"`
  NamesLen string `yaml:"NamesLen"`
  First string `yaml:"first"`
  Last string `yaml:"last"`
  Num int `yaml:"num"`
}


func main() {

    numArgs:= len(os.Args)

    if numArgs < 2 {
        fmt.Printf("error - exit: insufficient command line arguments\n")
        fmt.Printf("usage is: CvtTplToPdf \"yaml file\" [\\out=] [\\dbg]\n")
        os.Exit(-1)
    }

    inFilNam := os.Args[1]
    idx := bytes.Index([]byte(inFilNam),[]byte(".tpl"))
    if idx == -1 {fmt.Printf("invalid input file name \"%s\": no ext \".tpl\"!\n", inFilNam); os.Exit(-1);}
	BackOutFilNam := string(inFilNam[:idx]) + ".pdf"

    infil, err := os.Open(inFilNam)
    if err != nil {fmt.Printf("cannot open input file \"%s\": %v!\n", inFilNam, err); os.Exit(-1);}
    defer infil.Close()

    flags := [] string {"out", "dbg"}

    argmap, err := util.ParseFlagsStart(os.Args, flags,2)
    if err != nil {fmt.Printf("error ParseFlags: %v\n", err); os.Exit(-1);}

    outFilNamStr := ""


    outFilNam, ok := argmap["out"]
    if !ok {
        fmt.Printf("no output filename provided! using \"%s\"\n", BackOutFilNam)
		outFilNamStr = BackOutFilNam
//		os.Exit(-1)
    } else {
        outFilNamStr := outFilNam.(string)
        outIdx := bytes.Index([]byte(outFilNamStr),[]byte(".pdf"))
        if outIdx == -1 {fmt.Printf("invalid output file name \"%s\": no ext \".pdf\"!\n", outFilNamStr); os.Exit(-1);}
        if outFilNam == "none" {outFilNamStr = BackOutFilNam}
    }

fmt.Printf("out file: %s\n",outFilNamStr)

	// read Yaml File
    info, err := infil.Stat()
    if err != nil {fmt.Printf("error reading stat of input file: %v",err); os.Exit(-1)}

//  fmt.Printf("infile size: %d\n", info.Size())

    inbuf := make([]byte, info.Size())

    _, err = infil.Read(inbuf)
    if err != nil {fmt.Printf("error reading input file: %v",err); os.Exit(-1)}

  fmt.Println("inbuf:\n", string(inbuf[:100]))

	var tpl tplObj
	err = yaml.Unmarshal(inbuf, &tpl)
    if err != nil {
        fmt.Printf("unmarshall error: %v\n", err)
    }
  fmt.Printf("tpl: %v\n", tpl)

	printTpl(tpl)

	fmt.Println("*** success CvtTplToPdf ****")
}

func printTpl(tpl tplObj) {

	fmt.Println("******** tpl ********")
	fmt.Printf("Title: %s\n", tpl.Title)
	fmt.Printf("Id:    %s\n", tpl.Id)
	fmt.Printf("Fields: %s\n", tpl.NamesLen)
	fmt.Printf("First:  %s\n", tpl.First)
	fmt.Printf("Last:  %s\n", tpl.Last)
	fmt.Printf("Num:  %d\n", tpl.Num)

}    
