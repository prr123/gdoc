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
	"google.golang.org/api/drive/v3"
//    gdocUtil "google/gdoc/gdocUtil"
	gdocApi "google/gdoc/gdocApiRW"
	util "google/gdoc/util"
	gdApi "google/gdoc/gdApi"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)

type TxtGdocObj struct {
    doc *docs.Document
	svc *docs.Service
	drSvc *drive.Service
	inpFil *os.File
	InpFilPath string
	OutFilPath string
}

func InitTxtGdoc(title string) (dObj *TxtGdocObj, err error) {

	var doc docs.Document
	var ndoc *docs.Document
	var gd TxtGdocObj

	gd.InpFilPath = "inpTestTxt"

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

	inpfilnam:= gd.InpFilPath + "/" + title + ".txt"

	inpFil, err := os.Open(inpfilnam)
	if err != nil {
		return nil, fmt.Errorf("os.Open %v", err)
	}

	gd.inpFil = inpFil

    // need to create a minimal doc first
    svc, err := gdocApi.InitGdocRWApi()
	if err != nil {
		return nil, fmt.Errorf("initGocRWApi: %v", err)
	}

    ndoc, err = svc.Documents.Create(&doc).Do()
    if err != nil {
        fmt.Println("Unable to create document: ", err)
        os.Exit(1)
    }

	gd.svc = svc
	gd.doc = ndoc

	gdSvc, err := gdApi.InitDriveApi()
    if err != nil {
        fmt.Println("Unable to start drive service: ", err)
        os.Exit(1)
    }
	gd.drSvc = gdSvc

	// need to add parent

    fmt.Printf("*************** CvtTxtToGdoc ************\n")
    fmt.Printf("The doc title is: %s\n", ndoc.Title)
    fmt.Printf("The doc Id is:    %s\n", ndoc.DocumentId)

//    fmt.Printf("Destination folder: %s\n", outFilPath)

	// convert text into gdoc
	err = gd.CvtTxtFil(true)
	if err != nil {
		return &gd, fmt.Errorf("cvtTxtFil: %v", err)
	}
	inpFil.Close()
	return &gd, nil
}

func (dObj *TxtGdocObj) CvtTxtFil(cr bool) (err error) {
// function that reads input file and converts text into update requests
// the cr parameter is an option that, if set, combines text lines into one paragraph.
// the cr remove the cr at the end of a line and substitutes a white space
// paragaph ends are:
// 1. a text line followed by an empty line
// 2. a text line that is less than x% of the following line


	var eos docs.EndOfSegmentLocation
//	var insTxtReq docs.InsertText
	var updreqs [](*docs.Request)
//	var bUpdResp docs.BatchUpdateDocumentResponse

	infil := dObj.inpFil
	fileInfo, err := infil.Stat()
	if err != nil {return fmt.Errorf("infil.Stat: %v", err)}

	size := fileInfo.Size()

	fmt.Printf("input file size: %d\n", size)

	inBuf := make([]byte, size + 2)

	_,err = infil.Read(inBuf)
	if err != nil {return fmt.Errorf("cannot read input file: %v", err)}

	doc := dObj.doc
	svc := dObj.svc

	docId := doc.DocumentId
	eos.SegmentId = ""

	ilin := 0
	linSt:=0
	linEnd:=0
	txtStr :=""
	parEnd := false
	reqCount := 0
	for i:=0; i<int(size); i++ {

		if inBuf[i] == '\n' {
			linEnd = i
			if cr {
				if i < int(size) - 1 {
				// for all char but last char
					switch {
					case inBuf[i+1] == '\n':
						// next line is an empty line
						txtStr = string(inBuf[linSt:linEnd+1])
						parEnd = true
					case util.IsAlpha(inBuf[i+1]):
						// next line is a text line
						if linEnd > linSt {
							// if the line contains text and then add text to par string
							inBuf[linEnd] = ' '
							txtStr += string(inBuf[linSt:linEnd+1])
						} else {
							// if the  line only contains a cr , print it (if lineEnd !> linSt -> lineEnd = linSt)
							txtStr = string(inBuf[linSt:linEnd+1])
							parEnd = true
						}
					default:

					}
				} else {
					// last char
					txtStr += string(inBuf[linSt:])
					parEnd = true
				}
			} else {
				if i < int(size) - 1 {
					txtStr = string(inBuf[linSt:linEnd+1])
				} else {
					txtStr = string(inBuf[linSt:])
				}
				parEnd = true
			}

			linSt = i+1
			ilin++
		}

		if parEnd {
				reqCount++
//	fmt.Printf("req %d: %s\n", reqCount, txtStr)
				insTxtReq := new(docs.InsertTextRequest)
				(*insTxtReq).EndOfSegmentLocation = &eos
				(*insTxtReq).Text = txtStr
				req := new(docs.Request)
				(*req).InsertText = insTxtReq
				updreqs = append(updreqs, req)
				txtStr = ""
				parEnd = false
		}
	}

	fmt.Printf("udpreqs: %d\n", len(updreqs))
	bUpdReq := new(docs.BatchUpdateDocumentRequest)
	bUpdReq.Requests = updreqs

	bUpdResp, err := svc.Documents.BatchUpdate(docId, bUpdReq).Do()
	if err != nil {
		fmt.Printf("BatchUpdate: %v", err)
		return fmt.Errorf("BatchUpdate: %v", err)
	}
	fmt.Printf("batch Update Response: %v\n", bUpdResp)
	return nil
}
