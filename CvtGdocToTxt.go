// program that converts gdoc doc to a text doc
// author prr
// created 12/2021
//
// 9/12/21
// split into two packages
//
// CvtGdocToTxtv1.go from rd_Gdocv4.go
// v1 move test to Gdoc
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi"
		gdocTxt "google/gdoc/gdocTxt"
)


func main() {
	var gd gdocApi.GdocApiStruct

    numArgs := len(os.Args)

	cmd := os.Args[0]
    if numArgs < 2 {
        fmt.Printf("*** error %s:: -- no gdoc id!\n", string(cmd[2:]))
  		  fmt.Printf(" usage is: \"%s docId\"\n", cmd)
        os.Exit(2)
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

	outFilNam := fmt.Sprintf("output/%s", doc.Title)
	outfil, err := gd.CreOutFile(outFilNam, "txt")
	if err != nil {
		fmt.Println("error main -- cannot open out file: ", err)
		os.Exit(1)
	}

	err = gdocTxt.CvtGdocToTxt(outfil, doc)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("Success!")
	os.Exit(0)
}
