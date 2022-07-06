// CvtDocx.go
// author: prr, azul software
// date 6 July 2022
// copyright 2022 prr, azul software
//
// using word


package main

import (
	"fmt"
	"os"
//	"bufio"
	"archive/zip"
	"google/gdoc/wordlib"
)


func main() {

var file *zip.File
//	var buf bytes.Buffer

	numArg := len(os.Args)

//	htmlflag:= false
	inpfilNam := ""
	outfilNam := ""

	switch numArg {
	case 1:
		fmt.Println("error -- no output file name provided!")
		fmt.Println("usage is ./CvtDocx infil outfil!")
		os.Exit(1)
	case 2:
		inpfilNam = os.Args[1]

	case 3:
		inpfilNam = os.Args[1]
		outfilNam = os.Args[2]

	default:
		fmt.Println("error -- too many arguments!")
		fmt.Println("usage is ./CvtMdToHtml outfil!")
		os.Exit(1)

	}

	if outfilNam == "" {outfilNam = inpfilNam}

//	htmlFilNam := "output/htmlTest/" + outfilNam + ".html"
	docxFilNam := "inpTestDocx/" + inpfilNam + ".docx"


	r, err := zip.OpenReader(docxFilNam)
	if err != nil {
		fmt.Printf("error zip.OpenReader file %s: %v\n", docxFilNam, err)
		os.Exit(-1)
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)


        if f.Name == "word/document.xml" {
            file = f
			break
		}

//		fmt.Println("*** end of file ***")
	}

	if file == nil {
		fmt.Printf("error could not find file word/document!\n")
		os.Exit(-1)
	}

	docxFil, err := file.Open()

	if err != nil {
		fmt.Printf("error zip.Open file %s: %v\n", file.Name, err)
		os.Exit(-1)
	}
	defer docxFil.Close()

	tmpStr := docx.WordDocToString(docxFil)

	fmt.Printf("content:\n%s\n", tmpStr)

	var doc docx.Document

	err = doc.Extract(tmpStr)
	if err != nil {
		fmt.Printf("error extract XML: %v", err)
		os.Exit(-1)
	}

	fmt.Printf("doc: %v\n", doc.XMLName)
	for i:=0; i<len(doc.Body.Paragraph); i++ {
		fmt.Printf("doc par %d: %s\n", i, doc.Body.Paragraph[i])
	}

	fmt.Println("**** success *****")
}
