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
	"strings"
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


func (dObj *gdocObj) getColor(color  *docs.Color)(outstr string) {
	outstr = ""
	if color != nil {
        blue := int(color.RgbColor.Blue*255.0)
        red := int(color.RgbColor.Red*255.0)
        green := int(color.RgbColor.Green*255)
        outstr += fmt.Sprintf("rgb(%d, %d, %d);\n", red, green, blue)
		return outstr
	}
	outstr = "no color\n"
	return outstr
}


func (dObj *gdocObj) dispBodyEl(elStr *docs.StructuralElement)(outstr string, err error) {

	if elStr == nil {
		return "", fmt.Errorf("error dispContentEl -- elStr is nil")
	}
	if dObj == nil {
		return "", fmt.Errorf("error dispContentEl -- dObj is nil")
	}

//	doc := dObj.Doc
//	body := doc.Body
	notFound := true
	if elStr.Paragraph != nil {
		outstr += fmt.Sprintf(" Paragraph StartIndex: %d EndIndex: %d\n",  elStr.StartIndex, elStr.EndIndex)
		notFound = false
	}
	if elStr.SectionBreak != nil {
		outstr += fmt.Sprintf(" Section_Break StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
	}
	if elStr.Table != nil {
		outstr += fmt.Sprintf(" Table StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
	}
	if elStr.TableOfContents != nil {
		outstr += fmt.Sprintf(" TOC StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
	}
	if notFound {
		return outstr, fmt.Errorf(" type: unknown StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
	}
	return outstr, nil
}

func (dObj *gdocObj) SearchTemp() (outstr string, err error) {

	body := dObj.doc.Body
	numEl := len(body.Content)
	outstr = "\n**********************************\n"
	outstr += fmt.Sprintf("[Body] %d\n", numEl)

	for el:=0; el< numEl; el++ {
		bodyEl := body.Content[el]
		outstr += fmt.Sprintf("\n[element] %d ", el)
		if bodyEl.Paragraph == nil {
			outstr +="\n"
			continue
		}
		numPel := len(bodyEl.Paragraph.Elements)
		outstr += fmt.Sprintf("  [Par Elements] %d\n", numPel)
		for pel:=0; pel<numPel; pel++ {
			parEl := bodyEl.Paragraph.Elements[pel]
			parElStr := bodyEl.Paragraph.Elements[pel].TextRun.Content
			outstr += fmt.Sprintf("Par El[%d] len: %d Start: %d End: %d\n",pel, len(parElStr),parEl.StartIndex, parEl.EndIndex)
			outstr += parElStr +"\n"
			idx := strings.Index(parElStr,"<<")
//			outstr += fmt.Sprintf("index: %d\n", idx)
			if idx > -1 {
				edx := strings.Index(parElStr,">>")
				if edx > -1 {
					outstr += fmt.Sprintf("found [%s] %d:%d\n ", parElStr[idx+2:edx], idx+2, edx)
				} else {
					outstr += fmt.Sprintf("error no \">>\"!\n")
				}
			}
		}
	}

//	outfil.WriteString(outstr)
	return outstr, nil
}


/*
func CvtGdocToTxt(outfil *os.File, doc *docs.Document)(err error) {
	var outstr string

    docObj := new(gdocObj)
    docObj.doc = doc
    err = docObj.Init()
    if err != nil {
        return fmt.Errorf("error Cvt Txt Init %v", err)
    }
	_, err = outfil.WriteString("[Document_Title] " + doc.Title + "\n")
	if err != nil {
		return fmt.Errorf("error CvtDdoc -- cannot write to file: %v", err)
	}

	outstr = fmt.Sprintf("[Document_Id] %s \n", doc.DocumentId)
	outstr += fmt.Sprintf("[Revision_Id] %s \n", doc.RevisionId)

// Inline Objects
// Lists
	inObjLen := len(doc.InlineObjects)
	outstr += fmt.Sprintf("\n[Inline_Objects] %d\n", inObjLen)

	posObjLen := len(doc.PositionedObjects)
	outstr += fmt.Sprintf("\n[Positioned_Objects] %d\n",posObjLen)

	headLen := len(doc.Headers)
	outstr += fmt.Sprintf("\n[Headers] %d\n",headLen)
	knum := 0
	for key, header := range doc.Headers {
		knum++
		outstr += fmt.Sprintf("  header %d: %s || %s\n", knum, key,header.HeaderId )
	}

	footLen := len(doc.Footers)
	outstr += fmt.Sprintf("\n[Footers] %d\n",footLen)

	ftnoteLen := len(doc.Footnotes)
	outstr += fmt.Sprintf("\n[Footnotes] %d\n",ftnoteLen)

	listLen := len(doc.Lists)
	outstr += fmt.Sprintf("\n[Lists] %d\n",listLen)
	knum = 0
	for key, list := range doc.Lists {
		knum++
		nest := list.ListProperties.NestingLevels
		outstr += fmt.Sprintf("\nList %d: id: %s nest levels: %d\n", knum, key, len(nest) )
	}
//	outfil.WriteString(outstr)

	nrLen := len(doc.NamedRanges)
	outstr += fmt.Sprintf("\n[NamedRanges] %d\n",nrLen)
	outfil.WriteString(outstr)


	body := doc.Body
	numEl := len(body.Content)
	outstr = fmt.Sprintf("\n[Body] %d\n", numEl)

	for el:=0; el< numEl; el++ {
		elStr := body.Content[el]
		outstr += fmt.Sprintf("\n [element] %d ", el)
		tStr,err := docObj.dispBodyEl(elStr)
		if err != nil {
			outstr += fmt.Sprintf("/n error dispContent[%d]: %v\n", el, err) + tStr
		} else {
			outstr += tStr
		}
	}
	outfil.WriteString(outstr)

	tStr, err := docObj.SearchTemp()
	if err != nil {
		outstr = fmt.Sprintf("/n error SearchTemp: %v\n", err) + tStr
	} else {
		outstr = tStr
	}
	outfil.WriteString(outstr)

	outfil.Close()
	return nil
}
/*
func ReadGdocTpl(infilnam string)(tplList *Tpl_list, err error) {
	var tpl_list Tpl_list
	var inBt [1024]byte

	infil, err := os.Open(infilnam)
	if err != nil {
		return fmt.Errorf("error ReadGdocTpl - cannot open file: %s err: %v",infilnam, err)
	}

	nb, err := infil.Read(inBt[:])
	if err != nil {
        return fmt.Errorf("error ReadGdocTpl read %v", err)
	}
	fmt.Printf("read %d bytes!\n", nb)

	ist:=0
	il:= 0
	tpls := make([]tpl, 10, 20)
	tpl_list.tpls = &tpls

	for i:=0; i<nb; i++ {
		if inBt[i] == '\n' {
			tpls[il].tpl_line = string(inBt[ist:i])
			ist = i+1
			il++
		}
	}
	numLin := il
	for il:=0; il< numLin; il++ {
		fmt.Printf("line: %d tpl: %s\n", il, tpls[il].tpl_line)
	}

	tplList = &tpl_list
	return tplList, nil
}
*/

// create text file to dump document file
func CreTxtOutFile(filnam string, ext string)(outfil *os.File, err error) {

    // convert white spaces in file name to underscore
    bfil := []byte(filnam)
    for i:=0; i<len(filnam); i++ {
        if bfil[i] == ' ' {
            bfil[i] = '_'
        }
    }

    filinfo, err := os.Stat("output")
    if os.IsNotExist(err) {
        return nil, fmt.Errorf("error CreTxtOutFile: sub-dir \"output\" does not exist!")
    }
    if err != nil {
        return nil, fmt.Errorf("error CreTxtOutFile: %v \n", err)
    }
    if !filinfo.IsDir() {
        return nil, fmt.Errorf("error CreTxtOutFile -- file \"output\" is not a directory!")
    }

    path:= "output/" + string(bfil[:]) + "." +ext

// need to delete the previous version
    outfil, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return nil, fmt.Errorf("cannot open output text file: %v!", err)
    }
    return outfil, nil
}


