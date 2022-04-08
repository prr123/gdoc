// showGdoc
// author: prr
// date: 8. April 2022
// copyright 2022 prr
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
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
        fmt.Println("CvtGdocToHtmlv1 usage is:\n  CvtGdocToTxtv1 docId\n")
        os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.InitGdocApi()
	srv := gd.Svc
//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	err = gdocHtml.ShowGdocSummmary(doc)
	if err != nil {
		fmt.Println("error main -- cannot show gdoc file: ", err)
		os.Exit(1)
	}

	fmt.Println("success!")
	os.Exit(0)
}

