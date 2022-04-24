// 13/1/2022
// EditGdoc
// author prr
// from CvtGdocToHtml
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
	docId := gd.DocId
	err = gd.Gdoc_read2()
	if err != nil {
		fmt.Println("error Gdoc_read2: ", err)
		os.Exit(1)
	}
	fmt.Println("doc title: ",gd.Doc.Title)
	os.Exit(0)
}

