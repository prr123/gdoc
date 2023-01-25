// gdocEditLib.go
// golang library that makes changes to a gdoc document
// author: prr, azul software
// created: 23/2/2023
// copyright 2023 prr, Peter Riemenschneider, Azul Software
//
// license see github
//
//

package gdocEditLib

import (
	"fmt"
//	"os"
//	"unicode/utf8"
	"bytes"
	"google.golang.org/api/docs/v1"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdEditObj struct {
	DocSvc *docs.DocumentsService
	Doc *docs.Document
	DocId string
	diffStart int64
	diff int64
}

type stringObj struct {
	start int64
	txt string
	elIdx int64
	parEl int64
	newTxt string
}

type tblObj struct {
	Start int64
	Rows int
	Cols int
	El int
}

func initGdocEdit(docSvc *docs.DocumentsService, docId string) (gdEdObj *gdEditObj, err error) {
	var gdEd gdEditObj

	doc, err := docSvc.Get(docId).Do()
    if err != nil {
		return nil, fmt.Errorf("Unable to retrieve document: %v\n", err)
    }
	gdEd.DocSvc = docSvc
    gdEd.Doc = doc
	gdEd.DocId = docId
	return &gdEd, nil
}

func (edObj *gdEditObj) BatchUpd(docId string, updreq *docs.BatchUpdateDocumentRequest) (err error) {

	if len(docId) < 1 {return fmt.Errorf("no valid docId")}
	if len(updreq.Requests) < 1 {return fmt.Errorf("no valid update requests!")}

    _, err = edObj.DocSvc.BatchUpdate(docId, updreq).Do()
    if err != nil {
        return fmt.Errorf("cannot update document: %v\n", err)
   	}
	return nil
}

func (edObj *gdEditObj) FindText(strObj *stringObj) (stPos int64, err error) {

	if edObj.Doc == nil {return -1, fmt.Errorf("doc is nil!")}

    body := edObj.Doc.Body

    numEl := len(body.Content)
//fmt.Printf("************ Body Els: %d *************\n", numEl)

	stPos = -1
    for el:=0; el< numEl; el++ {
        bodyEl := body.Content[el]
        if bodyEl.Paragraph == nil {continue;}
		// found paragraph
        numPel := len(bodyEl.Paragraph.Elements)
//fmt.Printf("Par Elements: %d\n", numPel)

//		parStart := bodyEl.Paragraph.Elements[0].StartIndex
//		parEnd := bodyEl.Paragraph.Elements[numPel-1].EndIndex

		parStr:=""
        for pel:=0; pel<numPel; pel++ {
//          parEl := bodyEl.Paragraph.Elements[pel]
            parStr += bodyEl.Paragraph.Elements[pel].TextRun.Content
//          pelStart := bodyEl.Paragraph.Elements[pel].StartIndex
//          pelEnd := bodyEl.Paragraph.Elements[pel].EndIndex
//fmt.Printf("Par El[%d]: %s\n",pel, parElStr)

        }

		idx := bytes.Index([]byte(parStr), []byte(strObj.txt))
		if idx > -1 {
			stPos = bodyEl.Paragraph.Elements[0].StartIndex + int64(idx)
			break
		}
    }

	//
	if stPos == -1 {return stPos, fmt.Errorf("string not found!")}
	return stPos, nil
}

func (edObj *gdEditObj) FindTables() (tables *[]tblObj, err error) {
// method that finds all tables in google doc. returns an array of startindices

	var tbl tblObj
	var tableList []tblObj

	doc := edObj.Doc

    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Table != nil {
            tbl.Start = el.StartIndex
			tbl.Rows = int(el.Table.Rows)
			tbl.Cols = int(el.Table.Columns)
			tbl.El = i
			tableList = append(tableList, tbl)
        }
    }
	return &tableList, nil
}

func (edObj *gdEditObj) FindNextTable(start int64) (tblobj *tblObj, err error) {
// method to find the next table after start

	var tbl tblObj

    tbl.Start = -1
	doc := edObj.Doc

    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Table != nil {
			if el.StartIndex >= start {
            	tbl.Start = el.StartIndex
				tbl.Rows = int(el.Table.Rows)
				tbl.Cols = int(el.Table.Columns)
				tbl.El = i
            	break
			}
        }
    }

	if tbl.Start == -1 {return nil, fmt.Errorf("table not found!")}
	return &tbl, nil
}

func (edObj *gdEditObj) addTblRows(addRows int, tbl *tblObj) (updreq *docs.BatchUpdateDocumentRequest, err error) {
// method that adds rows to an exisiting table

    var loc docs.Location
    var tblCelLoc docs.TableCellLocation
    var insTblRowReq docs.InsertTableRowRequest
    var req docs.Request

//  tbl = el.Table
//  fmt.Printf("addRows: table: %d %d Index: %d\n", tbl.Rows, tbl.Columns, tblStart)


    if tbl == nil {return nil, fmt.Errorf("no tblObj provided!")}

	doc := edObj.Doc

	el := doc.Body.Content[tbl.El]

	if el.Table == nil {return nil, fmt.Errorf("el %d is not a table!", tbl.El)}

    loc.Index = el.StartIndex

    tblCelLoc.TableStartLocation = &loc
    tblCelLoc.ColumnIndex = 0

    insTblRowReq.InsertBelow = true

    updreq = new(docs.BatchUpdateDocumentRequest)

    updreq.Requests = make([]*docs.Request, addRows)

    for i:= 0; i< addRows; i++ {
        celLoc := tblCelLoc
        celLoc.RowIndex = int64(i+1)

        addRowReq :=  insTblRowReq
        addRowReq.TableCellLocation = &celLoc
        insReq := req
        insReq.InsertTableRow = &addRowReq
        updreq.Requests[i] = &insReq
    }

	return updreq, nil
}

func  (edObj *gdEditObj) FillTblContent(contTbl *[][]string, tbl *tblObj) (updreq *docs.BatchUpdateDocumentRequest, err error) {

    var parEl *docs.ParagraphElement

    if tbl == nil {return nil, fmt.Errorf("no tblObj provided!")}

	doc := edObj.Doc

	el := doc.Body.Content[tbl.El]

	if el.Table == nil {return nil, fmt.Errorf("el %d is not a table!", tbl.El)}

	elTbl := el.Table
	if len(*contTbl) != int(elTbl.Rows) {return nil, fmt.Errorf("contTbl and tbl row numbers do not match!")}
	if len((*contTbl)[0]) != int(elTbl.Columns) {return nil, fmt.Errorf("contTbl and tbl col numbers do not match!")}

    fmt.Printf("update Cell Content: table: %d %d Index: %d\n", elTbl.Rows, elTbl.Columns, el.StartIndex)

	updreq = new(docs.BatchUpdateDocumentRequest)
    updreq.Requests = make([]*docs.Request, el.Table.Rows*el.Table.Columns)

	reqCount:=0
	for row:=0; row<int(elTbl.Rows); row++ {

		for col:=0; col<int(elTbl.Columns); col++ {

			tblCel := elTbl.TableRows[row].TableCells[col]

    		celContItems := len(tblCel.Content)

			idx := -1
			for i:=0; i< celContItems; i++ {
        		celCont := tblCel.Content[i]
        		if celCont.Paragraph != nil {
            		idx = i
            		break
        		}
			}

//fmt.Printf("tblCel[%d, %d]: %d %d\n", row, col, celContItems, idx)


//			if idx < 0 {
        			// insert paragraph

 //   		} else {
        		celPar := tblCel.Content[idx].Paragraph
//        			parEls := len(celPar.Elements)
//fmt.Printf("celpars - idx: %d parEls: %d\n", idx, parEls)
//        for j:=0; j< parEls; j++ {
            	parEl = celPar.Elements[0]
//            	elStr = parEl.TextRun.Content
//fmt.Printf("parEl[%d]: %d \"%s\"\n", parEl.StartIndex, len(elStr), elStr)
//        }
//        			parEl0 = celPar.Elements[0]
//    		}

    		loc:= new(docs.Location)
			loc.Index = parEl.StartIndex

    		insTxtReq:= new(docs.InsertTextRequest)
    		insTxtReq.Location = loc

    		insTxtReq.Text = (*contTbl)[row][col]

    		insReq := new(docs.Request)

    		insReq.InsertText = insTxtReq
//      insReq.InsertTableRow = &addRowReq
			updreq.Requests[reqCount] = insReq
			reqCount++
		}
	}

    return updreq, nil
}



func (edObj *gdEditObj) FindString(strObj stringObj) (stPos int64, err error) {

	// first we need to find the text string
	start := strObj.start
	doc := edObj.Doc
	par := -1
    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Paragraph != nil {
			if start >=el.StartIndex && start < el.EndIndex {
            	par = i
            	break
			}
        }
    }
	if par == -1 {return -1, fmt.Errorf("text not found in doc")}

	stPos = -1
    for i:=par; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Paragraph == nil {continue}

		pEl := doc.Body.Content[i].Paragraph
		parStr := ""
		for j:=0; j<len(pEl.Elements); j++ {
			parel := pEl.Elements[i]
			parStr += parel.TextRun.Content
		}

		idx:= bytes.Index([]byte(parStr),[]byte(strObj.txt))
		if idx == -1 {continue}
		stPos = el.StartIndex + int64(idx)
	}

	if stPos == -1 {return stPos, fmt.Errorf("string not found!")}
	return stPos, nil
}

func (edObj *gdEditObj) ReplaceString(strObj stringObj) (updreq *docs.BatchUpdateDocumentRequest, err error) {

	var insTxtReq docs.InsertTextRequest
	var delTxtReq docs.DeleteContentRangeRequest
	var loc docs.Location
	var delRange docs.Range

	var stPos int64

	// first we need to find the text string
	start := strObj.start
	doc := edObj.Doc
	par := -1
    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Paragraph != nil {
			if start >=el.StartIndex && start < el.EndIndex {
            	par = i
            	break
			}
        }
    }
	if par == -1 {return nil, fmt.Errorf("text not found in doc")}

	stPos = -1
    for i:=par; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Paragraph == nil {continue}

		pEl := doc.Body.Content[i].Paragraph
		parStr := ""
		for j:=0; j<len(pEl.Elements); j++ {
			parel := pEl.Elements[i]
			parStr += parel.TextRun.Content
		}

		idx:= bytes.Index([]byte(parStr),[]byte(strObj.txt))
		if idx == -1 {continue}
		stPos = el.StartIndex + int64(idx)
	}

	if stPos == -1 {return nil, fmt.Errorf("string not found!")}

    loc.Index = int64(stPos)

    insTxtReq.Location = &loc
    insTxtReq.Text = strObj.newTxt

    updreq = new(docs.BatchUpdateDocumentRequest)
    updreq.Requests = make([]*docs.Request, 2)

    updreq.Requests[0].InsertText = &insTxtReq


    delst := stPos + int64(len(strObj.newTxt))
    delend := delst + int64(len(strObj.txt))
    delRange.StartIndex = int64(delst)
    delRange.EndIndex = int64(delend)
    delTxtReq.Range = &delRange

    updreq.Requests[1].DeleteContentRange = &delTxtReq

    return updreq, nil
}

func PrintTblObj (tbls *[]tblObj) {

	fmt.Printf("*** Tables: %d ****\n", len(*tbls))
	for i:=0; i< len(*tbls); i++ {
		fmt.Printf("table[%d]: rows: %d cols: %d\n", i, (*tbls)[i].Rows, (*tbls)[i].Cols)
	}
	fmt.Printf("********************\n")
}
