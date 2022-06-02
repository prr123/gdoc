// readYaml.go
// author: prr
//
package main

import (
	"os"
	"fmt"
	util "google/gdoc/gdocUtil"
)

func main() {

	inpFilPath := "input"
	filnam := "option.yaml"

	inpOpt, err := util.ReadYamlFil(inpFilPath, filnam)
	if err != nil {
		fmt.Printf("error readYamlFile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", inpOpt)
	util.PrintOptYaml(inpOpt)

	fmt.Println("success")
}
