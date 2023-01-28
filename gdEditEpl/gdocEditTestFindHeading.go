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
		fmt.Printf("usage is: gdocEditTestFindHeading docid heading\n")
		os.Exit(1)
	}

	if numArg != 3 {
		fmt.Printf("error: numArg [%d] not 3!\n", numArg)
		fmt.Printf("usage is: gdocEditTestFindHeading docid heading\n")
		os.Exit(-1)
	}

	docId := os.Args[1]

	heading := os.Args[2]

fmt.Printf("heading: %s\n", os.Args[2])

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

	posList, err := geObj.FindHeadingsAll(heading)
    if posList == nil {
        if err == nil {fmt.Printf("search text not found!\n")} else {fmt.Printf("error FindTxtNext: %v\n", err)}
    } else {
        fmt.Printf("found %d instances!\n", len(*posList))
        for i:=0; i< len(*posList); i++ {
            fmt.Printf("pos[%d]: %d\n", i, (*posList)[i])
        }
    }

    for i:=0; i< len(*posList); i++ {
        outstr, err := geObj.DispPar((*posList)[i])
        if err != nil {
            fmt.Printf("error DispText: %v\n", err); os.Exit(-1);
        }
		fmt.Printf("par[%d]: %s\n", i, outstr)
    }

}


