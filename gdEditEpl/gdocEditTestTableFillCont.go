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

	tblList, err := geObj.ListTables()
    if tblList == nil {
        if err == nil {fmt.Printf("search text not found!\n")} else {fmt.Printf("error FindTxtNext: %v\n", err)}
    } else {
        fmt.Printf("found %d tables!\n", len(*tblList))
        for i:=0; i< len(*tblList); i++ {
            fmt.Printf("table[%d]: el: %d pos: %d rows: %d cols:%d\n", i, (*tblList)[i].El, (*tblList)[i].Start, (*tblList)[i].Rows, (*tblList)[i].Cols)
        }
    }

	tbl:=0
	tblObj := (*tblList)[tbl]
fmt.Printf("\n***** get Content from table %d ****\n", tbl)
fmt.Printf("table[%d]: el: %d pos: %d rows: %d cols:%d\n", tbl, tblObj.El, tblObj.Start, tblObj.Rows, tblObj.Cols)

	contTab:= make([][]string,tblObj.Rows)
    for irow:=0; irow<tblObj.Rows; irow++ {contTab[irow] = make([]string, tblObj.Cols)}

    for irow:=0; irow<tblObj.Rows; irow++ {
		for icol:=0; icol<tblObj.Cols; icol++ {
			contTab[irow][icol] = fmt.Sprintf("R%dC%d", irow, icol)
		}
	}
	fmt.Printf("generated table")
	gdEdit.PrintTbl(&contTab)

	updreq, err :=geObj.FillTblContent(&contTab, &tblObj)
    if err != nil {
        fmt.Printf("Error GetTblContent: %v\n", err); os.Exit(-1);
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

}


