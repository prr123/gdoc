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
       		fmt.Println("error - no opt argument provided!")
			fmt.Println("assuming opt = all")
			opt = "all"
		case 3:
			opt = os.Args[2]
		default:
        	fmt.Println("error - too many arguments!")
          	fmt.Printf("%s usage is:\n  %s docId opt [sumary, main, all]\n", cmd[2:], cmd)
        	os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.InitGdocApi()
	srv := gd.Svc

	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	switch opt {
	case "summary":
		outFilNam := fmt.Sprintf("output/%s_summary", doc.Title)
		outfil, err := gd.CreOutFile(outFilNam, "html")
		if err != nil {
			fmt.Println("error creOutFile -- cannot open create outfile: ", err)
			os.Exit(1)
		}

		err = gdocHtml.CreGdocHtmlSection("summary", outfil, doc, nil)
		if err != nil {
			fmt.Println("error main: CreGdocHtmlSummary -- cannot convert gdoc doc: ", err)
			os.Exit(1)
		}
		outfil.Close()
		fmt.Println("*** success summary ***!")
		os.Exit(0)

	case "main":
		outFilNam := fmt.Sprintf("output/%s_main", doc.Title)
		outfil, err := gd.CreOutFile(outFilNam, "html")
		if err != nil {
			fmt.Println("error creOutFile-- cannot open create outfile: ", err)
			os.Exit(1)
		}

		err = gdocHtml.CreGdocHtmlMain(outfil, doc, nil)
		if err != nil {
			fmt.Println("error main CreGdocHtmlMain -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}
		outfil.Close()
		fmt.Println("*** success main ***!")
		os.Exit(0)

	case "all":
		outFilNam := fmt.Sprintf("output/%s", doc.Title)
		outfil, err := gd.CreOutFile(outFilNam, "html")
		if err != nil {
			fmt.Println("error creOutFil - cannot open create outfile: ", err)
			os.Exit(1)
		}

		err = gdocHtml.CreGdocHtmlDoc(outfil, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocHtmlSummary -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		outfil.Close()
		fmt.Println("*** success all ***!")
		os.Exit(0)

	default:
		fmt.Printf("did not provide a valid opt: %s\n", opt)
		fmt.Println("failure!")
		os.Exit(1)
	}
	
}

