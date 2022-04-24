// GdocTpl
// author prr
// 13/1/2022
// from EditGdoc
//  program that applies template file to convert template fields
//
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi2"
		gdocEdit "google/gdoc/gdocEdit"
)


func main() {
	var gd gdocApi.GdocApi2Obj

    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
          fmt.Println("CvtGdocToHtmlv1 usage is:\n  CvtGdocToTxtv1 docId\n")
        os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.Init()
	srv := gd.Svc
//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	infilnam := "output/TemplateTest.tpl"
	err = gdocEdit.ReadGdocTpl(infilnam, doc)
	if err != nil {
		fmt.Println("Unable to substitute template fields: ", err)
		os.Exit(1)
	}

	fmt.Println("success!")
	os.Exit(0)
}

