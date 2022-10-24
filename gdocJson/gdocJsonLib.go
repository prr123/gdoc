// golang library that creates a json file from a gdoc file
// adopted from gdocDomLib.go
// author: prr, azul software
// created: 10/10/2022
// copyright 2022 prr, Peter Riemenschneider
//
// for changes see github
//
// start: CvtGdocToJson
//
// fix cssRules
//

package gdocJson

import (
	"fmt"
	"os"
	"google.golang.org/api/docs/v1"
    util "google/gdoc/gdocUtil"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)

type GdocDomObj struct {
	doc *docs.Document
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
	counter string
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
	errCount int
	htmlFil *os.File
	jsonFil *os.File
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
	doAppend bool
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

type colScriptObj struct {
	cl1 string
	cl2 string
	idStr string
	parent string
	newEl string
	comment string
	spanCount int
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

func printLiStack(stack *[]cList, title string) {
var item cList
var	n int
	if stack == nil {
		fmt.Printf("*** no listStack: %s ***\n", title)
		return
	}
	n = len(*stack) -1
	if n>=0 {
		item = (*stack)[n]
	} else {
		n = -1
	}
	fmt.Printf("list stack %s: Nestlev: %d", title, n)
	if n >= 0 {fmt.Printf(" id: %s ordered: %t", item.cListId, item.cOrd)}
	fmt.Printf("\n")
	return
}

func printLiStackItem(listAtt cList, cNest int){
		fmt.Printf("\nNest Lev: %d", cNest)
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

func fillTxtMap (txtStyl *docs.TextStyle)(txtMapRef *textMap) {
    var txtMap textMap

    if txtStyl == nil { return nil}

    txtMap.baseOffset = "NONE"
    if len(txtStyl.BaselineOffset) >0 {
        if txtStyl.BaselineOffset != "BASELINE_OFFSET_UNSPECIFIED" {
            txtMap.baseOffset = txtStyl.BaselineOffset
        }
    }

    txtMap.fontWeight = 400
    if txtStyl.Bold {
        txtMap.fontWeight = 800
    }

    txtStyl.Italic = false
    if txtStyl.Italic {
        txtMap.italic = txtStyl.Italic
    }

    txtStyl.Underline = false
    if txtStyl.Underline {
        txtMap.underline = txtStyl.Underline
    }

    txtMap.strike = false
    if txtStyl.Strikethrough {
        txtMap.strike = txtStyl.Strikethrough
    }

    txtMap.fontType = "Calibri"
    if txtStyl.WeightedFontFamily != nil {
        if txtStyl.WeightedFontFamily.FontFamily != txtMap.fontType {
            txtMap.fontType = txtStyl.WeightedFontFamily.FontFamily
        }
        if txtStyl.WeightedFontFamily.Weight > 0 {
            if txtStyl.WeightedFontFamily.Weight != txtMap.fontWeight {
                txtMap.fontWeight = txtStyl.WeightedFontFamily.Weight
            }
        }
    }

    txtMap.fontSize = 0.0
    if txtStyl.FontSize != nil {
        if txtStyl.FontSize.Magnitude >0 {
            txtMap.fontSize = txtStyl.FontSize.Magnitude
        }
    }

    txtMap.txtColor = "rgb(0,0,0)"
    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            color := util.GetColor(txtStyl.ForegroundColor.Color)
            if color != txtMap.txtColor {
                txtMap.txtColor = color
            }
        }
    }

    txtMap.bckColor = ""
    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            color := util.GetColor(txtStyl.BackgroundColor.Color)
            if color != txtMap.bckColor {
                txtMap.bckColor = color
            }
        }
    }
    return &txtMap
}

func cvtTxtMapJson(txtMap *textMap)(cssStr string) {
// for css rule only
    cssStr =""
    if len(txtMap.baseOffset) > 0 {
        switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += " verticalAlign: sub;"
            case "SUBSCRIPT":
                cssStr += " verticalAlign: sup;"
            case "NONE":
                cssStr += " verticalAlign: baseline;"
            default:
            //error
                cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
        }
    }

    if txtMap.italic {
        cssStr += " fontStyle: italic;"
    } else {
        cssStr += " fontStyle: normal;"
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
    cssStr += fmt.Sprintf("textDecoration: %s;", textprop)

    if len(txtMap.fontType) >0 { cssStr += fmt.Sprintf("fontFamily: %s;", txtMap.fontType)}
    if txtMap.fontWeight > 0 {cssStr += fmt.Sprintf("fontWeight: %d;", txtMap.fontWeight)}
    if txtMap.fontSize >0 {cssStr += fmt.Sprintf("fontSize: %.1fpt;", txtMap.fontSize)}
    if len(txtMap.txtColor) >0 {cssStr += fmt.Sprintf("color: %s;", txtMap.txtColor)}
    if len(txtMap.bckColor) >0 {cssStr += fmt.Sprintf("backgroundColor: %s;", txtMap.bckColor)}

//	if len(cssStr)> 0 {cssStr = cssStr[:len(cssStr) -1]}
    return cssStr
}

func cvtTxtMapToCssJson(txtMap *textMap)(cssStr string) {
// for css rule only
    cssStr =""
    if len(txtMap.baseOffset) > 0 {
        switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += " vertical-align: sub;"
            case "SUBSCRIPT":
                cssStr += " vertical-align: sup;"
            case "NONE":
                cssStr += " vertical-align: baseline;"
            default:
            //error
                cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
        }
    }

    if txtMap.italic {
        cssStr += " fontStyle: italic;"
    } else {
        cssStr += " fontStyle: normal;"
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
    cssStr += fmt.Sprintf("text-decoration: %s;", textprop)

    if len(txtMap.fontType) >0 { cssStr += fmt.Sprintf("font-family: %s;", txtMap.fontType)}
    if txtMap.fontWeight > 0 {cssStr += fmt.Sprintf("font-weight: %d;", txtMap.fontWeight)}
    if txtMap.fontSize >0 {cssStr += fmt.Sprintf("font-size: %.1fpt;", txtMap.fontSize)}
    if len(txtMap.txtColor) >0 {cssStr += fmt.Sprintf("color: %s;", txtMap.txtColor)}
    if len(txtMap.bckColor) >0 {cssStr += fmt.Sprintf("background-color: %s;", txtMap.bckColor)}

//	if len(cssStr)> 0 {cssStr = cssStr[:len(cssStr) -1]}
    return cssStr
}

func cvtTxtMapStylToCssJson (txtMap *textMap, txtStyl *docs.TextStyle)(cssStr string) {

    if (len(txtStyl.BaselineOffset) > 0) && (txtStyl.BaselineOffset != "BASELINE_OFFSET_UNSPECIFIED") {        if txtStyl.BaselineOffset != txtMap.baseOffset {
            txtMap.baseOffset = txtStyl.BaselineOffset
            switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += "vertical-align: sub;"
            case "SUBSCRIPT":
                cssStr += "vertical-align: sup;"
            case "NONE":
                cssStr += "vertical-align: baseline;"
            default:
                cssStr += fmt.Sprintf("vertical-align: unknown %s;", txtMap.baseOffset)
            }
        }
    }

    switch {
    case txtStyl.Bold && (txtMap.fontWeight < 700):
        txtMap.fontWeight = 800
        cssStr += fmt.Sprintf("font-weight: %d;", txtMap.fontWeight)
    case !txtStyl.Bold && (txtMap.fontWeight > 500):
        txtMap.fontWeight = 400
        cssStr += fmt.Sprintf("font-weight: %d;", txtMap.fontWeight)
    default:

    }

    if txtStyl.Italic != txtMap.italic {
        txtMap.italic = txtStyl.Italic
        if txtMap.italic {
            cssStr += "font-style: italic;"
        } else {
            cssStr += "font-style: normal;"
        }
    }

    txtprop := ""

    if txtStyl.Underline != txtMap.underline {
        txtMap.underline = txtStyl.Underline
        if txtMap.underline {
            txtprop = "underline"
        } else {
            txtprop = "none"
        }
    }
//  if txtMap.underline { cssStr += "  text-decoration: underline;\n"}

    if txtStyl.Strikethrough != txtMap.strike {
        txtMap.strike = txtStyl.Strikethrough
        if txtMap.strike {
            if txtprop == "none" {
                txtprop = "line-through"
            } else {
                txtprop += " line-through"
            }
        }
    }

    if len(txtprop) > 0 {cssStr += fmt.Sprintf("text-decoration: %s;", txtprop)}

    if txtStyl.WeightedFontFamily != nil {
        if txtStyl.WeightedFontFamily.FontFamily != txtMap.fontType {
            txtMap.fontType = txtStyl.WeightedFontFamily.FontFamily
            cssStr += fmt.Sprintf("font-family: %s;", txtMap.fontType)
        }
/*
        if txtStyl.WeightedFontFamily.Weight != txtMap.fontWeight {
            txtMap.fontWeight = txtStyl.WeightedFontFamily.Weight
            alter = true
        }
*/
    }


    if txtStyl.FontSize != nil {
        if txtStyl.FontSize.Magnitude != txtMap.fontSize {
            txtMap.fontSize = txtStyl.FontSize.Magnitude
            cssStr += fmt.Sprintf("font-size: %.1fpt;", txtMap.fontSize)
        }
    }

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            color := util.GetColor(txtStyl.ForegroundColor.Color)
            if color != txtMap.txtColor {
                txtMap.txtColor = color
                cssStr += fmt.Sprintf("color: %s;", txtMap.txtColor)
            }
        }
    }

    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            color := util.GetColor(txtStyl.BackgroundColor.Color)
            if color != txtMap.bckColor {
                txtMap.bckColor = color
                cssStr += fmt.Sprintf("background-color: %s;", txtMap.bckColor)
            }
        }
    }

	ilen := len(cssStr)
	if ilen > 0 {
	    return cssStr[:ilen-1]
	}
	return ""
}

func cvtTxtMapStylToJson (txtMap *textMap, txtStyl *docs.TextStyle)(cssStr string) {

    if (len(txtStyl.BaselineOffset) > 0) && (txtStyl.BaselineOffset != "BASELINE_OFFSET_UNSPECIFIED") {
        if txtStyl.BaselineOffset != txtMap.baseOffset {
            txtMap.baseOffset = txtStyl.BaselineOffset
            switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += "\"verticalAlign\": \"sub\","
            case "SUBSCRIPT":
                cssStr += "\"verticalAlign\": \"sup\","
            case "NONE":
                cssStr += "\"verticalAlign\": \"baseline\","
            default:
            //error
                cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
            }
        }
    }

    switch {
    case txtStyl.Bold && (txtMap.fontWeight < 700):
        txtMap.fontWeight = 800
        cssStr += fmt.Sprintf("\"fontWeight\": \"%d\",", txtMap.fontWeight)
    case !txtStyl.Bold && (txtMap.fontWeight > 500):
        txtMap.fontWeight = 400
        cssStr += fmt.Sprintf("\"fontWeight\": \"%d\",", txtMap.fontWeight)
    default:

    }

    if txtStyl.Italic != txtMap.italic {
        txtMap.italic = txtStyl.Italic
        if txtMap.italic {
            cssStr += "\"fontStyle\": \"italic\","
        } else {
            cssStr += "\"fontStyle\": \"normal\","
        }
    }

    txtprop := ""

    if txtStyl.Underline != txtMap.underline {
        txtMap.underline = txtStyl.Underline
        if txtMap.underline {
            txtprop = "underline"
        } else {
            txtprop = "none"
        }
    }
//  if txtMap.underline { cssStr += "  text-decoration: underline;\n"}

    if txtStyl.Strikethrough != txtMap.strike {
        txtMap.strike = txtStyl.Strikethrough
        if txtMap.strike {
            if txtprop == "none" {
                txtprop = "line-through"
            } else {
                txtprop += " line-through"
            }
        }
    }

    if len(txtprop) > 0 {cssStr += fmt.Sprintf("\"textDecoration\": \"%s\",", txtprop)}

    if txtStyl.WeightedFontFamily != nil {
        if txtStyl.WeightedFontFamily.FontFamily != txtMap.fontType {
            txtMap.fontType = txtStyl.WeightedFontFamily.FontFamily
            cssStr += fmt.Sprintf("\"fontFamily\": \"%s\",", txtMap.fontType)
        }
/*
        if txtStyl.WeightedFontFamily.Weight != txtMap.fontWeight {
            txtMap.fontWeight = txtStyl.WeightedFontFamily.Weight
            alter = true
        }
*/
    }


    if txtStyl.FontSize != nil {
        if txtStyl.FontSize.Magnitude != txtMap.fontSize {
            txtMap.fontSize = txtStyl.FontSize.Magnitude
            cssStr += fmt.Sprintf("\"fontSize\": \"%.1fpt\",", txtMap.fontSize)
        }
    }

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            color := util.GetColor(txtStyl.ForegroundColor.Color)
            if color != txtMap.txtColor {
                txtMap.txtColor = color
                cssStr += fmt.Sprintf("\"color\": \"%s\",", txtMap.txtColor)
            }
        }
    }

    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            color := util.GetColor(txtStyl.BackgroundColor.Color)
            if color != txtMap.bckColor {
                txtMap.bckColor = color
                cssStr += fmt.Sprintf("\"backgroundColor\": \"%s\",", txtMap.bckColor)
            }
        }
    }

	ilen := len(cssStr)
	if ilen > 0 {
	    return cssStr[:ilen-1]
	}
	return ""
}


func cvtTxtStylJson (txtStyl *docs.TextStyle)(tcssStr string) {

    if len(txtStyl.BaselineOffset) > 0 {
        valStr := "\"verticalAlign\": "
        switch txtStyl.BaselineOffset {
            case "SUPERSCRIPT":
                valStr += "\"sub\""
            case "SUBSCRIPT":
                valStr += "\"sup\""
            case "NONE":
                valStr += "\"baseline\""
            default:
                valStr = fmt.Sprintf("// Baseline Offset unknown: %s \n", txtStyl.BaselineOffset)
        }
        tcssStr = valStr + ","
    }

    if txtStyl.Bold {
        tcssStr += "\"fontWeight\": \"800\","
    } else {
        tcssStr += "\"fontWeight\": \"400\","
    }

    if txtStyl.Italic { tcssStr += "\"fontStyle\": \"italic\","}

	txtprop1 := ""
	txtprop := ""
    if txtStyl.Underline {
		txtprop1 = "underline"
	} else {
		txtprop1 = "none"
    }

    if txtStyl.Strikethrough {
		if txtprop1 == "none" {
			txtprop = "line-through"
		} else {
			txtprop = txtprop1 + " line-through"
        }
    }

    if len(txtprop) > 0 {tcssStr += fmt.Sprintf("\"textDecoration\": \"%s\",", txtprop)}

    if txtStyl.WeightedFontFamily != nil {
        font := txtStyl.WeightedFontFamily.FontFamily
        tcssStr += fmt.Sprintf("\"fontFamily\": \"%s\",", font)
    }

    if txtStyl.FontSize != nil {
        mag := txtStyl.FontSize.Magnitude
        tcssStr += fmt.Sprintf("\"fontSize\": \"%.1fpt\",", mag)
    }

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            //0 to 1
            tcssStr += "\"color\": "
            tcssStr += "\"" + util.GetColor(txtStyl.ForegroundColor.Color) + "\","
        }
    }
    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            tcssStr += "\"backgroundColor\": "
            tcssStr += "\"" + util.GetColor(txtStyl.BackgroundColor.Color) + "\","
        }
    }

    if len(tcssStr) > 0 {
        tcssStr = tcssStr[:len(tcssStr)-1]
    }
    return tcssStr
}

func printParMap(parmap *parMap, parStyl *docs.ParagraphStyle) {

	alter := false
	fmt.Printf("*** align ***\n")
	halign:= parmap.halign
	if parStyl.Alignment != parmap.halign {
		fmt.Printf("align: %s %s ->", parmap.halign, parStyl.Alignment)
		halign = parStyl.Alignment
		alter = true
	}
	fmt.Printf("align: %s \n", halign)

//	parmap.direct = true
	fmt.Printf("*** indent ***\n")
	indStart := parmap.indStart
	if (parStyl.IndentStart != nil) {
		if parStyl.IndentStart.Magnitude != parmap.indStart {
			fmt.Printf("indent start: %.1fpt %.1fpt -> ", parmap.indStart, parStyl.IndentStart.Magnitude)
			indStart = parStyl.IndentStart.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent start: %.1fpt\n", indStart)

	fmt.Printf("*** indent end***\n")
	indEnd := parmap.indEnd
	if (parStyl.IndentEnd != nil) {
		if parStyl.IndentEnd.Magnitude != parmap.indEnd {
			fmt.Printf("indent end: %.1f %.1f -> ", parmap.indEnd, parStyl.IndentEnd.Magnitude)
			indEnd = parStyl.IndentEnd.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent end: %.1fpt\n", indEnd)

	fmt.Printf("*** indent first line ***\n")
	indFlin := parmap.indFlin
	if (parStyl.IndentFirstLine != nil) {
		if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
			fmt.Printf("indent first line: %.1f %.1f -> ", parmap.indFlin, parStyl.IndentFirstLine.Magnitude)
			indFlin = parStyl.IndentFirstLine.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent first line: %.1fpt\n", indFlin)

	fmt.Printf("*** line spacing ***\n")
	linSpac := parmap.linSpac
	if parStyl.LineSpacing/100.0 != parmap.linSpac {
		fmt.Printf("line spacing: %.2f %.2f -> ", parmap.linSpac, parStyl.LineSpacing/100.0)
		linSpac = parStyl.LineSpacing/100.0; alter = true;
	}
	fmt.Printf("line spacing: %.2fpt\n", linSpac)

	fmt.Printf("*** keep lines ***\n")
	keepLines := parmap.keepLines
	if parStyl.KeepLinesTogether != parmap.keepLines {
		fmt.Printf("keep Lines: %t %t -> ", parmap.keepLines, parStyl.KeepLinesTogether)
		keepLines = parStyl.KeepLinesTogether; alter = true;
	}
	fmt.Printf("keep Lines: %t\n", keepLines)

	fmt.Printf("*** keep next ***\n")
	keepNext := parmap.keepNext
	if parStyl.KeepWithNext != parmap.keepNext {
		fmt.Printf("keep With: %t %t -> ", parmap.keepNext, parStyl.KeepWithNext)
		keepNext = parStyl.KeepWithNext; alter = true;
	}
	fmt.Printf("keep With: %t\n", keepNext)

	fmt.Printf("*** space above ***\n")
	spaceTop := parmap.spaceTop
	if (parStyl.SpaceAbove != nil) {
		if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
			fmt.Printf("space above: %.1fpt %.1fpt -> ", parmap.spaceTop, parStyl.SpaceAbove.Magnitude)
			spaceTop = parStyl.SpaceAbove.Magnitude
			alter = true
		}
	}
	fmt.Printf("space above: %.1fpt\n", spaceTop)

	fmt.Printf("*** space below ***\n")
	spaceBelow := parmap.spaceBelow
	if (parStyl.SpaceBelow != nil) {
		if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
			fmt.Printf("space below: %.1f %.1f -> ", parmap.spaceBelow, parStyl.SpaceBelow.Magnitude)
			spaceBelow = parStyl.SpaceBelow.Magnitude
			alter = true
		}
	}
	fmt.Printf("space below: %.1fpt\n", spaceBelow)

	fmt.Printf("*** space mode ***\n")

    spaceMode := true
    switch parStyl.SpacingMode {
    case "NEVER_COLLAPSE":
        spaceMode = true
    case "COLLAPSE_LISTS":
        spaceMode = false
    default:
        spaceMode = true
    }

    if spaceMode != parmap.spaceMode {
       fmt.Printf("spacing mode: %t %t -> ", parmap.spaceMode, spaceMode)
	}
	fmt.Printf("spacing mode: %s\n", spaceMode)

	//tabs to do
//	parmap.hasBorders = true

	bb := true
	bb = bb && (parStyl.BorderBetween == nil)
	bb = bb && (parStyl.BorderTop == nil)
	bb = bb && (parStyl.BorderRight == nil)
	bb = bb && (parStyl.BorderBottom == nil)
	bb = bb && (parStyl.BorderLeft == nil)
	if bb {
		fmt.Printf("has no borders: %t %t -> ", parmap.hasBorders, !bb)
		fmt.Printf("mo borders: %t\n", false)
		return
	}

	fmt.Println("\n*** has Borders! ***\n")
	parmap.hasBorders = true

	alter = false
	bordalter := false
//	fmt.Printf("*** borders between ***\n")
	
	if parStyl.BorderBetween != nil {
		if parStyl.BorderBetween.Width != nil {
			if parStyl.BorderBetween.Width.Magnitude != parmap.bordBet.width {
				bordBet := parStyl.BorderBetween.Width.Magnitude
				alter = true
				fmt.Printf("width: %.1f\n",bordBet)
			}
		}
		if parStyl.BorderBetween.Padding != nil {
			if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
				pad := parStyl.BorderBetween.Padding.Magnitude
				alter = true
				fmt.Printf("padding: %.1f\n",pad)
			}
		}
		if parStyl.BorderBetween.Color != nil {
			if parStyl.BorderBetween.Color.Color != nil {
				color := util.GetColor(parStyl.BorderBetween.Color.Color)
				if color != parmap.bordBet.color {
					dispcolor := color
					alter = true
					fmt.Printf("color: %s\n",dispcolor)
				}
			}
		}

		if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {
			dash := parStyl.BorderBetween.DashStyle;
			alter = true;
			fmt.Printf("dash: %s\n", dash)
		}
	}
	fmt.Printf("*** border between alter: %t ***\n", alter)
	if alter {bordalter = true}

	alter = false
//	fmt.Printf("*** border top ***\n")
	if parStyl.BorderTop != nil {
		if parStyl.BorderTop.Width != nil {
			if parStyl.BorderTop.Width.Magnitude != parmap.bordTop.width {
				parmap.bordTop.width = parStyl.BorderTop.Width.Magnitude
				alter = true
				fmt.Printf("width: %.1f\n",parmap.bordTop.width)
			}
		}
		if parStyl.BorderTop.Padding != nil {
			if parStyl.BorderTop.Padding.Magnitude != parmap.bordTop.pad {
				parmap.bordTop.pad = parStyl.BorderTop.Padding.Magnitude
				alter = true
			}
		}
		if parStyl.BorderTop.Color != nil {
			if parStyl.BorderTop.Color.Color != nil {
				color := util.GetColor(parStyl.BorderTop.Color.Color)
				if color != parmap.bordTop.color {
					parmap.bordTop.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderTop.DashStyle != parmap.bordTop.dash {parmap.bordTop.dash = parStyl.BorderTop.DashStyle; alter = true;}
	}
	fmt.Printf("*** border top alter: %t ***\n", alter)
	if alter {bordalter = true}

	alter = false
//	fmt.Printf("*** border right ***\n")
	if parStyl.BorderRight != nil {
		if parStyl.BorderRight.Width != nil {
			if parStyl.BorderRight.Width.Magnitude != parmap.bordRight.width {
				parmap.bordRight.width = parStyl.BorderRight.Width.Magnitude
				alter = true
			}
		}
		if parStyl.BorderRight.Padding != nil {
			if parStyl.BorderRight.Padding.Magnitude != parmap.bordRight.pad {
				parmap.bordRight.pad = parStyl.BorderRight.Padding.Magnitude
				alter = true
			}
		}
		if parStyl.BorderRight.Color != nil {
			if parStyl.BorderRight.Color.Color != nil {
				color := util.GetColor(parStyl.BorderRight.Color.Color)
				if color != parmap.bordRight.color {
					parmap.bordRight.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderRight.DashStyle != parmap.bordRight.dash {
			parmap.bordRight.dash = parStyl.BorderRight.DashStyle
			alter = true
		}
	}
	fmt.Printf("*** border right alter: %t ***\n", alter)
	if alter {bordalter = true}

	alter = false
//	fmt.Printf("*** border bottom ***\n")
	if parStyl.BorderBottom != nil {
		if parStyl.BorderBottom.Width != nil {
			if parStyl.BorderBottom.Width.Magnitude != parmap.bordBot.width {
				parmap.bordBot.width = parStyl.BorderBottom.Width.Magnitude
				alter = true
			}
		}
		if parStyl.BorderBottom.Padding != nil {
			if parStyl.BorderBottom.Padding.Magnitude != parmap.bordBot.pad {
				parmap.bordBot.pad = parStyl.BorderBottom.Padding.Magnitude
				alter = true
			}
		}
		if parStyl.BorderBottom.Color != nil {
			if parStyl.BorderBottom.Color.Color != nil {
				color := util.GetColor(parStyl.BorderBottom.Color.Color)
				if color != parmap.bordBot.color {
					parmap.bordBot.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderBottom.DashStyle != parmap.bordBot.dash {parmap.bordBot.dash = parStyl.BorderBottom.DashStyle; alter = true;}
	}
	fmt.Printf("*** border bottom alter: %t ***\n", alter)
	if alter {bordalter = true}

	alter = false
//	fmt.Printf("*** border left ***\n")
	if parStyl.BorderLeft != nil {
		if parStyl.BorderLeft.Width != nil {
			if parStyl.BorderLeft.Width.Magnitude != parmap.bordLeft.width {
				parmap.bordLeft.width = parStyl.BorderLeft.Width.Magnitude
				alter = true
			}
		}
		if parStyl.BorderLeft.Padding != nil {
			if parStyl.BorderLeft.Padding.Magnitude != parmap.bordLeft.pad {
				parmap.bordLeft.pad = parStyl.BorderLeft.Padding.Magnitude
				alter = true
			}
		}
		if parStyl.BorderLeft.Color != nil {
			if parStyl.BorderLeft.Color.Color != nil {
				color := util.GetColor(parStyl.BorderLeft.Color.Color)
				if color != parmap.bordLeft.color {
					parmap.bordLeft.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderLeft.DashStyle != parmap.bordLeft.dash {parmap.bordLeft.dash = parStyl.BorderLeft.DashStyle; alter = true;}
	}
	fmt.Printf("*** border left alter: %t ***\n", alter)
	if alter {bordalter = true}

	alter = bordalter

	bb2 := true
	bb2 = bb2 && (parmap.bordBet.width == 0.0)
	bb2 = bb2 && (parmap.bordTop.width == 0.0)
	bb2 = bb2 && (parmap.bordRight.width == 0.0)
	bb2 = bb2 && (parmap.bordBot.width == 0.0)
	bb2 = bb2 && (parmap.bordLeft.width == 0.0)

	if bb2 {parmap.hasBorders = false; alter = false;}

	fmt.Printf("alter borders: %t\n", alter)

	return
}

func fillParMap(parStyl *docs.ParagraphStyle)(parMapRef *parMap) {
// function that converts a parameter style object into a parMap Object
    var parmap parMap

    if parStyl == nil { return nil}

    parmap.halign = "START"
    if len(parStyl.Alignment) > 0 {
        parmap.halign = parStyl.Alignment
    }

    parmap.direct = true


    parmap.indStart = 0
    if (parStyl.IndentStart != nil) {
        if parStyl.IndentStart.Magnitude != parmap.indStart {
            parmap.indStart = parStyl.IndentStart.Magnitude
        }
    }

    parmap.indEnd = -1
    if (parStyl.IndentEnd != nil) {
        if parStyl.IndentEnd.Magnitude != parmap.indEnd {
            parmap.indEnd = parStyl.IndentEnd.Magnitude
        }
    }

    parmap.indFlin = 0
    if (parStyl.IndentFirstLine != nil) {
        if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
            parmap.indFlin = parStyl.IndentFirstLine.Magnitude
        }
    }

    parmap.linSpac = 1.0
    if parStyl.LineSpacing/100.0 != parmap.linSpac {
        if parStyl.LineSpacing > 1.0 {
            parmap.linSpac = parStyl.LineSpacing/100.0
        }
    }

    // may have to introduce an exemption for title
    parmap.keepLines = false
    if parStyl.KeepLinesTogether != parmap.keepLines {
        parmap.keepLines = parStyl.KeepLinesTogether
    }

    parmap.keepNext = false
    if parStyl.KeepWithNext != parmap.keepNext {
        parmap.keepNext = parStyl.KeepWithNext
    }

    parmap.spaceTop = 0
    if (parStyl.SpaceAbove != nil) {
        if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
            parmap.spaceTop = parStyl.SpaceAbove.Magnitude
        }
    }

    parmap.spaceBelow = 0
    if (parStyl.SpaceBelow != nil) {
        if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
            parmap.spaceBelow = parStyl.SpaceBelow.Magnitude
        }
    }


    switch parStyl.SpacingMode {
    case "NEVER_COLLAPSE":
        parmap.spaceMode = true
    case "COLLAPSE_LISTS":
        parmap.spaceMode = false
    default:
        parmap.spaceMode = true
    }


   // check for zero width borders
    bb := true
    if (parStyl.BorderBetween != nil) {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if (parStyl.BorderTop != nil) {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }


    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    parmap.hasBorders = true
    if !bb {
        parmap.hasBorders = false
    }

    if !parmap.hasBorders {
        return &parmap
    }

//  bordDisp := false
    parmap.bordBet.width = 0
    if parStyl.BorderBetween != nil {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.0 {
//              bordDisp = true
                if parStyl.BorderBetween.Width.Magnitude != parmap.bordBet.width {
                    parmap.bordBet.width = parStyl.BorderBetween.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBetween.Padding != nil {
            if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
                parmap.bordBet.pad = parStyl.BorderBetween.Padding.Magnitude
            }
        }
        if parStyl.BorderBetween.Color != nil {
            if parStyl.BorderBetween.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBetween.Color.Color)
                if color != parmap.bordBet.color {
                    parmap.bordBet.color = color
                }
            }
        }
        parmap.bordBet.dash = "SOLID"
        if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {parmap.bordBet.dash = parStyl.BorderBetween.DashStyle;}
    }

//  bordDisp = false
    parmap.bordTop.width = 0
    if parStyl.BorderTop != nil {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.0 {
//              bordDisp = true
                if parStyl.BorderTop.Width.Magnitude != parmap.bordTop.width {
                    parmap.bordTop.width = parStyl.BorderTop.Width.Magnitude
                }
            }
        }
        if parStyl.BorderTop.Padding != nil {
            if parStyl.BorderTop.Padding.Magnitude != parmap.bordTop.pad {
                parmap.bordTop.pad = parStyl.BorderTop.Padding.Magnitude
            }
        }
        if parStyl.BorderTop.Color != nil {
            if parStyl.BorderTop.Color.Color != nil {
                color := util.GetColor(parStyl.BorderTop.Color.Color)
                if color != parmap.bordTop.color {
                    parmap.bordTop.color = color
                }
            }
        }
        parmap.bordTop.dash = "SOLID"
        if parStyl.BorderTop.DashStyle != parmap.bordTop.dash {parmap.bordTop.dash = parStyl.BorderTop.DashStyle;}
    }

//  bordDisp = false
    parmap.bordRight.width = 0
    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.0 {
//              bordDisp = true
                if parStyl.BorderRight.Width.Magnitude != parmap.bordRight.width {
                    parmap.bordRight.width = parStyl.BorderRight.Width.Magnitude
                }
            }
        }
        if parStyl.BorderRight.Padding != nil {
            if parStyl.BorderRight.Padding.Magnitude != parmap.bordRight.pad {
                parmap.bordRight.pad = parStyl.BorderRight.Padding.Magnitude
            }
        }
        if parStyl.BorderRight.Color != nil {
            if parStyl.BorderRight.Color.Color != nil {
                color := util.GetColor(parStyl.BorderRight.Color.Color)
                if color != parmap.bordRight.color {
                    parmap.bordRight.color = color
                }
            }
        }
        parmap.bordRight.dash = "SOLID"
        if parStyl.BorderRight.DashStyle != parmap.bordRight.dash {
            parmap.bordRight.dash = parStyl.BorderRight.DashStyle
        }
    }

//  bordDisp = false
    parmap.bordBot.width = 0
    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.0 {
//              bordDisp = true
                if parStyl.BorderBottom.Width.Magnitude != parmap.bordBot.width {
                    parmap.bordBot.width = parStyl.BorderBottom.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBottom.Padding != nil {
            if parStyl.BorderBottom.Padding.Magnitude != parmap.bordBot.pad {
                parmap.bordBot.pad = parStyl.BorderBottom.Padding.Magnitude
            }
        }
        if parStyl.BorderBottom.Color != nil {
            if parStyl.BorderBottom.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBottom.Color.Color)
                if color != parmap.bordBot.color {
                    parmap.bordBot.color = color
                }
            }
        }
        parmap.bordBot.dash = "SOLID"
        if parStyl.BorderBottom.DashStyle != parmap.bordBot.dash {parmap.bordBot.dash = parStyl.BorderBottom.DashStyle;}
    }

//  bordDisp = false
    parmap.bordLeft.width = 0
    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.0 {
//              bordDisp = true
                if parStyl.BorderLeft.Width.Magnitude != parmap.bordLeft.width {
                    parmap.bordLeft.width = parStyl.BorderLeft.Width.Magnitude
                }
            }
        }
        if parStyl.BorderLeft.Padding != nil {
            if parStyl.BorderLeft.Padding.Magnitude != parmap.bordLeft.pad {
                parmap.bordLeft.pad = parStyl.BorderLeft.Padding.Magnitude
            }
        }
        if parStyl.BorderLeft.Color != nil {
            if parStyl.BorderLeft.Color.Color != nil {
                color := util.GetColor(parStyl.BorderLeft.Color.Color)
                if color != parmap.bordLeft.color {
                    parmap.bordLeft.color = color
                }
            }
        }
        parmap.bordLeft.dash = "SOLID"
        if parStyl.BorderLeft.DashStyle != parmap.bordLeft.dash {parmap.bordLeft.dash = parStyl.BorderLeft.DashStyle;}
    }

    bb2 := true
    bb2 = bb2 && (parmap.bordBet.width < 0.01)
    bb2 = bb2 && (parmap.bordTop.width < 0.01)
    bb2 = bb2 && (parmap.bordRight.width < 0.01)
    bb2 = bb2 && (parmap.bordBot.width < 0.01)
    bb2 = bb2 && (parmap.bordLeft.width < 0.01)

    if bb2 {parmap.hasBorders = false}

    return &parmap
}


func cvtParMapStylToJson(parmap *parMap, parStyl *docs.ParagraphStyle, opt *util.OptObj)(cssStr string, alter bool) {
// function that creates the css attributes of a paragraph
// the function compares the values of the parMap and parStyl
    if parmap == nil {return "/* no parmap */", false}
    if parStyl == nil { return "/* no parStyl */", false}

    cssStr =""

    if (len(parStyl.Alignment) > 0) &&  (parmap.halign != parStyl.Alignment) {
		alter = true
        halign := parStyl.Alignment
        switch halign {
            case "START":
                cssStr += "\"textAlign\": \"left\","
            case "CENTER":
                cssStr += "\"textAlign\": \"center\","
            case "END":
                cssStr += "\"textAlign\": \"right\","
            case "JUSTIFIED":
                cssStr += "\"textAlign\": \"justify\","
            default:
                cssStr += fmt.Sprintf("\"textAlign\": \"error: %s\",", parmap.halign)
        }

    }

    // test direction skip for now
    parmap.direct = true

    if (parStyl.IndentFirstLine != nil) {
        if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
            indFlin := parStyl.IndentFirstLine.Magnitude
            cssStr += fmt.Sprintf("\"textIndent\": \"%.1fpt\",", indFlin)
			alter = true
        }
    }

    if parmap.linSpac < 0.1 {parmap.linSpac = 1.0}
    if parStyl.LineSpacing/100.0 != parmap.linSpac {
        if parStyl.LineSpacing/100 > 1.0 {
			alter = true
            linSpac := parStyl.LineSpacing/100.0
            if opt.DefLinSpacing > 0.0 {
                cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", opt.DefLinSpacing*linSpac)
            } else {
                cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", linSpac)
            }
        } else {
			linSpac := parmap.linSpac
            if opt.DefLinSpacing > 0.0 {
                cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", opt.DefLinSpacing*linSpac)
            } else {
                cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", linSpac)
            }
		}
    }

    margin := false
    lmarg := 0.0
    if (parStyl.IndentStart != nil) {
        if parStyl.IndentStart.Magnitude != parmap.indStart {
            lmarg = parStyl.IndentStart.Magnitude
            margin = true
        }
    }

    rmarg := 0.0
    if (parStyl.IndentEnd != nil) {
        if parStyl.IndentEnd.Magnitude != parmap.indEnd {
            rmarg = parStyl.IndentEnd.Magnitude
            margin = true
        }
    }

    tmarg := 0.0
    if (parStyl.SpaceAbove != nil) {
        if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
	        tmarg = parStyl.SpaceAbove.Magnitude
            margin = true
        }
    }

    bmarg := 0.0
    if (parStyl.SpaceBelow != nil) {
        if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
            bmarg = parStyl.SpaceBelow.Magnitude
            margin = true
        }
    }

    if margin {
		cssStr += fmt.Sprintf("\"margin\": \"%.0f %.0f %.0f %.0f\",", tmarg, rmarg, bmarg, lmarg)
		alter = true
	}

    // may have to introduce an exemption for title
    parmap.keepLines = false
    if parStyl.KeepLinesTogether != parmap.keepLines {
        parmap.keepLines = parStyl.KeepLinesTogether
    }

    parmap.keepNext = false
    if parStyl.KeepWithNext != parmap.keepNext {
        parmap.keepNext = parStyl.KeepWithNext
    }

    switch parStyl.SpacingMode {
    case "NEVER_COLLAPSE":
        parmap.spaceMode = true
    case "COLLAPSE_LISTS":
        parmap.spaceMode = false
    default:
        parmap.spaceMode = true
    }

    // check for zero width borders
    bb := true
    if (parStyl.BorderBetween != nil) {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if (parStyl.BorderTop != nil) {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }


    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    parmap.hasBorders = true
    if !bb {
        parmap.hasBorders = false
		cssStr = cssStr[:len(cssStr)-1]
        return cssStr, alter
    }

	// border between paragraphs
    bordDisp := false
    parmap.bordBet.width = 0
    if parStyl.BorderBetween != nil {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderBetween.Width.Magnitude != parmap.bordBet.width {
                    parmap.bordBet.width = parStyl.BorderBetween.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBetween.Padding != nil {
            if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
                parmap.bordBet.pad = parStyl.BorderBetween.Padding.Magnitude
            }
        }
        if parStyl.BorderBetween.Color != nil {
            if parStyl.BorderBetween.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBetween.Color.Color)
                if color != parmap.bordBet.color {
                    parmap.bordBet.color = color
                }
            }
        }
        parmap.bordBet.dash = "SOLID"
        if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {parmap.bordBet.dash = parStyl.BorderBetween.DashStyle;}
    }

	// top border
    parmap.bordTop.width = 0
    if parStyl.BorderTop != nil {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderTop.Width.Magnitude != parmap.bordTop.width {
                    parmap.bordTop.width = parStyl.BorderTop.Width.Magnitude
                }
            }
        }
        if parStyl.BorderTop.Padding != nil {
            if parStyl.BorderTop.Padding.Magnitude != parmap.bordTop.pad {
                parmap.bordTop.pad = parStyl.BorderTop.Padding.Magnitude
            }
        }
        if parStyl.BorderTop.Color != nil {
            if parStyl.BorderTop.Color.Color != nil {
                color := util.GetColor(parStyl.BorderTop.Color.Color)
                if color != parmap.bordTop.color {
                    parmap.bordTop.color = color
                }
            }
        }
        parmap.bordTop.dash = "SOLID"
        if parStyl.BorderTop.DashStyle != parmap.bordTop.dash {parmap.bordTop.dash = parStyl.BorderTop.DashStyle;}
    }

	// right border
    parmap.bordRight.width = 0
    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderRight.Width.Magnitude != parmap.bordRight.width {
                    parmap.bordRight.width = parStyl.BorderRight.Width.Magnitude
                }
            }
        }
        if parStyl.BorderRight.Padding != nil {
            if parStyl.BorderRight.Padding.Magnitude != parmap.bordRight.pad {
                parmap.bordRight.pad = parStyl.BorderRight.Padding.Magnitude
            }
        }
        if parStyl.BorderRight.Color != nil {
            if parStyl.BorderRight.Color.Color != nil {
                color := util.GetColor(parStyl.BorderRight.Color.Color)
                if color != parmap.bordRight.color {
                    parmap.bordRight.color = color
                }
            }
        }
        parmap.bordRight.dash = "SOLID"
        if parStyl.BorderRight.DashStyle != parmap.bordRight.dash {
            parmap.bordRight.dash = parStyl.BorderRight.DashStyle
        }
    }

    // bottom border
    parmap.bordBot.width = 0
    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderBottom.Width.Magnitude != parmap.bordBot.width {
                    parmap.bordBot.width = parStyl.BorderBottom.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBottom.Padding != nil {
            if parStyl.BorderBottom.Padding.Magnitude != parmap.bordBot.pad {
                parmap.bordBot.pad = parStyl.BorderBottom.Padding.Magnitude
            }
        }
        if parStyl.BorderBottom.Color != nil {
            if parStyl.BorderBottom.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBottom.Color.Color)
                if color != parmap.bordBot.color {
                    parmap.bordBot.color = color
                }
            }
        }
        parmap.bordBot.dash = "SOLID"
        if parStyl.BorderBottom.DashStyle != parmap.bordBot.dash {parmap.bordBot.dash = parStyl.BorderBottom.DashStyle;}
    }

    // left border
    parmap.bordLeft.width = 0
    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderLeft.Width.Magnitude != parmap.bordLeft.width {
                    parmap.bordLeft.width = parStyl.BorderLeft.Width.Magnitude
                }
            }
        }
        if parStyl.BorderLeft.Padding != nil {
            if parStyl.BorderLeft.Padding.Magnitude != parmap.bordLeft.pad {
                parmap.bordLeft.pad = parStyl.BorderLeft.Padding.Magnitude
            }
        }
        if parStyl.BorderLeft.Color != nil {
            if parStyl.BorderLeft.Color.Color != nil {
                color := util.GetColor(parStyl.BorderLeft.Color.Color)
                if color != parmap.bordLeft.color {
                    parmap.bordLeft.color = color
                }
            }
        }
        parmap.bordLeft.dash = "SOLID"
        if parStyl.BorderLeft.DashStyle != parmap.bordLeft.dash {parmap.bordLeft.dash = parStyl.BorderLeft.DashStyle;}
    }

    if !bordDisp {
        parmap.hasBorders = false
		if len(cssStr) > 0 {cssStr = cssStr[:len(cssStr) - 1]}
        return cssStr, alter
    }

    cssStr += fmt.Sprintf("\"padding\": \"%.1fpt %.1fpt %.1fpt %.1fpt\",", parmap.bordTop.pad, parmap.bordRight.pad, parmap.bordBot.pad, parmap.bordLeft.pad)
    cssStr += fmt.Sprintf("\"borderTop\": \"%.1fpt %s %s\",", parmap.bordTop.width, util.GetDash(parmap.bordTop.dash), parmap.bordTop.color)
    cssStr += fmt.Sprintf("\"borderRight\": \"%.1fpt %s %s\",", parmap.bordRight.width, util.GetDash(parmap.bordRight.dash), parmap.bordRight.color)
    cssStr += fmt.Sprintf("\"borderBottom\": \"%.1fpt %s %s\",", parmap.bordBot.width, util.GetDash(parmap.bordBot.dash), parmap.bordBot.color)
    cssStr += fmt.Sprintf("\"borderLeft\": \"%.1fpt %s %s\",", parmap.bordLeft.width, util.GetDash(parmap.bordLeft.dash), parmap.bordLeft.color)

	if len(cssStr) > 0 {cssStr = cssStr[:len(cssStr) - 1]}
//	fmt.Printf("css: %s\n %q %q %q", cssStr, cssStr[len(cssStr) -3], cssStr[len(cssStr) - 2], cssStr[len(cssStr) -1])
    return cssStr, alter
}

func cvtParMapStylToCssJson(parmap *parMap, parStyl *docs.ParagraphStyle, opt *util.OptObj)(cssStr string, alter bool) {
// function that creates the css attributes of a paragraph
// the function compares the values of the parMap and parStyl
    if parmap == nil {return "/* no parmap */", false}
    if parStyl == nil { return "/* no parStyl */", false}

    cssStr =""

    if (len(parStyl.Alignment) > 0) &&  (parmap.halign != parStyl.Alignment) {
		alter = true
        parmap.halign = parStyl.Alignment
        switch parmap.halign {
            case "START":
                cssStr += "text-align: left;"
            case "CENTER":
                cssStr += "text-align: center;"
            case "END":
                cssStr += "text-align: right;"
            case "JUSTIFIED":
                cssStr += "text-align: justify;"
            default:
                cssStr += "text-align: unrecognized;"
        }

    }

    // test direction skip for now
    parmap.direct = true

    if (parStyl.IndentFirstLine != nil) {
        if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
            parmap.indFlin = parStyl.IndentFirstLine.Magnitude
            cssStr += fmt.Sprintf("text-indent: %.1fpt;", parmap.indFlin)
			alter = true
        }
    }

    parmap.linSpac = 1.0
    if parStyl.LineSpacing/100.0 != parmap.linSpac {
        if parStyl.LineSpacing > 1.0 {
			alter = true
            parmap.linSpac = parStyl.LineSpacing/100.0
            if opt.DefLinSpacing > 0.0 {
                cssStr += fmt.Sprintf("line-height: %.2f;", opt.DefLinSpacing*parmap.linSpac)
            } else {
                cssStr += fmt.Sprintf("line-height: %.2f;", parmap.linSpac)
            }
        }
    }

    margin := false
    lmarg := 0.0
    if (parStyl.IndentStart != nil) {
        if parStyl.IndentStart.Magnitude != parmap.indStart {
            parmap.indStart = parStyl.IndentStart.Magnitude
            lmarg = parmap.indStart
            margin = true
        }
    }

    rmarg := 0.0
    if (parStyl.IndentEnd != nil) {
        if parStyl.IndentEnd.Magnitude != parmap.indEnd {
            parmap.indEnd = parStyl.IndentEnd.Magnitude
            rmarg = parmap.indEnd
            margin = true
        }
    }

    tmarg := 0.0
    if (parStyl.SpaceAbove != nil) {
        if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
	        parmap.spaceTop = parStyl.SpaceAbove.Magnitude
            tmarg = parmap.spaceTop
            margin = true
        }
    }

    bmarg := 0.0
    if (parStyl.SpaceBelow != nil) {
        if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
            parmap.spaceBelow = parStyl.SpaceBelow.Magnitude
            bmarg = parmap.spaceBelow
            margin = true
        }
    }

    if margin {
		cssStr += fmt.Sprintf("margin: %.0fpt %.0fpt %.0fpt %.0fpt;", tmarg, rmarg, bmarg, lmarg)
		alter = true
	}

    // may have to introduce an exemption for title
    parmap.keepLines = false
    if parStyl.KeepLinesTogether != parmap.keepLines {
        parmap.keepLines = parStyl.KeepLinesTogether
    }

    parmap.keepNext = false
    if parStyl.KeepWithNext != parmap.keepNext {
        parmap.keepNext = parStyl.KeepWithNext
    }

    switch parStyl.SpacingMode {
    case "NEVER_COLLAPSE":
        parmap.spaceMode = true
    case "COLLAPSE_LISTS":
        parmap.spaceMode = false
    default:
        parmap.spaceMode = true
    }

    // check for zero width borders
    bb := true
    if (parStyl.BorderBetween != nil) {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if (parStyl.BorderTop != nil) {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }


    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.1 {
                bb = false
            }
        }
    }

    parmap.hasBorders = true
    if !bb {
        parmap.hasBorders = false
		cssStr = cssStr[:len(cssStr)-1]
        return cssStr, alter
    }

	// border between paragraphs
    bordDisp := false
    parmap.bordBet.width = 0
    if parStyl.BorderBetween != nil {
        if parStyl.BorderBetween.Width != nil {
            if parStyl.BorderBetween.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderBetween.Width.Magnitude != parmap.bordBet.width {
                    parmap.bordBet.width = parStyl.BorderBetween.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBetween.Padding != nil {
            if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
                parmap.bordBet.pad = parStyl.BorderBetween.Padding.Magnitude
            }
        }
        if parStyl.BorderBetween.Color != nil {
            if parStyl.BorderBetween.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBetween.Color.Color)
                if color != parmap.bordBet.color {
                    parmap.bordBet.color = color
                }
            }
        }
        parmap.bordBet.dash = "SOLID"
        if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {parmap.bordBet.dash = parStyl.BorderBetween.DashStyle;}
    }

	// top border
    parmap.bordTop.width = 0
    if parStyl.BorderTop != nil {
        if parStyl.BorderTop.Width != nil {
            if parStyl.BorderTop.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderTop.Width.Magnitude != parmap.bordTop.width {
                    parmap.bordTop.width = parStyl.BorderTop.Width.Magnitude
                }
            }
        }
        if parStyl.BorderTop.Padding != nil {
            if parStyl.BorderTop.Padding.Magnitude != parmap.bordTop.pad {
                parmap.bordTop.pad = parStyl.BorderTop.Padding.Magnitude
            }
        }
        if parStyl.BorderTop.Color != nil {
            if parStyl.BorderTop.Color.Color != nil {
                color := util.GetColor(parStyl.BorderTop.Color.Color)
                if color != parmap.bordTop.color {
                    parmap.bordTop.color = color
                }
            }
        }
        parmap.bordTop.dash = "SOLID"
        if parStyl.BorderTop.DashStyle != parmap.bordTop.dash {parmap.bordTop.dash = parStyl.BorderTop.DashStyle;}
    }

	// right border
    parmap.bordRight.width = 0
    if parStyl.BorderRight != nil {
        if parStyl.BorderRight.Width != nil {
            if parStyl.BorderRight.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderRight.Width.Magnitude != parmap.bordRight.width {
                    parmap.bordRight.width = parStyl.BorderRight.Width.Magnitude
                }
            }
        }
        if parStyl.BorderRight.Padding != nil {
            if parStyl.BorderRight.Padding.Magnitude != parmap.bordRight.pad {
                parmap.bordRight.pad = parStyl.BorderRight.Padding.Magnitude
            }
        }
        if parStyl.BorderRight.Color != nil {
            if parStyl.BorderRight.Color.Color != nil {
                color := util.GetColor(parStyl.BorderRight.Color.Color)
                if color != parmap.bordRight.color {
                    parmap.bordRight.color = color
                }
            }
        }
        parmap.bordRight.dash = "SOLID"
        if parStyl.BorderRight.DashStyle != parmap.bordRight.dash {
            parmap.bordRight.dash = parStyl.BorderRight.DashStyle
        }
    }

    // bottom border
    parmap.bordBot.width = 0
    if parStyl.BorderBottom != nil {
        if parStyl.BorderBottom.Width != nil {
            if parStyl.BorderBottom.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderBottom.Width.Magnitude != parmap.bordBot.width {
                    parmap.bordBot.width = parStyl.BorderBottom.Width.Magnitude
                }
            }
        }
        if parStyl.BorderBottom.Padding != nil {
            if parStyl.BorderBottom.Padding.Magnitude != parmap.bordBot.pad {
                parmap.bordBot.pad = parStyl.BorderBottom.Padding.Magnitude
            }
        }
        if parStyl.BorderBottom.Color != nil {
            if parStyl.BorderBottom.Color.Color != nil {
                color := util.GetColor(parStyl.BorderBottom.Color.Color)
                if color != parmap.bordBot.color {
                    parmap.bordBot.color = color
                }
            }
        }
        parmap.bordBot.dash = "SOLID"
        if parStyl.BorderBottom.DashStyle != parmap.bordBot.dash {parmap.bordBot.dash = parStyl.BorderBottom.DashStyle;}
    }

    // left border
    parmap.bordLeft.width = 0
    if parStyl.BorderLeft != nil {
        if parStyl.BorderLeft.Width != nil {
            if parStyl.BorderLeft.Width.Magnitude > 0.0 {
                bordDisp = true
                if parStyl.BorderLeft.Width.Magnitude != parmap.bordLeft.width {
                    parmap.bordLeft.width = parStyl.BorderLeft.Width.Magnitude
                }
            }
        }
        if parStyl.BorderLeft.Padding != nil {
            if parStyl.BorderLeft.Padding.Magnitude != parmap.bordLeft.pad {
                parmap.bordLeft.pad = parStyl.BorderLeft.Padding.Magnitude
            }
        }
        if parStyl.BorderLeft.Color != nil {
            if parStyl.BorderLeft.Color.Color != nil {
                color := util.GetColor(parStyl.BorderLeft.Color.Color)
                if color != parmap.bordLeft.color {
                    parmap.bordLeft.color = color
                }
            }
        }
        parmap.bordLeft.dash = "SOLID"
        if parStyl.BorderLeft.DashStyle != parmap.bordLeft.dash {parmap.bordLeft.dash = parStyl.BorderLeft.DashStyle;}
    }

    if !bordDisp {
        parmap.hasBorders = false
        return cssStr, alter
    }

    cssStr += fmt.Sprintf("padding: %.1fpt %.1fpt %.1fpt %.1fpt;", parmap.bordTop.pad, parmap.bordRight.pad, parmap.bordBot.pad, parmap.bordLeft.pad)
    cssStr += fmt.Sprintf("border-top: %.1fpt %s %s;", parmap.bordTop.width, util.GetDash(parmap.bordTop.dash), parmap.bordTop.color)
    cssStr += fmt.Sprintf("border-right: %.1fpt %s %s;", parmap.bordRight.width, util.GetDash(parmap.bordRight.dash), parmap.bordRight.color)
    cssStr += fmt.Sprintf("border-bottom: %.1fpt %s %s;", parmap.bordBot.width, util.GetDash(parmap.bordBot.dash), parmap.bordBot.color)
    cssStr += fmt.Sprintf("border-left: %.1fpt %s %s;", parmap.bordLeft.width, util.GetDash(parmap.bordLeft.dash), parmap.bordLeft.color)

//	if len(cssStr) > 0 {cssStr = cssStr[:len(cssStr) - 1]}
    return cssStr, alter
}

func cvtParMapToJson(pMap *parMap, opt *util.OptObj)(cssStr string) {
    cssStr =""

    if len(pMap.halign) > 0 {
        switch pMap.halign {
            case "START":
                cssStr += "\"textAlign\": \"left\","
            case "CENTER":
                cssStr += "\"textAlign\": \"center\","
            case "END":
                cssStr += "\"textAlign\": \"right\","
            case "JUSTIFIED":
                cssStr += "\"textAlign\": \"justify\","
            default:
                cssStr += fmt.Sprintf("\n// unrecognized Alignment %s \n", pMap.halign)
        }

    }

    if pMap.linSpac > 0.0 {
        if opt.DefLinSpacing > 0.0 {
            cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", opt.DefLinSpacing*pMap.linSpac)
        } else {
            cssStr += fmt.Sprintf("\"lineHeight\": \"%.2f\",", pMap.linSpac)
        }
    }

    if pMap.indFlin > 0.0 {
        cssStr += fmt.Sprintf("\"textIndent\": \"%.1fpt\",", pMap.indFlin)
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

    if margin {cssStr += fmt.Sprintf("\"margin\": \"%.0fpt %.0fpt %.0fpt %.0fpt\",", tmarg, rmarg, bmarg, lmarg)}

    if !pMap.hasBorders {
		if len(cssStr) > 0 { cssStr = cssStr[:len(cssStr)-1] }
		return cssStr
	}
    cssStr += fmt.Sprintf("\"padding\": \"%.1fpt %.1fpt %.1fpt %.1fpt\",", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
    cssStr += fmt.Sprintf("\"borderTop\": \"%.1fpt %s %s\",", pMap.bordTop.width, util.GetDash(pMap.bordTop.dash), pMap.bordTop.color)
    cssStr += fmt.Sprintf("\"borderRight\": \"%.1fpt %s %s\",", pMap.bordRight.width, util.GetDash(pMap.bordRight.dash), pMap.bordRight.color)
    cssStr += fmt.Sprintf("\"borderBottom\": \"%.1fpt %s %s\",", pMap.bordBot.width, util.GetDash(pMap.bordBot.dash), pMap.bordBot.color)
    cssStr += fmt.Sprintf("\"borderLeft: \"%.1fpt %s %s\",", pMap.bordLeft.width, util.GetDash(pMap.bordLeft.dash), pMap.bordLeft.color)

	if len(cssStr) > 0 { cssStr = cssStr[:len(cssStr)-1] }
    return cssStr
}

func cvtParMapToCssJson(pMap *parMap, opt *util.OptObj)(cssStr string) {
    cssStr =""

    if len(pMap.halign) > 0 {
        switch pMap.halign {
            case "START":
                cssStr += "text-align: left;"
            case "CENTER":
                cssStr += "text-align: center;"
            case "END":
                cssStr += "text-align: right;"
            case "JUSTIFIED":
                cssStr += "text-align: justify;"
            default:
                cssStr += fmt.Sprintf("\n// unrecognized Alignment %s \n", pMap.halign)
        }

    }

    if pMap.linSpac > 0.0 {
        if opt.DefLinSpacing > 0.0 {
            cssStr += fmt.Sprintf(" line-height: %.2f;", opt.DefLinSpacing*pMap.linSpac)
        } else {
            cssStr += fmt.Sprintf(" line-height: %.2f;", pMap.linSpac)
        }
    }

    if pMap.indFlin > 0.0 {
        cssStr += fmt.Sprintf(" text-indent: %.1fpt;", pMap.indFlin)
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

    if margin {cssStr += fmt.Sprintf("margin: %.0fpt %.0fpt %.0fpt %.0fpt;", tmarg, rmarg, bmarg, lmarg)}

    if !pMap.hasBorders {
//		if len(cssStr) > 0 { cssStr = cssStr[:len(cssStr)-1] }
		return cssStr
	}
    cssStr += fmt.Sprintf(" padding: %.1fpt %.1fpt %.1fpt %.1fpt;", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
    cssStr += fmt.Sprintf(" border-top: %.1fpt %s %s;", pMap.bordTop.width, util.GetDash(pMap.bordTop.dash), pMap.bordTop.color)
    cssStr += fmt.Sprintf(" border-right: %.1fpt %s %s;", pMap.bordRight.width, util.GetDash(pMap.bordRight.dash), pMap.bordRight.color)
    cssStr += fmt.Sprintf(" border-bottom: %.1fpt %s %s;", pMap.bordBot.width, util.GetDash(pMap.bordBot.dash), pMap.bordBot.color)
    cssStr += fmt.Sprintf(" border-left: %.1fpt %s %s;", pMap.bordLeft.width, util.GetDash(pMap.bordLeft.dash), pMap.bordLeft.color)

    return cssStr
}

func cvtParMapToCssJsonDeprecated(pMap *parMap, opt *util.OptObj)(cssStr string) {
    cssStr =""

    if len(pMap.halign) > 0 {
        switch pMap.halign {
            case "START":
                cssStr += "textAlign: left;"
            case "CENTER":
                cssStr += "textAlign: center;"
            case "END":
                cssStr += "textAlign: right;"
            case "JUSTIFIED":
                cssStr += "textAlign: justify;"
            default:
                cssStr += fmt.Sprintf("\n// unrecognized Alignment %s \n", pMap.halign)
        }

    }

    if pMap.linSpac > 0.0 {
        if opt.DefLinSpacing > 0.0 {
            cssStr += fmt.Sprintf(" lineHeight: %.2f;", opt.DefLinSpacing*pMap.linSpac)
        } else {
            cssStr += fmt.Sprintf(" lineHeight: %.2f;", pMap.linSpac)
        }
    }

    if pMap.indFlin > 0.0 {
        cssStr += fmt.Sprintf(" textIndent: %.1fpt;", pMap.indFlin)
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

    if margin {cssStr += fmt.Sprintf("margin: %.0fpt %.0fpt %.0fpt %.0fpt;", tmarg, rmarg, bmarg, lmarg)}

    if !pMap.hasBorders {
//		if len(cssStr) > 0 { cssStr = cssStr[:len(cssStr)-1] }
		return cssStr
	}
    cssStr += fmt.Sprintf(" padding: %.1fpt %.1fpt %.1fpt %.1fpt;", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
    cssStr += fmt.Sprintf(" borderTop: %.1fpt %s %s;", pMap.bordTop.width, util.GetDash(pMap.bordTop.dash), pMap.bordTop.color)
    cssStr += fmt.Sprintf(" borderRight: %.1fpt %s %s;", pMap.bordRight.width, util.GetDash(pMap.bordRight.dash), pMap.bordRight.color)
    cssStr += fmt.Sprintf(" borderBottom: %.1fpt %s %s;", pMap.bordBot.width, util.GetDash(pMap.bordBot.dash), pMap.bordBot.color)
    cssStr += fmt.Sprintf(" borderLeft: %.1fpt %s %s;", pMap.bordLeft.width, util.GetDash(pMap.bordLeft.dash), pMap.bordLeft.color)

	if len(cssStr) > 0 { cssStr = cssStr[:len(cssStr)-1] }
    return cssStr
}

func creJsonHead (docNam string) (outstr string) {

	outstr = "{\n\"doc\":{\"docNam\": \"" + docNam + "\"},\n"
    return outstr
}

//todo
func creTocSecJson(docName string)(cssStr string) {

	cssStr = fmt.Sprintf("  {\"cssRule\": \".%sMain.top ", docName)
	cssStr += "{padding: 10px 0 10px 0;}\"\n"

	cssStr += fmt.Sprintf("\"cssRule\":\".%s_title.leftTitle_UL {", docName)
	cssStr += "  text-align: start;"
	cssStr += "  text-decoration-line: underline;"
	cssStr += "}\"},\n"

	cssStr += fmt.Sprintf("  {\"cssRule\": \".%s_title.leftTitle {", docName)
	cssStr += "  text-align: start;"
	cssStr += "  text-decoration-line: none;"
	cssStr += "}\"},\n"

	cssStr += fmt.Sprintf("  {\"cssRule\": \".%s_noUl {", docName)
	cssStr += "  text-decoration: none;"
	cssStr += "}\"},\n"

	return cssStr
}

func creTocJson(docName string)(cssStr string) {
	cssStr = fmt.Sprintf("  {\"cssRule:\": \".%_div.toc {", docName)

	cssStr += "}\"},\n"
	return cssStr
}

//todo
func creSecJson(docName string)(cssStr string){

	cssStr = fmt.Sprintf("  {\"cssRule\": \".%sMain.sec {", docName)

	cssStr += "}\"},\n"

	cssStr += fmt.Sprintf("  {\"cssRule\": \".%sPage {", docName)
	cssStr += "text-align: right;"
	cssStr += "margin: 0;"
	cssStr += "}\"},\n"
	return cssStr
}

func creFtnoteJson(docName string)(cssStr string){
	//css footnote
	cssStr = fmt.Sprintf("  {\"cssRule\":\".%sFtnote\" {", docName)
//	cssStr += "vertical-align: super;"
	cssStr += "color: purple;"
	cssStr += "}\"},\n"
	return cssStr
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


func (dObj *GdocDomObj) printHeadings() {

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



func (dObj *GdocDomObj) getNamedStyl(namedTyp string)(parStyl *docs.ParagraphStyle, txtStyl *docs.TextStyle, err error) {
	var namStyl *docs.NamedStyle

	doc:= dObj.doc
// initialise named styles
	namStyles := doc.NamedStyles
	stylIdx := -1

// find normal style first
	for istyl:=0; istyl<len(namStyles.Styles); istyl++ {
		if namStyles.Styles[istyl].NamedStyleType == namedTyp {
			stylIdx = istyl
			namStyl = namStyles.Styles[istyl]
			break
		}
	}

	if stylIdx < 0 {
		return nil, nil, fmt.Errorf("cannot find named style %s!", namedTyp)
	}

	parStyl = namStyl.ParagraphStyle
	txtStyl = namStyl.TextStyle
	return parStyl, txtStyl, nil
}

func (dObj *GdocDomObj) findListProp (listId string) (listProp *docs.ListProperties) {

	found := false
	doc := dObj.doc

	for key, listItem := range doc.Lists  {
		if listId == key {
			listProp = listItem.ListProperties
			found = true
			break
		}
	}

	if found { return listProp}

	return nil
}

func (dObj *GdocDomObj) initGdocJson(folderPath string, options *util.OptObj) (err error) {
	var listItem docList
	var heading headingTyp
	var sec secTyp
	var ftnote docFtnoteTyp
	var docPb pbTyp

	doc := dObj.doc
	if doc == nil {return fmt.Errorf("doc is nil in GdocObj!")}

	// need to transform file name
	// replace spaces with underscore
	dNam := doc.Title
	x := []byte(dNam)
	for i:=0; i<len(x); i++ {
		if x[i] == ' ' {
			x[i] = '_'
		}
	}
	dObj.errCount = 0
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

	// elements
	dObj.elCount = len(doc.Body.Content)

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


	for el:=0; el<dObj.elCount; el++ {
		elObj:= doc.Body.Content[el]
		if elObj.SectionBreak != nil {
			if elObj.SectionBreak.SectionStyle.SectionType == "NEXT_PAGE" {
				sec.secElStart = el
				dObj.sections = append(dObj.sections, sec)
				seclen := len(dObj.sections)
				//	fmt.Println("el: ", el, "section len: ", seclen, " secPtEnd: ", secPtEnd)
				if seclen > 1 {
					dObj.sections[seclen-2].secElEnd = secPtEnd
				}
			}
		}

		// paragraphs and lists
		if elObj.Paragraph != nil {
			if elObj.Paragraph.Bullet != nil {
			// lists
				listId := elObj.Paragraph.Bullet.ListId
				found := findDocList(dObj.docLists, listId)
				nestlev := elObj.Paragraph.Bullet.NestingLevel
				if found < 0 {
					listItem.listId = listId
					listItem.maxNestLev = elObj.Paragraph.Bullet.NestingLevel
					nestL := doc.Lists[listId].ListProperties.NestingLevels[nestlev]
					listItem.ord = util.GetGlyphOrd(nestL)
					dObj.docLists = append(dObj.docLists, listItem)
				} else {
					if dObj.docLists[found].maxNestLev < nestlev { dObj.docLists[found].maxNestLev = nestlev }
				}

			}

			// named styles
			namedStyl := elObj.Paragraph.ParagraphStyle.NamedStyleType
			if len(namedStyl) > 0 {
				if !dObj.namStylMap[namedStyl] {
					dObj.namStylMap[namedStyl] = true
				}
			}

			// headings
			// used for creating a table of content TOC

			text := ""
			if len(elObj.Paragraph.ParagraphStyle.HeadingId) > 0 {
				heading.id = elObj.Paragraph.ParagraphStyle.HeadingId
				heading.hdElStart = el

				for parel:=0; parel<len(elObj.Paragraph.Elements); parel++ {
					if elObj.Paragraph.Elements[parel].TextRun != nil {
						text += elObj.Paragraph.Elements[parel].TextRun.Content
					}
				}
				txtlen:= len(text)
				if text[txtlen -1] == '\n' { text = text[:txtlen-1] }

				// assign attribute to heading
				heading.text = text
				heading.namedStyl = namedStyl

				dObj.headings = append(dObj.headings, heading)
			} // end headings

           // paragraph elements
			for parEl:=0; parEl<len(elObj.Paragraph.Elements); parEl++ {
				parElObj := elObj.Paragraph.Elements[parEl]
				// footnotes
				if parElObj.FootnoteReference != nil {
					ftnote.el = el
					ftnote.parel = parEl
					ftnote.id = parElObj.FootnoteReference.FootnoteId
					ftnote.numStr = parElObj.FootnoteReference.FootnoteNumber
					dObj.docFtnotes = append(dObj.docFtnotes, ftnote)
				}
               // page break
                if parElObj.PageBreak != nil {
                    docPb.el = el
                    docPb.parel = parEl
                    dObj.docPb = append(dObj.docPb, docPb)
                    dObj.pbCount++
                }
			}

			parHdEnd = el
			secPtEnd = el
		} // end paragraph

		if elObj.Table != nil {
			dObj.tableCount++
		}

	} // end el loop

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
	dObj.inImgCount = len(doc.InlineObjects)
	dObj.posImgCount = len(doc.PositionedObjects)

// create folders
//todo add opt of creating docnam folder
    fPath, fexist, err := util.CreateFileFolder(folderPath, "")
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

    // create output file path/outfilNam.json
	outfilNam := dObj.docName
    outfil, err := util.CreateOutFil(fPath, outfilNam, "json")
    if err!= nil {
        return fmt.Errorf("error -- util.CreateOutFil: %v", err)
    }
    dObj.jsonFil = outfil

    totObjNum := dObj.inImgCount + dObj.posImgCount
//  if totObjNum == 0 {return nil}


    if dObj.Options.CreImgFolder && (totObjNum > 0) {
        imgFoldPath, err := util.CreateImgFolder(fPath ,dObj.docName)
        if err != nil {
            return fmt.Errorf("error -- CreateImgFolder: could create ImgFolder: %v!", err)
        }
        dObj.imgFoldNam = imgFoldPath
        err = util.DownloadImages(doc, imgFoldPath, dObj.Options)
        if err != nil {
            return fmt.Errorf("error -- downloadImages could download images: %v!", err)
        }
    }

//    dObj.parCount = len(doc.Body.Content)

	return nil
}

func (dObj *GdocDomObj) cvtGlyph(nLev *docs.NestingLevel)(cssStr string) {
var glyphTyp string

	glyphTyp = util.GetGlyphStr(nLev)
	if len(glyphTyp) == 0 {
		cssStr = fmt.Sprintf("/* unknown GlyphType: %s Symbol: %s */\n", nLev.GlyphType, nLev.GlyphSymbol)
	} else {
		cssStr = "  list-style-type: " + glyphTyp +";\n"
	}
	return cssStr
}

/*
func (dObj *GdocDomObj) cvtParTxtSingle(parElTxt *docs.TextRun, namedTyp string)(attStr string, crEnd bool, err error) {

	attStr = ""
	crEnd = false
	if parElTxt == nil {
		return "", crEnd
	}

	if !(len(parElTxt.Content) > 0)  {
		return "", crEnd, fmt.Errorf("no Content!")
	}

	// get namedTyp
	if !(len(namedTyp) >0) {namedTyp = "NORMAL_TEXT"}

	_, namedTxtStyl, err := dObj.getNamedStyl(namedTyp)
	if err != nil {	namedTyp = "NORMAL_TEXT"}
//NAMED_STYLE_TYPE_UNSPECIFIED

	if parElTxt.TextStyle.Link == nil {return "", crEnd, fmt.Errorf("link element")}

	txtStr := parElTxt.Content
	txtEnd := len(txtStr) -1
//	fmt.Printf("txt: %q %q %q\n",txtStr[txtEnd -2], txtStr[txtEnd -1], txtStr[txtEnd])
	if txtStr[txtEnd] == '\n' {
		txtStr = txtStr[:txtEnd]
		crEnd = true
	}

	txtMap := fillTxtMap(namedTxtStyl)

	stylStr := cvtTxtMapStylToJson(txtMap, parElTxt.TextStyle)

	attStr = " \"textContent\": \"" + txtStr + "\","
	if len(stylStr) > 0 {attStr += " \"style\": \"" + stylStr + "\","}
	return attStr, crEnd, nil
}
*/

func (dObj *GdocDomObj) cvtParTxtElToJson(parElTxt *docs.TextRun, namedTyp string)(elStr string, crEnd bool) {

	elStr = ""
	crEnd = false
	if parElTxt == nil {
		return "", crEnd
	}
	if !(len(parElTxt.Content) > 0)  {
		return "", crEnd
	}

	// get namedTyp
	if !(len(namedTyp) >0) {
		namedTyp = "NORMAL_TEXT"
	}

	_, namedTxtStyl, err := dObj.getNamedStyl(namedTyp)
	if err != nil {	namedTyp = "NORMAL_TEXT"}
//NAMED_STYLE_TYPE_UNSPECIFIED

	txtStr := parElTxt.Content
	txtEnd := len(txtStr) -1
//	fmt.Printf("txt: %q %q %q\n",txtStr[txtEnd -2], txtStr[txtEnd -1], txtStr[txtEnd])
	if txtStr[txtEnd] == '\n' {
		txtStr = txtStr[:txtEnd]
		crEnd = true
	}
	if len(txtStr)<1 { return "", true}

	txtMap := fillTxtMap(namedTxtStyl)

	spanStylStr := cvtTxtMapStylToJson(txtMap, parElTxt.TextStyle)

	elStr = ""
	if parElTxt.TextStyle.Link != nil {
		elStr = "{\"typ\":\"a\","
//		elStr += "\"parent\":\"" + dObj.parent + "\","
		elStr += "\"parent\":\"par\","
		elStr += "\"href\":\"" + parElTxt.TextStyle.Link.Url + "\","
		elStr += "\"textContent\":\"" + txtStr + "\""
		if len(spanStylStr) > 0 {
			elStr += ", \"style\":{" + spanStylStr +"}"
		}
		elStr += "},\n"
		return elStr, crEnd
	}

	elStr = "{\"typ\":\"span\","
	elStr += "\"parent\":\"par\","
	elStr += "\"textContent\":\"" + txtStr + "\""
	if len(spanStylStr) > 0 {
			elStr += ", \"style\":{" + spanStylStr +"}"
	}
	elStr += "},\n"
	return elStr, crEnd
}

func (dObj *GdocDomObj) closeList(nl int) {
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

func (dObj *GdocDomObj) cvtHrElToJson (hr *docs.HorizontalRule)(hrJsonEl string) {
    var cssStr string

	hrJsonEl = "{\"typ\":\"hr\",\"parent\":\"gdocMain\","
    if hr.TextStyle != nil {
        cssStr = "\"style\": {"
        cssStr += cvtTxtStylJson(hr.TextStyle)
        cssStr += "}"
		hrJsonEl += cssStr
    }
	hrJsonEl += "},"
    return hrJsonEl
}

func (dObj *GdocDomObj) cvtFtnoteToJson ()(ftnoteStr string) {
//        	htmlStr += fmt.Sprintf("<span class=\"%s_ftnote\">[%d]</span>",dObj.docName, dObj.ftnoteCount)
	ftnoteStr = "{\"typ\":\"span\",\"parent\":\"gdocMain\",\"className\":\""+ dObj.docName + "FtNote\"},"

	return ftnoteStr
}

func (dObj *GdocDomObj) renderInlineImg(imgEl *docs.InlineObjectElement)(imgElStr string, err error) {

	if imgEl == nil {
		return "", fmt.Errorf("imgEl is nil!")
	}
	doc := dObj.doc

	imgElId := imgEl.InlineObjectId
	if !(len(imgElId) > 0) {return "", fmt.Errorf("no InlineObjectId found!")}

	// need to remove first part of the id
	idx := 0
	for i:=0; i< len(imgElId); i++ {
		if imgElId[i] == '.' {
			idx = i+1
			break
		}
	}
	imgId :=""
	if (idx>0) && (idx<len(imgElId)-1) {
		imgId = "img_" + imgElId[idx:]
	}

	// need to change for imagefolder
//	htmlStr := fmt.Sprintf("<!-- inline image %s -->\n", imgElId)

	imgObj := doc.InlineObjects[imgElId].InlineObjectProperties.EmbeddedObject

	imgSrcUri :=""
	if dObj.Options.ImgFold {
    	imgSrcUri = dObj.imgFoldNam + "/" + imgId + ".jpeg"
		// html htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgSrc, imgId, imgObj.Title)
	} else {
		// html htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgObj.ImageProperties.SourceUri, imgId, imgObj.Title)
		imgSrcUri = imgObj.ImageProperties.SourceUri
	}

	elStylStr := fmt.Sprintf("width:%.0fpt; height:%.0fpt;", imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )

	// todo add margin

	imgElStr = "{\"typ\":\"img\","
	imgElStr += "\"parent\":\"gdocMain\","
//	imgElStr += "\"parent\":\"" + dObj.parent + "\","
	imgElStr += "\"id\":\"" + imgId + "\","
	imgElStr += "\"src\":\"" + imgSrcUri + "\","
	imgElStr += "\"style\":{" + elStylStr + "}"
	if len(imgObj.Title) > 0 {imgElStr += ", \"title\":\"" + imgObj.Title +"\""}
	if len(imgObj.Description) > 0 {imgElStr += ", \"desc\":\"" + imgObj.Description + "\""}
	imgElStr += "}\n"
	return imgElStr, nil
}


func (dObj *GdocDomObj) renderPosImg(posImg *docs.PositionedObject, posId string)(imgElStr string, err error) {

/*
	// html
	posObjProp := posImg.PositionedObjectProperties
	imgProp := posObjProp.EmbeddedObject

	// fmt.Sprintf("\n<!-- Positioned Image %s -->\n", posId)
	scriptStr = "// *** Positioned Image " + posId + " ***\n"

	imgId := posId[4:]

	layout := posObjProp.Positioning.Layout
	topPos := posObjProp.Positioning.TopOffset.Magnitude
//	leftPos := posObjProp.Positioning.LeftOffset.Magnitude

//	fmt.Printf("layout %s top: %.1fmm left:%.1fmm\n", layout, topPos*PtTomm, leftPos*PtTomm)


	imgSrc := imgProp.ImageProperties.ContentUri
	if dObj.Options.ImgFold {
		imgSrc = dObj.imgFoldNam + "/" + posId[4:] + ".jpeg"
	}

	//css
	cssStr = ""
	switch layout {
		case "WRAP_TEXT", "BREAK_LEFT":
			cssStr += fmt.Sprintf("#%s {\n", imgId)
			cssStr += fmt.Sprintf("float:left; clear:both;")
			cssStr += fmt.Sprintf("  width:%.1fpt; height:%.1fpt;\n",imgProp.Size.Width.Magnitude, imgProp.Size.Height.Magnitude)
			cssStr += fmt.Sprintf("  margin: %.1fpt %.1fpt %.1fpt %.1fpt;\n", imgProp.MarginTop.Magnitude, imgProp.MarginRight.Magnitude, imgProp.MarginBottom.Magnitude, imgProp.MarginLeft.Magnitude)
			cssStr += "}\n"
			cssStr += fmt.Sprintf("#%s {\n", pimgId)
			cssStr += fmt.Sprintf("  margin-left: %.1fpt; margin-right: %.1fpt;\n", imgProp.MarginLeft.Magnitude, imgProp.MarginRight.Magnitude)
			cssStr += "}\n"
			cssStr += fmt.Sprintf("#%s:before {\n", imgDivId)
			cssStr += fmt.Sprintf("content:''; display:block; float:left; height:%.1fmm;\n",topPos*PtTomm)
			cssStr += "}\n"

		case "BREAK_RIGHT":
			cssStr += fmt.Sprintf("#%s {\n", imgId)
			cssStr += fmt.Sprintf("float:right; clear:both;")
			cssStr += fmt.Sprintf("  width:%.1fpt; height:%.1fpt;\n",imgProp.Size.Width.Magnitude, imgProp.Size.Height.Magnitude)
			cssStr += fmt.Sprintf("  margin: %.1fpt %.1fpt %.1fpt %.1fpt;\n", imgProp.MarginTop.Magnitude, imgProp.MarginRight.Magnitude, imgProp.MarginBottom.Magnitude, imgProp.MarginLeft.Magnitude)
			cssStr += "}\n"
			cssStr += fmt.Sprintf("#%s {\n", pimgId)
			cssStr += fmt.Sprintf("  margin-left: %.1fpt; margin-right: %.1fpt;\n", imgProp.MarginLeft.Magnitude, imgProp.MarginRight.Magnitude)
			cssStr += "}\n"
			cssStr += fmt.Sprintf("#%s:before {\n", imgDivId)
			cssStr += fmt.Sprintf("content:''; display:block; float:right; height:%.1fmm;\n",topPos*PtTomm)
			cssStr += "}\n"

		case "BREAK_LEFT_RIGHT":

		case "IN_FRONT_OF_TEXT":
// absolute
		case "BEHIND_TEXT":
// absolute
		default:
			cssStr += fmt.Sprintf("#%s {\n", imgId)
			cssStr += fmt.Sprintf("  width:%.1fpt; height:%.1fpt;\n",imgProp.Size.Width.Magnitude, imgProp.Size.Height.Magnitude)
			cssStr += fmt.Sprintf("  margin: %.1fpt %.1fpt %.1fpt %.1fpt;\n", imgProp.MarginTop.Magnitude, imgProp.MarginRight.Magnitude, imgProp.MarginBottom.Magnitude, imgProp.MarginLeft.Magnitude)
			cssStr += "}\n"
			cssStr += fmt.Sprintf("#%s {\n", pimgId)
			cssStr += fmt.Sprintf("  margin-left: %.1fpt; margin-right: %.1fpt;\n", imgProp.MarginLeft.Magnitude, imgProp.MarginRight.Magnitude)
			cssStr += "}\n"
	}

	// html
	// fmt.Sprintf("  <div id=\"%s\">\n",imgDivId)
	divEl.parent = dObj.parent
	divEl.typ = "div"
	divEl.newEl = "imgDiv"
	divEl.idStr = imgDivId
	divEl.doAppend = true
	scriptStr += addElToDom(divEl)

	//fmt.Sprintf("     <img src=\"%s\" alt=\"%s\" id=\"%s\">\n", imgSrc, imgProp.Title, imgId)
	imgEl.parent = "imgDiv"
	imgEl.src = imgSrc
	imgEl.title = imgProp.Title
	imgEl.idStr = imgId
	scriptStr += addImgElToDom(imgEl)

	//	fmt.Sprintf("     <p id=\"%s\">%s</p>\n", pimgId, imgProp.Title)
	if len(imgProp.Title) > 0 {
		divEl.parent = "imgDiv"
		divEl.typ = "p"
		divEl.txt = imgProp.Title
		divEl.idStr = pimgId
		divEl.cl1 = dObj.docName + "_p"
		divEl.doAppend = true
		scriptStr += addElToDom(divEl)
	}

	imgDispObj.script = scriptStr
	imgDispObj.bodyCss = cssStr
*/
	return imgElStr, nil
}


func (dObj *GdocDomObj) cvtTableToJson(tbl *docs.Table)(tabStr string, err error) {
	// https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Traversing_an_HTML_table_with_JavaScript_and_DOM_Interfaces
	// table element
//	var tabWidth float64

/*
	var icol, trow int64
	var defcel tblCell
	var tblObj, elObj elScriptObj
	var colObj colScriptObj
//	var tblCelObj elScriptObj

//	var tabcelObj tblCellScriptObj

	dObj.tableCounter++


//    docStyl := doc.DocumentStyle
//    PgWidth := docStyl.PageSize.Width.Magnitude
//    NetPgWidth := PgWidth - (docStyl.MarginLeft.Magnitude + docStyl.MarginRight.Magnitude)

//   fmt.Printf("Default Table Width: %.1f", NetPgWidth)
//    tabWidth = NetPgWidth

// table cell default values
// define default cell classs

	tcelDef := tbl.TableRows[0].TableCells[0]
	tcelDefStyl := tcelDef.TableCellStyle
	tblNam := "tbl"

// default values which google does not set but uses
	defcel.vert_align = "top"
	defcel.bcolor = "black"
	defcel.bwidth = 1.0
	defcel.bdash = "solid"

	if tcelDefStyl != nil {
		defcel.vert_align = util.Get_vert_align(tcelDefStyl.ContentAlignment)

		// if left border is the only border specified, let's use it for default values
		tb := (tcelDefStyl.BorderTop == nil)&& (tcelDefStyl.BorderRight == nil)
		tb = tb&&(tcelDefStyl.BorderBottom == nil)

		if (tcelDefStyl.BorderLeft != nil) && tb {
			if tcelDefStyl.BorderLeft.Color != nil {
				if tcelDefStyl.BorderLeft.Color.Color != nil {
					defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)
				}
			}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}

		if tcelDefStyl.PaddingTop != nil {defcel.pad[0] = tcelDefStyl.PaddingTop.Magnitude}
		if tcelDefStyl.PaddingRight != nil {defcel.pad[1] = tcelDefStyl.PaddingRight.Magnitude}
		if tcelDefStyl.PaddingBottom != nil {defcel.pad[2] = tcelDefStyl.PaddingBottom.Magnitude}
		if tcelDefStyl.PaddingLeft != nil {defcel.pad[3] = tcelDefStyl.PaddingLeft.Magnitude}

		if tcelDefStyl.BackgroundColor != nil {
			if tcelDefStyl.BackgroundColor.Color != nil {
				defcel.bckcolor = util.GetColor(tcelDefStyl.BackgroundColor.Color)
			}
		}

		if tcelDefStyl.BorderTop != nil {
			if tcelDefStyl.BorderTop.Color != nil {
				if tcelDefStyl.BorderTop.Color.Color != nil {
					defcel.border[0].color = util.GetColor(tcelDefStyl.BorderTop.Color.Color)
				}
			}
			if tcelDefStyl.BorderTop.Width != nil {defcel.border[0].width = tcelDefStyl.BorderTop.Width.Magnitude}
			defcel.border[0].dash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
		}

		if tcelDefStyl.BorderRight != nil {
			if tcelDefStyl.BorderRight.Color != nil {
				if tcelDefStyl.BorderRight.Color.Color != nil {
					defcel.border[1].color = util.GetColor(tcelDefStyl.BorderRight.Color.Color)
				}
			}
			if tcelDefStyl.BorderRight.Width != nil {defcel.border[1].width = tcelDefStyl.BorderRight.Width.Magnitude}
			defcel.border[1].dash = util.GetDash(tcelDefStyl.BorderRight.DashStyle)
		}

		if tcelDefStyl.BorderBottom != nil {
			if tcelDefStyl.BorderBottom.Color != nil {
				if tcelDefStyl.BorderBottom.Color.Color != nil {
					defcel.border[2].color = util.GetColor(tcelDefStyl.BorderBottom.Color.Color)
				}
			}
			if tcelDefStyl.BorderBottom.Width != nil {defcel.border[2].width = tcelDefStyl.BorderBottom.Width.Magnitude}
			defcel.border[2].dash = util.GetDash(tcelDefStyl.BorderBottom.DashStyle)
		}

		if tcelDefStyl.BorderLeft != nil {
			if tcelDefStyl.BorderLeft.Color != nil {
				if tcelDefStyl.BorderLeft.Color.Color != nil {
					defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)
				}
			}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}

		if tcelDefStyl.BorderTop == tcelDefStyl.BorderRight {
//			fmt.Println("same border!")
			if tcelDefStyl.BorderTop != nil {
				if tcelDefStyl.BorderTop.Color != nil {
					if tcelDefStyl.BorderTop.Color.Color != nil {
						defcel.bcolor = util.GetColor(tcelDefStyl.BorderTop.Color.Color)
					}
				}
				defcel.bdash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
				if tcelDefStyl.BorderTop.Width != nil {defcel.bwidth = tcelDefStyl.BorderTop.Width.Magnitude}
			}
		}
	}

	//set up table

	// if there is an open list, close it
	if dObj.listStack != nil {
		dObj.closeList(-1)
	}

	// html fmt.Sprintf("<table class=\"%s_tbl tbl_%d\">\n", dObj.docName, dObj.tableCounter)
	tblObj.parent = dObj.parent
	tblObj.typ = "table"
	tblObj.doAppend = false
	tblObj.newEl = tblNam
	tblObj.cl1 = dObj.docName + "_tbl"
	tblObj.cl2 = fmt.Sprintf("tbl_%d", dObj.tableCounter)
	scriptStr += addElToDom(tblObj)

	// html "  <tbody>\n"
	tblObj.parent = tblNam
	tblObj.typ = "tbody"
	tblObj.doAppend = true
	tblObj.newEl = "tblBody"
//	tblObj.cl1 = dObj.docName + "_tbl"
//	tblObj.cl2 = fmt.Sprintf("tbl_%d", dObj.tableCounter)
	scriptStr += addElToDom(tblObj)

	// table columns
	// conundrum: tables either have evenly distributed columns or not
	// should not be possible to have a mixture of evenly distributed columns and specified width columns
	// thus it should be sufficient to check the first column for that property

	tabWtyp :=tbl.TableStyle.TableColumnProperties[0].WidthType
	switch tabWtyp {
	case "EVENLY_DISTRIBUTED":
		// no need for column groups

	case "WIDTH_TYPE_UNSPECIFIED":
		// to be determined

	case "FIXED_WIDTH":
		// html htmlStr +="<colgroup>\n"
		elObj.parent = "tblBody"
		elObj.typ = "colgroup"
		elObj.newEl = "colgrp"
		elObj.doAppend = true
		scriptStr += addElToDom(elObj)
		tblW := 0.0
		for icol = 0; icol < tbl.Columns; icol++ {
            colW := tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude
			tblW += colW
            cssStr += fmt.Sprintf(".%s_colgrp_%d.col_%d {width: %.0fpt;}\n", dObj.docName, dObj.tableCounter, icol, colW)
            //htmlStr += fmt.Sprintf("<col span=\"1\" class=\"%s_colgrp_%d col_%d\">\n", dObj.docName, dObj.tableCounter, icol)
			colObj.parent = "colgrp"
			colObj.cl1 = fmt.Sprintf("%s_colgrp_%d", dObj.docName, dObj.tableCounter)
			colObj.cl2 = fmt.Sprintf("col_%d", icol)
			colObj.spanCount = 1
			scriptStr += addColToDom(colObj)
		}
		//if tabw > 0.0 {tabWidth = tabw}
		// html htmlStr +="</colgroup>\n"
	}

// row styling
	parent := dObj.parent
	tblCellCount := 0
	for trow=0; trow < tbl.Rows; trow++ {
		// html fmt.Sprintf("  <tr>\n")
		elObj.typ ="tr"
		elObj.cl1 = fmt.Sprintf("%s_tblrow", dObj.docName)
		elObj.parent = "tblBody"
		elObj.newEl = "trow"
		elObj.doAppend = true
		scriptStr += addElToDom(elObj)

		trowobj := tbl.TableRows[trow]
//		mrheight := trowobj.TableRowStyle.MinRowHeight.Magnitude

		numCols := len(trowobj.TableCells)
		for tcol:=0; tcol< numCols; tcol++ {
			tcell := trowobj.TableCells[tcol]
			tblCellCount++
			cellStr := ""

			// check whether cell style is different from default
			if tcell.TableCellStyle != nil {
				tstyl := tcell.TableCellStyle
				if tstyl.BackgroundColor != nil {
					if tstyl.BackgroundColor.Color != nil {
						cellStr += fmt.Sprintf(" background-color:\"%s\";",util.GetColor(tstyl.BackgroundColor.Color))
					}
				}
				if util.Get_vert_align(tstyl.ContentAlignment) != defcel.vert_align {cellStr += fmt.Sprintf(" vertical-align: %s;", util.Get_vert_align(tstyl.ContentAlignment))}
				if tstyl.PaddingTop != nil {
					if tstyl.PaddingTop.Magnitude != defcel.pad[0] { cellStr += fmt.Sprintf(" padding-top: %5.1fpt;", tstyl.PaddingTop.Magnitude)}
				}
                if tstyl.PaddingLeft != nil {
                    if tstyl.PaddingLeft.Magnitude != defcel.pad[1] { cellStr += fmt.Sprintf(" padding-left: %5.1fpt;", tstyl.PaddingLeft.Magnitude)}
                }
                if tstyl.PaddingBottom != nil {
                    if tstyl.PaddingBottom.Magnitude != defcel.pad[2] { cellStr += fmt.Sprintf(" padding-bottom: %5.1fpt;", tstyl.PaddingBottom.Magnitude)}
                }
                if tstyl.PaddingRight != nil {
                    if tstyl.PaddingRight.Magnitude != defcel.pad[3] { cellStr += fmt.Sprintf(" padding-right: %5.1fpt;", tstyl.PaddingRight.Magnitude)}
                }

                if tstyl.BorderTop != nil {
                    // Color
                    if tstyl.BorderTop.Color != nil {
                        if tstyl.BorderTop.Color.Color != nil {
                            cellStr += fmt.Sprintf(" border-top-color: %s;", util.GetColor(tstyl.BorderTop.Color.Color))
                        }
                    }
                    //dash
                    if util.GetDash(tstyl.BorderTop.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-top-style: %s;",  util.GetDash(tstyl.BorderTop.DashStyle))}
                    //Width
                    if tstyl.BorderTop.Width != nil {
                        cellStr += fmt.Sprintf(" border-top-width: %5.1fpt;", tstyl.BorderTop.Width.Magnitude)
                    }
                }

                if tstyl.BorderLeft != nil {
                    // Color
                    if tstyl.BorderLeft.Color != nil {
                        if tstyl.BorderLeft.Color.Color != nil {
                            cellStr += fmt.Sprintf(" border-left-color: %s;", util.GetColor(tstyl.BorderLeft.Color.Color))
                        }
                    }
                    //dash
                    if util.GetDash(tstyl.BorderLeft.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-left-style: %s;",  util.GetDash(tstyl.BorderLeft.DashStyle))}
                    //Width
                    if tstyl.BorderTop.Width != nil {
                        cellStr += fmt.Sprintf(" border-left-width: %5.1fpt;", tstyl.BorderLeft.Width.Magnitude)
                    }
                }


                if tstyl.BorderBottom != nil {
                    // Color
                    if tstyl.BorderBottom.Color != nil {
                        if tstyl.BorderBottom.Color.Color != nil {
                            cellStr += fmt.Sprintf(" border-bottom-color: %s;", util.GetColor(tstyl.BorderBottom.Color.Color))
                        }
                    }
                    //dash
                    if util.GetDash(tstyl.BorderBottom.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-bottom-style: %s;",  util.GetDash(tstyl.BorderBottom.DashStyle))}
                    //Width
                    if tstyl.BorderBottom.Width != nil {
                        cellStr += fmt.Sprintf(" border-bottom-width: %5.1fpt;", tstyl.BorderBottom.Width.Magnitude)
                    }
                }

                if tstyl.BorderRight != nil {
                    // Color
                    if tstyl.BorderRight.Color != nil {
                        if tstyl.BorderRight.Color.Color != nil {
                            cellStr += fmt.Sprintf(" border-right-color: %s;", util.GetColor(tstyl.BorderRight.Color.Color))
                        }
                    }
                    //dash
                    if util.GetDash(tstyl.BorderRight.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-right-style: %s;",  util.GetDash(tstyl.BorderRight.DashStyle))}
                    //Width
                    if tstyl.BorderRight.Width != nil {
                        cellStr += fmt.Sprintf(" border-right-width: %5.1fpt;", tstyl.BorderRight.Width.Magnitude)
                    }
                }
			}

			// table cell [tab1_row_col]
            if len(cellStr) >0 {
                cssStr += fmt.Sprintf(".%s_tblcel.tbc%d_%d_%d {", dObj.docName, dObj.tableCounter, trow, tcol)
                cssStr += fmt.Sprintf("%s }\n", cellStr)
                //htmlStr += fmt.Sprintf("    <td class=\"%s_tblcel tbc%d_%d_%d\">\n", dObj.docName, dObj.tableCounter, trow, tcol)
				elObj.cl2 =  fmt.Sprintf("tbc%d_%d_%d\">\n", dObj.tableCounter, trow, tcol)
            } else {
                // default
                //htmlStr += fmt.Sprintf("    <td class=\"%s_tblcel\">\n", dObj.docName)
            }

			elObj.cl1 =  fmt.Sprintf("%s_tblcel", dObj.docName)
			elObj.typ ="td"
			elObj.parent = "trow"
			elObj.newEl = "tcel"
			elObj.doAppend = true
			scriptStr += addElToDom(elObj)

			elNum := len(tcell.Content)
			for el:=0; el< elNum; el++ {
				elObj := tcell.Content[el]
				dObj.parent = "tcel"
				tObj, err:=dObj.cvtContentElToJson(elObj)
				if err != nil {
					tabObj.script = scriptStr + fmt.Sprintf("\n// error cnvtContentEl: %v\n", err)
					tabObj.bodyCss = cssStr
					return tabObj, fmt.Errorf("cvtContentElToDom - ConvertTable: %v", err)
				}
				cssStr += tObj.bodyCss
				scriptStr += tObj.script
				//htmlStr += "    " + tObj.bodyHtml
			}
			//htmlStr += "  </td>\n"
		}
		//htmlStr += "</tr>\n"
	}

	//"</tbody>\n</table>\n"
	// attach table to Dom
	dObj.parent = parent
	scriptStr += "appendEl(" + tblNam + ", " + dObj.parent +");\n"

	tabObj.script = scriptStr
	tabObj.bodyCss = cssStr
*/
	return tabStr, nil
}

func (dObj *GdocDomObj) cvtParToJson(par *docs.Paragraph)(elStr string, err error) {
// paragraph element par
// - Bullet
// - Elements
// - ParagraphStyle
// - Positioned Objects
//

	var newList cList
	var errStr string

	elStr = ""
	if par == nil {
        return "", fmt.Errorf("cvtPar -- parEl is nil!")
    }

	// Positioned Objects
/*
	numPosObj := len(par.PositionedObjectIds)
	for i:=0; i< numPosObj; i++ {
		posId := par.PositionedObjectIds[i]
		posObj, ok := dObj.doc.PositionedObjects[posId]
		if !ok {return parObj, fmt.Errorf("cvtPar: could not find positioned Object with id: ", posId)}

		imgDisp, err := dObj.renderPosImg(&posObj, posId)
		if err != nil {
			imgDisp.bodyHtml = fmt.Sprintf("<!-- error cvtPar:: render pos img %v -->\n", err) + imgDisp.bodyHtml
		}
		addDispObj(&parObj, imgDisp)
	}
*/

	if par.Bullet == nil {

		if len(par.Elements) == 1 {
      		if par.Elements[0].TextRun != nil {
            	if par.Elements[0].TextRun.Content == "\n" {
//				elStr = "{\"typ\": \"br\",\"parent\":\"" + dObj.parent +"\"},\n"
					elStr = "{\"typ\": \"br\",\"parent\":\"gdocMain\"},\n"
					return elStr, nil
            	}
        	}
		}


	// first we need to check whether this is a cr-only paragraph

// close with <br> ? or two <br>
		if dObj.listStack != nil {dObj.closeList(-1)}

		pelStr, _, err := dObj.cvtGdocParToJson(par)
		if err != nil {
			errStr = fmt.Sprintf("*** error ***: cvtParStyl: %v\n", err)
			dObj.errCount++
		}
		elStr += errStr + pelStr

		parElStr, err := dObj.cvtParElsToJson(par)
		if err != nil {
			errStr += fmt.Sprintf("**error\":  \"cvtParElDom: %v\"},\n",err)
			dObj.errCount++
		}
		elStr += errStr + parElStr


		return elStr, err
	}

	// lists
    if par.Bullet != nil {
		// there is paragraph style for each ul and a text style for each list element
// still todo
// need to apply bulletTxtMap to marker


//		if dObj.Options.Verb {
			// htnm listHtml += fmt.Sprintf("<!-- List Element %d -->\n", dObj.parCount)

		// find list id of paragraph
		listid := par.Bullet.ListId
		nestIdx := int(par.Bullet.NestingLevel)

		// retrieve the list properties from the doc.Lists map
		nestL := dObj.doc.Lists[listid].ListProperties.NestingLevels[nestIdx]
		listOrd := util.GetGlyphOrd(nestL)

		// A. check whether need new <ul> or <ol>
		// conditions for new <ul><ol>
		// 1. beginning of a list
		// 2. increase in nesting level
		// 3. different listid -> old list ended; beginning of new list

		// condition for </ul></ol>
		// 1. decrease in nesting level

//		fmt.Println("*********** listStack **********")
		fmt.Printf("listid: %s \n", listid)
		listAtt, cNest := getLiStack(dObj.listStack)
//printLiStack(dObj.listStack, "first")

//printLiStackItem(listAtt, cNest)
		listStr := ""
		parent :=""
		switch listid == listAtt.cListId {
			case true:
//fmt.Println("listid matched")

				switch {
					case nestIdx > cNest:
						// for each nest level, we have to start a new list
						for nl:=cNest + 1; nl <= nestIdx; nl++ {
							newList.cListId = listid
							newList.cOrd = listOrd
							newStack := pushLiStack(dObj.listStack, newList)
							dObj.listStack = newStack

							if listOrd {
								// html listHtml = fmt.Sprintf("<ol class=\"%s_ol nL_%d\">\n", listid[4:], nl)

								if nl > 0 {parent = fmt.Sprintf("Ol%d", nl - 1)}
								listStr = "{\"typ\":\"ol\","
								listStr += " \"parent\":\"" + parent + "\","
								cNam := fmt.Sprintf("%s Nl%d ",listid[4:] ,nl)
								listStr += " \"className\":\"" + cNam + "\","
								listStr += fmt.Sprintf(" \"name\":\"Ol%d\",",nl)
								counter := fmt.Sprintf("%sOl%d",listid[4:],nl)
								dObj.counter = counter
								listStr += " \"style\": {\"counterReset\": \"" + counter + "\"}"

								// css class: add css Rule
							} else {

								// html listHtml = fmt.Sprintf("<ul class=\"%s_ul nL_%d\">\n", listid[4:], nl)
								if nl > 0 {parent = fmt.Sprintf("Ul%d", nl - 1)}
								listStr = "{\"typ\":\"ul\","
								listStr += " \"parent\":\"" + parent + "\","
								cNam := fmt.Sprintf("%s Nl%d ",listid[4:] ,nl)
//								cNam := fmt.Sprintf("%sUl",dObj.docName)
								listStr += " \"className\":\"" + cNam + "\","
								listStr += fmt.Sprintf(" \"name\":\"Ul%d\"}",nl)
								dObj.counter = ""

								// css none
							}
						}
						// html	listHtml += fmt.Sprintf("<!-- same list increase %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)
						//html fmt.Printf("<!-- same list increase %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)

					case nestIdx < cNest:
						// html

						if listOrd {
							counter := fmt.Sprintf("%sOl%d",listid[4:],nestIdx)
							dObj.counter = counter
						}
						dObj.closeList(nestIdx)

					case nestIdx == cNest:

				}

			case false:
				// new list
				// if there is a list we need to close it
//				if cNest > -1 {dObj.closeList(cNest) ?

				newList.cListId = listid
				newList.cOrd = listOrd
				newStack := pushLiStack(dObj.listStack, newList)
				dObj.listStack = newStack
				nl := nestIdx
				parent := ""
				if nl == 0 {parent = "gdocMain"}

				if listOrd {
					// html listHtml += fmt.Sprintf("<ol class=\"%s_ol nL_%d\">\n", listid[4:], nestIdx)

					if nl > 0 {parent = fmt.Sprintf("Ol%d", nl - 1)}
					listStr = "{\"typ\":\"ol\","
					listStr += " \"parent\":\"" + parent + "\","
//					cNam := fmt.Sprintf("%sOl",dObj.docName)
					cNam := fmt.Sprintf("%s Nl%d ",listid[4:] ,nl)
//					classNam := fmt.Sprintf("%sOlNl%d",listid[4:],nl)
					listStr += " \"className\":\"" + cNam + "\","
					listStr += fmt.Sprintf(" \"name\":\"Ol%d\",",nl)
					// css
					counter := fmt.Sprintf("%sOl%d",listid[4:],nl)
					dObj.counter = counter
					listStr += " \"style\": {\"counterReset\": \"" + counter + "\"}"

//					listCss = fmt.Sprintf("  {\"cssRule\": \".%sOlNl%d {", listid[4:], nl)
//					listCss += fmt.Sprintf(" counter-reset: %sOlNl%d;}\"},\n",listid[4:], nl)

				} else {
					// html listHtml += fmt.Sprintf("<ul class=\"%s_ul nL_%d\">\n", listid[4:], nestIdx)

					if nl > 0 {parent = fmt.Sprintf("Ul%d", nl - 1)}
					listStr = "{\"typ\":\"ul\","
					listStr += " \"parent\":\"" + parent + "\","
					cNam := fmt.Sprintf("%s Nl%d ",listid[4:] ,nl)
//					cNam := fmt.Sprintf("%sUl",dObj.docName)
					listStr += " \"className\":\"" + cNam + "\","
					listStr += fmt.Sprintf(" \"name\":\"Ul%d\"",nl)
					dObj.counter = ""

					//css
				}
			default:

		} // switch
		if len(listStr) > 0 {elStr += listStr + "},\n"}

//fmt.Printf("elstr list: %s \n",elStr)

		// html <li>
		// html listPrefix = fmt.Sprintf("<li class=\"%s_li nL_%d\">", listid[4:], nestIdx)
		// script
		listParent := ""
		nl := nestIdx
		cNam := ""
		if listOrd {
			listParent = fmt.Sprintf("Ol%d", nl)
			cNam = fmt.Sprintf("%sOl%d",listid[4:],nl)
		} else {
			listParent = fmt.Sprintf("Ul%d", nl)
			cNam = fmt.Sprintf("%sUl%d",listid[4:],nl)
		}

		cNam = listid[4:] + "Li " + fmt.Sprintf("Nl%d", nl)
		counter := dObj.counter
		listStr = "{\"typ\":\"li\","
		listStr += " \"parent\":\"" + listParent + "\","
		listStr += " \"className\":\"" + cNam + "\","
		listStr += " \"name\":\"list\","
//	cssStr += " display: list-item;"
//	cssStr += " text-align: start;"
//	cssStr += " padding-left: 6pt;"

		listStr += " \"style\": {\"display\": \"listItem\", \"paddingLeft\": \"6pt\", "
		listStr += "\"counterIncrement\": \"" + counter + "\"}},\n"
		// mark

		if par.Bullet.TextStyle != nil {
//      	    bulletTxtMap := fillTxtMap(par.Bullet.TextStyle)
//			txtCss := cvtTxtMapStylToCssJson(defTxtMap, par.Bullet.TextStyle)
//			cssStr := cvtTxtMapToCssJson(bulletTxtMap)
//			cssLiRule += "  {\"cssRule\": \"." + cNam + " {" + txtCss + "}\"},\n"
		}

		elStr += listStr
//fmt.Printf("elstr list li: %s \n",elStr)

		// get paragraph
		pelStr, _, err := dObj.cvtGdocParToJson(par)
		if err != nil {
			errStr = fmt.Sprintf("{\"error\": \"cvtGdocPar: %v\"},\n", err)
			dObj.errCount++
		}
		elStr += errStr + pelStr
//fmt.Printf("elstr list par: %s \n",elStr)

		// Heading Id refers to a heading paragraph not just a normal paragraph
		// headings are bookmarked for TOC

		// par elements: text and css for text
		parElsStr, err := dObj.cvtParElsToJson(par)
		if err != nil {
			errStr = fmt.Sprintf("{\"error\": \"cvtGdocPar: %v\"},\n", err)
			dObj.errCount++
		}
		elStr += errStr + parElsStr
	}
//fmt.Printf("elstr list end: %s \n",elStr)

	return elStr, nil
}

func (dObj *GdocDomObj) cvtParElsToJson(par *docs.Paragraph)(parElsStr string, err error) {


//	addBrStr := "{\"typ\":\"br\", \"parent\":\"gdocMain\"},"

	namedTyp := par.ParagraphStyle.NamedStyleType
    numParEl := len(par.Elements)

// todo
// if numParEl = 1 and textrun -> no need to create a span element

    for pEl:=0; pEl< numParEl; pEl++ {
        parEl := par.Elements[pEl]


		if parEl.InlineObjectElement != nil {
//    	    imgObj, err := dObj.renderInlineImg(parEl.InlineObjectElement)
//       	if err != nil { }
		}

		if parEl.TextRun != nil {
//			parElStr, crEnd := dObj.cvtParTxtElToJson(parEl.TextRun, namedTyp)
			parElStr, _ := dObj.cvtParTxtElToJson(parEl.TextRun, namedTyp)
			parElsStr += parElStr
//			if crEnd {parElsStr += addBrStr}
		}

		if parEl.FootnoteReference != nil {
			dObj.ftnoteCount++
			parElsStr += dObj.cvtFtnoteToJson()
//			parDisp.bodyHtml += htmlStr
		}

		if parEl.PageBreak != nil {

		}

		if parEl.HorizontalRule != nil {
			parElsStr += dObj.cvtHrElToJson(parEl.HorizontalRule)
		}

		if parEl.ColumnBreak != nil {

		}

		if parEl.Person != nil {

		}

		if parEl.RichLink != nil {

		}

	} //loop end parEl

	return parElsStr, nil
}



func (dObj *GdocDomObj) cvtDocNamedStyles()(cssStr string, err error) {
// method that creates the css for the named Styles used in the document

	// the normal_text style are already defined in div_main
	// so the css attributes for other named styles only need to show the difference 
	// to the normal style
    normalParStyl, normalTxtStyl, err := dObj.getNamedStyl("NORMAL_TEXT")
    if err != nil {
        cssStr += fmt.Sprintf("  /* cvtNamedStyle: no NORMAL_TEST Style */\n")
    }

	ruleEndStr := "}\"},\n"
	ruleStartStr := "  {\"cssRule\": "

	for namedTyp, res := range dObj.namStylMap {

		defTxtMap := fillTxtMap(normalTxtStyl)
    	defParMap := fillParMap(normalParStyl)

		if namedTyp == "NORMAL_TEXT" { continue}
		if !res {continue}

		namParStyl, namTxtStyl, err := dObj.getNamedStyl(namedTyp)
		if err != nil {
			cssStr += fmt.Sprintf("  /* cvtNamedStyle: named Style not recognized */\n")
			continue
		}

		cssPrefix := ""
		switch namedTyp {
		case "TITLE":
			cssPrefix = fmt.Sprintf("\".%sTitle {", dObj.docName)

		case "SUBTITLE":
			cssPrefix = fmt.Sprintf("\".%sSubtitle {",dObj.docName)

		case "HEADING_1":
			cssPrefix =fmt.Sprintf("\".%sH1 {",dObj.docName)

		case "HEADING_2":
			cssPrefix =fmt.Sprintf("\".%sH2 {",dObj.docName)

		case "HEADING_3":
			cssPrefix =fmt.Sprintf("\".%sH3 {",dObj.docName)

		case "HEADING_4":
			cssPrefix =fmt.Sprintf("\".%sH4 {",dObj.docName)

		case "HEADING_5":
			cssPrefix =fmt.Sprintf("\".%sH5 {",dObj.docName)

		case "HEADING_6":
			cssPrefix =fmt.Sprintf("\".%sH6 {",dObj.docName)

		case "NORMAL_TEXT":

		case "NAMED_STYLE_TYPE_UNSPECIFIED":

		default:

		}

		if len(cssPrefix) > 0 {
			parCss, _ := cvtParMapStylToCssJson(defParMap, namParStyl, dObj.Options)
			txtCss := cvtTxtMapStylToCssJson(defTxtMap, namTxtStyl)
			cssStr += ruleStartStr + cssPrefix + parCss + txtCss + ruleEndStr
		}
	}
	return cssStr, nil
}

//todo
func (dObj *GdocDomObj) cvtGdocParToJson(par *docs.Paragraph)(parStr string, alter bool, err error) {

	var namParStyl *docs.ParagraphStyle
	// changed from Html need to handle case if parStyl == nil
	// q: is there a case where parstyl == nil
	// if parstyl == nil lets assume normal_text


	isList := true
	if (par.Bullet == nil) {isList = false}

    // NamedStyle Type
	parStyl := par.ParagraphStyle
	if parStyl == nil {
		namParStyl,_,_ = dObj.getNamedStyl("NORMAL_TEXT")
	} else {
		namedTyp := parStyl.NamedStyleType
		namParStyl, _, err = dObj.getNamedStyl(namedTyp)
		if err != nil {	return "", false, fmt.Errorf("getNamedStyl: %s not a valid name: %v", namedTyp, err)}
	}

	alter = false
	cssParAtt := ""

	parmap := fillParMap(namParStyl)

	if parStyl == nil {
		// use named style that has been published
		cssParAtt = cvtParMapToCssJson(parmap, dObj.Options)
	} else {
		cssParAtt, alter = cvtParMapStylToJson(parmap, parStyl, dObj.Options)
	}

/*
 	if len(par.Elements) == 1 {
		parElTxt := par.Elements[0].TextRun
		attStr, crEnd, err := dObj.cvtParTxtSingle(parElTxt *docs.TextRun, namedTyp string)
		if err == nil {insText = true}
	}
*/
	headingId := parStyl.HeadingId
	className := ""
	parStr = ""
	hdStr := ""

	if len(headingId) > 0 {hdStr = "\"hd\": \"" + headingId[3:] + "\","}

	parent := "gdocMain"
	if isList { parent = "list"}
	idStr := fmt.Sprintf("\"id\":\"p%d\",", dObj.parCount) + " \"name\":\"par\"," + " \"parent\":\"" + parent + "\","
	dObj.parCount++


	switch parStyl.NamedStyleType {
		case "TITLE":
			parStr = "{\"typ\":\"p\","

			if dObj.namStylMap["TITLE"] {
				//html prefix = fmt.Sprintf("<p class=\"%s_title%s\"", dObj.docName, isListClass)
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sTitle", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "}"
				}
			} else {
				parStr += "\"className\":\"title\""
			}

		case "SUBTITLE":
			parStr = "{\"typ\":\"p\","

			if dObj.namStylMap["SUBTITLE"] {
				//html prefix = fmt.Sprintf("<p class=\"%s_subtitle\"", dObj.docName)
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sSubtitle", dObj.docName)
				parStr += "\"className\":\"" + className + "\""

				if alter {
					parStr += ", \"style\":{" + cssParAtt + "}"
				}
			} else {
				parStr += "\"className\":\"subtitle\","
			}

		case "HEADING_1":
				//html prefix = fmt.Sprintf("<h1 class=\"%s_h1\"", dObj.docName)
			parStr = "{\"typ\":\"h1\","
			if dObj.namStylMap["HEADING_1"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH1", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "HEADING_2":
			//html suffix = "<h2>"
			parStr = "{\"typ\":\"h2\","
			if dObj.namStylMap["HEADING_2"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH2", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "HEADING_3":
			//html suffix = "<h3>"
			parStr = "{\"typ\":\"h3\","
			if dObj.namStylMap["HEADING_2"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH3", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "HEADING_4":
			//html suffix = "<h4>"
			parStr = "{\"typ\":\"h4\","
			if dObj.namStylMap["HEADING_4"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH4", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "HEADING_5":
			//html suffix = "<h5>"
			parStr = "{\"typ\":\"h5\","
			if dObj.namStylMap["HEADING_5"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH5", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "HEADING_6":
			//html suffix = "<h6>"
			parStr = "{\"typ\":\"h6\","
			if dObj.namStylMap["HEADING_6"] {
				parStr += idStr + hdStr
				className = fmt.Sprintf("%sH6", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				if alter {
					parStr += ", \"style\":{" + cssParAtt + "},"
				}
			}

		case "NORMAL_TEXT":
			//html prefix = fmt.Sprintf("<p class=\"%s_p\" style = {}">, dObj.docName)

			parStr = "{\"typ\":\"p\","

 			if dObj.namStylMap["NORMAL_TEXT"] {
				parStr += idStr + hdStr
				parStr += "\"name\":\"par\","
				className = fmt.Sprintf("%sPar", dObj.docName)
				parStr += "\"className\":\"" + className + "\""
				listAtt := ""
				if isList {listAtt = ", \"display\": \"inline\""}
				if alter {
					parStr += ", \"style\":{" + cssParAtt + listAtt + "}"
				} else {
					if isList { parStr += ", \"style\":{\"display\": \"inline\"}"}
				}
			}

		case "NAMED_STYLE_TYPE_UNSPECIFIED":
//			namTypValid = false
			parStr = "{\"typ\":\"p\","
			parStr += idStr + hdStr
			className = fmt.Sprintf("%sPar", dObj.docName)
			parStr += "\"className\":\"" + className + "\""

		default:
//			namTypValid = false
	}
	parStr += "},\n"

	if isList {


	}

	return parStr, alter, nil
}


func (dObj *GdocDomObj) createDivHead(divName, idStr string) (divObj dispObj, err error) {
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

//todo
func (dObj *GdocDomObj) creSecDivJson() (secStr string) {

/*

	if !dObj.Options.Sections {return nil}

	if len(dObj.sections) < 2 {return nil}

	//html
	// fmt.Sprintf("<div class=\"%s_main top\" id=\"%s_sectoc\">\n", dObj.docName, dObj.docName)
	divObj.parent = dObj.parent
	divObj.typ="div"
	divObj.newEl = "divSec"
	divObj.cl1 = fmt.Sprintf("%s_main_top", dObj.docName)
	divObj.idStr = fmt.Sprintf("%s_sectoc", dObj.docName)
	divObj.doAppend = true
	scriptStr = addElToDom(divObj)

	// fmt.Sprintf("<p class=\"%s_title %s_leftTitle_UL\">Sections</p>\n",dObj.docName, dObj.docName)
	parObj.parent = "divSec"
	parObj.typ = "p"
	parObj.newEl = "pel"
	parObj.cl1 = dObj.docName + "_title"
	parObj.cl2 = dObj.docName + "_leftTitle_UL"
	parObj.txt = "Sections"
	parObj.doAppend = true
	scriptStr += addElToDom(parObj)

	for i:=0; i< len(dObj.sections); i++ {
		// fmt.Sprintf("  <p class=\"%s_p\"><a href=\"#%s_sec_%d\">Page: %3d</a></p>\n", dObj.docName, dObj.docName, i, i)
		parObj.parent = "divSec"
		parObj.typ = "p"
		parObj.newEl = "pel"
		parObj.cl1 = dObj.docName + "_p"
		parObj.doAppend = true
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
func (dObj *GdocDomObj) creSecHeadToDom(ipage int) (secObj dispObj) {
// method that creates a distinct html dvision per section with a page heading

	var divObj, parObj elScriptObj
	var linkObj elScriptObj

	//css
	prefixCss := fmt.Sprintf(".%s_main.sec_%d {\n", dObj.docName, ipage)
	secCss := ""
	suffixCss := "}\n"

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
	divObj.doAppend = true
	secObj.script += addElToDom(divObj)

	parObj.parent = dObj.parent
	parObj.typ = "p"
	parObj.newEl = "ptop"
	parObj.cl1 = fmt.Sprintf("%s_page", dObj.docName)
	parObj.doAppend = true
	secObj.script += addElToDom(parObj)

	linkObj.parent = "ptop"
	linkObj.txt = fmt.Sprintf("Page %d", ipage)
	linkObj.href = fmt.Sprintf("%s_sectoc", dObj.docName)
	secObj.script += addLinkToDom(linkObj)
*/
	return secStr
}


func (dObj *GdocDomObj) creCssDocHeadJson() (headCss string, err error) {

	errStr :="";
	headCss = "\"cssRules\": [\n"
	ruleEndStr := "}\"},\n"
	ruleStartStr := "  {\"cssRule\": "
	cssStr := ""
    docStyl := dObj.doc.DocumentStyle
    dObj.docWidth = (docStyl.PageSize.Width.Magnitude - docStyl.MarginRight.Magnitude - docStyl.MarginLeft.Magnitude)

    //gdoc default el css and doc css
    cssStr += fmt.Sprintf("\".%sDiv {", dObj.docName)
    cssStr += fmt.Sprintf("margin: %.1fmm  %.1fmm %.1fmm %.1fmm;",docStyl.MarginTop.Magnitude*PtTomm, docStyl.MarginBottom.Magnitude*PtTomm,docStyl.MarginRight.Magnitude*PtTomm,docStyl.MarginLeft.Magnitude*PtTomm)
    if dObj.docWidth > 0 {cssStr += fmt.Sprintf("width: %.1fmm;", dObj.docWidth*PtTomm)}
//	if dObj.Options.DivBorders {cssStr += "border: 1px solid red;"}
	headCss += ruleStartStr + cssStr + ruleEndStr


	//css default text style
	cssStr = fmt.Sprintf("\".%sDiv {", dObj.docName)
	parStyl, txtStyl, err := dObj.getNamedStyl("NORMAL_TEXT")
	if err != nil {return headCss, fmt.Errorf("creHeadCss: %v", err)}

	defParMap := fillParMap(parStyl)
	defTxtMap := fillTxtMap(txtStyl)

	cssStr += "display:block;"
	if dObj.Options.DivBorders {cssStr += " border: 1px solid green;"}
//fix
	cssStr += cvtTxtMapToCssJson(defTxtMap)
	headCss += ruleStartStr + cssStr + ruleEndStr

	errStr = ""
	hdCss, err := dObj.cvtDocNamedStyles()
	if err != nil {
		errStr = fmt.Sprintf("//cvtDocNamedStyles %v\n", err)
		dObj.errCount++
	}
	headCss += errStr
	if len(hdCss) > 0 {headCss += hdCss}

	// paragraph default style
    pCssStr := cvtParMapToCssJson(defParMap, dObj.Options)
	if len(pCssStr) > 0 {
		cssStr = fmt.Sprintf("\".%sPar {", dObj.docName)
		cssStr += "margin: 0;"
		headCss += ruleStartStr + cssStr + pCssStr + ruleEndStr
	}

	// list css strings
	cssListAtt := "list-style-type: none; list-style-position: outside;"
//	cssStr = "  {\"cssRule\": \"." + dObj.docName + "Ol {" + cssListatt + "}\"},\n"
//	cssStr += "  {\"cssRule\": \"." + dObj.docName + "Ul {" + cssListatt + "}\"},\n"

	// list class
	cssStr = "  {\"cssRule\": \"." + dObj.docName + "Li {"
	cssStr += " display: list-item;"
	cssStr += " text-align: start;"
	cssStr += " padding-left: 6pt;"
	cssStr += "}\"},\n"
	headCss += cssStr

//		nestLev0 := listProp.NestingLevels[0]
//		defGlyphTxtMap := fillTxtMap(nestLev0.TextStyle)

    for i:=0; i<len(dObj.docLists); i++ {
        listid := dObj.docLists[i].listId
        listClass := listid[4:]
        listProp := dObj.doc.Lists[listid].ListProperties
		cumIndent := 0.0

		for nl:=0; nl <= int(dObj.docLists[i].maxNestLev); nl++ {
    		cssStr = ""
			nestLev := listProp.NestingLevels[nl]
			cssStr += "  {\"cssRule\":"
			glyphStr := util.GetGlyphStr(nestLev)

//			if dObj.docLists[i].ord {
				cssStr += fmt.Sprintf(" \".%s.Nl%d {", listClass, nl)
//			} else {
//				cssStr += fmt.Sprintf(" \".%sUl%d {", listClass, nl)
//			}

			idFl := nestLev.IndentFirstLine.Magnitude - cumIndent
			idSt := nestLev.IndentStart.Magnitude - cumIndent
			cssStr += cssListAtt + fmt.Sprintf(" margin: 0pt 0pt 0pt %.0fpt;", idFl)
			cssStr += fmt.Sprintf(" padding-left: %.0fpt;", idSt-idFl - 6.0)
			cssStr += fmt.Sprintf("}\"},\n")

			cumIndent += idSt

/*
			// Css <li nest level>
			cssStr += "  {\"cssRule\":"
			cssStr += fmt.Sprintf(" \".%sLiNl%d {", listClass, nl)
			if dObj.docLists[i].ord {
				cssStr += fmt.Sprintf(" counter-increment: %sLiNl%d;", listClass, nl)
			} else {
				cssStr += fmt.Sprintf(" list-style-type: %s;", glyphStr)
//					cssStr += fmt.Sprintf dObj.cvtGlyph(nestLev)
			}
			cssStr += fmt.Sprintf("}\"},\n")
*/
			// Css marker
			cssStr += "  {\"cssRule\":"
//			cssStr += fmt.Sprintf(" \".%sOl%d::marker {", listClass, nl)
			cssStr += fmt.Sprintf(" \".%sLi.Nl%d::marker {", listClass, nl)
			if dObj.docLists[i].ord {
				cssStr += fmt.Sprintf(" content: counter(%sOl%d, %s);", listClass, nl, glyphStr)
			}
// else {
// 				cssStr += fmt.Sprintf(" content: %s", glyphStr)

//
//list
//            cssStr += cvtTxtMapStylToCssJson(defTxtMap,nestLev.TextStyle)
			cssStr += fmt.Sprintf("}\"},\n")
			headCss += cssStr
		}
	}
	xlen := len(headCss) -2
	headCss = headCss[:xlen]
	headCss = headCss + "],\n"
	return headCss, nil
}

/*
   // css default table
    if dObj.tableCount > 0 {

       //css default table styling (center aligned)
        cssStr = fmt.Sprintf(".%s_tbl {\n", dObj.docName)
        cssStr += "  width: 100%;\n"
        cssStr += "  border-collapse: collapse;\n"
        cssStr += "  border: 1px solid black;\n"
        cssStr += "  margin-left: auto;  margin-right: auto;\n"
        cssStr += "}\n"

		//css table row
        cssStr += fmt.Sprintf(".%s_tblrow {\n", dObj.docName)
		cssStr += "  min-height: 1em;\n"
        cssStr += "}\n"

        //css table cell
        cssStr += fmt.Sprintf(".%s_tblcel {\n", dObj.docName)
		cssStr += "  border-collapse: collapse;\n"
        cssStr += "  border: 1px solid black;\n"
//      cssStr += "  margin:auto;\n"
        cssStr += "  padding: 0.5pt;\n"
		cssStr += "  height: 1em;\n"
        cssStr += "}\n"

		// add Css
		headCss += cssStr
    }


//	xlen := len(headCss)-1
//	fmt.Printf("headCss last %q %q %q\n", headCss[xlen-2], headCss[xlen-1], headCss[xlen])

//	if headCss[xlen-1] == ',' {headCss = headCss[:xlen-1]}
//	headCss += "],\n"
	return headCss, nil
}
*/

func (dObj *GdocDomObj) cvtContentElToJson(contEl *docs.StructuralElement) (elStr string, err error) {
// method that parses a Structural Element and invokes further methods

	if dObj == nil {
		return "", fmt.Errorf("error -- dObj is nil")
	}
//	parent = dObj.eldiv

	elStr = ""

	if contEl.Paragraph != nil {
		parEl := contEl.Paragraph
		parElStr, err := dObj.cvtParToJson(parEl)
		if err != nil { return parElStr, fmt.Errorf("error par %v\n", err) }
		elStr += parElStr
		return elStr, err
	}

	if contEl.SectionBreak != nil {
//		secStr := "{"
//		errStr = "\"comment\":\"section not implemented\""
//		elStr += secStr + errStr + "},\n"
		return "", nil
	}

	if contEl.Table != nil {
		tableEl := contEl.Table
		tabStr, err := dObj.cvtTableToJson(tableEl)
		if err != nil { return tabStr, fmt.Errorf("error table %v\n", err) }
		elStr = tabStr
		return elStr, err
	}

	if contEl.TableOfContents != nil {
//		errStr = "\"comment\":\"toc not implemented\""
//		elStr = errStr
		return "", nil
	}

	err = fmt.Errorf("no contEl found!")
	return "", err
}

//footnote div
func (dObj *GdocDomObj) creFtnoteDivDom () (ftnoteStr string, err error) {
/*
	var ftnDiv dispObj
	var htmlStr, cssStr, scriptStr string
	var jselObj elScriptObj

	doc := dObj.doc
	if len(dObj.docFtnotes) == 0 {
		return nil, nil
	}

	if len(dObj.docFtnotes) == 0 {return nil, nil}

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
	jselObj.doAppend = true
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
	jselObj.doAppend = true
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
	jselObj.doAppend = true
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
	htmlStr = ""
	cssStr = ""
	for iFtn:=0; iFtn<len(dObj.docFtnotes); iFtn++ {
		idStr := dObj.docFtnotes[iFtn].id
		docFt, ok := doc.Footnotes[idStr]
		if !ok {
			htmlStr += fmt.Sprintf("<!-- error ftnote %d not found! -->\n", iFtn)
			continue
		}
		//htmlStr = fmt.Sprintf("<!-- FTnote: %d %s els: %d -->\n", iFtn, idStr, len(docFt.Content))
		//htmlStr +="<li>\n"

		// script
		jselObj.parent = "ft_OL"
		jselObj.typ = "li"
		//  jselObj.cl1 = dObj.docName + "_ftnOL"
		//	jselObj.cl2 = dObj.docName + "_ftndiv"
		jselObj.newEl = "liEl"
		//	jselObj.txt = "Footnotes"
		jselObj.doAppend = true
		scriptStr += addElToDom(jselObj)


		// presumably footnotes are paragraphs only
		for el:=0; el<len(docFt.Content); el++ {
			cssStr = ""
			elDocObj := docFt.Content[el]
			if elDocObj.Paragraph == nil {continue}
			par := elDocObj.Paragraph
			pidStr := idStr[5:]

			//html htmlStr += fmt.Sprintf("<p class=\"%s_p %s_pft\" id=\"%s\">\n", dObj.docName, dObj.docName, pidStr)

			jselObj.parent = "liEl"
			jselObj.typ = "p"
			jselObj.cl1 = dObj.docName + "_p"
			jselObj.cl2 = dObj.docName + "_pft"
			jselObj.idStr = pidStr
			jselObj.newEl = "pliEl"
			jselObj.doAppend = true
			scriptStr += addElToDom(jselObj)

			dObj.parent = "pliEl"
			tDisp, err := dObj.cvtParElToDom(par)
			if err != nil {scriptStr += fmt.Sprintf("// *** error cvtParElToDom el: %d %v\n", el, err)}
			ftnDiv.bodyCss += tDisp.bodyCss
			scriptStr += tDisp.script

//			ftnDiv.bodyHtml += htmlStr
		}

		ftnDiv.script += scriptStr
	}
*/
	return ftnoteStr, nil
}

//toc div
func (dObj *GdocDomObj) creJsonTocDiv () (tocStr string, err error) {

/*
	if dObj.Options.Toc != true { return nil, nil }

	if dObj.Options.Verb {
		if len(dObj.headings) < 2 {
			fmt.Printf("*** no TOC insufficient headings: %d ***\n", len(dObj.headings))
		} else {
			fmt.Printf("*** creating TOC Div ***\n")
		}
	}

	if len(dObj.headings) < 2 {
		tocDiv.bodyHtml = fmt.Sprintf("<!-- no toc insufficient headings -->")
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
	elObj.doAppend = true
	tocDiv.script = addElToDom(elObj)

	//fmt.Sprintf("<p class=\"%s_title %s_leftTitle\">Table of Contents</p>\n", dObj.docName, dObj.docName)
	elObj.parent = "divToc"
	elObj.typ = "p"
	elObj.newEl = "divToc"
	elObj.cl1 = dObj.docName + "_title"
	elObj.cl2 = dObj.docName + "_leftTitle"
	elObj.txt = "Table of Contents"
	elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
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
			elObj.doAppend = true
			tocDiv.script += addElToDom(elObj)
			elObj.parent = "parel"
			elObj.txt = text
			elObj.href = "#" + hdId
			tocDiv.script += addLinkToDom(elObj)

		case "NORMAL_TEXT":

		default:

		}

	} // end loop
*/
	return tocStr, nil
}

func (dObj *GdocDomObj) cvtBodyToJson() (jsonStr string, err error) {

	var errStr string

	if dObj == nil {
		return "", fmt.Errorf("-- no GdocObj!")
	}

	doc := dObj.doc
	body := doc.Body
	if body == nil {
		return "", fmt.Errorf("-- no doc.body!")
	}

	err = nil


//	bodyObj.bodyHtml = fmt.Sprintf("<div class=\"%s_main\">\n", dObj.docName)
/*
	var divMain elScriptObj
	divMain.comment = "create main div"
	divMain.typ = "div"
	divMain.parent = "divDoc"
	divMain.cl1 = dObj.docName + "_main"
	dObj.parent = "divMain"
	divMain.newEl = dObj.parent
	divMain.doAppend = true
	bodyObj.script = addElToDom(divMain)
*/

	jsonStr = "\"elements\": ["
	// divMain
	classNam := dObj.docName + "Div"
	dObj.parent = dObj.docName + "Main"
	elStr := fmt.Sprintf("{\"typ\":\"div\",\"className\":\"%s\",\"id\":\"%sMain\",\"name\":\"gdocMain\"},\n",classNam, dObj.docName)

	jsonStr +=elStr

	elNum := len(body.Content)
	for el:=0; el< elNum; el++ {
		bodyEl := body.Content[el]
		elstr, err1 := dObj.cvtContentElToJson(bodyEl)
		if err1 != nil {
			errStr = fmt.Sprintf("error: el %d cvtContentEl: %v", el, err1)
			fmt.Printf("*** %s\n", errStr)
//			err = fmt.Errorf("cvtContentEl: El %d %v\n", el, err)
			return jsonStr, fmt.Errorf("cvtContentElToJson: %v", err1)
		}
		jsonStr += elstr
	} // for el loop end

	if dObj.listStack != nil {dObj.closeList(-1)}

	ilen := len(jsonStr)
	if ilen > 0 { jsonStr = jsonStr[:ilen-2]}
	jsonStr += "]"

//	xlen := len(cssRuleSet)-1
//	if cssRuleSet[xlen-1] == ',' {cssRuleSet = cssRuleSet[:xlen-1]}
//	cssRuleSet += "],\n"


	return jsonStr, err
}

/*
func (dObj *GdocDomObj) cvtBodySecToDom(elSt, elEnd int) (bodyObj *dispObj, err error) {

	if dObj == nil {
		return nil, fmt.Errorf("cvtBody -- no GdocObj!")
	}


	doc := dObj.doc
	body := doc.Body
	if body == nil {
		return nil, fmt.Errorf("cvtBody -- no body!")
	}

	elCount := len(body.Content)

	if elSt > elCount { return nil, fmt.Errorf("cvtSelBody: elSt > elCount!") }
	if elEnd > elCount { return nil, fmt.Errorf("cvtSelBody: elEnd > elCount!") }
	if elSt > elEnd { return nil, fmt.Errorf("cvtSelBody: elSt > elElnd!") }

//	toc := dObj.Options.Toc
	bodyObj = new(dispObj)

	bodyObj.script = ""

	for el:=elSt; el<= elEnd; el++ {
		bodyEl := body.Content[el]
		tObj, err:=dObj.cvtContentElToDom(bodyEl)
		if err != nil {
			tObj.script += fmt.Sprintf("// error el %d cvtContentElToDom: %v\n", el, err)
		}
		addDispObj(bodyObj, tObj)
	} // for el loop end

	if dObj.listStack != nil {dObj.closeList(-1)}

	return bodyObj, nil
}
*/

func CreGdocDomDoc(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
	// function which converts the entire document into an hmlt file

	return nil
}

func CreGdocDomMain(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that converts the main part of a gdoc document into an html file
// excludes everything before the "main" heading or
// excludes sections titled "summary" and "keywords"
	return nil
}


func CreGdocJsonSection(heading, folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that creates a J from the named section

	var dObj GdocDomObj
	var elementsStr, cssListClasses string
	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization of dObj
	err = dObj.initGdocJson(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

/*
// footnotes
	ftnoteDiv, err := dObj.creFtnoteDivJson()
	if err != nil {
		fmt.Errorf("creFtnoteDivDom: %v", err)
	}
*/
//todo
//	dObj.sections
	secDivStr := dObj.creSecDivJson()
	if len(secDivStr)>0 {
/*
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.creSecHeadToJson(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBodyStr, err := dObj.cvtBodySecToJson(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
		}
*/

	} else {
		elementsStr, err = dObj.cvtBodyToJson()
		if err != nil {
			return fmt.Errorf("cvtBody: %v", err)
		}
	}


	//css for document head
	cssClasses, err := dObj.creCssDocHeadJson()
	if err != nil {
		return fmt.Errorf("creCssDocHead: %v", err)
	}

/*
	//html + css for Toc Div
	tocDiv, err := dObj.creTocDivDom()
	if err != nil {
		tocDiv.bodyHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
	}
*/
	//get html file pointer
	outfil := dObj.jsonFil
	if outfil == nil {
		return fmt.Errorf("outfil is nil!")
	}

	// creates json string
	docHeadStr := creJsonHead(dObj.docName)
	outfil.WriteString(docHeadStr)


	//css default css of document and document dimensions
	outfil.WriteString(cssClasses)
	outfil.WriteString(cssListClasses)


	// css of body elements
	outfil.WriteString(elementsStr)

	// closes json Object
	outfil.WriteString("}")

	outfil.Close()
	return nil
}


func CreGdocJsonAll(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that creates an html fil from the named section

	var dObj GdocDomObj
	var elementsStr, cssListClasses string

	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization of dObj
	err = dObj.initGdocJson(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

/*
// footnotes
	ftnoteDiv, err := dObj.creFtnoteDivJson()
	if err != nil {
		fmt.Errorf("creFtnoteDivDom: %v", err)
	}
*/
//todo
//	dObj.sections
	secDivStr := dObj.creSecDivJson()
	if len(secDivStr)>0 {
/*
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.creSecHeadToJson(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBodyStr, err := dObj.cvtBodySecToJson(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
		}
*/

	} else {
		elementsStr, err = dObj.cvtBodyToJson()
		if err != nil {
			return fmt.Errorf("cvtBody: %v", err)
		}
	}


	//css for document head
	cssClasses, err := dObj.creCssDocHeadJson()
	if err != nil {
		return fmt.Errorf("creCssDocHead: %v", err)
	}

/*
	//html + css for Toc Div
	tocDiv, err := dObj.creTocDivDom()
	if err != nil {
		tocDiv.bodyHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
	}
*/
	//get html file pointer
	outfil := dObj.jsonFil
	if outfil == nil {
		return fmt.Errorf("outfil is nil!")
	}

	// creates json string
	docHeadStr := creJsonHead(dObj.docName)
	outfil.WriteString(docHeadStr)


	//css default css of document and document dimensions
	outfil.WriteString(cssClasses)
	outfil.WriteString(cssListClasses)


	// css of body elements
	outfil.WriteString(elementsStr)

	// closes json Object
	outfil.WriteString("}")

	outfil.Close()
	return nil
}


