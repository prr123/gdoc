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
	hdMap map[string]string
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

type sectObj struct {
	Start int64
	End int64
	Rows int
	Cols int
	El int
}

type pbObj struct {
	Start int64
	End int64
	El int
	pEl int
}

type imgObj struct {
	Start int64
	Width int
	Height int
	Src string
	El int
}

type headObj struct {
	Typ string
	Start int64
	End int64
	El int
	Text string
}

func InitGdocEdit(docSvc *docs.DocumentsService, docId string) (gdEdObj *gdEditObj, err error) {
	var gdEd gdEditObj

	doc, err := docSvc.Get(docId).Do()
    if err != nil {
		return nil, fmt.Errorf("Unable to retrieve document: %v\n", err)
    }
	gdEd.DocSvc = docSvc
    gdEd.Doc = doc
	gdEd.DocId = docId
	gdEd.hdMap = map[string]string{
	"h1":"HEADING_1",
	"h2":"HEADING_2",
	"h3":"HEADING_3",
	"h4":"HEADING_4",
	"h5":"HEADING_5",
	"h6":"HEADING_6",
	}
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

func (edObj *gdEditObj) FindTextNext(text string, start int64) (stPos int64, err error) {

	if edObj.Doc == nil {return -1, fmt.Errorf("doc is nil!")}

//fmt.Printf("start: %d\n", start)
    body := edObj.Doc.Body
    numEl := len(body.Content)

	stPos = -1
    for el:=0; el< numEl; el++ {
        bodyEl := body.Content[el]
		if bodyEl.Paragraph == nil {continue}
        numPel := len(bodyEl.Paragraph.Elements)
		parStart := bodyEl.Paragraph.Elements[0].StartIndex
		parEnd := bodyEl.Paragraph.Elements[numPel-1].EndIndex
//fmt.Printf("par %d %d %d\n",el, parStart, parEnd)

		if start > parEnd {continue}
//			startEl = el

//		parStart := bodyEl.Paragraph.Elements[0].StartIndex
//		parEnd := bodyEl.Paragraph.Elements[numPel-1].EndIndex

		parStr:=""
        for pel:=0; pel<numPel; pel++ {
            parStr += bodyEl.Paragraph.Elements[pel].TextRun.Content
//          	pelStart := bodyEl.Paragraph.Elements[pel].StartIndex
//          pelEnd := bodyEl.Paragraph.Elements[pel].EndIndex
//fmt.Printf("Par %d El[%d]: %s\n",el, pel, parStr)

			idx := bytes.Index([]byte(parStr), []byte(text))
			if idx == -1 || start > parStart + int64(idx) { continue}
			stPos = parStart + int64(idx)
			break
		}
		if stPos > -1 {break}
    }

	//
	if stPos == -1 {return stPos, fmt.Errorf("string not found!")}
	return stPos, nil
}

func (edObj *gdEditObj) FindTextAll(text string) (posList *[]int64, err error) {

	var pos []int64

	if edObj.Doc == nil {return nil, fmt.Errorf("doc is nil!")}

    body := edObj.Doc.Body
    numEl := len(body.Content)

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
            parStr += bodyEl.Paragraph.Elements[pel].TextRun.Content
//          pelStart := bodyEl.Paragraph.Elements[pel].StartIndex
//          pelEnd := bodyEl.Paragraph.Elements[pel].EndIndex
//fmt.Printf("Par El[%d]: %s\n",pel, parElStr)
        }

		idx := bytes.Index([]byte(parStr), []byte(text))
		if idx > -1 {
			stPos := bodyEl.Paragraph.Elements[0].StartIndex + int64(idx)
			pos = append(pos, stPos)
		}
    }

	//
	if len(pos) <1 {return nil, fmt.Errorf("string not found!")}
	return &pos, nil
}

func (edObj *gdEditObj) GetEl(start int64) (elObj int, pelObj int, err error) {

	if edObj.Doc == nil {return -1, -1, fmt.Errorf("doc is nil!")}

//fmt.Printf("start: %d\n", start)
    body := edObj.Doc.Body
    numEl := len(body.Content)

	elObj = -1
	pelObj = -1

    for el:=0; el< numEl; el++ {
        bodyEl := body.Content[el]
		if bodyEl.Paragraph == nil {continue}
        numPel := len(bodyEl.Paragraph.Elements)
//		parStart := bodyEl.Paragraph.Elements[0].StartIndex
		parEnd := bodyEl.Paragraph.Elements[numPel-1].EndIndex

		if start > parEnd {continue}
		elObj = el

        for pel:=0; pel<numPel; pel++ {
          	pelStart := bodyEl.Paragraph.Elements[pel].StartIndex
			pelEnd := bodyEl.Paragraph.Elements[pel].EndIndex
//fmt.Printf("Par %d El[%d]: %s\n",el, pel, parStr)
			if start>= pelStart && start <= pelEnd {
				pelObj = pel
				break
			}
		}
		if elObj > 0 {break}
    }

	//
	if pelObj == -1 {return elObj, pelObj, fmt.Errorf("string not found!")}
	return elObj, pelObj, nil
}


func (edObj *gdEditObj) DispText(start int64) (err error) {

	if edObj.Doc == nil {fmt.Errorf("doc is nil!")}

    body := edObj.Doc.Body
    numEl := len(body.Content)

	elIdx := -1
    for el:=0; el< numEl; el++ {
        bodyEl := body.Content[el]
		if bodyEl.Paragraph == nil {continue}
        numPel := len(bodyEl.Paragraph.Elements)
		parStart := bodyEl.Paragraph.Elements[0].StartIndex
		parEnd := bodyEl.Paragraph.Elements[numPel-1].EndIndex

		if start >= parStart && start < parEnd {
			elIdx = el
			break
		}
	}

	if elIdx < 0 {return fmt.Errorf("paragraph El not found!")} 

	dispEl := body.Content[elIdx]
	numPel := len(dispEl.Paragraph.Elements)

	pelIdx :=-1
	for pel:=0; pel<numPel; pel++ {
		parEl := dispEl.Paragraph.Elements[pel]
		parElStart := parEl.StartIndex
		parElEnd := parEl.EndIndex
		if start >= parElStart && start < parElEnd {
			pelIdx = pel
            parElStr := parEl.TextRun.Content
			textByt := []byte(parElStr)
			relSt := start - parElStart
			relEnd := relSt + 10
			if start + relEnd > parElEnd { relEnd = parElEnd - parElStart}
			fmt.Printf("disp text: %s\n", string(textByt[relSt: relEnd]))
			return nil
		}
	}
	if pelIdx < 0 {return fmt.Errorf("text not found in pEl!")} 
	
	return nil
}

func (edObj *gdEditObj) DispPar(el int) (outstr string, err error) {

	if edObj.Doc == nil {fmt.Errorf("doc is nil!")}

    body := edObj.Doc.Body
    numEl := len(body.Content)

	if el > numEl {return "", fmt.Errorf("el > num els!")}

	bodyEl := body.Content[el]
	if bodyEl.Paragraph == nil {return "", fmt.Errorf("el not a paragraph!")}

	numPel := len(bodyEl.Paragraph.Elements)


	for i:=0; i< numPel; i++ {
		outstr += bodyEl.Paragraph.Elements[i].TextRun.Content
	}

	return outstr, nil
}

func (edObj *gdEditObj) ListTables() (tables *[]tblObj, err error) {
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

func (edObj *gdEditObj) AddTblRows(addRows int, tbl *tblObj) (updreq *docs.BatchUpdateDocumentRequest, err error) {
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

func  (edObj *gdEditObj) GetTblContent(tbl *tblObj) (contTbl *[][]string, err error) {

fmt.Printf("Get Tbk Cont: %d %d\n", (*tbl).Rows, (*tbl).Cols)

	table := make([][]string,(*tbl).Rows)

	for irow:=0; irow< (*tbl).Rows; irow++ {
		tblRow := make([]string, (*tbl).Cols)
		table[irow] = tblRow
	}

	doc := edObj.Doc
	if doc == nil {return nil, fmt.Errorf("doc not provided")}

	el := doc.Body.Content[tbl.El]

	if el.Table == nil {return nil, fmt.Errorf("el %d is not a table!", tbl.El)}

	elTbl := el.Table
	if len(table) != int(elTbl.Rows) {return nil, fmt.Errorf("contTbl and tbl row numbers do not match!")}
	if len(table[0]) != int(elTbl.Columns) {return nil, fmt.Errorf("contTbl and tbl col numbers do not match!")}


	for row:=0; row<int(elTbl.Rows); row++ {
		for col:=0; col<int(elTbl.Columns); col++ {

			tblCel := elTbl.TableRows[row].TableCells[col]

    		celElCount := len(tblCel.Content)
//fmt.Printf("cel[%d,%d]: %d ", row, col, celElCount)
			celStr :=""
			for el:=0; el< celElCount; el++ {
        		celEl := tblCel.Content[el]
        		if celEl.Paragraph == nil {continue}
				numPel := len(celEl.Paragraph.Elements)
				// celPel cell Paragraph Element
				for celPel:=0; celPel< numPel; celPel++ {
					celStr += celEl.Paragraph.Elements[celPel].TextRun.Content
				}
			} //for el
//fmt.Printf("%s\n", celStr)
			celB := []byte(celStr)
//for i:=0; i<len(celB); i++ {fmt.Printf("%q ",celB[i])}
//fmt.Println()
			if celB[len(celB) - 1] == '\n' {celStr = string(celB[:len(celB)-1])}
			table[row][col] = celStr
		} // for col
	} // for row
	return &table, nil
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

//    fmt.Printf("update Cell Content: table: %d %d Index: %d\n", elTbl.Rows, elTbl.Columns, el.StartIndex)

	updreq = new(docs.BatchUpdateDocumentRequest)
    updreq.Requests = make([]*docs.Request, el.Table.Rows*el.Table.Columns)

	reqCount:=0

	delta := int64(0)
	for row:=0; row<int(elTbl.Rows); row++ {
		for col:=0; col<int(elTbl.Columns); col++ {

			tblCel := elTbl.TableRows[row].TableCells[col]
    		celElCount := len(tblCel.Content)
			if celElCount > 1 {return nil, fmt.Errorf("table cell[%d, %d] not empty", row, col)}
//			for el:=0; el< celElCount; el++ {
			celEl := tblCel.Content[0]
			celPar := celEl.Paragraph
			if celPar == nil {return nil, fmt.Errorf("table cell[%d, %d] has no Paragraph", row, col)}

			parEl = celPar.Elements[0]

    		loc:= new(docs.Location)
			loc.Index = parEl.StartIndex + delta

			newTxt := (*contTbl)[row][col]

    		insTxtReq:= new(docs.InsertTextRequest)
    		insTxtReq.Location = loc
    		insTxtReq.Text = newTxt

    		insReq := new(docs.Request)

    		insReq.InsertText = insTxtReq
			updreq.Requests[reqCount] = insReq

			delta += int64(len(newTxt))
			reqCount++
		}
	}

    return updreq, nil
}

func  (edObj *gdEditObj) ClearTblContent(tbl *tblObj) (updreq *docs.BatchUpdateDocumentRequest, err error) {

    if tbl == nil {return nil, fmt.Errorf("no tblObj provided!")}

	doc := edObj.Doc

	el := doc.Body.Content[tbl.El]

	if el.Table == nil {return nil, fmt.Errorf("el %d is not a table!", tbl.El)}

	elTbl := el.Table

    fmt.Printf("clear Cell Content: table: %d %d Index: %d\n", elTbl.Rows, elTbl.Columns, el.StartIndex)

	updreq = new(docs.BatchUpdateDocumentRequest)
    updreq.Requests = make([]*docs.Request, 0, elTbl.Rows*elTbl.Columns)

//fmt.Printf("reqs: %d\n", len(updreq.Requests))

	reqCount:=0

	delta := int64(0)
	for row:=0; row<int(elTbl.Rows); row++ {
		for col:=0; col<int(elTbl.Columns); col++ {

			tblCel := elTbl.TableRows[row].TableCells[col]
    		celElCount := len(tblCel.Content)
			if celElCount > 1 {return nil, fmt.Errorf("table cell[%d, %d] not empty", row, col)}
//			for el:=0; el< celElCount; el++ {
			celEl := tblCel.Content[0]
			celPar := celEl.Paragraph
			if celPar == nil {return nil, fmt.Errorf("table cell[%d, %d] has no Paragraph", row, col)}

//			parEl := celPar.Elements[0]
			numPel := len(celPar.Elements)

				// celPel cell Paragraph Element
			celStr:=""
			for celPel:=0; celPel< numPel; celPel++ {
				celStr += celPar.Elements[celPel].TextRun.Content
			}

			stPos := celPar.Elements[0].StartIndex + delta
			endPos := celPar.Elements[numPel -1].EndIndex -1 + delta

			delta += stPos - endPos
			if (endPos - stPos) < 1 {continue}

			delRange:= new(docs.Range)
    		delRange.StartIndex =  stPos
   			delRange.EndIndex = endPos

			delTxtReq:= new(docs.DeleteContentRangeRequest)
    		delTxtReq.Range = delRange

    		delReq := new(docs.Request)
    		delReq.DeleteContentRange = delTxtReq

			updreq.Requests = append(updreq.Requests, delReq)
//    		updreq.Requests[reqCount] = delReq

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

func (edObj *gdEditObj) ListImg() (images *[]imgObj, err error) {
// method that finds all tables in google doc. returns an array of startindices

	var img imgObj
	var imgList []imgObj

	doc := edObj.Doc

    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.Table != nil {
            img.Start = el.StartIndex
//			img.Rows = int(el.Table.Rows)
//			img.Cols = int(el.Table.Columns)
			img.El = i
			imgList = append(imgList, img)
        }
    }
	return &imgList, nil
}

func (edObj *gdEditObj) ListSects() (sections *[]sectObj, err error) {
// method that finds all tables in google doc. returns an array of startindices

	var sect sectObj
	var sectList []sectObj

	doc := edObj.Doc

    for i:=0; i< len(doc.Body.Content); i++ {
        el := doc.Body.Content[i]
        if el.SectionBreak != nil {
            sect.Start = el.StartIndex
			sect.End = el.EndIndex
			sect.El = i
			sectList = append(sectList, sect)
        }
    }
	return &sectList, nil
}

func (edObj *gdEditObj) GetHeadTxt(hd *headObj) (outstr string, err error) {

	if edObj.Doc == nil {return "", fmt.Errorf("doc is nil!")}

    body := edObj.Doc.Body
    numEl := len(body.Content)
	if hd.El > numEl -1 {return "", fmt.Errorf("invalid header Element!")}

	el := body.Content[hd.El]
    if el.Paragraph == nil {return "", fmt.Errorf("el is not a Paragraph!")}

	for pel:=0; pel < len(el.Paragraph.Elements); pel++ {
		if el.Paragraph.Elements[pel].TextRun == nil {continue}
		outstr += el.Paragraph.Elements[pel].TextRun.Content
	}
	return outstr, nil
}

func (edObj *gdEditObj) ListHeadings(shortHead string) (headlist *[]headObj, err error) {

	var head headObj
	var hdList []headObj

	if edObj.Doc == nil {return nil, fmt.Errorf("doc is nil!")}

	heading, ok := edObj.hdMap[shortHead]
	if !ok {return nil, fmt.Errorf("invalid heading!")}

fmt.Printf("heading: %s %s\n", heading, shortHead)

    body := edObj.Doc.Body
    numEl := len(body.Content)

    for iEl:=0; iEl< numEl; iEl++ {
        bodyEl := body.Content[iEl]
        if bodyEl.Paragraph == nil {continue;}
		parstyl := bodyEl.Paragraph.ParagraphStyle

//fmt.Printf("parstyle [el: %d] %s\n", iEl, parstyl.NamedStyleType)
		// found paragraph
		outstr := ""
		if parstyl.NamedStyleType == heading {
			head.Typ = shortHead
			head.Start = bodyEl.StartIndex
			head.End = bodyEl.EndIndex
			for pel:=0; pel < len(bodyEl.Paragraph.Elements); pel++ {
				if bodyEl.Paragraph.Elements[pel].TextRun == nil {continue}
				outstr += bodyEl.Paragraph.Elements[pel].TextRun.Content
			}
			outByt := []byte(outstr)
			if outByt[len(outByt)-1] == '\n' {outstr = string(outByt[:len(outByt)-1])}
			head.Text = outstr
			hdList = append(hdList, head)
		}
    }

	return &hdList, nil
}


func (edObj *gdEditObj) ListHeadingsAlt(heading string) (ellist *[]int, err error) {

	var elList []int
	var tgtStyl string

	if edObj.Doc == nil {return nil, fmt.Errorf("doc is nil!")}

	switch heading {
	case "h1":
		tgtStyl = "HEADING_1"
	case "h2":
		tgtStyl = "HEADING_2"
	default:
		return nil, fmt.Errorf("not a valid heading supplied")
	}

fmt.Printf("style: %s\n", tgtStyl)

    body := edObj.Doc.Body
    numEl := len(body.Content)

    for el:=0; el< numEl; el++ {
        bodyEl := body.Content[el]
        if bodyEl.Paragraph == nil {continue;}

		parstyl := bodyEl.Paragraph.ParagraphStyle
		// found paragraph
		if parstyl.NamedStyleType == tgtStyl {
			elList = append(elList, el)
		}
    }

	return &elList, nil
}


func PrintTblObj (tbls *[]tblObj) {

	fmt.Printf("*** Tables: %d ****\n", len(*tbls))
	for i:=0; i< len(*tbls); i++ {
		fmt.Printf("table[%d]: rows: %d cols: %d\n", i, (*tbls)[i].Rows, (*tbls)[i].Cols)
	}
	fmt.Printf("********************\n")
}

func PrintUpdResp (resp *docs.BatchUpdateDocumentResponse) {

    fmt.Println("******* Batch Update Response ********")
    fmt.Println("document id: ", resp.DocumentId)

    wc := resp.WriteControl
    if wc == nil {return}
    fmt.Println("rev id: ", wc.RequiredRevisionId)


    fmt.Println("Replies: ", len(resp.Replies))

    for i:=0; i < len(resp.Replies); i++ {
        rpl := resp.Replies[i]
        fmt.Printf("reply [%d]:\n", i)
        if rpl.CreateFooter != nil {
            fmt.Println("  create footer id: ", rpl.CreateFooter.FooterId)
        }
        if rpl.ReplaceAllText != nil {
            fmt.Println("  replace text: ", rpl.ReplaceAllText.OccurrencesChanged)
        }

    }

}

func PrintTbl(contTbl *[][]string) {

    rows:= len(*contTbl)
    cols := len((*contTbl)[0])
fmt.Printf("rows: %d cols: %d\n", rows, cols)

    fmt.Printf("cols:    ")
    for ic:=0; ic<cols; ic++ {
        fmt.Printf(" %8d   |", ic+1)
    }
    fmt.Println()

    for ir:=0; ir<rows; ir++ {
        fmt.Printf("row[%d]: |", ir+1)
        for ic:=0; ic<cols; ic++ {
            fmt.Printf(" %10s |", (*contTbl)[ir][ic])
        }
        fmt.Println()
    }


}
