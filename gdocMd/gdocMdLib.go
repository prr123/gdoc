// gdocMdLib.go
// library to convert gdoc document to markdown document
//
// author: prr
// created 1/10/2022
// copyright 2022 prr, azul software
//
//

package gdocToMd

import (
	"fmt"
	"os"
	"strings"
	"google.golang.org/api/docs/v1"
    util "google/gdoc/gdocUtil"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdocMdObj struct {
	elState int
    parCount int
    posImgCount int
    inImgCount int
    imgId []string
    doc *docs.Document
    DocName string
	listmap  map[string][]int
	listid	string
	nestlev int
	folderPath string
	outfil *os.File
	imgFoldNam string
//	imgFoldPath string
	docLists []docList
    headings []heading
    sections []sect
    docFtnotes []docFtnote
    headCount int
    secCount int
    elCount int
    ftnoteCount int
    Options *util.OptObj
}

type sect struct {
    sNum int
    secElStart int
    secElEnd int
}

type heading struct {
    hdElEnd int
    hdElStart int
    id string
    text string
}

type docFtnote struct {
    el int
    parel int
    id string
    numStr string
}

type docList struct {
    listId string
    maxNestLev int64
    ord bool
}

const (
	typPar = iota
	typEm
	typLi
	typTbl
	typImg
	typHd
	typSp
)

func findDocList(list []docList, listid string) (res int) {

    res = -1
    for i:=0; i< len(list); i++ {
        if list[i].listId == listid {
            return i
        }
    }
    return res
}


func (dObj *gdocMdObj) InitGdocMd(folderPath string, options *util.OptObj) (err error) {
    var listItem docList
    var heading heading
    var sec sect
    var ftnote docFtnote

    if dObj == nil {
        return fmt.Errorf("error gdocMD::Init: dObj is nil!")
    }
	doc := dObj.doc
	dObj.inImgCount = len(doc.InlineObjects)
	dObj.posImgCount = len(doc.PositionedObjects)

    dObj.parCount = len(doc.Body.Content)

	dNam := doc.Title
    x := []byte(dNam)
    for i:=0; i<len(x); i++ {
        if x[i] == ' ' {
            x[i] = '_'
        }
    }
    dObj.DocName = string(x[:])

    if options == nil {
        defOpt := new(util.OptObj)
        util.GetDefOption(defOpt)
        if defOpt.Verb {util.PrintOptions(defOpt)}
        dObj.Options = defOpt
    } else {
        dObj.Options = options
    }

//    fPath, fexist, err := util.CreateFileFolder(folderPath, dObj.DocName)
    fPath, fexist, err := util.CreateFileFolder(folderPath, dObj.DocName)
    if err!= nil {
        return fmt.Errorf("error -- util.CreateFileFolder: %v", err)
    }
    dObj.folderPath = fPath

    // create output file path/outfilNam.txt
    outfil, err := util.CreateOutFil(fPath, dObj.DocName,"md")
    if err!= nil {
        return fmt.Errorf("error -- util.CreateOutFil: %v", err)
    }
    dObj.outfil = outfil

    if dObj.Options.Verb {
        fmt.Println("******************* Output File ************")
        fmt.Printf("folder path: %s ", fPath)
        fstr := "is new!"
        if fexist { fstr = "exists!" }
        fmt.Printf("%s\n", fstr)
        fmt.Println("********************************************")
    }


	totObjNum := dObj.inImgCount + dObj.posImgCount
//	if totObjNum == 0 {return nil}


	if dObj.Options.CreImgFolder && (totObjNum > 0) {
		imgFoldPath, err := util.CreateImgFolder(fPath ,dObj.DocName)
		if err != nil {
			return fmt.Errorf("error -- CreateImgFolder: could create ImgFolder: %v!", err)
		}
		dObj.imgFoldNam = imgFoldPath
		err = util.DownloadImages(doc, imgFoldPath, dObj.Options)
		if err != nil {
			return fmt.Errorf("error -- downloadImages could download images: %v!", err)
		}
	}

    dObj.elCount = len(doc.Body.Content)
    // footnotes
    dObj.ftnoteCount = 0

    // section breaks
    parHdEnd := 0
    // last element of section
    secPtEnd := 0
    // set up first page
    sec.secElStart = 0
    dObj.sections = append(dObj.sections, sec)
    seclen := len(dObj.sections)
//      fmt.Println("el: ", el, "section len: ", seclen, " secPtEnd: ", secPtEnd)

    for el:=0; el<dObj.elCount; el++ {
        elObj:= doc.Body.Content[el]

		// Section
        if elObj.SectionBreak != nil {
            if elObj.SectionBreak.SectionStyle.SectionType == "NEXT_PAGE" {
                sec.secElStart = el
                dObj.sections = append(dObj.sections, sec)
                seclen := len(dObj.sections)
                if seclen > 1 {
                    dObj.sections[seclen-2].secElEnd = secPtEnd
                }
            }
        }

		// paragraph, Lists
		if elObj.Paragraph != nil {

			// list
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

            // heading
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
//  fmt.Printf(" text: %s %d\n", text, txtlen)
                heading.text = text

                dObj.headings = append(dObj.headings, heading)
                hdlen := len(dObj.headings)
//      fmt.Println("el: ", el, "hdlen: ", hdlen, "parHdEnd: ", parHdEnd)
                if hdlen > 1 {
                    dObj.headings[hdlen-2].hdElEnd = parHdEnd
                }
            } // end headings

            // footnotes
            if len(doc.Footnotes)> 0 {
                for parEl:=0; parEl<len(elObj.Paragraph.Elements); parEl++ {
                    parElObj := elObj.Paragraph.Elements[parEl]
                    if parElObj.FootnoteReference != nil {
                        ftnote.el = el
                        ftnote.parel = parEl
                        ftnote.id = parElObj.FootnoteReference.FootnoteId
                        ftnote.numStr = parElObj.FootnoteReference.FootnoteNumber
                        dObj.docFtnotes = append(dObj.docFtnotes, ftnote)
                    }
                }
            }

            parHdEnd = el
            secPtEnd = el
        } // end paragraph

		// table
		if elObj.Table != nil {

		}

		if elObj.TableOfContents != nil {

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
        fmt.Printf("********** Headings in Document: %2d ***********\n", len(dObj.headings))
        for i:=0; i< len(dObj.headings); i++ {
            fmt.Printf("  heading %3d  Id: %-15s text: %-20s El Start:%3d End:%3d\n", i, dObj.headings[i].id, dObj.headings[i].text, 
				dObj.headings[i].hdElStart, dObj.headings[i].hdElEnd)
		}
        fmt.Printf("\n********** Pages in Document: %2d ***********\n", len(dObj.sections))
        for i:=0; i< len(dObj.sections); i++ {
            fmt.Printf("  Page %3d  El Start:%3d End:%3d\n", i, dObj.sections[i].secElStart, dObj.sections[i].secElEnd)
        }
        fmt.Printf("\n************ Lists in Document: %2d ***********\n", len(dObj.docLists))
        for i:=0; i< len(dObj.docLists); i++ {
            fmt.Printf("list %3d id: %s max level: %d ordered: %t\n", i, dObj.docLists[i].listId, 
			dObj.docLists[i].maxNestLev, dObj.docLists[i].ord)
		}
        fmt.Printf("\n************ Footnotes in Document: %2d ***********\n", len(dObj.docFtnotes))
        for i:=0; i< len(dObj.docFtnotes); i++ {
            ftn := dObj.docFtnotes[i]
            fmt.Printf("ft %3d: Number: %-4s id: %-15s el: %3d parel: %3d\n", i, ftn.numStr, ftn.id, ftn.el, ftn.parel)
        }

        fmt.Printf("**************************************************\n\n")
    }


    return nil
}


func (dObj *gdocMdObj) cvtTocName(parstr string)(tocstr string) {
	var y [400]byte
	parlStr := strings.ToLower(parstr)
	parlStr = strings.TrimSpace(parlStr)
	if len(parstr) < 1 {
		return ""
	}
	//fmt.Println("toc parstr: ", parstr," : ", len(parstr))
	//fmt.Println("toc parlStr: ", parlStr, " : ", len(parlStr))

	x:= []byte(parlStr)

	xSt :=-1
	for i:=0; i<len(x); i++ {
		switch x[i] {
			case ' ','*','_':

			default:
				xSt = i
		}
		if xSt > -1 {break}
	}
	if xSt < 0 {
		return ""
	}

	yPos := -1
	xEnd := -1
	last := false
	if len(x) < 20 {
//fmt.Printf("[%d:%d]:%s\n",xSt,len(x),x[xSt:len(x)])
	}
	for i:=xSt; i<len(x); i++ {
		switch x[i] {
		case ' ':
			if !last {
				yPos++
				y[yPos] = '-'
			}
			last = true
		case '*','_':
			xEnd = i
		default:
			yPos++
			y[yPos] = x[i]
			last = false
		}
		if xEnd > -1 {break}
	}

	tocstr = string(y[:yPos+1])+ "-1"

	//fmt.Println("toctr: ", tocstr, " : ", len(tocstr))
	//fmt.Println()

	return tocstr
}

func (dObj *gdocMdObj) cvtMdHeadStyl()(outstr string, err error) {
	outstr = ""
	if dObj == nil {
		return outstr, fmt.Errorf("error cvtMdHeadStyl -- no dObj!")
	}
	doc := dObj.doc
// Inline Objects
// Lists
	inObjLen := len(doc.InlineObjects)
	outstr += fmt.Sprintf("\nInline Objects: %d\n", inObjLen)

	posObjLen := len(doc.PositionedObjects)
	outstr += fmt.Sprintf("\nPositioned Objects: %d\n",posObjLen)

	headLen := len(doc.Headers)
	outstr += fmt.Sprintf("\nHeaders: %d\n",headLen)

	footLen := len(doc.Footers)
	outstr += fmt.Sprintf("\nFooters: %d\n",footLen)

	ftnoteLen := len(doc.Footnotes)
	outstr += fmt.Sprintf("\nFootnotes: %d\n",ftnoteLen)

	listLen := len(doc.Lists)
	outstr += fmt.Sprintf("\nLists: %d\n",listLen)

	nrLen := len(doc.NamedRanges)
	outstr += fmt.Sprintf("\nNamedRanges: %d\n",nrLen)

    return outstr, nil
}

func (dObj *gdocMdObj) renderPosImg(posObj *docs.PositionedObject)(outstr string, err error) {

	if posObj == nil {
        return "", fmt.Errorf("posImg is nil!")
	}

	imgElId := posObj.ObjectId

    idx := -1
    for i:=0; i< len(imgElId); i++ {
        if imgElId[i] == '.' {
            idx = i+1
            break
        }
    }

	if idx < 0 {return "", fmt.Errorf("not a valid id: %s!", imgElId)}
    imgId :=""
    if (idx>0) && (idx<len(imgElId)-1) {
        imgId = "img_" + imgElId[idx:]
    }

    imgObj := posObj.PositionedObjectProperties.EmbeddedObject
	imgSrc := dObj.imgFoldNam + "/" + imgId + ".jpeg"

//    	outstr=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\" width=\"%.1fpt\" height=\"%.1fpt\">\n",
//			imgSrc, imgId, imgObj.Title, imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )

	outstr = "![" + imgObj.Description + " " + imgObj.Title + "](" + imgSrc + ")"


	return outstr, nil
}


func (dObj *gdocMdObj) renderInlineImg(imgEl *docs.InlineObjectElement)(outstr string, err error) {

   	if imgEl == nil {
        return "", fmt.Errorf("error convertPelInlineImg:: imgEl is nil!")
    }

    doc := dObj.doc

    imgElId := imgEl.InlineObjectId
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

    imgObj := doc.InlineObjects[imgElId].InlineObjectProperties.EmbeddedObject
	imgSrc := dObj.imgFoldNam + "/" + imgId + ".jpeg"
//    outstr=fmt.Sprintf("<img src=\"%s\" alt=\"%s\" width=\"%.1fpt\" height=\"%.1fpt\">\n",
//    outstr=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\" width=\"%.1fpt\" height=\"%.1fpt\">\n",
//		imgSrc, imgId, imgObj.Title, imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )

		outstr = "![" + imgObj.Description + " " + imgObj.Title + "](" + imgSrc + ")"

    return outstr, nil
}

func (dObj *gdocMdObj) cvtPelTxt(parEl *docs.ParagraphElement, bold bool, italic bool)(outstr string, err error) {

    if parEl == nil {
        return "", fmt.Errorf("error cvtPelToMd -- parEl is nil!")
    }

//fmt.Println("maxChar: ", maChar)

	xnbLen := len(parEl.TextRun.Content)
	xnb := []byte(parEl.TextRun.Content)

// should be either <br> or nothing
	if (xnbLen == 1) && (xnb[0] == '\n') {
		return "", nil
	}

	if xnb[xnbLen-1] == '\n' {xnbLen --}

	// left wsp
	xSt := -1
	for i:=0; i< xnbLen; i++ {
		if xnb[i] == ' ' {
			continue
		} else {
			xSt = i
			break
		}
	}

	xEnd := xnbLen
	for i:= xnbLen - 1; i> -1; i-- {
		if xnb[i] == ' ' {
			continue
		} else {
			xEnd = i +1
			break
		}
	}

//fmt.Printf("xnbLen: %d [%d:%d]: %s\n", xnbLen, xSt, xEnd, xnb[xSt:xEnd])

    prefix := ""
    suffix := ""

    if italic {
        prefix = "_"
        suffix  ="_"
    }

    if bold {
        prefix = "**" + prefix
        suffix = suffix + "**"
    }

	 outstr = string(xnb[0:xSt]) + prefix + string(xnb[xSt:xEnd]) + suffix + string(xnb[xEnd:xnbLen])

    return outstr, nil
}


func (dObj *gdocMdObj) cvtParToMd(par *docs.Paragraph)(outstr string, tocstr string) {
	var prefix, suffix, tocPrefix, tocSuffix string
	var parStr, listStr, hdStr, errStr string
	var NamedTxtStyl *docs.TextStyle
	var doc *docs.Document

	doc = dObj.doc
//	opt := dObj.Options

	numParEl := len(par.Elements)
	if (numParEl == 1) && (!(len(par.Elements[0].TextRun.Content) >0)) {
		// html		return "<br>\n","", nil
		return "\n",""
	}

//	fmt.Printf("  Paragraph with %d Par-Elements\n", len(par.Elements))
	outstr = ""

	// lists
	if par.Bullet != nil {
		listStr = ""
		listid := par.Bullet.ListId
		nestIdx:= int(par.Bullet.NestingLevel)

//        if dObj.Options.Verb {hdStr = fmt.Sprintf("<!--- List Element %d --->", dObj.parCount)}

        // retrieve the list properties from the doc.Lists map
        nestL := dObj.doc.Lists[listid].ListProperties.NestingLevels[nestIdx]
        listOrd := util.GetGlyphOrd(nestL)

        // list item
		// we can do this differently start with an empty byte string
		listPrefix := ""
		for nl:=0; nl < nestIdx; nl++ {
			listPrefix += "  "
		}

		if listOrd {
	        listStr = listPrefix + "1. "
		} else {
			listStr = listPrefix + "*  "
		}
    } // end of lists


	// check for heading

// heading id
//	hdId := par.ParagraphStyle.HeadingId

// indentation
	parIndent := 0.0
	if(par.ParagraphStyle.IndentStart) != nil {
		parIndent = par.ParagraphStyle.IndentStart.Magnitude
	}

// temp solution
//	fmt.Printf("par indent: %.0f\n", parIndent)
	if parIndent > 0 {

	}

	parStylTyp := par.ParagraphStyle.NamedStyleType
	// doc styles
	for i:=0; i < len(doc.NamedStyles.Styles); i++ {
		nstyl := doc.NamedStyles.Styles[i].NamedStyleType;
		if nstyl  == parStylTyp {
//			NamedParStyl = doc.NamedStyles.Styles[i].ParagraphStyle
			NamedTxtStyl = doc.NamedStyles.Styles[i].TextStyle
		}
	}


//	fmt.Println("par style type: ", parStylTyp, " index: ", NamedStylIdx);

    prefix = ""
	suffix = ""
	tocPrefix = ""
	tocSuffix = ""

	decode := false
	titlestyl := false
	subtitlestyl := false

	// for title and subtitle only
	boldStyl := false
	italicStyl := false
	if par.Elements[0].TextRun != nil {
		parTxtStyl := par.Elements[0].TextRun.TextStyle
		boldStyl = parTxtStyl.Bold || NamedTxtStyl.Bold
		italicStyl = parTxtStyl.Italic || NamedTxtStyl.Italic
	}

//	headId := ""
//	if len(par.ParagraphStyle.HeadingId) > 0 {headId = string(par.ParagraphStyle.HeadingId[2:])}

    switch parStylTyp {
		case "TITLE":
			prefix ="<p style=\"font-size:20pt; text-align:center;"
			suffix = "</p>\n\n"
			if boldStyl {prefix += " font-weight: 800;"}
			if italicStyl {prefix += " font-style:italic;"}
			prefix += "\">"
			titlestyl = true
			tocPrefix = prefix
			tocSuffix = suffix

      	case "SUBTITLE":
			prefix ="<p style=\"font-size:16pt;text-align:center;"
			suffix = "</p>\n\n"
			if boldStyl {prefix += " font-weight:bold;"}
			if italicStyl {prefix += " font-style:italic;"}
			prefix += "\">"
			subtitlestyl = true

		case "HEADING_1":
            prefix = "# "
//			suffix = fmt.Sprintf("    \n")
			tocPrefix = "# ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

		case "HEADING_2":
            prefix = "## "
//			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n## ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

        case "HEADING_3":
			prefix = "### "
//			suffix = fmt.Sprintf("     \n")
			tocPrefix = "\n### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

        case "HEADING_4":
            prefix = "#### "
//			suffix = fmt.Sprintf("     \n")
			tocPrefix = "\n#### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

        case "HEADING_5":
            prefix = "##### "
//			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n##### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

        case "HEADING_6":
            prefix = "###### "
//			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n###### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true

        case "NORMAL_TEXT":
            prefix = ""
//			suffix = ""
			if par.Bullet != nil {suffix = "\n"}

        default:
            prefix = fmt.Sprintf("[//]: * (Name Style: %s unknown)\n", par.ParagraphStyle.NamedStyleType)
    }

	// need to reconsider for paragraphs without text

	maxChar := dObj.Options.MaxCharLine
	parStr = ""
	parTxtStr := ""
	parErrStr := ""
	for p:=0; p< numParEl; p++ {

		parEl := par.Elements[p]
//fmt.Printf("parEl %2d: %s\n", p, parEl.TextRun.Content)
//		outstr += fmt.Sprintf("\nPar-El[%d]: %d - %d \n", p, parEl.StartIndex, parEl.EndIndex)

		boldStyl := false
		italicStyl := false
		if (par.Elements[p].TextRun != nil)&& !(titlestyl || subtitlestyl) {
			parTxtStyl := par.Elements[p].TextRun.TextStyle
			boldStyl = parTxtStyl.Bold || NamedTxtStyl.Bold
			italicStyl = parTxtStyl.Italic || NamedTxtStyl.Italic
		}


		// text
		if parEl.TextRun != nil {
			tstr, err := dObj.cvtPelTxt(parEl, boldStyl, italicStyl)
			if err != nil {
				errStr = fmt.Sprintf("\n[//]: # (error el %d cvtPelTxt: %v)\n", p, err)
			}
//fmt.Println("tstr: ",tstr)
			parErrStr += errStr
			parTxtStr += tstr
		}


		if parEl.HorizontalRule != nil {
			parStr += "---\n"
		}
		if parEl.ColumnBreak != nil {
			errStr = fmt.Sprintf("\n[//]: # (error el %d: Column Break not implemented)\n", p)
			parStr += errStr
		}
   		if parEl.Person != nil {
			errStr = fmt.Sprintf("\n[//]: # (error el %d: Person not implemented)\n", p)
			parStr += errStr
		}

		// inline image
		if parEl.InlineObjectElement != nil {
            tstr, err := dObj.renderInlineImg(parEl.InlineObjectElement)
            if err != nil {
				errStr = fmt.Sprintf("\n[//]: # (error renderInlineImg: %v)\n",err)
            }
			parStr += tstr + errStr
		}
		//fmt.Printf("Img parstr %d: %s (%d)\n",p, parStr,len(parStr))

 		if parEl.RichLink != nil {
			errStr = fmt.Sprintf("\n[//]: # (error el %d: RichLink not implemented)\n", p)
			parStr += errStr
		}

	} // loop parEl

	if (len(parTxtStr) == 0) && (par.Bullet == nil) {
		outstr ="\n"
		tocstr = "\n"
		return outstr, tocstr
	}

//fmt.Printf("parstr raw (%d):\n%s\n", len(parTxtStr), parTxtStr)

	xnb := []byte(parTxtStr)
	xnbLen := len(xnb)
	for i := 0; i< xnbLen; i++ {
		if xnb[i] == '\n' {xnb[i] = ' '}
	}

	altLen := xnbLen
	for i:=xnbLen -1; i> xnbLen-5; i-- {
		if xnb[i] == ' ' {altLen --} else {break}
	}

	xnbLen = altLen
//fmt.Printf("parstr2 no cr (%d):\n%s\n", xnbLen, string(xnb))

	linSt := 0
	if maxChar > 10 {
		for lin := 0; lin <100; lin++ {
			linEnd := linSt + maxChar

			if linEnd > xnbLen {
				linEnd = xnbLen
				break
			}


			for i:= linEnd; i> linSt; i-- {
				if xnb[i] == ' ' {
					xnb[i] = '\n'
					linSt = i +1
					break
				}
			}
		}
	}

	parNTxtStr := string(xnb[:xnbLen])
//	parCrTxtStr := string(xnb[:xnbLen]) + "  \n"
	parTitleStr := string(xnb[:xnbLen])

/*
	pL := len(parNTxtStr)
	plStart := pL -10
	if plStart < 0 {plStart = 0}
fmt.Printf("end parCrTxtStr:")
for i := plStart; i< pL; i++ {
	fmt.Printf("%q", parNTxtStr[i])
}
fmt.Println()
*/
	tocOpt := dObj.Options.Toc
	switch {
		case titlestyl:
			if tocOpt {	tocstr+= tocPrefix + parTitleStr + tocSuffix}
	    	outstr = prefix + parTitleStr + suffix
		case subtitlestyl:
	    	outstr = prefix + parTitleStr + suffix
		default:
			suffix = "  \n"
//			if len(headId) > 0 {suffix = " {#" + headId + "}  \n"}
	    	outstr = listStr + prefix + parNTxtStr + suffix + parStr
			if decode && tocOpt {
//				tocstr+= tocPrefix + parNTxtStr + tocSuffix + headId + ")"
				tocParStr := dObj.cvtTocName(parNTxtStr)
				tocstr+= tocPrefix + parNTxtStr + tocSuffix + tocParStr + ")"
			}
		}

/*
	outL := len(outstr)
	stOut:= outL - 10
	if stOut < 0 {stOut = 0}
fmt.Printf("cvt par outstr:")
for i := stOut; i< outL; i++ {
	fmt.Printf("%q", outstr[i])
}
fmt.Println()
*/

	if par.PositionedObjectIds != nil {
		for id := 0; id < len(par.PositionedObjectIds); id++ {
			idStr := par.PositionedObjectIds[id]
			posObj, ok := doc.PositionedObjects[idStr]
			if !ok {
				errStr = fmt.Sprintf("\n[//]: # (error par: RichLink not implemented)\n")
				outstr += errStr
				continue
			}
			tstr, err := dObj.renderPosImg(&posObj)
			if err != nil {
				errStr = fmt.Sprintf("\n[//]: # (error renderPosImg: %v)\n",err)
			}
			outstr += tstr + errStr
		}
	}

	return hdStr + outstr, tocstr
}

func (dObj *gdocMdObj) cvtTableToMd(tbl *docs.Table)(outstr string) {
var errStr string

	outstr = fmt.Sprintf("[//]: * (Table Start)\n")

	cols := int(tbl.Columns)
	rows := int(tbl.Rows)

	// simple paragraphs

	for irow := 0; irow < rows; irow++ {
		trow := tbl.TableRows[irow]
		rowStr := "|"
		for icol := 0; icol < cols; icol++ {
			tcell := trow.TableCells[icol]
			ilen := len(tcell.Content)
			if ilen > 1 {

			}
			par := tcell.Content[0].Paragraph
			if par == nil {

			}
			numParEl := len(par.Elements)
			parStr := ""
			for p:=0; p< numParEl; p++ {
				parEl := par.Elements[p]
//fmt.Println("parEl: ", p, parEl.TextRun.Content)
//		outstr += fmt.Sprintf("\nPar-El[%d]: %d - %d \n", p, parEl.StartIndex, parEl.EndIndex)
				bold := false
				italic := false
				// text
				if parEl.TextRun != nil {
					tstr, err := dObj.cvtPelTxt(parEl, bold, italic)
					if err != nil {
						errStr = fmt.Sprintf("\n[//]: # (error el %d cvtPelTxt: %v)\n", p, err)
					}
					parStr += tstr + errStr
				//fmt.Printf("parstr %d: %s (%d)\n",p, parStr,len(parStr))
				}

				if parEl.HorizontalRule != nil {
					continue
				}
				if parEl.ColumnBreak != nil {
					continue
				}
   				if parEl.Person != nil {
					continue
//				errStr = fmt.Sprintf("\n[//]: # (error el %d: Person not implemented)\n", p)
//				parStr += errStr
				}

				// inline image
				if parEl.InlineObjectElement != nil {
    	        	tstr, err := dObj.renderInlineImg(parEl.InlineObjectElement)
        	    	if err != nil {
						errStr = fmt.Sprintf("\n[//]: # (error renderInlineImg: %v)\n",err)
            		}
					parStr += tstr + errStr
				}
				//fmt.Printf("Img parstr %d: %s (%d)\n",p, parStr,len(parStr))

 				if parEl.RichLink != nil {
					continue
//					errStr = fmt.Sprintf("\n[//]: # (error el %d: RichLink not implemented)\n", p)
//					parStr += errStr
				}
			} // loop parEl
			rowStr += parStr + "|"
		} // icol
		outstr += rowStr + "\n"
	} // irow

	return outstr
}

func CvtGdocToMd(folderPath string, doc *docs.Document, options *util.OptObj)(err error) {
	var outstr, tocstr string

    docObj := new(gdocMdObj)
    docObj.doc = doc

    err = docObj.InitGdocMd(folderPath, options)
    if err != nil {
        return fmt.Errorf("error CvtGdocToMd: could not initialise! %v", err)
    }

	outstr = "[//]: * (Document Title: " + doc.Title + ")\n"

	outstr += fmt.Sprintf("[//]: * (Document Id: %s)\n", doc.DocumentId)
//	outstr += fmt.Sprintf("Revision Id: %s \n", doc.RevisionId)

	outfil := docObj.outfil

	_, err = outfil.WriteString(outstr)
	if err != nil {
		return fmt.Errorf("error CvtGdocTpMd: cannot write to file! %v", err)
	}

	if (docObj.Options.Toc) && (len(docObj.headings)>2) {
		tocstr = "<p style=\"font-size:20pt; text-align:center\">Table of Contents</p>\n"
	}

	body := doc.Body
	elNum := len(body.Content)
//	outstr = "\n******************** Body *********************************\n"

	docObj.parCount = 0;
	bodyStr :=""
	for el:=0; el< elNum; el++ {
		bodyEl := body.Content[el]
//		outstr += fmt.Sprintf("\nelement: %d StartIndex: %d EndIndex: %d\n", el, bodyEl.StartIndex, bodyEl.EndIndex)

		// par, lists, headings
		if bodyEl.Paragraph != nil {
			par := bodyEl.Paragraph
			tstr, toctstr := docObj.cvtParToMd(par)
			bodyStr += tstr
			tocstr += toctstr
/*
ltstr := len(tstr)
fmt.Printf("tstr (%d):\n%s", ltstr, tstr)
ltstart := ltstr -5
if ltstart < 0 {ltstart = 0}
for i:= ltstart; i< ltstr; i++ {
	fmt.Printf("%q",tstr[i])
}
fmt.Println()
*/
		} // end par

		// sections
		if bodyEl.SectionBreak != nil {
			bodyStr += fmt.Sprintf("[//]: * (Section Break)\n")
		}

		// tables
		if bodyEl.Table != nil {
			bodyStr += fmt.Sprintf("[//]: * (Table)\n")
			tbl := bodyEl.Table
			tstr := docObj.cvtTableToMd(tbl)
			bodyStr += tstr
		}

		// table of contents
		if bodyEl.TableOfContents != nil {
			bodyStr += fmt.Sprintf("[//]: * (TableOfContents)\n")
		}

//		fmt.Printf("el %d [%d]:\n%s", el, len(bodyStr), bodyStr)

	}

	if (docObj.Options.Toc) && (len(docObj.headings)>2) {
//	if len(tocstr) > 0  {
		tocstr += "\n\n"
		outfil.WriteString(tocstr)
	}

	outfil.WriteString(bodyStr)

	outfil.Close()
	return nil
}
