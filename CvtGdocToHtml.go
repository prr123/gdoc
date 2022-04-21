// v1 change doc to Par NTest
// v2 start a text debug
// v3 loop through some elements
// 9/12/21
// split into two packages
// CvtGdocToTxt
// CvtGdocToHtml
//
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi"
		gdocHtml "google/gdoc/gdocHtml"
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

	outfilPath := "output" + os.Args[2]
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
