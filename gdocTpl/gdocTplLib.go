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
	"os"
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


type gdocTpl struct {
	doc *docs.Document
	tplFil *os.File
	tplFilNam string
	parCount int
	TplItemList *[]TplItem
	TplNum int
}

type TplItem struct {
	name string
	typ string
	par int
	parEl int
	Start int64
	End int64
}


func InitTpl(doc *docs.Document) (gdobj *gdocTpl) {

var tpl gdocTpl

	tpl.doc = doc
	tpl.parCount = 0
	return &tpl
}



func (tplObj *gdocTpl) creTplFil(tplFilNam string) (err error) {

	tplfil, err :=os.Create(tplFilNam)
	if err != nil {return fmt.Errorf("os.Create: %v", err)}
	tplObj.tplFil = tplfil

	return nil
}

func (tplObj *gdocTpl) ParseDoc(tplFilNam string) (err error) {

	var tplItem TplItem

 	tplList := make([]TplItem,0)

	body := tplObj.doc.Body
	numEl := len(body.Content)
//fmt.Printf("************ Body Els: %d *************\n", numEl)

	parState := 0
	for el:=0; el< numEl; el++ {
		bodyEl := body.Content[el]
//fmt.Printf("el [%d]: ", el)
		if bodyEl.Paragraph == nil {continue;}

//fmt.Printf(" no par\n"); 

		numPel := len(bodyEl.Paragraph.Elements)
//fmt.Printf("Par Elements: %d\n", numPel)


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
				tplItem.name = string(byteEl[stdx+1:edx])
				tplItem.par = el
				tplItem.parEl = pel
				tplList = append(tplList, tplItem)
//fmt.Printf("tpl field: %s\n", string(byteEl[stdx+1:edx]))
			}

			if parState == 1 {
				parState = 0
				return fmt.Errorf("no ending bracket found!")
			}
		}
	}

	if tplObj.tplFil != nil {tplObj.tplFil.Close()}

	tplObj.TplItemList = &tplList

	return nil
}

func (tplObj *gdocTpl) CreateTplFil(tplFilNam string) (err error) {

	tplfil, err :=os.Create(tplFilNam)
	if err != nil {return fmt.Errorf("os.Create: %v", err)}

	doc := tplObj.doc
	outstr := "Title: " + doc.Title + "\n"
	tplfil.WriteString(outstr)
	outstr = "Id: " + doc.DocumentId + "\n"
	tplfil.WriteString(outstr)

	tplList := (*tplObj.TplItemList)

	outstr = fmt.Sprintf("NamesLen: %d\n", len(tplList))
	tplfil.WriteString(outstr)

	errCount :=0
	for i:=0; i<len(tplList); i++ {
		tplItem := tplList[i]
//		outstr:= fmt.Sprintf("item: %2d\n   - name: %-10s\n   - par: %d\n   - pel: %d\n", i, tplItem.name, tplItem.par, tplItem.parEl)
		errProc := procTplItem(&tplItem)
		if errProc != nil {
			errCount++
			errStr := fmt.Sprintf("%v\n", errProc)
			outstr = fmt.Sprintf("%s:  # %s\n", tplItem.name, errStr)
		} else {
			if len(tplItem.typ) == 0 {
				outstr = fmt.Sprintf("%s:\n", tplItem.name)
			} else {
				outstr = fmt.Sprintf("%s:  #<%s>\n", tplItem.name, tplItem.typ)
			}
		}
		tplList[i] = tplItem
		tplfil.WriteString(outstr)
//fmt.Printf("item [%d]: %v\n", i, tplItem)
	}
	tplfil.Close()
	if errCount > 0 {return fmt.Errorf("process name errors %d", errCount)}
	*tplObj.TplItemList = tplList

	return nil
}

func procTplItem(item *TplItem) (err error) {

	namByt := []byte((*item).name)
//fmt.Println("name: ", string(namByt))

	typSt := -1
	typEnd:= -1
	for i:=0; i< len(namByt); i++ {
		switch namByt[i] {
		case '<':
			typSt = i
		case '>':
			typEnd = i
		}
	} //i

	if typSt < 0 {return nil}

	if typEnd == -1 {return fmt.Errorf("%s: no ending bracket", string(namByt))}
	if typEnd<=typSt {return fmt.Errorf("%s: ending bracket before opening bracket", string(namByt))}

	(*item).name = string(namByt[:typSt])
	(*item).typ = string(namByt[typSt+1: typEnd])
	return nil
}

func (tplObj *gdocTpl) PrintTpl() () {

	doc := tplObj.doc
	fmt.Printf("*********** Tpl Title: %s ***********\n", doc.Title)
	fmt.Printf("Document Id:   %s\n", doc.DocumentId)
	fmt.Printf("Template File: %s\n", tplObj.tplFilNam)

	tplList := (*tplObj.TplItemList)
	for i:=0; i<len(tplList); i++ {
		tplItem := tplList[i]
		fmt.Printf("item: %2d name: %-10s typ: %s par: %d pel: %d\n", i, tplItem.name, tplItem.typ, tplItem.par, tplItem.parEl)
	}
	return
}

