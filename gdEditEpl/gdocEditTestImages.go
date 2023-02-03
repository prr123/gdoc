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

//	hdTxt := "h1"
	imgList, err := geObj.ListImages()
    if imgList == nil {
        if err == nil {fmt.Printf("no images found!\n")} else {fmt.Printf("error ListImages: %v\n", err)}
    } else {
        fmt.Printf("found %d images!\n", len(*imgList))
        for i:=0; i< len(*imgList); i++ {
			imgTyp:= true
			if !(*imgList)[i].Typ {imgTyp = false}
            fmt.Printf("Image[%d]: Id: %s Typ: %t\n", i, (*imgList)[i].Id, imgTyp)
//            fmt.Printf("Image[%d]: el: %d Start: %d End:%d \"%s\"\n", i, (*secList)[i].El, (*secList)[i].Start, (*secList)[i].End, (*secList)[i].Typ)
        }
    }

	gdEdit.PrintImgList(imgList)

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


