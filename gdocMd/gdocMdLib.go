// gdocMdLib.go
// author: prr
// created 10/2021
// copyright 2022 prr
//
//

package gdocToMd

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"io"
	"google.golang.org/api/docs/v1"
)

const (
    PtTomm = 0.35277777777778
    MmTopt = 2.8346456692913207
)


type gdocMdObj struct {
    parCount int
    posImgCount int
    inImgCount int
    imgId []string
    doc *docs.Document
    DocName string
	listmap  map[string][]int
	listid	string
	nestlev int
	folder *os.File
	imgFoldNam string
	imgFoldPath string
}

func (dObj *gdocMdObj) downloadImg()(err error) {

	verb := false
	doc := dObj.doc
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
		if verb	{fmt.Printf("image path: %s\n", imgNam)}
		URL := imgProp.ContentUri
		httpResp, err := http.Get(URL)
    	if err != nil {
        	return fmt.Errorf("error downloadImg:: could not fetch %s! %v", URL, err)
    	}
    	defer httpResp.Body.Close()
//	fmt.Printf("http got %s!\n", URL)
    	if httpResp.StatusCode != 200 {
        	return fmt.Errorf("error downloadImg:: Received non 200 response code %d!", httpResp.StatusCode)
    	}
//	fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file
    	outfil, err := os.Create(imgNam)
    	if err != nil {
        	return fmt.Errorf("error downloadImg:: cannot create img file! %v", err)
    	}
    	defer outfil.Close()
//	fmt.Println("created dir")
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
		if verb	{fmt.Printf("image path: %s\n", imgNam)}
		URL := imgProp.ContentUri
		httpResp, err := http.Get(URL)
    	if err != nil {
        	return fmt.Errorf("error downloadImg:: could not fetch %s! %v", URL, err)
    	}
    	defer httpResp.Body.Close()
//	fmt.Printf("http got %s!\n", URL)
    	if httpResp.StatusCode != 200 {
        	return fmt.Errorf("error downloadImg:: Received non 200 response code %d!", httpResp.StatusCode)
    	}
//	fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file
    	outfil, err := os.Create(imgNam)
    	if err != nil {
        	return fmt.Errorf("error downloadImg:: cannot create img file! %v", err)
    	}
    	defer outfil.Close()
//	fmt.Println("created dir")
    	//Write the bytes to the fiel
    	_, err = io.Copy(outfil, httpResp.Body)
    	if err != nil {
        	return fmt.Errorf("error downloadImg:: cannot copy img file content! %v", err)
    	}
	}

    return nil
}


func (dObj *gdocMdObj) createImgFolder()(err error) {

	filnam :=dObj.DocName
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
//		return fmt.Errorf("error createImgFolder:: no sub folder!")
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


func (dObj *gdocMdObj) Init() (err error) {
    if dObj == nil {
        return fmt.Errorf("error gdocMD::Init: dObj is nil!")
    }
	doc := dObj.doc
	dObj.DocName = doc.Title
	dObj.inImgCount = len(doc.InlineObjects)
	dObj.posImgCount = len(doc.PositionedObjects)
	totObjNum := dObj.inImgCount + dObj.posImgCount
    dObj.parCount = len(doc.Body.Content)

	if totObjNum == 0 {return nil}

	err = dObj.createImgFolder()
	if err != nil {
		return fmt.Errorf("error gdocMd::Init: could create ImgFolder: %v!", err)
	}
	err = dObj.downloadImg()
	if err != nil {
		return fmt.Errorf("error gdocMd::Init: could download images: %v!", err)
	}

    return nil
}


func (dObj *gdocMdObj) getColor(color  *docs.Color)(outstr string) {
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

func (dObj *gdocMdObj) cvtTocName(parstr string)(tocstr string) {
	var yLen int
	var y [100]byte
	parlStr := strings.ToLower(parstr)
	parlStr = strings.TrimSpace(parlStr)
	if len(parstr) < 1 {
		return ""
	}
	//fmt.Println("toc parstr: ", parstr," : ", len(parstr))
	//fmt.Println("toc parlStr: ", parlStr, " : ", len(parlStr))

	x:= []byte(parlStr)

	last := false
	yLen =0
	for i:=0; i<len(x); i++ {
		switch x[i] {
		case ' ':
			if !last {
				y[yLen] = '-'
				yLen++
			}
			last = true
		case '\n':

		default:
			y[yLen] = x[i]
			yLen++
			last = false
		}
	}

	tocstr = string(y[:yLen])

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

func (dObj *gdocMdObj) cvtPelInlineImg(imgEl *docs.InlineObjectElement)(outstr string, err error) {

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
    outstr=fmt.Sprintf("<img src=\"%s\" id=\"%s\" alt=\"%s\" width=\"%.1fpt\" height=\"%.1fpt\">\n",
		imgSrc, imgId, imgObj.Title, imgObj.Size.Width.Magnitude, imgObj.Size.Height.Magnitude )

    return outstr, nil
}

func (dObj *gdocMdObj) cvtPelTxt(parEl *docs.ParagraphElement)(outstr string, err error) {
	var txtStr string
    if parEl == nil {
        return "", fmt.Errorf("error cvtPelToMd -- parEl is nil!")
    }


	cLen := len(parEl.TextRun.Content)
	x := []byte(parEl.TextRun.Content)

// should be either <br> or nothing
	if (cLen == 1) && (x[0] == '\n') {
		return "", nil
	}

	if x[cLen-1] == '\n' {
		txtStr = string(x[:cLen-1])
//		retStr = "    \n"
		cLen--
	} else {
		txtStr = string(x)
	}

//fmt.Println("txtstr: ", txtStr)
    prefix := ""
    suffix := ""

    if parEl.TextRun.TextStyle.Italic {
        prefix = "_"
        suffix  ="_"
    }
    if parEl.TextRun.TextStyle.Bold {
        prefix = "**" + prefix
        suffix = suffix + "**"
    }


    outstr = prefix + txtStr + suffix
    return outstr, nil
}

func (dObj *gdocMdObj) cvtParToMd(par *docs.Paragraph)(outstr string, tocstr string, err error) {
	var prefix, suffix, tocPrefix, tocSuffix string
	var parStr, liststr string
	var NamedTxtStyl *docs.TextStyle
	var ListProp *docs.NestingLevel
	var doc *docs.Document
	
	doc = dObj.doc

    if par == nil {
        return "","", fmt.Errorf("error dispParStyl: no par pt")
    }

//	fmt.Printf("  Paragraph with %d Par-Elements\n", len(par.Elements))
	outstr = ""
	if par.Bullet != nil {
		liststr = ""
		ListId := par.Bullet.ListId
		NestLev:= par.Bullet.NestingLevel
//		outstr += fmt.Sprintf("[//]: * (List: %s Nest: %d )\n", ListId, NestLev)
		// need to determine whether ordered or unordered
		ListProp = nil
		for key, list := range doc.Lists {
			if key == ListId {
				ListProp = list.ListProperties.NestingLevels[NestLev]
				break
			}
		}
		if ListProp == nil {
			outstr += "[//]: * (Error List " + ListId + " cannot find listid )\n"
		} else {
			liststr = ""
			if dObj.listid != ListId {
				// new list
				dObj.listmap[ListId] = make([]int,5)
				dObj.listmap[ListId][NestLev] = 1
				dObj.listid = ListId
				dObj.nestlev = int(NestLev)
			} else {
				// same list
				if dObj.nestlev == int(NestLev) {
					dObj.listmap[ListId][NestLev]++
				} else {
//	fmt.Printf("ListId: %s Nest: %d prev nest: %d val: %d\n", ListId, NestLev, dObj.nestlev, dObj.listmap[ListId][NestLev])
					if dObj.nestlev < int(NestLev) {
					// new sublist
						dObj.listmap[ListId][NestLev] = 1
						dObj.nestlev = int(NestLev)
					} else {
						dObj.nestlev = int(NestLev)
						dObj.listmap[ListId][NestLev]++
					}
				}
			}

			var i int64
			for i=0; i<NestLev+1; i++ {
				liststr += "   "
			}
			if len(ListProp.GlyphSymbol)>0 {
				liststr += "* "
			} else {
				switch ListProp.GlyphType {
					case "DECIMAL":
						liststr += fmt.Sprintf("%d. ",dObj.listmap[ListId][NestLev])
					default:
						liststr += fmt.Sprintf("%d. ",dObj.listmap[ListId][NestLev])
				}
			}
		}
	}

	parStylTyp := par.ParagraphStyle.NamedStyleType
	// doc styles
	NamedStylIdx := -1
	for i:=0; i < len(doc.NamedStyles.Styles); i++ {
		nstyl := doc.NamedStyles.Styles[i].NamedStyleType;
		if nstyl  == parStylTyp {
			NamedStylIdx = i
//			NamedParStyl = doc.NamedStyles.Styles[i].ParagraphStyle
			NamedTxtStyl = doc.NamedStyles.Styles[i].TextStyle
		}
	}

	if NamedStylIdx == -1 {
		return "","",fmt.Errorf("error cvtParToMd: par style not a named style!")
	}
//
//	fmt.Println("par style type: ", parStylTyp, " index: ", NamedStylIdx);

    prefix = ""
	suffix = ""
	tocPrefix = ""
	tocSuffix = ""
	decode := false
	titlestyl := false
//	subtitlestyl := false

	// need to reconsider for paragraphs without text
	boldStyl := false
	italicStyl := false
	if par.Elements[0].TextRun != nil {
		parTxtStyl := par.Elements[0].TextRun.TextStyle
		boldStyl = parTxtStyl.Bold || NamedTxtStyl.Bold
		italicStyl = parTxtStyl.Italic || NamedTxtStyl.Italic
	}

    switch parStylTyp {
		case "TITLE":
			prefix ="<p style=\"font-size:20pt; text-align:center;"
			suffix = "</p>\n\n"
			if boldStyl {prefix += " font-weight:bold;"}
			if italicStyl {prefix += " font-style:italic;"}
			titlestyl = true
			prefix += "\">"
			tocPrefix = prefix
			tocSuffix = suffix
      	case "SUBTITLE":
			prefix ="<p style=\"font-size:16pt;text-align:center;"
			suffix = "</p>\n\n"
			if boldStyl {prefix += " font-weight:bold;"}
			if italicStyl {prefix += " font-style:italic;"}
			prefix += "\">"
//			subtitlestyl = true

		case "HEADING_1":
            prefix = "# "
			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n# ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
		case "HEADING_2":
            prefix = "## "
			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n## ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
        case "HEADING_3":
			prefix = "### "
			suffix = fmt.Sprintf("     \n")
			tocPrefix = "\n### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
        case "HEADING_4":
            prefix = "#### "
			suffix = fmt.Sprintf("     \n")
			tocPrefix = "\n#### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
        case "HEADING_5":
            prefix = "##### "
			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n##### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
        case "HEADING_6":
            prefix = "###### "
			suffix = fmt.Sprintf("    \n")
			tocPrefix = "\n###### ["
			tocSuffix = fmt.Sprintf("](#")
			decode = true
        case "NORMAL_TEXT":
            prefix = ""
			suffix = "\n\n"
        default:
            prefix = fmt.Sprintf("[//]: * (Name Style: %s unknown)\n", par.ParagraphStyle.NamedStyleType)
    }

	numParEl := len(par.Elements)
	if (numParEl == 1) && (!(len(par.Elements[0].TextRun.Content) >0)) {
//		return "<br>\n","", nil
		return "\n","", nil
	}

	parStr = ""
	for p:=0; p< numParEl; p++ {

		parEl := par.Elements[p]
//fmt.Println("parEl: ", p, parEl.TextRun.Content)
//		outstr += fmt.Sprintf("\nPar-El[%d]: %d - %d \n", p, parEl.StartIndex, parEl.EndIndex)

		// text
		if parEl.TextRun != nil {
			tstr, err := dObj.cvtPelTxt(parEl)
			if err != nil {
				outstr = fmt.Sprintf("\n[//]: # (error cvtPelTxt: %v)\n",err)
			}
			parStr += tstr + outstr
		//fmt.Printf("parstr %d: %s (%d)\n",p, parStr,len(parStr))
		}

		if parEl.HorizontalRule != nil {

		}
		if parEl.ColumnBreak != nil {

		}
   		if parEl.Person != nil {

		}

		// inline image
		if parEl.InlineObjectElement != nil {
            tstr, err := dObj.cvtPelInlineImg(parEl.InlineObjectElement)
            if err != nil {
				outstr = fmt.Sprintf("\n[//]: # (error cvtPelInlineImg: %v)\n",err)
            }
			parStr += tstr + outstr
		}
		//fmt.Printf("Img parstr %d: %s (%d)\n",p, parStr,len(parStr))

 		if parEl.RichLink != nil {

		}
	} // loop parEl


//	parLen := len(parStr)
//	xb := []byte(parStr)
	nparStr := parStr
//fmt.Println("parstr: ",parStr," : ",len(parStr))

// case of string terminated by new line
/*
	if xb[parLen-1] == '\n' {
		nparStr = string(xb[:parLen-1])
	}
*/
//fmt.Println("nparstr: ",nparStr," : ",len(nparStr))

	if !(len(nparStr) > 0) {
		outstr ="\n"
		tocstr = "\n"
		return outstr, tocstr,nil
	}

// check of new line in the middle of the string
		xnb := []byte(nparStr)
		n2parStr := ""
		ist := 0
		for i:= 0; i< len(nparStr); i++ {
			if xnb[i] == '\n' {
				n2parStr += string(xnb[ist:i]) + "     \n"
				ist = i
			}
		}
		n2parStr += string(xnb[ist:])
//fmt.Println("n2parstr: ",n2parStr," : ",len(n2parStr))
		if decode {
			tocParStr := dObj.cvtTocName(n2parStr)
//    	outstr += prefix + parStr + suffix
			tocstr+= tocPrefix + n2parStr + tocSuffix + tocParStr + ")\n\n"
		}
		boldPrefix := "";
		itPrefix := "";
		if titlestyl {
			tocstr+= tocPrefix + n2parStr + tocSuffix + "\n"
	    	outstr += prefix + n2parStr + suffix
		} else {
//			if subtitlestyl {
//	    		outstr += prefix + n2parStr + suffix
//			} else {
				if italicStyl {itPrefix = "_"}
				if boldStyl {boldPrefix = "**"}
	    		outstr += "\n" + liststr + prefix + boldPrefix + itPrefix + n2parStr + itPrefix + boldPrefix + suffix
//			}
		}

	if par.PositionedObjectIds != nil {
		outstr += fmt.Sprintf("/nParagraph has as Positioned Objects: %d\n", len(par.PositionedObjectIds))
		for id:=0; id< len(par.PositionedObjectIds); id++ {
			outstr += fmt.Sprintf("posObject Id[%d]: %s\n", id, par.PositionedObjectIds[id])
		}
	}

	return outstr, tocstr,nil
}


func CvtGdocToMd(outfil *os.File, doc *docs.Document, toc bool)(err error) {
	var outstr, tocstr, liststr string
    docObj := new(gdocMdObj)
    docObj.doc = doc
	docObj.folder = outfil
    err = docObj.Init()
    if err != nil {
        return fmt.Errorf("error CvtGdocToMd: could not initialise! %v", err)
    }
	outstr = "[//]: * (Document Title: " + doc.Title + ")\n"

	outstr += fmt.Sprintf("[//]: * (Document Id: %s)\n", doc.DocumentId)
//	outstr += fmt.Sprintf("Revision Id: %s \n", doc.RevisionId)

	_, err = outfil.WriteString(outstr)
	if err != nil {
		return fmt.Errorf("error CvtGdocTpMd: cannot write to file! %v", err)
	}

	outstr = ""
	tocstr = "<p style=\"font-size:20pt; text-align:center\">Table of Contents</p>\n"

	body := doc.Body
	elNum := len(body.Content)
//	outstr = "\n******************** Body *********************************\n"

	for el:=0; el< elNum; el++ {
		bodyEl := body.Content[el]
//		outstr += fmt.Sprintf("\nelement: %d StartIndex: %d EndIndex: %d\n", el, bodyEl.StartIndex, bodyEl.EndIndex)
		if bodyEl.Paragraph != nil {
			par := bodyEl.Paragraph
//			outstr += fmt.Sprintf("  Paragraph with %d Par-Elements\n", len(par.Elements))
			if par.Bullet != nil {
			}
			tstr, toctstr, err :=docObj.cvtParToMd(par)
			if err != nil {
				outstr += fmt.Sprintf("error cvtPar: %v\n",err)
			} else {
				outstr += liststr + tstr
				tocstr += toctstr
			}
			if par.ParagraphStyle != nil {
//				outstr += "  Has Style Structure\n"
			}
			if par.PositionedObjectIds != nil {
				outstr += fmt.Sprintf("  Has Positioned Objects: %d\n", len(par.PositionedObjectIds))
			}

		}
		if bodyEl.SectionBreak != nil {
//			outstr += fmt.Sprintf("Section Break\n")
		}
		if bodyEl.Table != nil {
//			outstr += fmt.Sprintf("Table\n")
		}
		if bodyEl.TableOfContents != nil {
//			outstr += fmt.Sprintf("Table of Contents\n")
		}

	}

	if toc {
		tocstr += "\n\n"
		outfil.WriteString(tocstr)
	}

	outfil.WriteString(outstr)

	outfil.Close()
	return nil
}

