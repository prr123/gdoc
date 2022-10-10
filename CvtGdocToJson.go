//
// CvtGdocToJson
// adapted from GdocToDom
// date: 10 Ocober 2022
//
// author: prr
// copyright 2022 prr azul software
//
//		gdocHtml "google/gdoc/gdocHtml"
//

package main

import (
        "fmt"
        "os"
//		"strings"
		gdocApi "google/gdoc/gdocApi"
		gdocJson "google/gdoc/gdocJson"
)


func main() {
	var gd gdocApi.GdocApiStruct
	// initialise default values
	baseFolder := "output/"
//	baseFolderSlash := baseFolder + "/"
	sel:=""
	dbg := false
    numArgs := len(os.Args)
//	fmt.Printf("args: %d\n", numArgs)
	cmd := os.Args[0]
	outFolder := "json"

	switch numArgs {
		case 1:
       		fmt.Println("error - no docid provided!")
          	fmt.Printf("%s usage is:\n  %s docId [outfolder] [dbg] [summary, main, all]\n", cmd[2:], cmd)
        	os.Exit(1)
		case 2:
       		fmt.Println("no outFolder nor selection argument provided!")
			fmt.Println("assuming outFolder is 'default' selection is 'all'!")
			sel = "all"
		case 3:
			switch os.Args[2] {
				case "dbg":
					dbg = true
					sel = "all"
				case "summary", "main", "all":
					sel = os.Args[2]
				default:
					outFolder = os.Args[2]
			}
		case 4:
			outFolder = os.Args[2]

			switch os.Args[3] {
				case "dbg":
					dbg = true
					sel = "all"
				case "summary", "main", "all":
					sel = os.Args[3]
				default:
					fmt.Println("invalid argument: %s!\n", os.Args[2])
          			fmt.Printf("%s usage is:\n  %s docId [out] [dbg] [summary, main, all] [opt]\n", cmd[2:], cmd)
					os.Exit(1)
			}
		case 5:
			outFolder = os.Args[2]
			if os.Args[3] == "dbg" {
				dbg = true
			} else {
				fmt.Println("invalid argument: %s!\n", os.Args[3])
       			fmt.Printf("%s usage is:\n  %s docId [out] [dbg] [summary, main, all] [opt]\n", cmd[2:], cmd)
				os.Exit(1)
			}
			switch os.Args[4] {
				case "summary", "main", "all":
					sel = os.Args[4]
				default:
					fmt.Println("invalid argument: %s!\n", os.Args[4])
          			fmt.Printf("%s usage is:\n  %s docId [out] [dbg] [summary, main, all] [opt]\n", cmd[2:], cmd)
					os.Exit(1)
			}
		default:
        	fmt.Println("error - too many arguments!")
          	fmt.Printf("%s usage is:\n  %s docId [dbg] [summary, main, all] [opt]\n", cmd[2:], cmd)
        	os.Exit(1)
    }

    docId := os.Args[1]

	err := gd.InitGdocApi()
	if err != nil {
		fmt.Printf("error - InitGdocApi: %v!", err)
		os.Exit(1)
	}

	srv := gd.Svc

	outfilPath:= baseFolder + outFolder

	fmt.Printf("*************** CvtGdocToJSON ************\n")
	if (dbg) {fmt.Printf("*** debug ***\n")}
	fmt.Printf("output folder: %s selection: %s\n", outfilPath, sel)

	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from document: ", err)
		os.Exit(1)
	}
	fmt.Printf("Doc Title: %s Selection: %s\n", doc.Title, sel)

	switch sel {
	case "heading":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
//		err = gdocHtml.CreGdocHtmlSection("", outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error main: CvtGdocJsonSummary -- cannot convert gdoc doc: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success summary ***!")
		os.Exit(0)

	case "main":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
		if err != nil {
			fmt.Println("error main CvtGdocJsonMain -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}
		fmt.Println("*** success main ***!")
		os.Exit(0)

	case "doc":
		fmt.Printf("*** not implemented yet ***\n")
		os.Exit(1)
		if err != nil {
			fmt.Println("error CvtGdocJsonDoc -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success doc ***!")
		os.Exit(0)

	case "all":
		err = gdocJson.CreGdocJsonAll(outfilPath, doc, nil)
		if err != nil {
			fmt.Println("error CreGdocJsonAll -- cannot convert gdoc file: ", err)
			os.Exit(1)
		}

		fmt.Println("*** success all ***!")
		os.Exit(0)

	default:
		fmt.Printf("%s is not a valid comand line opt!\n", sel)
		fmt.Println("exiting!")
		os.Exit(1)
	}

	fmt.Println("Success!")
}
