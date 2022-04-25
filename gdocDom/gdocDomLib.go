// golang library that creates a html file from a gdoc file
// author: prr
// created: 22/04/2021
// copyright 2022 prr, Peter Riemenschneider
//
// for changes see github
//
// start: CreGdocHtmlTil
//

package gdocDom

import (
    "fmt"
    "os"
//    "net/http"
//    "io"
//    "unicode/utf8"
    "google.golang.org/api/docs/v1"
	gd "google/gdoc/gdocHtml"
	util "github.com/prr123/util"
)

type GdocDomObj struct {
	doc *docs.Document
	docName string
	folderPath string
	outfil *os.File
}

type jsDomObj struct {
	jsStr string
	cssStr string
}

func creHtmlHead(docName string)(outstr string) {
    outstr = "<!DOCTYPE html>\n"
  	outstr += fmt.Sprintf("<!-- file: %s -->\n", docName)
    outstr += "<head>\n<script>\n"
	outstr += "windows.onload"
    outstr += "</script>\n</head>\n<body>\n"
	outstr += fmt.Sprintf("<div class=\"%s_doc\">\n",docName)
	outstr += "</div>/n</body></html>/n"
    return outstr
}

func (domObj *GdocDomObj) initGdocDom (doc *docs.Document, option *gd.OptObj) {

	domObj.doc = doc

   // need to transform file name
    // replace spaces with underscore
    dNam := doc.Title
    x := []byte(dNam)
    for i:=0; i<len(x); i++ {
        if x[i] == ' ' {
            x[i] = '_'
        }
    }
    domObj.docName = string(x[:])
	return
}


func CreGdocDomAll(folderPath string, doc *docs.Document, options *gd.OptObj)(err error) {

// function that creates an html fil from the named section
    var tocDiv *jsDomObj
    var dObj gd.GdocHtmlObj
	var domObj GdocDomObj

    domObj.initGdocDom(doc, options)
    if err != nil {
        return fmt.Errorf("initGdocDom %v", err)
    }

    fPath, fexist, err := util.CreateFileFolder(folderPath, domObj.docName)
    if err!= nil {
        return fmt.Errorf("util.CreateFileFolder %v", err)
    }
    domObj.folderPath = fPath

    if dObj.Options.Verb {
        fmt.Println("******************* Output File ************")
        fmt.Printf("folder path: %s ", fPath)
        fstr := "is new!"
        if fexist { fstr = "exists!" }
        fmt.Printf("%s\n", fstr)
    }

    outfil, err := util.CreateOutFil(fPath, domObj.docName,"html")
    if err!= nil {
        return fmt.Errorf("util.CreateOutFil %v", err)
    }
	domObj.outfil = outfil
/*
    if dObj.Options.ImgFold {
        err = dObj.dlImages()
        if err != nil {
            fmt.Errorf("dlImages: %v", err)
        }
    }

// footnotes
    ftnoteDiv, err := dObj.createFootnoteDiv()
    if err != nil {
        fmt.Errorf("createFootnoteDiv: %v", err)
    }

//  dObj.sections
    var mainDiv, secDiv dispObj

    if len(dObj.sections) > 0 {
        secDiv = dObj.createSectionHead()

        for ipage:=0; ipage<len(dObj.sections); ipage++ {
//          pageStr := fmt.Sprintf("Pg_%d", ipage)
//          idStr := fmt.Sprintf("%s_pg_%d", dObj.docName, ipage)
//ppp
            pgHd := dObj.createSectionDiv(ipage)
            elStart := dObj.sections[ipage].secElStart
            elEnd := dObj.sections[ipage].secElEnd
            pgBody, err := dObj.cvtBodySec(elStart, elEnd)
            if err != nil {
                return fmt.Errorf("cvtBodySec %d %v", ipage, err)
            }
            mainDiv.headCss += pgBody.headCss
            mainDiv.bodyCss += pgBody.bodyCss
            mainDiv.bodyHtml += pgHd.bodyHtml + pgBody.bodyHtml
        }
    }

    if len(dObj.sections) == 0 {
        mBody, err := dObj.cvtBody()
        if err != nil {
            return fmt.Errorf("cvtBody: %v", err)
        }
        mainDiv.headCss += mBody.headCss
        mainDiv.bodyCss += mBody.bodyCss
        mainDiv.bodyHtml += mBody.bodyHtml
    }

    headObj, err := dObj.createHead()
    if err != nil {
        return fmt.Errorf("createHead: %v", err)
    }

    toc := dObj.Options.Toc
    if toc {
        tocDiv, err = dObj.createTocDiv()
        if err != nil {
            tocDiv.bodyHtml = fmt.Sprintf("<!--- error Toc Head: %v --->\n",err)
        }
    }

   // create html file
    outfil := dObj.htmlFil
    if outfil == nil {
        return fmt.Errorf("outfil is nil!")
    }
*/
    docHeadStr := creHtmlHead()
    outfil.WriteString(docHeadStr)

    //css
//    outfil.WriteString(headObj.bodyCss)

    //css of named styles
//    outfil.WriteString(mainDiv.headCss)
//    outfil.WriteString(mainDiv.bodyCss)

    //css footnotes
//    if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyCss)}

    //css toc
//    if toc {
//        outfil.WriteString(tocDiv.bodyCss)
//    }

    outfil.WriteString("</style>\n<script>\n")

    outfil.WriteString("/// begin of script\n")

    outfil.WriteString("</script>\n</head>\n<body>\n")

    // html
//    outfil.WriteString(headObj.bodyHtml)
    // html toc
//    if toc {outfil.WriteString(tocDiv.bodyHtml)}

//    if dObj.Options.Sections {outfil.WriteString(secDiv.bodyHtml)}

    // html main document
//    outfil.WriteString(mainDiv.bodyHtml)

    // html footnotes
//    if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyHtml)}

	outfil.WriteString("</body>\n</html>\n")
    outfil.Close()
    return nil
}
