// CvtMdToHtml.go
// author: prr, azul software
// date 20 June 2022
// copyright 2022 prr, azul software
//

package main

import (
	"fmt"
	"os"
   	mdParse "google/gdoc/mdParse"
)


func main() {

	numArg := len(os.Args)

//	htmlflag:= false
	inpfilNam := ""
	outfilNam := ""

	switch numArg {
	case 1:
		fmt.Println("error -- no output file name provided!")
		fmt.Println("usage is ./CvtMdToHtml infil outfil!")
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

	htmlFilNam := "output/htmlTest/" + outfilNam + ".html"
	mdFilNam := "inpTestMd/" + inpfilNam + ".md"

	outfil, err := os.Create(htmlFilNam)
	if err != nil {
		fmt.Printf("error os.Create file %s: %v\n", htmlFilNam, err)
		os.Exit(1)
	}

//	fmt.Printf("outfil: %v name: %s\n", outfil, htmlFilNam)

	mdp := mdParse.InitMdParse()

	// html
   	err = mdp.ParseMdFile(mdFilNam)
   	if err != nil {
		fmt.Printf("error - parseMdfile %s: %v\n", mdFilNam, err)
        os.Exit(1)
    }

	err = mdp.CvtMdToHtml(outfil)
	if err != nil {
		fmt.Printf("error CvtMdToHml: %v\n", err)
		os.Exit(1)
	}

	outfil.Close()

	fmt.Println("**** success *****")
}
