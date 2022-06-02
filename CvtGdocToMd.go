// CvtGdocToMd.go
// program that converts a gdoc file to a markdown file
//
// author: prr
// created 15/12/2021
// copyright 2021, 2022 prr
//
// 6/3/2022 add inline images
//

package main

import (
        "fmt"
        "os"
		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocMd "google/gdoc/gdocMd"
)


func main() {
	var gd gdocApi.GdocApiStruct
	opt := ""
    // initialise default values
	baseFolder := "output"
    baseFolderSlash := baseFolder + "/"

    numArgs := len(os.Args)
    cmd := os.Args[0]

    switch numArgs {
        case 1:
            fmt.Println("error - no comand line arguments!")
            fmt.Printf("%s usage is:\n  %s docId opt [sumary, main, all]\n", cmd[2:], cmd)
            os.Exit(1)
        case 2:
            fmt.Println("info - no option argument provided!")
            fmt.Println("assuming 'opt = all'!")
            opt = "all"
        case 3:

        case 4:
            opt = os.Args[3]

        default:
            fmt.Println("error - too many arguments!")
            fmt.Printf("%s usage is:\n  %s docId opt [sumary, main, all]\n", cmd[2:], cmd)
            os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.InitGdocApi()
    if err != nil {
        fmt.Printf("error - InitGdocApi: %v!", err)
        os.Exit(1)
    }

	srv := gd.Svc

   outfilPath:= ""
    switch {
        case numArgs == 2:
            outfilPath = baseFolder
        case os.Args[2] == baseFolder:
            outfilPath = os.Args[2]
        case strings.Index(os.Args[2], baseFolderSlash)< 0:
            outfilPath = baseFolderSlash + os.Args[2]
        case strings.Index(os.Args[2], baseFolderSlash) == 0:
            outfilPath = os.Args[2]
        case os.Args[2] == "":
            outfilPath = baseFolder
        default:
            fmt.Printf("error - no valid input folder: %s", os.Args[2])
            os.Exit(1)
    }

//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("error srv.Documents.Get: Unable to retrieve document! ", err)
		os.Exit(1)
	}

    fmt.Printf("*************** CvtGdocToMd ************\n")
    fmt.Printf("The title of the doc is: %s\n", doc.Title)
    fmt.Printf("Destination folder: %s\n", outfilPath)

	switch opt {
	case "new":
	
	default:
		err = gdocMd.CvtGdocToMd(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CvtGdocToMd -- cannot convert gdoc file! ", err)
			os.Exit(1)
		}
	}
	fmt.Println("Success!")
}
