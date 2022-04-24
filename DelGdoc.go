// 22/1/2022
// DelGdoc
// author prr
// from ReadGdoc.go
//

package main

import (
        "fmt"
        "os"
		gdocApi "google/gdoc/gdocApi2"
//		gdocEdit "google/gdoc/gdocEdit"
)


func main() {
	var gd gdocApi.GdocApi2Obj
	var docId, input string

    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
        fmt.Println("DelGdoc usage is: \"DelGdoc docId\"")
		docId = "1y6GbEKMhg_9WmRoAjSHGGUS1OXZfLqttB16MVaKHkXU"
		fmt.Printf("defaulting to docId: %s\n\n", docId)
    } else {
    	docId = os.Args[1]
	}

	fmt.Println("Fetching google doc with id: ", docId)

	err := gd.Init()
	if err != nil {
		fmt.Printf("error Init: %v \n", err)
		fmt.Println("exiting due to error!")
		os.Exit(1)
	}

	err = gd.Gdoc_read(docId)
	if err != nil {
		fmt.Println("error Gdoc_read: ", err)
		os.Exit(1)
	}
	fmt.Println("success reading Google Doc!")
	fmt.Println("doc title: ",gd.Doc.Title)

	fmt.Printf("deleting doc with id \"%s\"\n", docId)
	fmt.Print("Are you sure (Y/n):")
 	fmt.Scanln(&input)
	fmt.Println("inp: ", input)

	err = gd.Gdoc_delete(docId)
	if err != nil {
		fmt.Println("error Gdoc_delete: ", err)
		os.Exit(1)
	}

	fmt.Println("success!")
	os.Exit(0)

/*
	srv := gd.Svc
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
0	fmt.Printf("The title of the doc is: %s\n", doc.Title)


	outfil, err := gd.CreTxtOutFile(doc.Title, "toml")
	if err != nil {
		fmt.Println("error main -- cannot open out file: ", err)
		os.Exit(1)
	}
	err = gdocHtml.CvtGdocTxt(outfil, doc)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("success!")
	os.Exit(0)
*/
}

