// v1 change doc to Par NTest
// v2 start a text debug
// v3 loop through some elements
// 9/12/21
// split into two packages
// CvtGdocToTxt
// CvtGdocToHtml
// 
// CvtGdocToTxtv1.go from rd_Gdocv4.go
// v1 move test to Gdoc
// 19/12 - add docId to input arguments

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
    if numArgs < 2 {
		cmd := os.Args[0]
        fmt.Println("error - no comand line arguments!")
          fmt.Printf("%s usage is:\n  %s docId\n", cmd[2:], cmd)
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

	outFilNam := fmt.Sprintf("output/%s", doc.Title)
	outfil, err := gd.CreOutFile(outFilNam, "html")
	if err != nil {
		fmt.Println("error main -- cannot open create outfile: ", err)
		os.Exit(1)
	}
	err = gdocHtml.CreGdocHtmlFil(outfil, doc, nil)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("success!")
	os.Exit(0)
}

