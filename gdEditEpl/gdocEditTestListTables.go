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
		fmt.Printf("usage is: gdocEditTestListTables docid\n")
		os.Exit(1)
	}

	if numArg != 2 {
		fmt.Printf("error: numArg [%d] not 2!\n", numArg)
		fmt.Printf("usage is: gdocEditTestListTables docid\n")
		os.Exit(-1)
	}

	docId := os.Args[1]

//	heading := os.Args[2]

//fmt.Printf("heading: %s\n", os.Args[2])

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

	tblList, err := geObj.ListTables()
    if tblList == nil {
        if err == nil {fmt.Printf("search text not found!\n")} else {fmt.Printf("error FindTxtNext: %v\n", err)}
    } else {
        fmt.Printf("found %d tables!\n", len(*tblList))
        for i:=0; i< len(*tblList); i++ {
            fmt.Printf("table[%d]: el: %d pos: %d rows: %d cols:%d\n", i, (*tblList)[i].El, (*tblList)[i].Start, (*tblList)[i].Rows, (*tblList)[i].Cols)
        }
    }

/*
    for i:=0; i< len(*tblList); i++ {
        outstr, err := geObj.DispPar((*posList)[i])
        if err != nil {
            fmt.Printf("error DispText: %v\n", err); os.Exit(-1);
        }
		fmt.Printf("par[%d]: %s\n", i, outstr)
    }
*/
}


