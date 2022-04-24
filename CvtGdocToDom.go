//
// CvtGdocToDom
// date: 22. April 2022
// author: prr
// copyright 2022 prr azul software
//

package main

import (
        "fmt"
        "os"
		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocHtml "google/gdoc/gdocHtml"
		gdocDom "google/gdoc/gdocDom"
)


func main() {
	var gd gdocApi.GdocApiStruct

    numArgs := len(os.Args)
//	fmt.Printf("args: %d\n", numArgs)
	cmd := os.Args[0]
	opt:=""
	switch numArgs {
		case 1:
       		fmt.Println("error - no comand line arguments!")
 			fmt.Printf("%s usage is:\n  %s docId opt [sumary, main, all]\n", cmd[2:], cmd)
        	os.Exit(1)
		case 2:
       		fmt.Println("no option argument provided!")
			fmt.Println("assuming 'opt = all'!")
			opt = "all"
		case 3:

		case 4:
			opt = os.Args[3]

		default:
        	fmt.Println("error - too many arguments!")
          	fmt.Printf("%s usage is:\n  %s docId opt [sumary, main, all]\n", cmd[2:], cmd)
        	os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.InitGdocApi()
	srv := gd.Svc

	outfilPath:= ""
	switch {
		case os.Args[2] == "output":
			outfilPath = os.Args[2]
		case strings.Index(os.Args[2], "output/")< 0:
 			outfilPath = "output/" + os.Args[2]
		case strings.Index(os.Args[2], "output/") == 0:
			outfilPath = os.Args[2]
		case os.Args[2] == "":
			outfilPath = "output"
		default:
			fmt.Printf("no valid input folder: %s", os.Args[2])
			os.Exit(1)
	}

	fmt.Printf("*************** CctGdocToHtml ************\n")
	fmt.Printf("output folder: %s option: %s\n", outfilPath, opt)

	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("Doc Title: %s Option: %s\n", doc.Title, opt)

	switch opt {
	case "heading":
		err = gdocHtml.CreGdocHtmlSection("", outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main: CreGdocHtmlSummary -- cannot convert gdoc doc: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success summary ***!")
		os.Exit(0)

	case "main":
		err = gdocHtml.CreGdocHtmlMain(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main CreGdocHtmlMain -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success main ***!")
		os.Exit(0)

	case "doc":
		err = gdocHtml.CreGdocHtmlDoc(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocHtmlDoc -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success all ***!")
		os.Exit(0)

	case "all":
		err = gdocHtml.CreGdocHtmlAll(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocHtmlAll -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success all ***!")
		os.Exit(0)

	default:
		fmt.Printf("did not provide a valid opt: %s\n", opt)
		fmt.Println("failure!")
		os.Exit(1)
	}
	
}
