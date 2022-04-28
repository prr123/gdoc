package main

import (
	"os"
	"fmt"
	util "google/gdoc/gdocUtil"
)

func main() {

	folder := "gdocUtil"
	filnam := "testopt.yaml"
	folderPath, size, err := util.CheckFil(folder, filnam)
	if err != nil {
		fmt.Printf("error Checkfil: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("folderpath: %s size: %d\n", folderPath, size)
	opt, err := util.ReadYamlFil(folder,filnam)
	if err != nil {
		fmt.Printf("read yaml error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Base Font:\n")
	fmt.Printf("  Font Name: %s\n", opt.BaseFont.Name)
	fmt.Printf("  Font Size: %s\n", opt.BaseFont.Size)
	fmt.Printf("Doc:\n")
	fmt.Printf("  Doc Id:    %s\n", opt.Doc.DocId)


	fmt.Println("success!")
}
