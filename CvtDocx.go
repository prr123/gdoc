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
		rc, err := f.Open()
		if err != nil {
			fmt.Printf("error zip.Open file %s: %v\n", f.Name, err)
			os.Exit(-1)
		}
		tmpStr := docx.WordDocToString(rc)

		fmt.Printf("str: %s\n", tmpStr)

		rc.Close()
		fmt.Println("*** end of file ***")
	}

/*
	inpfil, err := os.Open(docxFilNam)
	if err != nil {
		fmt.Printf("error os.Open file %s: %v\n", docxFilNam, err)
		os.Exit(-1)
	}
	defer inpfil.Close()


	outfil, err := os.Create(htmlFilNam)
	if err != nil {
		fmt.Printf("error os.Create file %s: %v\n", htmlFilNam, err)
		os.Exit(-1)
	}
	defer outfil.Close()
	w := bufio.NewWriter(outfil)
//	fmt.Printf("outfil: %v name: %s\n", outfil, htmlFilNam)
*/
	// read file
/*
	inpfilInfo,_ := inpfil.Stat()
	inpSize := inpfilInfo.Size()

	inbufp := make([]byte, inpSize)
    nb, _ := inpfil.Read(inbufp)
    if nb != int(inpSize) {
		fmt.Printf("error could not read all input!")
		os.Exit(-1)
	}

//
	err = goldmark.AltConvert(inbufp, w)

	if err != nil {
		fmt.Printf("error goldmark convert: %v\n", err)
		os.Exit(-1)
	}
*/

	fmt.Println("**** success *****")
}
