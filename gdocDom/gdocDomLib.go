// golang library that creates a html / js file from a gdoc file
// author: prr
// created: 2/5/2022
// copyright 2022 prr, Peter Riemenschneider
//
// for changes see github
//
// start: CreGdocToDom
//

package gdocDom

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
    ImgCount int
    tableCount int
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

func cvtTxtMapStylCss (txtMap *textMap, txtStyl *docs.TextStyle)(cssStr string) {

    if (len(txtStyl.BaselineOffset) > 0) && (txtStyl.BaselineOffset != "BASELINE_OFFSET_UNSPECIFIED") {
        if txtStyl.BaselineOffset != txtMap.baseOffset {
            txtMap.baseOffset = txtStyl.BaselineOffset
            switch txtMap.baseOffset {
            case "SUPERSCRIPT":
                cssStr += "  vertical-align: sub;\n"
            case "SUBSCRIPT":
                cssStr += " vertical-align: sup;\n"
            case "NONE":
                cssStr += " vertical-align: baseline;\n"
            default:
            //error
                cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
            }
        }
    }

    switch {
    case txtStyl.Bold && (txtMap.fontWeight < 700):
        txtMap.fontWeight = 800
        cssStr += fmt.Sprintf("  font-weight: %d;\n", txtMap.fontWeight)
    case !txtStyl.Bold && (txtMap.fontWeight > 500):
        txtMap.fontWeight = 400
        cssStr += fmt.Sprintf("  font-weight: %d;\n", txtMap.fontWeight)
    default:

    }

    if txtStyl.Italic != txtMap.italic {
        txtMap.italic = txtStyl.Italic
        if txtMap.italic {
            cssStr += "  font-style: italic;\n"
        } else {
            cssStr += "  font-style: normal;\n"
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

    if len(txtprop) > 0 {cssStr += fmt.Sprintf("  text-decoration: %s;\n", txtprop)}

    if txtStyl.WeightedFontFamily != nil {
        if txtStyl.WeightedFontFamily.FontFamily != txtMap.fontType {
            txtMap.fontType = txtStyl.WeightedFontFamily.FontFamily
            cssStr += fmt.Sprintf("  font-family: %s;\n", txtMap.fontType)
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
            cssStr += fmt.Sprintf("  font-size: %.2fpt;\n", txtMap.fontSize)
        }
    }

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            color := util.GetColor(txtStyl.ForegroundColor.Color)
            if color != txtMap.txtColor {
                txtMap.txtColor = color
                cssStr += fmt.Sprintf("  color: %s;\n", txtMap.txtColor)
            }
        }
    }

    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            color := util.GetColor(txtStyl.BackgroundColor.Color)
            if color != txtMap.bckColor {
                txtMap.bckColor = color
                cssStr += fmt.Sprintf("  background-color: %s;\n", txtMap.bckColor)
            }
        }
    }

    return cssStr
}

func cvtTxtStylCss (txtStyl *docs.TextStyle)(cssStr string) {
    var tcssStr string

    if len(txtStyl.BaselineOffset) > 0 {
        valStr := "vertical-align: "
        switch txtStyl.BaselineOffset {
            case "SUPERSCRIPT":
                valStr += "sub"
            case "SUBSCRIPT":
                valStr += "sup"
            case "NONE":
                valStr += "baseline"
            default:
                valStr = fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtStyl.BaselineOffset)
        }
        tcssStr = valStr + ";\n"
    }

    if txtStyl.Bold {
        tcssStr += "  font-weight: 800;\n"
    } else {
        tcssStr += "  font-weight: 400;\n"
    }

    if txtStyl.Italic { tcssStr += "  font-style: italic;\n"}
    if txtStyl.Underline { tcssStr += "  text-decoration: underline;\n"}
    if txtStyl.Strikethrough { tcssStr += "  text-decoration: line-through;\n"}

    if txtStyl.WeightedFontFamily != nil {
        font := txtStyl.WeightedFontFamily.FontFamily
        tcssStr += fmt.Sprintf("  font-family: %s;\n", font)
    }

    if txtStyl.FontSize != nil {
        mag := txtStyl.FontSize.Magnitude
        tcssStr += fmt.Sprintf("  font-size: %.2fpt;\n", mag)
    }

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            //0 to 1
            tcssStr += "  color: "
            tcssStr += util.GetColor(txtStyl.ForegroundColor.Color)
        }
    }
    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            tcssStr += "  background-color: "
            tcssStr += util.GetColor(txtStyl.BackgroundColor.Color)
        }
    }

    if len(tcssStr) > 0 {
        cssStr = tcssStr
    }
    return cssStr
}

func printParMap(parmap *parMap, parStyl *docs.ParagraphStyle) {

	alter := false
	fmt.Printf("*** align ***\n")
	if parStyl.Alignment != parmap.halign {
		fmt.Printf("align: %s %s \n", parmap.halign, parStyl.Alignment)
		parmap.halign = parStyl.Alignment
		alter = true
	}
	fmt.Printf("align: %s \n", parmap.halign)

//	parmap.direct = true
	fmt.Printf("*** indent ***\n")
	if (parStyl.IndentStart != nil) {
		if parStyl.IndentStart.Magnitude != parmap.indStart {
			fmt.Printf("indent start: %.1fpt %.1fpt\n", parmap.indStart, parStyl.IndentStart.Magnitude)
			parmap.indStart = parStyl.IndentStart.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent start: %.1fpt\n", parmap.indStart)

	fmt.Printf("*** indent end***\n")
	if (parStyl.IndentEnd != nil) {
		if parStyl.IndentEnd.Magnitude != parmap.indEnd {
			fmt.Printf("indent end: %.1f %.1f \n", parmap.indEnd, parStyl.IndentEnd.Magnitude)
			parmap.indEnd = parStyl.IndentEnd.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent end: %.1fpt\n", parmap.indEnd)

	fmt.Printf("*** indent first line ***\n")
	if (parStyl.IndentFirstLine != nil) {
		if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
			fmt.Printf("indent first line: %.1f %.1f \n", parmap.indFlin, parStyl.IndentFirstLine.Magnitude)
			parmap.indFlin = parStyl.IndentFirstLine.Magnitude
			alter = true
		}
	}
	fmt.Printf("indent first line: %.1fpt\n", parmap.indFlin)

	fmt.Printf("*** line spacing ***\n")
	if parStyl.LineSpacing/100.0 != parmap.linSpac {
		fmt.Printf("line spacing: %.2f %.2f \n", parmap.linSpac, parStyl.LineSpacing/100.0)
		parmap.linSpac = parStyl.LineSpacing/100.0; alter = true;
	}
	fmt.Printf("line spacing: %.2fpt\n", parmap.linSpac)

	fmt.Printf("*** keep lines ***\n")
	if parStyl.KeepLinesTogether != parmap.keepLines {
		fmt.Printf("keep Lines: %t %t\n", parmap.keepLines, parStyl.KeepLinesTogether)
		parmap.keepLines = parStyl.KeepLinesTogether; alter = true;
	}
	fmt.Printf("keep Lines: %t\n", parmap.keepLines)

	fmt.Printf("*** keep next ***\n")
	if parStyl.KeepWithNext != parmap.keepNext {
		fmt.Printf("keep With: %t %t\n", parmap.keepNext, parStyl.KeepWithNext)
		parmap.keepNext = parStyl.KeepWithNext; alter = true;
	}
	fmt.Printf("keep With: %t\n", parmap.keepNext)

	fmt.Printf("*** space above ***\n")
	if (parStyl.SpaceAbove != nil) {
		if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
			fmt.Printf("space above: %.1fpt %.1fpt\n", parmap.spaceTop, parStyl.SpaceAbove.Magnitude)
			parmap.spaceTop = parStyl.SpaceAbove.Magnitude
			alter = true
		}
	}
	fmt.Printf("space above: %.1fpt\n", parmap.spaceTop)

	fmt.Printf("*** space below ***\n")
	if (parStyl.SpaceBelow != nil) {
		if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
			fmt.Printf("space below: %.1f %.1f \n", parmap.spaceBelow, parStyl.SpaceBelow.Magnitude)
			parmap.spaceBelow = parStyl.SpaceBelow.Magnitude
			alter = true
		}
	}
	fmt.Printf("space below: %.1fpt\n", parmap.spaceBelow)

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
       fmt.Printf("spacing mode: %t %t \n", parmap.spaceMode, spaceMode)
        parmap.spaceMode = spaceMode
	}
	fmt.Printf("spacing mode: %s\n", parmap.spaceMode)

	//tabs to do
//	parmap.hasBorders = true

	bb := true
	bb = bb && (parStyl.BorderBetween == nil)
	bb = bb && (parStyl.BorderTop == nil)
	bb = bb && (parStyl.BorderRight == nil)
	bb = bb && (parStyl.BorderBottom == nil)
	bb = bb && (parStyl.BorderLeft == nil)
	if bb {
		fmt.Printf("has no borders: %t %t \n", parmap.hasBorders, !bb)
		parmap.hasBorders = false
		fmt.Printf("alter 1: %t\n", alter)
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
				parmap.bordBet.width = parStyl.BorderBetween.Width.Magnitude
				alter = true
				fmt.Printf("width: %.1f\n",parmap.bordBet.width)
			}
		}
		if parStyl.BorderBetween.Padding != nil {
			if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
				parmap.bordBet.pad = parStyl.BorderBetween.Padding.Magnitude
				alter = true
				fmt.Printf("padding: %.1f\n",parmap.bordBet.pad)
			}
		}
		if parStyl.BorderBetween.Color != nil {
			if parStyl.BorderBetween.Color.Color != nil {
				color := util.GetColor(parStyl.BorderBetween.Color.Color)
				if color != parmap.bordBet.color {
					parmap.bordBet.color = color
					alter = true
					fmt.Printf("color: %s\n",parmap.bordBet.color)
				}
			}
		}
		if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {
			parmap.bordBet.dash = parStyl.BorderBetween.DashStyle;
			alter = true;
			fmt.Printf("dash: %s\n",parmap.bordBet.dash)
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

func cvtParMapStylCss(parmap *parMap, parStyl *docs.ParagraphStyle, opt *util.OptObj)(cssStr string, alter bool) {
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
                cssStr += "  text-align: left;\n"
            case "CENTER":
                cssStr += "  text-align: center;\n"
            case "END":
                cssStr += "  text-align: right;\n"
            case "JUSTIFIED":
                cssStr += "  text-align: justify;\n"
            default:
                cssStr += fmt.Sprintf("/* unrecognized Alignment %s */\n", parmap.halign)
        }

    }

    // test direction skip for now
    parmap.direct = true

    if (parStyl.IndentFirstLine != nil) {
        if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
            parmap.indFlin = parStyl.IndentFirstLine.Magnitude
            cssStr += fmt.Sprintf("  text-indent: %.1fpt;\n", parmap.indFlin)
			alter = true
        }
    }

    parmap.linSpac = 1.0
    if parStyl.LineSpacing/100.0 != parmap.linSpac {
        if parStyl.LineSpacing > 1.0 {
			alter = true
            parmap.linSpac = parStyl.LineSpacing/100.0
            if opt.DefLinSpacing > 0.0 {
                cssStr += fmt.Sprintf("  line-height: %.2f;\n", opt.DefLinSpacing*parmap.linSpac)
            } else {
                cssStr += fmt.Sprintf("  line-height: %.2f;\n", parmap.linSpac)
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
		cssStr += fmt.Sprintf("  margin: %.0f %.0f %.0f %.0f;\n", tmarg, rmarg, bmarg, lmarg)
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

    cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", parmap.bordTop.pad, parmap.bordRight.pad, parmap.bordBot.pad, parmap.bordLeft.pad)
    cssStr += fmt.Sprintf("  border-top: %.1fpt %s %s;\n", parmap.bordTop.width, util.GetDash(parmap.bordTop.dash), parmap.bordTop.color)
    cssStr += fmt.Sprintf("  border-right: %.1fpt %s %s;\n", parmap.bordRight.width, util.GetDash(parmap.bordRight.dash), parmap.bordRight.color)
    cssStr += fmt.Sprintf("  border-bottom: %.1fpt %s %s;\n", parmap.bordBot.width, util.GetDash(parmap.bordBot.dash), parmap.bordBot.color)
    cssStr += fmt.Sprintf("  border-left: %.1fpt %s %s;\n", parmap.bordLeft.width, util.GetDash(parmap.bordLeft.dash), parmap.bordLeft.color)

    return cssStr, alter
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
						//tblCell add text
		jsStr += "      tblRow.appendChild(tblCell);\n"
		jsStr += "	  }/n"
		jsStr += "	  tblBody.appendChild(tblRow);\n"
		jsStr += "	}/n"

		jsStr += "  tblp = tblObj.parent;\n"
		jsStr += "  tblp.appendChild(tab);\n"
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

func addTblElToDom(tableObj tableScriptObj)(script string) {

	script = "// addEl \n"
	script += "// " + tableObj.comment + "\n"
	if !(len(tableObj.parent) > 0) {
		script += "// error - no el parent provided!\n"
		return script
	}
	script = "  for (key in tableObj) {tableObj[key] = null;}\n"
	if len(tableObj.cl1) > 0 {script += fmt.Sprintf("  tableObj.cl1 = '%s';\n", tableObj.cl1)}
	if len(tableObj.cl2) > 0 {script += fmt.Sprintf("  tableObj.cl2 = '%s';\n", tableObj.cl2)}
	if len(tableObj.idStr) > 0 {script += fmt.Sprintf("  tableObj.idStr = '%s';\n", tableObj.idStr)}
	script += fmt.Sprintf("  tableObj.parent = %s;\n", tableObj.parent)
	script += fmt.Sprintf("  tbl = addTblEl(tableObj);\n")
	for irow:=0; irow < tableObj.rowCount; irow++ {
		for icol:=0; icol < tableObj.colCount; icol++ {

		}
	}
	script += fmt.Sprintf("  fillTblEl(tableObj);\n")

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

func cvtText(inp string) (out string) {
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


func (dObj *GdocDomObj) cvtParMapCss(pMap *parMap)(cssStr string) {
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

func (dObj *GdocDomObj) initGdocDom(folderPath string, options *util.OptObj) (err error) {
	var listItem docList
	var heading headingTyp
	var sec secTyp
	var ftnote docFtnoteTyp
	var docPb pbTyp

	doc := dObj.doc
	if doc == nil {return fmt.Errorf("doc is nil in GdocDomObj!")}

	// need to transform file name
	// replace spaces with underscore
	dNam := doc.Title
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
//				hdlen := len(dObj.headings)
//				if hdlen > 1 {
//					dObj.headings[hdlen-2].hdElEnd = parHdEnd
//				}
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


func (dObj *GdocDomObj) cvtParElTextold(parElTxt *docs.TextRun)(htmlStr string, cssStr string, err error) {

   if parElTxt == nil {
        return "","", fmt.Errorf("cvtPelText -- parElTxt is nil!")
    }

	// need to check whether <1
	if len(parElTxt.Content) < 2 { return "","",nil}

	// need to compare text style with the default style

	spanCssStr := cvtTxtStylCss(parElTxt.TextStyle)

	linkPrefix := ""
	linkSuffix := ""
	if parElTxt.TextStyle.Link != nil {
		if len(parElTxt.TextStyle.Link.Url)>0 {
			linkPrefix = "<a href = \"" + parElTxt.TextStyle.Link.Url + "\">"
			linkSuffix = "</a>"
		}
	}
	if len(spanCssStr)>0 {
		dObj.spanCount++
		spanIdStr := fmt.Sprintf("%s_sp%d", dObj.docName, dObj.spanCount)
		cssStr = fmt.Sprintf("#%s {\n", spanIdStr) + spanCssStr + "}\n"
		htmlStr = fmt.Sprintf("<span id=\"%s\">",spanIdStr) + linkPrefix + parElTxt.Content + linkSuffix + "</span>"
	} else {
		htmlStr = linkPrefix + parElTxt.Content + linkSuffix
	}
	return htmlStr, cssStr, nil
}

func (dObj *GdocDomObj) cvtParDomElText(parElTxt *docs.TextRun, namedTyp string)(parTxt dispObj) {
	var scriptStr, cssStr, spanCssStr string
	var spanEl, linkEl, txtEl elScriptObj

	if parElTxt == nil {
		parTxt.script = "//cvtPelText -- parElTxt is nil!"
		return parTxt
	}
	if !(len(parElTxt.Content) > 0)  {
		parTxt.script = "//cvtPelText -- no Content!"
		return parTxt
	}
	if !(len(namedTyp) >0) {
		parTxt.script = "//cvtPelText -- no Named Type!"
		namedTyp = "NORMAL_TEXT"
	}

//	namedTyp := parStyl.NamedStyleType

	_, namedTxtStyl, err := dObj.getNamedStyl(namedTyp)
	if err != nil {
		parTxt.script = fmt.Sprintf("//cvtPelText -- invalid Named Type! %v", err)
		return parTxt
	}

	txtMap := fillTxtMap(namedTxtStyl)

	spanCssStr = cvtTxtMapStylCss(txtMap, parElTxt.TextStyle)

	if len(spanCssStr) > 0 {
	// create span element
		dObj.spanCount++
		spanIdStr := fmt.Sprintf("%s_sp%d", dObj.docName, dObj.spanCount)
		// script
		spanEl.comment = "spanEl"
		spanEl.parent = "hdel"
		spanEl.idStr = spanIdStr
		spanEl.typ = "span"
		spanEl.newEl = "spanEl"
//		spanEl.txt = cvtText(parElTxt.Content)
		scriptStr = addElToDom(spanEl)
		//css
		cssStr = fmt.Sprintf("#%s {\n", spanIdStr) + spanCssStr + "}\n"

	// need to create an a element

	} else {
		spanEl.comment = "spanEl"
		spanEl.parent = "hdel"
		spanEl.typ = "span"
		spanEl.newEl = "spanEl"
		scriptStr = addElToDom(spanEl)
		//css
		cssStr = ""
	}

	if parElTxt.TextStyle.Link != nil {
		linkEl.parent = "spanEl"
		linkEl.typ = "a"
		linkEl.newEl = "anchor"
		linkEl.href = parElTxt.TextStyle.Link.Url
		linkEl.txt = parElTxt.Content
		scriptStr += addElToDom(linkEl)
	} else {
		txtEl.parent = "spanEl"
		txtEl.txt = cvtText(parElTxt.Content)
		scriptStr += addTxtElToDom(txtEl)
	}

	parTxt.bodyCss = cssStr
	parTxt.script = scriptStr
	return parTxt
}


func (dObj *GdocDomObj) closeList(nl int) {
	// ends a list

	if (dObj.listStack == nil) {return}

	stack := dObj.listStack
	n := len(*stack)

	for i := n -1; i >= nl; i-- {
//		ord := (*stack)[i].cOrd
//		if ord {
//			htmlStr += "</ol>\n"
//		} else {
//			htmlStr +="</ul>\n"
//		}
		nstack := popLiStack(stack)
		dObj.listStack = nstack
	}
	return
}

func (dObj *GdocDomObj) cvtHrElToDom (hr *docs.HorizontalRule)(hrObj dispObj) {
    var cssStr string
	var hrEl elScriptObj
    //html htmlStr = "<hr>\n"
	// script
	hrEl.parent = dObj.parent
	hrEl.typ = "hr"
	hrEl.newEl = "hrEl"
    if hr.TextStyle != nil {
        cssStr = fmt.Sprintf(".%s_hr_%d {\n", dObj.docName, dObj.hrCount)
        cssStr += cvtTxtStylCss(hr.TextStyle)
        cssStr += "}\n"
		// html  fmt.Sprintf("<hr class=\"%s_hr_%d\">\n", dObj.docName, dObj.hrCount)
		hrEl.cl1 = fmt.Sprintf("%s_hr_%d", dObj.docName, dObj.hrCount)
    }

    hrObj.script = addElToDom(hrEl)
    hrObj.bodyCss = cssStr
    return hrObj
}

func (dObj *GdocDomObj) renderInlineImg(imgEl *docs.InlineObjectElement)(imgDisp *dispObj, err error) {
	var imgDispObj dispObj
	var imgDomEl imgScriptObj

	if imgEl == nil {
		return nil, fmt.Errorf("imgEl is nil!")
	}
	doc := dObj.doc

	imgElId := imgEl.InlineObjectId
	if !(len(imgElId) > 0) {return nil, fmt.Errorf("no InlineObjectId found!")}

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

	if dObj.Options.ImgFold {
    	imgSrc := dObj.imgFoldNam + "/" + imgId + ".jpeg"
		// html htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgSrc, imgId, imgObj.Title)
		imgDomEl.src = imgSrc
	} else {
		// html htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgObj.ImageProperties.SourceUri, imgId, imgObj.Title)
		imgDomEl.src = imgObj.ImageProperties.SourceUri
	}
	cssStr := fmt.Sprintf("#%s {\n",imgId)
	cssStr += fmt.Sprintf(" width:%.0fpt; height:%.0fpt; \n}\n", imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )
	// todo add margin
	imgDomEl.parent = dObj.parent
	imgDomEl.idStr = imgId
	imgDomEl.width = int(imgObj.Size.Width.Magnitude)
	imgDomEl.height = int(imgObj.Size.Height.Magnitude)
	if len(imgObj.Title) > 0 {imgDomEl.title = imgObj.Title}
	if len(imgObj.Description) > 0 {imgDomEl.desc = imgObj.Description}

	imgDispObj.script = addImgElToDom(imgDomEl)
	imgDispObj.bodyCss = cssStr
	return &imgDispObj, nil
}


func (dObj *GdocDomObj) renderPosImg(posImg *docs.PositionedObject, posId string)(imgDisp *dispObj, err error) {
	var imgDispObj dispObj

	// html
	posObjProp := posImg.PositionedObjectProperties
	imgProp := posObjProp.EmbeddedObject
	htmlStr := fmt.Sprintf("\n<!-- Positioned Image %s -->\n", posId)
	imgDivId := fmt.Sprintf("%s_%s", dObj.docName, posId[4:])
	imgId := imgDivId + "_img"
	pimgId := imgDivId +"_p"

	layout := posObjProp.Positioning.Layout
	topPos := posObjProp.Positioning.TopOffset.Magnitude
	leftPos := posObjProp.Positioning.LeftOffset.Magnitude
	fmt.Printf("layout %s top: %.1fmm left:%.1fmm\n", layout, topPos*PtTomm, leftPos*PtTomm)

	imgSrc := imgProp.ImageProperties.ContentUri
	if dObj.Options.ImgFold {
		imgSrc = dObj.imgFoldNam + "/" + posId[4:] + ".jpeg"
	}

	//css
	cssStr := ""
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

	htmlStr += fmt.Sprintf("  <div id=\"%s\">\n",imgDivId)
	htmlStr += fmt.Sprintf("     <img src=\"%s\" alt=\"%s\" id=\"%s\">\n", imgSrc, imgProp.Title, imgId)
//	htmlStr += fmt.Sprintf("     <p id=\"%s\">%s</p>\n", pimgId, imgProp.Title)
	htmlStr += "  </div>\n"

	imgDispObj.bodyHtml = htmlStr
	imgDispObj.bodyCss = cssStr
	return &imgDispObj, nil
}


func (dObj *GdocDomObj) cvtTableToDom(tbl *docs.Table)(tabObj dispObj, err error) {
	// https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Traversing_an_HTML_table_with_JavaScript_and_DOM_Interfaces
	// table element
	var htmlStr, cssStr, scriptStr string
	var tabWidth float64
	var icol, trow int64
	var defcel tblCell


	doc := dObj.doc
	dObj.tableCount++
//	tblId := fmt.Sprintf("%s_tab_%d", dObj.docName, dObj.tableCount)

    docPg := doc.DocumentStyle
    PgWidth := docPg.PageSize.Width.Magnitude
    NetPgWidth := PgWidth - (docPg.MarginLeft.Magnitude + docPg.MarginRight.Magnitude)
//   fmt.Printf("Default Table Width: %.1f", NetPgWidth)
    tabWidth = NetPgWidth
	tabw := 0.0
    for icol=0; icol < tbl.Columns; icol++ {
        tcolObj :=tbl.TableStyle.TableColumnProperties[icol]
        if tcolObj.Width != nil {
            tabw += tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude
        }
    }
	if tabw > 0.0 {tabWidth = tabw}

// table cell default values
// define default cell classs
	tcelDef := tbl.TableRows[0].TableCells[0]
	tcelDefStyl := tcelDef.TableCellStyle

// default values which google does not set but uses
	defcel.vert_align = "top"
	defcel.bcolor = "black"
	defcel.bwidth = 1.0
	defcel.bdash = "solid"

// xxx
	if tcelDefStyl != nil {
		defcel.vert_align = util.Get_vert_align(tcelDefStyl.ContentAlignment)

// if left border is the only border specified, let's use it for default values
		tb := (tcelDefStyl.BorderTop == nil)&& (tcelDefStyl.BorderRight == nil)
		tb = tb&&(tcelDefStyl.BorderBottom == nil)
		if (tcelDefStyl.BorderLeft != nil) && tb {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}

		if tcelDefStyl.PaddingTop != nil {defcel.pad[0] = tcelDefStyl.PaddingTop.Magnitude}
		if tcelDefStyl.PaddingRight != nil {defcel.pad[1] = tcelDefStyl.PaddingRight.Magnitude}
		if tcelDefStyl.PaddingBottom != nil {defcel.pad[2] = tcelDefStyl.PaddingBottom.Magnitude}
		if tcelDefStyl.PaddingLeft != nil {defcel.pad[3] = tcelDefStyl.PaddingLeft.Magnitude}

		if tcelDefStyl.BackgroundColor != nil {defcel.bckcolor = util.GetColor(tcelDefStyl.BackgroundColor.Color)}

		if tcelDefStyl.BorderTop != nil {
			if tcelDefStyl.BorderTop.Color != nil {defcel.border[0].color = util.GetColor(tcelDefStyl.BorderTop.Color.Color)}
			if tcelDefStyl.BorderTop.Width != nil {defcel.border[0].width = tcelDefStyl.BorderTop.Width.Magnitude}
			defcel.border[0].dash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
		}
		if tcelDefStyl.BorderRight != nil {
			if tcelDefStyl.BorderRight.Color != nil {defcel.border[1].color = util.GetColor(tcelDefStyl.BorderRight.Color.Color)}
			if tcelDefStyl.BorderRight.Width != nil {defcel.border[1].width = tcelDefStyl.BorderRight.Width.Magnitude}
			defcel.border[1].dash = util.GetDash(tcelDefStyl.BorderRight.DashStyle)
		}
		if tcelDefStyl.BorderBottom != nil {
			if tcelDefStyl.BorderBottom.Color != nil {defcel.border[2].color = util.GetColor(tcelDefStyl.BorderBottom.Color.Color)}
			if tcelDefStyl.BorderBottom.Width != nil {defcel.border[2].width = tcelDefStyl.BorderBottom.Width.Magnitude}
			defcel.border[2].dash = util.GetDash(tcelDefStyl.BorderBottom.DashStyle)
		}
		if tcelDefStyl.BorderLeft != nil {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}
		if tcelDefStyl.BorderTop == tcelDefStyl.BorderRight {
//			fmt.Println("same border!")
			if tcelDefStyl.BorderTop != nil {
				if tcelDefStyl.BorderTop.Color != nil {defcel.bcolor = util.GetColor(tcelDefStyl.BorderTop.Color.Color)}
				defcel.bdash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
				if tcelDefStyl.BorderTop.Width != nil {defcel.bwidth = tcelDefStyl.BorderTop.Width.Magnitude}
			}
		}
	}

	//set up table
	tblClass := fmt.Sprintf("%s_tbl", dObj.docName)
	tblCellClass := fmt.Sprintf("%s_tcel", dObj.docName)
	//htmlStr = ""

	// if there is an open list, close it
	if len(*dObj.listStack) >= 0 {
		dObj.closeList(0)
	}

	// html fmt.Sprintf("<table class=\"%s\">\n", tblClass)

  	// table styling
  	cssStr = fmt.Sprintf(".%s {\n",tblClass)
 	cssStr += fmt.Sprintf("  border: 1px solid black;\n  border-collapse: collapse;\n")
 	cssStr += fmt.Sprintf("  width: %.1fpt;\n", tabWidth)
	cssStr += "   margin:auto;\n"
	cssStr += "}\n"

	// table columns
		// html htmlStr +="<colgroup>\n"

	for icol = 0; icol < tbl.Columns; icol++ {
		colId := fmt.Sprintf("tab%d_col%d", dObj.tableCount, icol)
		tabWtyp :=tbl.TableStyle.TableColumnProperties[icol].WidthType
		if tabWtyp == "FIXED_WIDTH" {
			cssStr += fmt.Sprintf("#%s {width: %.0fpt;}\n", colId, tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude)
			//html	htmlStr += fmt.Sprintf("<col span=\"1\" id=\"%s\">\n", colId)
		}
	}
	// html htmlStr +="</colgroup>\n"

	cssStr += fmt.Sprintf(".%s {\n",tblCellClass)
 	cssStr += fmt.Sprintf("  border: %.1fpt %s %s;\n", defcel.bwidth, defcel.bdash, defcel.bcolor)
	cssStr += fmt.Sprintf("  vertical-align: %s;\n", defcel.vert_align )
	cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", defcel.pad[0], defcel.pad[1], defcel.pad[2], defcel.pad[3])
 	cssStr += "}\n"


// row styling
	htmlStr += "  <tbody>\n"
	tblCellCount := 0
	for trow=0; trow < tbl.Rows; trow++ {
		htmlStr += fmt.Sprintf("  <tr>\n")
		trowobj := tbl.TableRows[trow]
//		mrheight := trowobj.TableRowStyle.MinRowHeight.Magnitude

		numCols := len(trowobj.TableCells)
		for tcol:=0; tcol< numCols; tcol++ {
			tcell := trowobj.TableCells[tcol]
			tblCellCount++
			cellStr := ""
			celId := fmt.Sprintf("tab%d_cell%d", dObj.tableCount, tblCellCount)
			// check whether cell style is different from default
			if tcell.TableCellStyle != nil {
				tstyl := tcell.TableCellStyle
				if tstyl.BackgroundColor != nil {cellStr += fmt.Sprintf(" background-color:\"%s\";",util.GetColor(tstyl.BackgroundColor.Color))}
				if util.Get_vert_align(tstyl.ContentAlignment) != defcel.vert_align {cellStr += fmt.Sprintf(" vertical-align: %s;", util.Get_vert_align(tstyl.ContentAlignment))}
				if tstyl.PaddingTop != nil {
					if tstyl.PaddingTop.Magnitude != defcel.pad[0] { cellStr += fmt.Sprintf(" padding-top: %5.1fpt;", tstyl.PaddingTop.Magnitude)}
				}

				if tstyl.BorderTop != nil {
					// Color
					if tstyl.BorderTop.Color != nil {cellStr += fmt.Sprintf(" border-top-color: %s;", util.GetColor(tstyl.BorderTop.Color.Color))}
					//dash
					if util.GetDash(tstyl.BorderTop.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-top-style: %s;",  util.GetDash(tstyl.BorderTop.DashStyle))}
					//Width
					if tstyl.BorderTop.Width != nil {cellStr += fmt.Sprintf(" border-top-width: %5.1fpt;", tstyl.BorderTop.Width.Magnitude)}
				}
			}
// xxxx
			if len(cellStr) >0 {
				cssStr += fmt.Sprintf("#%s {",celId)
				cssStr += fmt.Sprintf("%s }\n", cellStr)
				htmlStr += fmt.Sprintf("    <td id=\"%s\" class=\"%s\">\n", celId, tblCellClass)
			} else {
			// default
				htmlStr += fmt.Sprintf("    <td class=\"%s\">\n", tblCellClass)
			}
			elNum := len(tcell.Content)
			for el:=0; el< elNum; el++ {
				elObj := tcell.Content[el]
				tObj, err:=dObj.cvtContentElToDom(elObj)
				if err != nil {
					tabObj.script = scriptStr + fmt.Sprintf("\n// error cnvtContentEl: %v\n", err)
					tabObj.bodyCss = cssStr
					return tabObj, fmt.Errorf("cvtContentElToDom - ConvertTable: %v", err)
				}
				cssStr += tObj.bodyCss
				htmlStr += "    " + tObj.bodyHtml
			}
			htmlStr += "  </td>\n"

		}
		htmlStr += "</tr>\n"
	}

	htmlStr += "  </tbody>\n</table>\n"
	tabObj.bodyHtml = htmlStr
	tabObj.bodyCss = cssStr
	return tabObj, nil
}

func (dObj *GdocDomObj) cvtTableHtml(tbl *docs.Table)(tabObj dispObj, err error) {
	var htmlStr, cssStr string
	var tabWidth float64
	var icol, trow int64
	var defcel tblCell


	doc := dObj.doc
	dObj.tableCount++
//	parent := dObj.elDiv
//	tblId := fmt.Sprintf("%s_tab_%d", dObj.docName, dObj.tableCount)

    docPg := doc.DocumentStyle
    PgWidth := docPg.PageSize.Width.Magnitude
    NetPgWidth := PgWidth - (docPg.MarginLeft.Magnitude + docPg.MarginRight.Magnitude)
//   fmt.Printf("Default Table Width: %.1f", NetPgWidth)
    tabWidth = NetPgWidth
	tabw := 0.0
    for icol=0; icol < tbl.Columns; icol++ {
        tcolObj :=tbl.TableStyle.TableColumnProperties[icol]
        if tcolObj.Width != nil {
            tabw += tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude
        }
    }
	if tabw > 0.0 {tabWidth = tabw}

// table cell default values
// define default cell classs
	tcelDef := tbl.TableRows[0].TableCells[0]
	tcelDefStyl := tcelDef.TableCellStyle

// default values which google does not set but uses
	defcel.vert_align = "top"
	defcel.bcolor = "black"
	defcel.bwidth = 1.0
	defcel.bdash = "solid"

// xxx
	if tcelDefStyl != nil {
		defcel.vert_align = util.Get_vert_align(tcelDefStyl.ContentAlignment)

// if left border is the only border specified, let's use it for default values
		tb := (tcelDefStyl.BorderTop == nil)&& (tcelDefStyl.BorderRight == nil)
		tb = tb&&(tcelDefStyl.BorderBottom == nil)
		if (tcelDefStyl.BorderLeft != nil) && tb {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}

		if tcelDefStyl.PaddingTop != nil {defcel.pad[0] = tcelDefStyl.PaddingTop.Magnitude}
		if tcelDefStyl.PaddingRight != nil {defcel.pad[1] = tcelDefStyl.PaddingRight.Magnitude}
		if tcelDefStyl.PaddingBottom != nil {defcel.pad[2] = tcelDefStyl.PaddingBottom.Magnitude}
		if tcelDefStyl.PaddingLeft != nil {defcel.pad[3] = tcelDefStyl.PaddingLeft.Magnitude}

		if tcelDefStyl.BackgroundColor != nil {defcel.bckcolor = util.GetColor(tcelDefStyl.BackgroundColor.Color)}

		if tcelDefStyl.BorderTop != nil {
			if tcelDefStyl.BorderTop.Color != nil {defcel.border[0].color = util.GetColor(tcelDefStyl.BorderTop.Color.Color)}
			if tcelDefStyl.BorderTop.Width != nil {defcel.border[0].width = tcelDefStyl.BorderTop.Width.Magnitude}
			defcel.border[0].dash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
		}
		if tcelDefStyl.BorderRight != nil {
			if tcelDefStyl.BorderRight.Color != nil {defcel.border[1].color = util.GetColor(tcelDefStyl.BorderRight.Color.Color)}
			if tcelDefStyl.BorderRight.Width != nil {defcel.border[1].width = tcelDefStyl.BorderRight.Width.Magnitude}
			defcel.border[1].dash = util.GetDash(tcelDefStyl.BorderRight.DashStyle)
		}
		if tcelDefStyl.BorderBottom != nil {
			if tcelDefStyl.BorderBottom.Color != nil {defcel.border[2].color = util.GetColor(tcelDefStyl.BorderBottom.Color.Color)}
			if tcelDefStyl.BorderBottom.Width != nil {defcel.border[2].width = tcelDefStyl.BorderBottom.Width.Magnitude}
			defcel.border[2].dash = util.GetDash(tcelDefStyl.BorderBottom.DashStyle)
		}
		if tcelDefStyl.BorderLeft != nil {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = util.GetColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = util.GetDash(tcelDefStyl.BorderLeft.DashStyle)
		}
		if tcelDefStyl.BorderTop == tcelDefStyl.BorderRight {
//			fmt.Println("same border!")
			if tcelDefStyl.BorderTop != nil {
				if tcelDefStyl.BorderTop.Color != nil {defcel.bcolor = util.GetColor(tcelDefStyl.BorderTop.Color.Color)}
				defcel.bdash = util.GetDash(tcelDefStyl.BorderTop.DashStyle)
				if tcelDefStyl.BorderTop.Width != nil {defcel.bwidth = tcelDefStyl.BorderTop.Width.Magnitude}
			}
		}
	}

	//set up table
	tblClass := fmt.Sprintf("%s_tbl", dObj.docName)
	tblCellClass := fmt.Sprintf("%s_tcel", dObj.docName)
	htmlStr = ""

	// if there is an open list, close it
	if len(*dObj.listStack) >= 0 {
		dObj.closeList(0)
//fmt.Printf("table closing list!\n")
	}

	htmlStr += fmt.Sprintf("<table class=\"%s\">\n", tblClass)

  	// table styling
  	cssStr = fmt.Sprintf(".%s {\n",tblClass)
 	cssStr += fmt.Sprintf("  border: 1px solid black;\n  border-collapse: collapse;\n")
 	cssStr += fmt.Sprintf("  width: %.1fpt;\n", tabWidth)
	cssStr += "   margin:auto;\n"
	cssStr += "}\n"

	// table columns
	tabWtyp :=tbl.TableStyle.TableColumnProperties[0].WidthType
//fmt.Printf("table width type: %s\n", tabWtyp)
	if tabWtyp == "FIXED_WIDTH" {
		htmlStr +="<colgroup>\n"
		for icol = 0; icol < tbl.Columns; icol++ {
			colId := fmt.Sprintf("tab%d_col%d", dObj.tableCount, icol)
			cssStr += fmt.Sprintf("#%s {width: %.1fpt;}\n", colId, tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude)
			htmlStr += fmt.Sprintf("<col span=\"1\" id=\"%s\">\n", colId)
		}
		htmlStr +="</colgroup>\n"
	}


	cssStr += fmt.Sprintf(".%s {\n", tblCellClass)
 	cssStr += fmt.Sprintf("  border: %.1fpt %s %s;\n", defcel.bwidth, defcel.bdash, defcel.bcolor)
	cssStr += fmt.Sprintf("  vertical-align: %s;\n", defcel.vert_align )
	cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", defcel.pad[0], defcel.pad[1], defcel.pad[2], defcel.pad[3])
 	cssStr += "}\n"


// row styling
	htmlStr += "  <tbody>\n"
	tblCellCount := 0
	for trow=0; trow < tbl.Rows; trow++ {
		htmlStr += fmt.Sprintf("  <tr>\n")
		trowobj := tbl.TableRows[trow]
//		mrheight := trowobj.TableRowStyle.MinRowHeight.Magnitude

		numCols := len(trowobj.TableCells)
		for tcol:=0; tcol< numCols; tcol++ {
			tcell := trowobj.TableCells[tcol]
			tblCellCount++
			cellStr := ""
			celId := fmt.Sprintf("tab%d_cell%d", dObj.tableCount, tblCellCount)
			// check whether cell style is different from default
			if tcell.TableCellStyle != nil {
				tstyl := tcell.TableCellStyle
				if tstyl.BackgroundColor != nil {cellStr += fmt.Sprintf(" background-color:\"%s\";",util.GetColor(tstyl.BackgroundColor.Color))}
				if util.Get_vert_align(tstyl.ContentAlignment) != defcel.vert_align {cellStr += fmt.Sprintf(" vertical-align: %s;", util.Get_vert_align(tstyl.ContentAlignment))}
				if tstyl.PaddingTop != nil {
					if tstyl.PaddingTop.Magnitude != defcel.pad[0] { cellStr += fmt.Sprintf(" padding-top: %5.1fpt;", tstyl.PaddingTop.Magnitude)}
				}

				if tstyl.BorderTop != nil {
					// Color
					if tstyl.BorderTop.Color != nil {cellStr += fmt.Sprintf(" border-top-color: %s;", util.GetColor(tstyl.BorderTop.Color.Color))}
					//dash
					if util.GetDash(tstyl.BorderTop.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-top-style: %s;",  util.GetDash(tstyl.BorderTop.DashStyle))}
					//Width
					if tstyl.BorderTop.Width != nil {cellStr += fmt.Sprintf(" border-top-width: %5.1fpt;", tstyl.BorderTop.Width.Magnitude)}
				}
			}
			if len(cellStr) >0 {
				cssStr += fmt.Sprintf("#%s {",celId)
				cssStr += fmt.Sprintf("%s }\n", cellStr)
				htmlStr += fmt.Sprintf("    <td id=\"%s\" class=\"%s\">\n", celId, tblCellClass)
			} else {
			// default
				htmlStr += fmt.Sprintf("    <td class=\"%s\">\n", tblCellClass)
			}
			elNum := len(tcell.Content)
			for el:=0; el< elNum; el++ {
				elObj := tcell.Content[el]
				tObj, err:=dObj.cvtContentElToDom(elObj)
				if err != nil {
// xxxx
					tabObj.bodyHtml = htmlStr
					tabObj.bodyCss = cssStr
					return tabObj, fmt.Errorf("ConvertTable: %v", err)
				}
				cssStr += tObj.bodyCss
				htmlStr += "    " + tObj.bodyHtml
			}
			htmlStr += "  </td>\n"

		}
		htmlStr += "</tr>\n"
	}

	htmlStr += "  </tbody>\n</table>\n"
	tabObj.bodyHtml = htmlStr
	tabObj.bodyCss = cssStr
	return tabObj, nil
}

func (dObj *GdocDomObj) cvtParToDom(par *docs.Paragraph)(parObj dispObj, err error) {
// paragraph element par
// - Bullet
// - Elements
// - ParagraphStyle
// - Positioned Objects
//
	var listCss, scriptStr string
	var newList cList
	var	listEl elScriptObj
	var orList, unList elScriptObj
	var parent string

	if par == nil {
        return parObj, fmt.Errorf("cvtPar -- parEl is nil!")
    }

//	parent := dObj.elDiv

	dObj.parCount++

	// first we need to check whether this is a cr-only paragraph
	if len(par.Elements) == 1 {
       if par.Elements[0].TextRun != nil {
            if par.Elements[0].TextRun.Content == "\n" {
				brEl := elScriptObj{typ: "br", newEl: "noel",}
				brEl.parent = dObj.parent
                parObj.script = addElToDom(brEl)
                return parObj, nil
            }
        }

	}

	isList := false
	if par.Bullet != nil {isList = true}
//fmt.Printf("********** par %d list: %t ***********\n", dObj.parCount, isList)

	if par.Bullet == nil {
		// if we have a non-list paragraph. we assume any open lists need to be closed
		// in a Dom we wander back to div_main
		if dObj.listStack != nil {dObj.closeList(0)}
			//fmt.Printf("new par -> close list\n")
	}


	// Positioned Objects
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
//		return parObj, nil
	}

	// par elements includes text


	// get paragraph style
// need to fix we know there is no list
	if (par.ParagraphStyle != nil) && (par.Bullet == nil) {
		parStyl, _, err := dObj.cvtParStylToDom(par.ParagraphStyle, dObj.parent, isList)
		if err != nil {
			parStyl.bodyCss += fmt.Sprintf("/* error cvtParStyl: %v */\n", err)
		}
		addDispObj(&parObj,&parStyl)

	// Heading Id refers to a heading paragraph not just a normal paragraph
	// headings are bookmarked for TOC

	// par elements: text and css for text

		parElSumDisp, err := dObj.cvtParElToDom(par)
		if err != nil {parElSumDisp.script += fmt.Sprintf("// error cvtParElDom: %v\n",err)}
		addDispObj(&parObj, &parElSumDisp)
		return parObj, nil
	}

	// lists
    if par.Bullet != nil {
		// there is paragraph style for each ul and a text style for each list element
// still todo
// need to apply bulletTxtMap to marker


		if dObj.Options.Verb {
			// htnm listHtml += fmt.Sprintf("<!-- List Element %d -->\n", dObj.parCount)
			// script
			scriptStr += fmt.Sprintf("// List El %d\n", dObj.parCount)
		}
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
//		fmt.Printf("listid: %s \n", listid)
//		printLiStack(dObj.listStack)

		listAtt, cNest := getLiStack(dObj.listStack)
		//printLiStackItem(listAtt, cNest)
		scriptStr += fmt.Sprintf("// " + "listid: %s listAtt: %s\n", listid, listAtt.cListId)

		switch listid == listAtt.cListId {
			case true:
				switch {
					case nestIdx > cNest:
						for nl:=cNest + 1; nl <= nestIdx; nl++ {
							newList.cListId = listid
							newList.cOrd = listOrd
							newStack := pushLiStack(dObj.listStack, newList)
							dObj.listStack = newStack

							if listOrd {
								// html listHtml = fmt.Sprintf("<ol class=\"%s_ol nL_%d\">\n", listid[4:], nl)
								// script
								orList.parent = fmt.Sprintf("ol_%d", nl - 1)
								orList.newEl = fmt.Sprintf("ol_%d", nl)
								orList.typ = "ol"
								orList.cl1 = listid[4:] + "_ol"
								orList.cl2 = fmt.Sprintf("nL_%d", nl)
								scriptStr += addElToDom(orList)
								// css
								listCss = fmt.Sprintf(".%s_ol.nL_%d {\n", listid[4:], nl)
								listCss += fmt.Sprintf("  counter-reset: %s_nL_%d\n",listid[4:], nl)
								listCss += "}\n"
							} else {
								// html listHtml = fmt.Sprintf("<ul class=\"%s_ul nL_%d\">\n", listid[4:], nl)
								// script
								unList.parent = fmt.Sprintf("ul_%d", nl - 1)
								unList.newEl = fmt.Sprintf("ul_%d", nl)
								unList.typ = "ul"
								unList.cl1 = listid[4:] + "_ul"
								unList.cl2 = fmt.Sprintf("nL_%d", nl)
								scriptStr += addElToDom(orList)
								// css none
							}
						}
						// html	listHtml += fmt.Sprintf("<!-- same list increase %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)
						//html fmt.Printf("<!-- same list increase %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)

					case nestIdx < cNest:
						// html
						// listHtml = dObj.closeList(nestIdx)
						// listHtml += fmt.Sprintf("<!-- same list reduce %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)
//						scriptStr += fmt.Sprintf("// list id %s new NL: %d old NL: %d\n", list id, nestIdx, cNest)
						//html diag	fmt.Printf("<!-- same list reduce %s new NL %d  old Nl %d -->\n", listid, nestIdx, cNest)

						// script
						dObj.closeList(nestIdx)


					case nestIdx == cNest:
//						listHtml =""
				}

			case false:
				// new list
				// close list first
				// html listHtml = dObj.closeList(-1)
				// html listHtml += fmt.Sprintf("<!-- new list %s %s -->\n", listid, listAtt.cListId)
//fmt.Printf("<!-- new list %s %s -->\n", listid, listAtt.cListId)

				// start a new list
				scriptStr += "// " + "new list\n"
				newList.cListId = listid
				newList.cOrd = listOrd
				newStack := pushLiStack(dObj.listStack, newList)
				dObj.listStack = newStack
				nl := nestIdx
				parent = ""
				if nl == 0 {
					parent = dObj.parent
				}
		fmt.Printf("parent: %s nl: %d \n", parent, nl)
				if listOrd {
					// html listHtml += fmt.Sprintf("<ol class=\"%s_ol nL_%d\">\n", listid[4:], nestIdx)
					// script
					if len(parent) == 0 {parent = fmt.Sprintf("ol_%d", nl - 1)}

					orList.parent = parent
					orList.newEl = fmt.Sprintf("ol_%d", nl)
					orList.typ = "ol"
					orList.cl1 = listid[4:] + "_ol"
					orList.cl2 = fmt.Sprintf("nL_%d", nl)
					scriptStr += addElToDom(orList)
					// css
					listCss = fmt.Sprintf(".%s_ol.nL_%d {\n", listid[4:], nestIdx)
					listCss += fmt.Sprintf("  counter-reset: %s_nL_%d\n",listid[4:], nestIdx)
					listCss += "}\n"
				} else {
					// html listHtml += fmt.Sprintf("<ul class=\"%s_ul nL_%d\">\n", listid[4:], nestIdx)
					if len(parent) == 0 {parent = fmt.Sprintf("ul_%d", nl - 1)}

		fmt.Printf("ul parent: %s nl: %d \n", parent, nl)
					unList.parent = parent
					unList.newEl = fmt.Sprintf("ul_%d", nl)
					unList.typ = "ul"
					unList.cl1 = listid[4:] + "_ul"
					unList.cl2 = fmt.Sprintf("nL_%d", nl)
					scriptStr += addElToDom(unList)
					//css
				}
		} // switch

		parObj.script += scriptStr

		// html <li>
		// html listPrefix = fmt.Sprintf("<li class=\"%s_li nL_%d\">", listid[4:], nestIdx)
		// script
		listParent := ""
		nl := nestIdx
		if listOrd {
			listParent = fmt.Sprintf("ol_%d", nl)
		} else {
			listParent = fmt.Sprintf("ul_%d", nl)
		}

		listEl.parent = listParent
		listEl.cl1 = listid[4:] + "_li"
		listEl.cl2 = fmt.Sprintf("nL_%d", nestIdx)
		listEl.typ = "li"
		listEl.newEl = "lsIt"
		parObj.script += addElToDom(listEl)


		// mark
		if par.Bullet.TextStyle != nil {
//      	    bulletTxtMap = fillTxtMap(par.Bullet.TextStyle)
		}

		// get paragraph style
		parent = "lsIt"
		parStyl, _, err := dObj.cvtParStylToDom(par.ParagraphStyle, parent, isList)
		if err != nil {
			parStyl.bodyCss += fmt.Sprintf("/* error cvtParStyl: %v */\n", err)
		}
		addDispObj(&parObj,&parStyl)

		// Heading Id refers to a heading paragraph not just a normal paragraph
		// headings are bookmarked for TOC

		// par elements: text and css for text

		parElSumDisp, err := dObj.cvtParElToDom(par)
		if err != nil {parElSumDisp.script += fmt.Sprintf("// error cvtParElDom: %v\n",err)}
		addDispObj(&parObj, &parElSumDisp)

	}

	parObj.bodyCss += listCss
	return parObj, nil
}

func (dObj *GdocDomObj) cvtParElToDom(par *docs.Paragraph)(parDisp dispObj, err error) {

//	hasList := false
//	if par.Bullet != nil {hasList = true}

	namedTyp := par.ParagraphStyle.NamedStyleType
    numParEl := len(par.Elements)
    for pEl:=0; pEl< numParEl; pEl++ {
        parEl := par.Elements[pEl]

		if parEl.InlineObjectElement != nil {
    	    imgObj, err := dObj.renderInlineImg(parEl.InlineObjectElement)
        	if err != nil {
				imgObj.bodyHtml += fmt.Sprintf("<!-- error cvtPelToHtml: %v -->\n",err)
        	}
			addDispObj(&parDisp, imgObj)
		}

		if parEl.TextRun != nil {
			txtObj := dObj.cvtParDomElText(parEl.TextRun, namedTyp)
//        	if err != nil {
//            	txtObj.bodyHtml += fmt.Sprintf("<!-- error cvtPelToHtml: %v -->\n",err)
//        	}
			addDispObj(&parDisp, &txtObj)
		}

		if parEl.FootnoteReference != nil {
			dObj.ftnoteCount++
//        	htmlStr += fmt.Sprintf("<span class=\"%s_ftnote\">[%d]</span>",dObj.docName, dObj.ftnoteCount)
//			parDisp.bodyHtml += htmlStr
		}

		if parEl.PageBreak != nil {

		}

		if parEl.HorizontalRule != nil {
			horDisp := dObj.cvtHrElToDom(parEl.HorizontalRule)
			addDispObj(&parDisp, &horDisp)
		}

		if parEl.ColumnBreak != nil {

		}

		if parEl.Person != nil {

		}

		if parEl.RichLink != nil {

		}

	} //loop end parEl

	return parDisp, nil
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

	defTxtMap := fillTxtMap(normalTxtStyl)
    defParMap := fillParMap(normalParStyl)


	for namedTyp, res := range dObj.namStylMap {
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
			parCss, _ := cvtParMapStylCss(defParMap, namParStyl, dObj.Options)
			txtCss := cvtTxtMapStylCss(defTxtMap, namTxtStyl)
			cssStr += cssPrefix + parCss + txtCss + "}\n"
		}
	}
	return cssStr, nil
}

func (dObj *GdocDomObj) cvtParStylToDom(parStyl *docs.ParagraphStyle, parent string, isList bool)(parStylObj dispObj, alter bool, err error) {
	var elObj elScriptObj
	var cssStr string
	var namParStyl *docs.ParagraphStyle
	// changed from Html need to handle case if parStyl == nil
	// q: is there a case where parstyl == nil
	// if parstyl == nil lets assume normal_text

	if parStyl == nil {
		parStyl,_,_ = dObj.getNamedStyl("NORMAL_TEXT")
		namParStyl = parStyl
	} else {
		namedTyp := parStyl.NamedStyleType
		namParStyl, _, err = dObj.getNamedStyl(namedTyp)
		if err != nil {
			return parStylObj, false, fmt.Errorf("getNamedStyl: %v", err)
		}
	}

	// default style for each named style used in the document
	// add css for named style at the begining of the style sheet
	// normal_text is already defined as the default in the css for the <div>

	cssComment := ""
	if namParStyl == nil {
		// def error the default is that the normal_text paragraph style is passed
		cssComment = fmt.Sprintf("/* Paragraph Style: no named Style */\n")
		parStylObj.bodyCss = cssComment
		return parStylObj, false, nil
	}

	cssComment = fmt.Sprintf("/* Paragraph Style: %s */\n", parStyl.NamedStyleType )

	alter = false
	cssParAtt := ""

	parmap := fillParMap(namParStyl)

	if parStyl == nil || isList {
		// use named style that has been published
		cssParAtt = cvtParMapCss(parmap, dObj.Options)
	} else {
		cssParAtt, alter = cvtParMapStylCss(parmap, parStyl, dObj.Options)
	}

    // NamedStyle Type
	cssPrefix := ""
	headingId := parStyl.HeadingId

	switch parStyl.NamedStyleType {
		case "TITLE":
			elObj.typ = "p"
			if dObj.namStylMap["TITLE"] && !alter {
				//html prefix = fmt.Sprintf("<p class=\"%s_title%s\"", dObj.docName, isListClass)
				elObj.cl1 = fmt.Sprintf("%s_title", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_title.%s_title_%d {\n", dObj.docName, dObj.docName, dObj.title.count)
				//html prefix = fmt.Sprintf("<p class=\"%s_title %s_title_%d\"", dObj.docName, dObj.docName, dObj.title.count)
				elObj.cl1 = fmt.Sprintf("%s_title", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_title_%d", dObj.docName, dObj.title.count)
				dObj.title.count++
			}
			//html suffix = "</p>"

		case "SUBTITLE":
			if dObj.namStylMap["SUBTITLE"] && !alter {
				//html prefix = fmt.Sprintf("<p class=\"%s_subtitle\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_subtitle", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".s_subtitle.s_subtitle_%d {\n", dObj.docName, dObj.docName, dObj.subtitle.count)
				//html prefix = fmt.Sprintf("<p class=\"%s_subtitle %s_subtitle_%d\"", dObj.docName, dObj.docName, dObj.subtitle.count)
				elObj.cl1 = fmt.Sprintf("%s_subtitle", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_subtitle_%d", dObj.docName, dObj.subtitle.count)
				dObj.subtitle.count++
			}

		case "HEADING_1":
			elObj.typ = "h1"
			if dObj.namStylMap["HEADING_1"] && !alter {
				//html prefix = fmt.Sprintf("<h1 class=\"%s_h1\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h1", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_h1.%s_h1_%d {\n",dObj.docName, dObj.docName, dObj.h1.count)
				//html prefix = fmt.Sprintf("<h1 class=\"%s_h1 %s_h1_%d\"", dObj.docName, dObj.docName, dObj.h1.count)
				elObj.cl1 = fmt.Sprintf("%s_h1", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h1_%d", dObj.docName, dObj.h1.count)
				dObj.h1.count++
			}

		case "HEADING_2":
			elObj.typ = "h2"
			if dObj.namStylMap["HEADING_2"] && !alter {
				//html prefix = fmt.Sprintf("<h2 class=\"%s_h2\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h2", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_h2.%s_h2_%d {\n",dObj.docName, dObj.docName, dObj.h2.count)
				//html prefix = fmt.Sprintf("<h2 class=\"%s_h2 %s_h2_%d\"", dObj.docName, dObj.docName, dObj.h2.count)
				elObj.cl1 = fmt.Sprintf("%s_h2", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h2_%d", dObj.docName, dObj.h2.count)
				dObj.h2.count++
			}
			//html suffix = "</h2>"

		case "HEADING_3":
			elObj.typ = "h3"
			if dObj.namStylMap["HEADING_3"] && !alter {
				//html prefix = fmt.Sprintf("<h3 class=\"%s_h3\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h3", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_h3.%s_h3_%d {\n",dObj.docName, dObj.docName, dObj.h3.count)
				//html prefix = fmt.Sprintf("<h3 class=\"%s_h3 %s_h3_%d\"", dObj.docName, dObj.docName, dObj.h3.count)
				elObj.cl1 = fmt.Sprintf("%s_h3", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h3_%d", dObj.docName, dObj.h3.count)
				dObj.h3.count++
			}

		case "HEADING_4":
			elObj.typ = "h4"
			if dObj.namStylMap["HEADING_4"] && !alter {
				//html prefix = fmt.Sprintf("<h4 class=\"%s_h4\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h4", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_h4.%s_h4_%d {\n",dObj.docName, dObj.docName, dObj.h4.count)
				//html prefix = fmt.Sprintf("<h4 class=\"%s_h4 %s_h4_%d\"", dObj.docName, dObj.docName, dObj.h4.count)
				elObj.cl1 = fmt.Sprintf("%s_h4", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h4_%d", dObj.docName, dObj.h4.count)
				dObj.h4.count++
			}

		case "HEADING_5":
			elObj.typ = "h5"
			if dObj.namStylMap["HEADING_5"] && !alter {
				//html prefix = fmt.Sprintf("<h5 class=\"%s_h5\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h5", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf(".%s_h5.%s_h5_%d {\n",dObj.docName, dObj.docName, dObj.h5.count)
				//html prefix = fmt.Sprintf("<h5 class=\"%s_h5 %s_h5_%d\"", dObj.docName, dObj.docName, dObj.h5.count)
				elObj.cl1 = fmt.Sprintf("%s_h5", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h5_%d", dObj.docName, dObj.h5.count)
				dObj.h5.count++
			}

		case "HEADING_6":
			elObj.typ = "h6"
			if dObj.namStylMap["HEADING_6"] && !alter {
				//html prefix = fmt.Sprintf("<h6 class=\"%s_h6\"", dObj.docName)
				elObj.cl1 = fmt.Sprintf("%s_h6", dObj.docName)
			}
			if alter {
				cssPrefix = fmt.Sprintf("%s_h6.%s_h6_%d {\n",dObj.docName, dObj.docName, dObj.h6.count)
				//html prefix = fmt.Sprintf("<h6 class=\"%s_h6 %s_h6_%d\"", dObj.docName, dObj.docName, dObj.h6.count)
				elObj.cl1 = fmt.Sprintf("%s_h6", dObj.docName)
				elObj.cl2 = fmt.Sprintf("%s_h6_%d", dObj.docName, dObj.h6.count)
				dObj.h6.count++
			}

		case "NORMAL_TEXT":
			switch {
//				case isList:
					// html prefix = "<span>"
//					elObj.typ = "p"
				case alter:
					cssPrefix = fmt.Sprintf(".%s_p.%s_p_%d {\n",dObj.docName, dObj.docName, dObj.parCount)
					//html prefix = fmt.Sprintf("<p class=\"%s_p %s_p_%d\"",dObj.docName, dObj.docName, dObj.parCount)
					elObj.typ = "p"
					elObj.cl1 = fmt.Sprintf("%s_p", dObj.docName)
					elObj.cl2 = fmt.Sprintf("%s_p_%d", dObj.docName, dObj.parCount)

				default:
				//html prefix = fmt.Sprintf("<p class=\"%s_p\"", dObj.docName)
					elObj.typ = "p"
					elObj.cl1 = fmt.Sprintf("%s_p", dObj.docName)
			}
		case "NAMED_STYLE_TYPE_UNSPECIFIED":
//			namTypValid = false
			elObj.typ = "p"
			elObj.cl1 = fmt.Sprintf("%s_h1", dObj.docName)
		default:
//			namTypValid = false
	}

	if len(headingId) > 0 {
//		prefix = fmt.Sprintf("%s id=\"%s\">", prefix, headingId[3:])
		elObj.idStr = headingId[3:]
	} else {
//		prefix = prefix + ">"
	}

	//fmt.Printf("parstyl: %s %s %s\n", parStyl.NamedStyleType, prefix, suffix)
	if (len(cssPrefix) > 0) {cssStr = cssComment + cssPrefix + cssParAtt + "}\n"}
	elObj.comment ="cvtParStyl"
	elObj.parent = parent
	elObj.newEl = "hdel"
	parStylObj.script = addElToDom(elObj)
	parStylObj.bodyCss = cssStr
	return parStylObj, alter, nil
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

func (dObj *GdocDomObj) createSectionDiv() (secHd *dispObj) {
	var secHead dispObj

	if !dObj.Options.Sections {return nil}

	if len(dObj.sections) < 2 {return nil}

	htmlStr := fmt.Sprintf("<div class=\"%s_main top\" id=\"%s_sectoc\">\n", dObj.docName, dObj.docName)
	htmlStr += fmt.Sprintf("<p class=\"%s_title %s_leftTitle_UL\">Sections</p>\n",dObj.docName, dObj.docName)
	for i:=0; i< len(dObj.sections); i++ {
		htmlStr += fmt.Sprintf("  <p class=\"%s_p\"><a href=\"#%s_sec_%d\">Page: %3d</a></p>\n", dObj.docName, dObj.docName, i, i)
	}
	htmlStr +="</div>\n"
	secHead.bodyHtml = htmlStr
	return &secHead
}

func (dObj *GdocDomObj) createSectionHeading(ipage int) (secObj dispObj) {
// method that creates a distinct html dvision per section with a page heading

	secObj.bodyCss = fmt.Sprintf(".%s_main.sec_%d {\n", dObj.docName, ipage)

	// html
	secObj.bodyHtml = fmt.Sprintf("<div class=\"%s_main sec_%d\" id=\"%s_sec_%d\">\n", dObj.docName, ipage, dObj.docName, ipage)
	secObj.bodyHtml += fmt.Sprintf("<p class=\"%s_page\"><a href=\"#%s_sectoc\">Page %d</a></p>\n", dObj.docName, dObj.docName, ipage)

	return secObj
}

func (dObj *GdocDomObj) creCssDocHead() (headCss string, err error) {

	var cssStr, errStr string

    docStyl := dObj.doc.DocumentStyle
    dObj.docWidth = (docStyl.PageSize.Width.Magnitude - docStyl.MarginRight.Magnitude - docStyl.MarginLeft.Magnitude)

    //gdoc default el css and doc css
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

	//css default text style
	cssStr = fmt.Sprintf(".%s_main {\n", dObj.docName)
	parStyl, txtStyl, err := dObj.getNamedStyl("NORMAL_TEXT")
	if err != nil {
		return headCss, fmt.Errorf("creHeadCss: %v", err)
	}

	defParMap := fillParMap(parStyl)
	defTxtMap := fillTxtMap(txtStyl)

	cssStr += "  display:block;\n"
	cssStr += "  margin: 0;\n"
	if dObj.Options.DivBorders {
		cssStr += "  border: solid green;\n"
		cssStr += "  border-width: 1px;\n"
	}
	cssStr += cvtTxtMapCss(defTxtMap)
	cssStr += "}\n"
	headCss += cssStr

	hdcss, err := dObj.cvtDocNamedStyles()
	if err != nil {
		errStr = fmt.Sprintf("cvtDocNamedStyles %v", err)
	}
	headCss += hdcss + errStr

	// paragraph default style
    pCssStr := cvtParMapCss(defParMap, dObj.Options)
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
		listProp := dObj.doc.Lists[listid].ListProperties

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
			nestLev := listProp.NestingLevels[nl]

			glyphStr := util.GetGlyphStr(nestLev)
			switch dObj.docLists[i].ord {
				case true:
					cssStr += fmt.Sprintf(".%s_ol.nL_%d {\n", listClass, nl)
				case false:
					cssStr += fmt.Sprintf(".%s_ul.nL_%d {\n", listClass, nl)
			}

			idFl := nestLev.IndentFirstLine.Magnitude - cumIndent
			idSt := nestLev.IndentStart.Magnitude - cumIndent
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

            cssStr += cvtTxtMapStylCss(defTxtMap,nestLev.TextStyle)
			cssStr += fmt.Sprintf("}\n")
		}
	}
	headCss += cssStr

   // css default table
    if dObj.tableCount > 0 {
        //css default table styling
        cssStr = fmt.Sprintf(".%s_tbl {\n", dObj.docName)
        //def cell
        cssStr += fmt.Sprintf("  border: 1px solid black;\n  border-collapse: collapse;\n")
        cssStr += fmt.Sprintf("  width: %.1fpt;\n", dObj.docWidth)
        cssStr += "   margin:auto;\n"
        cssStr += "}\n"
    }

    headCss += cssStr
	return headCss, nil
}

func (dObj *GdocDomObj) cvtContentElToDom(contEl *docs.StructuralElement) (GdocDomObj *dispObj, err error) {
	if dObj == nil {
		return nil, fmt.Errorf("error -- dObj is nil")
	}
//	parent = dObj.eldiv

	bodyElObj := new(dispObj)

	if contEl.Paragraph != nil {
		parEl := contEl.Paragraph
		tObj, err := dObj.cvtParToDom(parEl)
		if err != nil { bodyElObj.bodyHtml += fmt.Sprintf("<!-- %v -->\n", err) }
		addDispObj(bodyElObj, &tObj)
	}

	if contEl.SectionBreak != nil {

	}
	if contEl.Table != nil {
		tableEl := contEl.Table
		tObj, err := dObj.cvtTableToDom(tableEl)
		if err != nil { bodyElObj.bodyHtml += fmt.Sprintf("<!-- %v -->\n", err) }
		addDispObj(bodyElObj, &tObj)
	}
	if contEl.TableOfContents != nil {

	}
//	fmt.Println(" ConvertEl: ",htmlObj)
	return bodyElObj, nil
}

//ootnote div
func (dObj *GdocDomObj) createFootnoteDiv () (ftnoteDiv *dispObj, err error) {
	var ftnDiv dispObj
	var htmlStr, cssStr, scriptStr string

	doc := dObj.doc
	if len(dObj.docFtnotes) == 0 {
		return nil, nil
	}

	if len(dObj.docFtnotes) == 0 {return nil, nil}

	//html div footnote
	htmlStr = fmt.Sprintf("<!-- Footnotes: %d -->\n", len(dObj.docFtnotes))
	htmlStr += fmt.Sprintf("<div class=\"%s_main %s_ftndiv\">\n", dObj.docName, dObj.docName)

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
	htmlStr += fmt.Sprintf("<p class=\"%s_main %s_title %s_ftTit\">Footnotes</p>\n", dObj.docName, dObj.docName, dObj.docName)
//	ftnDiv.bodyHtml = htmlStr

	//css footnote title
	cssStr += fmt.Sprintf(".%s_main.%s_title.%s_ftTit {\n", dObj.docName, dObj.docName, dObj.docName)
	cssStr += "  color: purple;\n"
	cssStr += "}\n"

	// list for footnotes
	htmlStr +=fmt.Sprintf("<ol class=\"%s_ftnOl\">\n", dObj.docName)
	cssStr += fmt.Sprintf(".%s_ftnOl {\n", dObj.docName)
	cssStr += "  display:block;\n"
	cssStr += "  list-style-type: decimal;\n"
	cssStr += "  padding-inline-start: 10pt;\n"
	cssStr += "  margin: 0;\n"
	cssStr += "}\n"

	// prefix for paragraphs
	cssStr += fmt.Sprintf(".%s_p.%s_pft {\n",dObj.docName, dObj.docName)
	cssStr += "text-indent: 10pt;"
	cssStr += "counter-increment:ftcounter;"
	cssStr += "}\n"
	cssStr += fmt.Sprintf(".%s_p.%s_pft::before {\n",dObj.docName, dObj.docName)
	cssStr += "counter(ftcounter) ' ';"
	cssStr += "}\n"
	ftnDiv.bodyCss = cssStr
	ftnDiv.bodyHtml = htmlStr

	// footnotes paragraph html
	htmlStr = ""
	cssStr = ""
	for iFtn:=0; iFtn<len(dObj.docFtnotes); iFtn++ {
		idStr := dObj.docFtnotes[iFtn].id
//		ftnDiv.bodyHtml += htmlStr
		// reset htmlStr
		docFt, ok := doc.Footnotes[idStr]
		if !ok {
			htmlStr += fmt.Sprintf("<!-- error ftnote %d not found! -->\n", iFtn)
			continue
		}
		htmlStr = fmt.Sprintf("<!-- FTnote: %d %s els: %d -->\n", iFtn, idStr, len(docFt.Content))
		htmlStr +="<li>\n"
		ftnDiv.bodyHtml += htmlStr
		// presumably footnotes are paragraphs only
		for el:=0; el<len(docFt.Content); el++ {
			htmlStr = ""
			cssStr = ""
			elObj := docFt.Content[el]
			if elObj.Paragraph == nil {continue}
			par := elObj.Paragraph
//			pidStr := idStr[5:]
			//html htmlStr += fmt.Sprintf("<p class=\"%s_p %s_pft\" id=\"%s\">\n", dObj.docName, dObj.docName, pidStr)

			var parElSumDisp *dispObj
// need to change
			tDisp, err := dObj.cvtParElToDom(par)
			addDispObj(parElSumDisp, &tDisp)

			tObj, err := dObj.cvtContentElToDom(elObj)
			if err != nil {
// xxxx
				ftnDiv.bodyHtml += fmt.Sprintf("<!-- error display el: %d -->\n", el)
			}
			addDispObj(&ftnDiv, tObj)

//			ftnDiv.bodyHtml += htmlStr
			ftnDiv.bodyCss += cssStr
			ftnDiv.script += scriptStr
		}

	}

	return &ftnDiv, nil
}

//toc div
func (dObj *GdocDomObj) createTocDiv () (tocObj *dispObj, err error) {
	var tocDiv dispObj
	var htmlStr, cssStr string

	if dObj.Options.Toc != true { return nil, nil }

	if dObj.Options.Verb {
		if len(dObj.headings) < 2 {
//			fmt.Printf("*** no TOC insufficient headings: %d ***\n", len(dObj.headings))
			return nil, nil
		}
		fmt.Printf("*** creating TOC Div ***\n")
	}

	if len(dObj.headings) < 2 {
		tocDiv.bodyHtml = fmt.Sprintf("<!-- no toc insufficient headings -->")
		return nil, nil
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
	htmlStr = fmt.Sprintf("<div class=\"%s_main top\">\n", dObj.docName)
	htmlStr += fmt.Sprintf("<p class=\"%s_title %s_leftTitle\">Table of Contents</p>\n", dObj.docName, dObj.docName)
	tocDiv.bodyHtml = htmlStr

	//html all headings are entries to toc table of content
	for ihead:=0; ihead<len(dObj.headings); ihead++ {
		cssStr = ""
		htmlStr = ""
//		elStart := dObj.headings[ihead].hdElStart
		namedStyl := dObj.headings[ihead].namedStyl
		hdId := dObj.headings[ihead].id[3:]
		text := dObj.headings[ihead].text
		switch namedStyl {
		case "TITLE":
			prefix := fmt.Sprintf("<p class=\"%s_title %s_leftTitle_UL\">", dObj.docName, dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\" class=\"%s_noUl\">%s</a>", hdId, dObj.docName, text)
			suffix := "</p>\n"
			htmlStr = prefix + middle + suffix

		case "SUBTITLE":
			prefix := fmt.Sprintf("<p class=\"%s_subtitle\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</p>\n"
			htmlStr = prefix + middle + suffix

		case "HEADING_1":
			//html
			prefix := fmt.Sprintf("<h1 class=\"%s_h1 toc_h1\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h1>\n"
			htmlStr = prefix + middle + suffix
		case "HEADING_2":
			prefix := fmt.Sprintf("<h2 class=\"%s_h2 toc_h2\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h2>\n"
			htmlStr = prefix + middle + suffix
		case "HEADING_3":
			prefix := fmt.Sprintf("<h3 class=\"%s_h3 toc_h3\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h3>\n"
			htmlStr = prefix + middle + suffix
		case "HEADING_4":
			prefix := fmt.Sprintf("<h4 class=\"%s_h4 toc_h4\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h4>\n"
			htmlStr = prefix + middle + suffix
		case "HEADING_5":
			prefix := fmt.Sprintf("<h5 class=\"%s_h5 toc_h5\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h5>\n"
			htmlStr = prefix + middle + suffix
		case "HEADING_6":
			prefix := fmt.Sprintf("<h6 class=\"%s_h6 toc_h6\">", dObj.docName)
			middle := fmt.Sprintf("<a href=\"#%s\">%s</a>", hdId, text)
			suffix := "</h6>\n"
			htmlStr = prefix + middle + suffix
		case "NORMAL_TEXT":

		default:

		}
		tocDiv.bodyHtml += htmlStr

	} // end loop

	tocDiv.bodyHtml += "</div>\n"

	return &tocDiv, nil
}

func (dObj *GdocDomObj) cvtBodyToDom() (bodyObj *dispObj, err error) {

	if dObj == nil {
		return nil, fmt.Errorf("-- no GdocObj!")
	}

	doc := dObj.doc
	body := doc.Body
	if body == nil {
		return nil, fmt.Errorf("-- no doc.body!")
	}

	err = nil

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

	elNum := len(body.Content)
	for el:=0; el< elNum; el++ {
		bodyEl := body.Content[el]
		tObj, err1 := dObj.cvtContentElToDom(bodyEl)
		if err1 != nil {
			tObj.script = fmt.Sprintf("//el %d cvtContentEl: %v\n", el, err) + tObj.script
			err = fmt.Errorf("cvtContentEl: El %d %v\n", el, err)
		}
		addDispObj(bodyObj, tObj)
	} // for el loop end

	if dObj.listStack != nil {
		dObj.closeList(0)
	//fmt.Printf("end of doc closing list!")
	}

	return bodyObj, err
}

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

	if dObj.listStack != nil {
		dObj.closeList(0)
	}

	return bodyObj, nil
}

func CreGdocDomDoc(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
	// function which converts the entire document into an hmlt file

    if doc == nil { return fmt.Errorf("error -- doc is nil!\n")}
	var mainDiv dispObj
	var dObj GdocDomObj

	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization
	err = dObj.initGdocDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocHtml %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.createFootnoteDiv()
	if err != nil {
		fmt.Errorf("createFootnoteDiv: %v", err)
	}

//	dObj.sections
	secDiv := dObj.createSectionDiv()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.createSectionHeading(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
//			mainDiv.headCss += pgBody.headCss
			mainDiv.bodyCss += pgBody.bodyCss
			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
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
	tocDiv, err := dObj.createTocDiv()
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
//	htmlStr := creHtmlDocDiv(dObj.docName)
//	outfil.WriteString(htmlStr)

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

func CreGdocDomMain(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that converts the main part of a gdoc document into an html file
// excludes everything before the "main" heading or
// excludes sections titled "summary" and "keywords"

	var mainDiv dispObj
	var dObj GdocDomObj

	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization
	err = dObj.initGdocDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.createFootnoteDiv()
	if err != nil {
		fmt.Errorf("createFootnoteDiv: %v", err)
	}

//	dObj.sections
	secDiv := dObj.createSectionDiv()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.createSectionHeading(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
//			mainDiv.headCss += pgBody.headCss
			mainDiv.bodyCss += pgBody.bodyCss
			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
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
	tocDiv, err := dObj.createTocDiv()
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


func CreGdocDomSection(heading, folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that creates an html fil from the named section

	var mainDiv dispObj
	var dObj GdocDomObj

	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization
	err = dObj.initGdocDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.createFootnoteDiv()
	if err != nil {
		fmt.Errorf("createFootnoteDiv: %v", err)
	}

//	dObj.sections
	secDiv := dObj.createSectionDiv()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.createSectionHeading(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
//			mainDiv.headCss += pgBody.headCss
			mainDiv.bodyCss += pgBody.bodyCss
			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
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
	tocDiv, err := dObj.createTocDiv()
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



func CreGdocDomAll(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
// function that creates an html fil from the named section
	var mainDiv dispObj
	var dObj GdocDomObj

	// initialize dObj with doc assignment
	dObj.doc = doc

	// further initialization
	err = dObj.initGdocDom(folderPath, options)
	if err != nil {
		return fmt.Errorf("initGdocDom %v", err)
	}

// footnotes
	ftnoteDiv, err := dObj.createFootnoteDiv()
	if err != nil {
		fmt.Errorf("createFootnoteDiv: %v", err)
	}

//	dObj.sections
	secDiv := dObj.createSectionDiv()
	if secDiv != nil {
		for ipage:=0; ipage<len(dObj.sections); ipage++ {
			pgHd := dObj.createSectionHeading(ipage)
			elStart := dObj.sections[ipage].secElStart
			elEnd := dObj.sections[ipage].secElEnd
			pgBody, err := dObj.cvtBodySecToDom(elStart, elEnd)
			if err != nil {
				return fmt.Errorf("cvtBodySecToDom %d %v", ipage, err)
			}
			mainDiv.bodyCss += pgBody.bodyCss
			mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
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
	tocDiv, err := dObj.createTocDiv()
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


