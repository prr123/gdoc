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
		fmt.Printf("usage is: gdocEditTestTableGetCont docid\n")
		os.Exit(1)
	}

	if numArg != 2 {
		fmt.Printf("error: numArg [%d] not 2!\n", numArg)
		fmt.Printf("usage is: gdocEditTestTableGetCont docid\n")
		os.Exit(-1)
	}

	docId := os.Args[1]

    err := gd.InitDriveApi()
    if err != nil {
        fmt.Printf("error main::Init gdriveApi: %v\n", err)
        os.Exit(1)
    }

	geObj, err := gdEdit.InitGdocEdit(gd.GdocSvc.Documents, docId)
    if err != nil {
        fmt.Printf("Unable to init GdEdit: %v\n", err); os.Exit(-1);
    }

    fmt.Printf("geObj: Doc Title: %s\n", geObj.Doc.Title)

	hdTxt := "h1"
	hdList, err := geObj.ListHeadings(hdTxt)
    if hdList == nil {
        if err == nil {fmt.Printf("no headings %s found!\n, hdTxt")} else {fmt.Printf("error FindTxtNext: %v\n", err)}
    } else {
        fmt.Printf("found %d headers!\n", len(*hdList))
        for i:=0; i< len(*hdList); i++ {
            fmt.Printf("Header[%d]: el: %d Start: %d End:%d \"%s\"\n", i, (*hdList)[i].El, (*hdList)[i].Start, (*hdList)[i].End, (*hdList)[i].Text)
        }
    }

/*
	tbl:=0
	tblObj := (*tblList)[tbl]

	updreq, err :=geObj.ClearTblContent(&tblObj)
    if err != nil {
        fmt.Printf("Error GetTblContent: %v\n", err); os.Exit(-1);
    }

	fmt.Printf("upd requests: %d\n", len(updreq.Requests))

	if len(updreq.Requests) == 0 {
		fmt.Printf("no update del requests!\n")
		os.Exit(0)
	}
    _, err = geObj.DocSvc.BatchUpdate(docId, updreq).Do()
    if err != nil {
        fmt.Printf("Unable to update document: %v\n", err)
        os.Exit(-1)
    }

	// need to pull updated doc

    geObj, err = gdEdit.InitGdocEdit(gd.GdocSvc.Documents, docId)
    if err != nil {
        fmt.Printf("Unable to init GdEdit: %v\n", err); os.Exit(-1);
    }

	contTbl, err :=geObj.GetTblContent(&tblObj)
    if err != nil {
        fmt.Printf("Error GetTblContent: %v\n", err); os.Exit(-1);
    }

	gdEdit.PrintTbl(contTbl)
*/
}


