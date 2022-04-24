// 22/1/2022
// CopyGdoc
// author prr
// from CreateGdoc.go
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi2"
//		gdocEdit "google/gdoc/gdocEdit"
)


func main() {
	var gd gdocApi.GdocApi2Obj
	var docId string
/*
    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
        fmt.Println("GetGdoc usage is: \"GetGdoc docId\"")
		docId = "1lEodX98Eq6_2elpgct_OOv-5L5Es_iGyZJqrIS2BznY"
		fmt.Printf("defaulting to docId: %s\n\n", docId)
    } else {
    	docId = os.Args[1]
	}

	fmt.Println("Fetching google doc with id: ", docId)
*/

	err := gd.Init()
	if err != nil {
		fmt.Printf("error Init: %v \n", err)
		fmt.Println("exiting due to error!")
		os.Exit(1)
	}

	err = gd.Gdoc_create("title")
	if err != nil {
		fmt.Println("error Gdoc_create: ", err)
		os.Exit(1)
	}
	fmt.Println("success creating Google Doc!")
	fmt.Println("reading new googe doc with Id: ", gd.DocId)

	err = gd.Gdoc_read(docId)
	fmt.Println("doc title: ",gd.Doc.Title)
	os.Exit(0)

/*
	srv := gd.Svc
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)


	outfil, err := gd.CreTxtOutFile(doc.Title, "toml")
	if err != nil {
		fmt.Println("error main -- cannot open out file: ", err)
		os.Exit(1)
	}
	err = gdocHtml.CvtGdocTxt(outfil, doc)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("success!")
	os.Exit(0)
*/
}

