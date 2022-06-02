//
// CvtGdocToDom
// date: 22. April 2022
// author: prr
// copyright 2022 prr azul software
//
//		gdocHtml "google/gdoc/gdocHtml"


package main

import (
        "fmt"
        "os"
		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocDom "google/gdoc/gdocDom"
//		gdocHtml "google/gdoc/gdocHtml"
)


func main() {
	var gd gdocApi.GdocApiStruct
	// initialise default values
	baseFolder := "output"
	baseFolderSlash := baseFolder + "/"
	opt:=""

    numArgs := len(os.Args)
//	fmt.Printf("args: %d\n", numArgs)
	cmd := os.Args[0]


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
	if err != nil {
		fmt.Printf("error - InitGdocApi: %v!", err)
		os.Exit(1)
	}

	srv := gd.Svc

	outfilPath:= ""
	switch {
		case numArgs == 2:
			outfilPath = baseFolder
		case os.Args[2] == baseFolder:
			outfilPath = os.Args[2]
		case strings.Index(os.Args[2], baseFolderSlash)< 0:
 			outfilPath = baseFolderSlash + os.Args[2]
		case strings.Index(os.Args[2], baseFolderSlash) == 0:
			outfilPath = os.Args[2]
		case os.Args[2] == "":
			outfilPath = baseFolder
		default:
			fmt.Printf("no valid input folder: %s", os.Args[2])
			os.Exit(1)
	}

	fmt.Printf("*************** CvtGdocToDom ************\n")
	fmt.Printf("output folder: %s option: %s\n", outfilPath, opt)

	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("Doc Title: %s Option: %s\n", doc.Title, opt)

	switch opt {
	case "heading":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
//		err = gdocHtml.CreGdocHtmlSection("", outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main: CreGdocDomSummary -- cannot convert gdoc doc: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success summary ***!")
		os.Exit(0)

	case "main":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
//		err = gdocHtml.CreGdocHtmlMain(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main CreGdocDomMain -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success main ***!")
		os.Exit(0)

	case "doc":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
//		err = gdocHtml.CreGdocHtmlDoc(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocDomDoc -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success doc ***!")
		os.Exit(0)

	case "all":
		err = gdocDom.CreGdocDomAll(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocDomAll -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success all ***!")
		os.Exit(0)

	default:
		fmt.Printf("%s is not a valid comand line opt!\n", opt)
		fmt.Println("exiting!")
		os.Exit(1)
	}
	fmt.Println("Success!")
}
