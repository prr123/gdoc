//
// CvtGdocToJson
// adapted from GdocToDom
// date: 10 Ocober 2022
//
// author: prr
// copyright 2022 prr azul software
//
//		gdocHtml "google/gdoc/gdocHtml"
//

package main

import (
        "fmt"
        "os"
//		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocJson "google/gdoc/gdocJson"
		util "google/gdoc/utilLib"
)


func main() {
	var gd gdocApi.GdocApiStruct
	// initialise default values
	baseFolder := "output/"
//	baseFolderSlash := baseFolder + "/"
	dbg := false
    numArgs := len(os.Args)
//	fmt.Printf("args: %d\n", numArgs)
	cmd := os.Args[0]

	if numArgs < 2 {
       		fmt.Println("error - no docid provided!")
          	fmt.Printf("%s usage is:\n  %s docId [outfolder] [dbg] [summary, main, all]\n", cmd[2:], cmd)
        	os.Exit(1)
	}

    docId := os.Args[1]
	flags := []string{"baseFolder", "out","dbg","sel"}

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

	if baseStr != "none" { baseFolder = baseStr}

	fmt.Printf("baseFolder: %s\n", baseFolder)

	outStr, ok := cliMap["out"].(string)
	if !ok {
		fmt.Println("invalid argument for outFolder: ", outStr)
		os.Exit(1)
	}
	outFolder := "json"
	if outStr != "none" { outFolder = outStr}

	fmt.Printf("outFolder: %s\n", outFolder)

	selStr, ok := cliMap["sel"].(string)
	if !ok {
		fmt.Println("invalid argument for selection: ", selStr)
		os.Exit(1)
	}

	sel := "all"
	if selStr != "none" { sel = selStr}

	fmt.Printf("select: %s\n", sel)


//	os.Exit(2)

	err = gd.InitGdocApi()
	if err != nil {
		fmt.Printf("error - InitGdocApi: %v!", err)
		os.Exit(1)
	}

	srv := gd.Svc

	outfilPath:= baseFolder + outFolder

	fmt.Printf("*************** CvtGdocToJSON ************\n")
	if (dbg) {fmt.Printf("*** debug ***\n")}
	fmt.Printf("output folder: %s selection: %s\n", outfilPath, sel)

	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("Doc Title: %s Selection: %s\n", doc.Title, sel)

	switch sel {
	case "heading":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
//		err = gdocHtml.CreGdocHtmlSection("", outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main: CvtGdocJsonSummary -- cannot convert gdoc doc: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success summary ***!")
		os.Exit(0)

	case "main":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
		if err != nil {
			fmt.Println("error main CvtGdocJsonMain -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success main ***!")
		os.Exit(0)

	case "doc":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
		if err != nil {
			fmt.Println("error CvtGdocJsonDoc -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success doc ***!")
		os.Exit(0)

	case "all":
		err = gdocJson.CreGdocJsonAll(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocJsonAll -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success all ***!")
		os.Exit(0)

	default:
		fmt.Printf("%s is not a valid comand line opt!\n", sel)
		fmt.Println("exiting!")
		os.Exit(1)
	}

	fmt.Println("Success!")
}
