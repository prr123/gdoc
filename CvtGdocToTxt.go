// program that converts gdoc doc to a text doc
// author prr
// created 12/2021
//
// 9/12/21
// split into two packages
//
// CvtGdocToTxtv1.go from rd_Gdocv4.go
// v1 move test to Gdoc
//

package main

import (
        "fmt"
        "os"
		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocTxt "google/gdoc/gdocTxt"
)


func main() {
	var gd gdocApi.GdocApiStruct
	// intialise
    baseFolder := "output"
    baseFolderSlash := baseFolder + "/"

    numArgs := len(os.Args)

	cmd := os.Args[0]

 	switch numArgs {
        case 1:
            fmt.Println("error - no comand line arguments!")
            fmt.Printf("%s usage is:\n  %s docId folder\n", cmd[2:], cmd)
            os.Exit(1)
        case 2:
		// doc id
		case 3:
		// output folder
		default:
            fmt.Println("error - too many arguments!")
            fmt.Printf("%s usage is:\n  %s folder docId\n", cmd[2:], cmd)
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
            fmt.Printf("no valid input folder: %s", os.Args[2])
            os.Exit(1)
    }


//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"


	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}

    fmt.Printf("*************** CvtGdocToTxt ************\n")
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	err = gdocTxt.CvtGdocToTxt(outfilPath, doc)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
	os.Exit(0)
}
