// CvtMdToHtml.go
// author: prr, azul software
// date 20 June 2022
// copyright 2022 prr, azul software
//

package main

import (
	"fmt"
	"os"
	"google/gdoc/htmlLib"
)


func main() {

	numArg := len(os.Args)

	outfilNam :=""
	switch numArg {
	case 1:
		fmt.Println("error -- no output file name provided!")
		fmt.Println("usage is ./CvtMdToHtml outfil!")
		os.Exit(1)
	case 2:
		outfilNam = os.Args[1]

	default:
		fmt.Println("error -- too many arguments!")
		fmt.Println("usage is ./CvtMdToHtml outfil!")
		os.Exit(1)

	}

	htmlFilNam := "output/htmlTest/" + outfilNam + ".html"

	outfil, err := os.Create(htmlFilNam)
	if err != nil {
		fmt.Printf("error os.Create file %s: %v\n", htmlFilNam, err)
		os.Exit(1)
	}

	// html
	outstr := htmlLib.CreHtmlHead()

	outstr += htmlLib.CreHtmlMid()

	outstr += htmlLib.CreHtmlEnd()

	outfil.WriteString(outstr)

	outfil.Close()

	fmt.Println("**** success *****")
}
