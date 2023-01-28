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
	var testStart int64

	numArg := len(os.Args)

	if numArg == 1 {
		fmt.Println("error: insufficient arguments")
		fmt.Printf("usage is: gdocEditTestFind docid start\n")
		os.Exit(1)
	}

	if numArg != 3 {
		fmt.Printf("error: numArg [%d] not 3!\n", numArg)
		fmt.Printf("usage is: gdocEditTestFind docid start\n")
		os.Exit(-1)
	}

	docId := os.Args[1]

//fmt.Printf("start: %s\n", os.Args[2])
	_, err := fmt.Sscanf(os.Args[2], "%d", &testStart)
	if err != nil {
		fmt.Printf("error converting cmd arg 2: %d\n", os.Args[2])
		os.Exit(-1)
	}

    err = gd.InitDriveApi()
    if err != nil {
        fmt.Printf("error main::Init gdriveApi: %v\n", err)
        os.Exit(1)
    }

	geObj, err := gdEdit.InitGdocEdit(gd.GdocSvc.Documents, docId)
    if err != nil {
        fmt.Printf("Unable to init GdEdit: %v\n", err); os.Exit(-1);
    }

    fmt.Printf("geObj: Doc Title: %s\n", geObj.Doc.Title)

	pos, err := geObj.FindTextNext("test", testStart)
	if pos == -1 {
		if err == nil {fmt.Printf("search text not found!\n")} else {fmt.Printf("error FindTxtNext: %v\n", err)}
	} else {fmt.Printf("pos: %d!\n", pos)}

	err = geObj.DispText(pos)
    if err != nil {
        fmt.Printf("error DispText: %v\n", err); os.Exit(-1);
    }

}


