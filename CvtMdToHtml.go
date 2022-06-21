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

	htmlflag:= false
	outfilNam :=""
	switch numArg {
	case 1:
		fmt.Println("error -- no output file name provided!")
		fmt.Println("usage is ./CvtMdToHtml outfil!")
		os.Exit(1)
	case 2:
		outfilNam = os.Args[1]

	case 3:
		if os.Args[2] == "html" {htmlflag = true}

	default:
		fmt.Println("error -- too many arguments!")
		fmt.Println("usage is ./CvtMdToHtml outfil!")
		os.Exit(1)

	}

	htmlFilNam := "output/htmlTest/" + outfilNam + ".html"
	mdFilNam := "inpTestMd/" + outfilNam + ".md"

	outfil, err := os.Create(htmlFilNam)
	if err != nil {
		fmt.Printf("error os.Create file %s: %v\n", htmlFilNam, err)
		os.Exit(1)
	}

	mdp := mdParse.InitMdParse()

	// html
	if !htmlflag {
    	err := mdp.ParseMdFile(mdFilNam)
    	if err != nil {
        	fmt.Printf("error - parseMdfile: %v\n", err)
        	os.Exit(1)
    	}
	}
	err = mdp.CvtMdToHtml(outfil)
	if err != nil {
		fmt.Printf("error CvtMdToHml: %v\n", err)
		os.Exit(1)
	}

	outfil.Close()

	fmt.Println("**** success *****")
}
