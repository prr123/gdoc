// golang program that updates a gdoc document
// author: prr, azul software
// created: 19/1/2023
// copyright 2023 prr, Peter Riemenschneider, Azul Software
//
// start: tableBatchUpd
//

package main

import (
    "os"
    "fmt"
//  "bytes"
//    gdocApi "google/gdoc/gdocApi"
//    gdocTpl "google/gdoc/gdocTpl"
//    "google.golang.org/api/docs/v1"
    gdrive "google/gdoc/gdriveLib"
	gdEdit "google/gdoc/gdocEdit"
//    util "google/gdoc/utilLib"
)

func main() {
    var gd gdrive.GdApiObj

	numArg := len(os.Args)

	if numArg == 1 {
		fmt.Println("error: insufficient arguments")
		fmt.Printf("usage is: testBatchUpd docid\n")
		os.Exit(1)
	}

	if numArg != 2 {
		fmt.Printf("error: numArg [%d] not 2!\n", numArg)
		fmt.Printf("usage is: testBatchUpd docid\n")
		os.Exit(-1)
	}

	docId := os.Args[1]

    err := gd.InitDriveApi()
    if err != nil {
        fmt.Printf("error main::Init gdriveApi: %v\n", err)
        os.Exit(1)
    }

    svc := gd.GdocSvc

    doc, err := svc.Documents.Get(docId).Do()
    if err != nil {
        fmt.Printf("Unable to retrieve document: %v\n", err); os.Exit(-1);
    }
    gd.Doc = doc

    fmt.Printf("GdocSvc Doc Title: %s\n", doc.Title)

	geObj, err := gdEdit.InitGdocEdit(svc.Documents, docId)
    if err != nil {
        fmt.Printf("Unable to init GdEdit: %v\n", err); os.Exit(-1);
    }

    fmt.Printf("geObj: Doc Title: %s\n", geObj.Doc.Title)

}


