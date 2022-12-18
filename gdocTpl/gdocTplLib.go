// golang library that parses a gdoc document template
// adopted from gdocEditLib.go
// author: prr, azul software
// created: 18/12/2022
// copyright 2022 prr, Peter Riemenschneider, Azul Software
//
// for changes see github
//
// start: GdocTplAnalyse
// GdocTplFill
//


package gdocTpl

import (
	"fmt"
//	"os"
//	"unicode/utf8"
	"bytes"
//	"strings"
	"google.golang.org/api/docs/v1"
//	gdoc "google/gdoc/gdocCommon"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdocObj struct {
	doc *docs.Document
	parCount int
	DocName string
	TplItemList *[]TplItem
	TplNum int
}

type TplItem struct {
	Par *docs.Paragraph
	ParEl *[]docs.ParagraphElement
	TplName string
	Start int64
	End int64
}


func InitTpl(doc *docs.Document) (gdobj *gdocObj) {

var gdObj gdocObj

	gdObj.doc = doc
	gdObj.parCount = 0
	return &gdObj
}


func (dObj *gdocObj) ParseDoc(tplFilNam string) (err error) {

	body := dObj.doc.Body
	numEl := len(body.Content)
fmt.Printf("************ Body Els: %d *************\n", numEl)

	parState := 0
	for el:=0; el< numEl; el++ {
		bodyEl := body.Content[el]
fmt.Printf("el [%d]: ", el)
		if bodyEl.Paragraph == nil {fmt.Printf(" no par\n"); continue;}

		numPel := len(bodyEl.Paragraph.Elements)
fmt.Printf("Par Elements: %d\n", numPel)


		for pel:=0; pel<numPel; pel++ {
//			parEl := bodyEl.Paragraph.Elements[pel]
			parElStr := bodyEl.Paragraph.Elements[pel].TextRun.Content
//fmt.Printf("Par El[%d]: %s\n",pel, parElStr)

			byteEl := []byte(parElStr)
			stdx := bytes.Index(byteEl,[]byte("{"))
			if stdx == -1 {continue}
			parState = 1
			edx := bytes.Index(byteEl,[]byte("}"))
			if edx == -1 {
				if parState != 1 {return fmt.Errorf("found end bracket without start bracket!")}
			} else {
				parState = 0
fmt.Printf("tpl field: %s\n", string(byteEl[stdx+1:edx]))
			}

			if parState == 1 {
				parState = 0
				return fmt.Errorf("no ending bracket found!")
			}
		}
	}

	return nil
}
