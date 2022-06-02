// CvtTxtToGdoc
// program that converts text doc to a gdoc file
// author prr
// created 02/06/20212
//
//

package main

import (
        "fmt"
        "os"
		"strings"
		txtGdoc "google/gdoc/txtGdoc"
)


func main() {

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
		// input file
		case 3:
		// output folder
		default:
            fmt.Println("error - too many arguments!")
            fmt.Printf("%s usage is:\n  %s folder docId\n", cmd[2:], cmd)
            os.Exit(1)
	}

    inpFil := os.Args[1]

	doc, err := txtGdoc.InitTxtGdoc(inpFil)
    if err != nil {
        fmt.Printf("error - InitTxtGdoc: %v!", err)
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

	// need to create a minimal doc first
	doc, _ := txtGdoc.CreMinDoc(inpfil)

	doc, err := srv.Documents.Create(doc).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}


    fmt.Printf("*************** CvtGdocToTxt ************\n")
	fmt.Printf("The title of the doc is: %s\n", doc.Title)
	fmt.Printf("Destination folder: %s\n", outfilPath)

/*
	err = txtGdoc.CvtGdocToTxt(outfilPath, doc, nil)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

*/
	fmt.Println("Success!")
}
