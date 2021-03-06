//  gdocTxtLib
//  reads gdoc file and writes text version to outfile
//  author PRR
//  created 1/1/2022
//
// copyright 2022 prr, azul software
//
// license see https://github.com/prr123/gdoc/tree/master/gdocTxt
//

package gdocTxtLib

import (
	"fmt"
	"os"
	"unicode/utf8"
	"google.golang.org/api/docs/v1"
    gdocUtil "google/gdoc/gdocUtil"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdocTxtObj struct {
	parCount int
	cpl int
	posImgCount int
 	inImgCount int
	TableCount int
	imgId []string
	doc *docs.Document
	DocName string
	headingId []string
	folderPath string
	outfil *os.File
    imgFoldNam string
    imgFoldPath string
	DocOpt	bool
    docLists []docList
//    headings []heading
    sections []sect
	pgCounter int
    docFtnotes []docFtnote
	ftnoteCount int
//	namStylMap map[string]bool
	Options *gdocUtil.OptObj
	txtfil *os.File
	listCount [10]int
}

type docList struct {
    listId string
    maxNestLev int64
    ord bool
}

type namStyl struct {
    count int
    exist bool
    tocExist bool
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

func findDocList(list []docList, listid string) (res int) {

    res = -1
    for i:=0; i< len(list); i++ {
        if list[i].listId == listid {
            return i
        }
    }
    return res
}

func (dObj *gdocTxtObj) initGdocTxt(folderPath string, options *gdocUtil.OptObj) (err error) {
    var listItem docList
    var sec sect
    var ftnote docFtnote

	doc := dObj.doc
	if dObj == nil {return fmt.Errorf("dObj is nil in FdocTxtObj!")}

	body := dObj.doc.Body


    dNam := doc.Title
    x := []byte(dNam)
    for i:=0; i<len(x); i++ {
        if x[i] == ' ' {
            x[i] = '_'
        }
    }
    dObj.DocName = string(x[:])

	if options == nil {
		defOpt := new(gdocUtil.OptObj)
		gdocUtil.GetDefOption(defOpt)
		if defOpt.Verb {gdocUtil.PrintOptions(defOpt)}
 		dObj.Options = defOpt
	} else {
		dObj.Options = options
	}

	dObj.parCount = 0
	dObj.posImgCount = 0
	dObj.inImgCount = 0
	dObj.TableCount = 0
	dObj.pgCounter = 0
	dObj.cpl = 70

    // footnotes
    dObj.ftnoteCount = 0

    // section breaks
//    parHdEnd := 0

    // last element of section
    secPtEnd := 0

    // set up first page
    sec.secElStart = 0
    dObj.sections = append(dObj.sections, sec)
    seclen := len(dObj.sections)
	// fmt.Println("el: ", el, "section len: ", seclen, " secPtEnd: ", secPtEnd)

	for el:=0; el< len(body.Content); el++ {
		elObj := body.Content[el]
		if elObj.Paragraph != nil {
            if elObj.Paragraph.Bullet != nil {
                listId := elObj.Paragraph.Bullet.ListId
                found := findDocList(dObj.docLists, listId)
                nestlev := elObj.Paragraph.Bullet.NestingLevel
                if found < 0 {
                    listItem.listId = listId
                    listItem.maxNestLev = elObj.Paragraph.Bullet.NestingLevel
                    nestL := doc.Lists[listId].ListProperties.NestingLevels[nestlev]
                    listItem.ord = gdocUtil.GetGlyphOrd(nestL)
                    dObj.docLists = append(dObj.docLists, listItem)
                } else {
                    if dObj.docLists[found].maxNestLev < nestlev { dObj.docLists[found].maxNestLev = nestlev }
                }

            }

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

//            parHdEnd = el
            secPtEnd = el
        } // end paragraph
    } // end el loop

    seclen = len(dObj.sections)
    if seclen > 0 {
        dObj.sections[seclen-1].secElEnd = secPtEnd
    }

    if dObj.Options.Verb {
        fmt.Printf("\n********** Pages in Document: %2d ***********\n", len(dObj.sections))
        for i:=0; i< len(dObj.sections); i++ {
            fmt.Printf("  Page %3d  El Start:%3d End:%3d\n", i, dObj.sections[i].secElStart, dObj.sections[i].secElEnd)
		}
        fmt.Printf("\n************ Lists in Document: %2d ***********\n", len(dObj.docLists))
        for i:=0; i< len(dObj.docLists); i++ {
            fmt.Printf("list %3d id: %s max level: %d ordered: %t\n", i, dObj.docLists[i].listId, dObj.docLists[i].maxNestLev,
            dObj.docLists[i].ord)
		}
        fmt.Printf("\n************ Footnotes in Document: %2d ***********\n", len(dObj.docFtnotes))
        for i:=0; i< len(dObj.docFtnotes); i++ {
            ftn := dObj.docFtnotes[i]
            fmt.Printf("ft %3d: Number: %-4s id: %-15s el: %3d parel: %3d\n", i, ftn.numStr, ftn.id, ftn.el, ftn.parel)
		}
        fmt.Printf("**************************************************\n\n")
    }

	// images
    dObj.inImgCount = len(doc.InlineObjects)
    dObj.posImgCount = len(doc.PositionedObjects)

	// create folders
    fPath, fexist, err := gdocUtil.CreateFileFolder(folderPath, dObj.DocName)
	//    fPath, _, err := gdocUtil.CreateFileFolder(folderPath, dObj.DocName)
    if err!= nil { return fmt.Errorf("error -- gdocUtil.CreateFileFolder: %v", err)}

    dObj.folderPath = fPath

    if dObj.Options.Verb {
        fmt.Println("******************* Output File ************")
        fmt.Printf("folder path: %s ", fPath)
        fstr := "is new!"
        if fexist { fstr = "exists!" }
        fmt.Printf("%s\n", fstr)
		fmt.Println("********************************************")
    }

    // create output file path/outfilNam.txt
    outfil, err := gdocUtil.CreateOutFil(fPath, dObj.DocName,"txt")
    if err != nil {
        return fmt.Errorf("error -- gdocUtil.CreateOutFil: %v", err)
    }
    dObj.txtfil = outfil

	totObjNum := dObj.inImgCount + dObj.posImgCount

    if dObj.Options.CreImgFolder && (totObjNum > 0) {
        imgFoldPath, err := gdocUtil.CreateImgFolder(fPath ,dObj.DocName)
        if err != nil {
            return fmt.Errorf("error -- CreateImgFolder: could create ImgFolder: %v!", err)
        }
        dObj.imgFoldNam = imgFoldPath
        err = gdocUtil.DownloadImages(doc, imgFoldPath, dObj.Options)
        if err != nil {
            return fmt.Errorf("error -- downloadImages could download images: %v!", err)
        }
    }

	return nil
}


func (dObj *gdocTxtObj) cvtTOC(toc *docs.TableOfContents)(outstr string, err error) {

	if toc == nil {
		return "", fmt.Errorf("error dispTOC: toc is nil")
	}

	numEl := len(toc.Content)
	outstr = fmt.Sprintf("\n  *** Table of Contents: (%d) ***\n", numEl)

	for el:=0; el< numEl; el++ {
//		tocEl := toc.Content[el]
//		outstr += fmt.Sprintf("\n *** element: %d ***", el)
/*
		tStr,err := dObj.dispContentEl(tocEl)
		if err != nil {
			outstr += fmt.Sprintf("\n error dispContentEl: %v\n", err) + tStr
		} else {
			outstr += tStr
		}
*/
	}
	return outstr, nil
}

func (dObj *gdocTxtObj) cvtTable(tbl *docs.Table)(outstr string, err error) {
	var tabWidth, tabHeight float64
	var icol, irow int64

	doc := dObj.doc

	if tbl == nil {
		return "", fmt.Errorf("error dispTable: no table pt")
	}

	dObj.TableCount++
	outstr = fmt.Sprintf("  *** table %d: rows: %d cols: %d ***\n",dObj.TableCount,tbl.Rows, tbl.Columns )
	// table rows

	outstr += "  Table Style Properties\n"
	tabWidth = 0.0
	if tbl.TableStyle != nil {
		numColProp := (int64)(len(tbl.TableStyle.TableColumnProperties))
		for icol=0; icol<numColProp; icol++ {
			tColProp := tbl.TableStyle.TableColumnProperties[icol]
			outstr += fmt.Sprintf("    col[%d]: w type: %s", icol, tColProp.WidthType)
			if tColProp.Width != nil {
				outstr += fmt.Sprintf(" width: %.1fpt", tColProp.Width.Magnitude)
			}
			outstr += "\n"
		}
	}

	docPg := doc.DocumentStyle
	PgWidth := docPg.PageSize.Width.Magnitude
	NetPgWidth := PgWidth - (docPg.MarginLeft.Magnitude + docPg.MarginRight.Magnitude)
	outstr += fmt.Sprintf("    Default Table Width: %.1f", NetPgWidth*PtTomm)
	tabWidth = NetPgWidth
    for icol=0; icol < tbl.Columns; icol++ {
		tcolObj :=tbl.TableStyle.TableColumnProperties[icol]
		if tcolObj.Width != nil {
			tabWidth += tbl.TableStyle.TableColumnProperties[icol].Width.Magnitude
		}
	}

	tabHeight = 0.0
	for irow=0; irow < tbl.Rows; irow++ {
		trowObj := tbl.TableRows[irow]
		tabHeight += trowObj.TableRowStyle.MinRowHeight.Magnitude
	}

	outstr += fmt.Sprintf("  Min Height: %.1fmm Width: %.1fmm\n\n", tabHeight*PtTomm, tabWidth*PtTomm)

	tblCellCount:=0
    for irow =0; irow < tbl.Rows; irow++ {
        trowobj := tbl.TableRows[irow]
        mrheight := trowobj.TableRowStyle.MinRowHeight.Magnitude
        numCols := (int64)(len(trowobj.TableCells))
		outstr += fmt.Sprintf("  table row[%d]: cols:%d min Height: %.1f\n", irow, numCols, mrheight)
		tcellDefWidth := tabWidth/(float64)(numCols)

        for icol =0; icol< numCols; icol++ {
			tcolObj :=tbl.TableStyle.TableColumnProperties[icol]
			tcellWidth := tcellDefWidth
			if tcolObj.Width != nil {
				tcellWidth = tcolObj.Width.Magnitude
			}
			outstr += fmt.Sprintf("    col[%d]: width: %6.1f type: %s", icol, tcellWidth, tcolObj.WidthType)

            tcell := trowobj.TableCells[icol]
			txtStr := ""
			numEl := len(tcell.Content)
			for el:=0; el<numEl; el++ {
				if tcell.Content[el].Paragraph != nil {
					for pel:=0; pel < len(tcell.Content[el].Paragraph.Elements); pel++ {
						pelObj := tcell.Content[el].Paragraph.Elements[pel]
						if pelObj.TextRun != nil {
							txtStr += pelObj.TextRun.Content
						}
					}
				}
			}
            tcellstyl := tcell.TableCellStyle
			if tcellstyl != nil {
				if tcellstyl.BackgroundColor.Color != nil {
					outstr += fmt.Sprintf(" color: %s", gdocUtil.GetColor(tcellstyl.BackgroundColor.Color))
				}
				outstr += fmt.Sprintf(" vert align: %s", tcellstyl.ContentAlignment)
				padTop := 0.0
				if tcellstyl.PaddingTop != nil {
					padTop = tcellstyl.PaddingTop.Magnitude
				}
				padRight := 0.0
				if tcellstyl.PaddingRight != nil {
					padRight = tcellstyl.PaddingRight.Magnitude
				}
				padBottom := 0.0
				if tcellstyl.PaddingBottom != nil {
					padBottom = tcellstyl.PaddingBottom.Magnitude
				}
				padLeft := 0.0
				if tcellstyl.PaddingLeft != nil {
					padLeft = tcellstyl.PaddingLeft.Magnitude
				}
				outstr += fmt.Sprintf(" pad: %.1f %.1f %.1f %.1f", padTop, padRight, padBottom, padLeft)
			}
			outstr += fmt.Sprintf(" text: %s", txtStr)

			if tcellstyl != nil {
				outstr += "     border top: "
				if tcellstyl.BorderTop != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderTop.DashStyle)

					if tcellstyl.BorderTop.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderTop.Width.Magnitude)
					}
					if tcellstyl.BorderTop.Color != nil {
						outstr += fmt.Sprintf(" color: %s", gdocUtil.GetColor(tcellstyl.BorderTop.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " right: "
				if tcellstyl.BorderRight != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderRight.DashStyle)

					if tcellstyl.BorderRight.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderRight.Width.Magnitude)
					}
					if tcellstyl.BorderRight.Color != nil {
						outstr += fmt.Sprintf(" color: %s", gdocUtil.GetColor(tcellstyl.BorderRight.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " bottom: "
				if tcellstyl.BorderBottom != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderBottom.DashStyle)

					if tcellstyl.BorderBottom.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderBottom.Width.Magnitude)
					}
					if tcellstyl.BorderBottom.Color != nil {
						outstr += fmt.Sprintf(" color: %s", gdocUtil.GetColor(tcellstyl.BorderBottom.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " left: "
				if tcellstyl.BorderLeft != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderLeft.DashStyle)

					if tcellstyl.BorderLeft.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderLeft.Width.Magnitude)
					}
					if tcellstyl.BorderLeft.Color != nil {
						outstr += fmt.Sprintf(" color: %s", gdocUtil.GetColor(tcellstyl.BorderLeft.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += "\n"

				for el:=0; el< len(tcell.Content); el++ {
					tstr := dObj.cvtPar(tcell.Content[el].Paragraph)
					outstr += tstr
				}
			}
            tblCellCount++
		}
	}
	return outstr, nil
}

func (dObj *gdocTxtObj) dispListProp(listp *docs.ListProperties)(outstr string, err error) {

	var nlObj *docs.NestingLevel

	if listp == nil {
		return "", fmt.Errorf(" error dispList: ListProperties is nil!")
	}
	if listp.NestingLevels == nil {
		return "", fmt.Errorf(" error dispList: List NestingLevels is nil!")
	}

	nLvl := len(listp.NestingLevels)
	outstr += fmt.Sprintf("  Nesting Levels: %d\n",nLvl)

	for nl:=0; nl<nLvl; nl++ {
		nlObj = listp.NestingLevels[nl]
		outstr += fmt.Sprintf(" NL: %d Start: %d Indent %4.2f First Line %4.2f\n",nl, nlObj.StartNumber, 
		nlObj.IndentStart.Magnitude, nlObj.IndentFirstLine.Magnitude )
		outstr += fmt.Sprintf("   Bullet Alignment: %s \n", nlObj.BulletAlignment)
		outstr += fmt.Sprintf("   Glyph Format:     %s \n", nlObj.GlyphFormat)
        r, size := utf8.DecodeRuneInString(nlObj.GlyphSymbol)
		outstr += fmt.Sprintf("   GlyphSymbol:      %s %d %v\n", nlObj.GlyphSymbol, r, size)
		outstr += fmt.Sprintf("   GlyphType:        %s \n", nlObj.GlyphType)
		if nlObj.TextStyle == nil {
			outstr += "    *** no text style ***\n"
		} else {
			outstr += "    *** has text style ***\n"
		}
	}

	return outstr, nil
}

func (dObj *gdocTxtObj) cvtPar(par *docs.Paragraph)(outstr string) {

	wsp := "  "
	if par.Bullet == nil {
		for i:=0; i< 10; i++ {dObj.listCount[i] = 0}
	}
	listPrefix := ""
	if par.Bullet != nil {
		// lists tbd
        listid := par.Bullet.ListId
        nestIdx := int(par.Bullet.NestingLevel)

        // retrieve the list properties from the doc.Lists map
        nestL := dObj.doc.Lists[listid].ListProperties.NestingLevels[nestIdx]
        listOrd := gdocUtil.GetGlyphOrd(nestL)

		for j:=0; j< nestIdx ; j++ {wsp += wsp}

		if listOrd {
			dObj.listCount[nestIdx]++
			listPrefix = wsp + fmt.Sprintf("%d.  ", dObj.listCount)
		} else {
			listPrefix = wsp + "*  "
		}

	}
/*
	// indent styles
	if par.ParagraphStyle != nil {
		tstr,err := dObj.dispParStyl(par.ParagraphStyle, 4)
		if err != nil {
			return outstr, fmt.Errorf("error dispParStyl %v", err)
		}
		outstr += tstr
	} else {
		outstr += "    *** no Pargraph Style ***\n"
	}
*/
	parStr :=""
	for pEl:=0; pEl< len(par.Elements); pEl++ {
		parSubEl := par.Elements[pEl]
		if parSubEl.TextRun == nil { continue }
		parStr += parSubEl.TextRun.Content
	}

	tbuf := []byte(parStr)
	wsPos := 0
	pEnd := 0
//	fmt.Printf("cpl: %d %d\n", dObj.cpl, len(tbuf))
	for i:= 0; i<len(tbuf); i++ {
		switch tbuf[i] {
		case '\n':
			wsPos = i
		case ' ':
			wsPos = i
		default:
			if (i - pEnd) > dObj.cpl {
				tbuf[wsPos] = '\n'
//	fmt.Printf("wsPos: %d, pEnd: %d\n", wsPos, pEnd)
				pEnd = wsPos
			}
		}
	}

/*
	// images tbd
	if par.PositionedObjectIds != nil {
		outstr += fmt.Sprintf("    *** Has Positioned Objects: %d ***\n", len(par.PositionedObjectIds))
		for id:=0; id< len(par.PositionedObjectIds); id++ {
			outstr += fmt.Sprintf("      posObject Id[%d]: %s\n", id, par.PositionedObjectIds[id])
		}
	}
*/
	outstr += listPrefix + string(tbuf)
	return outstr
}

func (dObj *gdocTxtObj) cvtParEl(parDet *docs.ParagraphElement)(outstr string, err error) {

	if parDet.TextRun != nil {
		cLen := len(parDet.TextRun.Content)
		outstr += fmt.Sprintf("cl: %d", cLen)
		tCount:=1
		for i:=0; i< cLen; i++ {
			if parDet.TextRun.Content[i] == '\t' {
				outstr += fmt.Sprintf(" tab %d: %d ", tCount, i)
				tCount++
			}
		}
		if cLen > 0 {
			if cLen > 21 {
				outstr += fmt.Sprintf("    \"%s ...\"",parDet.TextRun.Content[0:20])
			} else {
				if parDet.TextRun.Content[cLen-1:cLen] == "\n" {
					outstr += fmt.Sprintf("    \"%s\"",parDet.TextRun.Content[:cLen-1])
				} else {
					outstr += fmt.Sprintf("    \"%s\"",parDet.TextRun.Content[:cLen])
				}
			}
		}
		outstr += fmt.Sprintf("   Last Char: %q\n",parDet.TextRun.Content[cLen-1:cLen])

	} else { outstr +="\n" }

	if parDet.ColumnBreak != nil {
		outstr += "      *** Column Break ***\n"
	} else {
		outstr += "      *** no Column Break ***\n"
	}

	if parDet.InlineObjectElement != nil {
		outstr += fmt.Sprintf("      *** Inline Object with id: %s ***\n", parDet.InlineObjectElement.InlineObjectId)
	} else {
		outstr += "      *** no Inline Object ***\n"
	}

	if parDet.Person != nil {
		outstr += fmt.Sprintf("      *** Has Person with id ***\n", parDet.Person.PersonId)
	} else {
		outstr += "      *** no Person ***\n"
	}

	if parDet.RichLink != nil {
		outstr += fmt.Sprintf("      *** Has Rich Text ***\n")
	} else {
		outstr += "      *** no Rich Text ***\n"
	}

	if parDet.PageBreak != nil {
		insTxt := "without"
		if (parDet.PageBreak.TextStyle != nil) {
			insTxt = "with"
		}
		outstr += fmt.Sprintf("      *** Has Page Break %s txtstyle ***\n", insTxt)
		if (parDet.PageBreak.TextStyle != nil) {
//			outstr += tstr
		}
	} else {
		outstr += "      *** no Page Break ***\n"
	}

	if parDet.AutoText != nil {
		outstr += fmt.Sprintf("      *** Has AutoText ***\n")
	} else {
		outstr += "      *** no AutoText ***\n"
	}

	if parDet.Equation != nil {
		outstr += fmt.Sprintf("      *** Has Equation ***\n")
	} else {
		outstr += "      *** no Equation ***\n"
	}

	if parDet.HorizontalRule != nil {
		insTxt := "without"
		if (parDet.HorizontalRule.TextStyle != nil) {
			insTxt = "with"
		}
		outstr += fmt.Sprintf("      *** Has Horizontal Rule %s txtstyle ***\n", insTxt)
		if (parDet.HorizontalRule.TextStyle != nil) {
//			tstr, err := dObj.dispTxtStyl(parDet.HorizontalRule.TextStyle, 8)
//			outstr += tstr
		}
	} else {
		outstr += "      *** no Horizontal Rule ***\n"
	}

	if parDet.FootnoteReference != nil {
		ftref := parDet.FootnoteReference
		outstr += fmt.Sprintf("      *** Has Footnote Reference ***\n")
		outstr += fmt.Sprintf("          Id:     %s\n", ftref.FootnoteId)
		outstr += fmt.Sprintf("          Number: %s\n", ftref.FootnoteNumber)
//		tstr, err := dObj.dispTxtStyl(ftref.TextStyle, 8)
//		if err == nil {outstr += tstr} else { outstr += fmt.Sprintf("       *** error %v\n", err) }
	} else {
			outstr += "      *** no Footnote Reference ***\n"
	}

	return outstr, nil
}

func (dObj *gdocTxtObj) dispSecStyle(secStyl *docs.SectionStyle)(outstr string, err error) {

	outstr = "*** Section Style ***\n"
	if secStyl == nil {
		return outstr, fmt.Errorf("error dispSecStyl: no secStyl")
	}
	if len(secStyl.SectionType) > 0 {
		outstr += fmt.Sprintf("  SectionType: %s\n", secStyl.SectionType)
	} else {
		outstr += fmt.Sprintf("  SectionType: not specified\n")
	}

	outstr += fmt.Sprintf("  Column Properties: %d", len(secStyl.ColumnProperties))
	for i:=0; i< len(secStyl.ColumnProperties); i++ {
		col := secStyl.ColumnProperties[i]
		outstr += fmt.Sprintf("  Column [%d] Width: %d %s Padding End: %d %s\n",i, col.Width.Magnitude, col.Width.Unit, col.PaddingEnd.Magnitude, col.PaddingEnd.Unit)
	}
	outstr += "  Column Separador:     " + secStyl.ColumnSeparatorStyle + "\n"
	outstr += "  Default Header Id:    " + secStyl.DefaultHeaderId + "\n"
	outstr += "  Default Footer Id:    " + secStyl.DefaultFooterId + "\n"
	outstr += "  Even Page Header Id:  " + secStyl.EvenPageHeaderId + "\n"
	outstr += "  Even Page Footer Id:  " + secStyl.EvenPageFooterId + "\n"
	outstr += "  First Page Header Id: " + secStyl.FirstPageHeaderId + "\n"
	outstr += "  First Page Footer Id: " + secStyl.FirstPageFooterId + "\n"
	outstr += fmt.Sprintf("  Page Number Start:     %d \n", secStyl.PageNumberStart)
	outstr += fmt.Sprintf("  Use First Page H/F:   %t \n", secStyl.UseFirstPageHeaderFooter)
	outstr += "  *** Section Margins ***\n"
	if secStyl.MarginHeader != nil {
		outstr += fmt.Sprintf("    Margin Header: %.1f %s\n",secStyl.MarginHeader.Magnitude, secStyl.MarginHeader.Unit)
	}
	if secStyl.MarginFooter != nil {
		outstr += fmt.Sprintf("    Margin Footer: %.1f %s\n",secStyl.MarginFooter.Magnitude, secStyl.MarginFooter.Unit)
	}
	if secStyl.MarginTop != nil {
		outstr += fmt.Sprintf("    Margin Top:    %.1f %s\n",secStyl.MarginTop.Magnitude, secStyl.MarginTop.Unit)
	}
	if secStyl.MarginRight != nil {
		outstr += fmt.Sprintf("    Margin Right: %.1f %s\n",secStyl.MarginRight.Magnitude, secStyl.MarginRight.Unit)
	}
	if secStyl.MarginBottom != nil {
		outstr += fmt.Sprintf("    Margin Bottom: %.1f %s\n",secStyl.MarginBottom.Magnitude, secStyl.MarginBottom.Unit)
	}
	if secStyl.MarginLeft != nil {
		outstr += fmt.Sprintf("    Margin Left: %.1f %s\n",secStyl.MarginLeft.Magnitude, secStyl.MarginLeft.Unit)
	}
	return outstr, nil
}


func (dObj *gdocTxtObj) dispParStyl(parStyl *docs.ParagraphStyle, wsp int)(outstr string, err error) {

	if parStyl == nil {
		return "", fmt.Errorf("error dispParStyl: no ParStyl")
	}
	wspStr := ""
	for i:=0; i<wsp; i++ {wspStr += " "}

	outstr = wspStr + "*** Paragraph Style ***\n"
	if len(parStyl.HeadingId) > 0 {
		outstr += wspStr + fmt.Sprintf("  Heading Id:  %s \n", parStyl.HeadingId)
	} else {
		outstr += wspStr + fmt.Sprintf("  Heading Id:  none\n")
	}
	outstr += wspStr + fmt.Sprintf("  Named Style: %s \n", parStyl.NamedStyleType)
	if len(parStyl.Alignment) > 0 {
		outstr +=  wspStr + fmt.Sprintf("  Alignment:  %s \n", parStyl.Alignment)
	} else {
		outstr += wspStr + "  Alignment: not specified!\n"
	}
	outstr += wspStr + fmt.Sprintf("  Direction: %s \n", parStyl.Direction)
	if parStyl.LineSpacing > 0 {
		outstr += wspStr + fmt.Sprintf("  Line Spacing: %.2f \n", parStyl.LineSpacing/100.0)
	} else {
		outstr += wspStr + fmt.Sprintf("  Line Spacing not specified!\n")
	}
	outstr += wspStr + fmt.Sprintf("  Keep Lines together: %t \n", parStyl.KeepLinesTogether)
	outstr += wspStr + fmt.Sprintf("  Keep With Next:      %t \n", parStyl.KeepWithNext)

	if parStyl.IndentFirstLine != nil {
		outstr += wspStr + fmt.Sprintf("  Indent First Line: %f %s\n", parStyl.IndentFirstLine.Magnitude, parStyl.IndentFirstLine.Unit )
	}
	if parStyl.IndentStart !=nil {
		outstr += wspStr + fmt.Sprintf("  Indent Start:      %f %s\n", parStyl.IndentStart.Magnitude, parStyl.IndentStart.Unit)
	}
	if parStyl.IndentEnd !=nil {
		outstr += wspStr + fmt.Sprintf("  Indent End:        %f %s\n", parStyl.IndentEnd.Magnitude, parStyl.IndentEnd.Unit)
	}
	if parStyl.Shading !=nil {
		if parStyl.Shading.BackgroundColor != nil {
			if parStyl.Shading.BackgroundColor.Color != nil {
				color := parStyl.Shading.BackgroundColor.Color
    	        blue := int(color.RgbColor.Blue*255.0)
        	    red := int(color.RgbColor.Red*255.0)
            	green := int(color.RgbColor.Green*255)
            	outstr += wspStr + fmt.Sprintf(" Shading background-color: rgb(%d, %d, %d);\n", red, green, blue)
			}
		}
	}

	outstr += wspStr + fmt.Sprintf("  Tabs: %d\n",  len(parStyl.TabStops))

	for i:=0; i< len(parStyl.TabStops); i++ {
		tab := parStyl.TabStops[i]
		outstr += wspStr + fmt.Sprintf("  Tab: %2d ", i)
		outstr += fmt.Sprintf("align: %-10s Offset: %.0f\n", tab.Alignment, tab.Offset.Magnitude)
	}


	return outstr, nil
}


func (dObj *gdocTxtObj) dispBodySummary()(outstr string, err error) {
//    var listItem docList

	body := dObj.doc.Body

	parCount := 0
	listCount := 0
	secCount := 0
	tabCount := 0
	tocCount := 0

	for el:=0; el< len(body.Content); el++ {
		elObj := body.Content[el]
		if elObj.Paragraph != nil {
			parCount++
            if elObj.Paragraph.Bullet != nil {
				listCount++
            }

		}
		if elObj.SectionBreak != nil {secCount++}
		if elObj.Table != nil {tabCount++}
		if elObj.TableOfContents != nil {tocCount++}
	}

	outstr =  fmt.Sprintf("  Paragraphs: %3d\n", parCount)
	outstr += fmt.Sprintf("  Doc Lists:  %3d\n", len(dObj.docLists))
	outstr += fmt.Sprintf("  List Items: %3d\n", listCount)
	outstr += fmt.Sprintf("  Sections:   %3d\n", secCount)
	outstr += fmt.Sprintf("  Tables:     %3d\n", tabCount)
	outstr += fmt.Sprintf("  TOC:        %3d\n", tocCount)

	return outstr, nil
}

func (dObj *gdocTxtObj) cvtSec(sec *docs.SectionBreak)(outstr string, err error) {

	secStyl := sec.SectionStyle
	if secStyl.SectionType == "Next_Page" {
		dObj.pgCounter++
		outstr += fmt.Sprintf("<<page %d>>\n", dObj.pgCounter)
	}

	return outstr, err
}

func (dObj *gdocTxtObj) dispContentEl(elStr *docs.StructuralElement)(outstr string, err error) {

	if elStr == nil {
		return "", fmt.Errorf("error dispContentEl -- elStr is nil")
	}
	if dObj == nil {
		return "", fmt.Errorf("error dispContentEl -- dObj is nil")
	}

	notFound := true
	if elStr.Paragraph != nil {
		listStr := "Paragraph  "
		if elStr.Paragraph.Bullet != nil {listStr = fmt.Sprintf("List Id:%-18s NL:%3d", elStr.Paragraph.Bullet.ListId, elStr.Paragraph.Bullet.NestingLevel)}
		outstr += fmt.Sprintf(" type: %s StartIndex: %d EndIndex: %d\n",  listStr, elStr.StartIndex, elStr.EndIndex)
		notFound = false
//		par := elStr.Paragraph
//		tstr, err := dObj.dispPar(par)
//		if err != nil {return tstr, fmt.Errorf("**error** Par Style: %v", err)}
//		outstr +=tstr
	}
	if elStr.SectionBreak != nil {
		outstr += fmt.Sprintf(" type: Section Break StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
		secStyl := elStr.SectionBreak.SectionStyle
		tstr, err := dObj.dispSecStyle(secStyl)
		if err != nil {
			return tstr, fmt.Errorf("**error** Section Break Style: %v", err)
		}
		outstr += tstr
	}
	if elStr.Table != nil {
		outstr += fmt.Sprintf(" type: Table StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
		tstr, err := dObj.cvtTable(elStr.Table)
		if err != nil {
			return tstr, fmt.Errorf("**error** disp Table: %v", err)
		}
		outstr += tstr
	}
	if elStr.TableOfContents != nil {
		outstr += fmt.Sprintf(" type: TOC StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
		tstr, err := dObj.cvtTOC(elStr.TableOfContents)
		if err != nil {
			return tstr, fmt.Errorf("**error** disp TOC: %v", err)
		}
		outstr += tstr
	}
	if notFound {
		return outstr, fmt.Errorf(" type: unknown StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
	}
	return outstr, nil
}

func (dObj *gdocTxtObj) cvtBody()(outstr string, err error) {

	body := dObj.doc.Body
	err = nil
	errCount := 0
	for el:=0; el< len(body.Content); el++ {
		elBody := body.Content[el]
		if elBody.Paragraph != nil {
			tstr := dObj.cvtPar(elBody.Paragraph)
			outstr += tstr
		}
		if elBody.SectionBreak != nil {
			tstr, err := dObj.cvtSec(elBody.SectionBreak)
			if err != nil {
				errCount++
				outstr += fmt.Sprintf("error %d: %v\n", errCount, err)
			}
			outstr += tstr
		}
		if elBody.Table != nil {
			tstr, err := dObj.cvtTable(elBody.Table)
			if err != nil {
				errCount++
				outstr += fmt.Sprintf("error %d: %v\n", errCount, err)
			}
			outstr += tstr
		}

	}
	if errCount > 0 {err = fmt.Errorf("error cvtBody count: %d", errCount)}
	return outstr, err
}


func CvtGdocToTxt(folderPath string, doc *docs.Document, options *gdocUtil.OptObj)(err error) {
	var outstr string

	if doc == nil { return fmt.Errorf("error -- doc is nil!")}

    docObj := new(gdocTxtObj)
    docObj.doc = doc
	docObj.DocOpt = true

	if !(len(doc.Title) >0) {
		return fmt.Errorf("error CvtGdocToTxt:: the string doc.Title %s is too short!", doc.Title)
	}

	// initialise docObj
    err = docObj.initGdocTxt(folderPath, options)
    if err != nil {
        return fmt.Errorf("error Cvt Txt Init %v", err)
    }

	outfil := docObj.txtfil

	outstr = fmt.Sprintf("Document Id: %s \n", doc.DocumentId)
	outstr += fmt.Sprintf("Revision Id: %s \n\n", doc.RevisionId)
	outfil.WriteString(outstr)

	outstr = ""
	tstr, err := docObj.cvtBody()
	if err != nil {
		outstr += "*** error cvtBody ***\n"
	}
	outstr += tstr
	outfil.WriteString(outstr)

	// external objects
	inObjLen := len(doc.InlineObjects)
	posObjLen := len(doc.PositionedObjects)

	// Inline Objects
	outstr = fmt.Sprintf("\n*** Inline Objects: %d ***\n", inObjLen)
	for key,inlObj :=  range doc.InlineObjects {
		emObj := inlObj.InlineObjectProperties.EmbeddedObject
		outstr+= fmt.Sprintf("key: %s Title: %s H: %.2f W:%.2f\n", key, emObj.Title, emObj.Size.Height.Magnitude,emObj.Size.Width.Magnitude)
		outstr+= fmt.Sprintf("  Content Uri: %s Source Uri: %s\n", emObj.ImageProperties.ContentUri, emObj.ImageProperties.SourceUri )
	}
	outfil.WriteString(outstr)

	// positioned objects
	outstr = fmt.Sprintf("\n*** Positioned Objects: %d ***\n",posObjLen)
	for key, posObj :=  range doc.PositionedObjects {
		emObj := posObj.PositionedObjectProperties.EmbeddedObject
		objPos := posObj.PositionedObjectProperties.Positioning
		outstr += fmt.Sprintf("key: %s Title: %s H: %f W:%f\n", key, emObj.Title, emObj.Size.Height.Magnitude,emObj.Size.Width.Magnitude)
		outstr+= fmt.Sprintf("Content Uri: %s Source Uri: %s\n", emObj.ImageProperties.ContentUri, emObj.ImageProperties.SourceUri )
		outstr += fmt.Sprintf("    Margins (TRBL): %.2f %.2f %.2f %.2f\n", emObj.MarginTop.Magnitude, emObj.MarginRight.Magnitude, emObj.MarginBottom.Magnitude, 
			emObj.MarginLeft.Magnitude)
		outstr += fmt.Sprintf("    Layout: %s Pos Top: %.2f Left: %.2f\n", objPos.Layout, objPos.TopOffset.Magnitude, objPos.LeftOffset.Magnitude)
	}
	outfil.WriteString(outstr)


	// footnotes
	ftnoteLen := len(doc.Footnotes)
	outstr = fmt.Sprintf("\n*** Footnotes: %d ***\n",ftnoteLen)
	knum := 0
	for key, ftnote := range doc.Footnotes {
		knum++
		outstr += fmt.Sprintf("  ftnote %2d: key: %-16s elements: %2d", knum, key, len(ftnote.Content) )
		for i:= 0; i< len(ftnote.Content); i++ {
			if ftnote.Content[i].Paragraph == nil {continue}
			text := ""
			for j:=0; j<len(ftnote.Content[i].Paragraph.Elements); j++ {
				el := ftnote.Content[i].Paragraph.Elements[j]
				if el.TextRun == nil {continue}
				text += el.TextRun.Content
			}
			outstr += fmt.Sprintf(" text: %s", text)
		}
	}
	outfil.WriteString(outstr)

// footnote occurence
	icnt := 0
	outstr = ""
	for i:=0; i< len(doc.Body.Content); i++ {
		if doc.Body.Content[i].Paragraph == nil {continue}
		par:= doc.Body.Content[i].Paragraph
		for j:=0; j<len(par.Elements); j++ {
			if par.Elements[j].FootnoteReference == nil {continue}
			icnt++
			ftn := par.Elements[j].FootnoteReference
			outstr += fmt.Sprintf("  ftnote: %2d struct: %3d el %3d Id: %-16s Number: %-2s\n", icnt, i, j, ftn.FootnoteId, ftn.FootnoteNumber)
		}
	}
	outfil.WriteString(outstr)


	outfil.Close()
	return nil
}


