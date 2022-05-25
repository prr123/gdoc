// golang library that converts a markdown file into a html/js file
// author: prr
// created: 25/5/2022
// copyright 2022 prr, azul software
//
// for changes see github
//
// start: CvtMdToDom
//

package mdDom

import (
	"fmt"
	"os"
    util "google/gdoc/gdocUtil"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)

type mdDomObj struct {
	Title string
	docName string
    docWidth float64
	docHeight float64
	ImgFoldName string
    imgCount int
	imgCounter int
    tableCount int
	tableCounter int
    parCount int
	title namStylTyp
	subtitle namStylTyp
	h1 namStylTyp
	h2 namStylTyp
	h3 namStylTyp
	h4 namStylTyp
	h5 namStylTyp
	h6 namStylTyp
	parent string
	listStack *[]cList
	docLists []docList
	headings []headingTyp
	sections []secTyp
	docPb []pbTyp
	docFtnotes []docFtnoteTyp
	namStylMap map[string]bool
	headCount int
	secCount int
	pbCount int
	elCount int
	spanCount int
	ftnoteCount int
	inImgCount int
	posImgCount int
	hrCount int
	htmlFil *os.File
	cssFil *os.File
	jsFil *os.File
//	folderName string
	folderPath string
    imgFoldNam string
//    imgFoldPath string
	Options *util.OptObj
}

type namStylTyp struct {
	count int
}

type dispObj struct {
	bodyHtml string
	bodyCss string
	script string
}

type secTyp struct {
	sNum int
	secElStart int
	secElEnd int
}

type pbTyp struct {
    el int
    parel int
}

type headingTyp struct {
	hdElEnd int
	hdElStart int
	namedStyl string
	id string
	text string
}

type docFtnoteTyp struct {
	el int
	parel int
	id string
	numStr string
}

type cList struct {
	cListId string
	cListCount int
	cOrd bool
}

type docList struct {
	listId string
	maxNestLev int64
	ord bool
}

type nestLevel struct {
	GlAl string
	GlFmt string
    GlSym string
	GlTyp string
	GlOrd bool
    Count int64
    FlInd float64
    StInd float64
	glTxtmap *textMap
}



type tblBorder struct {
	color string
	dash string
	width float64
}

type parMap struct {
	halign string
	headingId string
	direct bool
	indEnd float64
	indFlin float64
	indStart float64
	keepLines bool
	keepNext bool
	linSpac float64
	shading string
	spaceTop float64
	spaceBelow float64
	spaceMode bool
	tabs []*tabStop
	pad [4]float64
	margin [4]float64
	hasBorders bool
	bordTop parBorder
	bordLeft parBorder
	bordRight parBorder
	bordBot parBorder
	bordBet parBorder
}

type parBorder struct {
	pad float64
	width float64
	color string
	dash string
}

type tabStop struct {
	tabAlign string
	offset float64
}

type textMap struct {
	bckColor string
	baseOffset string
	bold bool
	italic bool
	underline bool
	strike bool
	fontSize float64
	txtColor string
	link bool
	fontType string
	fontWeight	int64
}

type linkMap struct {
	url string
	bookmark string
}

type elScriptObj struct {
	typ string
	txt string
	cl1 string
	cl2 string
	idStr string
	href string
	parent string
	newEl string
	comment string
}

type imgScriptObj struct {
	cl1 string
	cl2 string
	idStr string
	height int
	width int
	src string
	parent string
	title string
	desc string
	comment string
}

type tableScriptObj struct {
	cl1 string
	cl2 string
	idStr string
	rowCount int
	colCount int
	parent string
	newEl string
	comment string
	tblRows []tblRowScriptObj
}

type tblRowScriptObj struct {
	cl1 string
	cl2 string
	idStr string
	parent string
	newEl string
	tblCells []tblCellScriptObj
}

type tblCellScriptObj struct {
	cl1 string
	cl2 string
	idStr string
	parent string
	newEl string
}

type tblCell struct {
	pad [4] float64
	spad float64
	vert_align string
	bckcolor string
	border [4] tblBorder
	bwidth float64
	bdash string
	bcolor string
	cspan int
	rspan int
}



func findDocList(list []docList, listid string) (res int) {

	res = -1
	for i:=0; i< len(list); i++ {
		if list[i].listId == listid {
			return i
		}
	}
	return res
}

func pushLiStack(stack *[]cList, item cList)(nstack *[]cList) {
	if (stack == nil) {
		xx := make([]cList,0)
		stack = &xx
	}
	*stack= append(*stack, item)
	return stack
}

func popLiStack(stack *[]cList)(nstack *[]cList) {

	if stack == nil {return nil}

	n := len(*stack) -1
	if n >= 0 {
		xx	:= (*stack)[:n]
		nstack = &xx
	} else {
		nstack = nil
	}
	return nstack
}

func getLiStack(stack *[]cList) (item cList, nl int) {

	if stack == nil {
		return item, -1
	}
	nl = len(*stack) -1
	// review n>1 if there is a  stack!
	if nl >= 0 {
		item = (*stack)[nl]
	} else {
		nl = -1
	}
	return item, nl
}

func printLiStack(stack *[]cList) {
var item cList
var	n int
	if stack == nil {
		fmt.Println("*** no listStack ***")
		return
	}
	n = len(*stack) -1
	if n>=0 {
		item = (*stack)[n]
	} else {
		n = -1
	}
	fmt.Printf("list stack Nlevel: %d", n)
	if n >= 0 {
		fmt.Printf(" id: %s ordered: %t", item.cListId, item.cOrd)
	}
	fmt.Printf("\n")
	return
}

func printLiStackItem(listAtt cList, cNest int){
		fmt.Printf("list Att Nlevel: %d", cNest)
		if cNest >= 0 {
			fmt.Printf(" id: %s ordered: %t", listAtt.cListId, listAtt.cOrd)
		}
		fmt.Printf("\n")
}


func printTxtMap(txtMap *textMap) {

	fmt.Println("********* text map ****************")
	fmt.Printf("Base Offset: %s\n", txtMap.baseOffset)
	fmt.Printf("Bold Text:   %t\n", txtMap.bold)
	fmt.Printf("Italic Text: %t\n", txtMap.italic)
	fmt.Printf("Underline:   %t\n", txtMap.underline)
	fmt.Printf("Text Strike: %t\n", txtMap.strike)
	fmt.Printf("Font:        %s\n", txtMap.fontType)
	fmt.Printf("Font Weight: %d\n", txtMap.fontWeight)
	fmt.Printf("Font Size:   %.1f\n", txtMap.fontSize)
	fmt.Printf("Font Color:  %s\n", txtMap.txtColor)
	fmt.Printf("Font BckCol: %s\n", txtMap.bckColor)

	return
}



func cvtTxtMapCss(txtMap *textMap)(cssStr string) {

    cssStr =""
    if len(txtMap.baseOffset) > 0 {
        switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += "  vertical-align: sub;\n"
            case "SUBSCRIPT":
                cssStr += "  vertical-align: sup;\n"
            case "NONE":
                cssStr += "  vertical-align: baseline;\n"
            default:
            //error
                cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
        }
    }

    if txtMap.italic {
        cssStr += "  font-style: italic;\n"
    } else {
        cssStr += "  font-style: normal;\n"
    }

    textprop := ""
    switch {
    case txtMap.underline && txtMap.strike:
        textprop = "underline line-through"
    case txtMap.underline && !txtMap.strike:
        textprop = "underline"
    case !txtMap.underline && txtMap.strike:
        textprop = "line-through"
    case !txtMap.underline && !txtMap.strike:
        textprop = "none"
    }
    cssStr += fmt.Sprintf("  text-decoration: %s;\n", textprop)

    if len(txtMap.fontType) >0 { cssStr += fmt.Sprintf("  font-family: %s;\n", txtMap.fontType)}
    if txtMap.fontWeight > 0 {cssStr += fmt.Sprintf("  font-weight: %d;\n", txtMap.fontWeight)}
    if txtMap.fontSize >0 {cssStr += fmt.Sprintf("  font-size: %.2fpt;\n", txtMap.fontSize)}
    if len(txtMap.txtColor) >0 {cssStr += fmt.Sprintf("  color: %s;\n", txtMap.txtColor)}
    if len(txtMap.bckColor) >0 {cssStr += fmt.Sprintf("  background-color: %s;\n", txtMap.bckColor)}

    return cssStr
}


func cvtParMapCss(pMap *parMap, opt *util.OptObj)(cssStr string) {
    cssStr =""

    if len(pMap.halign) > 0 {
        switch pMap.halign {
            case "START":
                cssStr += "  text-align: left;\n"
            case "CENTER":
                cssStr += "  text-align: center;\n"
            case "END":
                cssStr += "  text-align: right;\n"
            case "JUSTIFIED":
                cssStr += "  text-align: justify;\n"
            default:
                cssStr += fmt.Sprintf("/* unrecognized Alignment %s */\n", pMap.halign)
        }

    }

    if pMap.linSpac > 0.0 {
        if opt.DefLinSpacing > 0.0 {
            cssStr += fmt.Sprintf("  line-height: %.2f;\n", opt.DefLinSpacing*pMap.linSpac)
        } else {
            cssStr += fmt.Sprintf("  line-height: %.2f;\n", pMap.linSpac)
        }
    }

    if pMap.indFlin > 0.0 {
        cssStr += fmt.Sprintf("  text-indent: %.1fpt;\n", pMap.indFlin)
    }

    margin := false
    lmarg := 0.0
    if pMap.indStart > 0.0 {
		lmarg = pMap.indStart
		margin = true
    }

    rmarg := 0.0
    if pMap.indEnd > 0.0 {
		rmarg = pMap.indEnd
		margin = true
    }

    tmarg := 0.0
    if pMap.spaceTop > 0.0 {
		tmarg = pMap.spaceTop
		margin = true
    }

    bmarg := 0.0
    if pMap.spaceBelow > 0.0 {
		bmarg = pMap.spaceBelow
		margin = true
    }

    if margin {cssStr += fmt.Sprintf("  margin: %.0f %.0f %.0f %.0f;\n", tmarg, rmarg, bmarg, lmarg)}

    if !pMap.hasBorders { return cssStr }
    cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
    cssStr += fmt.Sprintf("  border-top: %.1fpt %s %s;\n", pMap.bordTop.width, util.GetDash(pMap.bordTop.dash), pMap.bordTop.color)
    cssStr += fmt.Sprintf("  border-right: %.1fpt %s %s;\n", pMap.bordRight.width, util.GetDash(pMap.bordRight.dash), pMap.bordRight.color)
    cssStr += fmt.Sprintf("  border-bottom: %.1fpt %s %s;\n", pMap.bordBot.width, util.GetDash(pMap.bordBot.dash), pMap.bordBot.color)
    cssStr += fmt.Sprintf("  border-left: %.1fpt %s %s;\n", pMap.bordLeft.width, util.GetDash(pMap.bordLeft.dash), pMap.bordLeft.color)

    return cssStr
}



func addDispObj(src, add *dispObj) {
	src.bodyHtml += add.bodyHtml
	src.bodyCss += add.bodyCss
	src.script += add.script
	return
}

func creHtmlDocHead(docNam string)(outstr string) {
    outstr = "<!DOCTYPE html>\n"
    outstr += fmt.Sprintf("<!-- file: %s -->\n", docNam + "Dom")
    outstr += "<head>\n<style>\n"
    return outstr
}

func creHtmlScript()(outstr string) {
    outstr = "</style><script>\n"
    return outstr
}

func creHtmlBody()(outstr string) {
    outstr = "</script>\n<body>\n"
    return outstr
}

func creHtmlDocEnd()(string) {
    return "</body></html>\n"
}

func creTocSecCss(docName string)(cssStr string) {

	cssStr = fmt.Sprintf(".%s_main.top {\n", docName)
	cssStr += "  padding: 10px 0 10px 0;\n"
	cssStr += "}\n"

	cssStr += fmt.Sprintf(".%s_title.leftTitle_UL {\n", docName)
	cssStr += "  text-align: start;\n"
	cssStr += "  text-decoration-line: underline;\n"
	cssStr += "}\n"

	cssStr += fmt.Sprintf(".%s_title.leftTitle {\n", docName)
	cssStr += "  text-align: start;\n"
	cssStr += "  text-decoration-line: none;\n"
	cssStr += "}\n"

	cssStr += fmt.Sprintf(".%s_noUl {\n", docName)
	cssStr += "  text-decoration: none;\n"
	cssStr += "}\n"

	return cssStr
}

func creTocCss(docName string)(cssStr string) {
	cssStr = fmt.Sprintf(".%_div.toc {\n", docName)

	cssStr += "}\n"
	return cssStr
}

func creSecCss(docName string)(cssStr string){

	cssStr = fmt.Sprintf(".%s_main.sec {\n", docName)
	cssStr += "}\n"

	cssStr += fmt.Sprintf(".%s_page {\n", docName)
	cssStr += "  text-align: right;\n"
	cssStr += "  margin: 0;\n"
	cssStr += "}\n"
	return cssStr
}

func creFtnoteCss(docName string)(cssStr string){
	//css footnote
	cssStr = fmt.Sprintf(".%s_ftnote {\n", docName)
//	cssStr += "vertical-align: super;"
	cssStr += "  color: purple;\n"
	cssStr += "}\n"
	return cssStr
}

func creHtmlDocDiv(docName string)(htmlStr string) {
	htmlStr = fmt.Sprintf("<div class=\"%s_doc\">\n", docName)
	return htmlStr
}

func creElFuncScript(imgFun bool, tableFun bool) (jsStr string) {
	jsStr = "function addEl(elObj) {\n"
	jsStr += "  let el = document.createElement(elObj.typ);\n"
	jsStr += "  if (elObj.cl1 != null) {el.classList.add(elObj.cl1);}\n"
	jsStr += "  if (elObj.cl2 != null) {el.classList.add(elObj.cl2);}\n"
	jsStr += "  if (elObj.idStr != null) {el.setAttribute(\"id\", elObj.idStr);}\n"
	jsStr += "  if (elObj.href != null) {el.href=elObj.href};\n"
	jsStr += "  if (elObj.txt != null) {\n"
	jsStr += "    var text =  document.createTextNode(elObj.txt);\n"
	jsStr += "    el.appendChild(text);\n"
	jsStr += "  }\n"
	jsStr += "  elp = elObj.parent;\n"
	jsStr += "  elp.appendChild(el);\n"
	jsStr += "  return el\n}\n"
	jsStr += "function clearObj(elObj) {\n"
	jsStr += "  for (key in elObj) {elObj[key] = null;}\n"
	jsStr += "  return\n}\n"
	jsStr += "function addTxt(elObj) {\n"
	jsStr += "  let el = elObj.parent;\n"
	jsStr += "  if (elObj.txt != null) {\n"
	jsStr += "    var text =  document.createTextNode(elObj.txt);\n"
	jsStr += "    el.appendChild(text);\n"
	jsStr += "  }\n}\n"
	jsStr = "function addLink(elObj) {\n"
	jsStr += "  let el = document.createElement('a');\n"
	jsStr += "  el.setAttribute('href', elObj.href);\n"
	jsStr += "  if (elObj.cl1 != null) {el.classList.add(elObj.cl1);}\n"
	jsStr += "  if (elObj.cl2 != null) {el.classList.add(elObj.cl2);}\n"
	jsStr += "  if (elObj.idStr != null) {el.setAttribute(\"id\", elObj.idStr);}\n"
	jsStr += "  if (elObj.txt != null) {\n"
	jsStr += "    var text =  document.createTextNode(elObj.txt);\n"
	jsStr += "    el.appendChild(text);\n"
	jsStr += "  }\n"
	jsStr += "  elp = elObj.parent;\n"
	jsStr += "  elp.appendChild(el);\n"
	jsStr += "  return\n}\n"
	if imgFun {
		jsStr += "function addImgEl(imgObj) {\n"
		jsStr += "  if (imgObj.src == null) {return\n}\n"
		jsStr += "  var img = new Image(imgObj.width, imgObj.height);\n"
		jsStr += "  if (imgObj.idStr != null) {img.setAttribute(\"id\", imgObj.idStr);}\n"
		jsStr += "  if (imgObj.cl1 != null) {img.classList.add(imgObj.cl1);}\n"
		jsStr += "  if (imgObj.cl2 != null) {img.classList.add(imgObj.cl2);}\n"
		jsStr += "  img.src = imgObj.src;\n"
		jsStr += "  img.alt = imgObj.alt;\n"
//		jsStr += ""
		jsStr += "  imgp = imgObj.parent;\n"
		jsStr += "  imgp.appendChild(img);\n"
		jsStr += "  return\n}\n"
	}
	if tableFun {
		jsStr += "function addTblEl(tblObj) {\n"
		jsStr += "  var tbl = document.createElement('table');\n"
		jsStr += "  var tblBody = document.createElement('tbody');\n"
		jsStr += "  var colgrp = document.createElement('colgroup');\n"

		jsStr += "  for (var ir = 0; i < tblObj.rows; i++) {\n"
		jsStr += "    var tblRow = document.createElement('tr');\n"
		jsStr += "    var row = tblObj.row[ir];\n"
		jsStr += "    if (row.idStr != null) {tblRow.setAttribute(\"id\", row.idStr);}\n"
		jsStr += "    if (row.cl1 != null) {tblRow.classList.add(row.cl1);}\n"
		jsStr += "    if (row.cl2 != null) {tblRow.classList.add(row.cl2);}\n"
		jsStr += "    for (var ic = 0; i < tblObj.cols; i++) {\n"
		jsStr += "      var tblCell = document.createElement('td');\n"
		jsStr += "      var col = tblObj.row[ir].col[ic];\n"
		jsStr += "      if (col.idStr != null) {tblCell.setAttribute(\"id\", col.idStr);}\n"
		jsStr += "      if (col.cl1 != null) {tblCell.classList.add(col.cl1);}\n"
		jsStr += "      if (col.cl2 != null) {tblCell.classList.add(col.cl2);}\n"
		jsStr += "      tblRow.appendChild(tblCell);\n"
		jsStr += "	  }/n"
		jsStr += "	  tblBody.appendChild(tblRow);\n"
		jsStr += "	}/n"
		jsStr += "  tblp = tblObj.parent;\n"
		jsStr += "  tblp.appendChild(tbl);\n"
		jsStr += "  return tbl\n}\n"
	}


	jsStr += "function addBodyElScript(divDoc) {\n"
	jsStr += "  const elObj = {};\n"
	jsStr += "  const imgObj = {};\n"
	jsStr += "  const tblObj = {};\n"

	return jsStr
}


func addTxtElToDom(elObj elScriptObj)(script string) {
	script = "//addTxtEl \n"
	script += "//" + elObj.comment + "\n"
	if !(len(elObj.parent) > 0) {
		script += "// no el parent provided!\n"
		return script
	}
	if !(len(elObj.txt) > 0) {
		script += "// no text provided!\n"
		return script
	}
	script = "    for (key in elObj) {elObj[key] = null;}\n"
	script += fmt.Sprintf("  elObj.txt = '%s';\n", elObj.txt)
	script += fmt.Sprintf("  elObj.parent = %s;\n", elObj.parent)
	script += fmt.Sprintf("  addTxt(elObj);\n")
	return script
}

func addElToDom(elObj elScriptObj)(script string) {

	script = "// addEl \n"
	script += "// " + elObj.comment + "\n"
	if !(len(elObj.parent) > 0) {
		script += "// error - no el parent provided!\n"
		return script
	}
	if !(len(elObj.typ) > 0) {
		script += "// error - no el type provided!\n"
		return script
	}
	script = "  for (key in elObj) {elObj[key] = null;}\n"
	if len(elObj.cl1) > 0 {script += fmt.Sprintf("  elObj.cl1 = '%s';\n", elObj.cl1)}
	if len(elObj.cl2) > 0 {script += fmt.Sprintf("  elObj.cl2 = '%s';\n", elObj.cl2)}
	if len(elObj.idStr) > 0 {script += fmt.Sprintf("  elObj.idStr = '%s';\n", elObj.idStr)}
	if len(elObj.txt) > 0 {script += fmt.Sprintf("  elObj.txt = '%s';\n", elObj.txt)}
	script += fmt.Sprintf("  elObj.parent = %s;\n", elObj.parent)
	script += fmt.Sprintf("  elObj.typ = '%s';\n", elObj.typ)
	script += fmt.Sprintf("  %s = addEl(elObj);\n", elObj.newEl)
	return script
}

func addLinkToDom(elObj elScriptObj)(script string) {

	script = "// addLinkEl \n"
	script += "// " + elObj.comment + "\n"
	if !(len(elObj.parent) > 0) {
		script += "// error - no el parent provided!\n"
		return script
	}
	if !(len(elObj.txt) > 0) {
		script += "// error - no text provided!\n"
		return script
	}
	if !(len(elObj.href) > 0) {
		script += "// error - no href provided!\n"
		return script
	}
	script = "  for (key in elObj) {elObj[key] = null;}\n"
	if len(elObj.cl1) > 0 {script += fmt.Sprintf("  elObj.cl1 = '%s';\n", elObj.cl1)}
	if len(elObj.cl2) > 0 {script += fmt.Sprintf("  elObj.cl2 = '%s';\n", elObj.cl2)}
	if len(elObj.idStr) > 0 {script += fmt.Sprintf("  elObj.idStr = '%s';\n", elObj.idStr)}
	if len(elObj.txt) > 0 {script += fmt.Sprintf("  elObj.txt = '%s';\n", elObj.txt)}
	script += fmt.Sprintf("  elObj.parent = %s;\n", elObj.parent)
	script += fmt.Sprintf("  elObj.typ = 'a';\n")
	script += fmt.Sprintf("  addLink(elObj);\n")
	return script
}

func addImgElToDom(imgObj imgScriptObj)(script string) {

	script = "// addEl \n"
	script += "// " + imgObj.comment + "\n"
	if !(len(imgObj.parent) > 0) {
		script += "// error - no el parent provided!\n"
		return script
	}
	script = "  for (key in imgObj) {imgObj[key] = null;}\n"
	if len(imgObj.cl1) > 0 {script += fmt.Sprintf("  imgObj.cl1 = '%s';\n", imgObj.cl1)}
	if len(imgObj.cl2) > 0 {script += fmt.Sprintf("  imgObj.cl2 = '%s';\n", imgObj.cl2)}
	if len(imgObj.idStr) > 0 {script += fmt.Sprintf("  imgObj.idStr = '%s';\n", imgObj.idStr)}
	script += fmt.Sprintf("  imgObj.parent = %s;\n", imgObj.parent)
	script += fmt.Sprintf("  addImgEl(imgObj);\n")
	return script
}

func addTblToDom(tblObj tableScriptObj)(script string) {

	script = "// *** addTbl ***\n"
	if len(tblObj.comment) > 0 {script += "// " + tblObj.comment + "\n"}
	if !(len(tblObj.parent) > 0) {
		script += "// error - no el parent provided!\n"
		return script
	}
	if !(tblObj.rowCount > 0) {
		script += "// error -- no table rows provided!\n"
		return script
	}
	if !(tblObj.colCount > 0) {
		script += "// error -- no table columns provided!\n"
		return script
	}

	script = "  for (key in tblObj) {tblObj[key] = null;}\n"

	if len(tblObj.cl1) > 0 {script += fmt.Sprintf("  tblObj.cl1 = '%s';\n", tblObj.cl1)}
	if len(tblObj.cl2) > 0 {script += fmt.Sprintf("  tblebj.cl2 = '%s';\n", tblObj.cl2)}
	if len(tblObj.idStr) > 0 {script += fmt.Sprintf("  tblObj.idStr = '%s';\n", tblObj.idStr)}

	script += fmt.Sprintf("  tblObj.parent = %s;\n", tblObj.parent)
//	for irow:=0; irow < tblObj.rowCount; irow++ {

//		for icol:=0; icol < tblObj.colCount; icol++ {
	script += fmt.Sprintf("  tblObj.rowCount = %d\n", tblObj.rowCount)
	script += fmt.Sprintf("  tblObj.colCount = %d\n", tblObj.colCount)


	script += fmt.Sprintf("  tbl = addTblEl(tblObj);\n")

//	script += fmt.Sprintf("  fillTblEl(tblObj);\n")

	return script
}

func addDivMainScript(docName string) (jsStr string) {
    jsStr += "  let divMain = document.createElement('div');\n"
    jsStr += fmt.Sprintf("  divMain.classList.add('%s_main');\n", docName)
    jsStr += "  divDoc.appendChild(divMain);\n"
	return jsStr
}

func creDocDivScript(docName string)(jsStr string) {

	jsStr = "  return\n}\n"
	jsStr += "function dispDoc() {\n"
    jsStr += "  let divDoc = document.createElement('div');\n"
    jsStr += fmt.Sprintf("  divDoc.classList.add('%s_doc');\n", docName)
    jsStr += "  document.body.appendChild(divDoc);\n"
	jsStr += "  addBodyElScript(divDoc);\n"
	jsStr += "}\n"
    jsStr += "document.addEventListener(\"DOMContentLoaded\", dispDoc);\n"
    return jsStr
}

func stripCrText(inp string) (out string) {
	ilen := len(inp)
	if inp[ilen -1] == '\n' {
		out = inp[:ilen-1]
	} else {
		out = inp
	}
	return out
}

func cvtTextjs(inp string) (out string) {
	ilen := len(inp)
	outb := make([]byte, ilen + 20, 100)
	j:=0
	ret :=0
	for i:=0; i<ilen; i++ {
		if inp[i] == '\n' {
			outb[j] = '\\'
			j++
			outb[j] = 'r'
			ret++
			if ret > 10 {
				outb = append(outb, make([]byte, 20)...)
				ret = 0
			}
		} else {
			outb[j] = inp[i]
		}
		j++
	}
	return string(outb)
}


func addParScript(docName string)(jsStr string) {

	jsStr = "function addPar(txt, idStr, cl1, cl2) {\n"
    jsStr += "  let par = document.createElement('p');\n"
    jsStr += fmt.Sprintf("  p.classList.add('%s_p');\n", docName)
    jsStr += "  divMain.appendChild(par);\n"
	jsStr += "}\n"
    return jsStr
}

func (dObj *mdDomObj) printHeadings() {

	if len(dObj.headings) == 0 {
		fmt.Println("*** no Headings ***")
		return
	}

	fmt.Printf("**** Headings: %d ****\n", len(dObj.headings))
	for i:=0; i< len(dObj.headings); i++ {
		fmt.Printf("  heading %3d  Id: %-15s named Styl: %-12s El Start:%3d End:%3d\n", i, dObj.headings[i].id,
			dObj.headings[i].namedStyl, dObj.headings[i].hdElStart, dObj.headings[i].hdElEnd)
	}
}


func (dObj *mdDomObj) cvtParMapCss(pMap *parMap)(cssStr string) {
	cssStr =""

	if len(pMap.halign) > 0 {
		switch pMap.halign {
			case "START":
				cssStr += "  text-align: left;\n"
			case "CENTER":
				cssStr += "  text-align: center;\n"
			case "END":
				cssStr += "  text-align: right;\n"
			case "JUSTIFIED":
				cssStr += "  text-align: justify;\n"
			default:
				cssStr += fmt.Sprintf("/* unrecognized Alignment %s */\n", pMap.halign)
		}

	}

	opt := dObj.Options
	if pMap.linSpac > 0.0 {
		if opt.DefLinSpacing > 0.0 {
			cssStr += fmt.Sprintf("  line-height: %.2f;\n", opt.DefLinSpacing*pMap.linSpac)
		} else {
			cssStr += fmt.Sprintf("  line-height: %.2f;\n", pMap.linSpac)
		}
	}

	if pMap.indFlin > 0.0 {
		cssStr += fmt.Sprintf("  text-indent: %.1fpt;\n", pMap.indFlin)
	}

	mlCss :=""
	if pMap.indStart > 0.0 {
		mlCss = fmt.Sprintf("%.1fpt", pMap.indStart)
	} else {
		mlCss = fmt.Sprintf("0")
	}
	mrCss:=""
	if pMap.indEnd > 0.0 {
		mrCss = fmt.Sprintf("%.1fpt", pMap.indEnd)
	} else {
		mrCss = fmt.Sprintf("0")
	}
	mtCss := ""
	if pMap.spaceTop > 0.0 {
		mtCss = fmt.Sprintf("%.1fpt", pMap.spaceTop)
	} else {
		mtCss = fmt.Sprintf("0")
	}
	mbCss := ""
	if pMap.spaceBelow > 0.0 {
		mbCss = fmt.Sprintf("%.1fpt", pMap.spaceBelow)
	} else {
		mbCss = fmt.Sprintf("0")
	}

	cssStr += fmt.Sprintf("  margin: %s %s %s %s;\n", mtCss, mrCss, mbCss, mlCss)

	if !pMap.hasBorders { return cssStr }
	cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
	cssStr += fmt.Sprintf("  border-top: %.1fpt %s %s;\n", pMap.bordTop.width, util.GetDash(pMap.bordTop.dash), pMap.bordTop.color)
	cssStr += fmt.Sprintf("  border-right: %.1fpt %s %s;\n", pMap.bordRight.width, util.GetDash(pMap.bordRight.dash), pMap.bordRight.color)
	cssStr += fmt.Sprintf("  border-bottom: %.1fpt %s %s;\n", pMap.bordBot.width, util.GetDash(pMap.bordBot.dash), pMap.bordBot.color)
	cssStr += fmt.Sprintf("  border-left: %.1fpt %s %s;\n", pMap.bordLeft.width, util.GetDash(pMap.bordLeft.dash), pMap.bordLeft.color)

	return cssStr
}

func (dObj *mdDomObj) initMdDom(folderPath string, options *util.OptObj) (err error) {
	var sec secTyp

	// need to transform file name
	// replace spaces with underscore

	dNam := dObj.Title
	x := []byte(dNam)
	for i:=0; i<len(x); i++ {
		if x[i] == ' ' {
			x[i] = '_'
		}
	}
	dObj.docName = string(x[:])

	if options == nil {
		defOpt := new(util.OptObj)
		util.GetDefOption(defOpt)
		if defOpt.Verb {util.PrintOptions(defOpt)}
		dObj.Options = defOpt
		dObj.Options.DivBorders = true
	} else {
		dObj.Options = options
	}


	dObj.namStylMap = make(map[string]bool, 8)

	dObj.namStylMap["NORMAL_TEXT"] = true
	dObj.namStylMap["TITLE"] = false
	dObj.namStylMap["SUBTITLE"] = false
	dObj.namStylMap["HEADING_1"] = false
	dObj.namStylMap["HEADING_2"] = false
	dObj.namStylMap["HEADING_3"] = false
	dObj.namStylMap["HEADING_4"] = false
	dObj.namStylMap["HEADING_5"] = false
	dObj.namStylMap["HEADING_6"] = false

	// footnotes
	dObj.ftnoteCount = 0

	// span
	dObj.spanCount = 0

	//horizntal rule
	dObj.hrCount = 0

	// section breaks
	parHdEnd := 0
	// last element of section
	secPtEnd := 0
	// set up first page
	sec.secElStart = 0
	dObj.sections = append(dObj.sections, sec)
	seclen := len(dObj.sections)
//		fmt.Println("el: ", el, "section len: ", seclen, " secPtEnd: ", secPtEnd)


	hdlen := len(dObj.headings)
	if hdlen > 0 {
		dObj.headings[hdlen-1].hdElEnd = parHdEnd
	}
	seclen = len(dObj.sections)
	if seclen > 0 {
		dObj.sections[seclen-1].secElEnd = secPtEnd
	}

	if dObj.Options.Verb {
		fmt.Printf("************ Headings: %2d **************\n", len(dObj.headings))
		for i:=0; i< len(dObj.headings); i++ {
			fmt.Printf("  heading %3d  Id: %-15s text: %-20s El Start:%3d End:%3d\n", i, dObj.headings[i].id, dObj.headings[i].text, 
				dObj.headings[i].hdElStart, dObj.headings[i].hdElEnd)
		}

		fmt.Printf("***********  Named Styles: %2d **********\n", len(dObj.headings))
		for namedTyp, val := range dObj.namStylMap {
			if val {
				fmt.Printf("Style: %s\n", namedTyp)
			}
		}

		fmt.Printf("************ Doc Pages: %2d *************\n", len(dObj.sections))
		if len(dObj.sections) > 1 {
			for i:=0; i< len(dObj.sections); i++ {
				fmt.Printf("  Page %3d  El Start:%3d End:%3d\n", i, dObj.sections[i].secElStart, dObj.sections[i].secElEnd)
			}
		}

		fmt.Printf("************ Lists: %2d *****************\n", len(dObj.docLists))
		for i:=0; i< len(dObj.docLists); i++ {
			fmt.Printf("list %3d id: %s max level: %d ordered: %t\n", i, dObj.docLists[i].listId, dObj.docLists[i].maxNestLev, 
			dObj.docLists[i].ord)
		}

		fmt.Printf("************ Footnotes: %2d *************\n", len(dObj.docFtnotes))
		for i:=0; i< len(dObj.docFtnotes); i++ {
			ftn := dObj.docFtnotes[i]
			fmt.Printf("ft %3d: Number: %-4s id: %-15s el: %3d parel: %3d\n", i, ftn.numStr, ftn.id, ftn.el, ftn.parel)
		}

		fmt.Printf("****************************************\n\n")
	}


// images
	dObj.inImgCount = 0
	dObj.posImgCount = 0

// create folders
    fPath, fexist, err := util.CreateFileFolder(folderPath, dObj.docName)
//    fPath, _, err := util.CreateFileFolder(folderPath, dObj.docName)
    if err!= nil {
        return fmt.Errorf("error -- util.CreateFileFolder: %v", err)
    }
    dObj.folderPath = fPath

	if dObj.Options.Verb {
		fmt.Println("*************** Output File ****************")
		fmt.Printf("folder path: %s ", fPath)
		fstr := "is new!"
		if fexist { fstr = "exists!" }
		fmt.Printf("%s\n", fstr)
		fmt.Println("********************************************")
	}

    // create output file path/outfilNam.html
	outfilNam := dObj.docName + "Dom"
    outfil, err := util.CreateOutFil(fPath, outfilNam, "html")
    if err!= nil {
        return fmt.Errorf("error -- util.CreateOutFil: %v", err)
    }
    dObj.htmlFil = outfil

    totObjNum := dObj.inImgCount + dObj.posImgCount
//  if totObjNum == 0 {return nil}


    if dObj.Options.CreImgFolder && (totObjNum > 0) {
        imgFoldPath, err := util.CreateImgFolder(fPath ,dObj.docName)
        if err != nil {
            return fmt.Errorf("error -- CreateImgFolder: could create ImgFolder: %v!", err)
        }
        dObj.imgFoldNam = imgFoldPath
//        err = util.DownloadImages(imgFoldPath, dObj.Options)
//        if err != nil {
//            return fmt.Errorf("error -- downloadImages could download images: %v!", err)
//        }
    }

//    dObj.parCount = len(doc.Body.Content)

	return nil
}


func (dObj *mdDomObj) closeList(nl int) {
	// ends a list

	if (dObj.listStack == nil) {return}

	stack := dObj.listStack
	n := len(*stack)

	for i := n -1; i > nl; i-- {
		nstack := popLiStack(stack)
		dObj.listStack = nstack
	}
	return
}



func (dObj *mdDomObj) cvtDocNamedStyles()(cssStr string, err error) {
// method that creates the css for the named Styles used in the document

	// the normal_text style are already defined in div_main
	// so the css attributes for other named styles only need to show the difference 
	// to the normal style


	for namedTyp, res := range dObj.namStylMap {
		if namedTyp == "NORMAL_TEXT" { continue}
		if !res {continue}

//		namParStyl, namTxtStyl, err := dObj.getNamedStyl(namedTyp)
		if err != nil {
			cssStr += fmt.Sprintf("  /* cvtNamedStyle: named Style not recognized */\n")
			continue
		}

		cssPrefix := ""
		switch namedTyp {
		case "TITLE":
			cssPrefix = fmt.Sprintf(".%s_title {\n", dObj.docName)

		case "SUBTITLE":
			cssPrefix = fmt.Sprintf(".%s_subtitle {\n",dObj.docName)

		case "HEADING_1":
			cssPrefix =fmt.Sprintf(".%s_h1 {\n",dObj.docName)

		case "HEADING_2":
			cssPrefix =fmt.Sprintf(".%s_h2 {\n",dObj.docName)

		case "HEADING_3":
			cssPrefix =fmt.Sprintf(".%s_h3 {\n",dObj.docName)

		case "HEADING_4":
			cssPrefix =fmt.Sprintf(".%s_h4 {\n",dObj.docName)

		case "HEADING_5":
			cssPrefix =fmt.Sprintf(".%s_h5 {\n",dObj.docName)

		case "HEADING_6":
			cssPrefix =fmt.Sprintf(".%s_h6 {\n",dObj.docName)

		case "NORMAL_TEXT":

		case "NAMED_STYLE_TYPE_UNSPECIFIED":

		default:

		}

		if len(cssPrefix) > 0 {
//			parCss, _ := cvtParMapStylCss(defParMap, namParStyl, dObj.Options)
//			txtCss := cvtTxtMapStylCss(defTxtMap, namTxtStyl)
//			cssStr += cssPrefix + parCss + txtCss + "}\n"
		}
	}
	return cssStr, nil
}


func (dObj *mdDomObj) createDivHead(divName, idStr string) (divObj dispObj, err error) {
	var htmlStr, cssStr, script string
	//gdoc division css

	if len(divName) == 0 { return divObj, fmt.Errorf("createDivHead: no divNam!") }

	// css
	cssStr = fmt.Sprintf(".%s_main.%s {\n", dObj.docName, divName)

	// script
	if len(divName) == 0 {
//		htmlStr = fmt.Sprintf("<div class=\"%s_main\"", dObj.docName)
//		script = 
	} else {
		
//		htmlStr = fmt.Sprintf("<div class=\"%s_main %s\"", dObj.docName, divName)
	}

	if len(idStr) > 0 {

//		htmlStr += fmt.Sprintf(" id=\"%s\"", idStr)
	}

	htmlStr += ">\n"

	divObj.bodyCss = cssStr
//	divObj.bodyHtml = htmlStr
	divObj.script = script

	return divObj, nil
}

func (dObj *mdDomObj) creSecDivDom() (secHd *dispObj) {
	var secHead dispObj
	var scriptStr string
	var divObj, parObj elScriptObj

	if !dObj.Options.Sections {return nil}

	if len(dObj.sections) < 2 {return nil}

	//html
	// fmt.Sprintf("<div class=\"%s_main top\" id=\"%s_sectoc\">\n", dObj.docName, dObj.docName)
	divObj.parent = dObj.parent
	divObj.typ="div"
	divObj.newEl = "divSec"
	divObj.cl1 = fmt.Sprintf("%s_main_top", dObj.docName)
	divObj.idStr = fmt.Sprintf("%s_sectoc", dObj.docName)
	scriptStr = addElToDom(divObj)

	// fmt.Sprintf("<p class=\"%s_title %s_leftTitle_UL\">Sections</p>\n",dObj.docName, dObj.docName)
	parObj.parent = "divSec"
	parObj.typ = "p"
	parObj.newEl = "pel"
	parObj.cl1 = dObj.docName + "_title"
	parObj.cl2 = dObj.docName + "_leftTitle_UL"
	parObj.txt = "Sections"
	scriptStr += addElToDom(parObj)

	for i:=0; i< len(dObj.sections); i++ {
		// fmt.Sprintf("  <p class=\"%s_p\"><a href=\"#%s_sec_%d\">Page: %3d</a></p>\n", dObj.docName, dObj.docName, i, i)
		parObj.parent = "divSec"
		parObj.typ = "p"
		parObj.newEl = "pel"
		parObj.cl1 = dObj.docName + "_p"
		scriptStr += addElToDom(parObj)

		parObj.parent = "pel"
		parObj.href = fmt.Sprintf("#%s_sec_%d", dObj.docName, i)
		parObj.txt = fmt.Sprintf("Page: %3d", i)
		scriptStr += addLinkToDom(parObj)
	}

	secHead.script = scriptStr
	return &secHead
}

//section
func (dObj *mdDomObj) creSecHeadToDom(ipage int) (secObj dispObj) {
// method that creates a distinct html dvision per section with a page heading

	var divObj, parObj elScriptObj
	var linkObj elScriptObj

	//css
	prefixCss := fmt.Sprintf(".%s_main.sec_%d {\n", dObj.docName, ipage)
	secCss := ""
	suffixCss := "}/n"

	if len(secCss) > 0 {secObj.bodyCss = prefixCss + secCss + suffixCss}

	// html
	// fmt.Sprintf("<div class=\"%s_main sec_%d\" id=\"%s_sec_%d\">\n", dObj.docName, ipage, dObj.docName, ipage)
	// fmt.Sprintf("<p class=\"%s_page\"><a href=\"#%s_sectoc\">Page %d</a></p>\n", dObj.docName, dObj.docName, ipage)

	// script
	divObj.parent = "divDoc"
	divObj.newEl = fmt.Sprintf("%s_main_%d", dObj.docName, ipage)
	divObj.typ = "div"
	dObj.parent = divObj.newEl
	divObj.cl1 = fmt.Sprintf("%s_main", dObj.docName)
	divObj.cl2 = fmt.Sprintf("sec_%d", ipage)
	divObj.idStr = fmt.Sprintf("%s_sec_%d", dObj.docName, ipage)
	secObj.script += addElToDom(divObj)

	parObj.parent = dObj.parent
	parObj.typ = "p"
	parObj.newEl = "ptop"
	parObj.cl1 = fmt.Sprintf("%s_page", dObj.docName)
	secObj.script += addElToDom(parObj)

	linkObj.parent = "ptop"
	linkObj.txt = fmt.Sprintf("Page %d", ipage)
	linkObj.href = fmt.Sprintf("%s_sectoc", dObj.docName)
	secObj.script += addLinkToDom(linkObj)

	return secObj
}

func (dObj *mdDomObj) creCssDocHead() (headCss string, err error) {

	var cssStr string


    //gdoc default el css and doc css
/*
    cssStr = fmt.Sprintf(".%s_doc {\n", dObj.docName)
    cssStr += fmt.Sprintf("  margin-top: %.1fmm; \n",docStyl.MarginTop.Magnitude*PtTomm)
    cssStr += fmt.Sprintf("  margin-bottom: %.1fmm; \n",docStyl.MarginBottom.Magnitude*PtTomm)
    cssStr += fmt.Sprintf("  margin-right: %.2fmm; \n",docStyl.MarginRight.Magnitude*PtTomm)
    cssStr += fmt.Sprintf("  margin-left: %.2fmm; \n",docStyl.MarginLeft.Magnitude*PtTomm)
    if dObj.docWidth > 0 {cssStr += fmt.Sprintf("  width: %.1fmm;\n", dObj.docWidth*PtTomm)}
	if dObj.Options.DivBorders {
		cssStr += "  border: solid red;\n"
		cssStr += "  border-width: 1px;\n"
	}
	cssStr += "}\n"
	headCss = cssStr
*/
	//css default text style
	cssStr = fmt.Sprintf(".%s_main {\n", dObj.docName)
//	parStyl, txtStyl, err := dObj.getNamedStyl("NORMAL_TEXT")
	if err != nil {
		return headCss, fmt.Errorf("creHeadCss: %v", err)
	}

	cssStr += "  display:block;\n"
	cssStr += "  margin: 0;\n"
	if dObj.Options.DivBorders {
		cssStr += "  border: solid green;\n"
		cssStr += "  border-width: 1px;\n"
	}
//	cssStr += cvtTxtMapCss(defTxtMap)
	cssStr += "}\n"
	headCss += cssStr

//	hdcss, err := dObj.cvtDocNamedStyles()
//	if err != nil {errStr := fmt.Sprintf("cvtDocNamedStyles %v", err)}
//	headCss += hdcss + errStr

	// paragraph default style
//    pCssStr := cvtParMapCss(defParMap, dObj.Options)
	pCssStr := ""
	cssStr =""
	if len(pCssStr) > 0 {
		cssStr += fmt.Sprintf(".%s_p {\n", dObj.docName)
		cssStr += "  margin: 0;\n"
		cssStr += pCssStr + "}\n"
	}
	headCss += cssStr

	// list css strings
	cssStr = ""
	for i:=0; i<len(dObj.docLists); i++ {
		listid := dObj.docLists[i].listId
		listClass := listid[4:]

		switch dObj.docLists[i].ord {
			case true:
				cssStr += fmt.Sprintf(".%s_ol {\n", listClass)
//				glyphNum := "none"
				cssStr += fmt.Sprintf("  list-style-type: none;\n")
				cssStr += fmt.Sprintf("  list-style-position: outside;\n")
				cssStr += fmt.Sprintf("}\n")

			case false:
				cssStr += fmt.Sprintf(".%s_ul {\n", listClass)
				cssStr += fmt.Sprintf("  list-style-type: none;\n")
				cssStr += fmt.Sprintf("  list-style-position: outside;\n")
				cssStr += fmt.Sprintf("}\n")
		}
		cssStr += fmt.Sprintf(".%s_li {\n", listClass)
		cssStr += fmt.Sprintf("  display: list-item;\n")
		cssStr += fmt.Sprintf("  text-align: start;\n")
		cssStr += fmt.Sprintf("  padding-left: 6pt;\n")
		cssStr += fmt.Sprintf("}\n")

//		nestLev0 := listProp.NestingLevels[0]
//		defGlyphTxtMap := fillTxtMap(nestLev0.TextStyle)

		cumIndent := 0.0

		for nl:=0; nl <= int(dObj.docLists[i].maxNestLev); nl++ {
//			nestLev := listProp.NestingLevels[nl]
//			nestLev := 10

//			glyphStr := util.GetGlyphStr(nestLev)
			switch dObj.docLists[i].ord {
				case true:
					cssStr += fmt.Sprintf(".%s_ol.nL_%d {\n", listClass, nl)
				case false:
					cssStr += fmt.Sprintf(".%s_ul.nL_%d {\n", listClass, nl)
			}


//			idFl := nestLev.IndentFirstLine - cumIndent
//			idSt := nestLev.IndentStart - cumIndent
			idFl := 18.0
			idSt := 36.0
			glyphStr := ""

			cssStr += fmt.Sprintf("  margin: 0 0 0 %.0fpt;\n", idFl)
			cssStr += fmt.Sprintf("  padding-left: %.0fpt;\n", idSt-idFl - 6.0)
			cssStr += fmt.Sprintf("}\n")

			cumIndent += idSt

			// Css <li nest level>
			cssStr += fmt.Sprintf(".%s_li.nL_%d {\n", listClass, nl)
			switch dObj.docLists[i].ord {
				case true:
					cssStr += fmt.Sprintf("  counter-increment: %s_li_nL_%d;\n", listClass, nl)
//					cssStr += fmt.Sprintf("list-style-type: %s;\n", )
				case false:
					cssStr += fmt.Sprintf("  list-style-type: %s;\n", glyphStr)
//					cssStr += fmt.Sprintf dObj.cvtGlyph(nestLev)
			}
			cssStr += fmt.Sprintf("}\n")

			// Css marker
			cssStr += fmt.Sprintf(".%s_li.nL_%d::marker {\n", listClass, nl)
			switch dObj.docLists[i].ord {
				case true:
					cssStr += fmt.Sprintf(" content: counter(%s_li_nL_%d, %s) \".\";", listClass, nl, glyphStr)
				case false:

			}

//            cssStr += cvtTxtMapStylCss(defTxtMap,nestLev.TextStyle)
			cssStr += fmt.Sprintf("}\n")
		}
	}
	headCss += cssStr

   // css default table
    if dObj.tableCount > 0 {

       //css default table styling (center aligned)
        cssStr = fmt.Sprintf(".%s_tbl {\n", dObj.docName)
        cssStr += "  width: 100%;\n"
        cssStr += "  border-collapse: collapse;\n"
        cssStr += "  border: 1px solid black;\n"
        cssStr += "  margin-left: auto;  margin-right: auto;\n"
        cssStr += "}\n"

        //css table cell
        cssStr = fmt.Sprintf(".%s_tblcell {\n", dObj.docName)
//      cssStr += "  border-collapse: collapse;\n"
        cssStr += "  border: 1px solid black;\n"
//      cssStr += "  margin:auto;\n"
        cssStr += "  padding: 0.5pt;\n"
        cssStr += "}\n"

    }

    headCss += cssStr
	return headCss, nil
}

//footnote div
func (dObj *mdDomObj) creFtnoteDivDom () (ftnoteDiv *dispObj, err error) {
	var ftnDiv dispObj
	var cssStr, scriptStr string
	var jselObj elScriptObj

	//html div footnote
	// fmt.Sprintf("<!-- Footnotes: %d -->\n", len(dObj.docFtnotes))
	// fmt.Sprintf("<div class=\"%s_main %s_ftndiv\">\n", dObj.docName, dObj.docName)

	// script
	scriptStr = fmt.Sprintf("// *** Footnotes: %d ***\n", len(dObj.docFtnotes))
	jselObj.parent = "divDoc"
	jselObj.typ = "div"
	jselObj.cl1 = dObj.docName + "_main"
	jselObj.cl2 = dObj.docName + "_ftndiv"
	jselObj.newEl = "divFtn"
	scriptStr += addElToDom(jselObj)

	//css div footnote
	cssStr = fmt.Sprintf(".%s_main.%s_ftndiv  {\n", dObj.docName, dObj.docName)

	if dObj.Options.DivBorders {
		cssStr += "  border: solid green;\n"
		cssStr += "  border-width: 1px;\n"
	}
	cssStr += "  padding-top:10px;\n"
	cssStr += "  counter-reset:ftcounter;\n"
	cssStr += "}\n"

	//html footnote title
	// fmt.Sprintf("<p class=\"%s_title %s_ftTit\">Footnotes</p>\n", dObj.docName, dObj.docName, dObj.docName)
	
	//script

	jselObj.parent = "divFtn"
	jselObj.typ = "p"
	jselObj.cl1 = dObj.docName + "_title"
	jselObj.cl2 = dObj.docName + "_ftndiv"
	jselObj.newEl = "ft_title"
	jselObj.txt = "Footnotes"
	scriptStr += addElToDom(jselObj)

	//css footnote title
	cssStr += fmt.Sprintf("%s_title.%s_ftTit {\n", dObj.docName, dObj.docName, dObj.docName)
	cssStr += "  color: purple;\n"
	cssStr += "}\n"

	// list for footnotes

	//css list
	cssStr += fmt.Sprintf(".%s_ftnOl {\n", dObj.docName)
	cssStr += "  display:block;\n"
	cssStr += "  list-style-type: decimal;\n"
	cssStr += "  padding-inline-start: 10pt;\n"
	cssStr += "  margin: 0;\n"
	cssStr += "}\n"

	// html
	// fmt.Sprintf("<ol class=\"%s_ftnOl\">\n", dObj.docName)

	// script
	jselObj.parent = "divFtn"
	jselObj.typ = "ol"
	jselObj.cl1 = dObj.docName + "_ftnOL"
//	jselObj.cl2 = dObj.docName + "_ftndiv"
	jselObj.newEl = "ft_Ol"
//	jselObj.txt = "Footnotes"
	scriptStr += addElToDom(jselObj)

	// prefix for paragraphs

	// css
	cssStr += fmt.Sprintf(".%s_p.%s_pft {\n",dObj.docName, dObj.docName)
	cssStr += "text-indent: 10pt;"
	cssStr += "counter-increment:ftcounter;"
	cssStr += "}\n"
	cssStr += fmt.Sprintf(".%s_p.%s_pft::before {\n",dObj.docName, dObj.docName)
	cssStr += "counter(ftcounter) ' ';"
	cssStr += "}\n"
	ftnDiv.bodyCss = cssStr

	// footnotes paragraph html
	for iFtn:=0; iFtn<len(dObj.docFtnotes); iFtn++ {
	}

	return &ftnDiv, nil
}

//toc div
func (dObj *mdDomObj) creTocDivDom () (tocObj *dispObj, err error) {
	var tocDiv dispObj
	var cssStr, scriptStr string
	var elObj elScriptObj

	if dObj.Options.Toc != true { return nil, nil }

	if dObj.Options.Verb {
		if len(dObj.headings) < 2 {
			fmt.Printf("*** no TOC insufficient headings: %d ***\n", len(dObj.headings))
		} else {
			fmt.Printf("*** creating TOC Div ***\n")
		}
	}

	if len(dObj.headings) < 2 {
//		tocDiv.bodyHtml = fmt.Sprintf("<!-- no toc insufficient headings -->")
		tocObj.script = "// *** no TOC insufficient headings ***\n"
		return tocObj, nil
	}

	//css
	cssStr = ""
	for namStyl, val := range dObj.namStylMap {

		if !val {continue}

		switch namStyl {
			case "HEADING_1":
				cssStr = fmt.Sprintf(".%s_h1.toc_h1 {\n",dObj.docName)
 				cssStr += "  padding-left: 10px;\n  margin: 0px;\n"
				cssStr += "}\n"

			case "HEADING_2":
				cssStr = fmt.Sprintf(".%s_h2.toc_h2 {\n",dObj.docName)
				cssStr += " padding-left: 20px;\n  margin: 0px;\n"
				cssStr += "}\n"

			case "HEADING_3":
				cssStr = fmt.Sprintf(".%s_h3.toc_h3 {\n",dObj.docName)
				cssStr += " padding-left: 40px;\n  margin: 0px;\n"
				cssStr += "}\n"

			case "HEADING_4":
				cssStr = fmt.Sprintf(".%s_h4.toc_h4 {\n",dObj.docName)
				cssStr += " padding-left: 60px;\n  margin: 0px;\n"
				cssStr += "}\n"

			case "HEADING_5":
				cssStr = fmt.Sprintf(".%s_h5.toc_h5 {\n",dObj.docName)
				cssStr += " padding-left: 80px;\n  margin: 0px;\n"
				cssStr += "}\n"

			case "HEADING_6":
				cssStr = fmt.Sprintf(".%s_h6.toc_h6 {\n",dObj.docName)
				cssStr += " padding-left: 100px;\n  margin: 0px;\n"
				cssStr += "}\n"

			default:

		}
	}

	tocDiv.bodyCss = cssStr

	//html
	//fmt.Sprintf("<div class=\"%s_main %s_top\">\n", dObj.docName)
	elObj.parent = "divDoc"
	elObj.typ = "div"
	elObj.newEl = "divToc"
	elObj.cl1 = dObj.docName + "_main"
	elObj.cl2 = dObj.docName + "_top"
	tocDiv.script = addElToDom(elObj)

	//fmt.Sprintf("<p class=\"%s_title %s_leftTitle\">Table of Contents</p>\n", dObj.docName, dObj.docName)
	elObj.parent = "divToc"
	elObj.typ = "p"
	elObj.newEl = "divToc"
	elObj.cl1 = dObj.docName + "_title"
	elObj.cl2 = dObj.docName + "_leftTitle"
	elObj.txt = "Table of Contents"
	tocDiv.script = addElToDom(elObj)

	tocDiv.script = scriptStr

	//html all headings are entries to toc table of content
	for ihead:=0; ihead<len(dObj.headings); ihead++ {

		namedStyl := dObj.headings[ihead].namedStyl
		hdId := dObj.headings[ihead].id[3:]
		text := dObj.headings[ihead].text


		switch namedStyl {
		case "TITLE":
			//prefix := fmt.Sprintf("<p class=\"%s_title %s_leftTitle_UL\">", dObj.docName, dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\" class=\"%s_noUl\">%s</a>", hdId, dObj.docName, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "p"
			elObj.cl1 = dObj.docName + "_title"
			elObj.cl2 = dObj.docName + "_leftTitle"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "SUBTITLE":
			//prefix := fmt.Sprintf("<p class=\"%s_subtitle\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "p"
			elObj.cl1 = dObj.docName + "_subtitle"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_1":
			//html
			//prefix := fmt.Sprintf("<h1 class=\"%s_h1 toc_h1\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h1"
			elObj.cl1 = dObj.docName + "_h1"
			elObj.cl2 = "toc_h1"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_2":
			//prefix := fmt.Sprintf("<h2 class=\"%s_h2 toc_h2\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h2"
			elObj.cl1 = dObj.docName + "_h2"
			elObj.cl2 = "toc_h2"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_3":
			//prefix := fmt.Sprintf("<h3 class=\"%s_h3 toc_h3\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h3"
			elObj.cl1 = dObj.docName + "_h3"
			elObj.cl2 = "toc_h3"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_4":
			//prefix := fmt.Sprintf("<h4 class=\"%s_h4 toc_h4\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h4"
			elObj.cl1 = dObj.docName + "_h4"
			elObj.cl2 = "toc_h4"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_5":
			//prefix := fmt.Sprintf("<h5 class=\"%s_h5 toc_h5\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h5"
			elObj.cl1 = dObj.docName + "_h5"
			elObj.cl2 = "toc_h5"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "HEADING_6":
			//prefix := fmt.Sprintf("<h6 class=\"%s_h6 toc_h6\">", dObj.docName)
			//middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)

			//script
			elObj.parent = "divToc"
			elObj.typ = "h6"
			elObj.cl1 = dObj.docName + "_h6"
			elObj.cl2 = "toc_h6"
			elObj.newEl = "parel"
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "NORMAL_TEXT":

		default:

		}

	} // end loop

	return &tocDiv, nil
}

func (dObj *mdDomObj) cvtBodyToDom() (bodyObj *dispObj, err error) {

	if dObj == nil {
		return nil, fmt.Errorf("-- no GdocObj!")
	}

	bodyObj = new(dispObj)

//	bodyObj.bodyHtml = fmt.Sprintf("<div class=\"%s_main\">\n", dObj.docName)
	var divMain elScriptObj
	divMain.comment = "create main div"
	divMain.typ = "div"
	divMain.parent = "divDoc"
	divMain.cl1 = dObj.docName + "_main"
	dObj.parent = "divMain"
	divMain.newEl = dObj.parent
	bodyObj.script = addElToDom(divMain)

	return bodyObj, err
}


func CvtMdToDomSections(heading, folderPath string, doc string, options *util.OptObj)(err error) {
// function that creates an html fil from the named section

	var mainDiv dispObj
	var dObj mdDomObj

	// initialize dObj with doc assignment
	err = dObj.initMdDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.creFtnoteDivDom()
	if err != nil {
		fmt.Errorf("creFtnoteDivDom: %v", err)
	}

//	dObj.sections
	secDiv := dObj.creSecDivDom()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
//			pgHd := dObj.creSecHeadToDom(ipage)
//			elStart := dObj.sections[ipage].secElStart
//			elEnd := dObj.sections[ipage].secElEnd
//			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
//			mainDiv.headCss += pgBody.headCss
//			mainDiv.bodyCss += pgBody.bodyCss
//			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
		}
	} else {
		mBody, err := dObj.cvtBodyToDom()
		if err != nil {
			return fmt.Errorf("cvtBody: %v", err)
		}
//		mainDiv.headCss += mBody.headCss
		mainDiv.bodyCss += mBody.bodyCss
		mainDiv.bodyHtml += mBody.bodyHtml
	}

	//css for document head
	headCss, err := dObj.creCssDocHead()
	if err != nil {
		return fmt.Errorf("creCssDocHead: %v", err)
	}

	//html + css for Toc Div
	tocDiv, err := dObj.creTocDivDom()
	if err != nil {
		tocDiv.bodyHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
	}

	//get html file pointer
	outfil := dObj.htmlFil
	if outfil == nil {
		return fmt.Errorf("outfil is nil!")
	}

	// assemble html document
	// html document file header
	docHeadStr := creHtmlDocHead(dObj.docName)
	outfil.WriteString(docHeadStr)

	//Css

	//css default css of document and document dimensions
	outfil.WriteString(headCss)

	// css of body elements
	outfil.WriteString(mainDiv.bodyCss)

	//css footnotes
	if ftnoteDiv != nil {
		cssStr := creFtnoteCss(dObj.docName)
		cssStr += ftnoteDiv.bodyCss
		outfil.WriteString(cssStr)
	}

	//css toc
	if tocDiv != nil {
		cssStr := creTocCss(dObj.docName)
		cssStr  += tocDiv.bodyCss
		outfil.WriteString(cssStr)
	}

	//css sec
	if secDiv != nil {
		cssStr := creSecCss(dObj.docName)
		if tocDiv == nil { cssStr += creTocSecCss(dObj.docName) }
		outfil.WriteString(cssStr)
	}

	// html start body
	outfil.WriteString("</style>\n</head>\n<body>\n")


	// html doc div
	htmlStr := creHtmlDocDiv(dObj.docName)
	outfil.WriteString(htmlStr)

	// html toc
	if tocDiv != nil  {outfil.WriteString(tocDiv.bodyHtml)}

	// html sections
	if secDiv != nil {outfil.WriteString(secDiv.bodyHtml)}

	// html main document
	outfil.WriteString(mainDiv.bodyHtml)

	// html footnotes
	if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyHtml)}

	// html ends doc div
	outfil.WriteString("</div>\n</body>\n</html>\n")
	outfil.Close()
	return nil
}



func CvtMdToDomAll(folderPath string, doc *os.File, options *util.OptObj)(err error) {
// function that creates an html fil from the named section
	var mainDiv dispObj
	var dObj mdDomObj

	// initialize dObj with doc assignment
	err = dObj.initMdDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.creFtnoteDivDom()
	if err != nil {
		fmt.Errorf("creFtnoteDivDom: %v", err)
	}

//	dObj.sections
	secDiv := dObj.creSecDivDom()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
//			pgHd := dObj.creSecHeadToDom(ipage)
//			elStart := dObj.sections[ipage].secElStart
//			elEnd := dObj.sections[ipage].secElEnd
//			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
//			mainDiv.bodyCss += pgBody.bodyCss
//			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
		}
	} else {
		mBody, err := dObj.cvtBodyToDom()
		if err != nil {
			return fmt.Errorf("cvtBody: %v", err)
		}
		mainDiv.bodyCss = mBody.bodyCss
		mainDiv.script = mBody.script
	}

	//css for document head
	headCss, err := dObj.creCssDocHead()
	if err != nil {
		return fmt.Errorf("creCssDocHead: %v", err)
	}

	//html + css for Toc Div
	tocDiv, err := dObj.creTocDivDom()
	if err != nil {
		tocDiv.bodyHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
	}

	//get html file pointer
	outfil := dObj.htmlFil
	if outfil == nil {
		return fmt.Errorf("outfil is nil!")
	}

	// assemble html document
	// html document file header
	docHeadStr := creHtmlDocHead(dObj.docName)
	outfil.WriteString(docHeadStr)

	//Css

	//css default css of document and document dimensions
	outfil.WriteString(headCss)

	// css of body elements
	outfil.WriteString(mainDiv.bodyCss)

	//css footnotes
	if ftnoteDiv != nil {
		cssStr := creFtnoteCss(dObj.docName)
		cssStr += ftnoteDiv.bodyCss
		outfil.WriteString(cssStr)
	}

	//css toc
	if tocDiv != nil {
		cssStr := creTocCss(dObj.docName)
		cssStr  += tocDiv.bodyCss
		outfil.WriteString(cssStr)
	}

	//css sec
	if secDiv != nil {
		cssStr := creSecCss(dObj.docName)
		if tocDiv == nil { cssStr += creTocSecCss(dObj.docName) }
		outfil.WriteString(cssStr)
	}

	// css end
	cssStr := "</style>\n"
	outfil.WriteString(cssStr)

	//script start
	jsStr := "<script>\n"
	outfil.WriteString(jsStr)

	//js create doc div
	imgfun := false
	if (dObj.inImgCount + dObj.posImgCount) > 0 {imgfun = true}
	tablefun := false
	if dObj.tableCount > 0 {tablefun = true}
	jsStr = creElFuncScript(imgfun, tablefun)
	outfil.WriteString(jsStr)

	jsStr = mainDiv.script
	outfil.WriteString(jsStr)

	jsStr = creDocDivScript(dObj.docName)
	outfil.WriteString(jsStr)

	//script end
	jsStr = "</script>\n"
	outfil.WriteString(jsStr)


	// html start body
	htmlStr := "<body>\n"
	outfil.WriteString(htmlStr)


	// html doc div
//	htmlStr = creHtmlDocDiv(dObj.docName)
//	outfil.WriteString(htmlStr)

	// html toc
	if tocDiv != nil  {outfil.WriteString(tocDiv.bodyHtml)}

	// html sections
	if secDiv != nil {outfil.WriteString(secDiv.bodyHtml)}

	// html main document
//	outfil.WriteString(mainDiv.bodyHtml)

	// html footnotes
	if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyHtml)}

	// html ends doc div
	htmlStr = creHtmlDocEnd()
	outfil.WriteString(htmlStr)
//	outfil.WriteString("</div>\n</body>\n</html>\n")
	outfil.Close()
	return nil
}


