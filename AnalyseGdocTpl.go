// golang program that parses a gdoc document template
// author: prr, azul software
// created: 18/12/2022
// copyright 2022 prr, Peter Riemenschneider, Azul Software
//
// start: AnalyseGdocTpl
//

package main

import (
	"os"
	"fmt"
    gdocApi "google/gdoc/gdocApi"
	gdocTpl "google/gdoc/gdocTpl"
    util "google/gdoc/utilLib"
)

func main() {
    var docId string

    numArgs:= len(os.Args)

    if numArgs < 2 {
        fmt.Printf("error - exit: insufficient command line arguments\n")
        fmt.Printf("usage is: AnalyseGdocTpl \"docId\" [\\out=] [\\dbg]\n")
        os.Exit(-1)
    }

    docId =os.Args[1]

    flags := [] string {"out", "dbg"}

    argmap, err := util.ParseFlagsStart(os.Args, flags,2)
    if err != nil {fmt.Printf("error ParseFlags: %v\n", err); os.Exit(-1);}

    outFilNam, ok := argmap["out"]
    if !ok {fmt.Printf("error no out Filnam provided!\n"); os.Exit(-1);}

    outFilNamStr := outFilNam.(string)
fmt.Printf("out file: %s\n",outFilNamStr)


    gdoc, err := gdocApi.InitGdocApiV2()
    if err != nil {
        fmt.Printf("error - InitGdocApiV2: %v!", err)
        os.Exit(1)
    }

	srv := gdoc.Svc

    doc, err := srv.Documents.Get(docId).Do()
    if err != nil {
        fmt.Println("Unable to retrieve data from document: ", err)
        os.Exit(1)
    }

    fmt.Printf("Doc Title: %s\n", doc.Title)

	tplObj := gdocTpl.InitTpl(doc)

	fmt.Printf("gdObj: %v\n", tplObj)


}
