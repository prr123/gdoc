// CvtTxtToGdoc.go
// program that converts text doc to a gdoc file
// author prr
// created 02/06/20212
//
// copyright 2022 prr, Azul Software//
//
// for license description and documentation: 
// see github.com/prr123/gdoc

package main

import (
        "fmt"
        "os"
//		"strings"
		txtGdoc "google/gdoc/txtGdoc"
)


func main() {

	// intialise
//    baseFolder := "azul"
//    baseFolderSlash := baseFolder + "/"

    numArgs := len(os.Args)

	cmd := os.Args[0]

 	switch numArgs {
        case 1:
            fmt.Println("error - no comand line arguments!")
            fmt.Printf("%s usage is:\n  %s \"input file\" folder\n", cmd[2:], cmd)
            os.Exit(1)
		case 2:
		// output folder
		default:
            fmt.Println("error - too many arguments!")
            fmt.Printf("%s usage is:\n  %s \"input file\" folder\n", cmd[2:], cmd)
            os.Exit(1)
	}

    inpFil := os.Args[1]

/*
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
            fmt.Printf("no valid output folder provided: %s", os.Args[2])
            os.Exit(1)
    }
*/

	_, err := txtGdoc.InitTxtGdoc(inpFil)
    if err != nil {
        fmt.Printf("error - InitTxtGdoc: %v!\n", err)
        os.Exit(1)
    }

//	gd.OutFilPath = outfilPath

	fmt.Println("Success!")
}
