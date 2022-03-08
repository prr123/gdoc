// CvtGdocToMdv1.go
// author: prr
// created 15/12/2021
// copyright 2021, 2022 prr
//
// 15/12/21
// split from GdocToTxtv1
//
// CvtGdocToTxtv1.go from rd_Gdocv4.go
// v1 move test to Gdoc
//
// 6/3/2022 add inline images
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi"
		gdocMd "google/gdoc/gdocMd"
)


func main() {
	var gd gdocApi.GdocApiStruct
    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
          fmt.Println("CvtGdocToTxtv1 usage is:\n  CvtGdocToTxtv1 docId\n")
        os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.Init()
	srv := gd.Svc
//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document! ", err)
		os.Exit(1)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	outfil, err := gd.CreTxtOutFile(doc.Title, "md")
	if err != nil {
		fmt.Println("error main -- cannot open md file! ", err)
		os.Exit(1)
	}
	err = gdocMd.CvtGdocToMd(outfil, doc, false)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file! ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("Success!")
	os.Exit(0)
}
