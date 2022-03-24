//  gdocTxtLib
//  reads gdoc file and writes text version to outfile
//  author PRR
//  created 2022
//
//  13/1/2022 added unicode for glyphtype
//	add empty line before each list
//
//  24/2/2022 inline objects
//
//  21/3/2022 fomatting update

package gdocToText

import (
	"fmt"
	"os"
	"unicode/utf8"
	"google.golang.org/api/docs/v1"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdocTxtObj struct {
	parCount int
	posImgCount int
 	inImgCount int
	TableCount int
	imgId []string
	doc *docs.Document
	DocName string
	headingId []string
}

func (dObj *gdocTxtObj) InitGdocTxt() (err error) {
	if dObj == nil {
		return fmt.Errorf("error Init: dObj is nil!")
	}
	dObj.parCount = 0
	dObj.posImgCount = 0
	dObj.inImgCount = 0
	dObj.TableCount = 0
	return nil
}

func (dObj *gdocTxtObj) dispTOC(toc *docs.TableOfContents)(outstr string, err error) {

	if toc == nil {
		return "", fmt.Errorf("error dispTOC: toc is nil")
	}

	if dObj == nil {
		return "", fmt.Errorf("error dispTOC: dObj is nil")
	}

	numEl := len(toc.Content)
	outstr = fmt.Sprintf("\n*** Table of Contents: (%d) ***\n", numEl)

	for el:=0; el< numEl; el++ {
		tocEl := toc.Content[el]
		outstr += fmt.Sprintf("\n *** element: %d ***", el)
		tStr,err := dObj.dispContentEl(tocEl)
		if err != nil {
			outstr += fmt.Sprintf("\n error dispContentEl: %v\n", err) + tStr
		} else {
			outstr += tStr
		}
	}
	return outstr, nil
}

func (dObj *gdocTxtObj) dispTable(tbl *docs.Table)(outstr string, err error) {
	var tabWidth, tabHeight float64
	var icol, irow int64

	doc := dObj.doc

	if tbl == nil {
		return "", fmt.Errorf("error dispTable: no table pt")
	}
	dObj.TableCount++
	outstr = fmt.Sprintf("*** table %d: rows: %d cols: %d ***\n",dObj.TableCount,tbl.Rows, tbl.Columns )
	// table rows

	outstr += "Table Style Properties\n"
	tabWidth = 0.0
	if tbl.TableStyle != nil {
		numColProp := (int64)(len(tbl.TableStyle.TableColumnProperties))
		for icol=0; icol<numColProp; icol++ {
			tColProp := tbl.TableStyle.TableColumnProperties[icol]
			outstr += fmt.Sprintf(" col[%d]: w type: %s", icol, tColProp.WidthType)
			if tColProp.Width !=nil {
				outstr += fmt.Sprintf(" width: %.1f", tColProp.Width.Magnitude*PtTomm)
			}
			outstr += "\n"
		}
	}

	docPg := doc.DocumentStyle
	PgWidth := docPg.PageSize.Width.Magnitude
	NetPgWidth := PgWidth - (docPg.MarginLeft.Magnitude + docPg.MarginRight.Magnitude)
	outstr += fmt.Sprintf("Default Table Width: %.1f", NetPgWidth*PtTomm)
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
					outstr += fmt.Sprintf(" color: %s", dObj.getColor(tcellstyl.BackgroundColor.Color))
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
						outstr += fmt.Sprintf(" color: %s", dObj.getColor(tcellstyl.BorderTop.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " right: "
				if tcellstyl.BorderRight != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderRight.DashStyle)

					if tcellstyl.BorderRight.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderRight.Width.Magnitude)
					}
					if tcellstyl.BorderRight.Color != nil {
						outstr += fmt.Sprintf(" color: %s", dObj.getColor(tcellstyl.BorderRight.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " bottom: "
				if tcellstyl.BorderBottom != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderBottom.DashStyle)

					if tcellstyl.BorderBottom.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderBottom.Width.Magnitude)
					}
					if tcellstyl.BorderBottom.Color != nil {
						outstr += fmt.Sprintf(" color: %s", dObj.getColor(tcellstyl.BorderBottom.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += " left: "
				if tcellstyl.BorderLeft != nil {
					outstr += fmt.Sprintf(" dash %s", tcellstyl.BorderLeft.DashStyle)

					if tcellstyl.BorderLeft.Width != nil {
						outstr += fmt.Sprintf(" width: %.1f", tcellstyl.BorderLeft.Width.Magnitude)
					}
					if tcellstyl.BorderLeft.Color != nil {
						outstr += fmt.Sprintf(" color: %s", dObj.getColor(tcellstyl.BorderLeft.Color.Color))
					}
				} else { outstr += " <no char>" }

				outstr += "\n"

				for el:=0; el< len(tcell.Content); el++ {
					tstr, err := dObj.dispPar(tcell.Content[el].Paragraph)
					if err != nil {
						outstr += fmt.Sprintf("error par of tablecel %d\n",el)
					} else {
						outstr += tstr
					}
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

func (dObj *gdocTxtObj) dispPar(par *docs.Paragraph)(outstr string, err error) {

	if par == nil {
		return "", fmt.Errorf("error dispPar: no par element provided! ")
	}

	if par.Bullet != nil {
		outstr += fmt.Sprintf(" *** List Paragraph with %d Sub-Elements ***\n", len(par.Elements))
		outstr += fmt.Sprintf("       list id: %s nest level %d \n", par.Bullet.ListId, par.Bullet.NestingLevel)
	} else {
		outstr += fmt.Sprintf(" *** Paragraph with %d Sub-Elements ***\n", len(par.Elements))
	}
//xxx
	if par.ParagraphStyle != nil {
		tstr,err := dObj.dispParStyl(par.ParagraphStyle, 4)
		if err != nil {
			return outstr, fmt.Errorf("error dispParStyl %v", err)
		}
		outstr += tstr
	} else {
		outstr += "    *** no Pargraph Style ***\n"
	}
	outstr += fmt.Sprintf("\n  *** Paragraph Elements: %d ***\n", len(par.Elements))
	for p:=0; p< len(par.Elements); p++ {
		parDet := par.Elements[p]
		outstr += fmt.Sprintf("  Par-El[%d]: (%d-%d): ", p, parDet.StartIndex, parDet.EndIndex)
		if parDet.TextRun != nil {
			cLen := len(parDet.TextRun.Content)
			outstr += fmt.Sprintf("cl: %d", cLen)
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
			outstr += fmt.Sprintf(" Last Char: %q\n",parDet.TextRun.Content[cLen-1:cLen])
			if (parDet.TextRun.TextStyle != nil) {
				tstr, err := dObj.dispTxtStyl(parDet.TextRun.TextStyle, 4)
				if err != nil {
					outstr += fmt.Sprintf("/* error disp Text Style: %v */\n", err)
				}
				outstr += tstr
			}
		} else { outstr +="\n" }
		if parDet.ColumnBreak != nil {
			outstr += "    *** Column Break ***\n"
		} else {
			outstr += "    *** no Column Break ***\n"
		}
		if parDet.InlineObjectElement != nil {
			outstr += "    *** Inline Object ***\n"
		} else {
			outstr += "    *** no Inline Object ***\n"
		}
		if parDet.Person != nil {
			outstr += fmt.Sprintf("    *** Has Person ***\n")
		} else {
			outstr += "    *** no Person ***\n"
		}
		if parDet.RichLink != nil {
			outstr += fmt.Sprintf("    *** Has Rich Text ***\n")
		} else {
			outstr += "    *** no Rich Text ***\n"
		}
		if parDet.PageBreak != nil {
			outstr += fmt.Sprintf("    *** Has Page Break ***\n")
		} else {
			outstr += "    *** no Page Break ***\n"
		}
		if parDet.AutoText != nil {
			outstr += fmt.Sprintf("    *** Has AutoText ***\n")
		} else {
			outstr += "    *** no AutoText ***\n"
		}
		if parDet.Equation != nil {
			outstr += fmt.Sprintf("    *** Has Equation ***\n")
		} else {
			outstr += "    *** no Equation ***\n"
		}
		if parDet.HorizontalRule != nil {
			outstr += fmt.Sprintf("    *** Has Horizontal Rule ***\n")
		} else {
			outstr += "    *** no Horizontal Rule ***\n"
		}
		if parDet.FootnoteReference != nil {
			ftref := parDet.FootnoteReference
			outstr += fmt.Sprintf("    *** Has Footnote Reference ***\n")
			outstr += fmt.Sprintf("       Id:     %s\n", ftref.FootnoteId)
			outstr += fmt.Sprintf("       Number: %s\n", ftref.FootnoteNumber)
			tstr, err := dObj.dispTxtStyl(ftref.TextStyle, 8)
			if err == nil {outstr += tstr} else { outstr += fmt.Sprintf("     *** error %v\n", err) }
		} else {
			outstr += "    *** no Footnote Reference ***\n"
		}
	}
	if par.PositionedObjectIds != nil {
		outstr += fmt.Sprintf("    *** Has Positioned Objects: %d ***\n", len(par.PositionedObjectIds))
		for id:=0; id< len(par.PositionedObjectIds); id++ {
			outstr += fmt.Sprintf("      posObject Id[%d]: %s\n", id, par.PositionedObjectIds[id])
		}
	}


	return outstr, nil
}

func (dObj *gdocTxtObj) dispSecStyle(secStyl *docs.SectionStyle)(outstr string, err error) {

	outstr = "*** Section Style ***\n"
	if secStyl == nil {
		return outstr, fmt.Errorf("error dispSecStyl: no secStyl")
	}
	if len(secStyl.SectionType) > 0 {
		outstr += fmt.Sprintf("SectionType:          %s\n", secStyl.SectionType)
	} else {
		outstr += fmt.Sprintf("SectionType:   		not specified\n")
	}

	outstr += fmt.Sprintf(" Column Properties: %d", len(secStyl.ColumnProperties))
	for i:=0; i< len(secStyl.ColumnProperties); i++ {
		col := secStyl.ColumnProperties[i]
		outstr += fmt.Sprintf(" Column [%d] Width: %d %s Padding End: %d %s\n",i, col.Width.Magnitude, col.Width.Unit, col.PaddingEnd.Magnitude, col.PaddingEnd.Unit)
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


func (dObj *gdocTxtObj) dispBorder(parBorder *docs.ParagraphBorder, wsp int)(outstr string) {

	if parBorder == nil {
		outstr += "error dispParStyl: no Par Border Style\n"
		return outstr
	}

	wspStr := ""
	for i:=0; i<wsp; i++ {wspStr += " "}
	outstr = wspStr + fmt.Sprintf("    Border Style:   %s\n", parBorder.DashStyle)
	outstr += wspStr + fmt.Sprintf("    Border Width:   %.1f %s\n", parBorder.Width.Magnitude, parBorder.Width.Unit)
	outstr += wspStr + fmt.Sprintf("    Border Padding: %.1f %s\n", parBorder.Padding.Magnitude, parBorder.Padding.Unit)
	if parBorder.Color != nil {
		if parBorder.Color.Color != nil {
			colStr := dObj.getColor(parBorder.Color.Color)
			outstr += wspStr + fmt.Sprintf("    Border Color: %s\n", colStr)
		}
	}
	return outstr
}

func (dObj *gdocTxtObj) getColor(color  *docs.Color)(outstr string) {
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

func (dObj *gdocTxtObj) dispParStyl(parStyl *docs.ParagraphStyle, wsp int)(outstr string, err error) {

	if parStyl == nil {
		return "", fmt.Errorf("error dispParStyl: no ParStyl")
	}
	wspStr := ""
	for i:=0; i<wsp; i++ {wspStr += " "}

	outstr = wspStr + "*** Paragraph Style ***\n"
	outstr += wspStr + fmt.Sprintf("  Heading Id:  %s \n", parStyl.HeadingId)
	outstr += wspStr + fmt.Sprintf("  Named Style: %s \n", parStyl.NamedStyleType)
	if len(parStyl.Alignment) > 0 {
		outstr +=  wspStr + fmt.Sprintf("  Alignment:  %s \n", parStyl.Alignment)
	}
	outstr += wspStr + fmt.Sprintf("  Direction: %s \n", parStyl.Direction)
	if parStyl.LineSpacing > 0 {
		outstr += wspStr + fmt.Sprintf("  Line Spacing: %.2f \n", parStyl.LineSpacing/100.0)
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
	if parStyl.BorderBetween != nil {
		outstr += wspStr + "  *** Border Between ***\n"
		outstr += dObj.dispBorder(parStyl.BorderBetween, 4)
		outstr += "\n"
	}

	if parStyl.BorderTop != nil {
		outstr += wspStr + "  *** Border Top ***\n"
		outstr += dObj.dispBorder(parStyl.BorderTop, 4)
		outstr += "\n"
	}

	if parStyl.BorderRight != nil {
		outstr += wspStr + "  *** Border Right ***\n"
		outstr += dObj.dispBorder(parStyl.BorderRight, 4)
		outstr += "\n"
	}

	if parStyl.BorderBottom != nil {
		outstr += wspStr + "  *** Border Bottom ***\n"
		outstr += dObj.dispBorder(parStyl.BorderBottom, 4)
		outstr += "\n"
	}

	if parStyl.BorderLeft != nil {
		outstr += wspStr + "  *** Border Left ***\n"
		outstr += dObj.dispBorder(parStyl.BorderLeft, 4)
		outstr += "\n"
	}

	return outstr, nil
}

func (dObj *gdocTxtObj) dispTxtStyl(txtStyl *docs.TextStyle, wsp int)(outstr string, err error) {

	if txtStyl == nil {
		return "", fmt.Errorf("error dispTxtStyl: TextStyle is nil!")
	}
	wspStr :=""
	for i:=0; i< wsp; i++ {wspStr +=" "}

	outstr = wspStr + "*** Text Style ***\n"
	if len(txtStyl.BaselineOffset) > 0 {
		outstr +=  wspStr + fmt.Sprintf("    BaseLine Offset: %s\n",txtStyl.BaselineOffset)
	}
	outstr += wspStr + fmt.Sprintf("    Bold:      %t\n", txtStyl.Bold)
	outstr += wspStr + fmt.Sprintf("    Italic:    %t\n", txtStyl.Italic)
	outstr += wspStr + fmt.Sprintf("    Underline: %t\n", txtStyl.Underline)
	outstr += wspStr + fmt.Sprintf("    strike through: %t\n", txtStyl.Strikethrough)
	outstr += wspStr + fmt.Sprintf("    small caps: %t\n", txtStyl.SmallCaps)

	if txtStyl.WeightedFontFamily != nil {
		outstr += wspStr + fmt.Sprintf("    Font : %s %d\n", txtStyl.WeightedFontFamily.FontFamily, txtStyl.WeightedFontFamily.Weight)
	}
	if txtStyl.FontSize != nil {
		outstr += wspStr + fmt.Sprintf("    Font Size: %f %s\n", txtStyl.FontSize.Magnitude, txtStyl.FontSize.Unit)
	}

    if txtStyl.ForegroundColor != nil {
        if txtStyl.ForegroundColor.Color != nil {
            outstr += wspStr + "     foreground-color: "
            outstr += dObj.getColor(txtStyl.ForegroundColor.Color)
            outstr += "\n"
        }
    }
    if txtStyl.BackgroundColor != nil {
        if txtStyl.BackgroundColor.Color != nil {
            outstr += wspStr + fmt.Sprintf("     background-color:  ")
            outstr += dObj.getColor(txtStyl.BackgroundColor.Color)
            outstr += "\n"
        }
    }
    if txtStyl.Link != nil {
		outstr += wspStr + "     Link TBD\n"
	}
	return outstr, nil
}
func (dObj *gdocTxtObj) dispBodySummary()(outstr string, err error) {

	body := dObj.doc.Body

	parCount := 0
	listCount := 0
	secCount := 0
	tabCount := 0
	tocCount := 0

	for el:=0; el< len(body.Content); el++ {
		elObj := body.Content[el]
		if elObj.Paragraph != nil {
			if elObj.Paragraph.Bullet == nil {parCount++} else {listCount++}
		}
		if elObj.SectionBreak != nil {secCount++}
		if elObj.Table != nil {tabCount++}
		if elObj.TableOfContents != nil {tocCount++}
	}

	outstr =  fmt.Sprintf("  Paragraphs: %3d\n", parCount)
	outstr += fmt.Sprintf("  Lists:      %3d\n", listCount)
	outstr += fmt.Sprintf("  Sections:   %3d\n", secCount)
	outstr += fmt.Sprintf("  Tables:     %3d\n", tabCount)
	outstr += fmt.Sprintf("  TOC:        %3d\n", tocCount)

	return outstr, nil
}

func (dObj *gdocTxtObj) dispContentEl(elStr *docs.StructuralElement)(outstr string, err error) {

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
		outstr += fmt.Sprintf(" type: Paragraph StartIndex: %d EndIndex: %d\n",  elStr.StartIndex, elStr.EndIndex)
		notFound = false
		par := elStr.Paragraph
		tstr, err := dObj.dispPar(par)
		if err != nil {
			return tstr, fmt.Errorf("**error** Par Style: %v", err)
		}
		outstr +=tstr
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
		tstr, err := dObj.dispTable(elStr.Table)
		if err != nil {
			return tstr, fmt.Errorf("**error** disp Table: %v", err)
		}
		outstr += tstr
	}
	if elStr.TableOfContents != nil {
		outstr += fmt.Sprintf(" type: TOC StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
		notFound = false
		tstr, err := dObj.dispTOC(elStr.TableOfContents)
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

func (dObj *gdocTxtObj) dispContentElShort(elStr *docs.StructuralElement)(outstr string, err error) {

	if elStr == nil {
		return "", fmt.Errorf("error dispContentEl -- elStr is nil")
	}
	if dObj == nil {
		return "", fmt.Errorf("error dispContentEl -- dObj is nil")
	}

	if elStr.Paragraph != nil {
		outstr += fmt.Sprintf(" type: Paragraph StartIndex: %d EndIndex: %d\n",  elStr.StartIndex, elStr.EndIndex)
	}

	if elStr.SectionBreak != nil {
		outstr += fmt.Sprintf(" type: Section Break StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
	}

	if elStr.Table != nil {
		outstr += fmt.Sprintf(" type: Table StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
	}

	if elStr.TableOfContents != nil {
		outstr += fmt.Sprintf(" type: TOC StartIndex: %d EndIndex: %d\n", elStr.StartIndex, elStr.EndIndex)
	}

	return outstr, nil
}

func CvtGdocToTxt(outfil *os.File, doc *docs.Document)(err error) {
	var outstr string

    docObj := new(gdocTxtObj)
    docObj.doc = doc
    err = docObj.InitGdocTxt()
    if err != nil {
        return fmt.Errorf("error Cvt Txt Init %v", err)
    }
	if !(len(doc.Title) >0) {
		return fmt.Errorf("error CvtGdocToTxt:: the string doc.Title %s is too short!", doc.Title)
	}
	_, err = outfil.WriteString("*** Document Title: " + doc.Title + " ***\n")
	if err != nil {
		return fmt.Errorf("error CvtGdocToTxt -- cannot write to file: %v", err)
	}

	outstr = fmt.Sprintf("Document Id: %s \n", doc.DocumentId)
	outstr += fmt.Sprintf("Revision Id: %s \n", doc.RevisionId)

// Inline Objects

	inObjLen := len(doc.InlineObjects)
	posObjLen := len(doc.PositionedObjects)
/*
	if (inObjLen + posObjLen) > 0 {
		err = createImgFolder(doc.Title)
		if err != nil {
			return fmt.Errorf("error CvtGdocToTxt:: cannot create image folder: %v",err)
		}
	}
*/
	outstr += fmt.Sprintf("\n*** Inline Objects: %d ***\n", inObjLen)
	for key,inlObj :=  range doc.InlineObjects {
		emObj := inlObj.InlineObjectProperties.EmbeddedObject
		outstr+= fmt.Sprintf("key: %s Title: %s H: %.2f W:%.2f\n", key, emObj.Title, emObj.Size.Height.Magnitude,emObj.Size.Width.Magnitude)
		outstr+= fmt.Sprintf("  Content Uri: %s Source Uri: %s\n", emObj.ImageProperties.ContentUri, emObj.ImageProperties.SourceUri )
	}

	outstr += fmt.Sprintf("\n*** Positioned Objects: %d ***\n",posObjLen)
	for key, posObj :=  range doc.PositionedObjects {
		emObj := posObj.PositionedObjectProperties.EmbeddedObject
		objPos := posObj.PositionedObjectProperties.Positioning
		outstr += fmt.Sprintf("key: %s Title: %s H: %f W:%f\n", key, emObj.Title, emObj.Size.Height.Magnitude,emObj.Size.Width.Magnitude)
		outstr+= fmt.Sprintf("Content Uri: %s Source Uri: %s\n", emObj.ImageProperties.ContentUri, emObj.ImageProperties.SourceUri )
		outstr += fmt.Sprintf("    Margins (TRBL): %.2f %.2f %.2f %.2f\n", emObj.MarginTop.Magnitude, emObj.MarginRight.Magnitude, emObj.MarginBottom.Magnitude, 
			emObj.MarginLeft.Magnitude)
		outstr += fmt.Sprintf("    Layout: %s Pos Top: %.2f Left: %.2f\n", objPos.Layout, objPos.TopOffset.Magnitude, objPos.LeftOffset.Magnitude)
	}

	headLen := len(doc.Headers)
	outstr += fmt.Sprintf("\n*** Headers: %d ****\n",headLen)
	knum := 0
	for key, header := range doc.Headers {
		knum++
		outstr += fmt.Sprintf("  header %d: key %s  id %s elements %s\n", knum, key,header.HeaderId, len(header.Content) )
		for i:= 0; i< len(header.Content); i++ {
			outstr += fmt.Sprintf("  content [%d]\n", i)
			tstr, err := docObj.dispContentEl(header.Content[i])
			if err == nil {outstr += tstr} else {outstr += fmt.Sprintf("error %v", err)}
		}
	}

	footLen := len(doc.Footers)
	outstr += fmt.Sprintf("\n*** Footers: %d ***\n",footLen)
	knum = 0
	for key, footer := range doc.Footers {
		knum++
		outstr += fmt.Sprintf("  footer %d: key %s id %s elements: %d\n", knum, key, footer.FooterId, len(footer.Content) )
		for i:= 0; i< len(footer.Content); i++ {
			outstr += fmt.Sprintf("  content [%d]\n", i)
			tstr, err := docObj.dispContentEl(footer.Content[i])
			if err == nil {outstr += tstr} else {outstr += fmt.Sprintf("error %v", err)}
		}
	}

	ftnoteLen := len(doc.Footnotes)
	outstr += fmt.Sprintf("\n*** Footnotes: %d ***\n",ftnoteLen)
	knum = 0
	for key, ftnote := range doc.Footnotes {
		knum++
		outstr += fmt.Sprintf("  ftnote %d: key %s id %s elements %d\n", knum, key, ftnote.FootnoteId, len(ftnote.Content) )
		for i:= 0; i< len(ftnote.Content); i++ {
			outstr += fmt.Sprintf("  content [%d]\n", i)
			tstr, err := docObj.dispContentEl(ftnote.Content[i])
			if err == nil {outstr += tstr} else {outstr += fmt.Sprintf("error %v", err)}
		}
	}
	outfil.WriteString(outstr)

	// Lists
	listLen := len(doc.Lists)
	outstr = fmt.Sprintf("\n*** Lists: %d ***\n",listLen)
	knum = 0
	for key, list := range doc.Lists {
		knum++
		nest := list.ListProperties.NestingLevels
		outstr += fmt.Sprintf("    List %2d: id: %-20s nest levels: %d\n", knum, key, len(nest) )
/*
		tstr, err := docObj.dispListProp(list.ListProperties)
		if err != nil {
			outstr += fmt.Sprintf("error dispLists: %v\n",err)
			break
		}
		outstr += tstr
*/
	}
	outfil.WriteString(outstr)

	nrLen := len(doc.NamedRanges)
	outstr = fmt.Sprintf("\n*** NamedRanges: %d ***\n",nrLen)
	knum = 0
	for key, namrange := range doc.NamedRanges {
		knum++
		outstr += fmt.Sprintf("  nam range %d: key %s id %s elements %d\n", knum, key, namrange.Name, len(namrange.NamedRanges) )
	}

	outfil.WriteString(outstr)

	hdstyles := doc.NamedStyles
	hdstyLen := len(hdstyles.Styles)
	outstr = fmt.Sprintf("\n*** Named Styles: %d ***\n", hdstyLen)
	for i:=0; i<hdstyLen; i++ {
		stylel := hdstyles.Styles[i]
		outstr += fmt.Sprintf("    Style[%1d]: %s \n", i,  stylel.NamedStyleType)
/*
		parStyl := stylel.ParagraphStyle
		parstr, err := docObj.dispParStyl(parStyl, 4)
		if err != nil {
			return fmt.Errorf("error dispParStyl %d: %v", i, err)
		}
		outstr += parstr
		txtStyl:= stylel.TextStyle
		txtstr,err := docObj.dispTxtStyl(txtStyl, 4)
		if err != nil {
			return fmt.Errorf("error dispTxtStyl %d: %v", i, err)
		}
		outstr += txtstr
*/
	}

	body := doc.Body
	hdCount:=0
	hdstr := ""
	for el:=0; el< len(body.Content); el++ {
		elObj := body.Content[el]
		if elObj.Paragraph != nil {
			if len(elObj.Paragraph.ParagraphStyle.HeadingId) > 0 {
				hdstr += fmt.Sprintf("heading [%d]: namedStyle: %s Heading Id: %s\n", hdCount,
					elObj.Paragraph.ParagraphStyle.NamedStyleType,
					elObj.Paragraph.ParagraphStyle.HeadingId)
				hdCount++
			}
		}
	}
	outstr = fmt.Sprintf("\n*** Headings: %d ***\n", hdCount)
	outstr += hdstr
	outfil.WriteString(outstr)

	outstr ="\n****************** Document Style **************************\n"
	docstyl := doc.DocumentStyle
	outstr += fmt.Sprintf("Height: %f %s  Width %f %s\n",docstyl.PageSize.Height.Magnitude,	docstyl.PageSize.Height.Unit,docstyl.PageSize.Width.Magnitude,docstyl.PageSize.Width.Unit)

	outstr += fmt.Sprintf("margin top: %f %s \n",docstyl.MarginTop.Magnitude, docstyl.MarginTop.Unit)
	outstr += fmt.Sprintf("margin right: %f %s \n",docstyl.MarginRight.Magnitude, docstyl.MarginTop.Unit)
	outstr += fmt.Sprintf("margin bottom: %f %s \n",docstyl.MarginBottom.Magnitude, docstyl.MarginTop.Unit)
	outstr += fmt.Sprintf("margin left: %f %s \n",docstyl.MarginLeft.Magnitude, docstyl.MarginTop.Unit)

	outfil.WriteString(outstr)


	outstr = "\n******************** Body *********************************\n"
	numEl := len(body.Content)
	outstr += fmt.Sprintf("*** Body - Number of Elements: %d ***\n", numEl)
	outstr += fmt.Sprintf("*** Body - Element Summary ***\n")

	sumStr, el := docObj.dispBodySummary()
	if err != nil {
		outstr += fmt.Sprintf("\n *** error dispContent[%d]: %v\n", el, err)
	}
	outstr += sumStr

	outfil.WriteString(outstr)

	outstr = fmt.Sprintf("*** Body - Element Detail ***\n")

	for el:=0; el< numEl; el++ {
		elObj := body.Content[el]
		outstr += fmt.Sprintf("   element: %d ", el)
		tStr,err := docObj.dispContentElShort(elObj)
		if err != nil {
			outstr += fmt.Sprintf("\n error dispContent[%d]: %v\n", el, err) + tStr
		} else {
			outstr += tStr
		}
	}
	outfil.WriteString(outstr)

	outstr = fmt.Sprintf("\n\n********************** Details *********************\n\n")
	outstr += fmt.Sprintf("*** Body - Elements: %d ***\n", numEl)

	for el:=0; el< numEl; el++ {
		elObj := body.Content[el]
		outstr += fmt.Sprintf("\n element: %d ", el)
		tStr,err := docObj.dispContentEl(elObj)
		if err != nil {
			outstr += fmt.Sprintf("\n error dispContent[%d]: %v\n", el, err) + tStr
		} else {
			outstr += tStr
		}
	}
	outfil.WriteString(outstr)

	outstr = fmt.Sprintf("\n*** Lists (Detail): %d ***\n",len(doc.Lists))
	knum = 0
	for key, list := range doc.Lists {
		knum++
		nest := list.ListProperties.NestingLevels
		outstr += fmt.Sprintf("\nList [%d]: id: %s nest levels: %d\n", knum, key, len(nest) )
		tstr, err := docObj.dispListProp(list.ListProperties)
		if err != nil {
			outstr += fmt.Sprintf("error dispLists: %v\n",err)
			break
		}
		outstr += tstr
	}
	outfil.WriteString(outstr)

	outstr = "\n\n*** Named Styles ***\n"
	hdstyles = doc.NamedStyles
	hdstyLen = len(hdstyles.Styles)
	outstr += fmt.Sprintf("document named styles: %d \n",hdstyLen)

	for i:=0; i<hdstyLen; i++ {
		stylel := hdstyles.Styles[i]
		outstr += fmt.Sprintf("\nstyle[%1d]: %s \n", i,  stylel.NamedStyleType)
		parStyl := stylel.ParagraphStyle
		parstr, err := docObj.dispParStyl(parStyl, 4)
		if err != nil {
			return fmt.Errorf("error dispParStyl %d: %v", i, err)
		}
		outstr += parstr
		txtStyl:= stylel.TextStyle
		txtstr,err := docObj.dispTxtStyl(txtStyl, 4)
		if err != nil {
			return fmt.Errorf("error dispTxtStyl %d: %v", i, err)
		}
		outstr += txtstr
	}
	outfil.WriteString(outstr)

	outstr = "\n******************** Body *********************************\n"
	outstr += fmt.Sprintf("Body - Number of Elements: %d\n", numEl)

	for el:=0; el< numEl; el++ {
		elObj := body.Content[el]
		outstr += fmt.Sprintf("\n element: %d ", el)
		tStr,err := docObj.dispContentEl(elObj)
		if err != nil {
			outstr += fmt.Sprintf("\n error dispContent[%d]: %v\n", el, err) + tStr
		} else {
			outstr += tStr
		}
	}
	outfil.WriteString(outstr)

	outfil.Close()
	return nil
}

func createImgFolder(filnam string)(err error) {

	if len(filnam) < 2 {
		return fmt.Errorf("error createIMgFolder:: filename %s too short!", filnam)
	}

	bf := []byte(filnam)
	// replace empty space with underscore
	for i:= 0; i< len(filnam); i++ {
		if bf[i] == ' ' {bf[i]='_'}
		if bf[i] == '.' {
			return fmt.Errorf("error createImgFolder:: filnam has period!")
		}
	}

	imgFoldNam := "output/img_" + string(bf)

	if _, err := os.Stat(imgFoldNam); os.IsNotExist(err) {
		err1:= os.Mkdir(imgFoldNam, os.ModePerm)
		if err1 != nil {
			return fmt.Errorf("error createImgFolder:: could not create folder! %v", err1)
		}
	}
	return nil
}
