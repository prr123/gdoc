// golang library that conversts gdoc to html
// author: prr
// created: 18/11/2021
// copyright 2022 Peter Riemenschneider
//
// v1
// 18/12  -- fix Color
//		-- add cvtPar
// 20/12  limit headings to four
// 21/12 introduce special ids for title and subtitle
//       made title and subtitle paragraphs instead of headings
//		assumption is that there is only one title and one subtitle
//		extended headings to six
//		added TOC division
// 22/12 added method CreateTocHead
//		need to add parstyle normal text to toc div
//
// links
//
// 15/2/2022 add ordered lists
//
// 28/2/2022 add inline images
//
// 2/3/2022 reduced css by eliminating paragraph element where default values are repeated
// 	 -indentfirst:0
//	  -indent:0
//  - line-height: 1.15
//
// 9/3/2022 add img sub folders
//          add positioned images
//		    add conversion options
// 12/3/022 add table
//

package gdocToHtml

import (
	"fmt"
	"os"
	"net/http"
	"io"
	"unicode/utf8"
	"google.golang.org/api/docs/v1"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)

type GdocHtmlObj struct {
	doc *docs.Document
//   	divClass string
	docName string
    docWidth float64
	docHeight float64
	ImgFoldName string
    ImgCount int
    tableCount int
    parCount int
    h1Count int
    h2Count int
    h3Count int
    h4Count int
    h5Count int
    h6Count int
	title namStyl
	subtitle namStyl
	h1 namStyl
	h2 namStyl
	h3 namStyl
	h4 namStyl
	h5 namStyl
	h6 namStyl
    spanCount int
//	pnorm namStyl
	hasList bool
//    subDivCount int
	listCount int
    lists *[]listObj
	listCssStr string
    cNestLev int64
    cListOr bool
	cListId string
	numHeaders int
	headers *[]string
	headCount int
	secCount int
	elCount int
	ftNoteCount int
	inImgCount int
	posImgCount int
//	df_ls float64
	folder *os.File
    imgFoldNam string
    imgFoldPath string
	Options *OptObj
}

type OptObj struct {
	CssFil bool
	ImgFold bool
    Verb bool
	Toc bool
	multDiv bool
	DivBorders bool
}

type dispObj struct {
	headCss string
	bodyHtml string
	bodyCss string
	tocHtml string
	tocCss string
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
}

type listObj struct {
    UnTyp bool
    Id string
	numNestLev int
	NestLev [9]nestLevel
}

type tabCell struct {
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

type tblBorder struct {
	color string
	dash string
	width float64
}

type parStyl struct {
	align string
	txtAl string
	lineHt float64
	hasBorder bool
}

type cssAtt struct {
	ValFl float64
	ValInt int
	AttNam string
	ValStr string
}

type namStyl struct {
	count int
	exist bool
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
	spaceMode string
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


func fillTextMap(txtMap *textMap, txtStyl *docs.TextStyle)(alter bool, err error) {

	alter = false
	if txtStyl == nil {
		return alter, fmt.Errorf("error decode txtstyle: -- no Style")
	}

	if txtStyl.BaselineOffset != txtMap.baseOffset {
		txtMap.baseOffset = txtStyl.BaselineOffset
		alter = true
	}

	if txtStyl.Bold != txtMap.bold {
		txtMap.bold = txtStyl.Bold
		alter = true
	}

	if txtStyl.Italic != txtMap.italic {
		txtMap.italic = txtStyl.Italic
		alter = true
	}

	if txtStyl.Underline != txtMap.underline {
 		txtMap.underline = txtStyl.Underline
		alter = true
	}

	if txtStyl.Strikethrough != txtMap.strike {
		txtMap.strike = txtStyl.Strikethrough
		alter = true
	}

	if txtStyl.WeightedFontFamily != nil {
		if txtStyl.WeightedFontFamily.FontFamily != txtMap.fontType {
			txtMap.fontType = txtStyl.WeightedFontFamily.FontFamily
			alter = true
		}
		if txtStyl.WeightedFontFamily.Weight != txtMap.fontWeight {
			txtMap.fontWeight = txtStyl.WeightedFontFamily.Weight
			alter = true
		}
	}
	if txtStyl.FontSize != nil {
		if txtStyl.FontSize.Magnitude != txtMap.fontSize {
			txtMap.fontSize = txtStyl.FontSize.Magnitude
			alter = true
		}
	}

	if txtStyl.ForegroundColor != nil {
		if txtStyl.ForegroundColor.Color != nil {
			color := getColor(txtStyl.ForegroundColor.Color)
			if color != txtMap.txtColor {
				txtMap.txtColor = color
				alter = true
			}
		}
	}

	if txtStyl.BackgroundColor != nil {
		if txtStyl.BackgroundColor.Color != nil {
			color := getColor(txtStyl.BackgroundColor.Color)
			if color != txtMap.bckColor {
				txtMap.bckColor = color
				alter = true
			}
		}
	}


	return alter, nil
}

func fillParMap(parmap *parMap, parStyl *docs.ParagraphStyle)(alter bool, err error) {

	alter = false
	if parStyl == nil {
		return alter, fmt.Errorf("error fillParMap: no parStyl!")
	}

	if parStyl.Alignment != parmap.halign {parmap.halign = parStyl.Alignment}
	parmap.direct = true
	if (parStyl.IndentStart != nil) {
		if parStyl.IndentStart.Magnitude != parmap.indStart {
			parmap.indStart = parStyl.IndentStart.Magnitude
			alter = true
		}
	}
	if (parStyl.IndentEnd != nil) {
		if parStyl.IndentEnd.Magnitude != parmap.indEnd {
			parmap.indEnd = parStyl.IndentEnd.Magnitude
			alter = true
		}
	}
	if (parStyl.IndentFirstLine != nil) {
		if parStyl.IndentFirstLine.Magnitude != parmap.indFlin {
			parmap.indFlin = parStyl.IndentFirstLine.Magnitude
			alter = true
		}
	}
	if parStyl.LineSpacing/100.0 != parmap.linSpac {parmap.linSpac = parStyl.LineSpacing/100.0; alter = true;}
	if parStyl.KeepLinesTogether != parmap.keepLines {parmap.keepLines = parStyl.KeepLinesTogether; alter = true;}
	if parStyl.KeepWithNext != parmap.keepNext {parmap.keepNext = parStyl.KeepWithNext; alter = true;}

	if (parStyl.SpaceAbove != nil) {
		if parStyl.SpaceAbove.Magnitude != parmap.spaceTop {
			parmap.spaceTop = parStyl.SpaceAbove.Magnitude
			alter = true
		}
	}
	if (parStyl.SpaceBelow != nil) {
		if parStyl.SpaceBelow.Magnitude != parmap.spaceBelow {
			parmap.spaceBelow = parStyl.SpaceBelow.Magnitude
			alter = true
		}
	}

	if parStyl.SpacingMode != parmap.spaceMode {parmap.spaceMode = parStyl.SpacingMode; alter = true;}

	//tabs to do
	parmap.hasBorders = true

	bb := true
	bb = bb && (parStyl.BorderBetween == nil)
	bb = bb && (parStyl.BorderTop == nil)
	bb = bb && (parStyl.BorderRight == nil)
	bb = bb && (parStyl.BorderBottom == nil)
	bb = bb && (parStyl.BorderLeft == nil)
	if bb {parmap.hasBorders = false}

	if bb { return alter, nil }

	if parStyl.BorderBetween != nil {
		if parStyl.BorderBetween.Width != nil {
			if parStyl.BorderBetween.Width.Magnitude != parmap.bordBet.width {
				parmap.bordBet.width = parStyl.BorderBetween.Width.Magnitude
				alter = true
			}
		}
		if parStyl.BorderBetween.Padding != nil {
			if parStyl.BorderBetween.Padding.Magnitude != parmap.bordBet.pad {
				parmap.bordBet.pad = parStyl.BorderBetween.Padding.Magnitude
				alter = true
			}
		}
		if parStyl.BorderBetween.Color != nil {
			if parStyl.BorderBetween.Color.Color != nil {
				color := getColor(parStyl.BorderBetween.Color.Color)
				if color != parmap.bordBet.color {
					parmap.bordBet.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderBetween.DashStyle != parmap.bordBet.dash {parmap.bordBet.dash = parStyl.BorderBetween.DashStyle; alter = true;}
	}

	if parStyl.BorderTop != nil {
		if parStyl.BorderTop.Width != nil {
			if parStyl.BorderTop.Width.Magnitude != parmap.bordTop.width {
				parmap.bordTop.width = parStyl.BorderTop.Width.Magnitude
				alter = true
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
				color := getColor(parStyl.BorderTop.Color.Color)
				if color != parmap.bordTop.color {
					parmap.bordTop.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderTop.DashStyle != parmap.bordTop.dash {parmap.bordTop.dash = parStyl.BorderTop.DashStyle; alter = true;}
	}

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
				color := getColor(parStyl.BorderRight.Color.Color)
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
				color := getColor(parStyl.BorderBottom.Color.Color)
				if color != parmap.bordBot.color {
					parmap.bordBot.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderBottom.DashStyle != parmap.bordBot.dash {parmap.bordBot.dash = parStyl.BorderBottom.DashStyle; alter = true;}
	}

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
				color := getColor(parStyl.BorderLeft.Color.Color)
				if color != parmap.bordLeft.color {
					parmap.bordLeft.color = color
					alter = true
				}
			}
		}
		if parStyl.BorderLeft.DashStyle != parmap.bordLeft.dash {parmap.bordLeft.dash = parStyl.BorderLeft.DashStyle; alter = true;}
	}

	bb2 := true
	bb2 = bb2 && (parmap.bordBet.width == 0.0)
	bb2 = bb2 && (parmap.bordTop.width == 0.0)
	bb2 = bb2 && (parmap.bordRight.width == 0.0)
	bb2 = bb2 && (parmap.bordBot.width == 0.0)
	bb2 = bb2 && (parmap.bordLeft.width == 0.0)

	if bb2 {parmap.hasBorders = false}


	return alter, nil
}

func cvtParMapCss(pMap *parMap)(cssStr string) {

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

	if pMap.indFlin > 0.0 {
		cssStr += fmt.Sprintf("  text-indent: %.2fpt;\n", pMap.indFlin)
	}
	if pMap.indStart > 0.0 {
		cssStr += fmt.Sprintf("  margin-left: %.2fpt;\n", pMap.indStart)
	}
	if pMap.indEnd > 0.0 {
		cssStr += fmt.Sprintf("  margin-right: %.2fpt;\n", pMap.indEnd)
	}

// need to investigate
	if pMap.linSpac > 0.0 {
		cssStr += fmt.Sprintf("  line-height: %.2f;\n", pMap.linSpac)
	}

	if pMap.spaceTop > 0.0 {
		cssStr += fmt.Sprintf("  margin-top: %.2fpt;\n", pMap.spaceTop)
	}

	if pMap.spaceBelow > 0.0 {
		cssStr += fmt.Sprintf("  margin-bottom: %.2fpt;\n", pMap.spaceBelow)
	}

	if !pMap.hasBorders { return cssStr }
	cssStr += fmt.Sprintf("  padding: %.1fpt %.1fpt %.1fpt %.1fpt;\n", pMap.bordTop.pad, pMap.bordRight.pad, pMap.bordBot.pad, pMap.bordLeft.pad)
	cssStr += fmt.Sprintf("  border-top: %.1fpt %s %s;\n", pMap.bordTop.width, getDash(pMap.bordTop.dash), pMap.bordTop.color)
	cssStr += fmt.Sprintf("  border-right: %.1fpt %s %s;\n", pMap.bordRight.width, getDash(pMap.bordRight.dash), pMap.bordRight.color)
	cssStr += fmt.Sprintf("  border-bottom: %.1fpt %s %s;\n", pMap.bordBot.width, getDash(pMap.bordBot.dash), pMap.bordBot.color)
	cssStr += fmt.Sprintf("  border-left: %.1fpt %s %s;\n", pMap.bordLeft.width, getDash(pMap.bordLeft.dash), pMap.bordLeft.color)

	return cssStr
}

func cvtTextMapCss(txtMap *textMap)(cssStr string) {

	cssStr =""
	if len(txtMap.baseOffset) > 0 {
		switch txtMap.baseOffset {
			case "SUPERSCRIPT":
				cssStr += "  vertical-align: sub;\n"
			case "SUBSCRIPT":
				cssStr += "	vertical-align: sup;\n"
			case "NONE":

			default:
			//error
				cssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtMap.baseOffset)
		}
	}
	if txtMap.bold {
		cssStr += "  font-weight: bold;\n"
	}
	if txtMap.italic { cssStr += "  font-style: italic;\n"}
	if txtMap.underline { cssStr += "  text-decoration: underline;\n"}
	if txtMap.strike { cssStr += "  text-decoration: line-through;\n"}

	if len(txtMap.fontType) >0 { cssStr += fmt.Sprintf("  font-family: %s;\n", txtMap.fontType)}
	if txtMap.fontWeight > 0 {cssStr += fmt.Sprintf("  font-weight: %d;\n", txtMap.fontWeight)}
	if txtMap.fontSize >0 {cssStr += fmt.Sprintf("  font-size: %.2fpt;\n", txtMap.fontSize)}
	if len(txtMap.txtColor) >0 {cssStr += fmt.Sprintf("  color: %s;\n", txtMap.txtColor)}
	if len(txtMap.bckColor) >0 {cssStr += fmt.Sprintf("  background-color: %s;\n", txtMap.bckColor)}

	return cssStr
}


func addDispObj(src, add *dispObj) {
	src.headCss += add.headCss
	src.bodyHtml += add.bodyHtml
	src.bodyCss += add.bodyCss
	src.tocHtml += add.tocHtml
	src.tocCss += add.tocCss
	return
}

func getColor(color  *docs.Color)(outstr string) {
    outstr = ""
        if color != nil {
            blue := int(color.RgbColor.Blue*255.0)
            red := int(color.RgbColor.Red*255.0)
            green := int(color.RgbColor.Green*255)
            outstr += fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
            return outstr
        }
    outstr = "/*no color*/\n"
    return outstr
}

func getDash(dashStyle string)(outstr string) {

	switch dashStyle {
		case "SOLID":
			outstr = "solid"
		case "DOT":
			outstr = "dotted"
		case "DASH":
			outstr = "dashed"
		default:
			outstr = "none"
	}

    return outstr
}

func getImgLayout (layout string) (ltyp int, err error) {

	switch layout {
		case "WRAP_TEXT":

		case "BREAK_LEFT":

		case "BREAK_RIGHT":

		case "BREAK_LEFT_RIGHT":

		case "IN_FRONT_OF_TEXT":

		case "BEHIND_TEXT":

		default:
			return -1, fmt.Errorf("error getImgLayout layout %s not implemented!", layout)
	}
	return ltyp, nil
}

func tcell_vert_align (alStr string) (outstr string) {
	switch alStr {
		case "TOP":
			outstr = "top"
		case "Middle":
		 	outstr = "middle"
		case "BOTTOM":
			outstr = "bottom"
		default:
			outstr = "baseline"
	}
	return outstr
}

func (dObj *GdocHtmlObj) creHtmlHead()(outstr string, err error) {
	if dObj == nil {return "", fmt.Errorf("error creHtmlHead:: no GdocHtml Object!") }
	outstr = "<!DOCTYPE html>\n"
	outstr += fmt.Sprintf("<!-- file: %s -->\n", dObj.docName)
	outstr += fmt.Sprintf("<!-- img folder: %s -->\n",dObj.ImgFoldName)
	outstr += "<head>\n<style>\n"
	return outstr, nil
}

//xxnam
func (dObj *GdocHtmlObj) getNamedStyl(namedTyp string)(parStyl *docs.ParagraphStyle, txtStyl *docs.TextStyle, err error) {
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
		return nil, nil, fmt.Errorf("error getNamedStyl: cannot find named style %s!", namedTyp)
	}

	parStyl = namStyl.ParagraphStyle
	txtStyl = namStyl.TextStyle
	return parStyl, txtStyl, nil
}


func (dObj *GdocHtmlObj) downloadImg()(err error) {

    doc := dObj.doc
	verb := dObj.Options.Verb
    if !(len(dObj.imgFoldNam) >0) {
        return fmt.Errorf("error downloadImg:: no imgfolder found!")
    }
    imgFoldPath := dObj.imgFoldPath + "/"
    fmt.Println("image dir: ", imgFoldPath)

    fmt.Printf("*** Inline Imgs: %d ***\n", len(doc.InlineObjects))
    for k, inlObj := range doc.InlineObjects {
        imgProp := inlObj.InlineObjectProperties.EmbeddedObject.ImageProperties
        if verb {
            fmt.Printf("Source: %s Obj %s\n", k, imgProp.SourceUri)
            fmt.Printf("Content: %s Obj %s\n", k, imgProp.ContentUri)
        }
        if !(len(imgProp.SourceUri) > 0) {
            return fmt.Errorf("error downloadImg:: image %s has no URI\n", k)
        }
        imgNam := imgFoldPath + k[4:] + ".jpeg"
        if verb {fmt.Printf("image path: %s\n", imgNam)}
        URL := imgProp.ContentUri
        httpResp, err := http.Get(URL)
        if err != nil {
            return fmt.Errorf("error downloadImg:: could not fetch %s! %v", URL, err)
        }
        defer httpResp.Body.Close()
//  fmt.Printf("http got %s!\n", URL)
        if httpResp.StatusCode != 200 {
            return fmt.Errorf("error downloadImg:: Received non 200 response code %d!", httpResp.StatusCode)
        }
//  fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file
        outfil, err := os.Create(imgNam)
        if err != nil {
            return fmt.Errorf("error downloadImg:: cannot create img file! %v", err)
        }
        defer outfil.Close()
//  fmt.Println("created dir")
        //Write the bytes to the fiel
        _, err = io.Copy(outfil, httpResp.Body)
        if err != nil {
            return fmt.Errorf("error downloadImg:: cannot copy img file content! %v", err)
        }
    }


   fmt.Printf("*** Positioned Imgs: %d ***\n", len(doc.PositionedObjects))
    for k, posObj := range doc.PositionedObjects {
        imgProp := posObj.PositionedObjectProperties.EmbeddedObject.ImageProperties
        if verb {
            fmt.Printf("Source: %s Obj %s\n", k, imgProp.SourceUri)
            fmt.Printf("Content: %s Obj %s\n", k, imgProp.ContentUri)
        }
        if !(len(imgProp.SourceUri) > 0) {
            return fmt.Errorf("error downloadImg:: image %s has no URI\n", k)
        }
        imgNam := imgFoldPath + k[4:] + ".jpeg"
        if verb {fmt.Printf("image path: %s\n", imgNam)}
        URL := imgProp.ContentUri
        httpResp, err := http.Get(URL)
        if err != nil {
            return fmt.Errorf("error downloadImg:: could not fetch %s! %v", URL, err)
        }
        defer httpResp.Body.Close()
//  fmt.Printf("http got %s!\n", URL)
        if httpResp.StatusCode != 200 {
            return fmt.Errorf("error downloadImg:: Received non 200 response code %d!", httpResp.StatusCode)
        }
//  fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file


        outfil, err := os.Create(imgNam)
        if err != nil {
            return fmt.Errorf("error downloadImg:: cannot create img file! %v", err)
        }
        defer outfil.Close()
//  fmt.Println("created dir")
        //Write the bytes to the fiel
        _, err = io.Copy(outfil, httpResp.Body)
        if err != nil {
            return fmt.Errorf("error downloadImg:: cannot copy img file content! %v", err)
        }
    }

    return nil
}


func (dObj *GdocHtmlObj) createImgFolder()(err error) {

    filnam :=dObj.docName
    if len(filnam) < 2 {
        return fmt.Errorf("error createImgFolder:: filename %s too short!", filnam)
    }

    bf := []byte(filnam)
    // replace empty space with underscore
    for i:= 0; i< len(filnam); i++ {
        if bf[i] == ' ' {bf[i]='_'}
        if bf[i] == '.' {
            return fmt.Errorf("error createImgFolder:: filnam has period!")
        }
    }

    imgFoldNam := "imgs_" + string(bf)

    fmt.Println("output file name: ", dObj.folder.Name())
    foldNamb := []byte(dObj.folder.Name())
    idx := 0
    for i:=len(foldNamb)-1; i> 0; i-- {
        if foldNamb[i] == '/' {
            idx = i
            break
        }
    }

    imgFoldPath := imgFoldNam
    if idx > 0 {
        imgFoldPath = string(foldNamb[:idx]) + "/" + imgFoldNam
    }

    fmt.Println("img folder path: ", imgFoldPath)

    // check whether dir folder exists, if not create one
    newDir := false
    _, err = os.Stat(imgFoldPath)
    if os.IsNotExist(err) {
        err1 := os.Mkdir(imgFoldPath, os.ModePerm)
        if err1 != nil {
            return fmt.Errorf("error createImgFolder:: could not create img folder! %v", err1)
        }
        newDir = true
    } else {
        if err != nil {
            return fmt.Errorf("error createImgFolder:: could not find img folder! %v", err)
        }
    }

    // open directory
    if !newDir {
        err = os.RemoveAll(imgFoldPath)
        if err != nil {
            return fmt.Errorf("error createImgFolder:: could not delete files in image folder! %v", err)
        }
        err = os.Mkdir(imgFoldPath, os.ModePerm)
        if err != nil {
            return fmt.Errorf("error createImgFolder:: could not create img folder! %v", err)
        }
    }
    dObj.imgFoldNam = imgFoldNam
    dObj.imgFoldPath = imgFoldPath

    return nil
}

func (dObj *GdocHtmlObj) disp_GdocHtmlObj (dbgfil *os.File) (err error) {
	var outstr string

  	if dObj == nil {
        return fmt.Errorf("error disp_GdocHtmlObj -- dObj is nil!")
    }
	if dbgfil == nil {
        return fmt.Errorf("error disp_GdocHtmlObj -- dbggil is nil!")
    }

	outstr = fmt.Sprintf("Document: %s\n", dObj.docName)
	outstr += "Lists:\n"

	dbgfil.WriteString(outstr)
	return nil
}

func (dObj *GdocHtmlObj) findListIndex (listId string) (listIdx int, err error) {

	if len(listId) < 1 {
		return -1, fmt.Errorf("error findListIndex: no listId string!")
	}

	listIdx = -1
	for i:=0; i< dObj.listCount; i++ {
		if (*dObj.lists)[i].Id == listId {
			listIdx = i
			break
		}
	}
	if !(listIdx >-1) {
		return listIdx, fmt.Errorf("error findListIndex: list index not found in Lists!")
	}

	return listIdx, nil
}

func (dObj *GdocHtmlObj) findListProp (listId string) (listProp *docs.ListProperties) {

	listIdx := 0
	doc := dObj.doc

	for key, listItem := range doc.Lists  {
		if listId == key {
			listProp = listItem.ListProperties
			break
		}
		listIdx++
	}
	if listIdx > 0 {
		return listProp
	}

	return nil
}

func (dObj *GdocHtmlObj) InitGdocHtmlLib (doc *docs.Document, opt *OptObj) (err error) {
//	var normStyl *docs.NamedStyle
	var defOpt OptObj
	dObj.doc = doc

	defOpt.Verb = true
	defOpt.ImgFold = true
	defOpt.Toc = false
	defOpt.CssFil = false

	if opt == nil {
		dObj.Options = &defOpt
	}
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

// initialise heading use
// each heading has a default paragraph style
// no reason to  create css if heading is not used
	dObj.title.exist = false
	dObj.subtitle.exist = false
	dObj.h1.exist = false
	dObj.h2.exist = false
	dObj.h3.exist = false
	dObj.h4.exist = false
	dObj.h5.exist = false
	dObj.h6.exist = false
	dObj.hasList = false

// section breaks
	dObj.elCount = len(doc.Body.Content)
	dObj.secCount = 0
	dObj.ftNoteCount = 0

	for el:=0; el<dObj.elCount; el++ {
		elObj:= doc.Body.Content[el]
		if elObj.SectionBreak != nil {
			if elObj.SectionBreak.SectionStyle.SectionType == "NEXT_PAGE" {dObj.secCount++}
		}
		if elObj.Paragraph != nil {
			for parEl:=0; parEl<len(elObj.Paragraph.Elements); parEl++ {
				parElObj := elObj.Paragraph.Elements[parEl]
				if parElObj.FootnoteReference != nil {dObj.ftNoteCount++}
			}
		}
	}

// list initalisation
// need to separate ordered and unordered lists

	dObj.listCount = len(doc.Lists)
	il := 0
	dlists := make([]listObj,dObj.listCount)
	dObj.lists = &dlists
	for idstr, list := range doc.Lists {
		dlists[il].Id = idstr

		numNestLev := len(list.ListProperties.NestingLevels)
		dlists[il].numNestLev = numNestLev

		for nl:=0; nl<numNestLev; nl++ {
			nestLevel := list.ListProperties.NestingLevels[nl]
			dlists[il].NestLev[nl].GlAl = nestLevel.BulletAlignment
			dlists[il].NestLev[nl].GlFmt = nestLevel.GlyphFormat
			dlists[il].NestLev[nl].GlSym = nestLevel.GlyphSymbol
			dlists[il].NestLev[nl].GlTyp = nestLevel.GlyphType
			dlists[il].NestLev[nl].FlInd = nestLevel.IndentFirstLine.Magnitude
			dlists[il].NestLev[nl].StInd = nestLevel.IndentStart.Magnitude
			dlists[il].NestLev[nl].Count = nestLevel.StartNumber
			ord := false
			switch nestLevel.GlyphType {
				case "DECIMAL":
					ord = true
				case "ZERO_DECIMAL":
					ord = true
				case "UPPER_ALPHA":
					ord = true
				case "ALPHA":
					ord = true
				case "UPPER_ROMAN":
					ord = true
				case "ROMAN":
					ord = true
				default:
					ord = false
			}
			dlists[il].NestLev[nl].GlOrd = ord
		}
		il++
	}

// Headers
	dObj.numHeaders = len(doc.Headers)
	dheaders := make(map[string]int,len(doc.Headers))
	for k, _ := range doc.Headers {
		dheaders[k] = 0
	}

// images
	dObj.inImgCount = len(doc.InlineObjects)
	dObj.posImgCount = len(doc.PositionedObjects)
    totObjNum := dObj.inImgCount + dObj.posImgCount
    dObj.parCount = len(doc.Body.Content)

	if (dObj.Options.ImgFold) && (totObjNum > 0){
		if dObj.Options.Verb {fmt.Printf("*** creating image folder for %d images ***\n", totObjNum)}

    	err = dObj.createImgFolder()
		if err != nil {
        	return fmt.Errorf("error InitGdocHtmlLib: could not create ImgFolder: %v!", err)
		}
    	err = dObj.downloadImg()
    	if err != nil {
        	return fmt.Errorf("error InitGdocHtmlLib: could not download images: %v!", err)
    	}
	} else {
		if dObj.Options.Verb {fmt.Printf("*** no image folder created ***\n")}
	}
	return nil
}


func (dObj *GdocHtmlObj) convertGlyph(nlev *docs.NestingLevel, ord bool)(cssStr string) {

	var glyphTyp string
	if ord {
	// ordered list
		switch nlev.GlyphType {
			case "DECIMAL":
				glyphTyp = "decimal"
			case "ZERO_DECIMAL":
				glyphTyp = "decimal-leading-zero"

			case "ALPHA":
				glyphTyp = "lower-alpha"

 			case "UPPER_ALPHA":
				glyphTyp = "upper-alpha"

			case "ROMAN":
				glyphTyp = "lower-roman"

			case "UPPER_ROMAN":
				glyphTyp = "upper-roman"

			default:
				cssStr = "/* unknown GlyphType */\n"
		}
		if len(glyphTyp) > 0 {
			cssStr = "  list-style-type:" + glyphTyp +";\n"
			cssStr +="  list-style-position: inside;\n"
			cssStr +="  padding-left: 0;\n"
		}
	} else {
	// unordered list
		cssStr =fmt.Sprintf("/*-Glyph Symbol:%x - */\n",nlev.GlyphSymbol)
		r, _ := utf8.DecodeRuneInString(nlev.GlyphSymbol)

		switch r {
			case 9679:
				glyphTyp = "disc"

			case 9675:
				glyphTyp = "circle"

			case 9632:
				glyphTyp = "square"

			default:

		}
		if len(glyphTyp) > 0 {
			cssStr = "  list-style-type:" + glyphTyp +";\n"
			cssStr +="  list-style-position: inside;\n"
			cssStr +="  padding-left: 0;\n"
		}
	}

	return cssStr
}

func (dObj *GdocHtmlObj) cvtInlineImg(imgEl *docs.InlineObjectElement)(htmlStr string, cssStr string, err error) {

	if imgEl == nil {
		return "","", fmt.Errorf("error cvtInlineImg:: imgEl is nil!")
	}
	doc := dObj.doc

	imgElId := imgEl.InlineObjectId
	if !(len(imgElId) > 0) {return "","", fmt.Errorf("error cvtInlineImg:: no InlineObjectId found!")}

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
	htmlStr = fmt.Sprintf("<!-- inline image %s -->\n", imgElId)
	imgObj := doc.InlineObjects[imgElId].InlineObjectProperties.EmbeddedObject

	if dObj.Options.ImgFold {
    	imgSrc := dObj.imgFoldNam + "/" + imgId + ".jpeg"
		htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgSrc, imgId, imgObj.Title)
	} else {
		htmlStr +=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\">\n", imgObj.ImageProperties.SourceUri, imgId, imgObj.Title)
	}
	cssStr = fmt.Sprintf("#%s {\n",imgId)
	cssStr += fmt.Sprintf(" width:%.1fpt; height:%.1fpt; \n}\n", imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )
	// todo add margin
	return htmlStr, cssStr, nil
}

func (dObj *GdocHtmlObj) cvtParElText(parEl *docs.ParagraphElement)(htmlStr string, cssStr string, err error) {

   if parEl == nil {
        return "","", fmt.Errorf("error cvtPelText -- parEl is nil!")
    }
   if dObj == nil {
        return "","", fmt.Errorf("error cvtPelText -- dObj is nil!")
    }


    cLen := len(parEl.TextRun.Content)
	if (cLen == 1) {
		let := parEl.TextRun.Content
		if (let == "\n") {
			htmlStr += "<br>"
			return htmlStr, "", nil
		}
	}
	spanCssStr, err := dObj.cvtTxtStylCss(parEl.TextRun.TextStyle, false)
	if err != nil {
		spanCssStr = fmt.Sprintf("/*error parEl Css %v*/\n", err) + spanCssStr
	}
	dObj.spanCount++
	linkPrefix := ""
	linkSuffix := ""
	if parEl.TextRun.TextStyle.Link != nil {
		if len(parEl.TextRun.TextStyle.Link.Url)>0 {
			linkPrefix = "<a href = \"" + parEl.TextRun.TextStyle.Link.Url + "\">"
			linkSuffix = "</a>"
		}
	}
	spanIdStr := fmt.Sprintf("%s_sp%d", dObj.docName, dObj.spanCount)
	if len(spanCssStr)>0 {
		cssStr = fmt.Sprintf("#%s {\n", spanIdStr) + spanCssStr + "}\n"
		htmlStr = fmt.Sprintf("<span id=\"%s\">",spanIdStr) + linkPrefix + parEl.TextRun.Content + linkSuffix + "</span>"
	} else {
		htmlStr = linkPrefix + parEl.TextRun.Content + linkSuffix
	}
	return htmlStr, cssStr, nil
}


func (dObj *GdocHtmlObj) closeList(nl int64)(htmlStr string) {
	var i int64
	for i =0; i <nl+1; i++ {
		if dObj.cListOr {
			htmlStr += "</ol>\n"
		} else {
			htmlStr +="</ul>\n"
		}
	}
	return htmlStr
}

func (dObj *GdocHtmlObj) renderPosImg(posImg docs.PositionedObject, posId string)(htmlStr, cssStr string, err error) {

	posObjProp := posImg.PositionedObjectProperties
	imgProp := posObjProp.EmbeddedObject
	htmlStr += fmt.Sprintf("\n<!-- Positioned Image %s -->\n", posId)
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

	return htmlStr, cssStr, nil
}

// table element
// 	tObj, _ := dObj.cvtTableToHtml(tableEl)

func (dObj *GdocHtmlObj) cvtTable(tbl *docs.Table)(tabObj dispObj, err error) {
	var htmlStr, cssStr string
	var tabWidth float64
	var icol, trow int64
	var defcel tabCell


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
		defcel.vert_align = tcell_vert_align(tcelDefStyl.ContentAlignment)

// if left border is the only border specified, let's use it for default values
		tb := (tcelDefStyl.BorderTop == nil)&& (tcelDefStyl.BorderRight == nil)
		tb = tb&&(tcelDefStyl.BorderBottom == nil)
		if (tcelDefStyl.BorderLeft != nil) && tb {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = getColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = getDash(tcelDefStyl.BorderLeft.DashStyle)
		}

		if tcelDefStyl.PaddingTop != nil {defcel.pad[0] = tcelDefStyl.PaddingTop.Magnitude}
		if tcelDefStyl.PaddingRight != nil {defcel.pad[1] = tcelDefStyl.PaddingRight.Magnitude}
		if tcelDefStyl.PaddingBottom != nil {defcel.pad[2] = tcelDefStyl.PaddingBottom.Magnitude}
		if tcelDefStyl.PaddingLeft != nil {defcel.pad[3] = tcelDefStyl.PaddingLeft.Magnitude}

		if tcelDefStyl.BackgroundColor != nil {defcel.bckcolor = getColor(tcelDefStyl.BackgroundColor.Color)}

		if tcelDefStyl.BorderTop != nil {
			if tcelDefStyl.BorderTop.Color != nil {defcel.border[0].color = getColor(tcelDefStyl.BorderTop.Color.Color)}
			if tcelDefStyl.BorderTop.Width != nil {defcel.border[0].width = tcelDefStyl.BorderTop.Width.Magnitude}
			defcel.border[0].dash = getDash(tcelDefStyl.BorderTop.DashStyle)
		}
		if tcelDefStyl.BorderRight != nil {
			if tcelDefStyl.BorderRight.Color != nil {defcel.border[1].color = getColor(tcelDefStyl.BorderRight.Color.Color)}
			if tcelDefStyl.BorderRight.Width != nil {defcel.border[1].width = tcelDefStyl.BorderRight.Width.Magnitude}
			defcel.border[1].dash = getDash(tcelDefStyl.BorderRight.DashStyle)
		}
		if tcelDefStyl.BorderBottom != nil {
			if tcelDefStyl.BorderBottom.Color != nil {defcel.border[2].color = getColor(tcelDefStyl.BorderBottom.Color.Color)}
			if tcelDefStyl.BorderBottom.Width != nil {defcel.border[2].width = tcelDefStyl.BorderBottom.Width.Magnitude}
			defcel.border[2].dash = getDash(tcelDefStyl.BorderBottom.DashStyle)
		}
		if tcelDefStyl.BorderLeft != nil {
			if tcelDefStyl.BorderLeft.Color != nil {defcel.border[3].color = getColor(tcelDefStyl.BorderLeft.Color.Color)}
			if tcelDefStyl.BorderLeft.Width != nil {defcel.border[3].width = tcelDefStyl.BorderLeft.Width.Magnitude}
			defcel.border[3].dash = getDash(tcelDefStyl.BorderLeft.DashStyle)
		}
		if tcelDefStyl.BorderTop == tcelDefStyl.BorderRight {
			fmt.Println("same border!")
			if tcelDefStyl.BorderTop != nil {
				if tcelDefStyl.BorderTop.Color != nil {defcel.bcolor = getColor(tcelDefStyl.BorderTop.Color.Color)}
				defcel.bdash = getDash(tcelDefStyl.BorderTop.DashStyle)
				if tcelDefStyl.BorderTop.Width != nil {defcel.bwidth = tcelDefStyl.BorderTop.Width.Magnitude}
			}
		}
	}

	//set up table
	tblClass := fmt.Sprintf("%s_tbl", dObj.docName)
	tblCellClass := fmt.Sprintf("%s_tcel", dObj.docName)
	htmlStr = fmt.Sprintf("<table class=\"%s\">\n", tblClass)

  // table styling
  	cssStr = fmt.Sprintf(".%s {\n",tblClass)
 	cssStr += fmt.Sprintf("  border: 1px solid black;\n  border-collapse: collapse;\n")
 	cssStr += fmt.Sprintf("  width: %.1fpt;\n", tabWidth)
	cssStr += "   margin:auto;\n"
	cssStr += "}\n"

// table columns
	tabWtyp :=tbl.TableStyle.TableColumnProperties[0].WidthType
//	if !((tabWtyp == "EVENLY_DISTRIBUTED")||(tabWtyp == "WIDTH_TYPE_UNSPECIFIED")) {
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
				if tstyl.BackgroundColor != nil {cellStr += fmt.Sprintf(" background-color:\"%s\";",getColor(tstyl.BackgroundColor.Color))}
				if tcell_vert_align(tstyl.ContentAlignment) != defcel.vert_align {cellStr += fmt.Sprintf(" vertical-align: %s;", tcell_vert_align(tstyl.ContentAlignment))}
				if tstyl.PaddingTop != nil {
					if tstyl.PaddingTop.Magnitude != defcel.pad[0] { cellStr += fmt.Sprintf(" padding-top: %5.1fpt;", tstyl.PaddingTop.Magnitude)}
				}

				if tstyl.BorderTop != nil {
					// Color
					if tstyl.BorderTop.Color != nil {cellStr += fmt.Sprintf(" border-top-color: %s;", getColor(tstyl.BorderTop.Color.Color))}
					//dash
					if getDash(tstyl.BorderTop.DashStyle) != defcel.bdash {cellStr += fmt.Sprintf(" border-top-style: %s;",  getDash(tstyl.BorderTop.DashStyle))}
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
				tObj, err:=dObj.cvtContentEl(elObj)
				if err != nil {
					tabObj.bodyHtml = htmlStr
					tabObj.bodyCss = cssStr
					return tabObj, fmt.Errorf("error ConvertTable: %v", err)
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
// paragraph element par
// - Bullet
// - Elements
// - ParagraphStyle
// - Positioned Objects
//

func (dObj *GdocHtmlObj) cvtPar(par *docs.Paragraph)(parObj dispObj, err error) {
	var parHtmlStr, parCssStr string
	var prefix, suffix string
	var tocPrefix, tocSuffix string
	var parIdStr string
	var tStr string
	var nestIdx int64
	var listStr, listCssStr, listSuffix string
	var nestinc, i, j int64

	if par == nil {
        return parObj, fmt.Errorf("error cvtPar -- parEl is nil!")
    }
	if dObj == nil {
        return parObj, fmt.Errorf("error cvttPar -- dObj is nil!")
    }
// Positioned Objects
	numPosObj := len(par.PositionedObjectIds)
	for i:=0; i< numPosObj; i++ {
		posId := par.PositionedObjectIds[i]
		posObj, ok := dObj.doc.PositionedObjects[posId]
		if !ok {return parObj, fmt.Errorf("error cvtPar: could not find positioned Object with id: ", posId)}

		imgHtmlStr, imgCssStr, err := dObj.renderPosImg(posObj, posId)
		if err != nil {
			parHtmlStr += fmt.Sprintf("<!-- error cvtPar:: render pos img %v -->\n", err) + imgHtmlStr
			parCssStr += imgCssStr
		} else {
			parHtmlStr += imgHtmlStr
			parCssStr += imgCssStr
		}
	}
// lists
    if par.Bullet != nil {
        parHtmlStr += "\n<!-- List Element -->\n"
		dObj.hasList = true
		// find list id
		listId := par.Bullet.ListId
		listIdx, err := dObj.findListIndex(listId)
		if err != nil {
			return parObj, fmt.Errorf("error cvtPar:: findListIndex: -- no list idx found -- %v!", err)
		}
		listProp := dObj.findListProp(listId)
		nestIdx = par.Bullet.NestingLevel
		listOrd := (*dObj.lists)[listIdx].NestLev[nestIdx].GlOrd
		liClaStr := listId[4:]
    	listStr += fmt.Sprintf("<!-- list id: %s OL: %t List Index: %d Nest Level: %d -->\n", listId, listOrd ,listIdx, nestIdx)
		listStr += fmt.Sprintf("<!-- CList id: %s COL: %t Nest Level: %d -->\n",dObj.cListId, dObj.cListOr, dObj.cNestLev)
		listStr += fmt.Sprintf("<!-- list class: %S -->\n", liClaStr)
		if len(dObj.cListId) == 0 {
		// create a list
			dObj.cListId = listId
			dObj.cNestLev = nestIdx
			for i=0; i< nestIdx +1; i++ {
				nestProp := listProp.NestingLevels[i]
				glyphStr := dObj.convertGlyph(nestProp, listOrd)
				idFl := nestProp.IndentFirstLine.Magnitude
				idSt := nestProp.IndentStart.Magnitude
				if listOrd {
					liClaOlStr := liClaStr + fmt.Sprintf("_ol%d",nestIdx)
					listStr +=fmt.Sprintf("<ol class=\"%s\">\n", liClaOlStr)
					listCssStr += fmt.Sprintf(".%s {\n", liClaOlStr)
					listCssStr += glyphStr
					listCssStr += "}\n"
					listCssStr += fmt.Sprintf(".%s li {\n", liClaOlStr)
					listCssStr += "  padding-top: 2pt;\n  padding-bottom: 2pt;\n"
					listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", idFl )
					listCssStr += "}\n"
					listCssStr += fmt.Sprintf(".%s li > * {\n", liClaOlStr)
					listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", (idSt - idFl - 14))
					listCssStr += "}\n"
					dObj.cListOr = true
				} else {
					liClaUlStr := liClaStr + fmt.Sprintf("_ul%d",nestIdx)
					listStr +=fmt.Sprintf("<ul class=\"%s\">\n", liClaUlStr)
					listCssStr += fmt.Sprintf(".%s {\n", liClaUlStr)
					listCssStr += glyphStr
					listCssStr += "}\n"
					listCssStr += fmt.Sprintf(".%s li {\n", liClaUlStr)
					listCssStr += "  padding-top: 2pt;\n  padding-bottom: 2pt;\n"
					listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", idFl )
					listCssStr += "}\n"
					listCssStr += fmt.Sprintf(".%s li > * {\n", liClaUlStr)
					listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", (idSt - idFl - 14))
					listCssStr += "}\n"
					dObj.cListOr = false
				}
			}
			listStr += "<li>"
			listSuffix = "</li>\n"
		} else {
			// a list already exists
			if dObj.cListId == listId {
				// more list entries
				if nestIdx < dObj.cNestLev {
					listStr += "<!-- end sub list -->\n"
					nestinc = dObj.cNestLev - nestIdx
					for j=0; j< nestinc; j++ {
						if dObj.cListOr {
							listStr +="</ol>"
						} else {
							listStr +="</ul>"
						}
					} // loop j
					listStr += "\n"
					dObj.cNestLev = nestIdx
				}
				if nestIdx > dObj.cNestLev {
					listStr += "<!-- new sub list -->\n"
					nestinc = nestIdx - dObj.cNestLev
					for j=0; j< nestinc; j++ {
						nestProp := listProp.NestingLevels[dObj.cNestLev + j+1]
						glyphStr := dObj.convertGlyph(nestProp, dObj.cListOr)
						idFl := nestProp.IndentFirstLine.Magnitude
						idSt := nestProp.IndentStart.Magnitude
						if dObj.cListOr {
							liClaOlStr := liClaStr + fmt.Sprintf("_ol%d",nestIdx)
							listStr +=fmt.Sprintf("<ol class=\"%s\">\n", liClaOlStr)
							listCssStr += fmt.Sprintf(".%s {\n", liClaOlStr)
							listCssStr += glyphStr
							listCssStr += "}\n"
							listCssStr += fmt.Sprintf(".%s li {\n", liClaOlStr)
							listCssStr += "  padding-top: 2pt;\n  padding-bottom: 2pt;\n"
							listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", idFl )
							listCssStr += "}\n"
							listCssStr += fmt.Sprintf(".%s li > * {\n", liClaOlStr)
							listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", (idSt - idFl - 14))
							listCssStr += "}\n"
						} else {
							liClaUlStr := liClaStr + fmt.Sprintf("_ul%d",nestIdx)
							listStr +=fmt.Sprintf("<ul class=\"%s\">\n", liClaUlStr)
							listCssStr += fmt.Sprintf(".%s {\n", liClaUlStr)
							listCssStr += glyphStr
							listCssStr += "}\n"
							listCssStr += fmt.Sprintf(".%s li {\n", liClaUlStr)
							listCssStr += "  padding-top: 2pt;\n  padding-bottom: 2pt;\n"
							listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", idFl )
							listCssStr += "}\n"
							listCssStr += fmt.Sprintf(".%s li > * {\n", liClaUlStr)
							listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", (idSt - idFl - 14))
							listCssStr += "}\n"
						}
					} // loop j
					dObj.cNestLev = nestIdx
				}
				// list with same nesting level
				listStr += "<li>"
				listSuffix = "</li>\n"
			} else {
			// a new list must be created
				listStr += dObj.closeList(dObj.cNestLev)
				dObj.cListId = listId
				dObj.cNestLev = nestIdx
//		nestIdx = par.Bullet.NestingLevel
//		listOrd := dObj.lists[listIdx].NestLev[nestIdx].GlOrd
				for i=0; i< nestIdx +1; i++ {
					nestProp := listProp.NestingLevels[i]
					glyphStr := dObj.convertGlyph(nestProp, listOrd)
					idFl := nestProp.IndentFirstLine.Magnitude
					idSt := nestProp.IndentStart.Magnitude
					if listOrd {
						liClaOlStr := liClaStr + fmt.Sprintf("_ol%d",nestIdx)
						listStr +=fmt.Sprintf("<ol class=\"%s\">\n", liClaOlStr)
						listCssStr += fmt.Sprintf(".%s {\n", liClaOlStr)
						listCssStr += glyphStr
						listCssStr += "}\n"
						dObj.cListOr = true
					} else {
						liClaUlStr := liClaStr + fmt.Sprintf("_ul%d",nestIdx)
						listStr +=fmt.Sprintf("<ol class=\"%s\">\n", liClaUlStr)
						listCssStr += fmt.Sprintf(".%s {\n", liClaUlStr)
						listCssStr += glyphStr
						listCssStr += "}\n"
						listCssStr += fmt.Sprintf(".%s li {\n", liClaUlStr)
						listCssStr += "  padding-top: 2pt;\n  padding-bottom: 2pt;\n"
						listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", idFl )
						listCssStr += "}\n"
						listCssStr += fmt.Sprintf(".%s li > * {\n", liClaUlStr)
						listCssStr += fmt.Sprintf("  padding-left: %.1fpt;\n", (idSt - idFl - 14))
						listCssStr += "}\n"
						dObj.cListOr = false
					}
				} // loop i
				listStr += "<li>"
				listSuffix = "</li>\n"
			} // if dObj.cListId == listId
		} //len(dObj.cListId) == 0

    } else {
		// if there was a list, close it
		if len(dObj.cListId) > 0 {
			parHtmlStr += dObj.closeList(dObj.cNestLev)
		}
		parHtmlStr += "\n<!-- Par Element -->\n"
	}

// we need to redo
	// Heading Id refers to a page heading
	if len(par.ParagraphStyle.HeadingId) > 0 {
		parHtmlStr += fmt.Sprintf("<!-- Heading Id: %s -->\n", par.ParagraphStyle.HeadingId)
//head
//		dObj.headers
	}

	decode := true
	prefix = ""
	suffix = ""
	switch par.ParagraphStyle.NamedStyleType {
        case "TITLE":
			parIdStr = fmt.Sprintf("%s_title",dObj.docName)
			titleStr := fmt.Sprintf("\"#%s\"",parIdStr)
			prefix = fmt.Sprintf("<p class=\"%s\">",parIdStr)
			tocPrefix = fmt.Sprintf("<p id=\"%s_TOC_title\"><a href = %s>",dObj.docName ,titleStr)
			suffix ="</p>"
			tocSuffix = "</a></p>"
			decode = true

        case "SUBTITLE":
			parIdStr = fmt.Sprintf("%s_subtitle",dObj.docName)
			prefix = fmt.Sprintf("<p class=\"%s\">",parIdStr)
			tocPrefix = ""
			suffix ="</p>"
			tocSuffix = ""
			decode = false

        case "HEADING_1":
			dObj.h1Count++
			parIdStr = fmt.Sprintf("%s_h1_%d", dObj.docName, dObj.h1Count)
			prefix = fmt.Sprintf("<h1 id=\"%s\" class=\"%s_h1\">", parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h1><a href = \"#%s\">",parIdStr)
			suffix ="</h1>"
			tocSuffix = "</a></h1>"
			decode = true

        case "HEADING_2":
			dObj.h2Count++
			parIdStr = fmt.Sprintf("%s_h2_%d",dObj.docName, dObj.h2Count)
			prefix = fmt.Sprintf("<h2 id=\"%s\" class=\"%s_h2\">", parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h2><a href = \"#%s\">",parIdStr)
			suffix ="</h2>"
			tocSuffix ="</a></h2>"
			decode = true

        case "HEADING_3":
			dObj.h3Count++
			parIdStr = fmt.Sprintf("%s_h3_%d",dObj.docName, dObj.h3Count)
			prefix = fmt.Sprintf("<h3 id=\"%s\" class = \"%s_h3\">",parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h3><a href = \"#%s\">",parIdStr)
			suffix ="</h3>"
			tocSuffix ="</a></h3>"
			decode = true

        case "HEADING_4":
			dObj.h4Count++
			parIdStr = fmt.Sprintf("%s_h4_%d",dObj.docName, dObj.h4Count)
			prefix = fmt.Sprintf("<h4 id=\"%s\" class=\"%s_h4\">", parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h4><a href = \"#%s\">",parIdStr)
			suffix ="</h4>"
			tocSuffix ="</a></h4>"
			decode = true

        case "HEADING_5":
			dObj.h5Count++
			parIdStr = fmt.Sprintf("%s_h5_%d",dObj.docName, dObj.h5Count)
			prefix = fmt.Sprintf("<h5 id=\"%s\" class=\"%s_h5\">", parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h5><a href = \"#%s\">",parIdStr)
			suffix ="</h5>"
			tocSuffix ="</a></h5>"
			decode = true

        case "HEADING_6":
			dObj.h6Count++
			parIdStr = fmt.Sprintf("%s_h6_%d",dObj.docName, dObj.h6Count)
			prefix = fmt.Sprintf("<h6 id=\"%s\" class=\"%s_h6\">",parIdStr, dObj.docName)
			tocPrefix = fmt.Sprintf("<h6><a href = \"#%s\">",parIdStr)
			suffix ="</h6>"
			tocSuffix ="</a></h6>"
			decode = true

		case "NORMAL_TEXT":
            prefix = fmt.Sprintf("<p class=\"%s_p\">", dObj.docName)
			suffix ="</p>"
			if par.ParagraphStyle != nil {
				// first check Style parameters against default
				dObj.parCount++
				parIdStr = fmt.Sprintf("%s_p%d",dObj.docName, dObj.parCount)
				prefix = fmt.Sprintf("<p id=\"%s\" class=\"%s_p\">",parIdStr, dObj.docName)
			}

		default:
		// no prefix
			prefix = ""
			suffix = ""
	}


	// a normal paragraph (not a list paragraph)

	if par.Bullet == nil {
		tStr, err = dObj.cvtParStylCss(par.ParagraphStyle)
		t2Str := tStr
		if err != nil {
			t2Str = fmt.Sprintf("/* error cvtParStyl: %v */\n",err) + tStr
		}
		if len(t2Str)>0 {
			parObj.bodyCss += fmt.Sprintf("#%s {\n",parIdStr) + t2Str + "}\n"
		}
	}

	// paragraph elements
	// test
//	parObj, err := dObj.cvtParEl(par)
	numParEl := len(par.Elements)
    for pEl:=0; pEl< numParEl; pEl++ {
        parEl := par.Elements[pEl]
		elHtmlStr, elCssStr, err := dObj.cvtParEl(parEl)
		if err != nil {
            	parHtmlStr += fmt.Sprintf("<!-- error cvtParEl: %v -->\n",err)
		}
      	parHtmlStr += elHtmlStr
		parCssStr += elCssStr

	} // loop par el
	parObj.bodyCss +=listCssStr + parCssStr
	parObj.bodyHtml += listStr + prefix + parHtmlStr + suffix + listSuffix + "\n"
	if decode {
		parObj.tocHtml += tocPrefix + parHtmlStr + tocSuffix + "\n"
	}
	return parObj, nil
}

func (dObj *GdocHtmlObj) cvtParEl(parEl *docs.ParagraphElement)(htmlStr string, cssStr string, err error) {

		if parEl.InlineObjectElement != nil {
        	imgHtmlStr, imgCssStr, err := dObj.cvtInlineImg(parEl.InlineObjectElement)
        	if err != nil {
            	htmlStr += fmt.Sprintf("<!-- error cvtPelToHtml: %v -->\n",err)
        	}
        	htmlStr += imgHtmlStr
			cssStr += imgCssStr
		}

		if parEl.TextRun != nil {
        	txtHtmlStr, txtCssStr, err := dObj.cvtParElText(parEl)
        	if err != nil {
            	htmlStr += fmt.Sprintf("<!-- error cvtPelToHtml: %v -->\n",err)
        	}
        	htmlStr += txtHtmlStr
			cssStr += txtCssStr
		}

        if parEl.HorizontalRule != nil {

        }

        if parEl.ColumnBreak != nil {

        }

        if parEl.Person != nil {

        }

        if parEl.RichLink != nil {

        }

	return htmlStr, cssStr, nil
}

func (dObj *GdocHtmlObj) creParStylCss(headParStyl, parParStyl *docs.ParagraphStyle)(cssStr string, hdStyl bool, err error) {
	// method that compares heading paragraph style to local paragraph style
	// if there is no local paragraph style, cssStr will be empty, otherwise the method will return a cssStr that will be unique for
	// this paragraph


	if (headParStyl == nil)&&(parParStyl == nil) {
		return "", false, fmt.Errorf("error creParStyl: no par style provided!")
	}
	if headParStyl == nil {
	// use pargraph Par Style

		return cssStr, false, nil
	}
	if parParStyl == nil {
	// use pargraph headPar Style

		return cssStr, true, nil
	}

	// both par styles are not nil
	// need to use Parstyle first and headPar as default
	// or just use Parstyle (need to test!!!!)


	return cssStr, false, nil
}

func (dObj *GdocHtmlObj) cvtParStylCss(parStyl *docs.ParagraphStyle)(cssStr string, err error) {

	cssStr = fmt.Sprintf("  /* Paragraph Style: %s */\n", parStyl.NamedStyleType )

	if parStyl == nil {
		return "", fmt.Errorf("error decode parstyle: -- no Style")
	}
//xxx
	tStr := ""
	switch parStyl.NamedStyleType {
		case "TITLE":
			if !dObj.title.exist {
				tStr =fmt.Sprintf(".%s_title {\n",dObj.docName)
				dObj.title.exist = true
			}
		case "SUBTITLE":
			if !dObj.subtitle.exist {
				tStr =fmt.Sprintf(".%s_subtitle {\n",dObj.docName)
				dObj.subtitle.exist = true
			}
		case "HEADING_1":
			if !dObj.h1.exist {
				tStr =fmt.Sprintf(".%s_h1 {\n",dObj.docName)
				dObj.h1.exist = true
			}
		case "HEADING_2":
			tStr =fmt.Sprintf(".%s_h2 {\n",dObj.docName)
		case "HEADING_3":
			tStr =fmt.Sprintf(".%s_h3 {\n",dObj.docName)
		case "HEADING_4":
			tStr =fmt.Sprintf(".%s_h4 {\n",dObj.docName)
		case "HEADING_5":
			tStr =fmt.Sprintf(".%s_h5 {\n",dObj.docName)
		case "HEADING_6":
			tStr =fmt.Sprintf(".%s_h6 {\n",dObj.docName)
		case "NORMAL_TEXT":

		case "NAMED_STYLE_TYPE_UNSPECIFIED":

		case "":

		default:
			tStr =fmt.Sprintf("/* error cvtParStylCss: unknown named Style: %s */", parStyl.NamedStyleType)
	}
	cssStr = tStr
//	fmt.Printf("heading css str: %s\n", tStr)

	tcssStr :=""
	if len(parStyl.Alignment) > 0 {
		switch parStyl.Alignment {
			case "START":
				tcssStr += "  text-align: left;\n"
			case "CENTER":
				tcssStr += "  text-align: center;\n"
			case "END":
				tcssStr += "  text-align: right;\n"
			case "JUSTIFIED":
				tcssStr += "  text-align: justify;\n"
			default:
				tcssStr += fmt.Sprintf("/* unrecognized Alignment %s */\n", parStyl.Alignment)
		}

	}

	if parStyl.IndentFirstLine != nil {
		mag := parStyl.IndentFirstLine.Magnitude
		if mag > 0.0 {
			tcssStr += fmt.Sprintf("  text-indent: %.2fpt;\n", mag)
		}
	}
	if parStyl.IndentStart != nil {
		mag := parStyl.IndentStart.Magnitude
		if mag > 0.0 {
			tcssStr += fmt.Sprintf("  padding-left: %.2fpt;\n", mag)
		}
	}
	if parStyl.IndentEnd != nil {
		mag := parStyl.IndentEnd.Magnitude
		tcssStr += fmt.Sprintf("  padding-right: %.2fpt;\n", mag)
	}

// need to investigate
	if parStyl.LineSpacing > 0 {
		ls := parStyl.LineSpacing
		tcssStr += fmt.Sprintf("  line-height: %.2f;\n", ls/100.0)
	}

	if parStyl.SpaceAbove != nil {
		mag := parStyl.SpaceAbove.Magnitude
		tcssStr += fmt.Sprintf("  padding-top: %.2fpt;\n", mag)
	}

	if parStyl.SpaceBelow != nil {
		mag := parStyl.SpaceBelow.Magnitude
		tcssStr += fmt.Sprintf("  padding-bottom: %.2fpt;\n", mag)
	}

//	if len(tcssStr) > 0 {
//		cssStr = cssCoStr + tcssStr
//	}
	return cssStr, nil
}

func (dObj *GdocHtmlObj) cvtTxtStylCss(txtStyl *docs.TextStyle, head bool)(cssStr string, err error) {
	var tcssStr string

	if txtStyl == nil {
		return "", fmt.Errorf("error decode txtstyle: -- no Style")
	}

	if len(txtStyl.BaselineOffset) > 0 {
		switch txtStyl.BaselineOffset {
			case "SUPERSCRIPT":
				tcssStr += "  vertical-align: sub;\n"
			case "SUBSCRIPT":
				tcssStr += "	vertical-align: sup;\n"
			case "NONE":

			default:
				tcssStr += fmt.Sprintf("/* Baseline Offset unknown: %s */\n", txtStyl.BaselineOffset)
		}
	}
	if txtStyl.Bold {
		tcssStr += "  font-weight: bold;\n"
	} else {
		if head {tcssStr += "  font-weight: normal;\n"}
	}
	if txtStyl.Italic { tcssStr += "  font-style: italic;\n"}
	if txtStyl.Underline { tcssStr += "  text-decoration: underline;\n"}
	if txtStyl.Strikethrough { tcssStr += "  text-decoration: line-through;\n"}

	if txtStyl.WeightedFontFamily != nil {
		font := txtStyl.WeightedFontFamily.FontFamily
		weight := txtStyl.WeightedFontFamily.Weight
		tcssStr += fmt.Sprintf("  font-family: %s;\n", font)
		tcssStr += fmt.Sprintf("  font-weight: %d;\n", weight)
	}
	if txtStyl.FontSize != nil {
		mag := txtStyl.FontSize.Magnitude
		tcssStr += fmt.Sprintf("  font-size: %.2fpt;\n", mag)
	}
	if txtStyl.ForegroundColor != nil {
		if txtStyl.ForegroundColor.Color != nil {
			//0 to 1
            tcssStr += "  color: "
            tcssStr += getColor(txtStyl.ForegroundColor.Color)
		}
	}
	if txtStyl.BackgroundColor != nil {
		if txtStyl.BackgroundColor.Color != nil {
            tcssStr += "  background-color: "
            tcssStr += getColor(txtStyl.BackgroundColor.Color)
		}
	}

//	if txtStyl.Link != nil {
	//txtHtmlStr = txtHtmlStr + '<a href = "' + partAtts.LINK_URL + '">'

	if len(tcssStr) > 0 {
		cssStr = tcssStr
	}
	return cssStr, nil
}


func (dObj *GdocHtmlObj) cvtNamedStylCss(NamStylTyp string) (cssStr string, err error) {
	var NamStyl *docs.NamedStyle
	var tStr string

	doc := dObj.doc
	nStyl := doc.NamedStyles
	idxStyl := -1
// find normal style first
	for istyl:=0; istyl<len(nStyl.Styles); istyl++ {
		if nStyl.Styles[istyl].NamedStyleType == NamStylTyp {
			idxStyl = istyl
			NamStyl = nStyl.Styles[istyl]
			break
		}
	}

	if idxStyl < 0 {return "", nil}

	tStr, err = dObj.cvtParStylCss(NamStyl.ParagraphStyle)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtParStyl: %v",err)
	}
	cssStr += tStr
	tStr, err = dObj.cvtTxtStylCss(NamStyl.TextStyle, true)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtTxtStyl: %v",err)
	}

	return cssStr, nil
}

func (dObj *GdocHtmlObj) cvtAllNamedStylCss() (cssStr string, err error) {
	var NamStyl *docs.NamedStyle
	var tStr string

	doc := dObj.doc
	nStyl := doc.NamedStyles
	normStyl := -1
// find normal style first
	for istyl:=0; istyl<len(nStyl.Styles); istyl++ {
		if nStyl.Styles[istyl].NamedStyleType == "NORMAL_TEXT" {
			normStyl = istyl
			NamStyl = nStyl.Styles[istyl]
			break
		}
	}

	if normStyl < 0 {
		return "", fmt.Errorf("error ConvertNStyl -- no NORMAL_TEXT")
	}

	tStr, err = dObj.cvtParStylCss(NamStyl.ParagraphStyle)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtParStyl: %v",err)
	}
	cssStr += tStr
	tStr, err = dObj.cvtTxtStylCss(NamStyl.TextStyle, true)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtTxtStyl: %v",err)
	}

	cssStr += tStr + "}\n"

	for istyl:=0; istyl<len(nStyl.Styles); istyl++ {
		tStr = ""
		NamStyl = nStyl.Styles[istyl]
		decode := false
		switch NamStyl.NamedStyleType {
		case "TITLE":
			tStr =fmt.Sprintf(".%s_title {\n",dObj.docName)
			decode = true
		case "SUBTITLE":
			tStr =fmt.Sprintf(".%s_subtitle {\n",dObj.docName)
			decode = true
		case "HEADING_1":
			tStr =fmt.Sprintf(".%s_h1 {\n",dObj.docName)
			decode = true
		case "HEADING_2":
			tStr =fmt.Sprintf(".%s_h2 {\n",dObj.docName)
			decode = true
		case "HEADING_3":
			tStr =fmt.Sprintf(".%s_h3 {\n",dObj.docName)
			decode = true
		case "HEADING_4":
			tStr =fmt.Sprintf(".%s_h4 {\n",dObj.docName)
			decode = true
		case "HEADING_5":
			tStr =fmt.Sprintf(".%s_h5 {\n",dObj.docName)
			decode = true
		case "HEADING_6":
			tStr =fmt.Sprintf(".%s_h6 {\n",dObj.docName)
			decode = true

		case "NORMAL_TEXT":

		default:
			tStr =fmt.Sprintf("/* error - header: %s */", NamStyl.NamedStyleType)
		}

		if decode {
			tpStr, err := dObj.cvtParStylCss(NamStyl.ParagraphStyle)
			if err != nil {
				return cssStr, fmt.Errorf("error cvtParStyl: %v",err)
			}
			ttStr, err := dObj.cvtTxtStylCss(NamStyl.TextStyle, true)
			if err != nil {
				return cssStr, fmt.Errorf("error cvtTxtStyl: %v",err)
			}

			cssStr += tStr + tpStr + ttStr + "}\n"
		} else {
			cssStr += tStr
		}
	}
	return cssStr, nil
}

func (dObj *GdocHtmlObj) creHeadCss() (cssStr string, err error) {

	//gdoc division css

	cssStr = fmt.Sprintf(".%s_main {\n", dObj.docName)

    docstyl := dObj.doc.DocumentStyle
	if dObj.Options.Toc {
		cssStr += "  margin-top: 0mm;\n"
		cssStr += "  margin-bottom: 0mm;\n"
	} else {
		cssStr += fmt.Sprintf("  margin-top: %.1fmm; \n",docstyl.MarginTop.Magnitude*PtTomm)
		cssStr += fmt.Sprintf("  margin-bottom: %.1fmm; \n",docstyl.MarginBottom.Magnitude*PtTomm)
	}
    cssStr += fmt.Sprintf("  margin-right: %.2fmm; \n",docstyl.MarginRight.Magnitude*PtTomm)
    cssStr += fmt.Sprintf("  margin-left: %.2fmm; \n",docstyl.MarginLeft.Magnitude*PtTomm)

	dObj.docWidth = (docstyl.PageSize.Width.Magnitude - docstyl.MarginRight.Magnitude - docstyl.MarginLeft.Magnitude)*PtTomm

	cssStr += fmt.Sprintf("  width: %.1fmm;\n", dObj.docWidth)
//xx
	if dObj.Options.DivBorders {
		cssStr += "  border: solid red;\n"
		cssStr += "  border-width: 1px;\n"
	}
	cssStr += "}\n"


//	cssStr += classStr + " > * {\n"
//	cssStr += "  margin: 0;\n}\n"

	listCssStr :=""
	listCssStr += fmt.Sprintf(".%s_li {\n", dObj.docName)
	listCssStr += "  display: inline;\n}\n"

	listCssStr += fmt.Sprintf(".%s_ul {\n", dObj.docName)
	listCssStr += "  margin: 0;\n  padding-left: 0;\n}\n"

	listCssStr += fmt.Sprintf(".%s_ol {\n", dObj.docName)
	listCssStr += "  margin: 0;\n  padding-left: 0;\n}\n"

	listCssStr += fmt.Sprintf(".%s_li_p {\n", dObj.docName)
	listCssStr += "  display: inline;\n  margin: 5pt;\n}\n"

	if dObj.hasList {cssStr += listCssStr}

	// div default
	cssStr += fmt.Sprintf(".%s_normal {\n", dObj.docName)

	defTxtMap := new(textMap)
//fmt.Println("getting txtstyl!")
	parStyl, txtStyl, err := dObj.getNamedStyl("NORMAL_TEXT")
	if err != nil {
		return cssStr, fmt.Errorf("error creHeadCss: %v", err)
	}
//fmt.Printf("txtstyl: %v\n", txtStyl)
	_, err = fillTextMap(defTxtMap, txtStyl)
	if err != nil {
//fmt.Printf("error %v\n", err)
		return cssStr, fmt.Errorf("error creHeadCss: %v", err)
	}
//fmt.Printf("txtmap: %v\n", defTxtMap)

	txtCssStr := cvtTextMapCss(defTxtMap)
//	fmt.Printf(" *** def css ***\n %s\n",tcssStr)
	cssStr += fmt.Sprintf("%s\n}\n", txtCssStr)

	// paragraph default
	cssStr += fmt.Sprintf(".%s_p {\n", dObj.docName)
	cssStr += "  margin: 0; display: block;\n"

	defParMap := new(parMap)

	_, err = fillParMap(defParMap, parStyl)
	if err != nil {
//fmt.Printf("error %v\n", err)
		return cssStr, fmt.Errorf("error creHeadCss: %v", err)
	}
//fmt.Printf("txtmap: %v\n", defTxtMap)

	cssStr += cvtParMapCss(defParMap)

	cssStr += txtCssStr

	cssStr += "}\n"

	if dObj.Options.Toc {
		cssStr += fmt.Sprintf(".%s_toc {\n", dObj.docName)
		cssStr += fmt.Sprintf("  margin-top: %.1fmm;\n", docstyl.MarginTop.Magnitude*PtTomm)
		cssStr += fmt.Sprintf("  margin-bottom: %.1fmm;\n}\n", docstyl.MarginBottom.Magnitude*PtTomm)
	}

	return cssStr, nil
}

func (dObj *GdocHtmlObj) cvtContentEl(contEl *docs.StructuralElement) (GdocHtmlObj *dispObj, err error) {
	if dObj == nil {
		return nil, fmt.Errorf("error cvtContentEl: -- dObj is nil")
	}

	bodyElObj := new(dispObj)

	if contEl.Paragraph != nil {
		parEl := contEl.Paragraph
		tObj, err := dObj.cvtPar(parEl)
		if err != nil { bodyElObj.bodyHtml += fmt.Sprintf("<!-- %v -->\n", err) }
		addDispObj(bodyElObj, &tObj)
	}

	if contEl.SectionBreak != nil {

	}
	if contEl.Table != nil {
		tableEl := contEl.Table
		tObj, err := dObj.cvtTable(tableEl)
		if err != nil { bodyElObj.bodyHtml += fmt.Sprintf("<!-- %v -->\n", err) }
		addDispObj(bodyElObj, &tObj)
	}
	if contEl.TableOfContents != nil {

	}
//	fmt.Println(" ConvertEl: ",htmlObj)
	return bodyElObj, nil
}

func (dObj *GdocHtmlObj) cvtTocHeadCss() (CssStr string, err error) {
	var cssStr, tStr string
	var NamStyl *docs.NamedStyle

	if dObj == nil {
		return "", fmt.Errorf("/* error convertTocHeadtoCss -- dObj is nil */")
	}

	tocCssId := fmt.Sprintf("%s_TOC", dObj.docName)

	doc := dObj.doc
    docstyl := doc.DocumentStyle
	nStyl := doc.NamedStyles

	cssStr = "." + tocCssId + " {\n"
	cssStr += "  margin-top: 10mm;\n"
	cssStr += "  margin-bottom: 10mm;\n"
    cssStr += fmt.Sprintf("  margin-right: %.2fmm; \n",docstyl.MarginRight.Magnitude*PtTomm)
    cssStr += fmt.Sprintf("  margin-left: %.2fmm; \n",docstyl.MarginLeft.Magnitude*PtTomm)

	dObj.docWidth = (docstyl.PageSize.Width.Magnitude - docstyl.MarginRight.Magnitude - docstyl.MarginLeft.Magnitude)*PtTomm

	cssStr += fmt.Sprintf("  width: %.1fmm;\n", dObj.docWidth)
//xxx
	if dObj.Options.DivBorders {
		cssStr += "  border: solid green;\n"
		cssStr += "  border-width: 1px;\n"
	}
	cssStr += "  padding-top:10px;\n  padding-bottom:10px;\n"
	normStyl := -1
// find normal first

	for istyl:=0; istyl<len(nStyl.Styles); istyl++ {
		if nStyl.Styles[istyl].NamedStyleType == "NORMAL_TEXT" {
			normStyl = istyl
			NamStyl = nStyl.Styles[istyl]
			break
		}
	}

	if normStyl < 0 {
		return "", fmt.Errorf("error ConvertNStyl -- no NORMAL_TEXT")
	}

	tStr, err = dObj.cvtParStylCss(NamStyl.ParagraphStyle)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtParStyl: %v",err)
	}
	cssStr += tStr
	tStr, err = dObj.cvtTxtStylCss(NamStyl.TextStyle, true)
	if err != nil {
		return cssStr, fmt.Errorf("error cvtTxtStyl: %v",err)
	}

	cssStr += tStr + "}\n"

	for istyl:=0; istyl<len(nStyl.Styles); istyl++ {
		tStr = ""
		NamStyl := nStyl.Styles[istyl]
		switch NamStyl.NamedStyleType {
		case "TITLE":
			tStr =fmt.Sprintf("#%s_title {\n",tocCssId)
		case "SUBTITLE":
			tStr =fmt.Sprintf("#%s_subtitle {\n",tocCssId)
		case "HEADING_1":
			tStr =fmt.Sprintf(".%s h1 {\n",tocCssId)
 			tStr += "  padding-left: 10px;\n  margin: 0px;"
		case "HEADING_2":
			tStr =fmt.Sprintf(".%s h2 {\n",tocCssId)
			tStr += " padding-left: 20px;\n  margin: 0px;"
		case "HEADING_3":
			tStr =fmt.Sprintf(".%s h3 {\n", tocCssId)
			tStr += " padding-left: 40px;\n  margin: 0px;"
		case "HEADING_4":
			tStr =fmt.Sprintf(".%s h4 {\n", tocCssId)
			tStr += " padding-left: 60px;\n  margin: 0px;"
		case "HEADING_5":
			tStr =fmt.Sprintf(".%s h5 {\n", tocCssId)
			tStr += " padding-left: 80px;\n  margin: 0px;"
		case "HEADING_6":
			tStr =fmt.Sprintf(".%s h6 {\n", tocCssId)
			tStr += " padding-left: 100px;\n  margin: 0px;"
		case "NORMAL_TEXT":

		default:
			tStr =fmt.Sprintf("/* error - header: %s */", NamStyl.NamedStyleType)
		}


		cssStr += tStr
	}
	return cssStr, nil
}


func (dObj *GdocHtmlObj) creTocHead() (hdObj *dispObj, err error) {
	if dObj == nil {
		return nil, fmt.Errorf("error creTocHead -- no GdocObj!")
	}

	hdObj = new(dispObj)
	hdObj.tocHtml = fmt.Sprintf("<div id=\"%s\"_TOC class=\"%s\">\n", dObj.docName, dObj.docName)
	hdObj.tocHtml += fmt.Sprintf("<p id=\"%s_TOC_subtitle\">Table of Contents</p>\n",dObj.docName)
	hdObj.tocCss, err = dObj.cvtTocHeadCss()
	if err != nil {
		return hdObj, fmt.Errorf("error creTocHead:: cvtTocHeadCss: %v", err)
	}
	return hdObj, nil
}

func (dObj *GdocHtmlObj) cvtBody() (bodyObj *dispObj, err error) {

	if dObj == nil {
		return nil, fmt.Errorf("error cvtBody -- no GdocObj!")
	}

	doc := dObj.doc
	body := doc.Body
	if body == nil {
		return nil, fmt.Errorf("error cvtBody -- no body!")
	}

//	toc := dObj.Options.Toc
	bodyObj = new(dispObj)

	bodyObj.bodyHtml = fmt.Sprintf("<div id=\"%s\" class=\"%s_main\">\n",dObj.docName, dObj.docName)

	elNum := len(body.Content)
	for el:=0; el< elNum; el++ {
		bodyEl := body.Content[el]
		tObj, err:=dObj.cvtContentEl(bodyEl)
		if err != nil {
			fmt.Println("error cvtContentEl: %v", err)
		}
//		fmt.Println("tObj:", tObj)
		addDispObj(bodyObj, tObj)
	} // for el loop end
	if len(dObj.cListId) > 0 {
		bodyObj.bodyHtml += dObj.closeList(dObj.cNestLev)
	}

	bodyObj.tocHtml += "</div>\n\n"
	bodyObj.bodyHtml += "</div>\n\n"

	return bodyObj, nil
}

func CreGdocHtmlFil(outfil *os.File, doc *docs.Document, options *OptObj)(err error) {
	var tocDiv *dispObj

	if outfil == nil {return fmt.Errorf("error CreGdocHtmlFil -- outfil is nil!")}

	dObj := new(GdocHtmlObj)
	dObj.folder = outfil
	err = dObj.InitGdocHtmlLib(doc, options)
	if err != nil {
		return fmt.Errorf("error CvtGdocHtml:: InitGdocHtml %v", err)
	}

	toc := dObj.Options.Toc

	mainDiv, err := dObj.cvtBody()
	if err != nil {
		return fmt.Errorf("error cc body %v", err)
	}

	headCssStr, err := dObj.creHeadCss()
	if err != nil {
		return fmt.Errorf("error CvtGdocHtml: ConvertDocHeatAttToCss %v", err)
	}
//	mainDiv.bodyCss = headCssStr + mainDiv.bodyCss

	if toc {
		tocDiv, err = dObj.creTocHead()
		if err != nil {
			tocDiv.tocHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
		}
	}

	// create html file
	docHeadStr,_ := dObj.creHtmlHead()
	outfil.WriteString(docHeadStr)
	// basic Css
	outfil.WriteString(headCssStr)
	// named styles
	outfil.WriteString(mainDiv.bodyCss)
	if toc {
		outfil.WriteString(tocDiv.tocCss)
		outfil.WriteString(mainDiv.tocCss)
	}
	outfil.WriteString("</style>\n</head>\n<body>\n")
	if toc {outfil.WriteString(tocDiv.tocHtml)}
	outfil.WriteString(mainDiv.bodyHtml)

	outfil.WriteString("</body>\n</html>\n")
	return nil
}

