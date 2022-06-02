// txtGdocLib.go
// golang library that creates a gdoc file from a txt file
// author: prr
// created: 2/6/2022
// copyright 2022 prr, Azul Software
//
// for license description and documentation:
// see github.com/prr123/gdoc
//
// start: CvtTxtToGdoc.go
//


package txtGdocLib

import (
    "fmt"
    "os"
    "google.golang.org/api/docs/v1"
//    gdocUtil "google/gdoc/gdocUtil"
	gdocApi "google/gdoc/gdocApiRW"
	util "google/gdoc/util"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)

type TxtGdocObj struct {
    doc *docs.Document
	inpFil *os.File
	OutfilPath string
}

func InitTxtGdoc(title string) (dObj *TxtGdocObj, err error) {

	var doc docs.Document
	var ndoc *docs.Document
	var gd TxtGdocObj

	if len(title) == 0 {
		return nil, fmt.Errorf("title has no string!")
	}

	inbyte := []byte(title)
	for i:=0; i< len(inbyte); i++ {
		switch inbyte[i] {
		case ' ':
			inbyte[i] = '_'
		case '.':
			return nil, fmt.Errorf("title string has period!")

		default:
			if !util.IsAlphaNumeric(inbyte[i]) {
				return nil, fmt.Errorf("title string has unacceptable char %c!", inbyte[i])
			}

		}
	}

	doc.Title = title

	inpfilnam:= title + ".txt"

	inpFil, err := os.Open(inpfilnam)
	if err != nil {
		return nil, fmt.Errorf("os.Open %v", err)
	}

	gd.inpFil = inpFil

    // need to create a minimal doc first
    srv, err := gdocApi.InitGdocRWApi()
	if err != nil {
		return nil, fmt.Errorf("initGocRWApi: %v", err)
	}

    ndoc, err = srv.Documents.Create(&doc).Do()
    if err != nil {
        fmt.Println("Unable to retrieve data from document: ", err)
        os.Exit(1)
    }
	gd.doc = ndoc


    fmt.Printf("*************** CvtGdocToTxt ************\n")
    fmt.Printf("The title of the doc is: %s\n", ndoc.Title)
//    fmt.Printf("Destination folder: %s\n", outfilPath)


	return &gd, nil
}
