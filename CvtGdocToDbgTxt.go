// program that converts gdoc doc to a text doc
// author prr
// created 12/2021
//
// 9/12/21
// split into two packages
//
// CvtGdocToTxtv1.go from rd_Gdocv4.go
// v1 move test to Gdoc
// v2 change cli interface
//

package main

import (
        "fmt"
        "os"
//		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocDbg "google/gdoc/gdocDbgTxt"
		util "google/gdoc/utilLib"
)


func main() {
	var gd gdocApi.GdocApiStruct
	// intialise
    baseFolder := "output/"

    numArgs := len(os.Args)
    dbg := false
	cmd := os.Args[0]

    if numArgs < 2 {
            fmt.Println("error - no docid provided!")
            fmt.Printf("%s usage is:\n  %s docId [outfolder] [dbg] [summary, main, all]\n", cmd[2:], cmd)
            os.Exit(1)
    }

    docId := os.Args[1]
    flags := []string{"baseFolder", "out","dbg"}

    cliMap, err :=util.ParseFlags(os.Args, flags)
    if err !=nil {
        fmt.Printf("error - CLI: ParseFlags: %v!\n", err)
        os.Exit(1)
    }

    dbgStr, ok := cliMap["dbg"].(string)
    if !ok {
        fmt.Println("invalid argument for dbg: ", dbgStr)
        os.Exit(1)
    }

    if dbgStr == "true" || dbgStr == "none" {dbg = true}

    fmt.Printf("cliMap: %v!\n", cliMap)
    fmt.Printf("dbg: %t\n", dbg)

    baseStr, ok := cliMap["baseFolder"].(string)
    if !ok {
        fmt.Println("invalid argument for baseFolder: ", baseStr)
        os.Exit(1)
    }

	if len(baseStr) > 0 {
		if baseStr != "none" {
			baseFolder = baseStr
		}
	}
    fmt.Printf("baseFolder: %s\n", baseFolder)

    outStr, ok := cliMap["out"].(string)
    if !ok {
        fmt.Println("invalid argument for outFolder: ", outStr)
        os.Exit(1)
    }

    outFolder := "json"
    if outStr != "none" { outFolder = outStr}

    fmt.Printf("outFolder: %s\n", outFolder)

	err = gd.InitGdocApi()
    if err != nil {
        fmt.Printf("error - InitGdocApi: %v!", err)
        os.Exit(1)
    }
	srv := gd.Svc

	outfilPath:= baseFolder + outFolder

//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"


	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}

    fmt.Printf("*************** Gdoc Debug ************\n")
	fmt.Printf("The title of the doc is: %s\n", doc.Title)
	fmt.Printf("Destination folder: %s\n", outfilPath)

	err = gdocDbg.CvtGdocToDbgTxt(outfilPath, doc, nil)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}
