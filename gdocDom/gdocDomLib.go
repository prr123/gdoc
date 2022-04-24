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
    "net/http"
    "io"
    "unicode/utf8"
    "google.golang.org/api/docs/v1"
	gd "gdocHtml/gdocHtmlLib"
)

type GdocDomObj struct {
	doc *doc.Document
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

}


func CreGdocDomAll(folderPath string, doc *docs.Document, options *gd.OptObj)(err error) {

// function that creates an html fil from the named section
    var tocDiv *gd.dispObj
    var dObj gd.GdocHtmlObj
	var dom GdocDomObj

    err = dom.initGdocDom(doc, options)
    if err != nil {
        return fmt.Errorf("initGdocHtml %v", err)
    }

    fPath, fexist, err := gd.CreateDocFolder(folderPath, dObj.docName)
    if err!= nil {
        return fmt.Errorf("createHtmlFolder %v", err)
    }
    dObj.folderPath = fPath

    if dObj.Options.Verb {
        fmt.Println("******************* Output File ************")
        fmt.Printf("folder path: %s ", fPath)
        fstr := "is new!"
        if fexist { fstr = "exists!" }
        fmt.Printf("%s\n", fstr)
    }

    err = dObj.createOutFil("")
    if err!= nil {
        return fmt.Errorf("createOutFil %v", err)
    }

    if dObj.Options.ImgFold {
        err = dObj.dlImages()
        if err != nil {
            fmt.Errorf("dlImages: %v", err)
        }
    }
/*
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

    docHeadStr,_ := creHtmlHead()
    outfil.WriteString(docHeadStr)

    //css
    outfil.WriteString(headObj.bodyCss)

    //css of named styles
    outfil.WriteString(mainDiv.headCss)
    outfil.WriteString(mainDiv.bodyCss)

    //css footnotes
    if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyCss)}

    //css toc
    if toc {
        outfil.WriteString(tocDiv.bodyCss)
    }
*/
    outfil.WriteString("</style>\n</head>\n<body>\n")

    // html
    outfil.WriteString(headObj.bodyHtml)
    // html toc
    if toc {outfil.WriteString(tocDiv.bodyHtml)}

    if dObj.Options.Sections {outfil.WriteString(secDiv.bodyHtml)}

    // html main document
    outfil.WriteString(mainDiv.bodyHtml)

    // html footnotes
    if ftnoteDiv != nil {outfil.WriteString(ftnoteDiv.bodyHtml)}
 
   outfil.WriteString("</body>\n</html>\n")
    outfil.Close()
    return nil
}
